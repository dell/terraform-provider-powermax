
# List fibre portgroups.
data "powermax_portgroups" "fibreportgroups" {
  # Optional filter to list specified Portgroups names and/or type
  filter {
    # type for which portgroups to be listed  - fibre or iscsi
    type = "fibre"
    # Optional list of IDs to filter
    names = [
      "tfacc_test1_fibre",
      #"test2_fibre",
    ]
  }
}

data "powermax_portgroups" "scsiportgroups" {
  filter {
    type = "iscsi"
    # Optional filter to list specified Portgroups Names
  }
}

# List all portgroups.
data "powermax_portgroups" "allportgroups" {
  #filter {
  # Optional list of IDs to filter
  #names = [
  #  "test1",
  #  "test2",
  #]
  #}
}

