// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.

package provider

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
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "name", "tfacc_masking_view"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "storage_group_id", "tfacc_masking_view_sg"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "host_id", "tfacc_masking_view_host"),
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
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_group_test", "name", "tfacc_masking_view_hg"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_group_test", "storage_group_id", "tfacc_masking_view_sg"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_group_test", "host_group_id", "tfacc_masking_view_hg"),
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
					assert.Equal(t, "tfacc_masking_view", s[0].Attributes["name"])
					assert.Equal(t, "tfacc_masking_view_host", s[0].Attributes["host_id"])
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
