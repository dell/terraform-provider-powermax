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
func TestAccPortDatasource(t *testing.T) {
	var portName = "data.powermax_port.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + PortDataSourceParamsAll,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(portName, "filter.#", "0"),
				),
			},
		},
	})
}

func TestAccPortDatasourceFiltered(t *testing.T) {
	var portID = "data.powermax_port.portFilter"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + PortDataSourceParamsFiltered,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(portID, "port_details.#", "1"),
					resource.TestCheckResourceAttr(portID, "port_details.0.aclx", "true"),
					resource.TestCheckResourceAttr(portID, "port_details.0.avoid_reset_broadcast", "false"),
					resource.TestCheckResourceAttr(portID, "port_details.0.identifier", "5000097200097802"),
					resource.TestCheckResourceAttr(portID, "port_details.0.max_speed", "32"),
					resource.TestCheckResourceAttr(portID, "port_details.0.port_id", "2"),
					resource.TestCheckResourceAttr(portID, "port_details.0.director_id", "OR-1C"),
					resource.TestCheckResourceAttr(portID, "port_details.0.scsi_support1", "true"),
					resource.TestCheckResourceAttr(portID, "port_details.0.scsi_3", "true"),
					resource.TestCheckResourceAttr(portID, "port_details.0.vcm_state", "Enabled"),
					resource.TestCheckResourceAttr(portID, "port_details.0.wwn_node", "5000097200097bff"),
				),
			},
		},
	})
}

func TestAccPortDatasourceFilteredError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + PortDataSourceFilterError,
				ExpectError: regexp.MustCompile(`.*invalid format for port filter*.`),
			},
		},
	})
}

var PortDataSourceParamsAll = `
# List all ports
data "powermax_port" "all" {}

output "all" {
  value = data.powermax_port.all
}
`

var PortDataSourceParamsFiltered = `
# List a specific ports
data "powermax_port" "portFilter" {
  filter {
    # Should be in the format ["directorId:portId"]
    port_ids = ["OR-1C:2"]
  }
}

output "portFilter" {
  value = data.powermax_port.portFilter
}
`

var PortDataSourceFilterError = `
# List a specific ports
data "powermax_port" "portFilter" {
  filter {
    # Should be in the format ["directorId:portId"]
    port_ids = ["bad_port_id"]
  }
}

output "portFilter" {
  value = data.powermax_port.portFilter
}
`
