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

# Available actions: Create, Update (name, compression, host_io_limit, workload, slo, srp_id, volume_ids), Delete and Import an existing storage group from the PowerMax Array.
# After `terraform apply` of this example file it will create a new storage group with the name set in `name` attribute on the PowerMax

# PowerMax storage groups are a collection of devices that are stored on the array.
# An application, a server, or a collection of servers use them.
resource "powermax_storagegroup" "test" {

  # Attributes which are able to be modified after create (name, compression, host_io_limit, workload, slo, srp_id, volume_ids)

  # Required the name of the new storage group
  name = "terraform_sg"

  # Required The Srp to be associated with the Storage Group. If you dont want an SRP the srp_id can be set to 'None'
  srp_id = "SRP_1"

  # Optional the service level of the storage group
  slo = "Gold"

  # Optional enable or disable compression on the storage group (Default to disabled)
  compression = false

  # Optional the workload of the storage group
  workload = "workload"

  # Optional Set Host I/O limits for the specified storage sroup
  host_io_limit = {
    # The MBs per Second Host IO limit for the specified storage group, NOLIMIT means no limits
    host_io_limit_io_sec = "1000"
    # The IOs per Second Host IO limit for the specified storage group , NOLIMIT means no limits
    host_io_limit_mb_sec = "1000"
    # The dynamic distribution type which can be "Never","Always" or "OnFailure"
    dynamic_distribution = "Never"
  }

  # Optional a list of volume ids to be added to the storage groups
  volume_ids = ["0008F"]
}

# After the execution of above resource block, a PowerMax storage group has been created at PowerMax array.
# For more information about the newly created resource use the `terraform show` command to review the current state