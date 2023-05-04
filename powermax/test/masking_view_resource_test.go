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
				Config: ProviderConfig + maskingViewCreateWithHostTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "name", "TestHostMaskingView"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "storage_group_id", "TestnewSG"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "host_id", "IG_Dell_198151"),
				),
			},
		},
	})
}
func TestAccMaskingView_CreateMaskingViewWithHostFailTest(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + maskingViewCreateFailedTest,
				ExpectError: regexp.MustCompile(`.*Specify either host_id or host_group_id*.`),
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
				Config: ProviderConfig + maskingViewCreateWithHostGroupTest,
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
				Config: ProviderConfig + maskingViewUpdateRenameTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "name", "maskingViewUpdate"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "storage_group_id", "esa_sg572"),
				),
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
				Config: ProviderConfig + maskingViewCreateWithHostTest,
			},
			{
				Config:            ProviderConfig + maskingViewCreateWithHostTest,
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
				Config:        ProviderConfig + maskingViewCreateWithHostTest,
				ResourceName:  "powermax_maskingview.masking_view_create_with_host_test",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`.*Error reading masking view*`),
				ImportStateId: "invalid-id",
			},
		},
	})
}

var maskingViewCreateWithHostTest = `
resource "powermax_maskingview" "masking_view_create_with_host_test" {
	name = "TestHostMaskingView"
	storage_group_id = "TestnewSG"
	host_id = "IG_Dell_198151"
	host_group_id = ""
	port_group_id = "TestnewSG_PG"
  }
`

var maskingViewCreateFailedTest = `
resource "powermax_maskingview" "masking_view_create_failed_test" {
	name = "terraform_MV_accTest_host"
	storage_group_id = "Tao_k8s_env2_SG"
	host_id = "Tao_k8s_env2_host"
	host_group_id = "test"
	port_group_id = "Tao_k8s_env2_PG"
  }
`

var maskingViewCreateWithHostGroupTest = `
resource "powermax_maskingview" "masking_view_create_with_host_group_test" {
	name = "TestHostGroupMaskingView"
	storage_group_id = "TestnewSG"
	host_id = ""
	host_group_id = "TestHostGroup"
	port_group_id = "TestnewSG_PG"
  }
`

var maskingViewUpdateRenameTest = `
resource "powermax_maskingview" "masking_view_create_with_host_test" {
	name = "maskingViewUpdate"
	storage_group_id = "esa_sg572"
	host_id = "Host173"
	host_group_id = ""
	port_group_id = "esa_vmax_portgroup572"
  }
`

var maskingViewUpdateFailedTest = `
resource "powermax_maskingview" "masking_view_create_with_host_test" {
	name = "terraform_MV_accTest_host_rename"
	storage_group_id = "Tao_k8s_env2_SG"
	host_id = ""
	host_group_id = "Tao_k8s_env2_host_group"
	port_group_id = "Tao_k8s_env2_PG"
  }
`
