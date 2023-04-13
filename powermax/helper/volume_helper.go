// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.

package helper

import (
	"context"
	"fmt"
	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"math/big"
	"reflect"
	"strconv"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/models"
)

const (
	// CapacityUnitTb represents the unit TeraBytes for capacity.
	CapacityUnitTb = "TB"
	// CapacityUnitGb represents the unit GigaBytes for capacity.
	CapacityUnitGb = "GB"
	// CapacityUnitMb represents the unit MegaBytes for capacity.
	CapacityUnitMb = "MB"
	// CapacityUnitCyl represents the unit CYL for capacity.
	CapacityUnitCyl = "CYL"
)

// GetVolumeSize returns the size of the volume in the correct format
func GetVolumeSize(volume models.VolumeResource) (interface{}, error) {
	var size interface{}
	if volume.CapUnit.ValueString() == CapacityUnitCyl {
		if volume.Size.ValueBigFloat().IsInt() {
			intVal, err := strconv.Atoi(volume.Size.ValueBigFloat().String())
			if err != nil {
				return size, err
			}
			size = intVal
		} else {
			return size, fmt.Errorf("when cap_unit is 'CYL', size should be defined only as an integer")
		}
	} else {
		size = volume.Size.ValueBigFloat().String()
	}
	return size, nil
}

// UpdateVolResourceState updates resource state given vol response from array.
func UpdateVolResourceState(ctx context.Context, volState *models.VolumeResource, volResponse *pmaxTypes.Volume, volPlan *models.VolumeResource) error {
	// Manually copy
	volState.ID = types.StringValue(volResponse.VolumeID)
	if volPlan != nil {
		volState.CapUnit = volPlan.CapUnit
		volState.StorageGroupName = volPlan.StorageGroupName
	}
	// Copy values with the same fields
	err := CopyFields(ctx, volResponse, volState)
	if err != nil {
		return err
	}
	// Convert size
	switch volState.CapUnit.ValueString() {
	case CapacityUnitCyl:
		volState.Size = types.NumberValue(big.NewFloat(float64(volResponse.CapacityCYL)))
	case CapacityUnitTb:
		volState.Size = types.NumberValue(big.NewFloat(volResponse.CapacityGB / 1024))
	case CapacityUnitGb:
		volState.Size = types.NumberValue(big.NewFloat(volResponse.CapacityGB))
	case CapacityUnitMb:
		volState.Size = types.NumberValue(big.NewFloat(volResponse.FloatCapacityMB))
	}
	// Handle symmetrix port key
	volState.SymmetrixPortKey, _ = GetSymmetrixPortKeyObjects(volResponse)

	return nil
}

func GetSymmetrixPortKeyObjects(volResponse *pmaxTypes.Volume) (types.List, diag.Diagnostics) {
	// handle symmetrix port key due to name rule
	var symmetrixPortKeyObjects []attr.Value
	pkType := map[string]attr.Type{
		"director_id": types.StringType,
		"port_id":     types.StringType,
	}
	for _, pk := range volResponse.SymmetrixPortKey {
		pkMap := make(map[string]attr.Value)
		pkMap["director_id"] = types.StringValue(pk.DirectorID)
		pkMap["port_id"] = types.StringValue(pk.PortID)
		pkObject, _ := types.ObjectValue(pkType, pkMap)
		symmetrixPortKeyObjects = append(symmetrixPortKeyObjects, pkObject)
	}
	return types.ListValue(types.ObjectType{AttrTypes: pkType}, symmetrixPortKeyObjects)
}

// UpdateVol updates the volume and return updated parameters, failed updated parameters and errors
func UpdateVol(ctx context.Context, client *client.Client, planVol, stateVol models.VolumeResource) ([]string, []string, []string) {
	var updatedParameters []string
	var updateFailedParameters []string
	var errorMessages []string

	if planVol.VolumeIdentifier.ValueString() != stateVol.VolumeIdentifier.ValueString() {
		_, err := client.PmaxClient.RenameVolume(ctx, client.SymmetrixID, stateVol.ID.ValueString(), planVol.VolumeIdentifier.ValueString())
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "name")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to rename volume: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "name")
		}
	}

	if planVol.MobilityIDEnabled.ValueBool() != stateVol.MobilityIDEnabled.ValueBool() {
		_, err := client.PmaxClient.ModifyMobilityForVolume(ctx, client.SymmetrixID, stateVol.ID.ValueString(), planVol.MobilityIDEnabled.ValueBool())
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "enable_mobility_id")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify mobility: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "enable_mobility_id")
		}
	}

	if planVol.Size.ValueBigFloat() != stateVol.Size.ValueBigFloat() || planVol.CapUnit.ValueString() != stateVol.CapUnit.ValueString() {
		size, err := GetVolumeSize(planVol)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "size")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify the volume size: %s", err.Error()))
			return updatedParameters, updateFailedParameters, errorMessages
		}
		_, err = client.PmaxClient.ExpandVolume(ctx, client.SymmetrixID, stateVol.ID.ValueString(), 0, size, planVol.CapUnit.ValueString())
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "size")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify the volume size: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "size")
		}
	}

	return updatedParameters, updateFailedParameters, errorMessages
}

func GetVolumeFilterParam(model models.VolumeDatasource) (map[string]string, error) {
	filter := model.VolumeFilter
	param := make(map[string]string)
	v := reflect.ValueOf(filter).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Type() {
		case reflect.TypeOf(basetypes.StringValue{}):
			stringValue, ok := field.Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", field.Type().Field(i).Name)
			}
			if len(stringValue.ValueString()) != 0 {
				param[v.Type().Field(i).Tag.Get("tfsdk")] = stringValue.ValueString()
			}
		case reflect.TypeOf(basetypes.BoolValue{}):
			boolValue, ok := field.Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", field.Type().Field(i).Name)
			}
			if !boolValue.IsNull() {
				param[v.Type().Field(i).Tag.Get("tfsdk")] = boolValue.String()
			}
		default:
			return param, fmt.Errorf("unexpected field %s is detected", field.Type().Field(i).Name)
		}
	}
	// Due to the rule of attribute name, need to handle storageGroupId separately
	if _, ok := param["storage_group_name"]; ok {
		param["storageGroupId"] = param["storage_group_name"]
		delete(param, "storage_group_name")
	}

	return param, nil
}
