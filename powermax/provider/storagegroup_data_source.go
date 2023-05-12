// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package provider

import (
	"context"
	"fmt"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/helper"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &StorageGroupDataSource{}
var _ datasource.DataSourceWithConfigure = &StorageGroupDataSource{}

func NewStorageGroupDataSource() datasource.DataSource {
	return &StorageGroupDataSource{}
}

// StorageGroupDataSource defines the data source implementation.
type StorageGroupDataSource struct {
	client *client.Client
}

// StorageGroupDataSourceModel describes the data source data model.
type StorageGroupDataSourceModel struct {
	ID            types.String                       `tfsdk:"id"`
	StorageGroups []models.StorageGroupResourceModel `tfsdk:"storage_groups"`
}

func (d *StorageGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storagegroup"
}

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
				Required:            true,
				Description:         "List of storage group attributes",
				MarkdownDescription: "List of storage group attributes",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "The ID of the storage group",
							MarkdownDescription: "The ID of the storage group",
						},
						"storage_group_id": schema.StringAttribute{
							Optional:            true,
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
						"host_io_limit": schema.MapAttribute{
							Computed:            true,
							Optional:            true,
							ElementType:         types.StringType,
							Description:         "Host IO limit of the storage group",
							MarkdownDescription: "Host IO limit of the storage group",
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
							MarkdownDescription: "SThe amount of unreducible data in Gb.",
						},
					},
				},
			},
		},
	}
}

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
	var data StorageGroupDataSourceModel
	var state StorageGroupDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var sgIds []string
	// Get storage group IDs from config or query all if not specified
	if len(data.StorageGroups) == 0 {
		storageGroupIDList, err := d.client.PmaxClient.GetStorageGroupIDList(ctx, d.client.SymmetrixID, "", false)
		if err != nil {
			resp.Diagnostics.AddError("Error reading storage group ids", err.Error())
			return
		}
		sgIds = storageGroupIDList.StorageGroupIDs
	} else {
		// get ids from storageGroups and assign to sgIds
		for _, sg := range data.StorageGroups {
			sgIds = append(sgIds, sg.StorageGroupID.ValueString())
		}
	}

	// iterate sgIds and GetStorageGroup with each id
	for _, sgId := range sgIds {
		storageGroup, err := d.client.PmaxClient.GetStorageGroup(ctx, d.client.SymmetrixID, sgId)
		if err != nil {
			resp.Diagnostics.AddError("Error reading storage group with id", err.Error())
			continue
		}
		var sg models.StorageGroupResourceModel
		// Copy fields from the provider client data into the Terraform state
		err = helper.CopyFields(ctx, storageGroup, &sg)
		if err != nil {
			resp.Diagnostics.AddError("Error copying storage group fields", err.Error())
			continue
		}
		if sg.HostIOLimit.IsNull() || sg.HostIOLimit.IsUnknown() {
			sg.HostIOLimit = types.MapNull(types.StringType)
		}
		state.StorageGroups = append(state.StorageGroups, sg)
	}
	state.ID = types.StringValue(fmt.Sprintf("%s-storage-groups", sgIds))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
