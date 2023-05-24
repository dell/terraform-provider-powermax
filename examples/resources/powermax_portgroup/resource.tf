resource "powermax_portgroup" "portgroup_1" {
  name     = "tfacc_pg_test_1"
  protocol = "SCSI_FC"
  ports = [
    {
      director_id = "OR-1C"
      port_id     = "0"
    }
  ]
}