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

# Available actions: Create, Update (name, ports), Delete and Import an existing portgroup from the PowerMax Array.
# After `terraform apply` of this example file it will create a new portgroup with the name set in `name` attribute on the PowerMax

# PowerMax port groups contain director and port identification and belong to a masking view. Ports can be added to and removed from the port group. Port groups that are no longer associated with a masking view can be deleted.
# Note the following recommendations:
# Port groups should contain four or more ports.
# Each port in a port group should be on a different director.
# A port can belong to more than one port group. However, for storage systems running HYPERMAX OS 5977 or higher, you cannot mix different types of ports (physical FC ports, virtual ports, and iSCSI virtual ports) within a single port group
resource "powermax_portgroup" "portgroup_1" {

  # Attributes which are able to be modified after create (name, ports)

  # Required The name of the portgroup. Only alphanumeric characters, underscores ( _ ), and hyphens (-) are allowed.
  name = "tfacc_pg_test_1"

  # Required The portgroup protocol. Protocols: SCSI_FC, iSCSI, NVMe_FC, NVMe_TCP
  protocol = "SCSI_FC"

  # Required The list of ports associated with the portgroup
  # Must include the director and port id in each object below
  ports = [
    {
      director_id = "OR-1C"
      port_id     = "0"
    }
  ]
}

# After the execution of above resource block, a PowerMax port group has been created at PowerMax array.
# For more information about the newly created resource use the `terraform show` command to review the current state