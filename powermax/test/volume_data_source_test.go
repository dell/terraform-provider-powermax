// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package test

import (
	"testing"

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
					resource.TestCheckResourceAttr("data.powermax_volume.volume_datasource_test", "volumes.0.id", "00160"),
					resource.TestCheckResourceAttr("data.powermax_volume.volume_datasource_test", "volumes.0.type", "TDEV"),
					resource.TestCheckResourceAttr("data.powermax_volume.volume_datasource_test", "filter.status", "Ready"),
					resource.TestCheckResourceAttr("data.powermax_volume.volume_datasource_test", "filter.volume_identifier", "test_acc_create_volume_1"),
				),
			},
		},
	})
}

var VolumeDatasourceConfig = `
data "powermax_volume" "volume_datasource_test" {
	filter {
		storage_group_name = "terraform_vol_sg"
		wwn = "60000970000197902572533030313630"
		status = "Ready"
		volume_identifier = "test_acc_create_volume_1"
		allocated_percent = "0"
		num_of_storage_groups = "1"
		num_of_masking_views = "0"
		num_of_front_end_paths = "0"
		mobility_id_enabled = false
		virtual_volumes = true
		private_volumes = false
		tdev = true
		vdev = false
		available_thin_volumes = false
		gatekeeper = false
		data_volume = false
		dld = false
		drv = false
		encapsulated = false
		associated = false
		reserved = false
		pinned = false
		mapped = false
		bound_tdev = true
		emulation = "FBA"
		has_effective_wwn = false
		effective_wwn = "60000970000197902572533030313630"
		type = "TDEV"
		unreducible_data_gb = "0"
	}
}
`
