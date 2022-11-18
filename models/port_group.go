package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// PortGroup holds portgroup schema attribute details
type PortGroup struct {
	// ID - defines portgroup ID
	ID types.String `tfsdk:"id"`
	// Name - The name of the portgroup
	Name types.String `tfsdk:"name"`
	// Ports - The ports associated with the portgroup
	Ports []PortKey `tfsdk:"ports"`
	// Protocol - The portgroup protocol
	Protocol types.String `tfsdk:"protocol"`
	// NumOfPorts - The number of ports associated with the portgroup
	NumOfPorts types.Int64 `tfsdk:"numofports"`
	// NumOfMaskingViews - The number of masking views associated with the portgroup
	NumOfMaskingViews types.Int64 `tfsdk:"numofmaskingviews"`
	// Type - The type of the portgroup
	Type types.String `tfsdk:"type"`
	// Maskingview - The list of masking views associated with the portgroup
	Maskingview types.List `tfsdk:"maskingview"`
	// TestID - The test ID of the portgroup
	TestID types.String `tfsdk:"test_id"`
}

// PortKey holds DirectorID and PortKey
type PortKey struct {
	DirectorID types.String `tfsdk:"director_id"`
	PortID     types.String `tfsdk:"port_id"`
}
