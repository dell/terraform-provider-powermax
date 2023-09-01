---
# Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
#
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://mozilla.org/MPL/2.0/
#
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

title: "powermax_snapshot resource"
linkTitle: "powermax_snapshot"
page_title: "powermax_snapshot Resource - terraform-provider-powermax"
subcategory: ""
description: |-
  Resource for managing Snapshots in PowerMax array. Supported Update (name, secure, timetolive, link, restore). PowerMax Snaphots is a local replication solution that is designed to nondisruptively create point-in-time copies (snapshots) of critical data.
---

# powermax_snapshot (Resource)

Resource for managing Snapshots in PowerMax array. Supported Update (name, secure, time_to_live, link, restore). PowerMax Snaphots is a local replication solution that is designed to nondisruptively create point-in-time copies (snapshots) of critical data.


## Example Usage

```terraform
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

# Available actions: Create, Update (name, secure, time_to_live, link, restore), Delete and Import an existing snapshot from the PowerMax Array.
# After `terraform apply` of this example file it will create a new snapshot with the name set in `name` attribute on the PowerMax for the storage group set in the storage group `name` attribute
# NOTE: that all of the PowerMax `snapshot_actions` are only available during modify of the snapshot after it has been created.

# PowerMax Snaphots is a local replication solution that is designed to nondisruptively create point-in-time copies (snapshots) of critical data.
resource "powermax_snapshot" "snapshot_1" {

  # Attributes which are able to be modified after create (secure, time_to_live, link, restore)

  # Required The storage group that the snapshot will be taken upon 
  storage_group {
    name = "test_sg"
  }

  # The different actions that can be taken on a snapshot ONLY after the snapshot has already been created
  snapshot_actions {

    # Required, name of new snapshot
    # Only alphanumeric characters, underscores ( _ ), and hyphens (-) are allowed.
    name = "new_test_snapshot"

    # Optional this is only available for modify after the resource is created
    # Secure will make it so the snapshot can not be delete until the time in the `secure` flag runs out.
    # Even if the user changes enabled to false after the secure flag is set
    # The snapshot will remmain secure and cannot be deleted
    secure = {
      # Set to true to it will enable this action
      enable = true
      # If set to true the time_to_live is in hours, otherwise it is in days
      time_in_hours = true
      # How long before the snapshot will no longer be secure
      secure = 3
    }

    # Optional this is only available for modify after the resource is created
    # When this flag is set it will delete the snapshot after the time_to_live expires
    time_to_live = {
      # Set to true to it will enable this action
      enable = true
      # If set to true the time_to_live is in hours, otherwise it is in days
      time_in_hours = false
      # How long before the snapshot is deleted
      time_to_live = 1
    }

    # Optional this is only available for modify after the resource is created
    # When this flag is set it will link the tareget_storage_group with the storage group above
    link = {
      # When set to false it will look through the already linked snapshots, if the one in the target group is currently linked it will remove that storage group
      # When set to true it will look through the already linked snapshots, if the one in the target group is not currently liked it will like that storage group
      enable = false
      # The storage group which will be linked
      target_storage_group = "test_target_sg"
      # The target storage group will not have compression turned on when the SRP is compression capable
      no_compression = true
      # Acknowledges that the data will be propagated to the remote mirror of the SRDF volume. This is not allowed on a nocopy linked target.
      remote = false
      # Sets the link copy mode to perform background copy to the target volume(s).
      copy = false
    }

    # Optional this is only available for modify after the resource is created
    # Will attempt to restore the snapshot to a pervious state
    restore = {
      # When enabled it will attempt to restore the snapshot
      enable = true
      # Acknowledges that the data will be propagated to the remote mirror of the RDF device. This is not allowed on a nocopy link target.
      remote = false
    }
  }
}

# After the execution of above resource block, a PowerMax snapshot has been created at PowerMax array.
# For more information about the newly created resource use the `terraform show` command to review the current state
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `generation` (Number) Number of generation for the snapshot
- `linked_storage_group` (Attributes List) Linked storage group and volume information. Only populated if the generation is linked (see [below for nested schema](#nestedatt--linked_storage_group))
- `non_shared_tracks` (Number) The number of tracks uniquely allocated for this snapshots delta. This is an approximate indication of the number of tracks that will be returned to the SRP if this snapshot is terminated.
- `num_source_volumes` (Number) The number of source volumes in the snapshot generation
- `persistent` (Boolean) Set if this snapshot is persistent.  Only applicable to policy based snapshots
- `secure_expiry_date` (String) When the snapshot will expire once it is not linked
- `snapid` (Number) Unique Snap ID for Snapshot
- `snapshot_actions` (Block, Optional) (see [below for nested schema](#nestedblock--snapshot_actions))
- `storage_group` (Block, Optional) (see [below for nested schema](#nestedblock--storage_group))
- `time_to_live_expiry_date` (String) When the snapshot will expire once it is not linked
- `tracks` (Number) The number of source tracks that have been overwritten by the host

### Read-Only

- `expired` (Boolean) Set if this generation secure has expired
- `id` (String) Identifier
- `linked` (Boolean) Set if this generation is SnapVX linked
- `linked_storage_group_names` (List of String) Linked storage group names. Only populated if the generation is linked
- `name` (String) Name of a snapshot
- `num_storage_group_volumes` (Number) The number of non-gatekeeper storage group volumes
- `restored` (Boolean) Set if this generation is SnapVX linked
- `source_volume` (Attributes List) The source volumes of the snapshot generation (see [below for nested schema](#nestedatt--source_volume))
- `state` (List of String) The state of the snapshot generation
- `timestamp` (String) Timestamp of the snapshot generation
- `timestamp_utc` (String) The timestamp of the snapshot generation in milliseconds since 1970

<a id="nestedatt--linked_storage_group"></a>
### Nested Schema for `linked_storage_group`

Optional:

- `background_define_in_progress` (Boolean) When the snapshot link is being defined
- `defined` (Boolean) When the snapshot link has been fully defined

Read-Only:

- `linked_creation_timestamp` (String) The average timestamp of all linked volumes that are linked
- `linked_volume_name` (String) The linked volumes name
- `name` (String) The storage group name
- `percentage_copied` (Number) Percentage of tracks copied
- `source_volume_name` (String) The source volumes name
- `track_size` (Number) Size of the tracks.
- `tracks` (Number) Number of tracks


<a id="nestedblock--snapshot_actions"></a>
### Nested Schema for `snapshot_actions`

Required:

- `name` (String) Name of the snapshot. (Update Supported)

Optional:

- `both_sides` (Boolean) both_sides defaults to false. Performs the operation on both locally and remotely associated snapshots.
- `link` (Attributes) Link a snapshot generation. (Update Supported) (see [below for nested schema](#nestedatt--snapshot_actions--link))
- `remote` (Boolean) remote defaults to false. If true, The target storage group will not have compression turned on when the SRP is compression capable.
- `restore` (Attributes) Restore a snapshot generation. (Update Supported) (see [below for nested schema](#nestedatt--snapshot_actions--restore))
- `secure` (Attributes) Set the number of days or hours for a snapshot generation to be secure before it auto-terminates (provided it is not linked). (Update Supported) (see [below for nested schema](#nestedatt--snapshot_actions--secure))
- `time_to_live` (Attributes) Set the number of days or hours for a snapshot generation before it auto-terminates (provided it is not linked). (Update Supported) (see [below for nested schema](#nestedatt--snapshot_actions--time_to_live))

<a id="nestedatt--snapshot_actions--link"></a>
### Nested Schema for `snapshot_actions.link`

Optional:

- `copy` (Boolean) copy defaults to false. If true Sets the link copy mode to perform background copy to the target volume(s).
- `enable` (Boolean) enable defaults to false. Flag to enable link on the snapshot
- `no_compression` (Boolean) no_compression defaults to false. If true, The target storage group will not have compression turned on when the SRP is compression capable. Option Used in Action Link
- `remote` (Boolean) remote defaults to false. If true, The target storage group will not have compression turned on when the SRP is compression capable. Option Used in Action Link
- `target_storage_group` (String) The target storage group to link the snapshot too


<a id="nestedatt--snapshot_actions--restore"></a>
### Nested Schema for `snapshot_actions.restore`

Optional:

- `enable` (Boolean) enable defaults to false. Flag to enable restore on the snapshot
- `remote` (Boolean) remote defaults to false. If true, The target storage group will not have compression turned on when the SRP is compression capable. Option Used in Action Link


<a id="nestedatt--snapshot_actions--secure"></a>
### Nested Schema for `snapshot_actions.secure`

Optional:

- `enable` (Boolean) enable defaults to false. Flag to enable link on the snapshot
- `secure` (Number) secure defaults to 1 day. The time that the snapshot generation is to be secure for.
- `time_in_hours` (Boolean) time_in_hours or Days defaults to Days. False is days, true is hours.


<a id="nestedatt--snapshot_actions--time_to_live"></a>
### Nested Schema for `snapshot_actions.time_to_live`

Optional:

- `enable` (Boolean) enable defaults to false. Flag to enable link on the snapshot
- `time_in_hours` (Boolean) time_in_hours or Days defaults to Days. False is days, true is hours.
- `time_to_live` (Number) time_to_live defaults to 1 day. Gives the total time before expiry for these actions.



<a id="nestedblock--storage_group"></a>
### Nested Schema for `storage_group`

Required:

- `name` (String) Name of the storage group you would like to take a snapshot.


<a id="nestedatt--source_volume"></a>
### Nested Schema for `source_volume`

Read-Only:

- `capacity` (Number) The capacity of the snapshot volume in cylinders
- `capacity_gb` (Number) The capacity of the snapshot volume in GB
- `name` (String) The name of the SnapVX snapshot generation source volume

## Import

Import is supported using the following syntax:

```shell
# Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://mozilla.org/MPL/2.0/


# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# The command is
# terraform import powermax_snapshot.snapshot_test storage_group.snapshot_name
# Example: must be storage_group.snapshot_name
terraform import powermax_snapshot.snapshot_test storage_group.snapshot_name
# after running this command, populate the name field in the config file to start managing this resource
```