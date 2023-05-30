// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var createPortGroupConfig = `
resource "powermax_portgroup" "test_portgroup" {
	name = "tfacc_pg_test_1"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "OR-1C"
			port_id = "0"
		}
	]
}
`
var updatePortGroupConfig = `
resource "powermax_portgroup" "test_portgroup" {
	# This will be updated 
	name = "tfacc_pg_test_1_upd"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "OR-2C"
			port_id = "2"
		}
	]
}
`

func TestAccPortgroupResource(t *testing.T) {
	var portgroupTerraformName = "powermax_portgroup.test_portgroup"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read test
			{
				Config: ProviderConfig + createPortGroupConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(portgroupTerraformName, "name", "tfacc_pg_test_1"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "protocol", "SCSI_FC"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "ports.0.director_id", "OR-1C"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "ports.0.port_id", "0"),

					// Verify Calculated values
					// numofmaskingviews
					resource.TestCheckResourceAttr(portgroupTerraformName, "numofmaskingviews", "0"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "numofports", "1"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "type", "SCSI_FC"),
				),
			},
			// Import testing
			{
				ResourceName:            portgroupTerraformName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"protocol"},
			},
			// Update testing
			{
				Config: ProviderConfig + updatePortGroupConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(portgroupTerraformName, "name", "tfacc_pg_test_1_upd"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "protocol", "SCSI_FC"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "ports.0.director_id", "OR-2C"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "ports.0.port_id", "2"),

					// Verify Calculated values
					// numofmaskingviews
					resource.TestCheckResourceAttr(portgroupTerraformName, "numofmaskingviews", "0"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "numofports", "1"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "type", "SCSI_FC"),
				),
			},
			// auto checks delete to clean up the test
		},
	})
}

func TestAccPortGroupResourceError(t *testing.T) {
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
var PortGroupRourceError = `
resource "powermax_portgroup" "test_portgroup" {
	name = "tfacc_error"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "OR-1F"
			port_id = "0"
		}
	]
}
`
