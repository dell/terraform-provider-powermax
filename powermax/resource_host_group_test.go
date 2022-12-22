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
	TestAccHostGroupName1        = "test_acc_hostgroup_1"
	TestAccHostGroupName2        = "test_acc_hostgroup_2"
	TestAccHostGroupNameUpdated1 = "test_acc_hostgroup_updated_1"
	TestAccHostForHG1            = "test_acc_host_hg_1"
	TestAccHostForHG2            = "test_acc_host_hg_2"
	ImportHostGroupResourceName1 = "powermax_host_group.import_host_group_success"
	ImportHostGroupResourceName2 = "powermax_host_group.import_host_group_failure"
)

func init() {
	resource.AddTestSweepers("powermax_host_group", &resource.Sweeper{
		Name:         "powermax_host_group",
		Dependencies: []string{"powermax_masking_view", "powermax_host"},
		F: func(region string) error {
			powermaxClient, err := getSweeperClient(region)
			if err != nil {
				log.Println("Error getting sweeper client")
				return nil
			}

			ctx := context.Background()

			hostgroups, err := powermaxClient.PmaxClient.GetHostGroupList(ctx, serialno)
			if err != nil {
				log.Println("Error getting hostgroup list")
				return nil
			}

			for _, hostGroup := range hostgroups.HostGroupIDs {
				if strings.Contains(hostGroup, SweepTestsTemplateIdentifier) {
					err := powermaxClient.PmaxClient.DeleteHostGroup(ctx, serialno, hostGroup)
					if err != nil {
						log.Println("Error deleting hostgroup")
					}
				}
			}
			return nil
		},
	})
}

func TestAccHostGroup_HostGroupCRUD(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	assertTFImportState := func(s []*terraform.InstanceState) error {
		assert.Equal(t, TestAccHostGroupNameUpdated1, s[0].Attributes["name"])
		assert.Equal(t, "1", s[0].Attributes["host_ids.#"])
		assert.Equal(t, TestAccHostForHG2, s[0].Attributes["host_ids.0"])
		assert.Equal(t, "true", s[0].Attributes["host_flags.spc2_protocol_version.enabled"])
		assert.Equal(t, "true", s[0].Attributes["host_flags.spc2_protocol_version.override"])
		assert.Equal(t, 1, len(s))
		return nil
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactory,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: CreateHostGroup,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_host_group.hostgroup_create_test_1", "name", TestAccHostGroupName1),
					resource.TestCheckResourceAttr("powermax_host_group.hostgroup_create_test_1", "host_flags.volume_set_addressing.enabled", "false"),
					resource.TestCheckResourceAttr("powermax_host_group.hostgroup_create_test_1", "host_flags.volume_set_addressing.override", "false"),
					resource.TestCheckResourceAttr("powermax_host_group.hostgroup_create_test_1", "host_flags.spc2_protocol_version.enabled", "false"),
					resource.TestCheckResourceAttr("powermax_host_group.hostgroup_create_test_1", "host_flags.spc2_protocol_version.override", "true"),
					resource.TestCheckResourceAttr("powermax_host_group.hostgroup_create_test_1", "host_ids.#", "1"),
					resource.TestCheckResourceAttr("powermax_host_group.hostgroup_create_test_1", "host_ids.0", TestAccHostForHG1),
				),
			},
			{
				Config:      UpdateHostGroupWithMixCaseHostIDFailure,
				ExpectError: regexp.MustCompile(UpdateHostGroupDetailsErrorMsg),
			},
			{
				Config:      UpdateHostGroupExistingNameFailure,
				ExpectError: regexp.MustCompile(UpdateHostGroupDetailsErrorMsg),
			},
			{
				Config: UpdateHostGroupSuccess,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_host_group.hostgroup_create_test_1", "name", TestAccHostGroupNameUpdated1),
					resource.TestCheckResourceAttr("powermax_host_group.hostgroup_create_test_1", "host_flags.volume_set_addressing.enabled", "false"),
					resource.TestCheckResourceAttr("powermax_host_group.hostgroup_create_test_1", "host_flags.volume_set_addressing.override", "false"),
					resource.TestCheckResourceAttr("powermax_host_group.hostgroup_create_test_1", "host_flags.spc2_protocol_version.enabled", "true"),
					resource.TestCheckResourceAttr("powermax_host_group.hostgroup_create_test_1", "host_flags.spc2_protocol_version.override", "true"),
					resource.TestCheckResourceAttr("powermax_host_group.hostgroup_create_test_1", "host_ids.#", "1"),
					resource.TestCheckResourceAttr("powermax_host_group.hostgroup_create_test_1", "host_ids.0", TestAccHostForHG2)),
			},
			{
				Config:           ImportHostGroupSuccess,
				ResourceName:     ImportHostGroupResourceName1,
				ImportState:      true,
				ImportStateCheck: assertTFImportState,
				ExpectError:      nil,
				ImportStateId:    TestAccHostGroupNameUpdated1,
			},
			{
				Config:        ImportHostGroupFailure,
				ResourceName:  ImportHostGroupResourceName2,
				ImportState:   true,
				ExpectError:   regexp.MustCompile(ImportHostGroupDetailsErrorMsg),
				ImportStateId: "TestInvalidHostGroup",
			},
		},
	})
}

func TestAccHostGroup_CreateHostGroupWithErrorCases(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      CreateHostGroupWithInvalidHostIDs,
				ExpectError: regexp.MustCompile(CreateHostGroupDetailErrorMsg),
			},
			{
				Config:      CreateHostGroupWithEmptyHostIDs,
				ExpectError: regexp.MustCompile(MinimumSizeValidationError),
			},
		},
	})
}

var CreateHostGroup = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "host_create_test_1" {
	name = "` + TestAccHostForHG1 + `"
	host_flags = {
	}
	initiators = ["` + InitiatorID1 + `"]
}

resource "powermax_host_group" "hostgroup_create_test_1" {
	name = "` + TestAccHostGroupName1 + `"
	host_flags = {
		spc2_protocol_version = {
			enabled = false
			override = true
		}
	}
	host_ids = [powermax_host.host_create_test_1.id]
}
`

var UpdateHostGroupWithMixCaseHostIDFailure = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "host_create_test_1" {
	name = "` + TestAccHostForHG1 + `"
	host_flags = {
	}
	initiators = ["` + InitiatorID1 + `"]
}

resource "powermax_host_group" "hostgroup_create_test_1" {
	name = "` + TestAccHostGroupName1 + `"
	host_flags = {
		spc2_protocol_version = {
			enabled = false
			override = true
		}
	}
	host_ids = ["Test_Acc_host_hg_1"]
}

resource "powermax_host_group" "hostgroup_create_test_2" {
	name = "` + TestAccHostGroupName2 + `"
	host_flags = {
		scsi_3 = {
			enabled = false
			override = true
		}
	}
	host_ids = [powermax_host.host_create_test_1.id]
}
`

var UpdateHostGroupExistingNameFailure = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "host_create_test_1" {
	name = "` + TestAccHostForHG1 + `"
	host_flags = {
	}
	initiators = ["` + InitiatorID1 + `"]
}

resource "powermax_host_group" "hostgroup_create_test_1" {
	name = "` + TestAccHostGroupName2 + `"
	host_flags = {
		spc2_protocol_version = {
			enabled = false
			override = true
		}
	}
	host_ids = [powermax_host.host_create_test_1.id]
}

resource "powermax_host_group" "hostgroup_create_test_2" {
	name = "` + TestAccHostGroupName2 + `"
	host_flags = {
		scsi_3 = {
			enabled = false
			override = true
		}
	}
	host_ids = [powermax_host.host_create_test_1.id]
}
`

var UpdateHostGroupSuccess = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host" "host_create_test_1" {
	name = "` + TestAccHostForHG1 + `"
	host_flags = {
	}
	initiators = ["` + InitiatorID1 + `"]
}

resource "powermax_host" "host_create_test_2" {
	name = "` + TestAccHostForHG2 + `"
	host_flags = {
	}
	initiators = ["` + InitiatorID2 + `"]
}

resource "powermax_host_group" "hostgroup_create_test_1" {
	name = "` + TestAccHostGroupNameUpdated1 + `"
	host_flags = {
		spc2_protocol_version = {
			enabled = true
			override = true
		}
	}
	host_ids = [powermax_host.host_create_test_2.id]
}
`

var CreateHostGroupWithInvalidHostIDs = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host_group" "hostgroup_create_test_1" {
	name = "` + TestAccHostGroupName1 + `"
	host_flags = {
		spc2_protocol_version = {
			enabled = true
			override = true
		}
	}
	host_ids = ["invalid-host-id"]
}
`

var CreateHostGroupWithEmptyHostIDs = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host_group" "hostgroup_create_test_1" {
	name = "` + TestAccHostGroupName1 + `"
	host_flags = {
		spc2_protocol_version = {
			enabled = true
			override = true
		}
	}
	host_ids = []
}
`

var ImportHostGroupSuccess = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host_group" "import_host_group_success" {
}
`

var ImportHostGroupFailure = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_host_group" "import_host_group_failure" {
}
`
