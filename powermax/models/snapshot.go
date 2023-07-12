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

// SnapshotDataSourceModel describes the hostgroup data source model.
type SnapshotDataSourceModel struct {
	ID           types.String          `tfsdk:"id"`
	Snapshots    []SnapshotDetailModal `tfsdk:"snapshots"`
	StorageGroup *filterTypeSnapshot   `tfsdk:"storage_group"`
}

type filterTypeSnapshot struct {
	Name types.String `tfsdk:"name"`
}

// SnapshotDetailModal describes the detail of snapshot data source.
type SnapshotDetailModal struct {
	// The name of the SnapVX snapshot.
	Name types.String `tfsdk:"name"`
	// The number of generation for the snapshot. Using snap IDs instead of generation numbers is preferred.
	Generation types.Int64 `tfsdk:"generation"`
	// The unique snap ID for the snapshot. This snap ID value does not change in the way generations can change when newer snapshots are created or terminated. Using snapset IDs instead of generation numbers is preferred.
	Snapid types.Int64 `tfsdk:"snapid"`
	// The timestamp of the snapshot generation.
	Timestamp types.String `tfsdk:"timestamp"`
	// The timestamp of the snapshot generation in milliseconds since 1970.
	TimestampUtc types.String `tfsdk:"timestamp_utc"`
	// The state of the snapshot generation.
	State types.List `tfsdk:"state"`
	// The number of source volumes in the snapshot generation.
	NumSourceVolumes types.Int64 `tfsdk:"num_source_volumes"`
	// The source volumes of the snapshot generation.
	SourceVolume types.List `tfsdk:"source_volume"`
	// The number of non-gatekeeper storage group volumes.
	NumStorageGroupVolumes types.Int64 `tfsdk:"num_storage_group_volumes"`
	// The number of source tracks that have been overwritten by the host.
	Tracks types.Int64 `tfsdk:"tracks"`
	// The number of tracks uniquely allocated for this snapshots delta. This is an approximate indication of the number of tracks that will be returned to the SRP if this snapshot is terminated.
	NonSharedTracks types.Int64 `tfsdk:"non_shared_tracks"`
	// When the snapshot will expire once it is not linked.
	TimeToLiveExpiryDate types.String `tfsdk:"time_to_live_expiry_date"`
	// When the snapshot will expire once it is not linked.
	SecureExpiryDate types.String `tfsdk:"secure_expiry_date"`
	// Set if this generation secure has expired.
	Expired types.Bool `tfsdk:"expired"`
	// Set if this generation is SnapVX linked.
	Linked types.Bool `tfsdk:"linked"`
	// Set if this generation is restored.
	Restored types.Bool `tfsdk:"restored"`
	// Linked storage group names. Only populated if the generation is linked.
	LinkedStorageGroupNames types.List `tfsdk:"linked_storage_group_names"`
	// Linked storage group and volume information. Only populated if the generation is linked.
	LinkedStorageGroup types.List `tfsdk:"linked_storage_group"`
	// Set if this snapshot is persistent.  Only applicable to policy based snapshots.
	Persistent types.Bool `tfsdk:"persistent"`
}

// SnapshotGenerationSource The source volumes of the snapshot generation.
type SnapshotGenerationSource struct {
	// The name of the SnapVX snapshot generation source volume.
	Name types.String `tfsdk:"name"`
	// The capacity of the snapshot volume in cylinders
	Capacity types.Int64 `tfsdk:"capacity"`
	// The capacity of the snapshot volume in GB
	CapacityGb types.Float64 `tfsdk:"capacity_gb"`
}

// LinkedSnapshot Linked snapshot.
type LinkedSnapshot struct {
	// The storage group name.
	Name types.String `tfsdk:"name"`
	// The source volumes name.
	SourceVolumeName types.String `tfsdk:"source_volume_name"`
	// The linked volumes name.
	LinkedVolumeName types.String `tfsdk:"linked_volume_name"`
	// Number of tracks.
	Tracks types.Int64 `tfsdk:"tracks"`
	// Size of the tracks.
	TrackSize types.Int64 `tfsdk:"track_size"`
	// Percentage of tracks copied.
	PercentageCopied types.Int64 `tfsdk:"percentage_copied"`
	// The average timestamp of all linked volumes that are linked.
	LinkedCreationTimestamp types.String `tfsdk:"linked_creation_timestamp"`
	// When the snapshot link has been fully defined.
	Defined types.Bool `tfsdk:"defined"`
	// When the snapshot link is being defined.
	BackgroundDefineInProgress types.Bool `tfsdk:"background_define_in_progress"`
}
