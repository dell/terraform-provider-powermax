// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
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
					resource.TestCheckResourceAttr(hostGroupName, "host_group_details.#", "4"),
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
					resource.TestCheckResourceAttr(hostGroupName, "host_group_details.0.host.0.host_id", "81"),
					resource.TestCheckResourceAttr(hostGroupName, "host_group_details.0.host.0.initiator.#", "1"),
					resource.TestCheckResourceAttr(hostGroupName, "host_group_details.0.host_group_id", "host_group_example_1"),
					resource.TestCheckResourceAttr(hostGroupName, "host_group_details.0.name", "host_group_example_1"),
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
    names = ["non-existent-host-group"]
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
    names = ["host_group_example_1", "host_group_example_2"]
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
