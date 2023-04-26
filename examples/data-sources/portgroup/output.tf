output "fibreportgroups" {
  value = data.powermax_portgroups.fibreportgroups
} 

output "scsiportgroups" {
  value = data.powermax_portgroups.scsiportgroups
} 

output "allportgroups" {
  value = data.powermax_portgroups.allportgroups.port_groups
}