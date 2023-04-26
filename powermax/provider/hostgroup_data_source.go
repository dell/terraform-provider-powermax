// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	_ datasource.DataSource              = &hostGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &hostGroupDataSource{}
)

// NewHostGroupDataSource is a helper function to simplify the provider implementation.
func NewHostGroupDataSource() datasource.DataSource {
	return &hostGroupDataSource{}
}

// coffeesDataSource is the data source implementation.
type hostGroupDataSource struct {
	client *client.Client
}

func (d *hostGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hostgroup"
}

func (d *hostGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier",
				Computed:    true,
			},
			"host_group_details": schema.ListNestedAttribute{
				Description:         "List of Hostgroups",
				MarkdownDescription: "List of Hostgroups",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host_group_id": schema.StringAttribute{
							Description:         "Id of a hostgroup",
							MarkdownDescription: "Id of a hostgroup",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							Description:         "Name of a hostgroup",
							MarkdownDescription: "Name of a hostgroup",
							Computed:            true,
						},
						"num_of_masking_views": schema.Int64Attribute{
							Description:         "Number of masking views related to a hostgroup",
							MarkdownDescription: "Number of masking views related to a hostgroup",
							Computed:            true,
						},
						"num_of_hosts": schema.Int64Attribute{
							Description:         "Number of hosts related to a hostgroup",
							MarkdownDescription: "Number of hosts related to a hostgroup",
							Computed:            true,
						},
						"num_of_initiators": schema.Int64Attribute{
							Description:         "Number of initiators related to a hostgroup",
							MarkdownDescription: "Number of initiators related to a hostgroup",
							Computed:            true,
						},
						"port_flags_override": schema.BoolAttribute{
							Description:         "Port flags are overwritten",
							MarkdownDescription: "Port flags are overwritten",
							Computed:            true,
						},
						"consistent_lun": schema.BoolAttribute{
							Description:         "Consistent lun flag set",
							MarkdownDescription: "Consistent lun flag set",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							Description:         "The host group type",
							MarkdownDescription: "The host group type",
							Computed:            true,
						},
						"host": schema.ListNestedAttribute{
							Description: "List of related host ids",
							Computed:    true,
							Optional:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"host_id": schema.StringAttribute{
										Description:         "The host id",
										MarkdownDescription: "The host id",
										Computed:            true,
									},
									"initiator": schema.ListAttribute{
										Description:         "The host initators associated with the host",
										MarkdownDescription: "The host initators associated with the host",
										Computed:            true,
										Optional:            true,
										ElementType:         types.StringType,
									},
								},
							},
						},
						"maskingview": schema.ListAttribute{
							Description: "List of masking views ids related to the host",
							Computed:    true,
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"ids": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (r *hostGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	r.client = client
}

// Read.
func (d *hostGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.HostGroupDataSourceModel
	var plan models.HostGroupDataSourceModel
	tflog.Info(ctx, "Attempting to read hostgroups")
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Apply Filter hostgroup filter
	hostgroupIds, err := helper.FilterHostGroupIds(ctx, &state, &plan, *d.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting the list of host group ids",
			constants.ReadHostGroupListDetailsErrorMsg+"with error: "+err.Error(),
		)
		return
	}

	// Get details of each of the hostgroups
	for _, hostGroupId := range hostgroupIds {
		tflog.Debug(ctx, hostGroupId)
		groupDetail, err := d.client.PmaxClient.GetHostGroupByID(ctx, d.client.SymmetrixID, hostGroupId)
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError(
				"Error getting the details of host group: "+hostGroupId,
				constants.ReadHostGroupListDetailsErrorMsg+"with error: "+err.Error(),
			)
			return
		}
		model, diag := helper.HostGroupDetailMapper(groupDetail)
		if diag.HasError() {
			resp.Diagnostics.AddError(
				"Error Unknown Error",
				constants.ReadHostGroupListDetailsErrorMsg+"with error: "+err.Error(),
			)
			return
		}
		state.HostGroupDetails = append(state.HostGroupDetails, model)
	}
	state.ID = types.StringValue("1")
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
