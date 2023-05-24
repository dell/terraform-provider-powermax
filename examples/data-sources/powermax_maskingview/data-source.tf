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
