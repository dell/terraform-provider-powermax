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

resource "powermax_storage_group" "storage_group_2" {
  name          = "sg_2"
  srpid         = "SRP_1"
  service_level = "Diamond"
  volume_ids    = ["0033D"]
  host_io_limits = {
    host_io_limit_mb_sec = "1"
    host_io_limit_io_sec = "100"
    dynamicdistribution  = "Always"
  }
  snapshot_policies = [
    {
      is_active   = true
      policy_name = "terraform_policy"
    }
  ]
}

resource "powermax_storage_group" "sg_import" {
}
