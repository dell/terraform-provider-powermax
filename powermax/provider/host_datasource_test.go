package provider

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
func TestAccHostDatasourceFilterEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + HostDataSourceFilterEmpty,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powermax_host.hostEmptyFilter", "hosts.#"),
				),
			},
		},
	})
}
func TestAccHostDatasourceGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + HostDataSourceGetAll,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powermax_host.hostGetAll", "hosts.#"),
				),
			},
		},
	})
}

var HostDataSourceParamsAll = `
data "powermax_host" "HostDs" {
	filter {
		# Optional list of IDs to filter
		names = [
		  "tfacc_host_test1",
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
var HostDataSourceFilterEmpty = `
data "powermax_host" "hostEmptyFilter" {
	   filter {
    		names = []
	   }
}`

var HostDataSourceGetAll = `
data "powermax_host" "hostGetAll" {
}`
