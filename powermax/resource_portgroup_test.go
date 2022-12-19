package powermax

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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

func TestAccPortGroup_CreatePortGroupUpdateExistingName(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	assertTFImportState := func(s []*terraform.InstanceState) error {
		assert.Equal(t, "test_acc_create_pg", s[0].Attributes["id"])
		assert.Equal(t, "test_acc_create_pg", s[0].Attributes["name"])
		assert.Equal(t, "SCSI_FC", s[0].Attributes["protocol"])
		resource.TestCheckResourceAttr("powermax_port_group.import_pg", "ports.#", "1")
		resource.TestCheckResourceAttr("powermax_port_group.import_pg", "ports.0.port_id", "1")
		resource.TestCheckResourceAttr("powermax_port_group.import_pg", "ports.0.director_id", "1")
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: CreatePortGroupParams,
				Check:  resource.ComposeTestCheckFunc(checkCreatePortGroup(t, testProvider, "test_acc_create_pg")),
			},
			{
				Config:      UpdatePortGroupParamsExistingName,
				ExpectError: regexp.MustCompile(UpdatePGDetailsErrMsg),
			},
			{
				Config:           ImportPortGroup,
				ResourceName:     "powermax_port_group.import_pg",
				ImportState:      true,
				ImportStateCheck: assertTFImportState,
				ExpectError:      nil,
				ImportStateId:    TestAccCreatePGName,
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
	assertTFImportState := func(s []*terraform.InstanceState) error {
		assert.Equal(t, TestAccCreatePGNameUpdated, s[0].Attributes["id"])
		assert.Equal(t, TestAccCreatePGNameUpdated, s[0].Attributes["name"])
		assert.Equal(t, "SCSI_FC", s[0].Attributes["protocol"])
		resource.TestCheckResourceAttr("powermax_port_group.import_pg", "ports.#", "2")
		resource.TestCheckResourceAttr("powermax_port_group.import_pg", "ports.0.port_id", "0")
		resource.TestCheckResourceAttr("powermax_port_group.import_pg", "ports.0.director_id", DirectorID1)
		resource.TestCheckResourceAttr("powermax_port_group.import_pg", "ports.1.port_id", "2")
		resource.TestCheckResourceAttr("powermax_port_group.import_pg", "ports.1.director_id", DirectorID2)
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: CreatePortGroupMultiplePorts,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_port_group.create_pg", "id", TestAccCreatePGName),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.#", "2"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.director_id", DirectorID1),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.port_id", "0"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.1.director_id", DirectorID1),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.1.port_id", "2")),
			},
			{
				Config: UpdatePortGroupParams,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_port_group.create_pg", "id", TestAccCreatePGNameUpdated),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.#", "2"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.director_id", DirectorID1),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.port_id", "0"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.1.director_id", DirectorID2),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.1.port_id", "2")),
			},
			{
				Config:           ImportPortGroup,
				ResourceName:     "powermax_port_group.import_pg",
				ImportState:      true,
				ImportStateCheck: assertTFImportState,
				ExpectError:      nil,
				ImportStateId:    TestAccCreatePGNameUpdated,
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
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.director_id", DirectorID1),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.port_id", "0"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.1.director_id", DirectorID1),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.1.port_id", "2")),
			},
			{
				Config:      UpdatePortGroupParamsFailure1,
				ExpectError: regexp.MustCompile("Failed to update ports"),
			},
		},
	})
}

// Update with duplicate ports
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
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.director_id", DirectorID1),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.port_id", "2")),
			},
			{
				Config:      UpdateDuplicatePortGroup,
				ExpectError: nil,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_port_group.create_pg", "id", TestAccCreatePGNameUpdated),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.#", "1"),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.director_id", DirectorID1),
					resource.TestCheckResourceAttr("powermax_port_group.create_pg", "ports.0.port_id", "2")),
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
	name = "` + TestAccCreatePGName + `"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "` + DirectorID1 + `"
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

resource "powermax_port_group" "create_pg" {
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

var ImportPortGroup = `
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

resource "powermax_port_group" "create_pg" {
	name = "` + TestAccCreatePGNameUpdated + `"
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
	name = "` + TestAccCreatePGNameUpdated + `"
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
var UpdateDuplicatePortGroup = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	timeout = "20m"
	insecure = true
}

resource "powermax_port_group" "create_pg" {
	name = "` + TestAccCreatePGNameUpdated + `"
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
