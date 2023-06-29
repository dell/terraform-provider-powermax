data "powermax_snapshot" "test" {
  # The storage group to which you want to see all the snapshots
  # Required
  storage_group {
    name = "example_storage_group"
  }
}

output "powermax_snapshot" {
  value = data.powermax_snapshot.test
}
