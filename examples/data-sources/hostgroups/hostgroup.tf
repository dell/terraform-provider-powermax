// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.

# List all hostgroups
data "powermax_hostgroup" "groups" {}

output "groups" {
  value = data.powermax_hostgroup.groups
}