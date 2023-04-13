// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// Volume holds volume schema attribute details.
type Volume struct {
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
	SymmetrixPortKeys     types.List    `tfsdk:"symmetrix_port_keys"`
	RDFGroupIDList        types.List    `tfsdk:"rdf_group_ids"`
}

// StorageGroupID holds information of StorageGroupName.
type StorageGroupID struct {
	StorageGroupName types.String `tfsdk:"storage_group_name"`
}
