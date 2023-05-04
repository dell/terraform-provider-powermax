// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// MaskingViewResourceModel describes the resource data model.
type MaskingViewResourceModel struct {
	Name           types.String `tfsdk:"name"`
	ID             types.String `tfsdk:"id"`
	StorageGroupID types.String `tfsdk:"storage_group_id"`
	HostID         types.String `tfsdk:"host_id"`
	HostGroupID    types.String `tfsdk:"host_group_id"`
	PortGroupID    types.String `tfsdk:"port_group_id"`
}

// MaskingViewDataSourceModel describes the data source data model.
type MaskingViewDataSourceModel struct {
	MaskingViews []MaskingViewModel `tfsdk:"masking_views"`
	ID           types.String       `tfsdk:"id"`
	//filter
	MaskingViewFilter *MaskingViewFilterType `tfsdk:"filter"`
}

type MaskingViewModel struct {
	MaskingViewName types.String  `tfsdk:"masking_view_name"`
	HostID          types.String  `tfsdk:"host_id"`
	HostGroupID     types.String  `tfsdk:"host_group_id"`
	PortGroupID     types.String  `tfsdk:"port_group_id"`
	StorageGroupID  types.String  `tfsdk:"storage_group_id"`
	CapacityGB      types.Float64 `tfsdk:"capacity_gb"`
	Volumes         types.List    `tfsdk:"volumes"`
	Initiators      types.List    `tfsdk:"initiators"`
	Ports           types.List    `tfsdk:"ports"`
}

type MaskingViewFilterType struct {
	Names []types.String `tfsdk:"names"`
}
