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

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StorageGroupResourceModel describes the resource data model.
type StorageGroupResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	StorageGroupID        types.String `tfsdk:"name"`
	Slo                   types.String `tfsdk:"slo"`
	Srp                   types.String `tfsdk:"srp_id"`
	ServiceLevel          types.String `tfsdk:"service_level"`
	Workload              types.String `tfsdk:"workload"`
	SloCompliance         types.String `tfsdk:"slo_compliance"`
	NumOfVols             types.Int64  `tfsdk:"num_of_vols"`
	NumOfChildSgs         types.Int64  `tfsdk:"num_of_child_sgs"`
	NumOfParentSgs        types.Int64  `tfsdk:"num_of_parent_sgs"`
	NumOfMaskingViews     types.Int64  `tfsdk:"num_of_masking_views"`
	NumOfSnapshots        types.Int64  `tfsdk:"num_of_snapshots"`
	NumOfSnapshotPolicies types.Int64  `tfsdk:"num_of_snapshot_policies"`
	CapGb                 types.Number `tfsdk:"cap_gb"`
	DeviceEmulation       types.String `tfsdk:"device_emulation"`
	Type                  types.String `tfsdk:"type"`
	Unprotected           types.Bool   `tfsdk:"unprotected"`
	ChildStorageGroup     types.List   `tfsdk:"child_storage_group"`
	ParentStorageGroup    types.List   `tfsdk:"parent_storage_group"`
	Maskingview           types.List   `tfsdk:"maskingview"`
	SnapshotPolicies      types.List   `tfsdk:"snapshot_policies"`
	HostIOLimit           types.Object `tfsdk:"host_io_limit"`
	Compression           types.Bool   `tfsdk:"compression"`
	CompressionRatio      types.String `tfsdk:"compression_ratio"`
	CompressionRatioToOne types.Number `tfsdk:"compression_ratio_to_one"`
	VpSavedPercent        types.Number `tfsdk:"vp_saved_percent"`
	Tags                  types.String `tfsdk:"tags"`
	UUID                  types.String `tfsdk:"uuid"`
	UnreducibleDataGb     types.Number `tfsdk:"unreducible_data_gb"`
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
	Timeout            timeouts.Value              `tfsdk:"timeouts"`
	StorageGroupFilter *sgFilterType               `tfsdk:"filter"`
}

type sgFilterType struct {
	IDs []types.String `tfsdk:"names"`
}
