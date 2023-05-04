# List a specific maskingView
data "powermax_maskingview" "maskingViewFilter" {
   filter {
    names = ["terraform_tao_testMV_rename", "Yulan_SG_MV"]
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
