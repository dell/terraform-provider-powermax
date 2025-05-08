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

func TestAccVolumeDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				Config: ProviderConfig + VolumeDatasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powermax_volume.volume_datasource_test", "volumes.#", "1"),
					resource.TestCheckResourceAttr("data.powermax_volume.volume_datasource_test", "volumes.0.type", "TDEV"),
					resource.TestCheckResourceAttr("data.powermax_volume.volume_datasource_test", "volumes.0.status", "Ready"),
					resource.TestCheckResourceAttr("data.powermax_volume.volume_datasource_test", "filter.volume_identifier", "tfacc_ds_vol_sOCTK"),
				),
			},
		},
	})
}

func TestAccVolumeDatasourceInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetVolumeFilterParam).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + VolumeDatasourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccVolumeDatasourceErrorUpdatingState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.UpdateVolumeState).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + VolumeDatasourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var VolumeDatasourceConfig = `

data "powermax_volume" "volume_datasource_test" {
	filter {
		storage_group_name = "tfacc_ds_vol_sg_sOCTK"
		volume_identifier = "tfacc_ds_vol_sOCTK"
		num_of_storage_groups = "1"
		type = "TDEV"
	}
}
`
