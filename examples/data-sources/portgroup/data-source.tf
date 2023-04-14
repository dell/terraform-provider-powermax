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

# List all portgroups.
data "powermax_portgroups" "fiberportgroups" {
    type = "fiber"
}
data "powermax_portgroups" "scsiportgroups" {
    type = "iscsi"
}


output "fiberportgroups" {
  value = data.powermax_portgroups.fiberportgroups
} 

output "scsiportgroups" {
  value = data.powermax_portgroups.scsiportgroups
} 