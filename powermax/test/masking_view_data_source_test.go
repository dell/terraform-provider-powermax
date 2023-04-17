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
				Config: ProviderConfig + maskingViewDataSourceparamsID,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powermax_maskingview.id", "masking_views.#", "1"),
					resource.TestCheckResourceAttr("data.powermax_maskingview.id", "masking_views.0.masking_view_id", "Tao_k8s_env2_SG_MV"),
					resource.TestCheckResourceAttr("data.powermax_maskingview.id", "masking_views.0.capacity_gb", "10"),
					resource.TestCheckResourceAttr("data.powermax_maskingview.id", "masking_views.0.host_id", "Tao_k8s_env2_host"),
					resource.TestCheckResourceAttr("data.powermax_maskingview.id", "masking_views.0.storage_group_id", "Tao_k8s_env2_SG"),
					resource.TestCheckResourceAttr("data.powermax_maskingview.id", "masking_views.0.initiators.#", "3"),
					resource.TestCheckTypeSetElemAttr("data.powermax_maskingview.id", "masking_views.0.ports.*", "OR-2C:000"),
				),
			},
			{
				Config: ProviderConfig + maskingViewDataSourceparamsIDList,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powermax_maskingview.idList", "masking_views.#", "2"),
				),
			},
			{
				Config: ProviderConfig + maskingViewDataSourceparamsIDEmpty,
			},
		},
	})
}

var maskingViewDataSourceparamsID = `
data "powermax_maskingview" "id" {
	id = "Tao_k8s_env2_SG_MV"
}
`

var maskingViewDataSourceparamsIDList = `
data "powermax_maskingview" "idList" {
	masking_view_ids = [ "csi-mv-YL1-worker-1-rvahi2ntosoe3", "Yulan_SG_Yiming_10838_MV" ]
}
`
var maskingViewDataSourceparamsIDEmpty = `
data "powermax_maskingview" "all" {}
`
