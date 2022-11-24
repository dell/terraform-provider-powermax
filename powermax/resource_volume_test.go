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
	"github.com/stretchr/testify/assert"
)

const (
	ImportVolumeResourceName1 = "powermax_volume.volume_import_success"
	ImportVolumeResourceName2 = "powermax_volume.volume_import_failure"
)

func TestAccVolume_CreateVolume(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: VolumeParams,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.volume_create_test", "name", "test_acc_cvol"),
					resource.TestCheckResourceAttr("powermax_volume.volume_create_test", "size", "2.32"),
					resource.TestCheckResourceAttr("powermax_volume.volume_create_test", "cap_unit", "GB")),
			},
		},
	})
}

func TestAccVolume_CreateVolumeWithMB(t *testing.T) {
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
		},
	})
}

func TestAccVolume_CreateVolumeWithTBInFloat(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: VolumeParamsWithTBInFloat,
				Check:  resource.ComposeTestCheckFunc(checkCreateVolume(t, testProvider, StorageGroup1, "test_acc_cvol_tb_float", "2.45", "TB")),
			},
		},
	})
}

func TestAccVolume_CreateVolumeWithTBInInt(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: VolumeParamsWithTBInInt,
				Check:  resource.ComposeTestCheckFunc(checkCreateVolume(t, testProvider, StorageGroup1, "test_acc_cvol_tb", "2", "TB")),
			},
		},
	})
}

func TestAccVolume_CreateVolumeWithCYL(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: VolumeParamsWithCYL,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.volume_create_test_cyl", "name", "test_acc_cvol_cyl"),
					resource.TestCheckResourceAttr("powermax_volume.volume_create_test_cyl", "size", "547"),
					resource.TestCheckResourceAttr("powermax_volume.volume_create_test_cyl", "cap_unit", "CYL")),
			},
		},
	})
}

func TestAccVolume_CreateVolumeWithInvalidCapUnit(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      VolumeParamsWithInvalidCapUnit,
				ExpectError: regexp.MustCompile("Unsupported capacity unit for volume size"),
			},
		},
	})
}

func TestAccVolume_UpdateVolumeCyl(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: VolumeCreateForUpdateCyl,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.volume_update_test_cyl", "name", "test_acc_uvol_cyl"),
					resource.TestCheckResourceAttr("powermax_volume.volume_update_test_cyl", "size", "500")),
			},
			{
				Config: VolumeUpdateCyl,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.volume_update_test_cyl", "name", "test_acc_uvol_cyl_updated"),
					resource.TestCheckResourceAttr("powermax_volume.volume_update_test_cyl", "size", "550"),
					resource.TestCheckResourceAttr("powermax_volume.volume_update_test_cyl", "enable_mobility_id", "true")),
			},
		},
	})
}

func TestAccVolume_UpdateVolumeRename(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: VolumeParams,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.volume_create_test", "name", "test_acc_cvol"),
					resource.TestCheckResourceAttr("powermax_volume.volume_create_test", "size", "2.32"),
					resource.TestCheckResourceAttr("powermax_volume.volume_create_test", "cap_unit", "GB")),
			},
			{
				Config: VolumeParamsRename,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.volume_create_test", "name", "test_acc_cvol_updated")),
			},
		},
	})
}

func TestAccVolume_UpdateVolumeCylError(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: VolumeCreateForUpdateCyl,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.volume_update_test_cyl", "name", "test_acc_uvol_cyl"),
					resource.TestCheckResourceAttr("powermax_volume.volume_update_test_cyl", "size", "500")),
			},
			{
				Config:      VolumeUpdateCylError,
				ExpectError: regexp.MustCompile("Failed to update all parameters"),
			},
		},
	})
}

func TestAccVolume_UpdateVolumeGb(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: VolumeCreateForUpdateGb,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.volume_update_test_gb", "name", "test_acc_uvol_gb"),
					resource.TestCheckResourceAttr("powermax_volume.volume_update_test_gb", "size", "2")),
			},
			{
				Config: VolumeUpdateGb,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.volume_update_test_gb", "name", "test_acc_uvol_gb_updated"),
					resource.TestCheckResourceAttr("powermax_volume.volume_update_test_gb", "size", "2.5"),
					resource.TestCheckResourceAttr("powermax_volume.volume_update_test_gb", "enable_mobility_id", "false")),
			},
		},
	})
}

func TestAccVolume_UpdateVolumeGbError1(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: VolumeCreateForUpdateGb,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.volume_update_test_gb", "name", "test_acc_uvol_gb"),
					resource.TestCheckResourceAttr("powermax_volume.volume_update_test_gb", "size", "2")),
			},
			{
				Config:      VolumeUpdateGbError,
				ExpectError: regexp.MustCompile("Current volume size exceeds new volume size"),
			},
		},
	})
}

func TestAccVolume_UpdateVolumeGbError2(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: VolumeCreateForUpdateInMaskingView,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.volume_update_test_gb_mv", "name", "test_acc_uvol_gb_mv"),
					resource.TestCheckResourceAttr("powermax_volume.volume_update_test_gb_mv", "size", "2")),
			},
			{
				Config:      VolumeCreateForUpdateInMaskingViewError,
				ExpectError: regexp.MustCompile("operation cannot be performed because the device is mapped"),
			},
		},
	})
}

func TestAccVolume_UpdateVolumeSizeGbToTb(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: VolumeParams,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.volume_create_test", "name", "test_acc_cvol"),
					resource.TestCheckResourceAttr("powermax_volume.volume_create_test", "size", "2.32"),
					resource.TestCheckResourceAttr("powermax_volume.volume_create_test", "cap_unit", "GB")),
			},
			{
				Config: VolumeUpdateGbToTb,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_volume.volume_create_test", "size", "0.5"),
					(resource.TestCheckResourceAttr("powermax_volume.volume_create_test", "cap_unit", "TB"))),
			},
		},
	})
}

func TestAccVolume_ImportVolumeSuccess(t *testing.T) {
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
		},
	})
}

func TestAccVolume_ImportVolumeFailure(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:        VolumeImportFailure,
				ResourceName:  ImportVolumeResourceName2,
				ImportState:   true,
				ExpectError:   regexp.MustCompile(ImportVolDetailsErrorMsg),
				ImportStateId: "testVolumeImport",
			},
		},
	})
}

func checkCreateVolume(t *testing.T, p tfsdk.Provider, sgName string, volName string, size string, capUnit string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providers := p.(*provider)
		_, err := providers.client.PmaxClient.GetVolumeByIdentifier(context.Background(), serialno, sgName, volName, size, capUnit)
		if err != nil {
			return fmt.Errorf("failed to fetch volume")
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

var VolumeParams = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "volume_create_test" {
	name = "test_acc_cvol"
	size = 2.32
	cap_unit = "GB"
	sg_name = "` + StorageGroupForVol1 + `"
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
	name = "test_acc_cvol_updated"
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
	name = "test_acc_cvol_mb"
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
	name = "test_acc_cvol_tb_float"
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

resource "powermax_volume" "volume_create_test_tb" {
	name = "test_acc_cvol_tb"
	size = 2
	cap_unit = "TB"
	sg_name = "` + StorageGroupForVol1 + `"
}
`

var VolumeParamsWithCYL = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "volume_create_test_cyl" {
	name = "test_acc_cvol_cyl"
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
	name = "test_acc_cvol"
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
	name = "test_acc_uvol_cyl"
	size = 500
	cap_unit = "CYL"
	sg_name = "` + StorageGroupForVol1 + `"
}
`

var VolumeUpdateCyl = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "volume_update_test_cyl" {
	name = "test_acc_uvol_cyl_updated"
	size = 550
	cap_unit = "CYL"
	sg_name = "` + StorageGroupForVol1 + `"
	enable_mobility_id = true
}
`

var VolumeUpdateGbToTb = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "volume_create_test" {
	name = "test_acc_cvol"
	size = 0.5
	cap_unit = "TB"
	sg_name = "` + StorageGroupForVol1 + `"
}
`

// Error scenario - size when cap_unit is 'CYL' cannot be in float
var VolumeUpdateCylError = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "volume_update_test_cyl" {
	name = "test_acc_uvol_cyl_updated"
	size = 500.5
	cap_unit = "CYL"
	sg_name = "` + StorageGroupForVol1 + `"
	enable_mobility_id = true
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
	name = "test_acc_uvol_gb"
	size = 2
	cap_unit = "GB"
	sg_name = "` + StorageGroupForVol1 + `"
	enable_mobility_id = true
}
`

var VolumeUpdateGb = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "volume_update_test_gb" {
	name = "test_acc_uvol_gb_updated"
	size = 2.5
	cap_unit = "GB"
	sg_name = "` + StorageGroupForVol1 + `"
	enable_mobility_id = false
	
}
`

// Error scenario - Powermax APIs throw error when size is reduced in update
var VolumeUpdateGbError = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_volume" "volume_update_test_gb" {
	name = "test_acc_uvol_gb"
	size = 1
	cap_unit = "GB"
	sg_name = "` + StorageGroupForVol1 + `"
	enable_mobility_id = false
	
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

resource "powermax_volume" "volume_update_test_gb_mv" {
	name = "test_acc_uvol_gb_mv"
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

resource "powermax_volume" "volume_update_test_gb_mv" {
	name = "test_acc_uvol_gb_mv"
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
	name = "` + ImportVolumeName1 + `"
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
