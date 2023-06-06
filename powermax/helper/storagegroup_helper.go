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
	"fmt"
	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/models"
)

// constants to annotate if a volume should be added or removed.
const (
	AddVolume    = 1
	RemoveVolume = 2
	NoOperation  = 0
)

// AddRemoveVolume add or remove a volume based on the config of plan and current state.
func AddRemoveVolume(ctx context.Context, plan, state *models.StorageGroupResourceModel, client *client.Client) error {
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
		err := client.PmaxClient.AddVolumesToStorageGroupS(ctx, client.SymmetrixID, plan.StorageGroupID.ValueString(), false, addVolumeArr...)
		if err != nil {
			return err
		}
	}
	if len(removeVolumeArr) > 0 {
		_, err := client.PmaxClient.RemoveVolumesFromStorageGroup(ctx, client.SymmetrixID, plan.StorageGroupID.ValueString(), false, removeVolumeArr...)
		if err != nil {
			return err
		}
	}
	state.VolumeIDs = plan.VolumeIDs
	return nil
}

// UpdateSgState update the state of storage group based on the current state of the storage group.
func UpdateSgState(ctx context.Context, client *client.Client, sgID string, state *models.StorageGroupResourceModel) error {
	// Update all fields of state
	storageGroup, err := client.PmaxClient.GetStorageGroup(ctx, client.SymmetrixID, sgID)
	if err != nil {
		return err
	}

	err = CopyFields(ctx, storageGroup, state)
	if err != nil {
		return err
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
				"host_io_limit_io_sec": types.StringValue(storageGroup.HostIOLimit.HostIOLimitIOSec),
				"host_io_limit_mb_sec": types.StringValue(storageGroup.HostIOLimit.HostIOLimitMBSec),
				"dynamic_distribution": types.StringValue(storageGroup.HostIOLimit.DynamicDistribution),
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
	volumeIDListInStorageGroup, err := client.PmaxClient.GetVolumeIDListInStorageGroup(ctx, client.SymmetrixID, state.StorageGroupID.ValueString())
	if err != nil {
		return err
	}
	state.VolumeIDs, _ = types.ListValueFrom(ctx, types.StringType, volumeIDListInStorageGroup)

	// set ID
	state.ID = types.StringValue(storageGroup.StorageGroupID)

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
