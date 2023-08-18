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
	"terraform-provider-powermax/powermax/helper"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &StorageGroupDataSource{}
var _ datasource.DataSourceWithConfigure = &StorageGroupDataSource{}

// NewStorageGroupDataSource creates a new data source for StorageGroup.
func NewStorageGroupDataSource() datasource.DataSource {
	return &StorageGroupDataSource{}
}

// StorageGroupDataSource defines the data source implementation.
type StorageGroupDataSource struct {
	client *client.Client
}

// Metadata defines the data source metadata.
func (d *StorageGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storagegroup"
}

// Schema defines the data source schema.
func (d *StorageGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data Source for reading StorageGroups in PowerMax array",
		Description:         "Data Source for reading StorageGroups in PowerMax array",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Placeholder value to run tests",
			},
			"storage_groups": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of storage group attributes",
				MarkdownDescription: "List of storage group attributes",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "The ID of the storage group",
							MarkdownDescription: "The ID of the storage group",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "The name of the storage group",
							MarkdownDescription: "The name of the storage group",
						},
						"slo": schema.StringAttribute{
							Computed:            true,
							Description:         "The service level associated with the storage group",
							MarkdownDescription: "The service level associated with the storage group",
						},
						"srp_id": schema.StringAttribute{
							Computed:            true,
							Description:         "The SRP to be associated with the Storage Group. An existing SRP or 'none' must be specified",
							MarkdownDescription: "The SRP to be associated with the Storage Group. An existing SRP or 'none' must be specified",
						},
						"service_level": schema.StringAttribute{
							Computed:            true,
							Description:         "The service level associated with the storage group",
							MarkdownDescription: "The service level associated with the storage group",
						},
						"workload": schema.StringAttribute{
							Computed:            true,
							Description:         "The workload associated with the storage group",
							MarkdownDescription: "The workload associated with the storage group",
						},
						"slo_compliance": schema.StringAttribute{
							Computed:            true,
							Description:         "The service level compliance status of the storage group",
							MarkdownDescription: "The service level compliance status of the storage group",
						},
						"num_of_vols": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of volumes associated with the storage group",
							MarkdownDescription: "The number of volumes associated with the storage group",
						},
						"num_of_child_sgs": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of child storage groups associated with the storage group",
							MarkdownDescription: "The number of child storage groups associated with the storage group",
						},
						"num_of_parent_sgs": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of parent storage groups associated with the storage group",
							MarkdownDescription: "The number of parent storage groups associated with the storage group",
						},
						"num_of_masking_views": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of masking views associated with the storage group",
							MarkdownDescription: "The number of masking views associated with the storage group",
						},
						"num_of_snapshots": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of snapshots associated with the storage group",
							MarkdownDescription: "The number of snapshots associated with the storage group",
						},
						"num_of_snapshot_policies": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of snapshot policies associated with the storage group",
							MarkdownDescription: "The number of snapshot policies associated with the storage group",
						},
						"cap_gb": schema.NumberAttribute{
							Computed:            true,
							Description:         "The capacity of the storage group",
							MarkdownDescription: "The capacity of the storage group",
						},
						"device_emulation": schema.StringAttribute{
							Computed:            true,
							Description:         "The emulation of the volumes in the storage group",
							MarkdownDescription: "The emulation of the volumes in the storage group",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							Description:         "The storage group type",
							MarkdownDescription: "The storage group type",
						},
						"unprotected": schema.BoolAttribute{
							Computed:            true,
							Description:         "States whether the storage group is protected",
							MarkdownDescription: "States whether the storage group is protected",
						},
						"child_storage_group": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							Description:         "The child storage group(s) associated with the storage group",
							MarkdownDescription: "The child storage group(s) associated with the storage group",
						},
						"parent_storage_group": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							Description:         "The parent storage group(s) associated with the storage group",
							MarkdownDescription: "The parent storage group(s) associated with the storage group",
						},
						"maskingview": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							Description:         "The masking views associated with the storage group",
							MarkdownDescription: "The masking views associated with the storage group",
						},
						"snapshot_policies": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							Description:         "The snapshot policies associated with the storage group",
							MarkdownDescription: "The snapshot policies associated with the storage group",
						},
						"host_io_limit": schema.ObjectAttribute{
							Computed:            true,
							Description:         "Host IO limit of the storage group",
							MarkdownDescription: "Host IO limit of the storage group",
							AttributeTypes: map[string]attr.Type{
								"host_io_limit_io_sec": types.StringType,
								"host_io_limit_mb_sec": types.StringType,
								"dynamic_distribution": types.StringType,
							},
						},
						"compression": schema.BoolAttribute{
							Computed:            true,
							Description:         "States whether compression is enabled on storage group",
							MarkdownDescription: "States whether compression is enabled on storage group",
						},
						"compression_ratio": schema.StringAttribute{
							Computed:            true,
							Description:         "States whether compression is enabled on storage group",
							MarkdownDescription: "States whether compression is enabled on storage group",
						},
						"compression_ratio_to_one": schema.NumberAttribute{
							Computed:            true,
							Description:         "Compression ratio numeric value of the storage group",
							MarkdownDescription: "Compression ratio numeric value of the storage group",
						},
						"vp_saved_percent": schema.NumberAttribute{
							Computed:            true,
							Description:         "VP saved percentage figure",
							MarkdownDescription: "VP saved percentage figure",
						},
						"tags": schema.StringAttribute{
							Computed:            true,
							Description:         "The tags associated with the storage group",
							MarkdownDescription: "The tags associated with the storage group",
						},
						"uuid": schema.StringAttribute{
							Computed:            true,
							Description:         "Storage Group UUID",
							MarkdownDescription: "Storage Group UUID",
						},
						"unreducible_data_gb": schema.NumberAttribute{
							Computed:            true,
							Description:         "The amount of unreducible data in Gb.",
							MarkdownDescription: "The amount of unreducible data in Gb.",
						},
						"volume_ids": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							Description:         "The IDs of the volume associated with the storage group.",
							MarkdownDescription: "The IDs of the volume associated with the storage group.",
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"names": schema.SetAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *StorageGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pmaxClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = pmaxClient
}

func (d *StorageGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Storage Group...")
	var data models.StorageGroupDataSourceModel
	var state models.StorageGroupDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var sgIDs []string
	// Get storage group IDs from config or query all if not specified
	if data.StorageGroupFilter == nil || len(data.StorageGroupFilter.IDs) == 0 {
		storageGroupIDList, _, err := helper.GetStorageGroupList(ctx, d.client)
		if err != nil {
			errStr := ""
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Error reading storage group ids:", message)
			return
		}
		sgIDs = storageGroupIDList.StorageGroupId
	} else {
		// get ids from data.StorageGroupFilter.IDs and assign to sgIDs
		for _, sg := range data.StorageGroupFilter.IDs {
			sgIDs = append(sgIDs, sg.ValueString())
		}
	}

	// iterate sgIDs and GetStorageGroup with each id
	for _, sgID := range sgIDs {
		var sg models.StorageGroupResourceModel
		err := helper.UpdateSgState(ctx, d.client, sgID, &sg)
		if err != nil {
			resp.Diagnostics.AddError("Error reading storage group", err.Error())
			return
		}
		state.StorageGroups = append(state.StorageGroups, sg)
	}
	state.ID = types.StringValue("storage-group-data-source")
	state.StorageGroupFilter = data.StorageGroupFilter

	if len(state.StorageGroups) > 0 {
		tflog.Info(ctx, fmt.Sprintf("State: %v", state.StorageGroups[0]))
		tflog.Info(ctx, fmt.Sprintf("State: %v", state.StorageGroups[0].VolumeIDs))
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
