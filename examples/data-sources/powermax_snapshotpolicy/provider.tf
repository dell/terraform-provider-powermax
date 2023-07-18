terraform {
  required_providers {
    powermax = {
      source = "dell/powermax"
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