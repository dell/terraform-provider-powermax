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
				Config: ProviderConfig + VolStorageGroupConfig + VolumeResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "sg_name", resourceVolSGName),
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "vol_name", resourceVolName),
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
					assert.Equal(t, resourceVolName, states[0].Attributes["vol_name"])
					assert.Equal(t, "2.45", states[0].Attributes["size"])
					assert.Equal(t, "GB", states[0].Attributes["cap_unit"])
					assert.Equal(t, "", states[0].Attributes["sg_name"])
					return nil
				},
			},
			// Update Name, Size, Mobility and Read testing
			{
				Config: ProviderConfig + VolStorageGroupConfig + VolumeUpdateNameSizeMobility,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_volume.volume_test", "vol_name", resourceVolName+"_2"),
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
				Config:      ProviderConfig + VolStorageGroupConfig + VolumeConfigNoSG,
				ExpectError: regexp.MustCompile("Error creating volume"),
			},
			// Config with invalid unit
			{
				Config:      ProviderConfig + VolStorageGroupConfig + VolumeConfigInvalidCYL,
				ExpectError: regexp.MustCompile("Error creating volume"),
			},
			// Config with invalid SG name
			{
				Config:      ProviderConfig + VolStorageGroupConfig + VolumeConfigInvalidSG,
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
				Config:      ProviderConfig + VolStorageGroupConfig + VolumeConfigWithCYL,
				ExpectError: nil,
			},
			// Invalid SG name
			{
				Config:      ProviderConfig + VolStorageGroupConfig + VolumeConfigInvalidSG,
				ExpectError: regexp.MustCompile("Failed to modify."),
			},
			// Invalid name
			{
				Config:      ProviderConfig + VolStorageGroupConfig + VolumeConfigInvalidName,
				ExpectError: regexp.MustCompile("Failed to modify."),
			},
		},
	})
}

var resourceVolSGName = fmt.Sprintf("tfacc_res_vol_sg_%s", ResourceSuffix)
var resourceVolName = fmt.Sprintf("tfacc_res_vol_%s", ResourceSuffix)

var VolStorageGroupConfig = fmt.Sprintf(`
resource "powermax_storagegroup" "sg_vol_test" {
  name             = "%s"
  srp_id           = "SRP_1"
  slo              = "Gold"
}
`, resourceVolSGName)

var VolumeResourceConfig = fmt.Sprintf(`
resource "powermax_volume" "volume_test" {
	vol_name = "%s"
	size = 2.45
	cap_unit = "GB"
	sg_name = powermax_storagegroup.sg_vol_test.name
  	depends_on = [
    	powermax_storagegroup.sg_vol_test
  	]
}
`, resourceVolName)

var VolumeUpdateNameSizeMobility = fmt.Sprintf(`
resource "powermax_volume" "volume_test" {
	vol_name = "%s_2"
	size = 1
	cap_unit = "TB"
	sg_name = powermax_storagegroup.sg_vol_test.name
	mobility_id_enabled = "true"
  	depends_on = [
    	powermax_storagegroup.sg_vol_test
  	]
}
`, resourceVolName)

var VolumeConfigNoSG = fmt.Sprintf(`
resource "powermax_volume" "volume_test" {
	vol_name = "%s"
	size = 0.5
	cap_unit = "CYL"
  	depends_on = [
    	powermax_storagegroup.sg_vol_test
  	]
}
`, resourceVolName)

var VolumeConfigInvalidCYL = fmt.Sprintf(`
resource "powermax_volume" "volume_test" {
	vol_name = "%s"
	sg_name = powermax_storagegroup.sg_vol_test.name
	size = 0.5
	cap_unit = "CYL"
  	depends_on = [
    	powermax_storagegroup.sg_vol_test
  	]
}
`, resourceVolName)

var VolumeConfigWithCYL = fmt.Sprintf(`
resource "powermax_volume" "volume_test" {
	vol_name = "%s"
	size = 1
	cap_unit = "CYL"
	sg_name = powermax_storagegroup.sg_vol_test.name
  	depends_on = [
    	powermax_storagegroup.sg_vol_test
  	]
}
`, resourceVolName)

var VolumeConfigInvalidSG = fmt.Sprintf(`
resource "powermax_volume" "volume_test" {
	vol_name = "%s"
	sg_name = "invalid#SG"
	size = 0.5
	cap_unit = "CYL"
  	depends_on = [
    	powermax_storagegroup.sg_vol_test
  	]
}
`, resourceVolName)

var VolumeConfigInvalidName = `
resource "powermax_volume" "volume_test" {
	vol_name = "!@#$%"
	sg_name = powermax_storagegroup.sg_vol_test.name
	size = 0.5
	cap_unit = "CYL"
  	depends_on = [
    	powermax_storagegroup.sg_vol_test
  	]
}
`
