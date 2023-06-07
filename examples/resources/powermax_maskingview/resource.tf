resource "powermax_maskingview" "test" {
  name             = "terraform_mv"
  storage_group_id = "terraform_sg"
  host_id          = "terraform_host"
  host_group_id    = ""
  port_group_id    = "terraform_pg"
}
