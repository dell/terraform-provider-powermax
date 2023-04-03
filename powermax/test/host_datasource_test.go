package test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test to Fetch Host details.
func TestAccHostDatasource(t *testing.T) {
	var hostName = "data.powermax_host.HostDs"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + HostDataSourceParamsAll,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(hostName, "hosts.#", "1"),
				),
			},
		},
	})
}

func TestAccHostDatasourceFilteredError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + HostDataSourceFilterError,
				ExpectError: regexp.MustCompile(`.*Error reading host with id*.`),
			},
		},
	})
}

var HostDataSourceParamsAll = `
data "powermax_host" "HostDs" {
	filter {
		# Optional list of IDs to filter
		names = [
		  "hostExample",
		]
	}
	
}
output "hostDsResult" {
	value = data.powermax_host.HostDs
 }
`

var HostDataSourceFilterError = `
# List a specific host
data "powermax_host" "hosts" {
  filter {
    names = ["non-existent-host"]
  }
}

output "hosts" {
  value = data.powermax_host.hosts
}
`
