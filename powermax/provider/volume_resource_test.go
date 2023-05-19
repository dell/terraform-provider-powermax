// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVolume(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + VolumeResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "sg_name", "terraform_vol_sg"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "vol_name", "test_acc_create_volume_1"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "size", "2.45"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "cap_unit", "GB"),

					resource.TestCheckResourceAttr("powermax_volume.volume_test", "type", "TDEV"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "emulation", "FBA"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "allocated_percent", "0"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "status", "Ready"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "reserved", "false"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "pinned", "false"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "reserved", "false"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "encapsulated", "false"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "num_of_storage_groups", "1"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "num_of_front_end_paths", "0"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "snapvx_source", "false"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "snapvx_target", "false"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "has_effective_wwn", "false"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "mobility_id_enabled", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName: "powermax_volume.volume_test",
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, "test_acc_create_volume_1", states[0].Attributes["vol_name"])
					assert.Equal(t, "2.45", states[0].Attributes["size"])
					assert.Equal(t, "GB", states[0].Attributes["cap_unit"])
					assert.Equal(t, "", states[0].Attributes["sg_name"])
					return nil
				},
			},
			// Update Name, Size, Mobility and Read testing
			{
				Config: ProviderConfig + VolumeUpdateNameSizeMobility,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "vol_name", "test_acc_create_volume_2"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "size", "1"),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "cap_unit", "TB"),
				),
			},
		},
	})
}

func TestAccVolume_Invalid_Config(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Config with no SG
			{
				Config:      ProviderConfig + VolumeConfigNoSG,
				ExpectError: regexp.MustCompile("Error creating volume"),
			},
			// Config with invalid unit
			{
				Config:      ProviderConfig + VolumeConfigInvalidCYL,
				ExpectError: regexp.MustCompile("Error creating volume"),
			},
			// Config with invalid SG name
			{
				Config:      ProviderConfig + VolumeConfigInvalidSG,
				ExpectError: regexp.MustCompile("Error creating volume"),
			},
		},
	})
}

func TestAccVolume_Error_Updating(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Normal Config
			{
				Config:      ProviderConfig + VolumeConfigWithCYL,
				ExpectError: nil,
			},
			// Invalid name
			{
				Config:      ProviderConfig + VolumeConfigInvalidSG,
				ExpectError: regexp.MustCompile("Failed to modify."),
			},
			// Config with invalid SG name
			{
				Config:      ProviderConfig + VolumeConfigInvalidName,
				ExpectError: regexp.MustCompile("Failed to modify."),
			},
		},
	})
}

var VolumeResourceConfig = `
resource "powermax_volume" "volume_test" {
	vol_name = "test_acc_create_volume_1"
	size = 2.45
	cap_unit = "GB"
	sg_name = "terraform_vol_sg"
}
`

var VolumeUpdateNameSizeMobility = `
resource "powermax_volume" "volume_test" {
	vol_name = "test_acc_create_volume_2"
	size = 1
	cap_unit = "TB"
	sg_name = "terraform_vol_sg"
	mobility_id_enabled = "true"	
}
`

var VolumeConfigNoSG = `
resource "powermax_volume" "volume_test" {
	vol_name = "test_acc_create_volume_2"
	size = 0.5
	cap_unit = "CYL"
}
`

var VolumeConfigInvalidCYL = `
resource "powermax_volume" "volume_test" {
	vol_name = "test_acc_create_volume_2"
	sg_name = "terraform_vol_sg"
	size = 0.5
	cap_unit = "CYL"
}
`

var VolumeConfigWithCYL = `
resource "powermax_volume" "volume_test" {
	vol_name = "test_acc_create_volume_1"
	size = 1
	cap_unit = "CYL"
	sg_name = "terraform_vol_sg"
}
`
var VolumeConfigInvalidSG = `
resource "powermax_volume" "volume_test" {
	vol_name = "test_acc_create_volume_2"
	sg_name = "invalid#SG"
	size = 0.5
	cap_unit = "CYL"
}
`

var VolumeConfigInvalidName = `
resource "powermax_volume" "volume_test" {
	vol_name = "!@#$%"
	sg_name = "terraform_vol_sg"
	size = 0.5
	cap_unit = "CYL"
}
`
