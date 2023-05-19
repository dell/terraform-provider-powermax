// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMaskingView_FetchMaskingViewAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + maskingViewDataSourceparamsIDEmpty,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powermax_maskingview.all", "masking_views.#", "2"),
				),
			},
		},
	})
}
func TestAccMaskingView_FetchMaskingViewSingle(t *testing.T) {
	var maskingView = "data.powermax_maskingview.single"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + maskingViewDataSourceparamsID,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(maskingView, "masking_views.#", "1"),
					resource.TestCheckResourceAttr(maskingView, "masking_views.0.masking_view_name", "TestHostMaskingView"),
					resource.TestCheckResourceAttr(maskingView, "masking_views.0.capacity_gb", "10"),
					resource.TestCheckResourceAttr(maskingView, "masking_views.0.host_id", "IG_Dell_198151"),
					resource.TestCheckResourceAttr(maskingView, "masking_views.0.storage_group_id", "TestnewSG"),
					resource.TestCheckResourceAttr(maskingView, "masking_views.0.initiators.#", "1"),
					resource.TestCheckTypeSetElemAttr(maskingView, "masking_views.0.ports.*", "OR-2C:000"),
				),
			},
		},
	})
}
func TestAccMaskingView_FetchMaskingViewList(t *testing.T) {
	var maskingView = "data.powermax_maskingview.idList"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + maskingViewDataSourceparamsIDList,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(maskingView, "masking_views.#", "2"),
				),
			},
		},
	})
}

func TestAccMaskingView_FetchMaskingViewListFailed(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + maskingViewDataSourceparamsIDListInvalid,
				ExpectError: regexp.MustCompile(`.*Failed to get MaskingView*.`),
			},
		},
	})
}

var maskingViewDataSourceparamsID = `
data "powermax_maskingview" "single" {
	filter {
		names = ["TestHostMaskingView"]
	}
}

output "single" {
	value = data.powermax_maskingview.single
  }
`
var maskingViewDataSourceparamsIDList = `
data "powermax_maskingview" "idList" {
	filter {
		names = [ "TestHostMaskingView", "TestHostGroupMaskingView" ]
	}
}
`
var maskingViewDataSourceparamsIDListInvalid = `
data "powermax_maskingview" "idList" {
	filter {
		names = [ "TestHostMaskingView", "InvalidID" ]
	}
}
`
var maskingViewDataSourceparamsIDEmpty = `
data "powermax_maskingview" "all" {}

output "all" {
	value = data.powermax_maskingview.all
}
`
