package test

import (
	"context"
	"github.com/dell/gopowermax/v2/types/v100"
	"terraform-provider-powermax/powermax/helper"
	"terraform-provider-powermax/powermax/models"
	"testing"
)

func Test_CopyFields(t *testing.T) {
	sgResponse := v100.StorageGroup{
		StorageGroupID:        "test_sg",
		SLO:                   "Diamond",
		ServiceLevel:          "Diamond",
		BaseSLOName:           "Diamond",
		SRP:                   "SRP_1",
		Workload:              "",
		SLOCompliance:         "STABLE",
		NumOfVolumes:          1,
		NumOfChildSGs:         1,
		NumOfParentSGs:        1,
		NumOfMaskingViews:     1,
		NumOfSnapshots:        1,
		NumOfSnapshotPolicies: 1,
		CapacityGB:            50.0,
		DeviceEmulation:       "",
		Type:                  "Standalone",
		Unprotected:           false,
		ChildStorageGroup:     []string{"child1", "child2"},
		ParentStorageGroup:    nil,
		MaskingView:           nil,
		SnapshotPolicies:      nil,
		HostIOLimit: &v100.SetHostIOLimitsParam{
			HostIOLimitMBSec:    "100",
			HostIOLimitIOSec:    "100",
			DynamicDistribution: "STABLE",
		},
		Compression:           false,
		CompressionRatio:      "29.2:1",
		CompressionRatioToOne: 29.2,
		VPSavedPercent:        99.0,
		Tags:                  "",
		UUID:                  "",
		UnreducibleDataGB:     0,
	}

	var sg models.StorageGroupResourceModel
	err := helper.CopyFields(context.Background(), sgResponse, &sg)
	if err != nil {
		t.Errorf("CopyFields() error = %v", err)
	}
	t.Log(sg)
}
