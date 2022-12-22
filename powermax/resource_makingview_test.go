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
	TestAccPGForMaskingView        = "test_acc_pg_maskingview"
	TestAccHostForMaskingView      = "test_acc_host_maskingview"
	TestAccCreateMaskingView       = "test_acc_create_maskingview"
	TestAccUpdateMaskingView       = "test_acc_update_maskingview"
	TestAccSGEmptyVolume           = "test_acc_sg_masking_view"
	ImportMaskingViewResourceName1 = "powermax_masking_view.import_masking_view_success"
	ImportMaskingViewResourceName2 = "powermax_masking_view.import_masking_view_failure"
)

func init() {
	resource.AddTestSweepers("powermax_masking_view", &resource.Sweeper{
		Name: "powermax_masking_view",
		F: func(region string) error {
			powermaxClient, err := getSweeperClient(region)
			if err != nil {
				log.Println("Error getting sweeper client: " + err.Error())
				return nil
			}

			ctx := context.Background()

			maskingViews, err := powermaxClient.PmaxClient.GetMaskingViewList(ctx, serialno)
			if err != nil {
				log.Println("Error getting masking view list: " + err.Error())
				return nil
			}

			for _, maskingView := range maskingViews.MaskingViewIDs {
				if strings.Contains(maskingView, SweepTestsTemplateIdentifier) {
					err := powermaxClient.PmaxClient.DeleteMaskingView(ctx, serialno, maskingView)
					if err != nil {
						log.Println("Error deleting maskingview: " + maskingView + "with error: " + err.Error())
					}
				}
			}
			return nil
		},
	})
}

func TestAccMaskingView_CreateUpdateMaskingView(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	assertTFImportState := func(s []*terraform.InstanceState) error {
		assert.Equal(t, TestAccUpdateMaskingView, s[0].Attributes["name"])
		assert.Equal(t, StorageGroupForMV1, s[0].Attributes["storage_group_id"])
		assert.Equal(t, PortGroupID1, s[0].Attributes["port_group_id"])
		assert.Equal(t, HostID1, s[0].Attributes["host_id"])
		assert.Equal(t, "", s[0].Attributes["host_group_id"])
		assert.Equal(t, 1, len(s))
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: CreateMaskingViewSuccess,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_masking_view.create_update_maskingview", "id", TestAccCreateMaskingView)),
			},
			{
				Config: RenameMaskingViewSuccess,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_masking_view.create_update_maskingview", "id", TestAccUpdateMaskingView)),
			},
			{
				Config:      UpdateMaskingViewHostNameFailure,
				ExpectError: regexp.MustCompile("Error updating maskingview"),
			},
			{
				Config:           ImportMaskingViewSuccess,
				ResourceName:     ImportMaskingViewResourceName1,
				ImportState:      true,
				ImportStateCheck: assertTFImportState,
				ExpectError:      nil,
				ImportStateId:    TestAccUpdateMaskingView,
			},
			{
				Config:        ImportMaskingViewFailure,
				ResourceName:  ImportMaskingViewResourceName2,
				ImportState:   true,
				ExpectError:   regexp.MustCompile(ImportMVDetailsErrorMsg),
				ImportStateId: "TestInvalidMaskingView",
			},
		},
	})
}

func TestAccMaskingView_CreateMaskingViewWithHostGroupID(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	assertTFImportState := func(s []*terraform.InstanceState) error {
		assert.Equal(t, TestAccCreateMaskingView, s[0].Attributes["name"])
		assert.Equal(t, TestAccHostGroupName1, s[0].Attributes["host_group_id"])
		assert.Equal(t, 1, len(s))
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: CreateMaskingViewWithHostGroupSuccess,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_masking_view.create_maskingview_with_host_group", "id", TestAccCreateMaskingView)),
			},
			{
				Config:           ImportMaskingViewSuccess,
				ResourceName:     ImportMaskingViewResourceName1,
				ImportState:      true,
				ImportStateCheck: assertTFImportState,
				ExpectError:      nil,
				ImportStateId:    TestAccCreateMaskingView,
			},
		},
	})
}

func TestAccMaskingView_CreateMaskingViewCaseInsensitive(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: CreateMaskingViewCaseInsensitive,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_masking_view.create_maskingview_with_host_group", "id", TestAccCreateMaskingView),
					resource.TestCheckResourceAttr("powermax_masking_view.create_maskingview_with_host_group", "storage_group_id", strings.ToLower(StorageGroupForMV1))),
			},
		},
	})
}

func TestAccMaskingView_CreateMaskingViewFailure(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				// Failure scenario -  Create masking view fails when storageGroup does not have any volumes associated with it.
				Config:      CreateMaskingViewFailure,
				ExpectError: regexp.MustCompile(CreateMVDetailErrorMsg),
			},
			{
				// Failure scenario -  Cannot create masking view when both host_id and host_group_id are given as input.
				Config:      CreateMaskingViewError,
				ExpectError: regexp.MustCompile(CreateMVDetailErrorMsg),
			},
			{
				// Failure scenario - cannot create masking view with empty host and host group ID
				Config:      CreateMaskingViewErrorEmptyHostAndHostGroupID,
				ExpectError: regexp.MustCompile(CreateMVDetailErrorMsg),
			},
		},
	})
}

var CreateMaskingViewSuccess = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}



resource "powermax_masking_view" "create_update_maskingview" {
	name = "` + TestAccCreateMaskingView + `"
	storage_group_id = "` + StorageGroupForMV1 + `"
	port_group_id = "` + PortGroupID1 + `"
	host_id = "` + HostID1 + `"
}
`

var CreateMaskingViewWithHostGroupSuccess = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_port_group" "pg_for_maskingview" {
	name = "` + TestAccPGForMaskingView + `"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "` + DirectorID1 + `"
			port_id = "2"
		}
	]
}

resource "powermax_host" "host_create_test_1" {
	name = "` + TestAccHostForHG1 + `"
	host_flags = {
	}
	initiators = ["` + InitiatorID1 + `"]
}

resource "powermax_host_group" "hg_for_maskingview" {
	name = "` + TestAccHostGroupName1 + `"
	host_flags = {
		spc2_protocol_version = {
			enabled = false
			override = true
		}
	}
	host_ids = [powermax_host.host_create_test_1.id]
}

resource "powermax_masking_view" "create_maskingview_with_host_group" {
	name = "` + TestAccCreateMaskingView + `"
	storage_group_id = "` + StorageGroupForMV1 + `"
	port_group_id = powermax_port_group.pg_for_maskingview.id
	host_group_id = powermax_host_group.hg_for_maskingview.id
}
`

var CreateMaskingViewCaseInsensitive = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_port_group" "pg_for_maskingview" {
	name = "` + TestAccPGForMaskingView + `"
	protocol = "SCSI_FC"
	ports = [
		{
			director_id = "` + DirectorID1 + `"
			port_id = "2"
		}
	]
}

resource "powermax_host" "host_for_maskingview" {
	name = "` + TestAccHostForMaskingView + `"
	initiators = ["` + InitiatorID1 + `"]
	host_flags = {}

}

resource "powermax_masking_view" "create_maskingview_with_host_group" {
	name = "` + TestAccCreateMaskingView + `"
	storage_group_id = lower("` + StorageGroupForMV1 + `")
	port_group_id = powermax_port_group.pg_for_maskingview.id
	host_id = powermax_host.host_for_maskingview.id
}
`

var RenameMaskingViewSuccess = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_masking_view" "create_update_maskingview" {
	name = "` + TestAccUpdateMaskingView + `"
	storage_group_id = "` + StorageGroupForMV1 + `"
	port_group_id = "` + PortGroupID1 + `"
	host_id = "` + HostID1 + `"
}
`

var UpdateMaskingViewHostNameFailure = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_masking_view" "create_update_maskingview" {
	name = "` + TestAccUpdateMaskingView + `"
	storage_group_id = "` + StorageGroupForMV1 + `"
	port_group_id = "` + PortGroupID1 + `"
	host_id = "test-host-id"
}
`

// Failure scenario -  Create masking view fails when storageGroup does not have any volumes associated with it.
var CreateMaskingViewFailure = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_for_maskingview" {
	name = "` + TestAccSGEmptyVolume + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
}

resource "powermax_masking_view" "create_maskingview" {
	name = "` + TestAccCreateMaskingView + `"
	storage_group_id = powermax_storage_group.sg_for_maskingview.id
	port_group_id = "` + PortGroupID1 + `"
	host_id = "` + HostID1 + `"
}
`

// Error scenario -  Cannot create masking view when both host_id and host_group_id are given as input.
var CreateMaskingViewError = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_masking_view" "create_maskingview" {
	name = "` + TestAccCreateMaskingView + `"
	storage_group_id = "test_storage_group"
	port_group_id = "test_port_group"
	host_id = "test_host"
	host_group_id = "test_host_group"
}
`

var CreateMaskingViewErrorEmptyHostAndHostGroupID = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_masking_view" "create_maskingview" {
	name = "` + TestAccCreateMaskingView + `"
	storage_group_id = "test_storage_group"
	port_group_id = "test_port_group"
}
`

var ImportMaskingViewSuccess = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_masking_view" "import_masking_view_success" {
}
`

var ImportMaskingViewFailure = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_masking_view" "import_masking_view_failure" {
}
`
