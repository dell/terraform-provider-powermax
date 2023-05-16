// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStorageGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: ProviderConfig + testAccSgDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powermax_storagegroup.test", "storage_groups.#", "1"),
					resource.TestCheckResourceAttr("data.powermax_storagegroup.test", "storage_groups.0.name", "esa_sg572"),
					resource.TestCheckResourceAttr("data.powermax_storagegroup.test", "storage_groups.0.id", "esa_sg572"),
					resource.TestCheckResourceAttr("data.powermax_storagegroup.test", "storage_groups.0.device_emulation", "FBA"),
					resource.TestCheckResourceAttr("data.powermax_storagegroup.test", "storage_groups.0.type", "Standalone"),
					resource.TestCheckResourceAttr("data.powermax_storagegroup.test", "storage_groups.0.unprotected", "true"),
				),
			},
			// Read testing
			{
				Config: ProviderConfig + testAccSgAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powermax_storagegroup.test", "storage_groups.#", "10"),
				),
			},
		},
	})
}

const testAccSgDataSourceConfig = `
data "powermax_storagegroup" "test" {
  filter {
    names = ["esa_sg572"]
  }
}
`

const testAccSgAllDataSourceConfig = `
data "powermax_storagegroup" "test" {
}
`
