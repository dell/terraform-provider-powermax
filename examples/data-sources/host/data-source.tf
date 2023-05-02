data "powermax_host" "HostDsAll" {
}

data "powermax_host" "HostDsFiltered" {
  filter {
    # Optional list of IDs to filter
    ids = [
      "Host124",
      "Host173",
    ]
  }
}

output "hostDsResultAll" {
  value = data.powermax_host.HostDsAll
}

output "hostDsResult" {
  value = data.powermax_host.HostDsFiltered
}