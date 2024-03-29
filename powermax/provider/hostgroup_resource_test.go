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

// Acceptance Tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccHostGroupResource(t *testing.T) {
	var hostGroupTerraformName = "powermax_hostgroup.test_hostgroup"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read test
			{
				Config: ProviderConfig + `
				resource "powermax_hostgroup" "test_hostgroup" {
					host_flags = {
						avoid_reset_broadcast = {
							enabled  = true
							override = true
						}
				  }
				  # This will be updated once host code is integrated to remove this from being hardcoded
				  host_ids = ["tfacc_host_group_host"]
				  name     = "test_host_group"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify host_flags are created and set correctly
					// avoid_reset_broadcast
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.avoid_reset_broadcast.enabled", "true"),
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.avoid_reset_broadcast.override", "true"),
					// disable_q_reset_on_ua
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.disable_q_reset_on_ua.enabled", "false"),
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.disable_q_reset_on_ua.override", "false"),
					// environ_set
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.environ_set.enabled", "false"),
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.environ_set.override", "false"),
					// openvms
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.openvms.enabled", "false"),
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.openvms.override", "false"),
					// scsi_3
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.scsi_3.enabled", "false"),
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.scsi_3.override", "false"),
					// spc2_protocol_version
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.spc2_protocol_version.enabled", "false"),
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.spc2_protocol_version.override", "false"),
					// scsi_support1
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.scsi_support1.enabled", "false"),
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.scsi_support1.override", "false"),
					// volume_set_addressing
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.volume_set_addressing.enabled", "false"),
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.volume_set_addressing.override", "false"),
					// Verify there is only 1 host attached
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_ids.#", "1"),
					// Verify the name
					resource.TestCheckResourceAttr(hostGroupTerraformName, "name", "test_host_group"),
					// Verify Calculated values
					// numofmaskingviews
					resource.TestCheckResourceAttr(hostGroupTerraformName, "numofmaskingviews", "0"),
					// numofinitiators
					resource.TestCheckResourceAttr(hostGroupTerraformName, "numofinitiators", "0"),
					// numofhosts
					resource.TestCheckResourceAttr(hostGroupTerraformName, "numofhosts", "1"),
				),
			},
			// Import testing
			{
				ResourceName:      hostGroupTerraformName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update testing
			{
				Config: ProviderConfig + `
				resource "powermax_hostgroup" "test_hostgroup" {
					host_flags = {
						avoid_reset_broadcast = {
							enabled  = false
							override = false
						}
						disable_q_reset_on_ua = {
							enabled  = true
							override = true
						}
						environ_set = {
							enabled  = true
							override = true
						}
						scsi_3 = {
							enabled  = true
							override = true
						}
						
				  }
				  # This will be updated once host code is integrated to remove this from being hardcoded
				  host_ids = ["tfacc_host_group_host", "tfacc_host_group_host_2"]
				  name     = "test_host_group_update"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify host_flags are created and set correctly
					// avoid_reset_broadcast
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.avoid_reset_broadcast.enabled", "false"),
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.avoid_reset_broadcast.override", "false"),
					// disable_q_reset_on_ua
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.disable_q_reset_on_ua.enabled", "true"),
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.disable_q_reset_on_ua.override", "true"),
					// environ_set
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.environ_set.enabled", "true"),
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.environ_set.override", "true"),
					// openvms
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.openvms.enabled", "false"),
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.openvms.override", "false"),
					// scsi_3
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.scsi_3.enabled", "true"),
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.scsi_3.override", "true"),
					// spc2_protocol_version
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.spc2_protocol_version.enabled", "false"),
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.spc2_protocol_version.override", "false"),
					// scsi_support1
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.scsi_support1.enabled", "false"),
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.scsi_support1.override", "false"),
					// volume_set_addressing
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.volume_set_addressing.enabled", "false"),
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_flags.volume_set_addressing.override", "false"),
					// Verify there is only 1 host attached
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_ids.#", "2"),
					// Verify the name
					resource.TestCheckResourceAttr(hostGroupTerraformName, "name", "test_host_group_update"),
					// Verify Calculated values
					// numofmaskingviews
					resource.TestCheckResourceAttr(hostGroupTerraformName, "numofmaskingviews", "0"),
					// numofinitiators
					resource.TestCheckResourceAttr(hostGroupTerraformName, "numofinitiators", "0"),
					// numofhosts
					resource.TestCheckResourceAttr(hostGroupTerraformName, "numofhosts", "2"),
				),
			},
			// auto checks delete to clean up the test
		},
	})
}

func TestAccHostGroupResourceEmptyHostIdInList(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + `
				resource "powermax_hostgroup" "test_hostgroup_create_err" {
					host_flags = {
						avoid_reset_broadcast = {
							enabled  = true
							override = true
						}
						disable_q_reset_on_ua = {
							enabled  = true
							override = true
						}
				  }
				  host_ids = [""]
				  name     = "tfacc_host_group_err"
				}
				`,
				ExpectError: regexp.MustCompile(`.*host_ids can not have an empty*.`),
			},
		},
	})
}

func TestAccHostGroupResourceNoHostFlagShouldStillWork(t *testing.T) {
	var hostGroupTerraformName = "powermax_hostgroup.test_hostgroup"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + `
				resource "powermax_hostgroup" "test_hostgroup" {
				  host_ids = ["tfacc_host_group_host"]
				  name     = "test_host_group_no_flag"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_ids.#", "1"),
					// Verify the name
					resource.TestCheckResourceAttr(hostGroupTerraformName, "name", "test_host_group_no_flag"),
					// Verify Calculated values
					// numofmaskingviews
					resource.TestCheckResourceAttr(hostGroupTerraformName, "numofmaskingviews", "0"),
					// numofinitiators
					resource.TestCheckResourceAttr(hostGroupTerraformName, "numofinitiators", "0"),
					// numofhosts
					resource.TestCheckResourceAttr(hostGroupTerraformName, "numofhosts", "1"),
				),
			},
		},
	})
}

func TestAccHostGroupResourceUpdateNoFlag(t *testing.T) {
	var hostGroupTerraformName = "powermax_hostgroup.test_hostgroup"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + `
				resource "powermax_hostgroup" "test_hostgroup" {
				  host_ids = ["tfacc_host_group_host"]
				  name     = "test_host_group_no_flag"
				}
				`,
			},
			{
				Config: ProviderConfig + `
				resource "powermax_hostgroup" "test_hostgroup" {
				  host_ids = ["tfacc_host_group_host", "tfacc_host_group_host_2"]
				  name     = "tfacc_host_group_update_no_flag"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(hostGroupTerraformName, "host_ids.#", "2"),
					// Verify the name
					resource.TestCheckResourceAttr(hostGroupTerraformName, "name", "tfacc_host_group_update_no_flag"),
					// Verify Calculated values
					// numofmaskingviews
					resource.TestCheckResourceAttr(hostGroupTerraformName, "numofmaskingviews", "0"),
					// numofinitiators
					resource.TestCheckResourceAttr(hostGroupTerraformName, "numofinitiators", "0"),
					// numofhosts
					resource.TestCheckResourceAttr(hostGroupTerraformName, "numofhosts", "2"),
				),
			},
		},
	})
}

func TestAccHostGroupResourceError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + `
				resource "powermax_hostgroup" "test_hostgroup_create_err" {
					host_flags = {
						avoid_reset_broadcast = {
							enabled  = true
							override = true
						}
						disable_q_reset_on_ua = {
							enabled  = true
							override = true
						}
				  }
				  host_ids = ["tfacc_host_group_host"]
				  name     = "tfacc_host_group_err"
				}
				`,
			},
			{
				Config: ProviderConfig + `
				resource "powermax_hostgroup" "test_hostgroup_err" {
					host_flags = {
						avoid_reset_broadcast = {
							enabled  = true
							override = true
						}
						disable_q_reset_on_ua = {
							enabled  = false
							override = false
						}
				  }
				  host_ids = ["tfacc_host_group_host"]
				  name     = "tfacc_host_group_err"
				}
				`,
				ExpectError: regexp.MustCompile(`.*Error creating hostgroup*.`),
			},
			{
				Config: ProviderConfig + `
				resource "powermax_hostgroup" "test_hostgroup_import_err" {
					host_flags = {
						avoid_reset_broadcast = {
							enabled  = true
							override = true
						}
				  }
				  host_ids = ["tfacc_host_group_host"]
				  name     = "tfacc_fake_host_group"
				}
				`,
				ResourceName:  "powermax_hostgroup.test_hostgroup_import_err",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`.*Error reading hostgroup*`),
				ImportStateId: "tfacc_fake_host_group",
			},
		},
	})
}
