resource "powermax_maskingview" "test" {
  name             = "terraform_testMV"
  storage_group_id = "Tao_k8s_env2_SG"
  host_id          = "Tao_k8s_env2_host"
  host_group_id    = ""
  port_group_id    = "Tao_k8s_env2_PG"
}
