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

func TestAccSnapshotPolicyResource(t *testing.T) {
	var snapPolicyTerraformName = "powermax_snapshotpolicy.terraform_test_sp"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + SnapshotPolicyResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "snapshot_policy_name", "terraform_test_sp"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "interval", "7 Days"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "compliance_count_critical", "29"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "compliance_count_warning", "47"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "type", "local"),
				),
			},
			// ImportState testing
			{
				ResourceName:      snapPolicyTerraformName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update
			{
				Config: ProviderConfig + SnapshotPolicyResourceUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "snapshot_policy_name", "terraform_test_sp1"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "interval", "1 Day"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "compliance_count_critical", "29"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "compliance_count_warning", "47"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "type", "local"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSnapshotPolicyResourceError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + `
				resource "powermax_snapshotpolicy" "terraform_test_sp_error" {
					snapshot_policy_name = "terraform_test_error"
					interval = "1 Minute"
				  }
				`,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value Match*.`),
			},
		},
	})
}

var SnapshotPolicyResourceConfig = `
resource "powermax_snapshotpolicy" "terraform_test_sp" {
	snapshot_policy_name = "terraform_test_sp"
	interval = "7 Days"
	compliance_count_critical = 29
  }
`
var SnapshotPolicyResourceUpdate = `
resource "powermax_snapshotpolicy" "terraform_test_sp" {
	snapshot_policy_name = "terraform_test_sp1"
	interval = "1 Day"
	compliance_count_critical = 29
  }
`
