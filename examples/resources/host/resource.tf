resource "powermax_host" "host_1" {
  id        = "host_1"
  initiator = ["10000000c9959b8e"]
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