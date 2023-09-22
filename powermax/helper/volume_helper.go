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
	"net/http"
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

// GetVolume on SG.
func GetVolume(ctx context.Context, client client.Client, volID string) (*powermax.Volume, *http.Response, error) {
	return client.PmaxOpenapiClient.SLOProvisioningApi.GetVolume(ctx, client.SymmetrixID, volID).Execute()
}

// ListVolumes on SG.
func ListVolumes(ctx context.Context, client client.Client, plan models.VolumeResource) (*powermax.Iterator, *http.Response, error) {
	param := client.PmaxOpenapiClient.SLOProvisioningApi.ListVolumes(ctx, client.SymmetrixID)
	param = param.StorageGroupId(plan.StorageGroupName.ValueString())
	return param.Execute()
}

// CreateVolume on SG.
func CreateVolume(ctx context.Context, client client.Client, plan models.VolumeResource) (*powermax.StorageGroup, *http.Response, error) {
	volumeAttributes := make([]powermax.VolumeAttribute, 0)
	num := int64(1)
	volumeAttributes = append(volumeAttributes, powermax.VolumeAttribute{
		CapacityUnit: plan.CapUnit.ValueString(),
		NumOfVols:    &num,
		VolumeSize:   plan.Size.ValueBigFloat().String(),
		VolumeIdentifier: &powermax.VolumeIdentifier{
			VolumeIdentifierChoice: "identifier_name",
			IdentifierName:         plan.VolumeIdentifier.ValueStringPointer(),
		},
	})
	tflog.Info(ctx, fmt.Sprintf("Create Volume att Param: %v", volumeAttributes))
	createNewVol := true
	emulation := "FBA"
	tflog.Debug(ctx, "calling create volume in storage groups on pmax client", map[string]interface{}{
		"symmetrixID":      client.SymmetrixID,
		"storageGroupName": plan.StorageGroupName.ValueString(),
		"name":             plan.VolumeIdentifier.ValueString(),
		"volumeAttributes": volumeAttributes,
	})
	createParam := client.PmaxOpenapiClient.SLOProvisioningApi.ModifyStorageGroup(ctx, client.SymmetrixID, plan.StorageGroupName.ValueString())
	createParam = createParam.EditStorageGroupParam(
		powermax.EditStorageGroupParam{
			EditStorageGroupActionParam: powermax.EditStorageGroupActionParam{
				ExpandStorageGroupParam: &powermax.ExpandStorageGroupParam{
					AddVolumeParam: &powermax.AddVolumeParam{
						CreateNewVolumes: &createNewVol,
						EnableMobilityId: plan.MobilityIDEnabled.ValueBoolPointer(),
						VolumeAttributes: volumeAttributes,
						Emulation:        &emulation,
						VolumeIdentifier: &powermax.VolumeIdentifier{
							VolumeIdentifierChoice: "identifier_name",
							IdentifierName:         plan.VolumeIdentifier.ValueStringPointer(),
						},
					},
				},
			},
		},
	)
	return createParam.Execute()
}

// UpdateVolumeState iterates over the volume list and update the state.
func UpdateVolumeState(ctx context.Context, p *client.Client, params powermax.ApiListVolumesRequest) (response []models.VolumeDatasourceEntity, err error) {
	volIDs, _, err := params.Execute()
	if err != nil {
		errStr := ""
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf(message)
	}
	for _, vol := range volIDs.ResultList.GetResult() {
		for _, volumeID := range vol {
			tflog.Info(ctx, fmt.Sprint(volumeID))
			volumeModel := p.PmaxOpenapiClient.SLOProvisioningApi.GetVolume(ctx, p.SymmetrixID, fmt.Sprint(volumeID))
			volResponse, _, err := volumeModel.Execute()
			if err != nil {
				errStr := ""
				message := GetErrorString(err, errStr)
				return nil, fmt.Errorf(message)

			}
			volState := models.VolumeDatasourceEntity{}
			err = CopyFields(ctx, volResponse, &volState)
			volState.SymmetrixPortKey, _ = GetSymmetrixPortKeyObjects(volResponse)
			volState.StorageGroups, _ = GetStorageGroupObjects(volResponse)
			volState.RfdGroupIDList, _ = GetRfdGroupIdsObjects(volResponse)
			if id, ok := volResponse.GetVolumeIdOk(); ok {
				volState.VolumeID = types.StringValue(*id)
			}
			if mobid, ok := volResponse.GetMobilityIdEnabledOk(); ok {
				volState.MobilityIDEnabled = types.BoolValue(*mobid)
			}
			if err != nil {
				return nil, err
			}
			volState.VolumeID = types.StringValue(volResponse.VolumeId)
			volState.MobilityIDEnabled = types.BoolValue(*volResponse.MobilityIdEnabled)
			response = append(response, volState)
		}
	}
	return response, nil
}

// GetVolumeFilterParam returns volume filter parameters.
func GetVolumeFilterParam(ctx context.Context, p *client.Client, model models.VolumeDatasource) (powermax.ApiListVolumesRequest, error) {
	filter := model.VolumeFilter
	param := p.PmaxOpenapiClient.SLOProvisioningApi.ListVolumes(ctx, p.SymmetrixID)

	if filter != nil {

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
					param = param.EncapsulatedWwn(stringValue.ValueString())
				}
			case "WWN":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.Wwn(stringValue.ValueString())
				}
			case "Symmlun":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.Symmlun(stringValue.ValueString())
				}
			case "Status":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.Status(stringValue.ValueString())
				}
			case "PhysicalName":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.PhysicalName(stringValue.ValueString())
				}
			case "VolumeIdentifier":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.VolumeIdentifier(stringValue.ValueString())
				}
			case "AllocatedPercent":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.AllocatedPercent(stringValue.ValueString())
				}
			case "CapTb":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.CapTb(stringValue.ValueString())
				}
			case "CapGb":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.CapGb(stringValue.ValueString())
				}
			case "CapMb":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.CapMb(stringValue.ValueString())
				}
			case "CapCYL":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.CapCyl(stringValue.ValueString())
				}
			case "NumOfStorageGroups":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.NumOfStorageGroups(stringValue.ValueString())
				}
			case "NumOfMaskingViews":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.NumOfMaskingViews(stringValue.ValueString())
				}
			case "NumOfFrontEndPaths":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.NumOfFrontEndPaths(stringValue.ValueString())
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
					param = param.Emulation(stringValue.ValueString())
				}
			case "SplitName":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.SplitName(stringValue.ValueString())
				}
			case "CuImageNum":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.CuImageNum(stringValue.ValueString())
				}
			case "Ssid":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.CuImageSsid(stringValue.ValueString())
				}
			case "RdfGroupNumber":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.RdfGroupNumber(stringValue.ValueString())
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
					param = param.EffectiveWwn(stringValue.ValueString())
				}
			case "Type":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.Type_(stringValue.ValueString())
				}
			case "OracleInstanceName":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.OracleInstanceName(stringValue.ValueString())
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
					param = param.UnreducibleDataGb(stringValue.ValueString())
				}
			case "Nguid":
				stringValue, ok := v.Field(i).Interface().(types.String)
				if !ok {
					return param, fmt.Errorf("failed to type assertion on field %s", v.Type().Field(i).Name)
				}
				if !stringValue.IsNull() {
					param = param.Nguid(stringValue.ValueString())
				}
			}
		}

	}
	tflog.Info(ctx, fmt.Sprintf("Param!!!! %v", param))
	return param, nil

}
