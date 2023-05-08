resource "powermax_host" "host_1" {
  name       = "host_1"
  initiators = ["0000000000000001"]
  host_flags = {
    volume_set_addressing = {
      override = true
      enabled  = true
    }
    openvms = {
      override = true
      enabled  = false
    }
    consistent_lun = false
  }
}
