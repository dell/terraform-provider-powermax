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

# This terraform DataSource is used to query the existing volume from PowerMax array.
# The information fetched from this data source can be used for getting the details / for further processing in resource block.

# Returns all of the PowerMax volumes and their details
# NOTE: PowerMax can have many volumes, running this command unfiltered can take several minutes
data "powermax_volume" "volume_datasource_all" {
}

output "volume_datasource_output" {
  value = data.powermax_volume.volume_datasource_all
}

# Returns a subset of the PowerMax volumes based on the filtered properties in the block
# All filter values are optional
# If you use more then one filter at a time, it will only show the subset of volumes which both of those filters satisfies 
data "powermax_volume" "volume_datasource_test" {
  filter {

    # Optional Volume ids from a single Storage Group only
    storage_group_name = "terraform_vol_sg"

    # Optional Volume ids that contain the specified volume wwn
    wwn = "wwn_num"

    # Optiona lVolume ids that contain the specified volume encapsulated_wwn
    encapsulated_wwn = "encapsulated_wwn_num"

    # Optional Volume ids that contain greater than(">1"), Less than("<1") or equal to the specified symmlun
    symmlun = "0"

    # Optional Volume ids that contain the specified volume status
    status = "Ready"

    # Optional Volume ids that contain the specified volume physical_name
    physical_name = "physical_name"

    # Optional Volume ids that contain the specified volume volume_identifier
    volume_identifier = "test_acc_create_volume"

    # Optional Volume ids that contain greater than(">1"), Less than("<1") or equal to the specified cap_tb
    cap_tb = "0"

    # Optional Volume ids that contain greater than(">1"), Less than("<1") or equal to the specified cap_gb
    cap_gb = "0"

    # Optional Volume ids that contain greater than(">1"), Less than("<1") or equal to the specified cap_mb
    cap_mb = "0"

    # Optional Volume ids that contain greater than(">1"), Less than("<1") or equal to the specified cap_cyl
    cap_cyl = "0"

    # Optional Volume ids that contain greater than(">1"), Less than("<1") or equal to the specified allocated_percent
    allocated_percent = "0"

    # Optional Volume ids that contain greater than(">1"), Less than("<1") or equal to the specified num_of_storage_groups
    num_of_storage_groups = "1"

    # Optional Volume ids that contain greater than(">1"), Less than("<1") or equal to the specified num_of_masking_views
    num_of_masking_views = "0"

    # Optional Volume ids that contain greater than(">1"), Less than("<1") or equal to the specified num_of_front_end_paths
    num_of_front_end_paths = "0"

    # Optional Volume ids that are mobility ID enabled (true/false)
    mobility_id_enabled = false

    # Optional Volume ids that are virtual_volumes (true/false)
    virtual_volumes = true

    # Optional Volume ids that are private_volumes (true/false)
    private_volumes = false

    # Optional Volume ids that are tdev (true/false)
    tdev = true

    # Optional Volume ids that are vdev (true/false)
    vdev = false

    # Optional Volume ids that are available_thin_volumes (true/false)
    available_thin_volumes = false

    # Optional Volume ids that are gatekeeper (true/false)
    gatekeeper = false

    # Optional Volume ids that are data_volume (true/false)
    data_volume = false

    # Optional Volume ids that are dld (true/false)
    dld = false

    # Optional Volume ids that are drv (true/false)
    drv = false

    # Optional Volume ids that are encapsulated (true/false)
    encapsulated = false

    # Optional Volume ids that are associated (true/false)
    associated = false

    # Optional Volume ids that are reserved (true/false)
    reserved = false

    # Optional	Volume ids that are pinned (true/false)
    pinned = false

    # Optional Volume ids that are mapped (true/false)
    mapped = false

    # Optional Volume ids that are bound_tdev (true/false)
    bound_tdev = true

    # Optional Volume ids that are of the specified emulation
    emulation = "FBA"

    # Optional Volume ids that are of the specified emulation.
    has_effective_wwn = false

    # Optional Volume ids that contain the specified effective_wwn
    effective_wwn = "effective_wwn"

    # Optional Volume ids that are mapped to CU images associated to the specified FICON split
    split_name = "split_name"

    # Optional Volume ids that contain the specified volume type
    type = "TDEV"

    # Optional Volume ids that contain greater than("unreducible_data_gb=>1"),Less than("unreducible_data_gb=<1") or equal to the unreducible_data_gb
    unreducible_data_gb = "0"

    # Optional Volume ids that are mapped to a CU image with the specified CU image number
    cu_image_num = "0"

    # Optional Volume ids that are mapped to a CU image with the specified CU SSID
    cu_image_ssid = "cu_image_ssid"

    # Optional Volume ids that are part of the specified rdf group
    rdf_group_number = "0"

    # Optional Volume ids that contain the specified Oracle Instance Name
    oracle_instance_name = "oracle_instance_name"

    # Optional Volumes Ids that correspond to Namespace Globally Unique Identifier that uses the EUI64 16-byte designator format. Used in conjunction with NVMe volumes
    nguid = "nguid"
  }
}

output "volume_datasource_output" {
  value = data.powermax_volume.volume_datasource_test
}

# After the successful execution of above said block, We can see the output value by executing 'terraform output' command.
# Also, we can use the fetched information by the variable data.powermax_volume.example