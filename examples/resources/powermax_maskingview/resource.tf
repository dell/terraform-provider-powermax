resource "powermax_maskingview" "test" {
  name             = "terraform_mv"
  storage_group_id = "TestnewSG"
  host_id          = "Host124"
  host_group_id    = ""
  port_group_id    = "TestnewSG_PG"
}
