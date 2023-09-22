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

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &portDataSource{}
	_ datasource.DataSourceWithConfigure = &portDataSource{}
)

// NewPortDataSource is a helper function to simplify the provider implementation.
func NewPortDataSource() datasource.DataSource {
	return &portDataSource{}
}

// hostGroupDataSource is the data source implementation.
type portDataSource struct {
	client *client.Client
}

func (d *portDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port"
}

func (d *portDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source for reading ports in PowerMax array. A port typically refers to the interface that allows for connections between the PowerMax system and other devices.",
		Description:         "Data source for reading ports in PowerMax array. A port typically refers to the interface that allows for connections between the PowerMax system and other devices.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx),
			"id": schema.StringAttribute{
				Description: "Identifier",
				Computed:    true,
			},
			"port_details": schema.ListNestedAttribute{
				Description:         "List of Ports",
				MarkdownDescription: "List of Ports",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"director_id": schema.StringAttribute{
							Description:         "Id of the director",
							MarkdownDescription: "Id of the director",
							Computed:            true,
						},
						"port_id": schema.StringAttribute{
							Description:         "Id of the port",
							MarkdownDescription: "Id of the port",
							Computed:            true,
						},
						"port_status": schema.StringAttribute{
							Description:         "Port Status",
							MarkdownDescription: "Port Status",
							Computed:            true,
						},
						"director_status": schema.StringAttribute{
							Description:         "Director Status",
							MarkdownDescription: "Director Status",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							Description:         "Port Type",
							MarkdownDescription: "Port Type",
							Computed:            true,
						},
						"num_of_cores": schema.Int64Attribute{
							Description:         "Total number of cors",
							MarkdownDescription: "Total number of cors",
							Computed:            true,
						},
						"identifier": schema.StringAttribute{
							Description:         "Port identifier",
							MarkdownDescription: "Port identifier",
							Computed:            true,
						},
						"negotiated_speed": schema.StringAttribute{
							Description:         "Negotiated speed",
							MarkdownDescription: "Negotiated speed",
							Computed:            true,
						},
						"mac_address": schema.StringAttribute{
							Description:         "Mac Address",
							MarkdownDescription: "Mac Address",
							Computed:            true,
						},
						"num_of_port_groups": schema.Int64Attribute{
							Description:         "Total number of port groups",
							MarkdownDescription: "Total number of port groups",
							Computed:            true,
						},
						"num_of_masking_views": schema.Int64Attribute{
							Description:         "Total number of masking views",
							MarkdownDescription: "Total number of masking views",
							Computed:            true,
						},
						"num_of_mapped_vols": schema.Int64Attribute{
							Description:         "Total number of volumes",
							MarkdownDescription: "Total number of volumes",
							Computed:            true,
						},
						"vcm_state": schema.StringAttribute{
							Description:         "VMC State",
							MarkdownDescription: "VMC State",
							Computed:            true,
						},
						"aclx": schema.BoolAttribute{
							Description:         "Has aclx",
							MarkdownDescription: "Has aclx",
							Computed:            true,
						},
						"common_serial_number": schema.BoolAttribute{
							Description:         "Common Serial Number",
							MarkdownDescription: "Common Serial Number",
							Computed:            true,
						},
						"unique_wwn": schema.BoolAttribute{
							Description:         "Unique WWN",
							MarkdownDescription: "Unique WWN",
							Computed:            true,
						},
						"init_point_to_point": schema.BoolAttribute{
							Description:         "Init point to point",
							MarkdownDescription: "Init point to point",
							Computed:            true,
						},
						"volume_set_addressing": schema.BoolAttribute{
							Description:         "Volume Set Addressing",
							MarkdownDescription: "Volume Set Addressing",
							Computed:            true,
						},
						"vnx_attached": schema.BoolAttribute{
							Description:         "VNX Attached",
							MarkdownDescription: "VNX Attached",
							Computed:            true,
						},
						"avoid_reset_broadcast": schema.BoolAttribute{
							Description:         "Avoid reset brodcast",
							MarkdownDescription: "Avoid reset brodcast",
							Computed:            true,
						},
						"negotiate_reset": schema.BoolAttribute{
							Description:         "Negotiate reset",
							MarkdownDescription: "Negotiate reset",
							Computed:            true,
						},
						"enable_auto_negotiate": schema.BoolAttribute{
							Description:         "Enable Auto Negotiate",
							MarkdownDescription: "Enable Auto Negotiate",
							Computed:            true,
						},
						"environ_set": schema.BoolAttribute{
							Description:         "Environ Set",
							MarkdownDescription: "Environ Set",
							Computed:            true,
						},
						"disable_q_reset_on_ua": schema.BoolAttribute{
							Description:         "Disable Q reset on ua",
							MarkdownDescription: "Disable Q reset on ua",
							Computed:            true,
						},
						"soft_reset": schema.BoolAttribute{
							Description:         "Soft reset",
							MarkdownDescription: "Soft reset",
							Computed:            true,
						},
						"scsi_3": schema.BoolAttribute{
							Description:         "SCSI 3",
							MarkdownDescription: "SCSI 3",
							Computed:            true,
						},
						"scsi_support1": schema.BoolAttribute{
							Description:         "SCSI support 1",
							MarkdownDescription: "SCSI support 1",
							Computed:            true,
						},
						"no_participating": schema.BoolAttribute{
							Description:         "No Participating",
							MarkdownDescription: "No Participating",
							Computed:            true,
						},
						"spc2_protocol_version": schema.BoolAttribute{
							Description:         "SPC2 Protocol Version",
							MarkdownDescription: "SPC2 Protocol Version",
							Computed:            true,
						},
						"hp_3000_mode": schema.BoolAttribute{
							Description:         "HP 3000 mode",
							MarkdownDescription: "HP 3000 mode",
							Computed:            true,
						},
						"sunapee": schema.BoolAttribute{
							Description:         "Sunapee",
							MarkdownDescription: "Sunapee",
							Computed:            true,
						},
						"siemens": schema.BoolAttribute{
							Description:         "Siemens",
							MarkdownDescription: "Siemens",
							Computed:            true,
						},
						"portgroup": schema.ListAttribute{
							Description:         "Portgroup",
							MarkdownDescription: "Portgroup",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"maskingview": schema.ListAttribute{
							Description:         "Masking Views",
							MarkdownDescription: "Masking Views",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"rx_power_level_mw": schema.Float64Attribute{
							Description:         "RX Power Level MW",
							MarkdownDescription: "RX Power Level MW",
							Computed:            true,
						},
						"tx_power_level_mw": schema.Float64Attribute{
							Description:         "TX Power Level MW",
							MarkdownDescription: "TX Power Level MW",
							Computed:            true,
						},
						"power_levels_last_sampled_date_milliseconds": schema.Int64Attribute{
							Description:         "Power Levels Last Sampled Date in Milliseconds",
							MarkdownDescription: "Power Levels Last Sampled Date in Milliseconds",
							Computed:            true,
						},
						"port_interface": schema.StringAttribute{
							Description:         "Port Interface",
							MarkdownDescription: "Port Interface",
							Computed:            true,
						},
						"num_of_hypers": schema.Int64Attribute{
							Description:         "TX Power Level MW",
							MarkdownDescription: "TX Power Level MW",
							Computed:            true,
						},
						"rdf_ra_group_attributes_farpoint": schema.BoolAttribute{
							Description:         "RDF RA group attributes farpoint",
							MarkdownDescription: "RDF RA group attributes farpoint",
							Computed:            true,
						},
						"prevent_automatic_rdf_link_recovery": schema.StringAttribute{
							Description:         "Prevent automatic rdf link recovery",
							MarkdownDescription: "Prevent automatic rdf link recovery",
							Computed:            true,
						},
						"prevent_ra_online_on_power_up": schema.StringAttribute{
							Description:         "Prevent RA Online on Power Up",
							MarkdownDescription: "Prevent RA Online on Power Up",
							Computed:            true,
						},
						"rdf_software_compression_supported": schema.StringAttribute{
							Description:         "RDF Software Compression Suppored",
							MarkdownDescription: "RDF Software Compression Suppored",
							Computed:            true,
						},
						"rdf_software_compression": schema.StringAttribute{
							Description:         "RDF Software Compression",
							MarkdownDescription: "RDF Software Compression",
							Computed:            true,
						},
						"rdf_hardware_compression_supported": schema.StringAttribute{
							Description:         "RDF Hardware Compression Supported",
							MarkdownDescription: "RDF Hardware Compression Supported",
							Computed:            true,
						},
						"rdf_hardware_compression": schema.StringAttribute{
							Description:         "RDF Hardware Compression",
							MarkdownDescription: "RDF Hardware Compression",
							Computed:            true,
						},
						"ipv4_address": schema.StringAttribute{
							Description:         "Ipv4 Address",
							MarkdownDescription: "Ipv4 Address",
							Computed:            true,
						},
						"ipv6_address": schema.StringAttribute{
							Description:         "Ipv6 Address",
							MarkdownDescription: "Ipv6 Address",
							Computed:            true,
						},
						"ipv6_prefix": schema.StringAttribute{
							Description:         "Ipv6 Prefix",
							MarkdownDescription: "Ipv6 Prefix",
							Computed:            true,
						},
						"ipv4_default_gateway": schema.StringAttribute{
							Description:         "Ipv4 Default Gateway",
							MarkdownDescription: "Ipv4 Default Gateway",
							Computed:            true,
						},
						"ipv4_domain_name": schema.StringAttribute{
							Description:         "Ipv4 Domain Name",
							MarkdownDescription: "Ipv4 Domain Name",
							Computed:            true,
						},
						"ipv4_netmask": schema.StringAttribute{
							Description:         "Ipv4 Netmask",
							MarkdownDescription: "Ipv4 Netmask",
							Computed:            true,
						},
						"max_speed": schema.StringAttribute{
							Description:         "Max Speed",
							MarkdownDescription: "Max Speed",
							Computed:            true,
						},
						"wwn_node": schema.StringAttribute{
							Description:         "WWN Node",
							MarkdownDescription: "WWN Node",
							Computed:            true,
						},
						"iscsi_target": schema.BoolAttribute{
							Description:         "iScsi Target",
							MarkdownDescription: "iScsi Target",
							Computed:            true,
						},
						"iscsi_endpoint": schema.BoolAttribute{
							Description:         "iScsi Endpoint",
							MarkdownDescription: "iScsi Endpoint",
							Computed:            true,
						},
						"nvmetcp_endpoint": schema.BoolAttribute{
							Description:         "NVME over TCP Endpoint",
							MarkdownDescription: "NVME over TCP Endpoint",
							Computed:            true,
						},
						"network_id": schema.Int64Attribute{
							Description:         "Network Id",
							MarkdownDescription: "Network Id",
							Computed:            true,
						},
						"tcp_port": schema.Int64Attribute{
							Description:         "TPC Port",
							MarkdownDescription: "TPC Port",
							Computed:            true,
						},
						"ip_addresses": schema.ListAttribute{
							Description:         "Ip Addresses",
							MarkdownDescription: "Ip Addresses",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"enabled_protocol": schema.ListAttribute{
							Description:         "Enabled Protocol",
							MarkdownDescription: "Enabled Protocol",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"capable_protocol": schema.ListAttribute{
							Description:         "Capable Protocol",
							MarkdownDescription: "Capable Protocol",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"z_hyperlink_port": schema.BoolAttribute{
							Description:         "Z Hyperlink Port",
							MarkdownDescription: "Z Hyperlink Port",
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"port_ids": schema.SetAttribute{
						Description:         "A set of port ids to filter on, should be look like the following ['directorId:portId']",
						MarkdownDescription: "A set of port ids to filter on, should be look like the following ['directorId:portId']",
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
			},
		},
	}
}

func (d *portDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read.
func (d *portDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Attempting to read ports")
	var state models.PortDataSourceModel
	var plan models.PortDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := helper.SetupTimeoutReadDatasource(ctx, resp, plan.Timeout)
	if resp.Diagnostics.HasError() {
		return
	}

	defer cancel()

	portIds, err := helper.FilterPortIds(ctx, &state, &plan, *d.client)
	if err != nil {
		errStr := constants.ReadPortDetailErrorMsg + "with error:"
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the list of ports",
			message,
		)
		return
	}
	for _, val := range portIds {
		port, _, err := helper.GetPort(ctx, *d.client, val.DirectorId, val.PortId)
		if err != nil {
			// Check to see if timeout was hit
			helper.ExceedTimeoutErrorCheck(err, resp)
			if resp.Diagnostics.HasError() {
				return
			}
			errStr := constants.ReadPortDetailErrorMsg + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError(
				"Error getting the details of port: "+val.DirectorId+":"+val.PortId,
				message,
			)
			return
		}
		model, err := helper.PortDetailMapper(ctx, port)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Unknown Error",
				constants.ReadPortDetailErrorMsg+"with error: "+err.Error(),
			)
			return
		}
		state.PortDetails = append(state.PortDetails, model)
	}
	state.ID = types.StringValue("port-datasource")
	state.Timeout = plan.Timeout
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
