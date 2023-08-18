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
	"net/http"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	// ActionSnapshotRestore is used as the restore action for snasphots.
	ActionSnapshotRestore = "Restore"
	// ActionSnapshotLink is used as the link action for snasphots.
	ActionSnapshotLink = "Link"
	// ActionSnapshotSetMode is used as the setMode action for snasphots.
	ActionSnapshotSetMode = "SetMode"
	// ActionSnapshotRename is used as the rename action for snasphots.
	ActionSnapshotRename = "Rename"
	// ActionSnapshotTimeToLive is used as the ttl action for snasphots.
	ActionSnapshotTimeToLive = "SetTimeToLive"
	// ActionSnapshotSecure is used as the secure action for snasphots.
	ActionSnapshotSecure = "SetSecure"
	// ActionSnapshotUnlink is used as the unlink action for snasphots.
	ActionSnapshotUnlink = "Unlink"
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

// UpdateSnapshotResourceState Update Snaposhot state.
func UpdateSnapshotResourceState(ctx context.Context, snapshotDetail *powermax.SnapVXSnapshotInstance, state *models.SnapshotResourceModel) error {
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

// ModifySnapshot Do the modify action.
func ModifySnapshot(ctx context.Context, client client.Client, plan *models.SnapshotResourceModel, state *models.SnapshotResourceModel) error {

	modifyParam := client.PmaxOpenapiClient.ReplicationApi.UpdateSnapshotSnapID(ctx, client.SymmetrixID, state.StorageGroup.Name.ValueString(), state.Snapshot.Name.ValueString(), state.Snapid.ValueInt64())
	actionList := []string{ActionSnapshotRestore, ActionSnapshotLink, ActionSnapshotSetMode, ActionSnapshotTimeToLive, ActionSnapshotSecure}
	// Do the rename modification first so it does not interfere with tho other modifications if there are multiple which need the snapshot name value
	if plan.Snapshot.Name.ValueString() != state.Snapshot.Name.ValueString() {
		modifyParam := modifyParam.StorageGroupSnapshotInstanceUpdate(powermax.StorageGroupSnapshotInstanceUpdate{
			Action: ActionSnapshotRename,
			Rename: &powermax.SnapVXRenameOptions{
				NewSnapshotName: plan.Snapshot.Name.ValueString(),
			},
		})
		_, _, err := modifyParam.Execute()
		if err != nil {
			return err
		}
	}

	// Update the modify Param after the possible rename
	modifyParam = client.PmaxOpenapiClient.ReplicationApi.UpdateSnapshotSnapID(ctx, client.SymmetrixID, state.StorageGroup.Name.ValueString(), plan.Snapshot.Name.ValueString(), state.Snapid.ValueInt64())
	// Loop through and check to see if any actions need to be done
	for _, v := range actionList {

		switch v {
		case ActionSnapshotRestore:
			if plan.Snapshot.Restore != nil && (state.Snapshot.SetMode == nil || plan.Snapshot.Restore.Enable.ValueBool() != state.Snapshot.Restore.Enable.ValueBool()) {
				modifyParam := modifyParam.StorageGroupSnapshotInstanceUpdate(powermax.StorageGroupSnapshotInstanceUpdate{
					Action: ActionSnapshotRestore,
					Restore: &powermax.SnapVxRestoreOptions{
						Remote: plan.Snapshot.Remote.ValueBoolPointer(),
					},
				})
				_, _, err := modifyParam.Execute()
				if err != nil {
					return err
				}
			}
		case ActionSnapshotLink:
			if plan.Snapshot.Link != nil && (state.Snapshot.SetMode == nil || plan.Snapshot.Link.Enable.ValueBool() != state.Snapshot.Link.Enable.ValueBool()) {
				if plan.Snapshot.Link.Enable.ValueBool() {
					modifyParam := modifyParam.StorageGroupSnapshotInstanceUpdate(powermax.StorageGroupSnapshotInstanceUpdate{
						Action: ActionSnapshotLink,
						Link: &powermax.SnapVxLinkOptions{
							StorageGroupName: plan.Snapshot.Link.TargetStorageGroup.ValueString(),
							NoCompression:    plan.Snapshot.Link.NoCompression.ValueBoolPointer(),
							Copy:             plan.Snapshot.Link.Copy.ValueBoolPointer(),
							Remote:           plan.Snapshot.Link.Remote.ValueBoolPointer(),
						},
					})
					_, _, err := modifyParam.Execute()
					if err != nil {
						return err
					}
				} else {
					modifyParam := modifyParam.StorageGroupSnapshotInstanceUpdate(powermax.StorageGroupSnapshotInstanceUpdate{
						Action: ActionSnapshotUnlink,
						Unlink: &powermax.SnapVxUnlinkOptions{
							StorageGroupName: plan.Snapshot.Link.TargetStorageGroup.ValueString(),
						},
					})
					_, _, err := modifyParam.Execute()
					if err != nil {
						return err
					}
				}
			}

		case ActionSnapshotSetMode:
			if plan.Snapshot.SetMode != nil && (state.Snapshot.SetMode == nil || plan.Snapshot.SetMode.Enable.ValueBool() != state.Snapshot.SetMode.Enable.ValueBool()) {
				modifyParam := modifyParam.StorageGroupSnapshotInstanceUpdate(powermax.StorageGroupSnapshotInstanceUpdate{
					Action: ActionSnapshotSetMode,
					SetMode: &powermax.SnapVXSetModeOptions{
						StorageGroupName: plan.Snapshot.SetMode.TargetStorageGroup.ValueString(),
						Copy:             plan.Snapshot.SetMode.Copy.ValueBoolPointer(),
					},
				})
				_, _, err := modifyParam.Execute()
				if err != nil {
					return err
				}
			}
		case ActionSnapshotRename:

		case ActionSnapshotTimeToLive:
			if plan.Snapshot.TimeToLive != nil && (state.Snapshot.TimeToLive == nil || plan.Snapshot.TimeToLive.Enable.ValueBool() != state.Snapshot.TimeToLive.Enable.ValueBool()) {
				ttl := int32(plan.Snapshot.TimeToLive.TimeToLive.ValueInt64())
				modifyParam := modifyParam.StorageGroupSnapshotInstanceUpdate(powermax.StorageGroupSnapshotInstanceUpdate{
					Action: ActionSnapshotTimeToLive,
					TimeToLive: &powermax.SnapVxTimeToLiveOptions{
						TimeToLive:  &ttl,
						TimeInHours: plan.Snapshot.TimeToLive.TimeInHours.ValueBoolPointer(),
					},
				})
				_, _, err := modifyParam.Execute()
				if err != nil {
					return err
				}
			}
		case ActionSnapshotSecure:
			if plan.Snapshot.Secure != nil && (state.Snapshot.Secure == nil || plan.Snapshot.Secure.Enable.ValueBool() != state.Snapshot.Secure.Enable.ValueBool()) {
				secure := int32(plan.Snapshot.Secure.Secure.ValueInt64())
				modifyParam := modifyParam.StorageGroupSnapshotInstanceUpdate(powermax.StorageGroupSnapshotInstanceUpdate{
					Action: ActionSnapshotSecure,
					Secure: &powermax.SnapVxSecureOptions{
						Secure:      &secure,
						TimeInHours: plan.Snapshot.Secure.TimeInHours.ValueBoolPointer(),
					},
				})
				_, _, err := modifyParam.Execute()
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// GetStorageGroupSnapshots get SG snapshots
func GetStorageGroupSnapshots(ctx context.Context, client client.Client, sgName string) (*powermax.StorageGroupSnapshotList, *http.Response, error) {
	return client.PmaxOpenapiClient.ReplicationApi.GetStorageGroupSnapshots(ctx, client.SymmetrixID, sgName).Execute()
}

// GetStorageGroupSnapshotSnapIDs get SG snapshots snap ids
func GetStorageGroupSnapshotSnapIDs(ctx context.Context, client client.Client, sgName string, snapIDName string) (*powermax.StorageGroupSnapshotSnapIDList, *http.Response, error) {
	return client.PmaxOpenapiClient.ReplicationApi.GetStorageGroupSnapshotSnapIDs(ctx, client.SymmetrixID, sgName, snapIDName).Execute()
}

// GetSnapshotSnapIDSG get SG snapshots snap details
func GetSnapshotSnapIDSG(ctx context.Context, client client.Client, sgName string, snapIDName string, snapID int64) (*powermax.SnapVXSnapshotInstance, *http.Response, error) {
	return client.PmaxOpenapiClient.ReplicationApi.GetSnapshotSnapIDSG(ctx, client.SymmetrixID, sgName, snapIDName, snapID).Execute()
}

// CreateSnapshot creates a snapshot on a particular SG
func CreateSnapshot(ctx context.Context, client client.Client, sgName string, plan models.SnapshotResourceModel) (*powermax.SnapVXSnapshotGeneration, *http.Response, error) {
	// Create Param Attributes
	snapshotCreateParam := powermax.StorageGroupSnapshotCreate{}
	if plan.Snapshot.Secure != nil && plan.Snapshot.Secure.Secure.ValueInt64() != 0 {
		secure := int32(plan.Snapshot.Secure.Secure.ValueInt64())
		snapshotCreateParam.Secure = &secure
	}
	if plan.Snapshot.TimeToLive != nil && plan.Snapshot.TimeToLive.TimeToLive.ValueInt64() != 0 {
		ttl := int32(plan.Snapshot.TimeToLive.TimeToLive.ValueInt64())
		snapshotCreateParam.TimeToLive = &ttl
		snapshotCreateParam.TimeInHours = plan.Snapshot.TimeToLive.TimeInHours.ValueBoolPointer()
	}
	snapshotCreateParam.Bothsides = plan.Snapshot.Bothsides.ValueBoolPointer()
	snapshotCreateParam.SnapshotName = plan.Snapshot.Name.ValueString()

	createParam := client.PmaxOpenapiClient.ReplicationApi.CreateSnapshot1(ctx, client.SymmetrixID, plan.StorageGroup.Name.ValueString())
	createParam = createParam.StorageGroupSnapshotCreate(snapshotCreateParam)

	return createParam.Execute()
}
