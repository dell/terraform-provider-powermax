// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// HostGroup holds hostgroup schema attribute details
type HostGroupModal struct {
	// ID - defines host ID
	ID types.String `tfsdk:"id"`
	// Name - The name of the hostgroup
	Name types.String `tfsdk:"name"`
	// HostFlags - Specifies the flags set for a hostgroup
	HostFlags HostFlags `tfsdk:"host_flags"`
	// ConsistentLun - Specifies whether the consistent_lun flag is set or not for a hostgroup
	ConsistentLun types.Bool `tfsdk:"consistent_lun"`
	// HostIDs - Specifies the host IDs associated with the hostgroup
	HostIDs types.Set `tfsdk:"host_ids"`
	// NumOfMaskingViews - Specifies the number of masking views for a hostgroup
	NumOfMaskingViews types.Int64 `tfsdk:"numofmaskingviews"`
	// NumOfInitiators - Specifies the number of initiators for a hostgroup
	NumOfInitiators types.Int64 `tfsdk:"numofinitiators"`
	// NumOfHosts - Specifies the number of hosts in the hostgroup
	NumOfHosts types.Int64 `tfsdk:"numofhosts"`
	// PortFlagsOverride - Specifies whether port flags override is enabled on the hostgroup
	PortFlagsOverride types.Bool `tfsdk:"port_flags_override"`
	// Type - Specifies the type of hostgroup
	Type types.String `tfsdk:"type"`
	// Maskingview - Specifies the list of maskingviews for a hostgroup
	Maskingviews types.List `tfsdk:"maskingviews"`
}
