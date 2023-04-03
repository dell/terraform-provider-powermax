resource "powermax_storagegroup" "test" {
  storage_group_id = "terraform_sg"
  srp_id           = "SRP_1"
  slo              = "Gold"
  host_io_limit = {
    host_io_limit_io_sec = "1000"
    host_io_limit_mb_sec = "1000"
    dynamicDistribution  = "Never"
  }
  volume_size            = "100"
  capacity_unit          = "CYL"
  volume_identifier_name = "terraform_volume"
  num_of_vols            = 1
}