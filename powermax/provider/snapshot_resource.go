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

package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/constants"
	"terraform-provider-powermax/powermax/helper"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &snapshotResource{}
var _ resource.ResourceWithConfigure = &snapshotResource{}
var _ resource.ResourceWithImportState = &snapshotResource{}

// NewSnapshotResource is a helper function to simplify the provider implementation.
func NewSnapshotResource() resource.Resource {
	return &snapshotResource{}
}

// snapshotResource defines the resource implementation.
type snapshotResource struct {
	client *client.Client
}

// Metadata Resource metadata.
func (r *snapshotResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot"
}

// Schema Resource schema.
func (r *snapshotResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Resource for managing Snapshots in PowerMax array. Supported Update (name, secure, time_to_live, link, restore). PowerMax Snaphots is a local replication solution that is designed to nondisruptively create point-in-time copies (snapshots) of critical data. ",
		Description:         "Resource for managing Snapshots in PowerMax array. Supported Update (name, secure, time_to_live, link, restore). PowerMax Snaphots is a local replication solution that is designed to nondisruptively create point-in-time copies (snapshots) of critical data.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier",
				Computed:    true,
			},

			"name": schema.StringAttribute{
				Description:         "Name of a snapshot",
				MarkdownDescription: "Name of a snapshot",
				Computed:            true,
			},
			"generation": schema.Int64Attribute{
				Description:         "Number of generation for the snapshot",
				MarkdownDescription: "Number of generation for the snapshot",
				Computed:            true,
				Optional:            true,
			},
			"snapid": schema.Int64Attribute{
				Description:         "Unique Snap ID for Snapshot",
				MarkdownDescription: "Unique Snap ID for Snapshot",
				Computed:            true,
				Optional:            true,
			},
			"timestamp": schema.StringAttribute{
				Description:         "Timestamp of the snapshot generation",
				MarkdownDescription: "Timestamp of the snapshot generation",
				Computed:            true,
			},
			"timestamp_utc": schema.StringAttribute{
				Description:         "The timestamp of the snapshot generation in milliseconds since 1970",
				MarkdownDescription: "The timestamp of the snapshot generation in milliseconds since 1970",
				Computed:            true,
			},
			"state": schema.ListAttribute{
				Description:         "The state of the snapshot generation",
				MarkdownDescription: "The state of the snapshot generation",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"num_source_volumes": schema.Int64Attribute{
				Description:         "The number of source volumes in the snapshot generation",
				MarkdownDescription: "The number of source volumes in the snapshot generation",
				Computed:            true,
				Optional:            true,
			},
			"source_volume": schema.ListNestedAttribute{
				Description: " The source volumes of the snapshot generation",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "The name of the SnapVX snapshot generation source volume",
							MarkdownDescription: "The name of the SnapVX snapshot generation source volume",
							Computed:            true,
						},
						"capacity": schema.Int64Attribute{
							Description:         "The capacity of the snapshot volume in cylinders",
							MarkdownDescription: "The capacity of the snapshot volume in cylinders",
							Computed:            true,
						},
						"capacity_gb": schema.Float64Attribute{
							Description:         "The capacity of the snapshot volume in GB",
							MarkdownDescription: "The capacity of the snapshot volume in GB",
							Computed:            true,
						},
					},
				},
			},
			"num_storage_group_volumes": schema.Int64Attribute{
				Description:         "The number of non-gatekeeper storage group volumes",
				MarkdownDescription: "The number of non-gatekeeper storage group volumes",
				Computed:            true,
			},
			"tracks": schema.Int64Attribute{
				Description:         "The number of source tracks that have been overwritten by the host",
				MarkdownDescription: "The number of source tracks that have been overwritten by the host",
				Computed:            true,
				Optional:            true,
			},
			"non_shared_tracks": schema.Int64Attribute{
				Description:         "The number of tracks uniquely allocated for this sn.apshots delta. This is an approximate indication of the number of tracks that will be returned to the SRP if this snapshot is terminated.",
				MarkdownDescription: "The number of tracks uniquely allocated for this snapshots delta. This is an approximate indication of the number of tracks that will be returned to the SRP if this snapshot is terminated.",
				Computed:            true,
				Optional:            true,
			},
			"time_to_live_expiry_date": schema.StringAttribute{
				Description:         "When the snapshot will expire once it is not linked",
				MarkdownDescription: "When the snapshot will expire once it is not linked",
				Computed:            true,
				Optional:            true,
			},
			"secure_expiry_date": schema.StringAttribute{
				Description:         "When the snapshot will expire once it is not linked",
				MarkdownDescription: "When the snapshot will expire once it is not linked",
				Computed:            true,
				Optional:            true,
			},
			"expired": schema.BoolAttribute{
				Description:         "Set if this generation secure has expired",
				MarkdownDescription: "Set if this generation secure has expired",
				Computed:            true,
			},
			"linked": schema.BoolAttribute{
				Description:         "Set if this generation is SnapVX linked",
				MarkdownDescription: "Set if this generation is SnapVX linked",
				Computed:            true,
			},
			"restored": schema.BoolAttribute{
				Description:         "Set if this generation is SnapVX linked",
				MarkdownDescription: "Set if this generation is SnapVX linked",
				Computed:            true,
			},
			"linked_storage_group_names": schema.ListAttribute{
				Description:         "Linked storage group names. Only populated if the generation is linked",
				MarkdownDescription: "Linked storage group names. Only populated if the generation is linked",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"linked_storage_group": schema.ListNestedAttribute{
				Description:         "Linked storage group and volume information. Only populated if the generation is linked",
				MarkdownDescription: "Linked storage group and volume information. Only populated if the generation is linked",
				Computed:            true,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "The storage group name",
							MarkdownDescription: "The storage group name",
							Computed:            true,
						},
						"source_volume_name": schema.StringAttribute{
							Description:         "The source volumes name",
							MarkdownDescription: "The source volumes name",
							Computed:            true,
						},
						"linked_volume_name": schema.StringAttribute{
							Description:         "The linked volumes name",
							MarkdownDescription: "The linked volumes name",
							Computed:            true,
						},
						"tracks": schema.Int64Attribute{
							Description:         "Number of tracks",
							MarkdownDescription: "Number of tracks",
							Computed:            true,
						},
						"track_size": schema.Int64Attribute{
							Description:         "Size of the tracks",
							MarkdownDescription: "Size of the tracks.",
							Computed:            true,
						},
						"percentage_copied": schema.Int64Attribute{
							Description:         "Percentage of tracks copied",
							MarkdownDescription: "Percentage of tracks copied",
							Computed:            true,
						},
						"linked_creation_timestamp": schema.StringAttribute{
							Description:         "The average timestamp of all linked volumes that are linked",
							MarkdownDescription: "The average timestamp of all linked volumes that are linked",
							Computed:            true,
						},
						"defined": schema.BoolAttribute{
							Description:         "When the snapshot link has been fully defined",
							MarkdownDescription: "When the snapshot link has been fully defined",
							Computed:            true,
							Optional:            true,
						},
						"background_define_in_progress": schema.BoolAttribute{
							Description:         "When the snapshot link is being defined",
							MarkdownDescription: "When the snapshot link is being defined",
							Computed:            true,
							Optional:            true,
						},
					},
				},
			},
			"persistent": schema.BoolAttribute{
				Description:         "Set if this snapshot is persistent.  Only applicable to policy based snapshots",
				MarkdownDescription: "Set if this snapshot is persistent.  Only applicable to policy based snapshots",
				Computed:            true,
				Optional:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"storage_group": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description:         "Name of the storage group you would like to take a snapshot.",
						MarkdownDescription: "Name of the storage group you would like to take a snapshot.",
						Required:            true,
					},
				},
			},
			"snapshot_actions": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description:         "Name of the snapshot. (Update Supported)",
						MarkdownDescription: "Name of the snapshot. (Update Supported)",
						Required:            true,
					},
					"restore": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "Restore a snapshot generation. (Update Supported)",
						MarkdownDescription: "Restore a snapshot generation. (Update Supported)",
						Attributes: map[string]schema.Attribute{
							"enable": schema.BoolAttribute{
								Description:         "enable defaults to false. Flag to enable restore on the snapshot",
								MarkdownDescription: "enable defaults to false. Flag to enable restore on the snapshot",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"remote": schema.BoolAttribute{
								Description:         "remote defaults to false. If true, The target storage group will not have compression turned on when the SRP is compression capable. Option Used in Action Link",
								MarkdownDescription: "remote defaults to false. If true, The target storage group will not have compression turned on when the SRP is compression capable. Option Used in Action Link",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
						},
					},
					"link": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "Link a snapshot generation. (Update Supported)",
						MarkdownDescription: "Link a snapshot generation. (Update Supported)",
						Attributes: map[string]schema.Attribute{
							"target_storage_group": schema.StringAttribute{
								Description:         "The target storage group to link the snapshot too",
								MarkdownDescription: "The target storage group to link the snapshot too",
								Computed:            true,
								Optional:            true,
							},
							"enable": schema.BoolAttribute{
								Description:         "enable defaults to false. Flag to enable link on the snapshot",
								MarkdownDescription: "enable defaults to false. Flag to enable link on the snapshot",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"no_compression": schema.BoolAttribute{
								Description:         "no_compression defaults to false. If true, The target storage group will not have compression turned on when the SRP is compression capable. Option Used in Action Link",
								MarkdownDescription: "no_compression defaults to false. If true, The target storage group will not have compression turned on when the SRP is compression capable. Option Used in Action Link",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"remote": schema.BoolAttribute{
								Description:         "remote defaults to false. If true, The target storage group will not have compression turned on when the SRP is compression capable. Option Used in Action Link",
								MarkdownDescription: "remote defaults to false. If true, The target storage group will not have compression turned on when the SRP is compression capable. Option Used in Action Link",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"copy": schema.BoolAttribute{
								Description:         "copy defaults to false. If true Sets the link copy mode to perform background copy to the target volume(s).",
								MarkdownDescription: "copy defaults to false. If true Sets the link copy mode to perform background copy to the target volume(s).",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
						},
					},
					"time_to_live": schema.SingleNestedAttribute{
						Description:         "Set the number of days or hours for a snapshot generation before it auto-terminates (provided it is not linked). (Update Supported)",
						MarkdownDescription: "Set the number of days or hours for a snapshot generation before it auto-terminates (provided it is not linked). (Update Supported)",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"enable": schema.BoolAttribute{
								Description:         "enable defaults to false. Flag to enable link on the snapshot",
								MarkdownDescription: "enable defaults to false. Flag to enable link on the snapshot",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"time_in_hours": schema.BoolAttribute{
								Description:         "time_in_hours or Days defaults to Days. False is days, true is hours.",
								MarkdownDescription: "time_in_hours or Days defaults to Days. False is days, true is hours.",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"time_to_live": schema.Int64Attribute{
								Description:         "time_to_live defaults to 1 day. Gives the total time before expiry for these actions.",
								MarkdownDescription: "time_to_live defaults to 1 day. Gives the total time before expiry for these actions.",
								Optional:            true,
								Computed:            true,
								Default:             int64default.StaticInt64(1),
							},
						},
					},
					"secure": schema.SingleNestedAttribute{
						Description:         "Set the number of days or hours for a snapshot generation to be secure before it auto-terminates (provided it is not linked). (Update Supported)",
						MarkdownDescription: "Set the number of days or hours for a snapshot generation to be secure before it auto-terminates (provided it is not linked). (Update Supported)",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"enable": schema.BoolAttribute{
								Description:         "enable defaults to false. Flag to enable link on the snapshot",
								MarkdownDescription: "enable defaults to false. Flag to enable link on the snapshot",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"time_in_hours": schema.BoolAttribute{
								Description:         "time_in_hours or Days defaults to Days. False is days, true is hours.",
								MarkdownDescription: "time_in_hours or Days defaults to Days. False is days, true is hours.",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"secure": schema.Int64Attribute{
								Description:         "secure defaults to 1 day. The time that the snapshot generation is to be secure for.",
								MarkdownDescription: "secure defaults to 1 day. The time that the snapshot generation is to be secure for.",
								Optional:            true,
								Computed:            true,
								Default:             int64default.StaticInt64(1),
							},
						},
					},
					// Options during create snapshot
					"both_sides": schema.BoolAttribute{
						Description:         "both_sides defaults to false. Performs the operation on both locally and remotely associated snapshots.",
						MarkdownDescription: "both_sides defaults to false. Performs the operation on both locally and remotely associated snapshots.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},

					"remote": schema.BoolAttribute{
						Description:         "remote defaults to false. If true, The target storage group will not have compression turned on when the SRP is compression capable.",
						MarkdownDescription: "remote defaults to false. If true, The target storage group will not have compression turned on when the SRP is compression capable.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
				},
			},
		},
	}
}

// Configure the resource.
func (r *snapshotResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pmaxClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = pmaxClient
}

// Create a snapshot.
func (r *snapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "creating snapshot")
	var plan models.SnapshotResourceModel
	var state models.SnapshotResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.StorageGroup.Name.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Error creating snapshot",
			fmt.Sprintf("Could not create snapshot %s with error: %s", plan.Snapshot.Name, "storage group name cannot be empty"),
		)
		return
	}
	_, _, err := helper.CreateSnapshot(ctx, *r.client, plan.StorageGroup.Name.ValueString(), plan)
	if err != nil {
		errStr := fmt.Sprintf("Could not create snapshot %s with error:", plan.Snapshot.Name)
		msgStr := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating snapshot",
			msgStr,
		)
		return
	}

	// Get the new snapID Id
	val, _, err := helper.GetStorageGroupSnapshotSnapIDs(ctx, *r.client, plan.StorageGroup.Name.ValueString(), plan.Snapshot.Name.ValueString())
	if err != nil {
		errStr := constants.ReadSnapshots + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the new snapID",
			message,
		)
		return
	}

	// Get the new Snapshot
	snapDetail, _, err := helper.GetSnapshotSnapIDSG(ctx, *r.client, plan.StorageGroup.Name.ValueString(), plan.Snapshot.Name.ValueString(), val.Snapids[0])
	if err != nil {
		errStr := fmt.Sprintf("Could not find snapshot %s after create with error:", plan.Snapshot.Name)
		msgStr := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating snapshot",
			msgStr,
		)
		return
	}
	errState := helper.UpdateSnapshotResourceState(ctx, snapDetail, &state)
	if errState != nil {
		resp.Diagnostics.AddError(
			"Error creating snapshot",
			errState.Error(),
		)
		return
	}
	state.ID = types.StringValue("snapshot-resource")
	state.StorageGroup = plan.StorageGroup
	state.Snapshot = plan.Snapshot
	// Save plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read a snapshot.
func (r *snapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "reading snapshot")
	var state models.SnapshotResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	snapDetail, _, err := helper.GetSnapshotSnapIDSG(ctx, *r.client, state.StorageGroup.Name.ValueString(), state.Name.ValueString(), state.Snapid.ValueInt64())
	if err != nil {
		errStr := fmt.Sprintf("Could not find snapshot %s with error:", state.Name)
		msgStr := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading snapshot",
			msgStr,
		)
		return
	}
	errState := helper.UpdateSnapshotResourceState(ctx, snapDetail, &state)
	if errState != nil {
		resp.Diagnostics.AddError(
			"Error reading snapshot",
			errState.Error(),
		)
		return
	}
	// Save plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a snapshot.
func (r *snapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "updating snapshot")
	var plan models.SnapshotResourceModel
	var state models.SnapshotResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diagsState := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diagsState...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := helper.ModifySnapshot(ctx, *r.client, &plan, &state)
	if err != nil {
		errStr := constants.UpdateSnapshot + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating snapshot",
			message,
		)
		return
	}
	// Read and update state after the modification
	getParam := r.client.PmaxOpenapiClient.ReplicationApi.GetSnapshotSnapIDSG(ctx, r.client.SymmetrixID, state.StorageGroup.Name.ValueString(), plan.Snapshot.Name.ValueString(), state.Snapid.ValueInt64())
	snapDetail, _, err := getParam.Execute()
	if err != nil {
		errStr := fmt.Sprintf("Error reading snapshot %s after update with error:", state.Name)
		msgStr := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading snapshot",
			msgStr,
		)
		return
	}
	errState := helper.UpdateSnapshotResourceState(ctx, snapDetail, &state)
	if errState != nil {
		resp.Diagnostics.AddError(
			"Error reading snapshot",
			errState.Error(),
		)
		return
	}
	state.StorageGroup = plan.StorageGroup
	state.Snapshot = plan.Snapshot
	// Save plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes a snapshot.
func (r *snapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting snapshot")
	var state models.SnapshotResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteParam := r.client.PmaxOpenapiClient.ReplicationApi.DeleteSnapshotSnapID(ctx, r.client.SymmetrixID, state.StorageGroup.Name.ValueString(), state.Name.ValueString(), state.Snapid.ValueInt64())
	_, err := deleteParam.Execute()
	if err != nil {
		errStr := fmt.Sprintf("Could not delete snapshot %s with error:", state.Name)
		msgStr := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error deleting snapshot",
			msgStr,
		)
		return
	}
}

// ImportState imports a Snapshot.
func (r *snapshotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "importing snapshot")
	id := req.ID
	ids := strings.Split(id, ".")
	tflog.Info(ctx, fmt.Sprintf("id: %s ids %v length %v", id, ids, len(ids)))
	sgName := ""
	snapshotName := ""
	if len(ids) >= 2 {
		sgName = ids[0]
		snapshotName = ids[1]
	} else {
		resp.Diagnostics.AddError(
			"Error importing snapshot",
			"The import ID must be 'storage_group_name.snapshot_name'",
		)
		return
	}

	var state models.SnapshotResourceModel
	// Get the snapID Id
	snapIDParam := r.client.PmaxOpenapiClient.ReplicationApi.GetStorageGroupSnapshotSnapIDs(ctx, r.client.SymmetrixID, sgName, snapshotName)
	val, _, err := snapIDParam.Execute()
	if err != nil {
		errStr := constants.ReadSnapshots + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error importing snapshot",
			message,
		)
		return
	}
	// Get the details
	snapDetail, _, err := helper.GetSnapshotSnapIDSG(ctx, *r.client, sgName, snapshotName, val.Snapids[0])
	if err != nil {
		errStr := fmt.Sprintf("Could not find snapshot %s with error:", state.Name)
		msgStr := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error importing snapshot",
			msgStr,
		)
		return
	}
	errState := helper.UpdateSnapshotResourceState(ctx, snapDetail, &state)
	if errState != nil {
		resp.Diagnostics.AddError(
			"Error importing snapshot",
			errState.Error(),
		)
		return
	}
	state.ID = types.StringValue("snapshot-resource")
	state.Snapshot = &models.SnapshotResourceFields{
		Name:      basetypes.NewStringValue(snapshotName),
		Bothsides: basetypes.NewBoolValue(false),
		Remote:    basetypes.NewBoolValue(false),
	}
	state.StorageGroup = &models.FilterTypeSnapshot{
		Name: basetypes.NewStringValue(sgName),
	}
	// Save plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
