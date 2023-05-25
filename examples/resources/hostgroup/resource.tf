// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
resource "powermax_hostgroup" "test_host_group" {
  # Optional
  host_flags = {
    avoid_reset_broadcast = {
      enabled  = true
      override = true
    }
  }
  host_ids = ["testHost"]
  name     = "host_group"
}