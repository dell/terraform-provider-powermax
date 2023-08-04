/*
Copyright (c) 2022-2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// PortDataSourceModel describes the port data source model.
type PortDataSourceModel struct {
	ID          types.String      `tfsdk:"id"`
	PortDetails []PortDetailModal `tfsdk:"port_details"`
	PortFilter  *portFilterType   `tfsdk:"filter"`
}

type portFilterType struct {
	IDs []types.String `tfsdk:"port_ids"`
}

// PortDetailModal the details of the port resource.
type PortDetailModal struct {
	// Director ID
	DirectorID types.String `tfsdk:"director_id"`
	// Port ID
	PortID types.String `tfsdk:"port_id"`
	// port_status
	PortStatus types.String `tfsdk:"port_status"`
	// director_status
	DirectorStatus types.String `tfsdk:"director_status"`
	// type
	Type types.String `tfsdk:"type"`
	// num_of_cores
	NumOfCores types.Int64 `tfsdk:"num_of_cores"`
	// identifier
	Identifier types.String `tfsdk:"identifier"`
	// negotiated_speed
	NegotiatedSpeed types.String `tfsdk:"negotiated_speed"`
	// mac_address
	MacAddress types.String `tfsdk:"mac_address"`
	// num_of_port_groups
	NumOfPortGroups types.Int64 `tfsdk:"num_of_port_groups"`
	// num_of_masking_views
	NumOfMaskingViews types.Int64 `tfsdk:"num_of_masking_views"`
	// num_of_mapped_vols
	NumOfMappedVols types.Int64 `tfsdk:"num_of_mapped_vols"`
	// vcm_state
	VcmState types.String `tfsdk:"vcm_state"`
	// aclx
	Aclx types.Bool `tfsdk:"aclx"`
	// common_serial_number
	CommonSerialNumber types.Bool `tfsdk:"common_serial_number"`
	// unique_wwn
	UniqueWwn types.Bool `tfsdk:"unique_wwn"`
	// init_point_to_point
	InitPointToPoint types.Bool `tfsdk:"init_point_to_point"`
	// volume_set_addressing
	VolumeSetAddressing types.Bool `tfsdk:"volume_set_addressing"`
	// vnx_attached
	VnxAttached types.Bool `tfsdk:"vnx_attached"`
	// avoid_reset_broadcast
	AvoidResetBroadcast types.Bool `tfsdk:"avoid_reset_broadcast"`
	// negotiate_reset
	NegotiateReset types.Bool `tfsdk:"negotiate_reset"`
	// enable_auto_negotiate
	EnableAutoNegotiate types.Bool `tfsdk:"enable_auto_negotiate"`
	// environ_set
	EnvironSet types.Bool `tfsdk:"environ_set"`
	// disable_q_reset_on_ua
	DisableQResetOnUa types.Bool `tfsdk:"disable_q_reset_on_ua"`
	// soft_reset
	SoftReset types.Bool `tfsdk:"soft_reset"`
	// scsi_3
	Scsi3 types.Bool `tfsdk:"scsi_3"`
	// scsi_support1
	ScsiSupport1 types.Bool `tfsdk:"scsi_support1"`
	// no_participating
	NoParticipating types.Bool `tfsdk:"no_participating"`
	// spc2_protocol_version
	Spc2ProtocolVersion types.Bool `tfsdk:"spc2_protocol_version"`
	// hp_3000_mode
	Hp3000Mode types.Bool `tfsdk:"hp_3000_mode"`
	// sunapee
	Sunapee types.Bool `tfsdk:"sunapee"`
	// siemens
	Siemens types.Bool `tfsdk:"siemens"`
	// portgroup
	Portgroup types.List `tfsdk:"portgroup"`
	// maskingview
	Maskingview types.List `tfsdk:"maskingview"`
	// rx_power_level_mw
	RxPowerLevelMw types.Float64 `tfsdk:"rx_power_level_mw"`
	// tx_power_level_mw
	TxPowerLevelMw types.Float64 `tfsdk:"tx_power_level_mw"`
	// power_levels_last_sampled_date_milliseconds
	PowerLevelsLastSampledDateMilliseconds types.Int64 `tfsdk:"power_levels_last_sampled_date_milliseconds"`
	// port_interface
	PortInterface types.String `tfsdk:"port_interface"`
	// num_of_hypers
	NumOfHypers types.Int64 `tfsdk:"num_of_hypers"`
	// rdf_ra_group_attributes_farpoint
	RdfRaGroupAttributesFarpoint types.Bool `tfsdk:"rdf_ra_group_attributes_farpoint"`
	// prevent_automatic_rdf_link_recovery
	PreventAutomaticRdfLinkRecovery types.String `tfsdk:"prevent_automatic_rdf_link_recovery"`
	// prevent_ra_online_on_power_up
	PreventRaOnlineOnPowerUp types.String `tfsdk:"prevent_ra_online_on_power_up"`
	// rdf_software_compression_supported
	RdfSoftwareCompressionSupported types.String `tfsdk:"rdf_software_compression_supported"`
	// rdf_software_compression
	RdfSoftwareCompression types.String `tfsdk:"rdf_software_compression"`
	// rdf_hardware_compression_supported
	RdfHardwareCompressionSupported types.String `tfsdk:"rdf_hardware_compression_supported"`
	// rdf_hardware_compression
	RdfHardwareCompression types.String `tfsdk:"rdf_hardware_compression"`
	// ipv4_address
	Ipv4Address types.String `tfsdk:"ipv4_address"`
	// ipv6_address
	Ipv6Address types.String `tfsdk:"ipv6_address"`
	// ipv6_prefix
	Ipv6Prefix types.String `tfsdk:"ipv6_prefix"`
	// ipv4_default_gateway
	Ipv4DefaultGateway types.String `tfsdk:"ipv4_default_gateway"`
	// ipv4_domain_name
	Ipv4DomainName types.String `tfsdk:"ipv4_domain_name"`
	// ipv4_netmask
	Ipv4Netmask types.String `tfsdk:"ipv4_netmask"`
	// max_speed
	MaxSpeed types.String `tfsdk:"max_speed"`
	// wwn_node
	WwnNode types.String `tfsdk:"wwn_node"`
	// iscsi_target
	IscsiTarget types.Bool `tfsdk:"iscsi_target"`
	// iscsi_endpoint
	IscsiEndpoint types.Bool `tfsdk:"iscsi_endpoint"`
	// nvmetcp_endpoint
	NvmetcpEndpoint types.Bool `tfsdk:"nvmetcp_endpoint"`
	// network_id
	NetworkID types.Int64 `tfsdk:"network_id"`
	// tcp_port
	TPCPort types.Int64 `tfsdk:"tcp_port"`
	// ip_addresses
	IPAddresses types.List `tfsdk:"ip_addresses"`
	// enabled_Protocols   Enumeration values: * **None** * **RDF_FC** * **RDF_GigE** * **Host_FICON** * **SCSI_FC** * **iSCSI** * **NVMe/FC** * **NVMe/TCP**
	EnabledProtocol types.List `tfsdk:"enabled_protocol"`
	// capable_Protocols   Enumeration values: * **None** * **RDF_FC** * **RDF_GigE** * **Host_FICON** * **SCSI_FC** * **iSCSI** * **NVMe/FC** * **NVMe/TCP**
	CapableProtocol types.List `tfsdk:"capable_protocol"`
	// z_hyperlink_port
	ZHyperlinkPort types.Bool `tfsdk:"z_hyperlink_port"`
}
