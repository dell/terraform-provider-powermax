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
  insecure      = false
}

resource "powermax_volume" "volume_2" {
  name               = "volume_2"
  size               = 2.5
  cap_unit           = "GB"
  sg_name            = "storage-group-1"
  enable_mobility_id = false
}


resource "powermax_volume" "volume_import" {
}
