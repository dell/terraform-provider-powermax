package powermax

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	TestAccPGForMaskingView   = "test_acc_pg_maskingview"
	TestAccVolForMaskingView  = "test_acc_vol_maskingview"
	TestAccHostForMaskingView = "test_acc_host_maskingview"
	TestAccCreateMaskingView  = "test_acc_create_maskingview"
	TestAccUpdateMaskingView  = "test_acc_update_maskingview"
)

func TestAccMaskingView_CreateUpdateMaskingView(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
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
		},
	})
}

func TestAccMaskingView_CreateMaskingViewWithHostGroupID(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: CreateMaskingViewWithHostGroupSuccess,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_masking_view.create_maskingview_with_host_group", "id", TestAccCreateMaskingView)),
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
				Config:      CreateMaskingViewFailure,
				ExpectError: regexp.MustCompile(CreateMVDetailErrorMsg),
			},
		},
	})
}

func TestAccMaskingView_CreateMaskingViewError(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      CreateMaskingViewError,
				ExpectError: regexp.MustCompile(CreateMVDetailErrorMsg),
			},
		},
	})
}

func TestAccMaskingView_CreateMaskingViewErrorEmptyHostAndHostGroupID(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
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

resource "powermax_storage_group" "sg_for_maskingview" {
	name = "` + StorageGroupForMV2 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
}

resource "powermax_volume" "volume_for_maskingview" {
	name = "` + TestAccVolForMaskingView + `"
	size = 1
	cap_unit = "GB"
	sg_name = powermax_storage_group.sg_for_maskingview.name
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
	initiators = ["` + InitiatorID + `"]
	host_flags = {}

}

resource "powermax_masking_view" "create_update_maskingview" {
	name = "` + TestAccCreateMaskingView + `"
	storage_group_id = powermax_storage_group.sg_for_maskingview.id
	port_group_id = powermax_port_group.pg_for_maskingview.id
	host_id = powermax_host.host_for_maskingview.id
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

resource "powermax_masking_view" "create_maskingview_with_host_group" {
	name = "` + TestAccCreateMaskingView + `"
	storage_group_id = "` + StorageGroupForMV1 + `"
	port_group_id = powermax_port_group.pg_for_maskingview.id
	host_group_id = "` + HostGroupID1 + `"
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

resource "powermax_storage_group" "sg_for_maskingview" {
	name = "` + StorageGroupForMV2 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
}

resource "powermax_volume" "volume_for_maskingview" {
	name = "` + TestAccVolForMaskingView + `"
	size = 1
	cap_unit = "GB"
	sg_name = powermax_storage_group.sg_for_maskingview.name
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
	initiators = ["` + InitiatorID + `"]
	host_flags = {}

}

resource "powermax_masking_view" "create_update_maskingview" {
	name = "` + TestAccUpdateMaskingView + `"
	storage_group_id = powermax_storage_group.sg_for_maskingview.id
	port_group_id = powermax_port_group.pg_for_maskingview.id
	host_id = powermax_host.host_for_maskingview.id
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
	name = "` + StorageGroupForMV2 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
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
	initiators = ["` + InitiatorID + `"]
	host_flags = {}

}

resource "powermax_masking_view" "create_maskingview" {
	name = "` + TestAccCreateMaskingView + `"
	storage_group_id = powermax_storage_group.sg_for_maskingview.id
	port_group_id = powermax_port_group.pg_for_maskingview.id
	host_id = powermax_host.host_for_maskingview.id
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
