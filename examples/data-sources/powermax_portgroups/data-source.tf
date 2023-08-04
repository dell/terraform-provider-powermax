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

# List fibre portgroups.
data "powermax_portgroups" "fibreportgroups" {
  # Optional filter to list specified Portgroups names and/or type
  filter {
    # type for which portgroups to be listed  - fibre or iscsi
    type = "fibre"
    # Optional list of IDs to filter
    names = [
      "tfacc_test1_fibre",
      #"test2_fibre",
    ]
  }
}

data "powermax_portgroups" "scsiportgroups" {
  filter {
    type = "iscsi"
    # Optional filter to list specified Portgroups Names
  }
}

# List all portgroups.
data "powermax_portgroups" "allportgroups" {
  #filter {
  # Optional list of IDs to filter
  #names = [
  #  "test1",
  #  "test2",
  #]
  #}
}

