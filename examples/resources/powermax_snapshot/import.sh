# The command is
# terraform import powermax_snapshot.snapshot_test storage_group.snapshot_name
# Example: must be storage_group.snapshot_name
terraform import powermax_snapshot.snapshot_test storage_group.snapshot_name
# after running this command, populate the name field in the config file to start managing this resource