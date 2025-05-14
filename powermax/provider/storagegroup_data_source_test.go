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

func TestAccStorageGroupDataSource(t *testing.T) {
	var storageGroupTerraformName = "data.powermax_storagegroup.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter read testing
			{
				Config: ProviderConfig + StorageGroupResourceConfig + SgDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(storageGroupTerraformName, "storage_groups.#", "1"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "storage_groups.0.name", "tfacc_sg_1"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "storage_groups.0.id", "tfacc_sg_1"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "storage_groups.0.srp_id", "SRP_1"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "storage_groups.0.slo", "Gold"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "storage_groups.0.service_level", "Gold"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "storage_groups.0.slo_compliance", "STABLE"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "storage_groups.0.num_of_child_sgs", "0"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "storage_groups.0.num_of_masking_views", "0"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "storage_groups.0.num_of_parent_sgs", "0"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "storage_groups.0.num_of_vols", "0"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "storage_groups.0.compression", "true"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "storage_groups.0.unprotected", "true"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "storage_groups.0.cap_gb", "0"),
				),
			},
			// read all testing
			{
				Config: ProviderConfig + SgAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powermax_storagegroup.test", "storage_groups.#"),
				),
			},
		},
	})
}

func TestAccStorageGroupDataSourceErrorNotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + SgDataSourceConfigError,
				ExpectError: regexp.MustCompile(`.*StorageGroup error_fake_sg_datasoure is not on the powermax*.`),
			},
		},
	})
}

func TestAccStorageGroupDataSourceError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetStorageGroupList).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SgAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccStorageGroupDataSourceMapperError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.UpdateSgState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SgDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var SgDataSourceConfig = `
data "powermax_storagegroup" "test" {
  filter {
    names = ["tfacc_sg_1"]
  }
}
`

var SgDataSourceConfigError = `
data "powermax_storagegroup" "test" {
  filter {
    names = ["error_fake_sg_datasoure"]
  }
}
`

var SgAllDataSourceConfig = `
data "powermax_storagegroup" "test" {
}
`
