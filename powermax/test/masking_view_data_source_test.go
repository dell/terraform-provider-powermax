// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMaskingView_FetchMaskingView(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + maskingViewDataSourceFilter,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powermax_maskingview.maskingViewFilter", "masking_views.#", "1"),
					resource.TestCheckResourceAttr("data.powermax_maskingview.maskingViewFilter", "masking_views.0.masking_view_name", "Yulan_SG_MV"),
					resource.TestCheckResourceAttr("data.powermax_maskingview.maskingViewFilter", "masking_views.0.capacity_gb", "10"),
					resource.TestCheckResourceAttr("data.powermax_maskingview.maskingViewFilter", "masking_views.0.host_id", "csi-node-YL1-worker-1-onenvvkgfbkhm"),
					resource.TestCheckResourceAttr("data.powermax_maskingview.maskingViewFilter", "masking_views.0.storage_group_id", "Yulan_SG"),
					resource.TestCheckResourceAttr("data.powermax_maskingview.maskingViewFilter", "masking_views.0.initiators.#", "1"),
					resource.TestCheckTypeSetElemAttr("data.powermax_maskingview.maskingViewFilter", "masking_views.0.ports.*", "OR-2C:000"),
				),
			},
			{
				Config: ProviderConfig + maskingViewDataSourceparamsIDEmpty,
			},
		},
	})
}

var maskingViewDataSourceFilter = `
data "powermax_maskingview" "maskingViewFilter" {
	filter {
	 names = ["Yulan_SG_MV"]
   }
 }
`

var maskingViewDataSourceparamsIDEmpty = `
data "powermax_maskingview" "all" {}
`
