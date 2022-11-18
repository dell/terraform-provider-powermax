package powermax

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	TestAccCreatePGName        = "test_acc_create_pg"
	TestAccCreatePGNameUpdated = "test_acc_create_pg_updated"
)

func TestAccPortGroup_CreatePortGroup(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: CreatePortGroupParams,
				Check:  resource.ComposeTestCheckFunc(checkCreatePortGroup(t, testProvider, "test_acc_create_pg")),
			},
		},
	})
}

func TestAccPortGroup_CreatePortGroupWithInvalidPort(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      CreatePortGroupParamsWithInvalidPort,
				ExpectError: regexp.MustCompile(CreatePGDetailErrorMsg),
			},
		},
	})
}

func TestAccPortGroup_CreatePortGroupWithInvalidCombination(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      CreatePortGroupParamsWithInvalidCombination,
				ExpectError: regexp.MustCompile("The port number is out of range or does not exist"),
			},
		},
	})
}

func TestAccPortGroup_UpdatePortGroupSuccess(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: CreatePortGroupMultiplePorts,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_port_group.create_pg", "id", TestAccCreatePGName),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.#", "2"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.director_id", "OR-1C"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.port_id", "1"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.1.director_id", "OR-1C"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.1.port_id", "2")),
			},
			{
				Config: UpdatePortGroupParams,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_port_group.create_pg", "id", TestAccCreatePGNameUpdated),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.#", "2"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.director_id", "OR-1C"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.port_id", "1"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.1.director_id", "OR-2C"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.1.port_id", "2")),
			},
		},
	})
}

// Failure scenario : Failing to add ports- Non existent port
func TestAccPortGroup_UpdatePortGroupFailure1(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: CreatePortGroupMultiplePorts,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_port_group.create_pg", "id", TestAccCreatePGName),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.#", "2"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.director_id", "OR-1C"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.port_id", "1"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.1.director_id", "OR-1C"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.1.port_id", "2")),
			},
			{
				Config:      UpdatePortGroupParamsFailure1,
				ExpectError: regexp.MustCompile("Failed to update ports"),
			},
		},
	})
}

// Failure scenario : Failing to add ports- Duplicate port, API layer ignores duplicate ports and considers only 1 port. Hence it will cause plan inconsistency
func TestAccPortGroup_UpdatePortGroupFailure2(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: CreatePortGroupParams,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_port_group.create_pg", "id", TestAccCreatePGName),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.#", "1"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.director_id", "OR-1C"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.port_id", "2")),
			},
			{
				Config:      UpdatePortGroupParamsFailure2,
				ExpectError: regexp.MustCompile("Provider produced inconsistent result after apply"),
			},
		},
	})
}

func checkCreatePortGroup(t *testing.T, p tfsdk.Provider, portGroupID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providers := p.(*provider)
		_, err := providers.client.PmaxClient.GetPortGroupByID(context.Background(), serialno, portGroupID)
		if err != nil {
			return fmt.Errorf("failed to fetch portgroup")
		}
		if !providers.configured {
			return fmt.Errorf("provider not configured")
		}

		if providers.client.PmaxClient == nil {
			return fmt.Errorf("provider not configured")
		}
		return nil
	}
}

var CreatePortGroupParams = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	timeout = "20m"
	insecure = true
}

resource "powermax_port_group" "create_pg" {
	name = "test_acc_create_pg"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "` + DirectorID1 + `"
			port_id = "2"
		}
	]
}
`

var CreatePortGroupParamsWithInvalidPort = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_port_group" "create_pg" {
	name = "test_acc_create_pg"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "` + DirectorID1 + `"
			port_id = "InvalidPort"
		}
	]
}
`

var CreatePortGroupParamsWithInvalidCombination = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_port_group" "create_pg" {
	name = "test_acc_create_pg"
	protocol = "iSCSI"
	ports = [
		{
			director_id = "` + DirectorID1 + `"
			port_id = "2"
		}
	]
}
`
var CreatePortGroupMultiplePorts = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	timeout = "20m"
	insecure = true
}

resource "powermax_port_group" "create_pg" {
	name = "test_acc_create_pg"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "` + DirectorID1 + `"
			port_id = "1"
		},
		{
			director_id = "` + DirectorID1 + `"
			port_id = "2"
		}
	]
}
`

var UpdatePortGroupParams = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	timeout = "20m"
	insecure = true
}

resource "powermax_port_group" "create_pg" {
	name = "test_acc_create_pg_updated"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "` + DirectorID1 + `"
			port_id = "1"
		},
		{
			director_id = "` + DirectorID2 + `"
			port_id = "2"
		}
	]
}
`

var UpdatePortGroupParamsFailure1 = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	timeout = "20m"
	insecure = true
}

resource "powermax_port_group" "create_pg" {
	name = "test_acc_create_pg_updated"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "` + DirectorID1 + `"
			port_id = "2"
		},
		{
			director_id = "OR-20C"
			port_id = "2"
		}
	]
}
`
var UpdatePortGroupParamsFailure2 = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	timeout = "20m"
	insecure = true
}

resource "powermax_port_group" "create_pg" {
	name = "test_acc_create_pg_updated"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "` + DirectorID1 + `"
			port_id = "2"
		},
		{
			director_id = "` + DirectorID1 + `"
			port_id = "2"
		}
	]
}
`
