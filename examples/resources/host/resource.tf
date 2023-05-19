resource "powermax_host" "host_1" {
  name      = "host_1"
  initiator = ["10000000c9fc4b7e"]
  host_flags = {
    volume_set_addressing = {
      override = true
      enabled  = true
    }
    openvms = {
      override = true
      enabled  = false
    }
  }
}