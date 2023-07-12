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
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// UpdateSnapshotDatasourceState Update Snaposhot state.
func UpdateSnapshotDatasourceState(ctx context.Context, snapshotDetail *powermax.SnapVXSnapshotInstance, state *models.SnapshotDetailModal) error {
	// Copy values with the same fields
	err := CopyFields(ctx, snapshotDetail, state)
	state.LinkedStorageGroup, _ = GetLinkedSgList(snapshotDetail)
	state.SourceVolume, _ = GetSnapshotGenerationVolume(snapshotDetail)
	tflog.Debug(ctx, fmt.Sprintf("Snapshot Detail State: %v", state))
	if err != nil {
		return err
	}
	return nil
}

// GetSnapshotGenerationVolume Get snapshot generation volume.
func GetSnapshotGenerationVolume(snapshotDetail *powermax.SnapVXSnapshotInstance) (types.List, diag.Diagnostics) {
	var genObjects []attr.Value
	typeKey := map[string]attr.Type{
		"name":        types.StringType,
		"capacity":    types.Int64Type,
		"capacity_gb": types.Float64Type,
	}
	for _, gen := range snapshotDetail.SourceVolume {
		genMap := make(map[string]attr.Value)
		genMap["name"] = types.StringValue(gen.Name)
		genMap["capacity"] = types.Int64Value(gen.Capacity)
		genMap["capacity_gb"] = types.Float64Value(float64(gen.CapacityGb))

		genObject, _ := types.ObjectValue(typeKey, genMap)
		genObjects = append(genObjects, genObject)
	}
	return types.ListValue(types.ObjectType{AttrTypes: typeKey}, genObjects)
}

// GetLinkedSgList Get linked storage group list.
func GetLinkedSgList(snapshotDetail *powermax.SnapVXSnapshotInstance) (types.List, diag.Diagnostics) {
	var sgObjects []attr.Value
	typeKey := map[string]attr.Type{
		"name":                          types.StringType,
		"source_volume_name":            types.StringType,
		"linked_volume_name":            types.StringType,
		"tracks":                        types.Int64Type,
		"track_size":                    types.Int64Type,
		"percentage_copied":             types.Int64Type,
		"linked_creation_timestamp":     types.StringType,
		"defined":                       types.BoolType,
		"background_define_in_progress": types.BoolType,
	}
	for _, sg := range snapshotDetail.LinkedStorageGroup {
		sgMap := make(map[string]attr.Value)
		sgMap["name"] = types.StringValue(sg.Name)
		sgMap["source_volume_name"] = types.StringValue(sg.SourceVolumeName)
		sgMap["linked_volume_name"] = types.StringValue(sg.LinkedVolumeName)
		sgMap["tracks"] = types.Int64Value(sg.Tracks)
		sgMap["track_size"] = types.Int64Value(sg.TrackSize)
		sgMap["percentage_copied"] = types.Int64Value(sg.PercentageCopied)
		sgMap["linked_creation_timestamp"] = types.StringValue(sg.LinkedCreationTimestamp)
		sgMap["defined"] = types.BoolValue(*sg.Defined)
		sgMap["background_define_in_progress"] = types.BoolValue(*sg.BackgroundDefineInProgress)

		sgbject, _ := types.ObjectValue(typeKey, sgMap)
		sgObjects = append(sgObjects, sgbject)
	}
	return types.ListValue(types.ObjectType{AttrTypes: typeKey}, sgObjects)
}
