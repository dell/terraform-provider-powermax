// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStorageGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + StorageGroupResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "storage_group_id", "terraform_sg"),
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "srp_id", "SRP_1"),
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "slo", "Gold"),
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "service_level", "Gold"),
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "slo_compliance", "STABLE"),
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "num_of_child_sgs", "0"),
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "num_of_masking_views", "0"),
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "num_of_parent_sgs", "0"),
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "compression", "true"),
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "unprotected", "true"),
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "cap_gb", "0"),
					// Check map value host_io_limit
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "host_io_limit.host_io_limit_io_sec", "1000"),
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "host_io_limit.host_io_limit_mb_sec", "1000"),
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "host_io_limit.dynamicDistribution", "Never"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "powermax_storagegroup.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update storage_group_id and Read testing
			{
				Config: ProviderConfig + StorageGroupRenameResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "storage_group_id", "terraform_sg_2"),
				),
			},
			// Update compression and Read testing
			{
				Config: ProviderConfig + StorageGroupUpdateCompressionResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "storage_group_id", "terraform_sg"),
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "compression", "false"),
				),
			},
			// Update and Read testing
			{
				Config: ProviderConfig + StorageGroupUpdateHostIOResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "storage_group_id", "terraform_sg"),
					// Check map value host_io_limit
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "host_io_limit.host_io_limit_io_sec", "2000"),
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "host_io_limit.host_io_limit_mb_sec", "2000"),
					resource.TestCheckResourceAttr("powermax_storagegroup.test", "host_io_limit.dynamicDistribution", "Never"),
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
  	host_io_limit = {
    	host_io_limit_io_sec = "1000"
    	host_io_limit_mb_sec = "1000"
    	dynamicDistribution  = "Never"
  	}
}
`

var StorageGroupUpdateCompressionResourceConfig = `
resource "powermax_storagegroup" "test" {
	storage_group_id = "terraform_sg"
  	srp_id           = "SRP_1"
  	slo              = "Gold"
  	compression              = false
  	host_io_limit = {
    	host_io_limit_io_sec = "1000"
    	host_io_limit_mb_sec = "1000"
    	dynamicDistribution  = "Never"
  	}
}
`

var StorageGroupUpdateHostIOResourceConfig = `
resource "powermax_storagegroup" "test" {
	storage_group_id = "terraform_sg"
  	srp_id           = "SRP_1"
  	slo              = "Gold"
  	compression              = false
  	host_io_limit = {
    	host_io_limit_io_sec = "2000"
    	host_io_limit_mb_sec = "2000"
    	dynamicDistribution  = "Never"
  	}
}
`
