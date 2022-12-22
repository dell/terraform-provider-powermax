package powermax

import (
	"context"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

// It is mandatory to create `test` resources with a prefix - 'test_acc_'
const (
	TestAccCreatePGName        = "test_acc_create_pg"
	TestAccCreatePGNameUpdated = "test_acc_create_pg_updated"
)

func init() {
	resource.AddTestSweepers("powermax_port_group", &resource.Sweeper{
		Name:         "powermax_port_group",
		Dependencies: []string{"powermax_masking_view"},
		F: func(region string) error {
			powermaxClient, err := getSweeperClient(region)
			if err != nil {
				log.Println("Error getting sweeper client: " + err.Error())
				return nil
			}

			ctx := context.Background()

			portGroups, err := powermaxClient.PmaxClient.GetPortGroupList(ctx, serialno, "")
			if err != nil {
				log.Println("Error getting portgroup list: " + err.Error())
				return nil
			}

			for _, portGroup := range portGroups.PortGroupIDs {
				if strings.Contains(portGroup, SweepTestsTemplateIdentifier) {
					err := powermaxClient.PmaxClient.DeletePortGroup(ctx, serialno, portGroup)
					if err != nil {
						log.Println("Error deleting portgroup: " + portGroup + "with error: " + err.Error())
					}
				}
			}
			return nil
		},
	})
}

func TestAccPortGroup_CRUDPortGroup(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	assertTFImportState := func(s []*terraform.InstanceState) error {
		assert.Equal(t, TestAccCreatePGNameUpdated, s[0].Attributes["id"])
		assert.Equal(t, TestAccCreatePGNameUpdated, s[0].Attributes["name"])
		assert.Equal(t, "SCSI_FC", s[0].Attributes["protocol"])
		resource.TestCheckResourceAttr("powermax_port_group.import_pg", "ports.#", "2")
		resource.TestCheckResourceAttr("powermax_port_group.import_pg", "ports.0.port_id", "2")
		resource.TestCheckResourceAttr("powermax_port_group.import_pg", "ports.0.director_id", DirectorID1)
		resource.TestCheckResourceAttr("powermax_port_group.import_pg", "ports.1.port_id", "0")
		resource.TestCheckResourceAttr("powermax_port_group.import_pg", "ports.1.director_id", DirectorID2)
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: CreatePortGroupParams,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "id", TestAccCreatePGName),
					resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "ports.#", "2"),
					resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "ports.0.director_id", DirectorID1),
					resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "ports.0.port_id", "0"),
					resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "ports.1.director_id", DirectorID2),
					resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "ports.1.port_id", "2")),
			},
			{
				Config:      UpdatePortGroupParamsExistingName,
				ExpectError: regexp.MustCompile(UpdatePGDetailsErrMsg),
			},
			{
				// Failure scenario : Failing to add ports- Non existent port
				Config:      UpdatePortGroupParamsFailure1,
				ExpectError: regexp.MustCompile("Failed to update ports"),
			},
			{
				// Update with duplicate ports
				Config:      UpdateDuplicatePortGroup,
				ExpectError: nil,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "id", TestAccCreatePGName),
					resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "ports.#", "2"),
					resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "ports.0.director_id", DirectorID1),
					resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "ports.0.port_id", "0"),
					resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "ports.1.director_id", DirectorID2),
					resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "ports.1.port_id", "2")),
			},
			{
				Config: UpdatePortGroupParams,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "id", TestAccCreatePGNameUpdated),
					resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "ports.#", "2"),
					resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "ports.0.director_id", DirectorID1),
					resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "ports.0.port_id", "2"),
					resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "ports.1.director_id", DirectorID2),
					resource.TestCheckResourceAttr("powermax_port_group.crud_pg", "ports.1.port_id", "0")),
			},
			{
				Config:           ImportPortGroupSuccess,
				ResourceName:     "powermax_port_group.import_pg",
				ImportState:      true,
				ImportStateCheck: assertTFImportState,
				ExpectError:      nil,
				ImportStateId:    TestAccCreatePGNameUpdated,
			},
			{
				Config:        ImportPortGroupFailure,
				ResourceName:  "powermax_port_group.import_pg",
				ImportState:   true,
				ExpectError:   regexp.MustCompile("Error importing portgroup"),
				ImportStateId: "InvalidPortGroup",
			},
		},
	})
}

func TestAccPortGroup_CreatePortGroupFailure(t *testing.T) {
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
			{
				Config:      CreatePortGroupParamsWithInvalidCombination,
				ExpectError: regexp.MustCompile("The port number is out of range or does not exist"),
			},
		},
	})
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

resource "powermax_port_group" "crud_pg" {
	name = "` + TestAccCreatePGName + `"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "` + DirectorID1 + `"
			port_id = "0"
		},
		{
			director_id = "` + DirectorID2 + `"
			port_id = "2"
		}
	]
}
`
var UpdatePortGroupParamsExistingName = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	timeout = "20m"
	insecure = true
}

resource "powermax_port_group" "crud_pg" {
	name = "` + PortGroupID1 + `"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "` + DirectorID1 + `"
			port_id = "2"
		}
	]
}
`

var ImportPortGroupSuccess = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	timeout = "20m"
	insecure = true
}

resource "powermax_port_group" "import_pg" {
}
`

var ImportPortGroupFailure = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	timeout = "20m"
	insecure = true
}

resource "powermax_port_group" "import_pg" {
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
	name = "` + TestAccCreatePGName + `"
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
	name = "` + TestAccCreatePGName + `"
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
	name = "` + TestAccCreatePGName + `"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "` + DirectorID1 + `"
			port_id = "0"
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

resource "powermax_port_group" "crud_pg" {
	name = "` + TestAccCreatePGNameUpdated + `"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "` + DirectorID2 + `"
			port_id = "0"
		},
		{
			director_id = "` + DirectorID1 + `"
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

resource "powermax_port_group" "crud_pg" {
	name = "` + TestAccCreatePGName + `"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "` + DirectorID1 + `"
			port_id = "0"
		},
		{
			director_id = "` + DirectorID2 + `"
			port_id = "2"
		},
		{
			director_id = "INVALID_DIRECTOR_ID"
			port_id = "2"
		}
	]
}
`
var UpdateDuplicatePortGroup = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	timeout = "20m"
	insecure = true
}

resource "powermax_port_group" "crud_pg" {
	name = "` + TestAccCreatePGName + `"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "` + DirectorID1 + `"
			port_id = "0"
		},
		{
			director_id = "` + DirectorID2 + `"
			port_id = "2"
		},
		{
			director_id = "` + DirectorID2 + `"
			port_id = "2"
		}
	]
}
`
