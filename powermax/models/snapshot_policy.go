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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SnapshotPolicyModel Snapshot Policy datasource Model structure.
type SnapshotPolicyModel struct {
	// The name of the snapshot policy on this System
	SnapshotPolicyName types.String `tfsdk:"snapshot_policy_name"`
	// The number of snapshots that will be taken before the oldest ones are no longer required.Max value is 1024.
	SnapshotCount types.Int64 `tfsdk:"snapshot_count"`
	// The number of minutes between each policy execution.
	IntervalMinutes types.Int64 `tfsdk:"interval_minutes"`
	// The number of minutes after 00:00 on Monday morning that the policy will execute
	OffsetMinutes types.Int64 `tfsdk:"offset_minutes"`
	// The name of the cloud provider associated with this policy. Only applies to cloud policies.
	ProviderName types.String `tfsdk:"provider_name"`
	// The number of days that snapshots will be retained in the cloud for. Only applies to cloud policies.
	RetentionDays types.Int64 `tfsdk:"retention_days"`
	// Set if the snapshot policy has been suspended
	Suspended types.Bool `tfsdk:"suspended"`
	// Set if the snapshot policy creates secure snapshots
	Secure types.Bool `tfsdk:"secure"`
	// The last time that the snapshot policy was run
	LastTimeUsed types.String `tfsdk:"last_time_used"`
	// The total number of storage groups that this snapshot policy is associated with
	StorageGroupCount types.Int64 `tfsdk:"storage_group_count"`
	// The threshold of good snapshots which are not failed/bad for compliance to change from normal to warning.
	ComplianceCountWarning types.Int64 `tfsdk:"compliance_count_warning"`
	// The threshold of good snapshots which are not failed/bad for compliance to change from warning to critical.
	ComplianceCountCritical types.Int64 `tfsdk:"compliance_count_critical"`
	// The type of Snapshots that are created with the policy, local or cloud.
	Type types.String `tfsdk:"type"`
}

// SnapshotPolicyDataSourceModel describes the snapshot policy data source model.
type SnapshotPolicyDataSourceModel struct {
	ID               types.String          `tfsdk:"id"`
	SnapshotPolicies []SnapshotPolicyModel `tfsdk:"snapshot_policies"`
	//filter
	SnapshotPolicyFilter *SnapshotPolicyFilterType `tfsdk:"filter"`
}

// SnapshotPolicyFilterType describes the filter data model.
type SnapshotPolicyFilterType struct {
	Names []types.String `tfsdk:"names"`
}

// SnapshotPolicyResource structure.
type SnapshotPolicyResource struct {
	ID types.String `tfsdk:"id"`
	// The name of the snapshot policy on this System
	SnapshotPolicyName types.String `tfsdk:"snapshot_policy_name"`
	// The number of snapshots that will be taken before the oldest ones are no longer required.Max value is 1024.
	SnapshotCount types.Int64 `tfsdk:"snapshot_count"`
	// The number of minutes between each policy execution.
	IntervalMinutes types.Int64 `tfsdk:"interval_minutes"`
	// The number of minutes after 00:00 on Monday morning that the policy will execute
	OffsetMinutes types.Int64 `tfsdk:"offset_minutes"`
	// The name of the cloud provider associated with this policy. Only applies to cloud policies.
	ProviderName types.String `tfsdk:"provider_name"`
	// The number of days that snapshots will be retained in the cloud for. Only applies to cloud policies.
	RetentionDays types.Int64 `tfsdk:"retention_days"`
	// Set if the snapshot policy has been suspended
	Suspended types.Bool `tfsdk:"suspended"`
	// Set if the snapshot policy creates secure snapshots
	Secure types.Bool `tfsdk:"secure"`
	// The last time that the snapshot policy was run
	LastTimeUsed types.String `tfsdk:"last_time_used"`
	// Time inetrval between each policy execution. Valid values:10 Minutes, 12 Minutes,
	// 15 Minutes,20 Minutes, 30 Minutes, 1 Hour, 2 Hours, 3 Hours, 4 Hours, 6 Hours, 8 Hours, 12 Hours, 1 Day, 7 Days
	Interval types.String `tfsdk:"interval"`
	// The total number of storage groups that this snapshot policy is associated with
	StorageGroupCount types.Int64 `tfsdk:"storage_group_count"`
	// The threshold of good snapshots which are not failed/bad for compliance to change from normal to warning.
	ComplianceCountWarning types.Int64 `tfsdk:"compliance_count_warning"`
	// The threshold of good snapshots which are not failed/bad for compliance to change from warning to critical.
	ComplianceCountCritical types.Int64 `tfsdk:"compliance_count_critical"`
	// The type of Snapshots that are created with the policy, local or cloud.
	Type types.String `tfsdk:"type"`
}
