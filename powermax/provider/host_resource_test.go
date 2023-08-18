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
	"fmt"
	"regexp"
	"terraform-provider-powermax/powermax/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccHostResource(t *testing.T) {
	var hostTerraformName = "powermax_host.Test_Host"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read test
			{
				Config: ProviderConfig + `
				resource "powermax_host" "Test_Host" {
					host_flags = {
						avoid_reset_broadcast = {
							enabled  = true
							override = true
						}
				  }
				  name     = "tfacc_host_test_cr"
				  initiator = ["21000024ff3efed6"]
				  consistent_lun = false
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify host_flags are created and set correctly
					// avoid_reset_broadcast
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.avoid_reset_broadcast.enabled", "true"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.avoid_reset_broadcast.override", "true"),
					// disable_q_reset_on_ua
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.disable_q_reset_on_ua.enabled", "false"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.disable_q_reset_on_ua.override", "false"),
					// environ_set
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.environ_set.enabled", "false"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.environ_set.override", "false"),
					// openvms
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.openvms.enabled", "false"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.openvms.override", "false"),
					// scsi_3
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.scsi_3.enabled", "false"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.scsi_3.override", "false"),
					// spc2_protocol_version
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.spc2_protocol_version.enabled", "false"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.spc2_protocol_version.override", "false"),
					// scsi_support1
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.scsi_support1.enabled", "false"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.scsi_support1.override", "false"),
					// volume_set_addressing
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.volume_set_addressing.enabled", "false"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.volume_set_addressing.override", "false"),
					// Verify the name
					resource.TestCheckResourceAttr(hostTerraformName, "name", "tfacc_host_test_cr"),
				),
			},
			// Import testing
			{
				ResourceName:      hostTerraformName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update testing
			{
				Config: ProviderConfig + `
				resource "powermax_host" "Test_Host" {
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
				  name     = "tfacc_host_test_up"
				  initiator = []
				  consistent_lun = true
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify host_flags are created and set correctly
					// avoid_reset_broadcast
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.avoid_reset_broadcast.enabled", "false"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.avoid_reset_broadcast.override", "false"),
					// disable_q_reset_on_ua
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.disable_q_reset_on_ua.enabled", "true"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.disable_q_reset_on_ua.override", "true"),
					// environ_set
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.environ_set.enabled", "true"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.environ_set.override", "true"),
					// openvms
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.openvms.enabled", "false"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.openvms.override", "false"),
					// scsi_3
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.scsi_3.enabled", "true"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.scsi_3.override", "true"),
					// spc2_protocol_version
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.spc2_protocol_version.enabled", "false"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.spc2_protocol_version.override", "false"),
					// scsi_support1
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.scsi_support1.enabled", "false"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.scsi_support1.override", "false"),
					// volume_set_addressing
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.volume_set_addressing.enabled", "false"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.volume_set_addressing.override", "false"),
					// Verify the name
					resource.TestCheckResourceAttr(hostTerraformName, "name", "tfacc_host_test_up"),
					// Verify Consistent_lun flag
					resource.TestCheckResourceAttr(hostTerraformName, "consistent_lun", "true"),
				),
			},
			{
				Config: ProviderConfig + `
				resource "powermax_host" "Test_Host" {
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
				  name     = "tfacc_host_test_up"
				  initiator = []
				  consistent_lun = true
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify host_flags are created and set correctly
					// avoid_reset_broadcast
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.avoid_reset_broadcast.enabled", "false"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.avoid_reset_broadcast.override", "false"),
					// disable_q_reset_on_ua
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.disable_q_reset_on_ua.enabled", "true"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.disable_q_reset_on_ua.override", "true"),
					// environ_set
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.environ_set.enabled", "true"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.environ_set.override", "true"),
					// openvms
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.openvms.enabled", "false"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.openvms.override", "false"),
					// scsi_3
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.scsi_3.enabled", "true"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.scsi_3.override", "true"),
					// spc2_protocol_version
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.spc2_protocol_version.enabled", "false"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.spc2_protocol_version.override", "false"),
					// scsi_support1
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.scsi_support1.enabled", "false"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.scsi_support1.override", "false"),
					// volume_set_addressing
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.volume_set_addressing.enabled", "false"),
					resource.TestCheckResourceAttr(hostTerraformName, "host_flags.volume_set_addressing.override", "false"),
					// Verify the name
					resource.TestCheckResourceAttr(hostTerraformName, "name", "tfacc_host_test_up"),
					// Verify Consistent_lun flag
					resource.TestCheckResourceAttr(hostTerraformName, "consistent_lun", "true"),
				),
			},
			// auto checks delete to clean up the test
		},
	})
}

func TestAccHostResourceReadError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read test
			{
				Config: ProviderConfig + `
				resource "powermax_host" "Test_Host" {
					host_flags = {
						avoid_reset_broadcast = {
							enabled  = true
							override = true
						}
				  }
				  name     = "tfacc_host_test_cr"
				  initiator = ["21000024ff3efed6"]
				  consistent_lun = false
				}
				`,
			},
			// Read error
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetHost).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config: ProviderConfig + `
				resource "powermax_host" "Test_Host" {
					host_flags = {
						avoid_reset_broadcast = {
							enabled  = true
							override = true
						}
				  }
				  name     = "tfacc_host_test_cr"
				  initiator = []
				  consistent_lun = false
				}
				`,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccHostResourceModifyError(t *testing.T) {
	var errorString []string
	errorString = append(errorString, "mock error")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read test
			{
				Config: ProviderConfig + `
				resource "powermax_host" "Test_Host" {
					host_flags = {
						avoid_reset_broadcast = {
							enabled  = true
							override = true
						}
				  }
				  name     = "tfacc_host_test_cr"
				  initiator = ["21000024ff3efed6"]
				  consistent_lun = false
				}
				`,
			},
			// Modify error
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.UpdateHost).Return(nil, nil, errorString).Build()
				},
				Config: ProviderConfig + `
				resource "powermax_host" "Test_Host" {
					host_flags = {
						avoid_reset_broadcast = {
							enabled  = true
							override = true
						}
				  }
				  name     = "tfacc_host_test_cr"
				  initiator = []
				  consistent_lun = false
				}
				`,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccHostResourceCreateReadError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read Error
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetHost).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config: ProviderConfig + `
				resource "powermax_host" "Test_Host" {
					host_flags = {
						avoid_reset_broadcast = {
							enabled  = true
							override = true
						}
				  }
				  name     = "tfacc_host_test_cr"
				  initiator = ["21000024ff3efed6"]
				  consistent_lun = false
				}
				`,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccHostResourceError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + `
				resource "powermax_host" "test_host_err" {
					host_flags = {
						avoid_reset_broadcast = {
							enabled  = true
							override = true
						}
				  }
				  name     = "non-existent-host"
				  initiator = [""]
				}
				`,
				ExpectError: regexp.MustCompile(`.*Could not create host*.`),
			},
		},
	})
}

func TestAccHostResourceImportError(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + `
				resource "powermax_host" "test_host_err" {
					host_flags = {
						avoid_reset_broadcast = {
							enabled  = true
							override = true
						}
				  }
				  name     = "non-existent-host"
				  initiator = ["10000000c9959b8e"]
				}
				`,
				ResourceName:  "powermax_host.test_host_err",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`.*Error reading host*`),
				ImportStateId: "non-existent-host",
			},
		},
	})
}
