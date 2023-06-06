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

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test to Fetch all Hostgroup details.
func TestAccHostGroupDatasource(t *testing.T) {
	var hostGroupName = "data.powermax_hostgroup.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + HostGroupDataSourceParamsAll,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(hostGroupName, "filter.#", "0"),
				),
			},
		},
	})
}

func TestAccHostGroupDatasourceFiltered(t *testing.T) {
	var hostGroupName = "data.powermax_hostgroup.groups"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + HostGroupDataSourceParamsFiltered,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(hostGroupName, "host_group_details.#", "2"),
					resource.TestCheckResourceAttr(hostGroupName, "host_group_details.0.consistent_lun", "false"),
					resource.TestCheckResourceAttr(hostGroupName, "host_group_details.0.host.#", "2"),
					resource.TestCheckResourceAttr(hostGroupName, "host_group_details.0.host.0.host_id", "Example_1_host_host_group"),
					resource.TestCheckResourceAttr(hostGroupName, "host_group_details.0.host.0.initiator.#", "1"),
					resource.TestCheckResourceAttr(hostGroupName, "host_group_details.0.host_group_id", "tfacc_host_group_example_1"),
					resource.TestCheckResourceAttr(hostGroupName, "host_group_details.0.name", "tfacc_host_group_example_1"),
					resource.TestCheckResourceAttr(hostGroupName, "host_group_details.0.num_of_hosts", "2"),
					resource.TestCheckResourceAttr(hostGroupName, "host_group_details.0.num_of_initiators", "1"),
					resource.TestCheckResourceAttr(hostGroupName, "host_group_details.0.num_of_masking_views", "1"),
					resource.TestCheckResourceAttr(hostGroupName, "host_group_details.0.port_flags_override", "false"),
					resource.TestCheckResourceAttr(hostGroupName, "host_group_details.0.type", "Fibre"),
				),
			},
		},
	})
}

func TestAccHostGroupDatasourceFilteredError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + HostGroupDataSourceFilterError,
				ExpectError: regexp.MustCompile(`.*Error getting the details of host group*.`),
			},
		},
	})
}

var HostGroupDataSourceFilterError = `
# List a specific hostgroup
data "powermax_hostgroup" "groups" {
  filter {
    names = ["tfacc_fake_host_group"]
  }
}

output "groups" {
  value = data.powermax_hostgroup.groups
}
`

var HostGroupDataSourceParamsFiltered = `
# List a specific hostgroup
data "powermax_hostgroup" "groups" {
  filter {
    names = ["tfacc_host_group_example_1", "tfacc_host_group_example_2"]
  }
}

output "groups" {
  value = data.powermax_hostgroup.groups
}
`
var HostGroupDataSourceParamsAll = `
# List all hostgroups
data "powermax_hostgroup" "all" {}

output "all" {
  value = data.powermax_hostgroup.all
}
`
