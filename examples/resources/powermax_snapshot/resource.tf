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
resource "powermax_snapshot" "snapshot_1" {
  # Required, what storage group to take the snapshot
  storage_group {
    name = "test_sg"
  }
  snapshot_actions {
    # Required, name of new snapshot
    name = "new_test_snapshot"
    # Optional
    # secure = {
    #   enable = true
    #   time_in_hours = true
    #   secure = 3
    # }
    # Optional
    # time_to_live = {
    #   enable = true
    #   time_in_hours = false
    #   time_to_live = 1
    # }
    # Optional
    # link = {
    #   enable = false
    #   target_storage_group = "test_target_sg"
    #   no_compression = true
    #   remote = false
    #   copy = false
    # }
    # Optional
    # restore = {
    #   enable = true
    #   remote = false
    # }
  }
}