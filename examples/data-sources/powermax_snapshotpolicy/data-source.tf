data "powermax_snapshotpolicy" "SnapshotPolicyAll" {
}

data "powermax_snapshotpolicy" "SnapshotPolicyFiltered" {
  filter {
    # Optional list of IDs to filter
    names = [
      "tfacc_snapshotPolicy1",
    ]
  }
}

output "SnapshotPolicyAll" {
  value = data.powermax_snapshotpolicy.SnapshotPolicyAll
}

output "SnapshotPolicyFiltered" {
  value = data.powermax_snapshotpolicy.SnapshotPolicyFiltered
}