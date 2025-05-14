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
	"dell/powermax-go-client"
	"fmt"
	"regexp"
	"terraform-provider-powermax/powermax/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccMaskingViewResourceCreateMaskingViewWithHost(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + maskingViewCreateFailed,
				ExpectError: regexp.MustCompile(`.*Specify either host_id or host_group_id*.`),
			},
			{
				Config:      ProviderConfig + maskingViewCreateError,
				ExpectError: regexp.MustCompile(`.*Error creating masking view*.`),
			},
			{
				Config: ProviderConfig + maskingViewCreateWithHost,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "name", "tfacc_masking_view"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "storage_group_id", "tfacc_masking_view_sg"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "host_id", "tfacc_masking_view_host"),
				),
			},
		},
	})
}

func TestAccMaskingViewResourceCreateMaskingViewWithHostGroup(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + maskingViewCreateWithHostGroup,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_group_test", "name", "tfacc_masking_view_hg"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_group_test", "storage_group_id", "tfacc_masking_view_sg"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_group_test", "host_group_id", "tfacc_masking_view_hg"),
				),
			},
		},
	})
}

func TestAccMaskingViewResourceUpdateMaskingView(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + maskingViewCreateWithHost,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "name", "tfacc_masking_view"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "storage_group_id", "tfacc_masking_view_sg"),
				),
			},
			{
				Config:      ProviderConfig + maskingViewUpdateError,
				ExpectError: regexp.MustCompile(`.*Error renaming masking view*.`),
			},
			{
				Config: ProviderConfig + maskingViewUpdateRename,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "name", "tfacc_masking_view_update"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "storage_group_id", "tfacc_masking_view_sg"),
				),
			},
			{
				Config:      ProviderConfig + maskingViewUpdateFailed,
				ExpectError: regexp.MustCompile(`.*maskingView's host, hostGroup, portGroup or storageGroup cannot be update after creation*.`),
			},
		},
	})
}

func TestAccMaskingViewResourceCreateErrors(t *testing.T) {
	idError := "someErrorId"
	sgID := "someSg"
	hostGroupErrorResponse := powermax.MaskingView{
		HostId:         &idError,
		StorageGroupId: &sgID,
	}
	hostErrorResponse := powermax.MaskingView{
		HostGroupId:    &idError,
		StorageGroupId: &sgID,
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.CreateMaskingView).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + maskingViewCreateWithHostGroup,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + maskingViewCreateWithHostGroup,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = mockey.Mock(helper.GetMaskingView).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + maskingViewCreateWithHostGroup,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = mockey.Mock(helper.GetMaskingView).Return(&hostGroupErrorResponse, nil, nil).Build()
				},
				Config:      ProviderConfig + maskingViewCreateWithHostGroup,
				ExpectError: regexp.MustCompile(`.*Error creating masking view*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = mockey.Mock(helper.GetMaskingView).Return(&hostErrorResponse, nil, nil).Build()
				},
				Config:      ProviderConfig + maskingViewCreateWithHost,
				ExpectError: regexp.MustCompile(`.*Error creating masking view*.`),
			},
		},
	})
}

func TestAccMaskingViewResourceImportSuccess(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + maskingViewCreateWithHost,
			},
			{
				Config:       ProviderConfig + maskingViewCreateWithHost,
				ResourceName: "powermax_maskingview.masking_view_create_with_host_test",
				ImportState:  true,
				ExpectError:  nil,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					assert.Equal(t, "tfacc_masking_view", s[0].Attributes["name"])
					assert.Equal(t, "tfacc_masking_view_host", s[0].Attributes["host_id"])
					return nil
				},
			},
		},
	})
}

func TestAccMaskingViewResourceImportFailure(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        ProviderConfig + maskingViewCreateWithHost,
				ResourceName:  "powermax_maskingview.masking_view_create_with_host_test",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`.*Error reading masking view*`),
				ImportStateId: "invalid-id",
			},
		},
	})
}

var maskingViewCreateWithHost = `
resource "powermax_maskingview" "masking_view_create_with_host_test" {
	name = "tfacc_masking_view"
	storage_group_id = "tfacc_masking_view_sg"
	host_id = "tfacc_masking_view_host"
	host_group_id = ""
	port_group_id = "tfacc_masking_view_pg"
  }
`

var maskingViewCreateFailed = `
resource "powermax_maskingview" "masking_view_create_failed_test" {
	name = "tfacc_masking_view"
	storage_group_id = "tfacc_masking_view_sg"
	host_id = "tfacc_masking_view_host"
	host_group_id = "tfacc_masking_view_hg"
	port_group_id = "tfacc_masking_view_pg"
  }
`

var maskingViewCreateError = `
resource "powermax_maskingview" "masking_view_create_failed_test" {
	name = "tfacc_masking_view_ds"
	storage_group_id = "tfacc_masking_view_sg"
	host_id = "tfacc_masking_view_host"
	host_group_id = ""
	port_group_id = "tfacc_masking_view_pg"
  }
`

var maskingViewCreateWithHostGroup = `
resource "powermax_maskingview" "masking_view_create_with_host_group_test" {
	name = "tfacc_masking_view_hg"
	storage_group_id = "tfacc_masking_view_sg"
	host_id = ""
	host_group_id = "tfacc_masking_view_hg"
	port_group_id = "tfacc_masking_view_pg"
  }
`

var maskingViewUpdateRename = `
resource "powermax_maskingview" "masking_view_create_with_host_test" {
	name = "tfacc_masking_view_update"
	storage_group_id = "tfacc_masking_view_sg"
	host_id = "tfacc_masking_view_host"
	host_group_id = ""
	port_group_id = "tfacc_masking_view_pg"
  }
`

var maskingViewUpdateFailed = `
resource "powermax_maskingview" "masking_view_create_with_host_test" {
	name = "tfacc_masking_view_update"
	storage_group_id = "tfacc_masking_view_sg_update"
	host_id = "tfacc_masking_view_host"
	host_group_id = ""
	port_group_id = "tfacc_masking_view_pg"
  }
`

var maskingViewUpdateError = `
resource "powermax_maskingview" "masking_view_create_with_host_test" {
	name = "tfacc_masking_view_ds"
	storage_group_id = "tfacc_masking_view_sg"
	host_id = "tfacc_masking_view_host"
	host_group_id = ""
	port_group_id = "tfacc_masking_view_pg"
  }
`
