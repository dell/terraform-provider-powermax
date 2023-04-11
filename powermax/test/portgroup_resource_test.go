// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var createPortGroupConfig = `
resource "powermax_portgroup" "test_portgroup" {
	name = "tf_pg_test_1"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "FA-2D"
			port_id = "11"
		}
	]
}
`
var updatePortGroupConfig = `
resource "powermax_portgroup" "test_portgroup" {
	# This will be updated 
	name = "tf_pg_test_1_upd"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "FA-2D"
			port_id = "11"
		}
	]
}
`

func TestAccPortgroupResource(t *testing.T) {
	var portgroupTerraformName = "powermax_portgroup.test_portgroup"
	Init()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read test
			{
				Config: ProviderConfig + createPortGroupConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(portgroupTerraformName, "name", "tf_pg_test_1"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "protocol", "SCSI_FC"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "ports.0.director_id", "FA-2D"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "ports.0.port_id", "11"),

					// Verify Calculated values
					// numofmaskingviews
					resource.TestCheckResourceAttr(portgroupTerraformName, "numofmaskingviews", "0"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "numofports", "1"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "type", "Fibre"),
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
					resource.TestCheckResourceAttr(portgroupTerraformName, "name", "tf_pg_test_1_upd"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "protocol", "SCSI_FC"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "ports.0.director_id", "FA-2D"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "ports.0.port_id", "11"),

					// Verify Calculated values
					// numofmaskingviews
					resource.TestCheckResourceAttr(portgroupTerraformName, "numofmaskingviews", "0"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "numofports", "1"),
					resource.TestCheckResourceAttr(portgroupTerraformName, "type", "Fibre"),
				),
			},
			// auto checks delete to clean up the test
		},
	})
}

/*func testAccPortgroupResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "portgroup" "test" {
  configurable_attribute = %[1]q
}
`, configurableAttribute)
}*/
