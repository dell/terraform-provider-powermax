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

# Available actions: Create, Update (vol_name, mobility_id_enabled, size, cap_unit), Delete and Import an existing volume from the PowerMax Array.
# After `terraform apply` of this example file it will create a new volume with the name set in `vol_name` attribute on the PowerMax

# PowerMax volumes is an identifiable unit of data storage. Storage groups are sets of volumes.
resource "powermax_volume" "test" {

  # Attributes which are able to be modified after create (vol_name, mobility_id_enabled, size, cap_unit)

  # Required name of the volume to be created
  vol_name = "terraform_volume"

  # Required size of the volume
  size = 2.45

  # Required name of the storage group which the volume will be created with
  sg_name = "terraform_sg"

  # Optional Default Unit is GB
  # Possible units are MB, GB, TB, and CYL
  cap_unit = "GB"

  # Optional enable the mobility id 
  mobility_id_enabled = false

}

# After the execution of above resource block, a PowerMax volume has been created at PowerMax array.
# For more information about the newly created resource use the `terraform show` command to review the current state