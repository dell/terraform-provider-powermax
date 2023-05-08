package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// Host holds host schema attribute details.
type Host struct {
	// ID - defines host ID
	ID types.String `tfsdk:"id"`
	// Name - The name of the host
	Name types.String `tfsdk:"name"`
	// HostFlags - Specifies the flags set for a host
	HostFlags HostFlags `tfsdk:"host_flags"`
	// ConsistentLun - Specifies whether the consistent_lun flag is set or not for a host
	ConsistentLun types.Bool `tfsdk:"consistent_lun"`
	// Initiators - (Set of String) Specifies the list of initiators for a host
	Initiators types.Set `tfsdk:"initiators"`
	// NumOfMaskingViews - Specifies the number of masking views for a host
	NumOfMaskingViews types.Int64 `tfsdk:"numofmaskingviews"`
	// NumOfInitiators - Specifies the number of initiators for a host
	NumOfInitiators types.Int64 `tfsdk:"numofinitiators"`
	// NumOfHostGroups - Specifies the number of host groups for a host
	NumOfHostGroups types.Int64 `tfsdk:"numofhostgroups"`
	// PortFlagsOverride - Specifies whether port flags override is enabled on the host
	PortFlagsOverride types.Bool `tfsdk:"port_flags_override"`
	// Type - Specifies the type of host
	Type types.String `tfsdk:"type"`
	// Maskingview - Specifies the list of masking view for a host
	Maskingview types.List `tfsdk:"maskingview"`
	// PowerpathHosts - Specifies powerpath hosts associated with the host
	PowerpathHosts types.List `tfsdk:"powerpath_hosts"`
	// NumOfPowerpathHosts - Specifies the number of powerpath hosts for a host
	NumOfPowerpathHosts types.Int64 `tfsdk:"numofpowerpathhosts"`
	// BWLimit - Specifies the bandwidth limit for a host
	BWLimit types.Int64 `tfsdk:"bw_limit"`
}

// HostFlags - group of flags used as part of host creation.
type HostFlags struct {
	VolumeSetAddressing HostFlag `tfsdk:"volume_set_addressing"`
	DisableQResetOnUa   HostFlag `tfsdk:"disable_q_reset_on_ua"`
	EnvironSet          HostFlag `tfsdk:"environ_set"`
	AvoidResetBroadcast HostFlag `tfsdk:"avoid_reset_broadcast"`
	Openvms             HostFlag `tfsdk:"openvms"`
	Scsi3               HostFlag `tfsdk:"scsi_3"`
	Spc2ProtocolVersion HostFlag `tfsdk:"spc2_protocol_version"`
	ScsiSupport1        HostFlag `tfsdk:"scsi_support1"`
}

// HostFlag holds overwrite info for individual flag.
type HostFlag struct {
	Enabled  types.Bool `tfsdk:"enabled"`
	Override types.Bool `tfsdk:"override"`
}
