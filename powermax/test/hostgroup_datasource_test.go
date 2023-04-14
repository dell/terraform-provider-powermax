// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test to Fetch Host details
func TestAccHostGroupDatasource(t *testing.T) {
	var hostGroupName = "data.powermax_hostgroup.groups"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + HostGroupDataSourceParamsAll,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(hostGroupName, "host_group_id.#", "9"),
				),
			},
		},
	})
}

var HostGroupDataSourceParamsAll = `
# List all hostgroups
data "powermax_hostgroup" "groups" {}

output "groups" {
  value = data.powermax_hostgroup.groups
}
`
