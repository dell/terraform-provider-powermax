resource "powermax_host" "host_1" {
  name      = "host_1"
  initiator = ["10000000c9fc4b7e"]
  host_flags = {
    disable_q_reset_on_ua = {
      override = true
      enabled  = true
    }

  }
}