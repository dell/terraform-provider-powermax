// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package provider

import (
	"regexp"
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
					resource.TestCheckResourceAttr(portGroupName, "port_groups.0.type", "SCSI_FC"),
				),
			},
		},
	})
}
func TestAccPortGroupDatasourceFilteredError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + PortGroupDataSourceFilterError,
				ExpectError: regexp.MustCompile(`.*Name of already created portgroup must be provided.*.`),
			},
		},
	})
}

// List a specific portgroup.
var PortGroupDataSourceFilterError = `
data "powermax_portgroups" "errgroups" {
  filter {
    names = ["tfacc_test1_fibre", "non-existent-port-group"]
  }
}`

var PortGroupDataSourceParamsAll = `
data "powermax_portgroups" "fibreportgroups" {
	filter {
		# Optional list of names to filter
		#names = [
		#  "tfacc_test1_fibre",
		#  "tfacc_test2_fibre",
		#]
		type = "fibre"
	}
}


output "fibreportgroups" {
  value = data.powermax_portgroups.fibreportgroups
} `
