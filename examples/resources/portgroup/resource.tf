terraform {
  required_providers {
    powermax = {
      version = "0.0.1"
      source  = "dell/powermax"
    }
  }
}
variable "username" {
  type = string
  default="smc"
}
variable "password" {
  type = string
  default="smc"
}
variable "endpoint" {
  type = string
  default= "https://10.225.104.33:8443"
}
variable "serial_number" {
  type = string
  default = "000197902572"
}
variable "pmax_version" {
  type = string
  default = "100"
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
	name = "tf_pg_1_1"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "FA-2D"
			port_id = "11"
		},
		{
			director_id = "FA-2D"
			port_id = "10"
		}
	]
}