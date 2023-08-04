/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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

# resource "powermax_volume" "test_vol" {
#   sg_name  = "terraform_sg"
#   vol_name = "test_acc_create_volume_1"
#   size     = 2.45
#   cap_unit = "GB"
# }

#data "powermax_volume" "volume_datasource_test" {
#  filter {
#    storage_group_name = "terraform_sg"
#  }
#}

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
#   name             = "terraform_sgnewup"
#   srp_id           = "SRP_1"
#   slo              = "Gold"
#   host_io_limit = {
#     host_io_limit_io_sec = "1000"
#     host_io_limit_mb_sec = "1000"
#     dynamic_distribution  = "Never"
#   }
#   volume_ids = ["0008F"]
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
# }

# resource "powermax_maskingview" "test_masking_view" {
#   name = "terraform_testMV"
#   storage_group_id = "Tao_k8s_env2_SG"
#   host_id = "Tao_k8s_env2_host"
#   host_group_id = ""
#   port_group_id = "Tao_k8s_env2_PG"
# }

# data "powermax_host" "HostDs" {
# }

# output "hostDsResult" {
#    value = data.powermax_host.HostDs
# }

# data "powermax_maskingview" "maskingViewData" {
#    filter {
#     names = ["terraform_tao_testMV_rename", "Yulan_SG_MV"]
#   }
# }

# output "maskingViewDataResult" {
#   value = data.powermax_maskingview.maskingViewData.masking_views
# }
