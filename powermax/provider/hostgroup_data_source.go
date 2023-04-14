// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package provider

import (
	"context"
	"fmt"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/constants"

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

type hostGroupDataSourceModel struct {
	ID           types.String   `tfsdk:"id"`
	HostGroupIds []types.String `tfsdk:"host_group_id"`
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
			"host_group_id": schema.ListAttribute{
				Description: "List of Host Group Ids",
				Computed:    true,
				ElementType: types.StringType,
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

// Read
func (d *hostGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state hostGroupDataSourceModel

	hostGroupResponse, err := d.client.PmaxClient.GetHostGroupList(ctx, d.client.SymmetrixID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting the list of host group ids",
			constants.ReadHostGroupListDetailsErrorMsg+"with error: "+err.Error(),
		)
		return
	}
	for _, hostGroupId := range hostGroupResponse.HostGroupIDs {
		tflog.Info(ctx, hostGroupId)
		state.HostGroupIds = append(state.HostGroupIds, types.StringValue(hostGroupId))
	}

	state.ID = types.StringValue("1")
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
