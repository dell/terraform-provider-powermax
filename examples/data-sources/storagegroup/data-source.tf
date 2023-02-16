data "powermax_storagegroup" "test" {
  storage_groups: [{
    "storage_group_id": "esa_sg572",
  }]
}

output "storagegroup_data" {
  value = data.powermax_storagegroup.test.storage_groups
}