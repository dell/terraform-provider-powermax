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

// PortGroup holds portgroup schema attribute details.
type PortGroup struct {
	// ID - defines portgroup ID
	ID types.String `tfsdk:"id"`
	// Name - The name of the portgroup
	Name types.String `tfsdk:"name"`
	// Ports - (Set of Ports) The ports associated with the portgroup
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
}

// PortKey holds DirectorID and PortKey.
type PortKey struct {
	DirectorID types.String `tfsdk:"director_id"`
	PortID     types.String `tfsdk:"port_id"`
}

// PortgroupsDataSourceModel describes the data source data model.
type PortgroupsDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	PortGroups []PortGroup  `tfsdk:"port_groups"`
	//filter
	PgFilter *portGroupFilterType `tfsdk:"filter"`
}

type portGroupFilterType struct {
	Names []types.String `tfsdk:"names"`
	// Type - The type of the portgroup
	Type types.String `tfsdk:"type"`
}
