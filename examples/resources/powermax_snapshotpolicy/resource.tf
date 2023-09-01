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

# Available actions: Create, Update (snapshot_policy_name, storage_groups, interval, snapshot_count, compliance_count_critical, compliance_count_warning, offset_minutes, secure), Delete and Import an existing snapshot policy from the PowerMax Array.
# After `terraform apply` of this example file it will create a new snapshot policy with the name set in `snapshot_policy_name` attribute on the PowerMax

# PowerMax snapshot policy feature provides snapshot orchestration at scale (1,024 snaps per storage group).
# The resource simplifies snapshot management for standard and cloud snapshots.
# This resouce will take snapshots on a periodic basic based on the configuration below.

resource "powermax_snapshotpolicy" "terraform_sp" {

  # Attributes which are able to be modified after create (snapshot_policy_name, storage_groups, interval, snapshot_count, compliance_count_critical, compliance_count_warning, offset_minutes, secure)

  # Required Field will become the name of the snapshot policy
  snapshot_policy_name = "terraform_sp"

  # should only be set for modify/edit operation , not supported during create. 
  # Also the destroy/delete will also unlink any associted storage groups from Snapshot Policy before deleting the snapshot policy.
  # storage_groups =  ["tfacc_sp_sg1", "tfacc_sp_sg2"]

  # Default values are defined below and can be modifed before or after create 

  # The interval between snapshots
  # interval             = "1 Hour"

  # The number of the snapshots that will be maintained by the snapshot policy
  # snapshot_count       = "48"

  # The number of snapshots which are not failed or bad when compliance changes to critical.
  # compliance_count_critical = 46

  # The number of snapshots which are not failed or bad when compliance changes to warning.
  # compliance_count_warning  = 47

  # The number of minutes from 00:00 on a Monday morning when the policy should run. Default is 0 if not specified.
  # offset_minutes            = 420

  # The snapshot policy will create secure snapshots
  # secure = false

}

# After the execution of above resource block, a PowerMax snapshot policy has been created at PowerMax array.
# For more information about the newly created resource use the `terraform show` command to review the current state