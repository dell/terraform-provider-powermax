// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package helper

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/models"
)

const (
	AddVolume    = 1
	RemoveVolume = 2
	NoOperation  = 0
)

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
