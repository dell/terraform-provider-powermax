package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// MaskingView  holds maskingview schema attribute details.
type MaskingView struct {
	// ID - defines maskingview ID
	ID types.String `tfsdk:"id"`
	// Name - The name of the maskingview
	Name types.String `tfsdk:"name"`
	// StorageGroupID - ID of the storage group associated with maskingview
	StorageGroupID types.String `tfsdk:"storage_group_id"`
	// PortGroupID - ID of the port group associated with maskingview
	PortGroupID types.String `tfsdk:"port_group_id"`
	// HostID - ID of the host associated with maskingview
	HostID types.String `tfsdk:"host_id"`
	// HostGroupID - ID of the hostgroup associated with maskingview
	HostGroupID types.String `tfsdk:"host_group_id"`
}
