package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
	"strconv"
)

const (
	AddVolume    = 1
	RemoveVolume = 2
	NoOperation  = 0
)

func UpdateVolume(ctx context.Context, plan, state *StorageGroupResourceModel, r *StorageGroup) error {
	err := addRemoveVolume(ctx, plan, state, r)
	if err != nil {
		return err
	}

	err = createVolume(ctx, plan, state, r)
	if err != nil {
		return err
	}
	return nil
}

func addRemoveVolume(ctx context.Context, plan, state *StorageGroupResourceModel, r *StorageGroup) error {
	var planVolumeIDs []string
	var stateVolumeIDs []string

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

	if reflect.DeepEqual(planVolumeIDs, stateVolumeIDs) {
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
		err := r.client.PmaxClient.AddVolumesToStorageGroupS(ctx, r.client.SymmetrixID, plan.StorageGroupID.ValueString(), false, addVolumeArr...)
		if err != nil {
			return err
		}
	}
	if len(removeVolumeArr) > 0 {
		_, err := r.client.PmaxClient.RemoveVolumesFromStorageGroup(ctx, r.client.SymmetrixID, plan.StorageGroupID.ValueString(), false, removeVolumeArr...)
		if err != nil {
			return err
		}
	}
	state.VolumeIDs = plan.VolumeIDs
	num := state.NumOfVolumes.ValueInt64() + int64(len(addVolumeArr)-len(removeVolumeArr))
	state.NumOfVolumes = types.Int64Value(num)
	return nil
}

func createVolume(ctx context.Context, plan, state *StorageGroupResourceModel, r *StorageGroup) error {
	// Create new volumes to the storage group based on the attribute "num_of_vols"
	planNumOfVolumes := plan.NumOfVolumes.ValueInt64()
	stateNumOfVolumes := state.NumOfVolumes.ValueInt64()
	if planNumOfVolumes == stateNumOfVolumes {
		return nil
	}

	var stateVolumeIDs []string
	if !state.VolumeIDs.IsNull() && !state.VolumeIDs.IsUnknown() {
		diags := state.VolumeIDs.ElementsAs(ctx, &stateVolumeIDs, true)
		if diags.HasError() {
			return fmt.Errorf("unable to parse volume ids from state")
		}
	}
	for planNumOfVolumes > stateNumOfVolumes {
		volumeOptions := make(map[string]interface{})
		if !plan.CapacityUnit.IsNull() && !plan.CapacityUnit.IsUnknown() {
			volumeOptions["capacityUnit"] = plan.CapacityUnit.ValueString()
		}
		volumeSizeStr := plan.VolumeSize.ValueString()
		volumeSizeInt, err := strconv.Atoi(volumeSizeStr)
		if err != nil {
			return err
		}
		addedVolume, err := r.client.PmaxClient.CreateVolumeInStorageGroupS(ctx, r.client.SymmetrixID, plan.StorageGroupID.ValueString(), plan.VolumeIdentifierName.ValueString(), volumeSizeInt, volumeOptions)
		if err != nil {
			return err
		}
		stateNumOfVolumes++
		stateVolumeIDs = append(stateVolumeIDs, addedVolume.VolumeID)
	}
	state.NumOfVolumes = plan.NumOfVolumes
	volumeSet, diagnostics := types.SetValueFrom(ctx, types.StringType, stateVolumeIDs)
	if diagnostics.HasError() {
		return fmt.Errorf("failed to create volume and parse volume ids: %s", stateVolumeIDs)
	}
	state.VolumeIDs = volumeSet
	return nil
}
