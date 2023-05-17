// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// StorageGroupResourceModel describes the resource data model.
type StorageGroupResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	StorageGroupID        types.String `tfsdk:"name"`
	SLO                   types.String `tfsdk:"slo"`
	SRP                   types.String `tfsdk:"srp_id"`
	ServiceLevel          types.String `tfsdk:"service_level"`
	Workload              types.String `tfsdk:"workload"`
	SLOCompliance         types.String `tfsdk:"slo_compliance"`
	NumOfVolumes          types.Int64  `tfsdk:"num_of_vols"`
	NumOfChildSGs         types.Int64  `tfsdk:"num_of_child_sgs"`
	NumOfParentSGs        types.Int64  `tfsdk:"num_of_parent_sgs"`
	NumOfMaskingViews     types.Int64  `tfsdk:"num_of_masking_views"`
	NumOfSnapshots        types.Int64  `tfsdk:"num_of_snapshots"`
	NumOfSnapshotPolicies types.Int64  `tfsdk:"num_of_snapshot_policies"`
	CapacityGB            types.Number `tfsdk:"cap_gb"`
	DeviceEmulation       types.String `tfsdk:"device_emulation"`
	Type                  types.String `tfsdk:"type"`
	Unprotected           types.Bool   `tfsdk:"unprotected"`
	ChildStorageGroup     types.List   `tfsdk:"child_storage_group"`
	ParentStorageGroup    types.List   `tfsdk:"parent_storage_group"`
	MaskingView           types.List   `tfsdk:"maskingview"`
	SnapshotPolicies      types.List   `tfsdk:"snapshot_policies"`
	HostIOLimit           types.Object `tfsdk:"host_io_limit"`
	Compression           types.Bool   `tfsdk:"compression"`
	CompressionRatio      types.String `tfsdk:"compression_ratio"`
	CompressionRatioToOne types.Number `tfsdk:"compression_ratio_to_one"`
	VPSavedPercent        types.Number `tfsdk:"vp_saved_percent"`
	Tags                  types.String `tfsdk:"tags"`
	UUID                  types.String `tfsdk:"uuid"`
	UnreducibleDataGB     types.Number `tfsdk:"unreducible_data_gb"`
	VolumeIDs             types.List   `tfsdk:"volume_ids"`
}

// SetHostIOLimitsParam describes the data model for setting host IO limits.
type SetHostIOLimitsParam struct {
	HostIOLimitMBSec    types.String `tfsdk:"host_io_limit_mb_sec"`
	HostIOLimitIOSec    types.String `tfsdk:"host_io_limit_io_sec"`
	DynamicDistribution types.String `tfsdk:"dynamic_distribution"`
}

// StorageGroupDataSourceModel describes the data source data model.
type StorageGroupDataSourceModel struct {
	ID                 types.String                `tfsdk:"id"`
	StorageGroups      []StorageGroupResourceModel `tfsdk:"storage_groups"`
	StorageGroupFilter *sgFilterType               `tfsdk:"filter"`
}

type sgFilterType struct {
	IDs []types.String `tfsdk:"names"`
}