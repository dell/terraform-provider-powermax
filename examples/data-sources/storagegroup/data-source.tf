data "powermax_storagegroup" "test" {
  filter {
    names = ["esa_sg572"]
  }
}

output "storagegroup_data" {
  value = data.powermax_storagegroup.test
}

data "powermax_storagegroup" "testall" {
}

output "storagegroup_data_all" {
  value = data.powermax_storagegroup.testall
}