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
func TestAccHostDatasource(t *testing.T) {
	var hostName = "data.powermax_host.HostDs"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + HostDataSourceParamsAll,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(hostName, "hosts.#", "1"),
				),
			},
		},
	})
}

func TestAccHostDatasourceFilteredError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + HostDataSourceFilterError,
				ExpectError: regexp.MustCompile(`.*Error reading host with id*.`),
			},
		},
	})
}
func TestAccHostDatasourceFilterEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + HostDataSourceFilterEmpty,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powermax_host.hostEmptyFilter", "hosts.#"),
				),
			},
		},
	})
}
func TestAccHostDatasourceGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + HostDataSourceGetAll,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powermax_host.hostGetAll", "hosts.#"),
				),
			},
		},
	})
}

func TestAccHostDatasourceError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetHostList).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + HostDataSourceGetAll,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var HostDataSourceParamsAll = `
data "powermax_host" "HostDs" {
	filter {
		# Optional list of IDs to filter
		names = [
		  "tfacc_host_test1",
		]
	}
	
}
output "hostDsResult" {
	value = data.powermax_host.HostDs
 }
`

var HostDataSourceFilterError = `
# List a specific host
data "powermax_host" "hosts" {
  filter {
    names = ["non-existent-host"]
  }
}

output "hosts" {
  value = data.powermax_host.hosts
}
`
var HostDataSourceFilterEmpty = `
data "powermax_host" "hostEmptyFilter" {
	   filter {
    		names = []
	   }
}`

var HostDataSourceGetAll = `
data "powermax_host" "hostGetAll" {
}`
