# Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
#
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://mozilla.org/MPL/2.0/
#
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# v1.0.0
## Release Summary
The release supports resources and data sources mentioned in the Features section for Dell PowerMax.
## Features

### Data Sources:
* `powermax_host` for reading hosts in PowerMax.
* `powermax_hostgroup` for reading host groups in PowerMax.
* `powermax_maskingview` for reading masking views in PowerMax.
* `powermax_portgroups` for reading port groups in PowerMax.
* `powermax_storagegroup` for reading storage groups in PowerMax.
* `powermax_volume` for reading volumes in PowerMax.
* `powermax_snapshotpolicy` for reading snapshot policies on PowerMax.
* `powermax_port` for reading ports on PowerMax.
* `powermax_snapshot` for reading snapshots on PowerMax.

### Resources
* `powermax_host` for managing hosts in PowerMax.
* `powermax_hostgroup` for managing host groups in PowerMax.
* `powermax_maskingview` for managing masking views in PowerMax.
* `powermax_portgroup` for managing port groups in PowerMax.
* `powermax_storagegroup` for managing storage groups in PowerMax.
* `powermax_volume` for managing volumes in PowerMax.
* `powermax_snapshotpolicy` for readomg snapshot policies on PowerMax.
* `powermax_snapshot` for reading snapshots on PowerMax.

### Others
N/A

## Enhancements
N/A

## Bug Fixes
N/A


# v1.0.0-beta
## Release Summary
The release supports HostGroup resource to manage HostGroup in PowerMax .
## Features

### Data Sources:
N/A

### Resources
* `powermax_host_group` for managing hostgroup in PowerMax.

### Others
N/A
## Enhancements
* `powermax_host` extracted consistent_lun from host_flags and now it can be configured as an independent parameter

## Bug Fixes
N/A


# v1.0.0-alpha

## Release Summary

The release supports a terraform plugin to manage StorageGroups,Volumes,Hosts,PortGroups and MaskingViews in PowerMax.

## Features

### Data Sources:
N/A

### Resources:
* `powermax_storage_group` for managing storage groups in PowerMax. Currently only Standalone storage groups are supported.
* `powermax_volume` for managing volumes in PowerMax. Currently volumes are not supported without storage group.
* `powermax_host` for managing hosts in PowerMax.
* `powermax_port_group` for managing port groups in PowerMax.
* `powermax_masking_view` for managing masking views in PowerMax.

### Others
N/A

## Enhancements
N/A

## Bug Fixes
N/A