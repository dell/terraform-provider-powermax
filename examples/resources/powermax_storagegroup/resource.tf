resource "powermax_storagegroup" "test" {
  name   = "terraform_sg"
  srp_id = "SRP_1"
  slo    = "Gold"
  host_io_limit = {
    host_io_limit_io_sec = "1000"
    host_io_limit_mb_sec = "1000"
    dynamic_distribution = "Never"
  }
  volume_ids = ["0008F"]
}