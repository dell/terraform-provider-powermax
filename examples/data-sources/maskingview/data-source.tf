data "powermax_maskingview" "id" {
  id = "terraform_tao_testMV_rename"
}

output "id" {
    value = data.powermax_maskingview.id.masking_views
}

data "powermax_maskingview" "idList" {
  masking_view_ids = [ "terraform_tao_testMV_rename", "Yulan_SG_MV" ]
}

output "idList" {
    value = data.powermax_maskingview.idList.masking_views
}

data "powermax_maskingview" "all" {}

output "all" {
    value = data.powermax_maskingview.all.maskingviews
}