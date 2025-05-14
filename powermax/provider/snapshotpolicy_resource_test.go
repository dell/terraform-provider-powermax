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
	"regexp"
	"terraform-provider-powermax/powermax/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var deleteMocker *mockey.Mocker

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
			// Add storage Groups to Snapshot Policy
			{
				Config: ProviderConfig + SnapshotPolicyResourceSGAdd,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "snapshot_policy_name", "terraform_test_sp_add"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "interval", "1 Day"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "compliance_count_critical", "29"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "compliance_count_warning", "47"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "type", "local"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "storage_groups.#", "2"),
				),
			},
			// Remove storage Groups from Snapshot Policy
			{
				Config: ProviderConfig + SnapshotPolicyResourceSGRemove,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "snapshot_policy_name", "terraform_test_sp_remove"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "interval", "1 Day"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "compliance_count_critical", "29"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "compliance_count_warning", "47"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "type", "local"),
					resource.TestCheckResourceAttr(snapPolicyTerraformName, "storage_groups.#", "0"),
				),
			},
			// Modify Error Check
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ModifySnapshotPolicy).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotPolicyResourceUpdateErr,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Read Policy Error Check
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = mockey.Mock(helper.GetSnapshotPolicy).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotPolicyResourceUpdateErr,
				ExpectError: regexp.MustCompile(`.*Error reading snapshot policy*.`),
			},
			// Read SG Error Check
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = mockey.Mock(helper.GetSnapshotPolicyStorageGroups).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotPolicyResourceUpdateErr,
				ExpectError: regexp.MustCompile(`.*Error getting snapshot policy storage groups*.`),
			},
			// Read Mapping Error Check
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = mockey.Mock(helper.UpdateSnapshotPolicyResourceState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotPolicyResourceUpdateErr,
				ExpectError: regexp.MustCompile(`.*Error reading snapshot policy*.`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccSnapshotPolicyResourceError(t *testing.T) {
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

func TestAccSnapshotPolicyResourceCreateError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					deleteMocker = mockey.Mock(helper.DeleteSnapshotPolicy).Return(nil, fmt.Errorf("mock error")).Build()
					FunctionMocker = mockey.Mock(helper.CreateSnapshotPolicy).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotPolicyResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Do the delete successfully to cleanup after the test
			{
				PreConfig: func() {
					if deleteMocker != nil {
						deleteMocker.UnPatch()
					}
				},
				Config:      ProviderConfig + SnapshotPolicyResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccSnapshotPolicyResourceSgError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetSnapshotPolicyStorageGroups).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotPolicyResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccSnapshotPolicyResourceMapperError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.UpdateSnapshotPolicyResourceState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotPolicyResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
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
var SnapshotPolicyResourceUpdateErr = `
resource "powermax_snapshotpolicy" "terraform_test_sp" {
	snapshot_policy_name = "terraform_test_sp3"
	interval = "1 Day"
	compliance_count_critical = 29
  }
`
var SnapshotPolicyResourceSGAdd = `
resource "powermax_snapshotpolicy" "terraform_test_sp" {
	snapshot_policy_name = "terraform_test_sp_add"
	interval = "1 Day"
	compliance_count_critical = 29
	storage_groups = ["tfacc_sp_sg1", "tfacc_sp_sg2"]
  }
`
var SnapshotPolicyResourceSGRemove = `
resource "powermax_snapshotpolicy" "terraform_test_sp" {
	snapshot_policy_name = "terraform_test_sp_remove"
	interval = "1 Day"
	compliance_count_critical = 29
	storage_groups = []
  }
`
