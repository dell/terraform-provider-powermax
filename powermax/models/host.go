// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// HostFlags - group of flags used as part of host creation
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

// HostFlag holds overwrite info for individual flag
type HostFlag struct {
	Enabled  types.Bool `tfsdk:"enabled"`
	Override types.Bool `tfsdk:"override"`
}
