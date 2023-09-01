---
page_title: "Creating Multiple Volumes"
title: "Creating Multiple Volumes"
linkTitle: "Creating Multiple Volumes"
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

You can use the for_each meta-argument to create multiple volumes.

## Example

To create 3 different volumes using terraform_for_each using the following configuration:

```terraform
resource "powermax_volume" "test" {
  for_each = toset(["vol_1", "vol_2", "vol_3"])
  vol_name = each.key
  size = 2.45
  sg_name = "terraform_sg"
  cap_unit = "GB"
}
```

You can use the count meta-argument to create multiple volumes.

## Example

To create 3 different volumes using the following configuration:

```terraform
resource "powermax_volume" "test" {
  count = 3
  vol_name = "terraform_volume-${count.index}"
  size = 2.45
  sg_name = "terraform_sg"
  cap_unit = "GB"
}