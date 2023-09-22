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

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &HostDataSource{}
	_ datasource.DataSourceWithConfigure = &HostDataSource{}
)

// NewHostDataSource returns the host data source object.
func NewHostDataSource() datasource.DataSource {
	return &HostDataSource{}
}

// HostDataSource configures client for host data source.
type HostDataSource struct {
	client *client.Client
}

// Metadata returns the metadata for host data source.
func (d *HostDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

// Schema returns the schema for host data source.
func (d *HostDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	hostFlagNestedAttr := map[string]schema.Attribute{
		"override": schema.BoolAttribute{
			Optional: true,
			Computed: true,
		},
		"enabled": schema.BoolAttribute{
			Optional: true,
			Computed: true,
		},
	}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source for reading Hosts in PowerMax array. PowerMax hosts systems are storage hosts that use storage system LUN resources. A logical unit number (LUN) is an identifier that is used for labeling and designating subsystems of physical or virtual storage",
		Description:         "Data source for reading Hosts in PowerMax array. PowerMax hosts systems are storage hosts that use storage system LUN resources. A logical unit number (LUN) is an identifier that is used for labeling and designating subsystems of physical or virtual storage",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx),
			"id": schema.StringAttribute{
				Description:         "Unique identifier of the host instance.",
				MarkdownDescription: "Unique identifier of the host instance.",
				Computed:            true,
			},
			"hosts": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of host attributes",
				MarkdownDescription: "List of host attributes",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "The ID of the host.",
							MarkdownDescription: "The ID of the host.",
						},
						"name": schema.StringAttribute{
							Required:            true,
							Description:         "The name of the host.",
							MarkdownDescription: "The name of the host.",
						},
						"num_of_masking_views": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of masking views associated with the host.",
							MarkdownDescription: "The number of masking views associated with the host.",
						},
						"num_of_initiators": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of initiators associated with the host.",
							MarkdownDescription: "The number of initiators associated with the host.",
						},
						"num_of_host_groups": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of hostgroups associated with the host.",
							MarkdownDescription: "The number of hostgroups associated with the host.",
						},
						"port_flags_override": schema.BoolAttribute{
							Computed:            true,
							Description:         "States whether port flags override is enabled on the host.",
							MarkdownDescription: "States whether port flags override is enabled on the host.",
						},
						"consistent_lun": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "It enables the rejection of any masking operation involving this host that would result in inconsistent LUN values.",
							MarkdownDescription: "It enables the rejection of any masking operation involving this host that would result in inconsistent LUN values.",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							Description:         "Specifies the type of host.",
							MarkdownDescription: "Specifies the type of host.",
						},
						"initiator": schema.ListAttribute{
							ElementType:         types.StringType,
							Required:            true,
							Description:         "The initiators associated with the host.",
							MarkdownDescription: "The initiators associated with the host.",
						},
						"hostgroup": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							Description:         "The host group associated with the host.",
							MarkdownDescription: "The host group associated with the host.",
						},
						"maskingview": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							Description:         "The masking views associated with the host.",
							MarkdownDescription: "The masking views associated with the host.",
						},
						"powerpathhosts": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							Description:         "The powerpath hosts associated with the host.",
							MarkdownDescription: "The powerpath hosts associated with the host.",
						},
						"numofpowerpathhosts": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of powerpath hosts associated with the host.",
							MarkdownDescription: "The number of powerpath hosts associated with the host.",
						},
						"bw_limit": schema.Int64Attribute{
							Computed:            true,
							Description:         "Specifies the bandwidth limit for a host.",
							MarkdownDescription: "Specifies the bandwidth limit for a host.",
						},
						"host_flags": schema.SingleNestedAttribute{
							Required: true,
							Attributes: map[string]schema.Attribute{
								"volume_set_addressing": schema.SingleNestedAttribute{
									Optional:            true,
									Computed:            true,
									Attributes:          hostFlagNestedAttr,
									Description:         "It enables the volume set addressing mode.",
									MarkdownDescription: "It enables the volume set addressing mode.",
								},
								"disable_q_reset_on_ua": schema.SingleNestedAttribute{
									Optional:            true,
									Computed:            true,
									Attributes:          hostFlagNestedAttr,
									Description:         "It is used for hosts that do not expect the queue to be flushed on a 0629 sense.",
									MarkdownDescription: "It is used for hosts that do not expect the queue to be flushed on a 0629 sense.",
								},
								"environ_set": schema.SingleNestedAttribute{
									Optional:            true,
									Computed:            true,
									Attributes:          hostFlagNestedAttr,
									Description:         "It enables the environmental error reporting by the storage system to the host on the specific port.",
									MarkdownDescription: "It enables the environmental error reporting by the storage system to the host on the specific port.",
								},
								"openvms": schema.SingleNestedAttribute{
									Optional:            true,
									Computed:            true,
									Attributes:          hostFlagNestedAttr,
									Description:         "This attribute enables an Open VMS fibre connection.",
									MarkdownDescription: "This attribute enables an Open VMS fibre connection.",
								},
								"avoid_reset_broadcast": schema.SingleNestedAttribute{
									Optional:            true,
									Computed:            true,
									Attributes:          hostFlagNestedAttr,
									Description:         "It enables a SCSI bus reset to only occur to the port that received the reset.",
									MarkdownDescription: "It enables a SCSI bus reset to only occur to the port that received the reset.",
								},
								"scsi_3": schema.SingleNestedAttribute{
									Optional:            true,
									Computed:            true,
									Attributes:          hostFlagNestedAttr,
									Description:         "Alters the inquiry data to report that the storage system supports the SCSI-3 protocol.",
									MarkdownDescription: "Alters the inquiry data to report that the storage system supports the SCSI-3 protocol.",
								},
								"spc2_protocol_version": schema.SingleNestedAttribute{
									Optional:            true,
									Computed:            true,
									Attributes:          hostFlagNestedAttr,
									Description:         "When setting this flag, the port must be offline.",
									MarkdownDescription: "When setting this flag, the port must be offline.",
								},
								"scsi_support1": schema.SingleNestedAttribute{
									Optional:            true,
									Computed:            true,
									Attributes:          hostFlagNestedAttr,
									Description:         "This attribute provides a stricter compliance with SCSI standards.",
									MarkdownDescription: "This attribute provides a stricter compliance with SCSI standards.",
								},
							},
							Description:         "Flags set for the host. When host_flags = {} then default flags will be considered.",
							MarkdownDescription: "Flags set for the host. When host_flags = {} then default flags will be considered.",
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

// Configure configure client.
func (d *HostDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read host data source.
func (d *HostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.HostsDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := helper.SetupTimeoutReadDatasource(ctx, resp, state.Timeout)
	if resp.Diagnostics.HasError() {
		return
	}

	defer cancel()

	var hostIds []string
	// Get host IDs from config or query all if not specified
	if state.HostFilter == nil || len(state.HostFilter.Names) == 0 {
		// Read all the hosts
		hostIDList, _, err := helper.GetHostList(ctx, *d.client)

		if err != nil {
			errStr := ""
			msgStr := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Error reading host ids", msgStr)
			return
		}
		hostIds = hostIDList.HostId
	} else {
		// get ids from filter and assign to hostIds
		for _, ids := range state.HostFilter.Names {
			hostIds = append(hostIds, ids.ValueString())
		}
	}

	// iterate Host IDs and Get Host with each id
	for _, id := range hostIds {
		getHostReq := d.client.PmaxOpenapiClient.SLOProvisioningApi.GetHost(ctx, d.client.SymmetrixID, id)
		hostResponse, _, err := getHostReq.Execute()
		if err != nil || hostResponse == nil {
			// Check to see if timeout was hit
			helper.ExceedTimeoutErrorCheck(err, resp)
			if resp.Diagnostics.HasError() {
				return
			}
			errStr := ""
			msgStr := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Error reading host with id", msgStr)
			continue
		}
		var host models.HostModel
		tflog.Debug(ctx, "Updating host state")
		helper.UpdateHostState(&host, []string{}, hostResponse)
		state.Hosts = append(state.Hosts, host)
	}

	state.ID = types.StringValue("1")

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
