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

func TestAccStorageGroupA(t *testing.T) {
	var storageGroupTerraformName = "powermax_storagegroup.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + StorageGroupResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(storageGroupTerraformName, "name", "tfacc_sg_1"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "id", "tfacc_sg_1"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "srp_id", "SRP_1"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "slo", "Gold"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "service_level", "Gold"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "slo_compliance", "STABLE"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "num_of_child_sgs", "0"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "num_of_masking_views", "0"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "num_of_parent_sgs", "0"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "num_of_vols", "0"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "compression", "true"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "unprotected", "true"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "cap_gb", "0"),
					// Check map value host_io_limit
					resource.TestCheckResourceAttr(storageGroupTerraformName, "host_io_limit.host_io_limit_io_sec", "1000"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "host_io_limit.host_io_limit_mb_sec", "1000"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "host_io_limit.dynamic_distribution", "Never"),
				),
			},
			// ImportState testing
			{
				ResourceName:      storageGroupTerraformName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: ProviderConfig + StorageGroupResourceConfig + StorageGroupVolumeResourceConfig,
			},
			// Update compression, volume_id and host_io_limit, then Read testing
			{
				Config: ProviderConfig + StorageGroupVolumeResourceConfig + StorageGroupVolumeDataResourceConfig + StorageGroupUpdateResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(storageGroupTerraformName, "name", "tfacc_sg_rename"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "id", "tfacc_sg_rename"),
					// check slo
					resource.TestCheckResourceAttr(storageGroupTerraformName, "slo", "Silver"),
					// check map value host_io_limit
					resource.TestCheckResourceAttr(storageGroupTerraformName, "host_io_limit.host_io_limit_io_sec", "2000"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "host_io_limit.host_io_limit_mb_sec", "2000"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "host_io_limit.dynamic_distribution", "Never"),
					// check Compression
					resource.TestCheckResourceAttr(storageGroupTerraformName, "compression", "false"),
					// check volume_ids
					resource.TestCheckResourceAttr(storageGroupTerraformName, "num_of_vols", "1"),
				),
			},
			{
				// Remove volume ahead of storage group
				Config: ProviderConfig + StorageGroupUpdateVolumeResourceConfig,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccStorageGroupNoHostIOLimit(t *testing.T) {
	var storageGroupTerraformName = "powermax_storagegroup.tfacc_sg_no_host_io_limit"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + StorageGroupNoHostIOLimitResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(storageGroupTerraformName, "name", "tfacc_sg_no_host_io_limit"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "id", "tfacc_sg_no_host_io_limit"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "srp_id", "SRP_1"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "slo", "Gold"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "service_level", "Gold"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "slo_compliance", "STABLE"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "num_of_child_sgs", "0"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "num_of_masking_views", "0"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "num_of_parent_sgs", "0"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "num_of_vols", "0"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "compression", "true"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "unprotected", "true"),
					resource.TestCheckResourceAttr(storageGroupTerraformName, "cap_gb", "0"),
				),
			},
		},
	})
}

func TestAccStorageGroupCreateError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + StorageGroupErrorCreateResourceConfig,
				ExpectError: regexp.MustCompile(".*Client Error*."),
			},
		},
	})
}

func TestAccStorageGroupUpdateError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + StorageGroupErrorUpdateResourceConfig,
			},
			{
				Config:      ProviderConfig + StorageGroupErrorUpdateResourceConfig2,
				ExpectError: regexp.MustCompile(".*Failed to update*."),
			},
		},
	})
}

var StorageGroupResourceConfig = `
resource "powermax_storagegroup" "test" {
	name             = "tfacc_sg_1"
  	srp_id           = "SRP_1"
  	slo              = "Gold"
  	host_io_limit = {
    	host_io_limit_io_sec = "1000"
    	host_io_limit_mb_sec = "1000"
    	dynamic_distribution  = "Never"
  	}
}
`

var StorageGroupVolumeResourceConfig = `
resource "powermax_volume" "volume_test" {
	vol_name = "tfacc_sg_vol_1"
	size = 2.45
	cap_unit = "GB"
	sg_name = "tfacc_sg_1"
}
`

var StorageGroupVolumeDataResourceConfig = `
data "powermax_volume" "volume_datasource_test" {
	filter {
		volume_identifier = "tfacc_sg_vol_1"
	}
}
`

var StorageGroupUpdateResourceConfig = `
resource "powermax_storagegroup" "test" {
	name             = "tfacc_sg_rename"
  	slo              = "Silver"
	srp_id           = "SRP_1"
	compression      = false
  	host_io_limit = {
    	host_io_limit_io_sec = "2000"
    	host_io_limit_mb_sec = "2000"
    	dynamic_distribution  = "Never"
  	}
	volume_ids = [data.powermax_volume.volume_datasource_test.volumes.0.id]
}
`

var StorageGroupUpdateVolumeResourceConfig = `
resource "powermax_storagegroup" "test" {
	name             = "tfacc_sg_rename"
  	slo              = "Silver"
	srp_id           = "SRP_1"
	compression      = false
  	host_io_limit = {
    	host_io_limit_io_sec = "2000"
    	host_io_limit_mb_sec = "2000"
    	dynamic_distribution  = "Never"
  	}
}
`

var StorageGroupNoHostIOLimitResourceConfig = `
resource "powermax_storagegroup" "tfacc_sg_no_host_io_limit" {
	name             = "tfacc_sg_no_host_io_limit"
  	srp_id           = "SRP_1"
  	slo              = "Gold"
}
`

var StorageGroupErrorCreateResourceConfig = `
resource "powermax_storagegroup" "test_error_1" {
	name             = "tfacc_sg_error_create"
  	srp_id           = "SRP_1"
  	slo              = "Gold"
}

resource "powermax_storagegroup" "test_error_2" {
	name             = "tfacc_sg_error_create"
  	srp_id           = "SRP_1"
  	slo              = "Gold"
}
`

var StorageGroupErrorUpdateResourceConfig = `
resource "powermax_storagegroup" "test_error_update" {
	name             = "tfacc_sg_error_update"
  	srp_id           = "SRP_1"
  	slo              = "Gold"
}

resource "powermax_storagegroup" "test_error_update_2" {
	name             = "tfacc_sg_error_update_rename"
  	srp_id           = "SRP_1"
  	slo              = "Gold"
}
`

var StorageGroupErrorUpdateResourceConfig2 = `
resource "powermax_storagegroup" "test_error_update" {
	name             = "tfacc_sg_error_update_rename"
  	srp_id           = "srp-non-existent"
  	slo              = "slo-non-existent"
	compression = false
	host_io_limit = {
    	host_io_limit_io_sec = "non-existent"
    	host_io_limit_mb_sec = ""
    	dynamic_distribution  = ""
  	}
	workload = "workload-non-existent"
	volume_ids = ["non_existent_vol_id"]
}
`
