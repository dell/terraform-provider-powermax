# List a specific ports
data "powermax_port" "portFilter" {
  filter {
    # Should be in the format ["directorId:portId"]
    port_ids = ["OR-1C:2"]
  }
}

output "portFilter" {
  value = data.powermax_port.portFilter
}

# List all ports
data "powermax_port" "all" {}

output "all" {
  value = data.powermax_port.all
}
