resource "powermax_volume" "volume_1" {
	name = "volume-1"
	size = 1
	cap_unit = "GB"
	sg_name = "StorageGroup-1"
}
