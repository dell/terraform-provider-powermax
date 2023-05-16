// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.

package test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccMaskingView_CreateMaskingViewWithHost(t *testing.T) {

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
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "name", "TestHostMaskingView"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "storage_group_id", "TestnewSG"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "host_id", "IG_Dell_198151"),
				),
			},
		},
	})
}

func TestAccMaskingView_CreateMaskingViewWithHostGroup(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + maskingViewCreateWithHostGroup,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_group_test", "name", "TestHostGroupMaskingView"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_group_test", "storage_group_id", "TestnewSG"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_group_test", "host_group_id", "TestHostGroup"),
				),
			},
		},
	})
}

func TestAccMaskingView_UpdateMaskingView(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + maskingViewCreateWithHost,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "name", "TestHostMaskingView"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "storage_group_id", "TestnewSG"),
				),
			},
			{
				Config:      ProviderConfig + maskingViewUpdateError,
				ExpectError: regexp.MustCompile(`.*Error renaming masking view*.`),
			},
			{
				Config: ProviderConfig + maskingViewUpdateRename,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "name", "maskingViewUpdate"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "storage_group_id", "TestnewSG"),
				),
			},
			{
				Config:      ProviderConfig + maskingViewUpdateFailed,
				ExpectError: regexp.MustCompile(`.*maskingView's host, hostGroup, portGroup or storageGroup cannot be update after creation*.`),
			},
		},
	})
}

func TestAccMaskingView_ImportSuccess(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + maskingViewCreateWithHost,
			},
			{
				Config:            ProviderConfig + maskingViewCreateWithHost,
				ResourceName:      "powermax_maskingview.masking_view_create_with_host_test",
				ImportState:       true,
				ExpectError:       nil,
				ImportStateVerify: true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					assert.Equal(t, "TestHostMaskingView", s[0].Attributes["name"])
					assert.Equal(t, "IG_Dell_198151", s[0].Attributes["host_id"])
					return nil
				},
			},
		},
	})
}

func TestAccMaskingView_ImportFailure(t *testing.T) {

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
	name = "TestHostMaskingView"
	storage_group_id = "TestnewSG"
	host_id = "IG_Dell_198151"
	host_group_id = ""
	port_group_id = "TestnewSG_PG"
  }
`

var maskingViewCreateFailed = `
resource "powermax_maskingview" "masking_view_create_failed_test" {
	name = "TestHostMaskingView"
	storage_group_id = "TestnewSG"
	host_id = "IG_Dell_198151"
	host_group_id = "test"
	port_group_id = "TestnewSG_PG"
  }
`

var maskingViewCreateError = `
resource "powermax_maskingview" "masking_view_create_failed_test" {
	name = "CreateMaskingViewError"
	storage_group_id = "TestnewSG"
	host_id = "IG_Dell_198151"
	host_group_id = ""
	port_group_id = "TestnewSG_PG"
  }
`

var maskingViewCreateWithHostGroup = `
resource "powermax_maskingview" "masking_view_create_with_host_group_test" {
	name = "TestHostGroupMaskingView"
	storage_group_id = "TestnewSG"
	host_id = ""
	host_group_id = "TestHostGroup"
	port_group_id = "TestnewSG_PG"
  }
`

var maskingViewUpdateRename = `
resource "powermax_maskingview" "masking_view_create_with_host_test" {
	name = "maskingViewUpdate"
	storage_group_id = "TestnewSG"
	host_id = "IG_Dell_198151"
	host_group_id = ""
	port_group_id = "TestnewSG_PG"
  }
`

var maskingViewUpdateFailed = `
resource "powermax_maskingview" "masking_view_create_with_host_test" {
	name = "maskingViewUpdate"
	storage_group_id = "TestnewSG_rename"
	host_id = "IG_Dell_198151"
	host_group_id = ""
	port_group_id = "TestnewSG_PG"
  }
`

var maskingViewUpdateError = `
resource "powermax_maskingview" "masking_view_create_with_host_test" {
	name = "maskingViewUpdateError"
	storage_group_id = "TestnewSG"
	host_id = "IG_Dell_198151"
	host_group_id = ""
	port_group_id = "TestnewSG_PG"
  }
`
