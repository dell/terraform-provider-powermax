// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStorageGroup(t *testing.T) {
	var storageGroupTerraformName = "powermax_storagegroup.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + StorageGroupResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(storageGroupTerraformName, "storage_group_id", "terraform_sg"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "srp_id", "SRP_1"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "slo", "Gold"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "service_level", "Gold"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "slo_compliance", "STABLE"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "num_of_child_sgs", "0"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "num_of_masking_views", "0"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "num_of_parent_sgs", "0"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "compression", "true"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "unprotected", "true"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "cap_gb", "0"),
					// Check map value host_io_limit
					resource.TestCheckResourceAttr(storageGroupTerraformName, "host_io_limit.host_io_limit_io_sec", "1000"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "host_io_limit.host_io_limit_mb_sec", "1000"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "host_io_limit.dynamicDistribution", "Never"),
				),
			},
			// ImportState testing
			{
				ResourceName:      storageGroupTerraformName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update storage_group_id, compression, and hostio_limit, then Read testing
			{
				Config: ProviderConfig + StorageGroupRenameResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(storageGroupTerraformName, "storage_group_id", "terraform_sg_2"),
					// Check map value host_io_limit
					resource.TestCheckResourceAttr(storageGroupTerraformName, "host_io_limit.host_io_limit_io_sec", "2000"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "host_io_limit.host_io_limit_mb_sec", "2000"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "host_io_limit.dynamicDistribution", "Never"),
					// Check Compression
					resource.TestCheckResourceAttr(storageGroupTerraformName, "compression", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

var StorageGroupResourceConfig = `
resource "powermax_storagegroup" "test" {
	storage_group_id = "terraform_sg"
  	srp_id           = "SRP_1"
  	slo              = "Gold"
  	host_io_limit = {
    	host_io_limit_io_sec = "1000"
    	host_io_limit_mb_sec = "1000"
    	dynamicDistribution  = "Never"
  	}
}
`

var StorageGroupRenameResourceConfig = `
resource "powermax_storagegroup" "test" {
	storage_group_id = "terraform_sg_2"
  	srp_id           = "SRP_1"
  	slo              = "Gold"
	compression      = false
  	host_io_limit = {
    	host_io_limit_io_sec = "2000"
    	host_io_limit_mb_sec = "2000"
    	dynamicDistribution  = "Never"
  	}
}
`
