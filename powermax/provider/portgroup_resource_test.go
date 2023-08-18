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
	"fmt"
	"regexp"
	"terraform-provider-powermax/powermax/helper"
	"testing"

	. "github.com/bytedance/mockey"
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
	var errorString []string
	errorString = append(errorString, "mock error")
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
			// Read test Error
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.ReadPortgroupByID).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + createPortGroupConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Modify Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.UpdatePortGroup).Return(nil, nil, errorString).Build()
				},
				Config:      ProviderConfig + createPortGroupConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// auto checks delete to clean up the test
		},
	})
}

func TestAccPortgroupResourceCreateError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read test
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.CreatePortGroup).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + createPortGroupConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
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
