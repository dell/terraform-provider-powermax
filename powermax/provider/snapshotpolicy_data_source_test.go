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

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSnapshotPolicyDs(t *testing.T) {
	var snapshotPolicyTerraformName = "data.powermax_snapshotpolicy.SnapshotPolicyAll"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + snapshotPolicyAllDatasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(snapshotPolicyTerraformName, "snapshot_policies.#"),
				),
			},
		},
	})
}
func TestAccSnapshotPolicyDsFiltered(t *testing.T) {
	var snapshotPolicyTerraformName = "data.powermax_snapshotpolicy.SnapshotPolicyFiltered"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + snapshotPolicyDsFilteredConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(snapshotPolicyTerraformName, "snapshot_policies.#", "1"),
				),
			},
		},
	})
}

func TestAccSnapshotPolicyDsError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + snapshotPolicyDsError,
				ExpectError: regexp.MustCompile(`.*Error reading snapshot policy with id*.`),
			},
		},
	})
}

func TestAccSnapshotPolicyDsGetListError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetSnapshotPolicies).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + snapshotPolicyAllDatasourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccSnapshotPolicyDsGetSpecificPolicyError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetSnapshotPolicy).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + snapshotPolicyAllDatasourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccSnapshotPolicyDsMappingError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + snapshotPolicyAllDatasourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var snapshotPolicyDsFilteredConfig = `
 data "powermax_snapshotpolicy" "SnapshotPolicyFiltered" {
	filter {
	  # Optional list of IDs to filter
	  names = [
		"tfacc_snapshotPolicy1",
	  ]
	}
  }
`
var snapshotPolicyAllDatasourceConfig = `
 data "powermax_snapshotpolicy" "SnapshotPolicyAll" {
  }
`

var snapshotPolicyDsError = `
 data "powermax_snapshotpolicy" "SnapshotPolicyError" {
	filter {
	  # Optional list of IDs to filter
	  names = [
		"tfacc_snapshotPolicy_err",
	  ]
	}
  }
`
