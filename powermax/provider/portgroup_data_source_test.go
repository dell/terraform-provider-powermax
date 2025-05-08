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

	"github.com/bytedance/mockey"
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

func TestAccPortGroupDatasourceError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ReadPortgroupByID).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + PortGroupDataSourceParamsAll,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccPortGroupDatasourceListError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetPortGroupList).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + PortGroupDataSourceParamsAll,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
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
