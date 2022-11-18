package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Volume holds volume schema attribute details.
type Volume struct {
	ID                 types.String  `tfsdk:"id"`
	Name               types.String  `tfsdk:"name"`
	Size               types.Number  `tfsdk:"size"`
	CapUnit            types.String  `tfsdk:"cap_unit"`
	SGName             types.String  `tfsdk:"sg_name"`
	Type               types.String  `tfsdk:"type"`
	Emulation          types.String  `tfsdk:"emulation"`
	SSID               types.String  `tfsdk:"ssid"`
	AllocatedPercent   types.Int64   `tfsdk:"allocated_percent"`
	Status             types.String  `tfsdk:"status"`
	Reserved           types.Bool    `tfsdk:"reserved"`
	Pinned             types.Bool    `tfsdk:"pinned"`
	WWN                types.String  `tfsdk:"wwn"`
	Encapsulated       types.Bool    `tfsdk:"encapsulated"`
	NumOfStorageGroups types.Int64   `tfsdk:"num_of_storage_groups"`
	NumOfFrontEndPaths types.Int64   `tfsdk:"num_of_front_end_paths"`
	StorageGroupIDs    types.List    `tfsdk:"storagegroup_ids"`
	SymmetrixPortKeys  types.List    `tfsdk:"symmetrix_port_keys"`
	RDFGroupIDs        types.List    `tfsdk:"rdf_group_ids"`
	SnapSource         types.Bool    `tfsdk:"snap_source"`
	SnapTarget         types.Bool    `tfsdk:"snap_target"`
	HasEffectiveWWN    types.Bool    `tfsdk:"has_effective_wwn"`
	EffectiveWWN       types.String  `tfsdk:"effective_wwn"`
	EncapsulatedWWN    types.String  `tfsdk:"encapsulated_wwn"`
	OracleInstanceName types.String  `tfsdk:"oracle_instance_name"`
	EnableMobilityID   types.Bool    `tfsdk:"enable_mobility_id"`
	UnreducibleDataGB  types.Float64 `tfsdk:"unreducible_data_gb"`
	NGUID              types.String  `tfsdk:"nguid"`
}

// SymmetrixPortKey holds director ID and port ID.
type SymmetrixPortKey struct {
	DirectorID types.String `tfsdk:"director_id"`
	PortID     types.String `tfsdk:"port_id"`
}

// RDFGroupID holds information of RDFGroupNumber and label.
type RDFGroupID struct {
	RDFGroupNumber types.Int64  `tfsdk:"rdf_group_number"`
	Label          types.String `tfsdk:"label"`
}
