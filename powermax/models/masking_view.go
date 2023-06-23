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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// MaskingViewResourceModel describes the resource data model.
type MaskingViewResourceModel struct {
	Name           types.String `tfsdk:"name"`
	ID             types.String `tfsdk:"id"`
	StorageGroupId types.String `tfsdk:"storage_group_id"`
	HostId         types.String `tfsdk:"host_id"`
	HostGroupId    types.String `tfsdk:"host_group_id"`
	PortGroupId    types.String `tfsdk:"port_group_id"`
}

// MaskingViewDataSourceModel describes the data source data model.
type MaskingViewDataSourceModel struct {
	MaskingViews []MaskingViewModel `tfsdk:"masking_views"`
	ID           types.String       `tfsdk:"id"`
	//filter
	MaskingViewFilter *MaskingViewFilterType `tfsdk:"filter"`
}

// MaskingViewModel holds masking view data source schema attribute details.
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

// MaskingViewFilterType holds filter attribute for masking view.
type MaskingViewFilterType struct {
	Names []types.String `tfsdk:"names"`
}
