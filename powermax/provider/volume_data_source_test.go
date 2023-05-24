// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package provider

import (
	"fmt"
	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"regexp"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/helper"
	"testing"
)

func TestAccVolumeDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				Config: ProviderConfig + VolStorageGroupDSConfig + VolumeDatasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powermax_volume.volume_datasource_test", "volumes.#", "1"),
					resource.TestCheckResourceAttr("data.powermax_volume.volume_datasource_test", "volumes.0.type", "TDEV"),
					resource.TestCheckResourceAttr("data.powermax_volume.volume_datasource_test", "volumes.0.status", "Ready"),
					resource.TestCheckResourceAttr("data.powermax_volume.volume_datasource_test", "filter.volume_identifier", datasourceVolName),
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
					FunctionMocker = Mock(helper.GetVolumeFilterParam).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + VolStorageGroupDSConfig + VolumeDatasourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			return nil
		},
	})
}

func TestAccVolumeDatasourceWithErrorVolList(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = Mock(GetMethod(client.Client{}.PmaxClient, "GetVolumeIDListWithParams")).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + VolStorageGroupDSConfig + VolumeDatasourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			return nil
		},
	})
}

var datasourceVolSGName = fmt.Sprintf("tfacc_ds_vol_sg_%s", ResourceSuffix)
var datasourceVolName = fmt.Sprintf("tfacc_ds_vol_%s", ResourceSuffix)

var VolStorageGroupDSConfig = fmt.Sprintf(`
resource "powermax_storagegroup" "sg_vol_ds_test" {
  name             = "%s"
  srp_id           = "SRP_1"
  slo              = "Gold"
}
`, datasourceVolSGName)

var VolumeDatasourceConfig = fmt.Sprintf(`
resource "powermax_volume" "volume_ds_resource_test" {
	sg_name = powermax_storagegroup.sg_vol_ds_test.name
	vol_name = "%s"
	size = 554
	cap_unit = "MB"
  	depends_on = [
    	powermax_storagegroup.sg_vol_ds_test
  	]
}

data "powermax_volume" "volume_datasource_test" {
	filter {
		storage_group_name = powermax_volume.volume_ds_resource_test.sg_name
		volume_identifier = powermax_volume.volume_ds_resource_test.vol_name
		num_of_storage_groups = "1"
		type = "TDEV"
	}
}
`, datasourceVolName)
