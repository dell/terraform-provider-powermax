package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test to Fetch Host details
func TestAccHostDatasource(t *testing.T) {
	var hostName = "data.powermax_host.HostDs"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + HostDataSourceParamsAll,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(hostName, "host_ids.#", "19"),
				),
			},
		},
	})
}

var HostDataSourceParamsAll = `
data "powermax_host" "HostDs" {
	
}
output "hostDsResult" {
	value = data.powermax_host.HostDs
 }
`
