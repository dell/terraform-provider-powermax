data "powermax_host" "HostDs" {
}

output "hostDsResult" {
   value = data.powermax_host.HostDs
}