terraform {
  required_providers {
    powermax = {
      version = "0.0.1"
      source  = "dell/powermax"
    }
  }
}

provider "powermax" {
  username      = var.username
  password      = var.password
  endpoint      = var.endpoint
  serial_number = var.serial_number
  pmax_version  = var.pmax_version
  insecure      = true
}

resource "powermax_portgroup" "portgroup_1" {
  name     = "tf_pg_test_1"
  protocol = "SCSI_FC"
  ports = [
    {
      director_id = "FA-2D"
      port_id     = "11"
    }
  ]
}