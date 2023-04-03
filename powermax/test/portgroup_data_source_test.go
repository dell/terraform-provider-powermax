// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test to Fetch Host details.
func TestAccPortGroupDatasource(t *testing.T) {
	var portGroupName = "data.powermax_portgroups.fibreportgroups"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + PortGroupDataSourceParamsAll,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(portGroupName, "port_groups.0.type", "Fibre"),
				),
			},
		},
	})
}

var PortGroupDataSourceParamsAll = `
data "powermax_portgroups" "fibreportgroups" {
	filter {
		# Optional list of names to filter
		names = [
		  "test1_fibre",
		  "test2_fibre",
		]
		type = "fibre"
	}
}


output "fibreportgroups" {
  value = data.powermax_portgroups.fibreportgroups
} `
