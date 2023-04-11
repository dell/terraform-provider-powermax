terraform {
  required_providers {
    powermax = {
      source = "registry.terraform.io/dell/powermax"
    }
  }
}

variable "username" {
  type = string
}

variable "password" {
  type = string
}

variable "endpoint" {
  type = string
}

variable "serial_number" {
  type = string
}

variable "pmax_version" {
  type = string
}

provider "powermax" {
  username      = var.username
  password      = var.password
  endpoint      = var.endpoint
  serial_number = var.serial_number
  pmax_version  = var.pmax_version
  insecure      = true
}

# resource "powermax_host" "host_02_test" {
#   name = "host_02_test"
# 	initiator = ["10000000c9959b8e"]
# 	host_flags = {
# 		volume_set_addressing = {
# 			override = true
# 			enabled = true
# 		}  
# 		openvms = {
# 			override = true
# 			enabled = false
# 		}
#     avoid_reset_broadcast = {
#       override = true
#       enabled = false
#     }
#   }
# }

# resource "powermax_storagegroup" "sg56new" {
#   storage_group_id = "terraform_sgnewup"
#   srp_id           = "SRP_1"
#   slo              = "Gold"
#   host_io_limit = {
#     host_io_limit_io_sec = "1000"
#     host_io_limit_mb_sec = "1000"
#     dynamicDistribution  = "Never"
#   }
#   volume_size            = "100"
#   capacity_unit          = "CYL"
#   volume_identifier_name = "terraform_volume"
#   num_of_vols            = 1
# }

# resource "powermax_hostgroup" "test_host_group" {
#   host_flags = {
#         avoid_reset_broadcast = {
#             enabled  = true
#             override = true
#         }
#   }
#   host_ids = ["testHost"]
#   name     = "host_group"
#   # volume_size            = "100"
#   # capacity_unit          = "CYL"
#   # volume_identifier_name = "terraform_volume"
#   # num_of_vols            = 1
# }
