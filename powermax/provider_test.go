package powermax

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// var testAccProviders map[string]func() tfsdk.Provider
var testProvider tfsdk.Provider
var testProviderFactory map[string]func() (tfprotov6.ProviderServer, error)
var endpoint = os.Getenv("UNISPHERE_HOST")
var serialno = os.Getenv("POWERMAX_SERIALNO")
var username = os.Getenv("UNISPHERE_USERNAME")
var password = os.Getenv("UNISPHERE_PASSWORD")
var InitiatorID1 = os.Getenv("HOST_INITIATOR_ID1")
var InitiatorID2 = os.Getenv("HOST_INITIATOR_ID2")
var StorageGroupID1 = os.Getenv("STORAGE_GROUP_ID1")
var VolumeID1 = os.Getenv("VOLUME_ID1")
var VolumeID2 = os.Getenv("VOLUME_ID2")
var VolumeID3 = os.Getenv("VOLUME_ID3")
var VolumeID4 = os.Getenv("VOLUME_ID4")
var ValidSrpID1 = os.Getenv("SRP_ID1")
var ImportVolumeName1 = os.Getenv("IMPORT_VOLUME_NAME1")
var SnapshotPolicy1 = os.Getenv("SNAPSHOT_POLICY1")
var SnapshotPolicy2 = os.Getenv("SNAPSHOT_POLICY2")
var DirectorID1 = os.Getenv("DIRECTOR_ID1")
var DirectorID2 = os.Getenv("DIRECTOR_ID2")
var StorageGroupForVol1 = os.Getenv("STORAGE_GROUP_VOL1")
var StorageGroupForMV1 = os.Getenv("STORAGE_GROUP_MV1")
var StorageGroupForMV2 = os.Getenv("STORAGE_GROUP_MV2")
var StorageGroup1 = os.Getenv("STORAGE_GROUP1")
var StorageGroup2 = os.Getenv("STORAGE_GROUP2")
var HostGroupID1 = os.Getenv("HOST_GROUP_ID1")
var ImportHostInitiatorID = os.Getenv("HOST_INITIATOR_ID3")

func init() {
	testProvider = New("test")()
	testProvider.Configure(context.Background(), tfsdk.ConfigureProviderRequest{}, &tfsdk.ConfigureProviderResponse{})
	testProviderFactory = map[string]func() (tfprotov6.ProviderServer, error){
		"powermax": providerserver.NewProtocol6WithError(testProvider),
	}
}

//lint:ignore U1000 used by the internal provider, to be checked
func testAccProvider(t *testing.T, p tfsdk.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providers := p.(*provider)
		if !providers.configured {
			return fmt.Errorf("provider not configured")
		}

		if providers.client.PmaxClient == nil {
			return fmt.Errorf("provider not configured")
		}
		return nil
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("UNISPHERE_USERNAME"); v == "" {
		t.Fatal("UNISPHERE_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("UNISPHERE_PASSWORD"); v == "" {
		t.Fatal("UNISPHERE_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("UNISPHERE_HOST"); v == "" {
		t.Fatal("UNISPHERE_HOST must be set for acceptance tests")
	}

	if v := os.Getenv("POWERMAX_SERIALNO"); v == "" {
		t.Fatal("POWERMAX_SERIALNO must be set for acceptance tests")
	}

	if v := os.Getenv("HOST_INITIATOR_ID1"); v == "" {
		t.Fatal("HOST_INITIATOR_ID1 must be set for acceptance tests")
	}

	if v := os.Getenv("HOST_INITIATOR_ID2"); v == "" {
		t.Fatal("HOST_INITIATOR_ID2 must be set for acceptance tests")
	}

	if v := os.Getenv("STORAGE_GROUP_ID1"); v == "" {
		t.Fatal("STORAGE_GROUP_ID1 must be set for acceptance tests")
	}

	if v := os.Getenv("VOLUME_ID1"); v == "" {
		t.Fatal("VOLUME_ID1 must be set for acceptance tests")
	}

	if v := os.Getenv("VOLUME_ID2"); v == "" {
		t.Fatal("VOLUME_ID2 must be set for acceptance tests")
	}

	if v := os.Getenv("VOLUME_ID3"); v == "" {
		t.Fatal("VOLUME_ID3 must be set for acceptance tests")
	}

	if v := os.Getenv("VOLUME_ID4"); v == "" {
		t.Fatal("VOLUME_ID4 must be set for acceptance tests")
	}

	if v := os.Getenv("SRP_ID1"); v == "" {
		t.Fatal("SRP_ID1 must be set for acceptance tests")
	}

	if v := os.Getenv("IMPORT_VOLUME_NAME1"); v == "" {
		t.Fatal("IMPORT_VOLUME_NAME1 must be set for acceptance tests")
	}

	if v := os.Getenv("SNAPSHOT_POLICY1"); v == "" {
		t.Fatal("SNAPSHOT_POLICY1 must be set for acceptance tests")
	}

	if v := os.Getenv("SNAPSHOT_POLICY2"); v == "" {
		t.Fatal("SNAPSHOT_POLICY2 must be set for acceptance tests")
	}

	if v := os.Getenv("DIRECTOR_ID1"); v == "" {
		t.Fatal("DIRECTOR_ID1 must be set for acceptance tests")
	}

	if v := os.Getenv("DIRECTOR_ID2"); v == "" {
		t.Fatal("DIRECTOR_ID2 must be set for acceptance tests")
	}

	if v := os.Getenv("STORAGE_GROUP_VOL1"); v == "" {
		t.Fatal("STORAGE_GROUP_VOL1 must be set for acceptance tests")
	}

	if v := os.Getenv("STORAGE_GROUP_MV1"); v == "" {
		t.Fatal("STORAGE_GROUP_MV1 must be set for acceptance tests")
	}

	if v := os.Getenv("STORAGE_GROUP_MV2"); v == "" {
		t.Fatal("STORAGE_GROUP_MV2 must be set for acceptance tests")
	}

	if v := os.Getenv("STORAGE_GROUP1"); v == "" {
		t.Fatal("STORAGE_GROUP1 must be set for acceptance tests")
	}

	if v := os.Getenv("STORAGE_GROUP2"); v == "" {
		t.Fatal("STORAGE_GROUP2 must be set for acceptance tests")
	}

	if v := os.Getenv("HOST_GROUP_ID1"); v == "" {
		t.Fatal("HOST_GROUP_ID1 must be set for acceptance tests")
	}
}

const EmptyEndpointConfig = `
provider "powermax" {
	username = "username"
	password = "password"
	serial_number = "serial_number"
}
`
