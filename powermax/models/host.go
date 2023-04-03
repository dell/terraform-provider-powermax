// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// HostModel describes the resource data model.
type HostModel struct {
	HostID             types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	NumberMaskingViews types.Int64  `tfsdk:"num_of_masking_views"`
	NumberInitiators   types.Int64  `tfsdk:"num_of_initiators"`
	NumberHostGroups   types.Int64  `tfsdk:"num_of_host_groups"`
	PortFlagsOverride  types.Bool   `tfsdk:"port_flags_override"`
	ConsistentLun      types.Bool   `tfsdk:"consistent_lun"`
	HostType           types.String `tfsdk:"type"`
	Initiators         types.List   `tfsdk:"initiator"`
	MaskingviewIDs     types.List   `tfsdk:"maskingview"`
	PowerPathHosts     types.List   `tfsdk:"powerpathhosts"`
	NumPowerPathHosts  types.Int64  `tfsdk:"numofpowerpathhosts"`
	BWLimit            types.Int64  `tfsdk:"bw_limit"`
	// HostFlags - Specifies the flags set for a host
	HostFlags HostFlags `tfsdk:"host_flags"`
}

// HostFlags - group of flags used as part of host creation.
type HostFlags struct {
	VolumeSetAddressing HostFlag `tfsdk:"volume_set_addressing"`
	DisableQResetOnUA   HostFlag `tfsdk:"disable_q_reset_on_ua"`
	EnvironSet          HostFlag `tfsdk:"environ_set"`
	AvoidResetBroadcast HostFlag `tfsdk:"avoid_reset_broadcast"`
	OpenVMS             HostFlag `tfsdk:"openvms"`
	SCSI3               HostFlag `tfsdk:"scsi_3"`
	Spc2ProtocolVersion HostFlag `tfsdk:"spc2_protocol_version"`
	SCSISupport1        HostFlag `tfsdk:"scsi_support1"`
}

// HostFlag holds overwrite info for individual flag.
type HostFlag struct {
	Enabled  types.Bool `tfsdk:"enabled"`
	Override types.Bool `tfsdk:"override"`
}

// HostsDataSourceModel describes the data source data model.
type HostsDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Hosts []HostModel  `tfsdk:"hosts"`

	//filter
	HostFilter *HostFilterType `tfsdk:"filter"`
}

// HostFilterType describes the filter data model.
type HostFilterType struct {
	Names []types.String `tfsdk:"names"`
}
