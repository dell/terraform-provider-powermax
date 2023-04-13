resource "powermax_volume" "test" {
  sg_name  = "terraform_sg"
  vol_name = "terraform_volume"
  size     = 2.45
  cap_unit = "GB"
}