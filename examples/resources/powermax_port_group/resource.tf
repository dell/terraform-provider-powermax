resource "powermax_port_group" "portgroup_1" {
  name     = "pg_1"
  protocol = "SCSI_FC"
  ports = [
    {
      director_id = "dir-1"
      port_id     = "2"
    },
    {
      director_id = "dir-2"
      port_id     = "2"
    }
  ]
}
