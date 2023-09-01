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