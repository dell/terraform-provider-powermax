// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package test

// Acceptance Tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccHostResource(t *testing.T) {
	var hostTerraformName = "powermax_host.test_host_2"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read test
			{
				Config: ProviderConfig + `
				resource "powermax_host" "test_host_2" {
					host_flags = {
						avoid_reset_broadcast = {
							enabled  = true
							override = true
						}
				  }
				  name     = "test_host_2"
				  initiator = ["10000000c9959b8e"]
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
					resource.TestCheckResourceAttr(hostTerraformName, "name", "test_host_2"),
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
				resource "powermax_host" "test_host_2" {
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
				  name     = "test_host_2_update"
				  initiator = ["10000000c9959b8e"]
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
					resource.TestCheckResourceAttr(hostTerraformName, "name", "test_host_2_update"),
				),
			},
			// auto checks delete to clean up the test
		},
	})
}