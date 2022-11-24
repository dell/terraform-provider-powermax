# v1.0.0

## Release Summary

The release supports a terraform client to query Unisphere and the resources to manage StorageGroups,Volumes,Hosts,PortGroups and MaskingViews in PowerMax.

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