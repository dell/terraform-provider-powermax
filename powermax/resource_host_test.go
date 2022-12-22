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
	TestAccHostName1        = "test_acc_chost1"
	TestAccHostName1Updated = "test_acc_chost1_updated"
	TestAccHostName4        = "test_acc_chost3"
	TestAccHostName5        = "test_acc_chost4"
	InvalidInitiatorID      = "0000000000000000"
	ImportHostResourceName1 = "powermax_host.import_host_success"
	ImportHostResourceName2 = "powermax_host.import_host_failure"
)

func init() {
	resource.AddTestSweepers("powermax_host", &resource.Sweeper{
		Name:         "powermax_host",
		Dependencies: []string{"powermax_masking_view"},
		F: func(region string) error {
			powermaxClient, err := getSweeperClient(region)
			if err != nil {
				log.Println("Error getting sweeper client")
				return nil
			}

			ctx := context.Background()

			hosts, err := powermaxClient.PmaxClient.GetHostList(ctx, serialno)
			if err != nil {
				log.Println("Error getting host list")
				return nil
			}

			for _, host := range hosts.HostIDs {
				if strings.Contains(host, SweepTestsTemplateIdentifier) {
					err := powermaxClient.PmaxClient.DeleteHost(ctx, serialno, host)
					if err != nil {
						log.Println("Error deleting host")
					}
				}
			}
			return nil
		},
	})
}

func TestAccHost_CRUDHost(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	assertTFImportState := func(s []*terraform.InstanceState) error {
		assert.Equal(t, TestAccHostName1Updated, s[0].Attributes["name"])
		assert.Equal(t, "1", s[0].Attributes["initiators.#"])
		assert.Equal(t, InitiatorID2, s[0].Attributes["initiators.0"])
		assert.Equal(t, "true", s[0].Attributes["host_flags.volume_set_addressing.enabled"])
		assert.Equal(t, "true", s[0].Attributes["host_flags.volume_set_addressing.override"])
		assert.Equal(t, "true", s[0].Attributes["host_flags.spc2_protocol_version.enabled"])
		assert.Equal(t, "true", s[0].Attributes["host_flags.spc2_protocol_version.override"])
		assert.Equal(t, "1", s[0].Attributes["numofinitiators"])
		assert.Equal(t, "0", s[0].Attributes["numofmaskingviews"])
		assert.Equal(t, 1, len(s))
		return nil
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactory,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: CreateHostParams,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_host.host_crud_test", "name", TestAccHostName1),
					resource.TestCheckResourceAttr("powermax_host.host_crud_test", "initiators.#", "1"),
					resource.TestCheckResourceAttr("powermax_host.host_crud_test", "initiators.0", InitiatorID1),
					resource.TestCheckResourceAttr("powermax_host.host_crud_test", "host_flags.spc2_protocol_version.enabled", "false"),
					resource.TestCheckResourceAttr("powermax_host.host_crud_test", "host_flags.spc2_protocol_version.override", "false"),
					resource.TestCheckResourceAttr("powermax_host.host_crud_test", "host_flags.volume_set_addressing.enabled", "true"),
					resource.TestCheckResourceAttr("powermax_host.host_crud_test", "host_flags.volume_set_addressing.override", "true"),
					resource.TestCheckResourceAttr("powermax_host.host_crud_test", "consistent_lun", "false")),
			},
			{
				Config:      UpdateHostExistingName,
				ExpectError: regexp.MustCompile(UpdateHostDetailsErrorMsg),
			},
			{
				Config: UpdateHostParams,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_host.host_crud_test", "name", TestAccHostName1Updated),
					resource.TestCheckResourceAttr("powermax_host.host_crud_test", "initiators.#", "1"),
					resource.TestCheckResourceAttr("powermax_host.host_crud_test", "initiators.0", InitiatorID2),
					resource.TestCheckResourceAttr("powermax_host.host_crud_test", "host_flags.spc2_protocol_version.enabled", "true"),
					resource.TestCheckResourceAttr("powermax_host.host_crud_test", "host_flags.spc2_protocol_version.override", "true"),
					resource.TestCheckResourceAttr("powermax_host.host_crud_test", "host_flags.volume_set_addressing.enabled", "true"),
					resource.TestCheckResourceAttr("powermax_host.host_crud_test", "host_flags.volume_set_addressing.override", "true"),
					resource.TestCheckResourceAttr("powermax_host.host_crud_test", "consistent_lun", "true")),
			},
			{
				Config:           ImportHostSuccess,
				ResourceName:     ImportHostResourceName1,
				ImportState:      true,
				ImportStateCheck: assertTFImportState,
				ExpectError:      nil,
				ImportStateId:    TestAccHostName1Updated,
			},
			{
				Config:        ImportHostFailure,
				ResourceName:  ImportHostResourceName2,
				ImportState:   true,
				ExpectError:   regexp.MustCompile(ImportHostDetailsErrorMsg),
				ImportStateId: "TestInvalidHost",
			},
		},
	})
}

func TestAccHost_CreateHostWithInvalidInitiator(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      CreateHostParamsWithInvalidInitiator,
				ExpectError: regexp.MustCompile(CreateHostDetailErrorMsg),
			},
		},
	})
}

var CreateHostParams = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "host_crud_test" {
	name = "` + TestAccHostName1 + `"
	host_flags = {
		volume_set_addressing = {
			override = true
			enabled = true
		}
		disable_q_reset_on_ua = {
			enabled = false
			override = false
		}
		environ_set = {
			enabled = false
			override = false
		}
		avoid_reset_broadcast = {
			enabled = true
			override = true
		}
		scsi_3 = {
			enabled = true
			override = true
		}
		openvms = {
			override = true
			enabled = false
		}
		spc2_protocol_version = {
			enabled = false
			override = false
		}
		scsi_support1 = {
			enabled = false
			override = true
		}
	}
	consistent_lun = false
	initiators = ["` + InitiatorID1 + `"]
}
`

var UpdateHostParams = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "host_crud_test" {
	name = "` + TestAccHostName1Updated + `"
	host_flags = {
		volume_set_addressing = {
			override = true
			enabled = true
		}
		disable_q_reset_on_ua = {
			enabled = false
			override = false
		}
		environ_set = {
			enabled = false
			override = false
		}
		avoid_reset_broadcast = {
			enabled = true
			override = true
		}
		scsi_3 = {
			enabled = true
			override = true
		}
		openvms = {
			override = true
			enabled = false
		}
		spc2_protocol_version = {
			enabled = true
			override = true
		}
		scsi_support1 = {
			enabled = false
			override = true
		}
	}
	consistent_lun = true
	initiators = ["` + InitiatorID2 + `"]
}
`

var UpdateHostExistingName = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "host_crud_test" {
	name = "` + HostID1 + `"
	host_flags = {
		volume_set_addressing = {
			override = true
			enabled = true
		}
		disable_q_reset_on_ua = {
			enabled = false
			override = false
		}
		environ_set = {
			enabled = false
			override = false
		}
		avoid_reset_broadcast = {
			enabled = true
			override = true
		}
		scsi_3 = {
			enabled = true
			override = true
		}
		openvms = {
			override = true
			enabled = false
		}
		spc2_protocol_version = {
			enabled = false
			override = true
		}
		scsi_support1 = {
			enabled = false
			override = true
		}
	}
	consistent_lun = false
	initiators = ["` + InitiatorID1 + `"]
}
`

var CreateHostParamsWithInvalidInitiator = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "host_create_test" {
	name = "` + TestAccHostName1 + `"
	host_flags = {
		volume_set_addressing = {
			override = true
			enabled = false
		}
		disable_q_reset_on_ua = {
			enabled = false
			override = false
		}
		environ_set = {
			enabled = false
			override = false
		}
		avoid_reset_broadcast = {
			enabled = false
			override = false
		}
		scsi_3 = {
			enabled = false
			override = false
		}
		openvms = {
			override = true
			enabled = true
		}
		spc2_protocol_version = {
			enabled = true
			override = true
		}
		scsi_support1 = {
			enabled = true
			override = true
		}
	}
	consistent_lun = true
	initiators = ["` + InvalidInitiatorID + `"]	    
}
`
var HostParams = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "host_create_rename_test" {
	name = "` + TestAccHostName4 + `"
	initiators = ["` + InitiatorID1 + `","` + InitiatorID2 + `"]
	host_flags = {
		volume_set_addressing = {
			override = true
			enabled = true
		}
		openvms = {
			override = true
			enabled = false
		}
	}
}
`

var HostParamsRename = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "host_create_rename_test" {
	name = "` + TestAccHostName5 + `"
	initiators = ["` + InitiatorID1 + `"]
	host_flags = {
		volume_set_addressing = {
			override = true
			enabled = true
		}
		openvms = {
			override = true
			enabled = false
		}
	}
}
`

var HostParamsChangeInitiatorRemove = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "host_create_rename_test" {
	name = "` + TestAccHostName4 + `"
	initiators = ["` + InitiatorID1 + `"]
	host_flags = {
		volume_set_addressing = {
			override = true
			enabled = true
		}
		openvms = {
			override = true
			enabled = false
		}
	}
}
`

var UpdateHostWithEmptyInitiatorFailure = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "host_create_rename_test" {
	name = "` + TestAccHostName4 + `"
	initiators = [""]
	host_flags = {
		volume_set_addressing = {
			override = true
			enabled = true
		}
		openvms = {
			override = true
			enabled = false
		}
	}
}
`

var HostParamsForUpdate = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "host_create_rename_test" {
	name = "` + TestAccHostName5 + `"
	initiators = ["` + InitiatorID1 + `"]
	host_flags = {
		volume_set_addressing = {
			override = true
			enabled = true
		}
		openvms = {
			override = true
			enabled = false
		}
	}
}
`

// Test duplicate initiators update along with host flags update
var HostParamsChangeAddInitiatorAndNameChange = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "host_create_rename_test" {
	name = "` + TestAccHostName4 + `"
	initiators = ["` + InitiatorID1 + `", "` + InitiatorID1 + `", "` + InitiatorID2 + `", "` + InitiatorID2 + `"]
	host_flags = {
		volume_set_addressing = {
			override = true
			enabled = true
		}
		openvms = {
			override = true
			enabled = false
		}
	}
}
`

var HostParamsChangeForUpdateFlags = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "host_create_rename_test" {
	name = "` + TestAccHostName5 + `"
	initiators = ["` + InitiatorID1 + `"]
	host_flags = {
		volume_set_addressing = {
			override = true
			enabled = true
		}
		openvms = {
			override = true
			enabled = false
		}
		avoid_reset_broadcast = {
			override = true
			enabled = true
		}
	}
}
`

var HostParamsChangeForUpdateFlagsInitiatorAndName = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "host_create_rename_test" {
	name = "` + TestAccHostName5 + `"
	initiators = ["` + InitiatorID1 + `"]
	host_flags = {
		volume_set_addressing = {
			override = true
			enabled = true
		}
		openvms = {
			override = true
			enabled = false
		}
		avoid_reset_broadcast = {
			override = true
			enabled = true
		}
		scsi_3 = {
			enabled = false
			override = true
		}
	}
}
`

var ImportHostSuccess = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "import_host_success" {
}
`

var ImportHostFailure = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "import_host_failure" {
}
`
