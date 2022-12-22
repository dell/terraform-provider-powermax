package powermax

import (
	"context"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

// It is mandatory to create `test` resources with a prefix - 'test_acc_'
const (
	TestAccSGName1        = "test_acc_csg_with_no_option_1"
	TestAccSGName2        = "test_acc_csg_1"
	TestAccSGName2Updated = "test_acc_csg_1_updated"
	TestAccSGName4        = "test_acc_sg_srp_id"
	TestAccSGName5        = "test_acc_sg_err"
	TestAccSGName8        = "test_acc_sg_volume_id_err"
	TestAccSGName9        = "test_acc_srp_slo_none"
	TestAccSGName10       = "test_acc_sg_with_already_attached_volume_id"
	TestAccSGName13       = "test_acc_valid_sg_with_vol"
	ResourceName1         = "powermax_storage_group.sg_import_success"
	ResourceName2         = "powermax_storage_group.sg_import_failure"
)

func init() {
	resource.AddTestSweepers("powermax_storage_group", &resource.Sweeper{
		Name:         "powermax_storage_group",
		Dependencies: []string{"powermax_masking_view", "powermax_volume"},
		F: func(region string) error {
			powermaxClient, err := getSweeperClient(region)
			if err != nil {
				log.Println("Error getting sweeper client: " + err.Error())
				return nil
			}

			ctx := context.Background()

			storageGroups, err := powermaxClient.PmaxClient.GetStorageGroupIDList(ctx, serialno)
			if err != nil {
				log.Println("Error getting storage group list: " + err.Error())
				return nil
			}

			for _, storageGroupID := range storageGroups.StorageGroupIDs {
				if strings.Contains(storageGroupID, SweepTestsTemplateIdentifier) {
					storageGroup, err := powermaxClient.PmaxClient.GetStorageGroup(ctx, serialno, storageGroupID)
					if err != nil {
						log.Println("Error getting storage group list: " + err.Error())
						return nil
					}

					if len(storageGroup.SnapshotPolicies) > 0 {
						payload := pmaxTypes.UpdateStorageGroupPayload{
							ExecutionOption: pmaxTypes.ExecutionOptionSynchronous,
							EditStorageGroupActionParam: pmaxTypes.EditStorageGroupActionParam{
								EditSnapshotPoliciesParam: &pmaxTypes.EditSnapshotPoliciesParam{
									DisassociateSnapshotPolicyParam: &pmaxTypes.SnapshotPolicies{
										SnapshotPolicies: storageGroup.SnapshotPolicies,
									},
								},
							},
						}

						err := powermaxClient.PmaxClient.UpdateStorageGroupS(ctx, serialno, storageGroupID, payload)
						if err != nil {
							log.Println("Error removing snapshot policies from storage group: " + storageGroupID + "with error: " + err.Error())
							return nil
						}
					}

					err = powermaxClient.PmaxClient.DeleteStorageGroup(ctx, serialno, storageGroupID)
					if err != nil {
						log.Println("Error deleting storage group: " + storageGroupID + "with error: " + err.Error())
					}
				}
			}
			return nil
		},
	})
}

func TestAccStorageGroup_CreateSGWithVolumeIdAttachedToAnotherSGSuccess(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: CreateValidSGWithSRPServiceLevelNone,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_storage_group.sg_with_srp_service_level_none", "id", TestAccSGName9),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_srp_service_level_none", "volume_ids.0", VolumeID3),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_srp_service_level_none", "srpid", "none"),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_srp_service_level_none", "service_level", "none")),
			},
			{
				// Test to verify create SG with volume ID which has already been attached to another SG, whose srpId and service_level is none
				Config: CreateValidSGWithAttachedVolumeID,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_storage_group.sg_with_already_attached_volume_id", "id", TestAccSGName10),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_already_attached_volume_id", "volume_ids.0", VolumeID3)),
			},
		},
	})
}

func TestAccStorageGroup_CreateSG_ExpectErr(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      StorageGroupWithExistingSGNameErr,
				ExpectError: regexp.MustCompile(CreateSgErrorMsg),
			},
			{
				// scenario: create SG
				// when: attaching a volume id to multiple SG
				// expected: error
				Config:      CreateStorageGroupAttachVolumeIDToMultipleSGErr,
				ExpectError: regexp.MustCompile(CreateSGAddVolumeErrMsg),
			},
			{
				Config:      CreateSGWithSuspendedSnapshotPolicyFailure,
				ExpectError: regexp.MustCompile(CreateSGDetailErrorMsg),
			},
			{
				Config:      StorageGroupExistingParam,
				ExpectError: regexp.MustCompile(CreateSGDetailErrorMsg),
			},
			{
				Config:      StorageGroupWithInvalidSRPID,
				ExpectError: regexp.MustCompile(CreateSGDetailErrorMsg),
			},
		},
	})
}

func TestAccStorageGroup_CRUD(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	assertTFImportState := func(s []*terraform.InstanceState) error {
		assert.Equal(t, TestAccSGName2Updated, s[0].Attributes["name"])
		assert.Equal(t, "false", s[0].Attributes["enable_compression"])
		assert.Equal(t, "Platinum", s[0].Attributes["service_level"])
		assert.Equal(t, ValidSrpID1, s[0].Attributes["srpid"])
		assert.Equal(t, SnapshotPolicy1, s[0].Attributes["snapshot_policies.0.policy_name"])
		assert.Equal(t, SnapshotPolicy2, s[0].Attributes["snapshot_policies.1.policy_name"])
		assert.Equal(t, "false", s[0].Attributes["snapshot_policies.0.is_active"])
		assert.Equal(t, "true", s[0].Attributes["snapshot_policies.1.is_active"])
		assert.Equal(t, "2", s[0].Attributes["snapshot_policies.#"])
		assert.Equal(t, "2", s[0].Attributes["volume_ids.#"])
		assert.Equal(t, 1, len(s))
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				// Test SG creation with snapshot policies, volume ids, service level and host_io_limits
				Config: StorageGroupOptionalParam,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "id", TestAccSGName2),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "service_level", "DiamoND"),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "volume_ids.0", VolumeID2),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "snapshot_policies.0.policy_name", SnapshotPolicy1),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "snapshot_policies.0.is_active", "true"),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "enable_compression", "true"),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "host_io_limits.dynamicdistribution", "Always")),
			},
			{
				// Test to verify SG update with snapshot policies(duplicates), volume ids(duplicates), service level, compressioion, host io limits, name,
				Config: StorageGroupUpdate,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "id", TestAccSGName2Updated),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "service_level", "Platinum"),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "volume_ids.1", VolumeID2),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "volume_ids.0", VolumeID3),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "volume_ids.#", "2"),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "snapshot_policies.0.policy_name", SnapshotPolicy1),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "snapshot_policies.0.is_active", "false"),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "snapshot_policies.1.policy_name", SnapshotPolicy2),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "snapshot_policies.1.is_active", "true"),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "enable_compression", "false"),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "host_io_limits.dynamicdistribution", "Never")),
			},
			{
				Config:           testAccImportStorageGroupSuccess,
				ResourceName:     ResourceName1,
				ImportState:      true,
				ImportStateCheck: assertTFImportState,
				ExpectError:      nil,
				ImportStateId:    TestAccSGName2Updated,
			},
			{
				Config:        testAccImportStorageGroupFailure,
				ResourceName:  ResourceName2,
				ImportState:   true,
				ExpectError:   regexp.MustCompile(ImportSGDetailsErrorMsg),
				ImportStateId: "InvalidStorageGroup",
			},
			{
				// Test to disassociate snapshot policies from SG so that it can be destroyed at end of test cycle
				Config: SnapshotPolicyRemove,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "id", TestAccSGName2Updated),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_with_optional_param", "snapshot_policies.#", "0")),
			},
		},
	})
}

func TestAccStorageGroup_UpdateSG_ExpectErr(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: CreateValidStorageGroup,
			},
			{
				Config:      StorageGroupWithInvalidServiceLevel,
				ExpectError: regexp.MustCompile(UpdateSGDetailsErrorMsg),
			},
			{
				// when: srpId = "none"
				// expected: service_level cannot be modified
				Config:      StorageGroupWithServiceLevelErr,
				ExpectError: regexp.MustCompile(UpdateSGDetailsErrorMsg),
			},
			{
				// when: srpId = "none", service_level = "none"
				// expected: enable_compression = false
				Config:      StorageGroupWithSrpCompressionErr,
				ExpectError: regexp.MustCompile(UpdateSGDetailsErrorMsg),
			},
			{
				// scenario: update SG
				// when: attaching a volume id to multiple SG
				// expected: error
				Config:      UpdateStorageGroupAttachVolumeIDToMultipleSGErr,
				ExpectError: regexp.MustCompile(UpdateSGDetailsErrorMsg),
			},
			{
				// scenario: update SG
				// when: associating a suspended snapshot policy to SG
				// expected: error
				Config:      UpdateStorageGroupAssociateSuspendedSnapshotPolicyErr,
				ExpectError: regexp.MustCompile(UpdateSGDetailsErrorMsg),
			},
			{
				Config:      StorageGroupUpdateNameWithExistingSG,
				ExpectError: regexp.MustCompile(UpdateSGDetailsErrorMsg),
			},
		},
	})
}

func TestAccStorageGroup_UpdateSrpID(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				// Test to verify create SG with SRP set to "none", which explicitly requires "enable_compression" set to false
				Config: StorageGroupWithoutSRPID,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_storage_group.sg_wo_srp_id", "enable_compression", "false"),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_wo_srp_id", "srpid", "none"),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_wo_srp_id", "service_level", "none")),
			},
			{
				// Test to verify update of SRP
				Config: StorageGroupWithSRPID,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("powermax_storage_group.sg_wo_srp_id", "enable_compression", "true"),
					resource.TestCheckResourceAttr("powermax_storage_group.sg_wo_srp_id", "srpid", ValidSrpID1)),
			},
		},
	})
}

var testAccImportStorageGroupFailure = `
	provider "powermax" {
		username = "` + username + `"
		password = "` + password + `"
		endpoint = "` + endpoint + `"
		serial_number = "` + serialno + `"
		insecure = true
	}

	resource "powermax_storage_group" "sg_import_failure" {
	}
`

var testAccImportStorageGroupSuccess = `
	provider "powermax" {
		username = "` + username + `"
		password = "` + password + `"
		endpoint = "` + endpoint + `"
		serial_number = "` + serialno + `"
		insecure = true
	}

	resource "powermax_storage_group" "sg_import_success" {
		name = "` + StorageGroupID1 + `"
		srpid = "` + ValidSrpID1 + `"
		service_level = "Diamond"
	}
`

var StorageGroupUpdateNameWithExistingSG = `

provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_with_invalid_update" {
	name = "` + StorageGroupID1 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
}

resource "powermax_storage_group" "sg_with_invalid_update1" {
	name = "` + TestAccSGName13 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
	volume_ids = ["` + VolumeID4 + `"]
}
`

var StorageGroupOptionalParam = `

provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_with_optional_param" {
	name = "` + TestAccSGName2 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "DiamoND"
	volume_ids = ["` + VolumeID2 + `"]
	host_io_limits = {
		host_io_limit_mb_sec = "1"
		host_io_limit_io_sec = "100"
		dynamicdistribution = "Always"
	}
	snapshot_policies = [
    {
      is_active = true
	  policy_name = "` + SnapshotPolicy1 + `"
    }
  ]
	
}
`

var StorageGroupUpdate = `

provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_with_optional_param" {
	name = "` + TestAccSGName2Updated + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Platinum"
	enable_compression = false
	volume_ids = ["` + VolumeID2 + `", "` + VolumeID2 + `", "` + VolumeID3 + `"]
	host_io_limits = {
		host_io_limit_mb_sec = "1"
		host_io_limit_io_sec = "100"
		dynamicdistribution = "Never"
	}
	snapshot_policies = [
		{
			is_active = false
			policy_name = "` + SnapshotPolicy1 + `"
		},
		{
			is_active = true
			policy_name = "` + SnapshotPolicy2 + `"
		},
		{
			is_active = true
			policy_name = "` + SnapshotPolicy2 + `"
		}
  ]
}
`

var SnapshotPolicyRemove = `

provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_with_optional_param" {
	name = "` + TestAccSGName2Updated + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
	volume_ids = ["` + VolumeID3 + `", "` + VolumeID2 + `"]
	host_io_limits = {
		host_io_limit_mb_sec = "1"
		host_io_limit_io_sec = "100"
		dynamicdistribution = "Never"
	}
	snapshot_policies = []
}
`

var StorageGroupExistingParam = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_existing_create1" {
	name = "` + StorageGroupID1 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
}
`
var StorageGroupWithoutSRPID = `

provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_wo_srp_id" {
	name = "` + TestAccSGName4 + `"
	srpid = "none"
	service_level = "none"
	enable_compression = false
}
`

var StorageGroupWithSRPID = `

provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_wo_srp_id" {
	name = "` + TestAccSGName4 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "none"
	host_io_limits = {
		host_io_limit_mb_sec = "1"
		host_io_limit_io_sec = "100"
		dynamicdistribution = "Always"
	}
}
`

var StorageGroupWithInvalidSRPID = `

provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_with_invalid_srp_id" {
	name = "` + TestAccSGName4 + `"
	srpid = "Invalid_SRP"
	service_level = "Diamond"
}
`

var CreateValidStorageGroup = `

provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_with_invalid_update" {
	name = "` + TestAccSGName5 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
}

resource "powermax_storage_group" "sg_with_invalid_update1" {
	name = "` + TestAccSGName13 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
	volume_ids = ["` + VolumeID4 + `"]
}
`

var CreateValidSGWithSRPServiceLevelNone = `

provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_with_srp_service_level_none" {
	name = "` + TestAccSGName9 + `"
	srpid = "none"
	service_level = "none"
	enable_compression = false
	volume_ids = ["` + VolumeID3 + `"]
}
`

var CreateValidSGWithAttachedVolumeID = `

provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_with_already_attached_volume_id" {
	name = "` + TestAccSGName10 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
	volume_ids = ["` + VolumeID3 + `"]
}
`

var CreateSGWithSuspendedSnapshotPolicyFailure = `
provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "storage_group_suspended_snapshot_policy" {
	name = "` + TestAccSGName2 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
	snapshot_policies = [
		{
			is_active = false
			policy_name = "` + SnapshotPolicy1 + `"
		}
	]
}
`

var StorageGroupWithInvalidServiceLevel = `

provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_with_invalid_update" {
	name = "` + TestAccSGName5 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "InvalidServiceLevel"
}

resource "powermax_storage_group" "sg_with_invalid_update1" {
	name = "` + TestAccSGName13 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
	volume_ids = ["` + VolumeID4 + `"]
}
`

var StorageGroupWithServiceLevelErr = `

provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_with_invalid_update" {
	name = "` + TestAccSGName5 + `"
	srpid = "none"
	service_level = "Platinum"
}

resource "powermax_storage_group" "sg_with_invalid_update1" {
	name = "` + TestAccSGName13 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
	volume_ids = ["` + VolumeID4 + `"]
}
`

var StorageGroupWithSrpCompressionErr = `

provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_with_invalid_update" {
	name = "` + TestAccSGName5 + `"
	srpid = "none"
	service_level = "none"
	enable_compression = true
}

resource "powermax_storage_group" "sg_with_invalid_update1" {
	name = "` + TestAccSGName13 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
	volume_ids = ["` + VolumeID4 + `"]
}
`

var StorageGroupWithExistingSGNameErr = `

provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_update_with_existing_name" {
	name = "` + StorageGroupForVol1 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
}

resource "powermax_storage_group" "sg_with_valid_volume" {
	name = "` + TestAccSGName13 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
	volume_ids = ["` + VolumeID4 + `"]
}
`

var CreateStorageGroupAttachVolumeIDToMultipleSGErr = `

provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_create_with_volume_id_err" {
	name = "` + TestAccSGName8 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
	volume_ids = ["` + VolumeID4 + `"]
}

resource "powermax_storage_group" "sg_with_valid_volume" {
	name = "` + TestAccSGName13 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
	volume_ids = ["` + VolumeID4 + `"]
}
`

var UpdateStorageGroupAttachVolumeIDToMultipleSGErr = `

provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_with_invalid_update" {
	name = "` + TestAccSGName5 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
	volume_ids = ["` + VolumeID4 + `"]
}


resource "powermax_storage_group" "sg_with_invalid_update1" {
	name = "` + TestAccSGName13 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
	volume_ids = ["` + VolumeID4 + `"]
}
`

var UpdateStorageGroupAssociateSuspendedSnapshotPolicyErr = `

provider "powermax" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	serial_number = "` + serialno + `"
	insecure = true
}

resource "powermax_storage_group" "sg_with_invalid_update" {
	name = "` + TestAccSGName5 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
	snapshot_policies = [
		{
			is_active = false
			policy_name = "` + SnapshotPolicy1 + `"
		},
  ]
}

resource "powermax_storage_group" "sg_with_invalid_update1" {
	name = "` + TestAccSGName13 + `"
	srpid = "` + ValidSrpID1 + `"
	service_level = "Diamond"
	volume_ids = ["` + VolumeID4 + `"]
}
`
