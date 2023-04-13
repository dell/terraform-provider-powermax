// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// MaskingViewResourceModel describes the resource data model.
type MaskingViewResourceModel struct {
	ID             types.String `tfsdk:"id"`
	StorageGroupID types.String `tfsdk:"storage_group_id"`
	HostID         types.String `tfsdk:"host_id"`
	HostGroupID    types.String `tfsdk:"host_group_id"`
	PortGroupID    types.String `tfsdk:"port_group_id"`
}
