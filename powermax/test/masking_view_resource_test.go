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
				Config:      ProviderConfig + masking_view_create_failed_test,
				ExpectError: regexp.MustCompile(`.*The host_id or host_group_id only needs to be specified one.*.`),
			},
			{
				Config: ProviderConfig + masking_view_create_with_host_test,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "name", "terraform_MV_accTest_host"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "storage_group_id", "Tao_k8s_env2_SG"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "host_id", "Tao_k8s_env2_host"),
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
				Config: ProviderConfig + masking_view_create_with_host_group_test,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_group_test", "name", "terraform_MV_accTest_hostGroup"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_group_test", "storage_group_id", "Tao_k8s_env2_SG"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_group_test", "host_group_id", "Tao_k8s_env2_host_group"),
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
				Config: ProviderConfig + masking_view_create_with_host_test,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "name", "terraform_MV_accTest_host"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "storage_group_id", "Tao_k8s_env2_SG"),
				),
			},
			{
				Config:      ProviderConfig + masking_view_update_failed_test,
				ExpectError: regexp.MustCompile(`.*maskingView's host, hostGroup, portGroup or storageGroup cannot be update after creation*.`),
			},
			{
				Config: ProviderConfig + masking_view_update_rename_test,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "name", "terraform_MV_accTest_host_rename"),
					resource.TestCheckResourceAttr("powermax_maskingview.masking_view_create_with_host_test", "storage_group_id", "Tao_k8s_env2_SG"),
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
				Config: ProviderConfig + masking_view_create_with_host_test,
			},
			{
				Config:            ProviderConfig + masking_view_create_with_host_test,
				ResourceName:      "powermax_maskingview.masking_view_create_with_host_test",
				ImportState:       true,
				ExpectError:       nil,
				ImportStateVerify: true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					assert.Equal(t, "terraform_MV_accTest_host", s[0].Attributes["name"])
					assert.Equal(t, "Tao_k8s_env2_host", s[0].Attributes["host_id"])
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
				Config:        ProviderConfig + masking_view_create_with_host_test,
				ResourceName:  "powermax_maskingview.masking_view_create_with_host_test",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`.*Error reading masking view*`),
				ImportStateId: "invalid-id",
			},
		},
	})
}

var masking_view_create_with_host_test = `
resource "powermax_maskingview" "masking_view_create_with_host_test" {
	name = "terraform_MV_accTest_host"
	storage_group_id = "Tao_k8s_env2_SG"
	host_id = "Tao_k8s_env2_host"
	host_group_id = ""
	port_group_id = "Tao_k8s_env2_PG"
  }
`

var masking_view_create_failed_test = `
resource "powermax_maskingview" "masking_view_create_failed_test" {
	name = "terraform_MV_accTest_host"
	storage_group_id = "Tao_k8s_env2_SG"
	host_id = "Tao_k8s_env2_host"
	host_group_id = "test"
	port_group_id = "Tao_k8s_env2_PG"
  }
`

var masking_view_create_with_host_group_test = `
resource "powermax_maskingview" "masking_view_create_with_host_group_test" {
	name = "terraform_MV_accTest_hostGroup"
	storage_group_id = "Tao_k8s_env2_SG"
	host_id = ""
	host_group_id = "Tao_k8s_env2_host_group"
	port_group_id = "Tao_k8s_env2_PG"
  }
`

var masking_view_update_rename_test = `
resource "powermax_maskingview" "masking_view_create_with_host_test" {
	name = "terraform_MV_accTest_host_rename"
	storage_group_id = "Tao_k8s_env2_SG"
	host_id = "Tao_k8s_env2_host"
	host_group_id = ""
	port_group_id = "Tao_k8s_env2_PG"
  }
`

var masking_view_update_failed_test = `
resource "powermax_maskingview" "masking_view_create_with_host_test" {
	name = "terraform_MV_accTest_host_rename"
	storage_group_id = "Tao_k8s_env2_SG"
	host_id = ""
	host_group_id = "Tao_k8s_env2_host_group"
	port_group_id = "Tao_k8s_env2_PG"
  }
`
