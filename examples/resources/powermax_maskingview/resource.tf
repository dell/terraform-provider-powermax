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

# Available actions: Create, Update (name), Delete and Import an existing maskingview from the PowerMax Array.
# After `terraform apply` of this example file it will create a new masking_view with the name set in `name` attribute on the PowerMax

# PowerMax masking views are a container of a storage group, a port group, and an initiator group, and makes the storage group visible to the host. 
# Devices are masked and mapped automatically. The groups must contain some devices entries.
resource "powermax_maskingview" "test" {

  # Attributes which are able to be modified after create (name)

  # Required the name of the new masking view
  # Only alphanumeric characters, underscores ( _ ), and hyphens (-) are allowed
  name = "terraform_mv"

  # Required the storage group id to be assoiciated with the new masking view
  storage_group_id = "TestnewSG"

  # Required host to be assoiciated with the new masking view
  # NOTE if host_id is set then host_group_id must be empty string ""
  host_id = "Host124"

  # Required host group to be assoiciated with the new masking view
  # NOTE if host_group_id is set then host_id must be an empty string ""
  host_group_id = ""

  # Required the port group to be assoiciated with the new masking view
  port_group_id = "TestnewSG_PG"
}

# After the execution of above resource block, a PowerMax maskingview has been created at PowerMax array.
# For more information about the newly created resource use the `terraform show` command to review the current state