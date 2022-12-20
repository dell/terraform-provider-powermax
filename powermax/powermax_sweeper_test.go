package powermax

import (
	"fmt"
	"terraform-provider-powermax/client"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func getSweeperClient(region string) (*client.Client, error) {
	powermaxClient, err := client.NewClient(endpoint, username, password, serialno, "100", true)
	if err != nil {
		return nil, fmt.Errorf("Unable to create sweeper client %s", err.Error())
	}
	return powermaxClient, nil
}
