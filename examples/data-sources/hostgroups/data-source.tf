// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.

# List all hostgroups
data "powermax_hostgroup" "all" {}

output "all" {
  value = data.powermax_hostgroup.all
}

# List a specific hostgroup
data "powermax_hostgroup" "groups" {
  filter {
    names = ["host_group_example_1", "host_group_example_2"]
  }
}

output "groups" {
  value = data.powermax_hostgroup.groups
}