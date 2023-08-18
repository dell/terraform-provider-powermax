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

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PortgroupDataSource{}
var _ datasource.DataSourceWithConfigure = &PortgroupDataSource{}

// PortgroupDataSource defines the data source implementation.
type PortgroupDataSource struct {
	client *client.Client
}

// NewPortgroupDataSource is a helper function to simplify the provider implementation.
func NewPortgroupDataSource() datasource.DataSource {
	return &PortgroupDataSource{}
}

// Metadata returns the metadata for the data source.
func (d *PortgroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_portgroups"
}

// Schema returns the schema for the data source.
func (d *PortgroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source for reading PortGroups in PowerMax array.",
		Description:         "Data source for reading PortGroups in PowerMax array.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier",
				Computed:    true,
			},
			"port_groups": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of port group attributes",
				MarkdownDescription: "List of port group attributes",
				NestedObject: schema.NestedAttributeObject{
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
						"name": schema.StringAttribute{
							Required:            true,
							Description:         "The name of the portgroup.",
							MarkdownDescription: "The name of the portgroup.",
						},
						"ports": schema.ListNestedAttribute{
							Required: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"director_id": schema.StringAttribute{
										Required: true,
									},
									"port_id": schema.StringAttribute{
										Required: true,
									},
								},
							},
							Description:         "The list of ports associated with the portgroup.",
							MarkdownDescription: "The list of ports associated with the portgroup.",
						},
						"protocol": schema.StringAttribute{
							Required:            true,
							Description:         "The portgroup protocol.",
							MarkdownDescription: "The portgroup protocol.",
						},
						"numofports": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of ports associated with the portgroup.",
							MarkdownDescription: "The number of ports associated with the portgroup.",
						},
						"numofmaskingviews": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of masking views associated with the portgroup.",
							MarkdownDescription: "The number of masking views associated with the portgroup.",
						},
						"maskingview": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							Description:         "The masking views associated with the portgroup.",
							MarkdownDescription: "The masking views associated with the portgroup.",
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
					"type": schema.StringAttribute{
						Optional:            true,
						Description:         "The Type of the portgroup.",
						MarkdownDescription: "The Type of the portgroup.",
					},
				},
			},
		},
	}
}

// Configure configures the data source.
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
	var pgPlan models.PortgroupsDataSourceModel
	var pgState models.PortgroupsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &pgPlan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var pgNames []string

	portGroupIDList, _, err := helper.GetPortGroupList(ctx, *d.client, pgPlan)
	if err != nil {
		errStr := ""
		msgStr := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Unable to Read PowerMax Port Groups", msgStr,
		)
	}
	// Get portgroup IDs from config or query all if not specified
	if pgPlan.PgFilter == nil || len(pgPlan.PgFilter.Names) == 0 {
		pgNames = portGroupIDList.GetPortGroupId()
	} else {
		// get ids from portGroups and assign to pgNames
		for _, pg := range pgPlan.PgFilter.Names {
			for _, pgFromList := range portGroupIDList.GetPortGroupId() {
				if pg.ValueString() == pgFromList {
					pgNames = append(pgNames, pg.ValueString())
				}
			}
		}
		if len(pgNames) != len(pgPlan.PgFilter.Names) {
			resp.Diagnostics.AddError("Invalid name(s) provided.", "Name of already created portgroup must be provided.")
			return
		}
	}
	var portGroups []models.PortGroup

	// iterate Portgroup IDs and GetPortGroup with each id
	for _, elemid := range pgNames {
		pgResponse, _, err := helper.ReadPortgroupByID(ctx, *d.client, elemid)
		if err != nil || pgResponse == nil {
			errStr := fmt.Sprintf("Error reading port group with id %s", elemid)
			msgStr := helper.GetErrorString(err, "")
			resp.Diagnostics.AddError(errStr, msgStr)
			return
		}
		var pg models.PortGroup
		// Copy fields from the provider client data into the Terraform state
		helper.UpdatePGState(&pg, &pg, pgResponse)
		portGroups = append(portGroups, pg)
	}
	pgState.PortGroups = portGroups
	//check if there is any error while getting the port group
	pgState.ID = types.StringValue("1")
	pgState.PgFilter = pgPlan.PgFilter

	tflog.Trace(ctx, "read PortGroup data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &pgState)...)
}
