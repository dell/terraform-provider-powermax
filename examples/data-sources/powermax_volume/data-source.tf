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
data "powermax_volume" "volume_datasource_test" {
  filter {
    storage_group_name     = "terraform_vol_sg"
    wwn                    = "wwn_num"
    encapsulated_wwn       = "encapsulated_wwn_num"
    symmlun                = "symmlun"
    status                 = "Ready"
    physical_name          = "physical_name"
    volume_identifier      = "test_acc_create_volume"
    cap_tb                 = "0"
    cap_gb                 = "0"
    cap_mb                 = "0"
    cap_cyl                = "0"
    allocated_percent      = "0"
    num_of_storage_groups  = "1"
    num_of_masking_views   = "0"
    num_of_front_end_paths = "0"
    mobility_id_enabled    = false
    virtual_volumes        = true
    private_volumes        = false
    tdev                   = true
    vdev                   = false
    available_thin_volumes = false
    gatekeeper             = false
    data_volume            = false
    dld                    = false
    drv                    = false
    encapsulated           = false
    associated             = false
    reserved               = false
    pinned                 = false
    mapped                 = false
    bound_tdev             = true
    emulation              = "FBA"
    has_effective_wwn      = false
    effective_wwn          = "effective_wwn"
    split_name             = "split_name"
    type                   = "TDEV"
    unreducible_data_gb    = "0"
    cu_image_num           = "0"
    cu_image_ssid          = "cu_image_ssid"
    rdf_group_number       = "0"
    oracle_instance_name   = "oracle_instance_name"
    nguid                  = "nguid"
  }
}

output "volume_datasource_output" {
  value = data.powermax_volume.volume_datasource_test
}
