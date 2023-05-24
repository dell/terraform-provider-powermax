// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStorageGroupDataSource(t *testing.T) {
	var storageGroupTerraformName = "data.powermax_storagegroup.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create the storage group to be tested as a data source
			{
				Config: ProviderConfig + StorageGroupResourceConfig,
			},
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

var SgDataSourceConfig = `
data "powermax_storagegroup" "test" {
  filter {
    names = ["tfacc_sg_1"]
  }
}
`

var SgAllDataSourceConfig = `
data "powermax_storagegroup" "test" {
}
`
