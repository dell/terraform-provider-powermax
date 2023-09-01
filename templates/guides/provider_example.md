---
page_title: "How to use our Provider"
title: "How to use our Provider"
linkTitle: "How to use our Provider"
---

<!--
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
-->

This provider gives Terraform the ablility to work with a PowerMax Storage System.
It can be used to manage many aspects of the PowerMax including hosts, host groups, masking views, port groups, snapshots, storage groups and volumes.
For more information about each particular resource or datasource view those specific examples and documentation.

Below is an example of the PowerMax Provider and how to use it.
## Example 

The following example demonstrates basic usage of the provider to provision a volume to a storage group. 

```
terraform {
  required_providers {
    powermax = {
      source = "dell/powermax"
    }
  }
}

provider "powermax" {
  username      = var.username
  password      = var.password
  endpoint      = var.endpoint
  serial_number = var.serial_number
  pmax_version  = var.pmax_version
  insecure      = true
}

resource "powermax_volume" "test" {

  # Required name of the volume to be created
  vol_name = "terraform_volume"

  # Required size of the volume
  size = 2.45

  # Required name of the storage group which the volume will be created with
  sg_name = "terraform_sg"

  # Optional Default Unit is GB
  # Possible units are MB, GB, TB, and CYL
  cap_unit = "GB"

  # Optional enable the mobility id 
  mobility_id_enabled = false

}
```

*Refer to the specific provider and resource documetation for more information*

## Argument Reference

The following arguments are used to configure the provider

### Required

- `endpoint` (String) IP or FQDN of the PowerMax host
- `password` (String, Sensitive) The password of the PowerMax host.
- `pmax_version` (String) The version of the PowerMax host.
- `serial_number` (String) The serial_number of the PowerMax host.
- `username` (String) The username of the PowerMax host.

### Optional

- `insecure` (Boolean) Boolean variable to specify whether to validate SSL certificate or not.

## Bug Reports and Contributing

If you need to sumbit a bug please do so on our github issues page https://github.com/dell/terraform-provider-powermax/issues
If you would like to contribute please follow the README here: https://github.com/dell/terraform-provider-powermax