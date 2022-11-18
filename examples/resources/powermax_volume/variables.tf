variable "username" {
  type = string
  description = "stores the username of unisphere."
}

variable "password" {
  type = string
  description = "stores the password of unisphere."
}

variable "endpoint" {
    type = string
    description = "stores the endpoint of unisphere instance"
}

variable "serial_number" {
    type = string
    description = "stores the serial number of the storage array"
}

variable "initiator_ids" {
  type = list(string)
  description = "list of initiator ids."
}
