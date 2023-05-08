// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.

package helper

import (
	"context"
	"fmt"
	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math/big"
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
func GetVolumeSize(volume models.Volume) (interface{}, error) {
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

// UpdateVolState updates the volume state with the response from the array
func UpdateVolState(ctx context.Context, volState *models.Volume, volResponse *pmaxTypes.Volume, volPlan *models.Volume) error {
	// Manually copy
	volState.ID = types.StringValue(volResponse.VolumeID)
	if volPlan != nil {
		volState.CapUnit = volPlan.CapUnit
		volState.StorageGroupName = volPlan.StorageGroupName
	}
	// Copy values with the same fields
	err := CopyFields(ctx, volResponse, volState)

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

	// copy storage groups
	var sgObjects []attr.Value
	sgType := map[string]attr.Type{
		"storage_group_name": types.StringType,
	}
	for _, sg := range volResponse.StorageGroupIDList {
		sgMap := make(map[string]attr.Value)
		sgMap["storage_group_name"] = types.StringValue(sg)
		sgObject, _ := types.ObjectValue(sgType, sgMap)
		sgObjects = append(sgObjects, sgObject)
	}
	volState.StorageGroups, _ = types.ListValue(types.ObjectType{AttrTypes: sgType}, sgObjects)

	var symmetrixPortKeysObjects []attr.Value
	pkType := map[string]attr.Type{
		"director_id": types.StringType,
		"port_id":     types.NumberType,
	}
	for _, pk := range volResponse.SymmetrixPortKey {
		pkMap := make(map[string]attr.Value)
		pkMap["director_id"] = types.StringValue(pk.DirectorID)
		pkMap["port_id"] = types.StringValue(pk.PortID)
		pkObject, _ := types.ObjectValue(pkType, pkMap)
		symmetrixPortKeysObjects = append(symmetrixPortKeysObjects, pkObject)
	}
	volState.SymmetrixPortKeys, _ = types.ListValue(types.ObjectType{AttrTypes: pkType}, symmetrixPortKeysObjects)

	var rdfGroupsObjects []attr.Value
	rdfType := map[string]attr.Type{
		"rdf_group_number": types.Int64Type,
		"label":            types.StringType,
	}
	for _, rdf := range volResponse.RDFGroupIDList {
		rdfMap := make(map[string]attr.Value)
		rdfMap["rdf_group_number"] = types.Int64Value(int64(rdf.RDFGroupNumber))
		rdfMap["label"] = types.StringValue(rdf.Label)
		rdfGroupsObject, _ := types.ObjectValue(rdfType, rdfMap)
		rdfGroupsObjects = append(rdfGroupsObjects, rdfGroupsObject)
	}
	volState.RDFGroupIDList, _ = types.ListValue(types.ObjectType{AttrTypes: rdfType}, rdfGroupsObjects)
	return err
}

// UpdateVol updates the volume and return updated parameters, failed updated parameters and errors
func UpdateVol(ctx context.Context, client *client.Client, planVol, stateVol models.Volume) ([]string, []string, []string) {
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
