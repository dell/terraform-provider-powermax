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
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/models"

	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// constants to annotate if a volume should be added or removed.
const (
	AddVolume    = 1
	RemoveVolume = 2
	NoOperation  = 0
)

// AddRemoveVolume add or remove a volume based on the config of plan and current state.
func AddRemoveVolume(ctx context.Context, plan *models.StorageGroupResourceModel, state *models.StorageGroupResourceModel, client *client.Client, sgID string) error {

	var planVolumeIDs []string
	var stateVolumeIDs []string

	// Initialize field VolumeIDs
	if state.VolumeIDs.IsNull() || state.VolumeIDs.IsUnknown() {
		state.VolumeIDs, _ = types.ListValueFrom(ctx, types.StringType, []string{})
	}

	if plan.VolumeIDs.IsNull() || plan.VolumeIDs.IsUnknown() {
		return nil
	}

	if !plan.VolumeIDs.IsNull() && !plan.VolumeIDs.IsUnknown() {
		diags := plan.VolumeIDs.ElementsAs(ctx, &planVolumeIDs, true)
		if diags.HasError() {
			return fmt.Errorf("unable to parse volume ids from plan")
		}
	}
	if !state.VolumeIDs.IsNull() && !state.VolumeIDs.IsUnknown() {
		diags := state.VolumeIDs.ElementsAs(ctx, &stateVolumeIDs, true)
		if diags.HasError() {
			return fmt.Errorf("unable to parse volume ids from state")
		}
	}

	if CompareStringSlice(planVolumeIDs, stateVolumeIDs) {
		return nil
	}

	payload := client.PmaxOpenapiClient.SLOProvisioningApi.ModifyStorageGroup(ctx, client.SymmetrixID, sgID)
	// Add or remove existing volumes to the storage group based on the attribute "volume_ids"
	volumeIDMap := make(map[string]int)
	for _, elem := range planVolumeIDs {
		volumeIDMap[elem] = AddVolume
	}
	for _, elem := range stateVolumeIDs {
		if _, found := volumeIDMap[elem]; found {
			volumeIDMap[elem] = NoOperation
		} else {
			volumeIDMap[elem] = RemoveVolume
		}
	}
	var addVolumeArr []string
	var removeVolumeArr []string
	for val, oper := range volumeIDMap {
		if oper == AddVolume {
			addVolumeArr = append(addVolumeArr, val)
		} else if oper == RemoveVolume {
			removeVolumeArr = append(removeVolumeArr, val)
		}
	}
	if len(addVolumeArr) > 0 {
		payload = payload.EditStorageGroupParam(powermax.EditStorageGroupParam{
			EditStorageGroupActionParam: powermax.EditStorageGroupActionParam{
				ExpandStorageGroupParam: &powermax.ExpandStorageGroupParam{
					AddSpecificVolumeParam: &powermax.AddSpecificVolumeParam{
						VolumeId: addVolumeArr,
					},
				},
			},
		})
		_, _, err := payload.Execute()
		if err != nil {
			return err
		}
	}
	if len(removeVolumeArr) > 0 {
		payload = payload.EditStorageGroupParam(powermax.EditStorageGroupParam{
			EditStorageGroupActionParam: powermax.EditStorageGroupActionParam{
				RemoveVolumeParam: &powermax.RemoveVolumeParam{
					VolumeId: removeVolumeArr,
				},
			},
		})
		_, _, err := payload.Execute()
		if err != nil {
			return err
		}
	}
	state.VolumeIDs = plan.VolumeIDs
	return nil
}

// CreateSloParam Create SLO param.
func CreateSloParam(plan models.StorageGroupResourceModel) []powermax.SloBasedStorageGroupParam {

	hostIOLimit := ConstructHostIOLimit(plan)
	workload := "None"
	thickVolumes := false
	num := int64(0)

	if hostIOLimit != nil {
		return []powermax.SloBasedStorageGroupParam{
			{
				SloId:                      plan.Slo.ValueStringPointer(),
				WorkloadSelection:          &workload,
				AllocateCapacityForEachVol: &thickVolumes,
				NoCompression:              &thickVolumes,
				VolumeAttributes: []powermax.VolumeAttribute{
					{
						VolumeSize:   "0",
						CapacityUnit: "CYL",
						NumOfVols:    &num,
					},
				},
				SetHostIOLimitsParam: &powermax.SetHostIOLimitsParam{
					HostIoLimitMbSec:    &hostIOLimit.HostIOLimitMBSec,
					HostIoLimitIoSec:    &hostIOLimit.HostIOLimitIOSec,
					DynamicDistribution: &hostIOLimit.DynamicDistribution,
				},
			},
		}
	}
	return []powermax.SloBasedStorageGroupParam{
		{
			SloId:                      plan.Slo.ValueStringPointer(),
			WorkloadSelection:          &workload,
			AllocateCapacityForEachVol: &thickVolumes,
			NoCompression:              &thickVolumes,
			VolumeAttributes: []powermax.VolumeAttribute{
				{
					VolumeSize:   "0",
					CapacityUnit: "CYL",
					NumOfVols:    &num,
				},
			},
		},
	}

}

// UpdateSgState update the state of storage group based on the current state of the storage group.
func UpdateSgState(ctx context.Context, client *client.Client, sgID string, state *models.StorageGroupResourceModel) error {
	// Update all fields of state
	sgModel := client.PmaxOpenapiClient.SLOProvisioningApi.GetStorageGroup2(ctx, client.SymmetrixID, sgID)
	storageGroup, _, err := sgModel.Execute()
	if err != nil {
		return err
	}

	err = CopyFields(ctx, storageGroup, state)
	if err != nil {
		return err
	}
	if id, ok := storageGroup.GetStorageGroupIdOk(); ok {
		state.StorageGroupID = types.StringValue(*id)
	}

	if uuid, ok := storageGroup.GetUuidOk(); ok {
		state.UUID = types.StringValue(*uuid)
	}

	// set HostIOLimit
	if storageGroup.HostIOLimit != nil {
		state.HostIOLimit, _ = types.ObjectValue(
			map[string]attr.Type{
				"host_io_limit_io_sec": types.StringType,
				"host_io_limit_mb_sec": types.StringType,
				"dynamic_distribution": types.StringType,
			},
			map[string]attr.Value{
				"host_io_limit_io_sec": types.StringValue(*storageGroup.HostIOLimit.HostIoLimitIoSec),
				"host_io_limit_mb_sec": types.StringValue(*storageGroup.HostIOLimit.HostIoLimitMbSec),
				"dynamic_distribution": types.StringValue(*storageGroup.HostIOLimit.DynamicDistribution),
			})
	} else {
		state.HostIOLimit, _ = types.ObjectValue(
			map[string]attr.Type{
				"host_io_limit_io_sec": types.StringType,
				"host_io_limit_mb_sec": types.StringType,
				"dynamic_distribution": types.StringType,
			},
			map[string]attr.Value{
				"host_io_limit_io_sec": types.StringNull(),
				"host_io_limit_mb_sec": types.StringNull(),
				"dynamic_distribution": types.StringNull(),
			})
	}

	// Read volume list in storage group
	volIDModel := client.PmaxOpenapiClient.SLOProvisioningApi.ListVolumes(ctx, client.SymmetrixID)
	// Set the storage group id
	volIDModel = volIDModel.StorageGroupId(storageGroup.StorageGroupId)
	volumeIDListInStorageGroup, _, err := volIDModel.Execute()
	vol := make([]string, 0, len(volumeIDListInStorageGroup.GetResultList().Result))
	for _, v := range volumeIDListInStorageGroup.ResultList.Result {
		for _, v2 := range v {
			vol = append(vol, fmt.Sprint(v2))
		}
	}
	if err != nil {
		return err
	}

	state.VolumeIDs, _ = types.ListValueFrom(ctx, types.StringType, vol)
	// set ID
	state.ID = types.StringValue(storageGroup.StorageGroupId)

	return nil
}

// ConstructHostIOLimit constructs the host io limit param based on the plan.
func ConstructHostIOLimit(plan models.StorageGroupResourceModel) *pmaxTypes.SetHostIOLimitsParam {
	if !plan.HostIOLimit.IsNull() && !plan.HostIOLimit.IsUnknown() {
		hostIOLimit := models.SetHostIOLimitsParam{}
		tfsdk.ValueAs(context.Background(), plan.HostIOLimit, &hostIOLimit)
		hostIOLimitParam := &pmaxTypes.SetHostIOLimitsParam{
			HostIOLimitIOSec:    hostIOLimit.HostIOLimitIOSec.ValueString(),
			HostIOLimitMBSec:    hostIOLimit.HostIOLimitMBSec.ValueString(),
			DynamicDistribution: hostIOLimit.DynamicDistribution.ValueString(),
		}
		return hostIOLimitParam
	}
	return nil
}
