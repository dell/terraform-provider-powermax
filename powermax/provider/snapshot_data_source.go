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
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/constants"
	"terraform-provider-powermax/powermax/helper"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &snapshotDataSource{}
	_ datasource.DataSourceWithConfigure = &snapshotDataSource{}
)

// NewSnapshotDataSource is a helper function to simplify the provider implementation.
func NewSnapshotDataSource() datasource.DataSource {
	return &snapshotDataSource{}
}

// snapshotDataSource is the data source implementation.
type snapshotDataSource struct {
	client *client.Client
}

func (d *snapshotDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot"
}

func (d *snapshotDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source for a specific StorageGroup Snapshots in PowerMax array.",
		Description:         "Data source for a specific StorageGroup Snapshots in PowerMax array.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier",
				Computed:    true,
			},
			"snapshots": schema.ListNestedAttribute{
				Description:         "List of Snapshots",
				MarkdownDescription: "List of Snapshots",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
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
				},
			},
		},
		Blocks: map[string]schema.Block{
			"storage_group": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required: true,
					},
				},
			},
		},
	}
}

func (d *snapshotDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if provider is not config
	if req.ProviderData == nil {
		return
	}

	client, err := req.ProviderData.(*client.Client)

	if !err {
		resp.Diagnostics.AddError(
			"Unexpected Resource Config Failure",
			fmt.Sprintf("Expected client, %T. Please report this issue to the provider developers", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *snapshotDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.SnapshotDataSourceModel
	var plan models.SnapshotDataSourceModel
	tflog.Info(ctx, "Attempting to read snapshots")
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	list, _, err := helper.GetStorageGroupSnapshots(ctx, *d.client, plan.StorageGroup.Name.ValueString())
	if err != nil {
		errStr := constants.ReadSnapshots + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the list of snapshots",
			message,
		)
		return
	}

	// Get the list of snapids
	for _, sngc := range list.SnapshotNamesAndCounts {
		val, _, err := helper.GetStorageGroupSnapshotSnapIDs(ctx, *d.client, plan.StorageGroup.Name.ValueString(), *sngc.Name)
		if err != nil {
			errStr := constants.ReadSnapshots + " with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError(
				"Error getting the list of snapshots Ids",
				message,
			)
			return
		}
		for _, id := range val.Snapids {
			var detail models.SnapshotDetailModal
			snapDetail, _, err := helper.GetSnapshotSnapIDSG(ctx, *d.client, plan.StorageGroup.Name.ValueString(), *sngc.Name, id)
			if err != nil {
				errStr := constants.ReadSnapshots + " with error: "
				message := helper.GetErrorString(err, errStr)
				resp.Diagnostics.AddError(
					"Error getting the list of snapshots snapIds",
					message,
				)
				return
			}
			errState := helper.UpdateSnapshotDatasourceState(ctx, snapDetail, &detail)
			if errState != nil {
				errStr := constants.ReadSnapshots + " with error: "
				message := helper.GetErrorString(errState, errStr)
				resp.Diagnostics.AddError(
					"Error getting the list of snapshots details",
					message,
				)
				return
			}
			state.Snapshots = append(state.Snapshots, detail)
		}
	}
	state.ID = types.StringValue("snapshot-datasource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}
