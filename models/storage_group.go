package models

import (
	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StorageGroup holds storage group schema attribute details
type StorageGroup struct {
	// ID - defines storage group ID
	ID types.String `tfsdk:"id"`
	// Name - The name of the storage group
	Name types.String `tfsdk:"name"`
	// SRPID - The storage resource pool ID associated with the storage group
	SRPID types.String `tfsdk:"srpid"`
	// ServiceLevel - The service level associated with the storage group
	ServiceLevel types.String `tfsdk:"service_level"`
	// SnapshotPolicies - (Set of Snapshot Policy Objects) The snapshot policies associated with the storage group
	SnapshotPolicies types.Set `tfsdk:"snapshot_policies"`
	// EnableCompression - States whether compression is enabled on storage group
	EnableCompression types.Bool `tfsdk:"enable_compression"`
	// Workload - The workload associated with the storage group
	Workload types.String `tfsdk:"workload"`
	// VolumeIDs - (Set of String) The IDs of the volumes associated with this storagegroup
	VolumeIDs types.Set `tfsdk:"volume_ids"`
	// NumOfVols - The number of volumes associated with the storage group
	NumOfVols types.Int64 `tfsdk:"numofvols"`
	// NumOfChildSgs - The number of child storage groups associated with the storage group
	NumOfChildSgs types.Int64 `tfsdk:"numofchildsgs"`
	// NumOfParentSgs - The number of parent storage groups associated with the storage group
	NumOfParentSgs types.Int64 `tfsdk:"numofparentsgs"`
	// NumOfMaskingViews - The number of masking views associated with the storage group
	NumOfMaskingViews types.Int64 `tfsdk:"numofmaskingviews"`
	// NumOfSnapshots - The number of snapshots associated with the storage group
	NumOfSnapshots types.Int64 `tfsdk:"numofsnapshots"`
	// CapGB - The capacity of the storage group in GB
	CapGB types.Number `tfsdk:"cap_gb"`
	// SloCompliance - The service level compliance status of the storage group
	SloCompliance types.String `tfsdk:"slo_compliance"`
	// DeviceEmulation - The emulation of the volumes in the storage group
	DeviceEmulation types.String `tfsdk:"device_emulation"`
	// Type - The storage group type
	Type types.String `tfsdk:"type"`
	//Unprotected - States whether the storage group is protected
	Unprotected types.Bool `tfsdk:"unprotected"`
	// CompressionRatio - Compression ratio of the storage group
	CompressionRatio types.String `tfsdk:"compression_ratio"`
	// CompressionRatioToOne - Compression ratio numeric value of the storage group
	CompressionRatioToOne types.Number `tfsdk:"compression_ratio_to_one"`
	// VPsavedPercent - VP saved percentage figure
	VPsavedPercent types.Number `tfsdk:"vp_saved_percent"`
	MaskingView    types.List   `tfsdk:"maskingview"`
	// ChildStorageGroup types.ListType `tfsdk:"childstoragegroup"`
	// UUID - Storage Group UUID
	UUID types.String `tfsdk:"uuid"`
	// UnreducibleDataGB - The amount of unreducible data in Gb.
	UnreducibleDataGB types.Number `tfsdk:"unreducible_data_gb"`
	// ParentStorageGroup types.ListType `tfsdk:"parentstoragegroup"`
	// HostIOLimits - Host IO limit of the storage group
	HostIOLimits types.Map `tfsdk:"host_io_limits"`
}

// SnapshotPolicy holds snapshot policy details associated with storage group
type SnapshotPolicy struct {
	PolicyName types.String `tfsdk:"policy_name"`
	IsActive   types.Bool   `tfsdk:"is_active"`
}

// StorageGroupWithData is a wrapper storageGroupModel which has the storageGroup returned by powermax along with VolumeIDs and snapshot policies associated with this storageGroup
type StorageGroupWithData struct {
	PmaxStorageGroup *pmaxTypes.StorageGroup
	VolumeIDs        []string
	SnapshotPolicies []pmaxTypes.StorageGroupSnapshotPolicy
}
