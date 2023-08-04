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
# List a specific maskingView
data "powermax_maskingview" "maskingViewFilter" {
  filter {
    names = ["terraform_mv_1", "terraform_mv_2"]
  }
}

output "maskingViewFilterResult" {
  value = data.powermax_maskingview.maskingViewFilter.masking_views
}

# List all maskingviews
data "powermax_maskingview" "allMaskingViews" {}

output "allMaskingViewsResult" {
  value = data.powermax_maskingview.allMaskingViews.masking_views
}
