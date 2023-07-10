resource "powermax_snapshot" "snapshot_1" {
  # Required, what storage group to take the snapshot
  storage_group {
    name = "test_sg"
  }
  snapshot_actions {
    # Required, name of new snapshot
    name = "new_test_snapshot"
    # Optional
    # secure = {
    #   enable = true
    #   time_in_hours = true
    #   secure = 3
    # }
    # Optional
    # time_to_live = {
    #   enable = true
    #   time_in_hours = false
    #   time_to_live = 1
    # }
    # Optional
    # link = {
    #   enable = false
    #   target_storage_group = "test_target_sg"
    #   no_compression = true
    #   remote = false
    #   copy = false
    # }
    # Optional
    # set_mode = {
    #   enable = true
    #   target_storage_group = "test_target_sg"
    #   copy = false
    # }
    # Optional
    # restore = {
    #   enable = true
    #   remote = false
    # }
  }
}