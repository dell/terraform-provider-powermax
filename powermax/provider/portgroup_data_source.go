// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package provider

import (
	"context"
	"fmt"
	"terraform-provider-powermax/client"

	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PortgroupDataSource{}
var _ datasource.DataSourceWithConfigure = &PortgroupDataSource{}

func NewPortgroupDataSource() datasource.DataSource {
	return &PortgroupDataSource{}
}

// PortgroupDataSource defines the data source implementation.
type PortgroupDataSource struct {
	client *client.Client
}

// PortgroupDataSourceModel describes the data source data model.
type portgroupsDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	Portgroups []string     `tfsdk:"portgroups_id"`
	// Type - The type of the portgroup
	Type types.String `tfsdk:"type"`
}

func (d *PortgroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_portgroups"
}

func (d *PortgroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source for reading PortGroups in PowerMax array.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Required:            true,
				Description:         "The Type of the portgroup.",
				MarkdownDescription: "The Type of the portgroup.",
			},
			"portgroups_id": schema.ListAttribute{
				Description: "List of Host Group Ids",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *PortgroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pmaxclient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = pmaxclient
}

func (d *PortgroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state portgroupsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	//Read the portgroup based on portgroup type and if nothing is mentioned, then it returns all the port groups
	portGroups, err := d.client.PmaxClient.GetPortGroupList(ctx, d.client.SymmetrixID, state.Type.ValueString())

	//check if there is any error while getting the port group
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read PowerMax Port Groups",
			err.Error(),
		)
		return
	}
	state.Portgroups = getPortgroupListData(portGroups)
	state.ID = types.StringValue("1")

	tflog.Trace(ctx, "read PortGroup data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getPortgroupListData(pgs *pmaxTypes.PortGroupList) []string {
	var pgIDList []string
	for _, elem := range pgs.PortGroupIDs {
		pgIDList = append(pgIDList, elem)
	}
	return pgIDList
}
