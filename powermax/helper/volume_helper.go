/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helper

import (
	"context"
	"dell/powermax-go-client"
	"fmt"
	"math/big"
	"reflect"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

// UpdateVolResourceState updates resource state given vol response from array.
func UpdateVolResourceState(ctx context.Context, volState *models.VolumeResource, volResponse *powermax.Volume, volPlan *models.VolumeResource) error {
	// Manually copy
	volState.ID = types.StringValue(volResponse.VolumeId)
	if volPlan != nil {
		volState.CapUnit = volPlan.CapUnit
		volState.StorageGroupName = volPlan.StorageGroupName
		volState.VolumeIdentifier = volPlan.VolumeIdentifier
	}
	// Copy values with the same fields
	err := CopyFields(ctx, volResponse, volState)
	if err != nil {
		return err
	}
	tflog.Info(ctx, fmt.Sprintf("Capacity %v", volResponse))
	// Convert size
	switch volState.CapUnit.ValueString() {
	case CapacityUnitCyl:
		volState.Size = types.NumberValue(big.NewFloat(float64(*volResponse.CapCyl)))
	case CapacityUnitTb:
		volState.Size = types.NumberValue(big.NewFloat(*volResponse.CapGb / 1024))
	case CapacityUnitGb:
		volState.Size = types.NumberValue(big.NewFloat(*volResponse.CapGb))
	case CapacityUnitMb:
		// pmax returns 1 MB less than actual cap
		volState.Size = types.NumberValue(big.NewFloat(*volResponse.CapMb - 1.0))
	}
	volState.MobilityIDEnabled = types.BoolValue(*volResponse.MobilityIdEnabled)
	// Handle symmetrix port key Storage Groups and RDF Group
	volState.SymmetrixPortKey, _ = GetSymmetrixPortKeyObjects(volResponse)
	volState.StorageGroups, _ = GetStorageGroupObjects(volResponse)
	volState.RDFGroupIDList, _ = GetRfdGroupIdsObjects(volResponse)

	return nil
}

// GetSymmetrixPortKeyObjects returns symmetrix port key objects.
func GetSymmetrixPortKeyObjects(volResponse *powermax.Volume) (types.List, diag.Diagnostics) {
	// handle symmetrix port key due to name rule
	var symmetrixPortKeyObjects []attr.Value
	pkType := map[string]attr.Type{
		"director_id": types.StringType,
		"port_id":     types.StringType,
	}
	for _, pk := range volResponse.SymmetrixPortKey {
		pkMap := make(map[string]attr.Value)
		pkMap["director_id"] = types.StringValue(pk.DirectorId)
		pkMap["port_id"] = types.StringValue(pk.PortId)
		pkObject, _ := types.ObjectValue(pkType, pkMap)
		symmetrixPortKeyObjects = append(symmetrixPortKeyObjects, pkObject)
	}
	return types.ListValue(types.ObjectType{AttrTypes: pkType}, symmetrixPortKeyObjects)
}

// GetRfdGroupIdsObjects Rfd group key objects.
func GetRfdGroupIdsObjects(volResponse *powermax.Volume) (types.List, diag.Diagnostics) {
	// handle symmetrix port key due to name rule
	var rfdKeyObjects []attr.Value
	pkType := map[string]attr.Type{
		"rdf_group_number": types.Int64Type,
		"label":            types.StringType,
	}
	for _, pk := range volResponse.GetRdfGroupId() {
		pkMap := make(map[string]attr.Value)
		pkMap["rdf_group_number"] = types.Int64Value(pk.GetRdfGroupNumber())
		pkMap["label"] = types.StringValue(pk.GetLabel())
		pkObject, _ := types.ObjectValue(pkType, pkMap)
		rfdKeyObjects = append(rfdKeyObjects, pkObject)
	}
	return types.ListValue(types.ObjectType{AttrTypes: pkType}, rfdKeyObjects)
}

// GetStorageGroupObjects returns storage group key objects.
func GetStorageGroupObjects(volResponse *powermax.Volume) (types.List, diag.Diagnostics) {
	// handle symmetrix port key due to name rule
	var sgKeyObjects []attr.Value
	pkType := map[string]attr.Type{
		"storage_group_name":        types.StringType,
		"parent_storage_group_name": types.StringType,
	}
	for _, pk := range volResponse.StorageGroups {
		pkMap := make(map[string]attr.Value)
		pkMap["storage_group_name"] = types.StringValue(pk.GetStorageGroupName())
		pkMap["parent_storage_group_name"] = types.StringValue(pk.GetParentStorageGroupName())
		pkObject, _ := types.ObjectValue(pkType, pkMap)
		sgKeyObjects = append(sgKeyObjects, pkObject)
	}
	return types.ListValue(types.ObjectType{AttrTypes: pkType}, sgKeyObjects)
}

// UpdateVol updates the volume and return updated parameters, failed updated parameters and errors.
func UpdateVol(ctx context.Context, client *client.Client, planVol, stateVol models.VolumeResource) ([]string, []string, []string) {
	var updatedParameters []string
	var updateFailedParameters []string
	var errorMessages []string

	if planVol.VolumeIdentifier.ValueString() != stateVol.VolumeIdentifier.ValueString() {
		modifyParam := client.PmaxOpenapiClient.SLOProvisioningApi.ModifyVolume(ctx, client.SymmetrixID, stateVol.ID.ValueString())
		modifyParam = modifyParam.EditVolumeParam(powermax.EditVolumeParam{
			EditVolumeActionParam: &powermax.EditVolumeActionParam{
				ModifyVolumeIdentifierParam: &powermax.ModifyVolumeIdentifierParam{
					VolumeIdentifier: &powermax.VolumeIdentifier{
						VolumeIdentifierChoice: "identifier_name",
						IdentifierName:         planVol.VolumeIdentifier.ValueStringPointer(),
					},
				},
			},
		})
		_, _, err := modifyParam.Execute()
		if err != nil {
			errStr := ""
			message := GetErrorString(err, errStr)
			updateFailedParameters = append(updateFailedParameters, "name")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to rename volume: %s", message))
		} else {
			updatedParameters = append(updatedParameters, "name")
		}
	}

	if planVol.MobilityIDEnabled.ValueBool() != stateVol.MobilityIDEnabled.ValueBool() {
		modifyParam := client.PmaxOpenapiClient.SLOProvisioningApi.ModifyVolume(ctx, client.SymmetrixID, stateVol.ID.ValueString())
		modifyParam = modifyParam.EditVolumeParam(powermax.EditVolumeParam{
			EditVolumeActionParam: &powermax.EditVolumeActionParam{
				EnableMobilityIdParam: &powermax.EnableMobilityIdParam{
					EnableMobilityId: planVol.MobilityIDEnabled.ValueBool(),
				},
			},
		})
		_, _, err := modifyParam.Execute()
		if err != nil {
			errStr := ""
			message := GetErrorString(err, errStr)
			updateFailedParameters = append(updateFailedParameters, "enable_mobility_id")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify mobility: %s", message))
		} else {
			updatedParameters = append(updatedParameters, "enable_mobility_id")
		}
	}
	if planVol.Size.ValueBigFloat().Cmp(stateVol.Size.ValueBigFloat()) != 0 || planVol.CapUnit.ValueString() != stateVol.CapUnit.ValueString() {
		modifyParam := client.PmaxOpenapiClient.SLOProvisioningApi.ModifyVolume(ctx, client.SymmetrixID, stateVol.ID.ValueString())
		modifyParam = modifyParam.EditVolumeParam(powermax.EditVolumeParam{
			EditVolumeActionParam: &powermax.EditVolumeActionParam{
				ExpandVolumeParam: &powermax.ExpandVolumeParam{
					VolumeAttribute: powermax.VolumeAttribute{
						CapacityUnit: planVol.CapUnit.ValueString(),
						VolumeSize:   planVol.Size.String(),
					},
				},
			},
		})
		_, _, err := modifyParam.Execute()
		if err != nil {
			errStr := ""
			message := GetErrorString(err, errStr)
			updateFailedParameters = append(updateFailedParameters, "size")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify the volume size: %s", message))
		} else {
			updatedParameters = append(updatedParameters, "size")
		}
	}

	return updatedParameters, updateFailedParameters, errorMessages
}

// GetVolumeFilterParam returns volume filter parameters.
func GetVolumeFilterParam(ctx context.Context, p *client.Client, model models.VolumeDatasource) (powermax.ApiListVolumesRequest, error) {
	filter := model.VolumeFilter
	param := p.PmaxOpenapiClient.SLOProvisioningApi.ListVolumes(ctx, p.SymmetrixID)

	v := reflect.ValueOf(filter).Elem()
	for i := 0; i < v.NumField(); i++ {
		key := v.Type().Field(i).Name
		switch key {
		case "StorageGroupID":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.StorageGroupId(stringValue.ValueString())
			}
		case "EncapsulatedWwn":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.EncapsulatedWwn(val)
			}
		case "WWN":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.Wwn(val)
			}
		case "Symmlun":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.Symmlun(val)
			}
		case "Status":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.Status(val)
			}
		case "PhysicalName":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.PhysicalName(val)
			}
		case "VolumeIdentifier":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.VolumeIdentifier(val)
			}
		case "AllocatedPercent":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.AllocatedPercent(val)
			}
		case "CapTb":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.CapTb(val)
			}
		case "CapGb":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.CapGb(val)
			}
		case "CapMb":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.CapMb(val)
			}
		case "CapCYL":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.CapCyl(val)
			}
		case "NumOfStorageGroups":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.NumOfStorageGroups(val)
			}
		case "NumOfMaskingViews":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.NumOfMaskingViews(val)
			}
		case "NumOfFrontEndPaths":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.NumOfFrontEndPaths(val)
			}
		case "VirtualVolumes":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.VirtualVolumes(stringValue.String())
			}
		case "PrivateVolumes":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.PrivateVolumes(stringValue.String())
			}
		case "AvailableThinVolumes":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.AvailableThinVolumes(stringValue.String())
			}
		case "Tdev":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.Tdev(stringValue.String())
			}
		case "ThinBcv":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.ThinBcv(stringValue.String())
			}
		case "Vdev":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.Vdev(stringValue.String())
			}
		case "Gatekeeper":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.Gatekeeper(stringValue.String())
			}
		case "DataVolume":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.DataVolume(stringValue.String())
			}
		case "Dld":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.Dld(stringValue.String())
			}
		case "Drv":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.Drv(stringValue.String())
			}
		case "Mapped":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.Mapped(stringValue.String())
			}
		case "BoundTdev":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.BoundTdev(stringValue.String())
			}
		case "Reserved":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.Reserved(stringValue.String())
			}
		case "Pinned":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.Pinned(stringValue.String())
			}
		case "Encapsulated":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.Encapsulated(stringValue.String())
			}
		case "Associated":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.Associated(stringValue.String())
			}
		case "Emulation":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.Emulation(val)
			}
		case "SplitName":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.SplitName(val)
			}
		case "CuImageNum":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.CuImageNum(val)
			}
		case "CuImageSsid":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.CuImageSsid(val)
			}
		case "RdfGroupNumber":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.RdfGroupNumber(val)
			}
		case "HasEffectiveWwn":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.HasEffectiveWwn(stringValue.String())
			}
		case "EffectiveWwn":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.EffectiveWwn(val)
			}
		case "Type":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.Type_(val)
			}
		case "OracleInstanceName":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.OracleInstanceName(val)
			}
		case "MobilityIDEnabled":
			stringValue, ok := v.Field(i).Interface().(types.Bool)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				param = param.MobilityIdEnabled(stringValue.String())
			}
		case "UnreducibleDataGb":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.UnreducibleDataGb(val)
			}
		case "Nguid":
			stringValue, ok := v.Field(i).Interface().(types.String)
			if !ok {
				return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
			}
			if !stringValue.IsNull() {
				val := make([]string, 0)
				val = append(val, stringValue.ValueString())
				param = param.Nguid(val)
			}
		}
	}

	return param, nil
}
