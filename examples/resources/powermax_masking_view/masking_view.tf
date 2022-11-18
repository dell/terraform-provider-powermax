terraform {
  required_providers {
    powermax = {
       version = "0.0.1"
      source  = "dell/powermax"
    }
  }
}

provider "powermax" {
  username = var.username
  password = var.password
  endpoint = var.endpoint
  serial_number = var.serial_number
  insecure = false
}

resource "powermax_storage_group" "sg_for_masking_view" {
	name = "sg_maskingview"
	srpid = "SRP_1"
	service_level = "Diamond"
}

resource "powermax_volume" "volume_for_masking_view" {
	name = "vol_maskingview"
	size = 1
	cap_unit = "GB"
	sg_name = powermax_storage_group.sg_for_masking_view.name
}

resource "powermax_port_group" "pg_for_masking_view" {
	name = "pg_maskingview"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "OR-2C"
			port_id = "2"
		},
		{
			director_id = "OR-1C"
			port_id = "2"
		}
	]
}

resource "powermax_host" "host_for_masking_view" {
	name = "host_maskingview"
	initiators = var.initiator_ids
	host_flags = {
		volume_set_addressing = {
			override = true
			enabled = true
		}
		openvms = {
			override = true
			enabled = false
		}
		consistent_lun = false
	}

}

resource "powermax_masking_view" "create_masking_view" {
	name = "create_maskingview"
	storage_group_id = powermax_storage_group.sg_for_masking_view.id
	port_group_id = powermax_port_group.pg_for_masking_view.id
	host_id = powermax_host.host_for_masking_view.id
}

output "storage-group-id" {
  value =  powermax_storage_group.sg_for_masking_view.id
}

output "volume-id" {
  value =  powermax_volume.volume_for_masking_view.id
}

output "host-id" {
  value =  powermax_host.host_for_masking_view.id
}

output "port-group-id" {
  value =  powermax_port_group.pg_for_masking_view.id
}

output "masking-view-id" {
  value =  powermax_masking_view.create_masking_view.id
}
