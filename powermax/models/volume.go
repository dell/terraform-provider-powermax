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

// VolumeResource holds volume schema attribute details.
type VolumeResource struct {
	ID                    types.String  `tfsdk:"id"`
	StorageGroupName      types.String  `tfsdk:"sg_name"`
	VolumeIdentifier      types.String  `tfsdk:"vol_name"`
	Size                  types.Number  `tfsdk:"size"`
	CapUnit               types.String  `tfsdk:"cap_unit"`
	StorageGroups         types.List    `tfsdk:"storage_groups"`
	Type                  types.String  `tfsdk:"type"`
	Emulation             types.String  `tfsdk:"emulation"`
	SSID                  types.String  `tfsdk:"ssid"`
	AllocatedPercent      types.Int64   `tfsdk:"allocated_percent"`
	Status                types.String  `tfsdk:"status"`
	Reserved              types.Bool    `tfsdk:"reserved"`
	Pinned                types.Bool    `tfsdk:"pinned"`
	WWN                   types.String  `tfsdk:"wwn"`
	Encapsulated          types.Bool    `tfsdk:"encapsulated"`
	NumberOfStorageGroups types.Int64   `tfsdk:"num_of_storage_groups"`
	NumberOfFrontEndPaths types.Int64   `tfsdk:"num_of_front_end_paths"`
	SnapSource            types.Bool    `tfsdk:"snapvx_source"`
	SnapTarget            types.Bool    `tfsdk:"snapvx_target"`
	HasEffectiveWWN       types.Bool    `tfsdk:"has_effective_wwn"`
	EffectiveWWN          types.String  `tfsdk:"effective_wwn"`
	EncapsulatedWWN       types.String  `tfsdk:"encapsulated_wwn"`
	MobilityIDEnabled     types.Bool    `tfsdk:"mobility_id_enabled"`
	UnreducibleDataGB     types.Float64 `tfsdk:"unreducible_data_gb"`
	NGUID                 types.String  `tfsdk:"nguid"`
	OracleInstanceName    types.String  `tfsdk:"oracle_instance_name"`
	SymmetrixPortKey      types.List    `tfsdk:"symmetrix_port_key"`
	RDFGroupIDList        types.List    `tfsdk:"rdf_group_ids"`
}

// VolumeDatasourceFilter holds volume datasource filter schema attribute details.
type VolumeDatasourceFilter struct {
	StorageGroupID       types.String `tfsdk:"storage_group_name"`
	EncapsulatedWwn      types.String `tfsdk:"encapsulated_wwn"`
	WWN                  types.String `tfsdk:"wwn"`
	Symmlun              types.String `tfsdk:"symmlun"`
	Status               types.String `tfsdk:"status"`
	PhysicalName         types.String `tfsdk:"physical_name"`
	VolumeIdentifier     types.String `tfsdk:"volume_identifier"`
	AllocatedPercent     types.String `tfsdk:"allocated_percent"`
	CapTb                types.String `tfsdk:"cap_tb"`
	CapGb                types.String `tfsdk:"cap_gb"`
	CapMb                types.String `tfsdk:"cap_mb"`
	CapCYL               types.String `tfsdk:"cap_cyl"`
	NumOfStorageGroups   types.String `tfsdk:"num_of_storage_groups"`
	NumOfMaskingViews    types.String `tfsdk:"num_of_masking_views"`
	NumOfFrontEndPaths   types.String `tfsdk:"num_of_front_end_paths"`
	VirtualVolumes       types.Bool   `tfsdk:"virtual_volumes"`
	PrivateVolumes       types.Bool   `tfsdk:"private_volumes"`
	AvailableThinVolumes types.Bool   `tfsdk:"available_thin_volumes"`
	Tdev                 types.Bool   `tfsdk:"tdev"`
	ThinBcv              types.Bool   `tfsdk:"thin_bcv"`
	Vdev                 types.Bool   `tfsdk:"vdev"`
	Gatekeeper           types.Bool   `tfsdk:"gatekeeper"`
	DataVolume           types.Bool   `tfsdk:"data_volume"`
	Dld                  types.Bool   `tfsdk:"dld"`
	Drv                  types.Bool   `tfsdk:"drv"`
	Mapped               types.Bool   `tfsdk:"mapped"`
	BoundTdev            types.Bool   `tfsdk:"bound_tdev"`
	Reserved             types.Bool   `tfsdk:"reserved"`
	Pinned               types.Bool   `tfsdk:"pinned"`
	Encapsulated         types.Bool   `tfsdk:"encapsulated"`
	Associated           types.Bool   `tfsdk:"associated"`
	Emulation            types.String `tfsdk:"emulation"`
	SplitName            types.String `tfsdk:"split_name"`
	CuImageNum           types.String `tfsdk:"cu_image_num"`
	CuImageSsid          types.String `tfsdk:"cu_image_ssid"`
	RdfGroupNumber       types.String `tfsdk:"rdf_group_number"`
	HasEffectiveWwn      types.Bool   `tfsdk:"has_effective_wwn"`
	EffectiveWwn         types.String `tfsdk:"effective_wwn"`
	Type                 types.String `tfsdk:"type"`
	OracleInstanceName   types.String `tfsdk:"oracle_instance_name"`
	MobilityIDEnabled    types.Bool   `tfsdk:"mobility_id_enabled"`
	UnreducibleDataGb    types.String `tfsdk:"unreducible_data_gb"`
	Nguid                types.String `tfsdk:"nguid"`
}

// VolumeDatasource holds volume datasource schema attribute details.
type VolumeDatasource struct {
	// placeholder for acc testing
	ID           types.String             `tfsdk:"id"`
	Volumes      []VolumeDatasourceEntity `tfsdk:"volumes"`
	VolumeFilter *VolumeDatasourceFilter  `tfsdk:"filter"`
}

// VolumeDatasourceEntity holds volume datasource entity schema attribute details.
type VolumeDatasourceEntity struct {
	VolumeID              types.String `tfsdk:"id"`
	VolumeIdentifier      types.String `tfsdk:"volume_identifier"`
	StorageGroups         types.List   `tfsdk:"storage_groups"`
	Type                  types.String `tfsdk:"type"`
	Emulation             types.String `tfsdk:"emulation"`
	SSID                  types.String `tfsdk:"ssid"`
	AllocatedPercent      types.Int64  `tfsdk:"allocated_percent"`
	PhysicalName          types.String `tfsdk:"physical_name"`
	Status                types.String `tfsdk:"status"`
	Reserved              types.Bool   `tfsdk:"reserved"`
	Pinned                types.Bool   `tfsdk:"pinned"`
	WWN                   types.String `tfsdk:"wwn"`
	Encapsulated          types.Bool   `tfsdk:"encapsulated"`
	NumberOfStorageGroups types.Int64  `tfsdk:"num_of_storage_groups"`
	NumberOfFrontEndPaths types.Int64  `tfsdk:"num_of_front_end_paths"`
	SnapSource            types.Bool   `tfsdk:"snapvx_source"`
	SnapTarget            types.Bool   `tfsdk:"snapvx_target"`
	HasEffectiveWWN       types.Bool   `tfsdk:"has_effective_wwn"`
	EffectiveWWN          types.String `tfsdk:"effective_wwn"`
	EncapsulatedWWN       types.String `tfsdk:"encapsulated_wwn"`
	MobilityIDEnabled     types.Bool   `tfsdk:"mobility_id_enabled"`
	UnreducibleDataGB     types.Number `tfsdk:"unreducible_data_gb"`
	NGUID                 types.String `tfsdk:"nguid"`
	OracleInstanceName    types.String `tfsdk:"oracle_instance_name"`
	SymmetrixPortKey      types.List   `tfsdk:"symmetrix_port_key"`
	RDFGroupIDList        types.List   `tfsdk:"rdf_group_ids"`
	CapacityGB            types.Number `tfsdk:"cap_gb"`
	FloatCapacityMB       types.Number `tfsdk:"cap_mb"`
	CapacityCYL           types.Int64  `tfsdk:"cap_cyl"`
}

// StorageGroupName holds information of StorageGroupName, ParentStorageGroupName.
type StorageGroupName struct {
	StorageGroupName       types.String `tfsdk:"storage_group_name"`
	ParentStorageGroupName types.String `tfsdk:"parent_storage_group_name"`
}
