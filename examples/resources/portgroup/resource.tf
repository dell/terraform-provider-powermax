resource "powermax_portgroup" "portgroup_1" {
  name     = "tf_pg_test_1"
  protocol = "SCSI_FC"
  ports = [
    {
      director_id = "FA-2D"
      port_id     = "11"
    }
  ]
}