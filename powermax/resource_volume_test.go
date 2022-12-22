package powermax

import (
	"context"
	"log"
	"os"
	"regexp"
	"strings"
	"terraform-provider-powermax/client"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

// It is mandatory to create `test` resources with a prefix - 'test_acc_'
const (
	ImportVolumeResourceName1     = "powermax_volume.volume_import_success"
	ImportVolumeResourceName2     = "powermax_volume.volume_import_failure"
	TestAccCreateVolumeGB         = "test_acc_create_volume_gb"
	TestAccCreateVolumeCYL        = "test_acc_create_volume_cyl"
	TestAccCreateVolumeCYLUpdated = "test_acc_create_volume_cyl_updated"
	TestAccCreateVolumeTB1        = "test_acc_create_volume_tb1"
	TestAccCreateVolumeTB2        = "test_acc_create_volume_tb2"
	TestAccCreateVolumeGBUpdated  = "test_acc_create_volume_gb_updated"
	TestAccVolumeMobilityErr      = "test_acc_uvolume_gb_mv"
)

// Currently the storage groups used for test cases are - StorageGroupForMV1 and StorageGroupForVol1
// The sweeper is implemented according to the storage groups mentioned above
// In case of any changes in the usage of storage groups, changes must be implemented accordingly.
func init() {
	resource.AddTestSweepers("powermax_volume", &resource.Sweeper{
		Name:         "powermax_volume",
		Dependencies: []string{"powermax_masking_view"},
		F: func(region string) error {
			powermaxClient, err := getSweeperClient(region)
			if err != nil {
				log.Println("Error getting sweeper client")
				return nil
			}

			ctx := context.Background()
			deleteVolumeForSG(ctx, powermaxClient, StorageGroupForMV1)
			deleteVolumeForSG(ctx, powermaxClient, StorageGroupForVol1)

			return nil
		},
	})
}

func deleteVolumeForSG(ctx context.Context, powermaxClient *client.Client, storageGroup string) {
	volumeIDsForSG, err := powermaxClient.PmaxClient.GetVolumesInStorageGroupIterator(ctx, serialno, storageGroup)
	if err != nil {
		log.Println("Error getting volume list")
	}

	var volumeIDs []string
	for _, volumeIDList := range volumeIDsForSG.ResultList.VolumeList {
		volume, err := powermaxClient.PmaxClient.GetVolumeByID(ctx, serialno, volumeIDList.VolumeIDs)
		if err != nil {
			log.Println("Error getting volume")
			continue
		}
		if strings.Contains(volume.VolumeIdentifier, SweepTestsTemplateIdentifier) {
			volumeIDs = append(volumeIDs, volumeIDList.VolumeIDs)
		}
	}

	_, err = powermaxClient.PmaxClient.RemoveVolumesFromStorageGroup(ctx, serialno, storageGroup, true, volumeIDs...)
	if err != nil {
		log.Println("Error removing volume from storage group")
	}

	for _, volumeID := range volumeIDs {
		err := powermaxClient.PmaxClient.DeleteVolume(ctx, serialno, volumeID)
		if err != nil {
			log.Println("Error deleting volume")
		}
	}
}

func TestAccVolume_CRUDVolumeGB(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: CreateVolumeGB,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.crud_vol_gb", "name", TestAccCreateVolumeGB),
					resource.TestCheckResourceAttr("powermax_volume.crud_vol_gb", "size", "2.32"),
					resource.TestCheckResourceAttr("powermax_volume.crud_vol_gb", "cap_unit", "GB"),
					resource.TestCheckResourceAttr("powermax_volume.crud_vol_gb", "enable_mobility_id", "true")),
			},
			{
				// Failure Scenario: cannot update volume with lesser volume size
				Config:      UpdateVolumeGBError,
				ExpectError: regexp.MustCompile("Current volume size exceeds new volume size"),
			},
			{
				Config: UpdateVolumeGB,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.crud_vol_gb", "name", TestAccCreateVolumeGBUpdated),
					resource.TestCheckResourceAttr("powermax_volume.crud_vol_gb", "size", "2.5"),
					resource.TestCheckResourceAttr("powermax_volume.crud_vol_gb", "enable_mobility_id", "false")),
			},
			{
				Config: VolumeUpdateGbToTbCapUnit,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.crud_vol_gb", "name", TestAccCreateVolumeGBUpdated),
					resource.TestCheckResourceAttr("powermax_volume.crud_vol_gb", "size", "2.5"),
					resource.TestCheckResourceAttr("powermax_volume.crud_vol_gb", "cap_unit", "TB"),
					resource.TestCheckResourceAttr("powermax_volume.crud_vol_gb", "enable_mobility_id", "false")),
			},
		},
	})
}

func TestAccVolume_CreateVolumeCapUnitFailures(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      VolumeParamsWithMB,
				ExpectError: regexp.MustCompile("Unsupported capacity unit for volume size"),
			},
			{
				Config:      VolumeParamsWithInvalidCapUnit,
				ExpectError: regexp.MustCompile("Unsupported capacity unit for volume size"),
			},
		},
	})
}

func TestAccVolume_CreateVolumeWithTB(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: VolumeParamsWithTBInFloat,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.volume_create_test_tb_float", "name", TestAccCreateVolumeTB1),
					resource.TestCheckResourceAttr("powermax_volume.volume_create_test_tb_float", "size", "2.45"),
					resource.TestCheckResourceAttr("powermax_volume.volume_create_test_tb_float", "cap_unit", "TB")),
			},
			{
				Config: VolumeParamsWithTBInInt,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.volume_create_test_tb_int", "name", TestAccCreateVolumeTB2),
					resource.TestCheckResourceAttr("powermax_volume.volume_create_test_tb_int", "size", "2"),
					resource.TestCheckResourceAttr("powermax_volume.volume_create_test_tb_int", "cap_unit", "TB")),
			},
		},
	})
}

func TestAccVolume_CRUDVolumeWithCYL(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: CreateVolumeWithCYL,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.crud_volume_cyl", "name", TestAccCreateVolumeCYL),
					resource.TestCheckResourceAttr("powermax_volume.crud_volume_cyl", "size", "547"),
					resource.TestCheckResourceAttr("powermax_volume.crud_volume_cyl", "cap_unit", "CYL"),
					resource.TestCheckResourceAttr("powermax_volume.crud_volume_cyl", "enable_mobility_id", "false")),
			},
			{
				// Failure Scenario: cannot update size in decimal when cap unit is CYL
				Config:      UpdateVolumeWithCYLError,
				ExpectError: regexp.MustCompile("Failed to update all parameters"),
			},
			{
				Config: UpdateVolumeWithCYL,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.crud_volume_cyl", "name", TestAccCreateVolumeCYLUpdated),
					resource.TestCheckResourceAttr("powermax_volume.crud_volume_cyl", "size", "550"),
					resource.TestCheckResourceAttr("powermax_volume.crud_volume_cyl", "enable_mobility_id", "true")),
			},
		},
	})
}

func TestAccVolume_UpdateVolumeMobilityFailure(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: VolumeCreateForUpdateInMaskingView,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.update_volume_mobility_err", "name", TestAccVolumeMobilityErr),
					resource.TestCheckResourceAttr("powermax_volume.update_volume_mobility_err", "size", "2")),
			},
			{
				// Error scenario - mobility cannot be enabled when the volume is part of a storage group which is in masking view
				Config:      VolumeCreateForUpdateInMaskingViewError,
				ExpectError: regexp.MustCompile("operation cannot be performed because the device is mapped"),
			},
		},
	})
}

func TestAccVolume_ImportVolume(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	assertTFImportState := func(s []*terraform.InstanceState) error {
		assert.Equal(t, ImportVolumeName1, s[0].Attributes["name"])
		assert.Equal(t, "1", s[0].Attributes["size"])
		assert.Equal(t, "GB", s[0].Attributes["cap_unit"])
		assert.Equal(t, StorageGroupID1, s[0].Attributes["storagegroup_ids.0"])
		assert.Equal(t, 1, len(s))
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:           VolumeImportSuccess,
				ResourceName:     ImportVolumeResourceName1,
				ImportState:      true,
				ImportStateCheck: assertTFImportState,
				ExpectError:      nil,
				ImportStateId:    VolumeID1,
			},
			{
				Config:        VolumeImportFailure,
				ResourceName:  ImportVolumeResourceName2,
				ImportState:   true,
				ExpectError:   regexp.MustCompile(ImportVolDetailsErrorMsg),
				ImportStateId: "InvalidVolume",
			},
		},
	})
}

var CreateVolumeGB = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "crud_vol_gb" {
	name = "` + TestAccCreateVolumeGB + `"
	size = 2.32
	cap_unit = "GB"
	sg_name = "` + StorageGroupForVol1 + `"
	enable_mobility_id = true
}
`

var VolumeParamsRename = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "volume_create_test" {
	name = "test_acc_cvolume_updated"
	size = 2.32
	cap_unit = "GB"
	sg_name = "` + StorageGroupForVol1 + `"
}
`

var VolumeParamsWithMB = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "volume_create_test_mb" {
	name = "test_acc_cvolume_mb"
	size = 800
	cap_unit = "MB"
	sg_name = "` + StorageGroupForVol1 + `"
}
`

var VolumeParamsWithTBInFloat = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "volume_create_test_tb_float" {
	name = "` + TestAccCreateVolumeTB1 + `"
	size = 2.45
	cap_unit = "TB"
	sg_name = "` + StorageGroupForVol1 + `"
}
`

var VolumeParamsWithTBInInt = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "volume_create_test_tb_int" {
	name = "` + TestAccCreateVolumeTB2 + `"
	size = 2
	cap_unit = "TB"
	sg_name = "` + StorageGroupForVol1 + `"
}
`

var CreateVolumeWithCYL = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "crud_volume_cyl" {
	name = "` + TestAccCreateVolumeCYL + `"
	size = 547
	cap_unit = "CYL"
	sg_name = "` + StorageGroupForVol1 + `"
}
`

var VolumeParamsWithInvalidCapUnit = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "volume_create_test" {
	name = "test_acc_cvolume"
	size = 3
	cap_unit = "PB"
	sg_name = "` + StorageGroupForVol1 + `"
}
`

var VolumeCreateForUpdateCyl = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "volume_update_test_cyl" {
	name = "test_acc_uvolume_cyl"
	size = 500
	cap_unit = "CYL"
	sg_name = "` + StorageGroupForVol1 + `"
}
`

var UpdateVolumeWithCYL = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "crud_volume_cyl" {
	name = "` + TestAccCreateVolumeCYLUpdated + `"
	size = 550
	cap_unit = "CYL"
	sg_name = "` + StorageGroupForVol1 + `"
	enable_mobility_id = true
}
`

var VolumeUpdateGbToTbCapUnit = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "crud_vol_gb" {
	name = "` + TestAccCreateVolumeGBUpdated + `"
	size = 2.5
	cap_unit = "TB"
	sg_name = "` + StorageGroupForVol1 + `"
}
`

var VolumeUpdateGbToTbSize = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "volume_create_test" {
	name = "test_acc_cvolume"
	size = 2.5
	cap_unit = "TB"
	sg_name = "` + StorageGroupForVol1 + `"
}
`

// Error scenario - size when cap_unit is 'CYL' cannot be in float
var UpdateVolumeWithCYLError = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "crud_volume_cyl" {
	name = "` + TestAccCreateVolumeCYL + `"
	size = 500.5
	cap_unit = "CYL"
	sg_name = "` + StorageGroupForVol1 + `"
}
`

var VolumeCreateForUpdateGb = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "volume_update_test_gb" {
	name = "test_acc_uvolume_gb"
	size = 2
	cap_unit = "GB"
	sg_name = "` + StorageGroupForVol1 + `"
	enable_mobility_id = true
}
`

// Updates: size, enable_mobility_id
var UpdateVolumeGB = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "crud_vol_gb" {
	name = "` + TestAccCreateVolumeGBUpdated + `"
	size = 2.5
	cap_unit = "GB"
	sg_name = "` + StorageGroupForVol1 + `"
	enable_mobility_id = false
	
}
`

// Error scenario - Powermax APIs throw error when size is reduced in update
var UpdateVolumeGBError = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "crud_vol_gb" {
	name = "` + TestAccCreateVolumeGB + `"
	size = 1
	cap_unit = "GB"
	sg_name = "` + StorageGroupForVol1 + `"
	enable_mobility_id = true
	
}
`

var VolumeCreateForUpdateInMaskingView = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "update_volume_mobility_err" {
	name = "` + TestAccVolumeMobilityErr + `"
	size = 2
	cap_unit = "GB"
	sg_name = "` + StorageGroupForMV1 + `"
}
`

// Error scenario - mobility cannot be enabled when the volume is part of a storage group which is in masking view
var VolumeCreateForUpdateInMaskingViewError = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "update_volume_mobility_err" {
	name = "` + TestAccVolumeMobilityErr + `"
	size = 2
	cap_unit = "GB"
	sg_name = "` + StorageGroupForMV1 + `"
	enable_mobility_id = true
}
`

var VolumeImportSuccess = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "volume_import_success" {
}
`

var VolumeImportFailure = `
	provider "powermax" {
		username = "` + username + `"
		password = "` + password + `"
		endpoint = "` + endpoint + `"
		serial_number = "` + serialno + `"
		insecure = true
	}

	resource "powermax_volume" "volume_import_failure" {
	}
`
