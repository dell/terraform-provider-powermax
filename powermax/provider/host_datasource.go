// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package provider

import (
	"context"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &hostDataSource{}
	_ datasource.DataSourceWithConfigure = &hostDataSource{}
)

// NewHostDataSource returns the host data source object
func NewHostDataSource() datasource.DataSource {
	return &hostDataSource{}
}

type hostDataSource struct {
	client *client.Client
}

func (d *hostDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (d *hostDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Host DataSource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier of the host instance.",
				MarkdownDescription: "Unique identifier of the host instance.",
				Optional:            true,
				Computed:            true,
			},
			"host_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Description:         "The list of host ids.",
				MarkdownDescription: "The list of host ids.",
			},
		},
	}
}

func (d *hostDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *hostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.HostDataSourceModel
	var err error

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Read all the hosts
	ids, err := d.client.PmaxClient.GetHostList(ctx, d.client.SymmetrixID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read PowerMax Hosts",
			err.Error(),
		)
		return
	}

	hostList := []attr.Value{}
	for _, id := range ids.HostIDs {
		hostList = append(hostList, types.StringValue(string(id)))
	}
	state.HostIDs, _ = types.ListValue(types.StringType, hostList)
	state.ID = types.StringValue("1")

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
