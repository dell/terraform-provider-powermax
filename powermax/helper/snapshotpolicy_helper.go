/*
Copyright (c) 2022-2023 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helper

import (
	"context"
	pmax "dell/powermax-go-client"
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ConvertToTimeString converts an int value to a string representation of time
func ConvertToTimeString(minutes int64) string {
	switch minutes {
	case 10, 12, 15, 20, 30:
		return strconv.FormatInt(minutes, 10) + " Minutes"
	case 60:
		return "1 Hour"
	case 120, 180, 240, 360, 480, 720:
		hours := minutes / 60
		return strconv.FormatInt(hours, 10) + " Hours"
	case 1440:
		return "1 Day"
	case 10080:
		return "7 Days"
	default:
		days := minutes / 1440
		if days > 0 {
			return strconv.FormatInt(days, 10) + " Days"
		}
		return "Invalid Input"
	}
}

// ConvertTimeStringToMinutes converts a string representation of time to an int value in minutes
func ConvertTimeStringToMinutes(timeStr string) (int64, error) {
	timeStr = strings.TrimSpace(timeStr)
	if strings.HasSuffix(timeStr, "Minutes") {
		minutesStr := strings.TrimSuffix(timeStr, " Minutes")
		return strconv.ParseInt(minutesStr, 10, 64)
	} else if strings.HasSuffix(timeStr, "Hour") {
		return 60, nil
	} else if strings.HasSuffix(timeStr, "Hours") {
		hoursStr := strings.TrimSuffix(timeStr, " Hours")
		hours, err := strconv.ParseInt(hoursStr, 10, 64)
		if err != nil {
			return 0, err
		}
		return hours * 60, nil
	} else if strings.HasSuffix(timeStr, "Day") {
		return 1440, nil
	} else if strings.HasSuffix(timeStr, "Days") {
		daysStr := strings.TrimSuffix(timeStr, " Days")
		days, err := strconv.ParseInt(daysStr, 10, 64)
		if err != nil {
			return 0, err
		}
		return days * 1440, nil
	} else {
		return 0, fmt.Errorf("Invalid Input")
	}
}

// UpdateSnapshotPolicyResourceState updates snapshot policy state
func UpdateSnapshotPolicyResourceState(ctx context.Context, snapshotPolicyDetail *pmax.SnapshotPolicy, state *models.SnapshotPolicyResource) error {
	err := CopyFields(ctx, snapshotPolicyDetail, state)
	state.Interval = types.StringValue(ConvertToTimeString(*snapshotPolicyDetail.IntervalMinutes))
	state.ID = types.StringValue(snapshotPolicyDetail.SnapshotPolicyName)
	if err != nil {
		return err
	}
	return nil
}

// ModifySnapshotPolicy modifies snapshot policy
func ModifySnapshotPolicy(ctx context.Context, client client.Client, plan *models.SnapshotPolicyResource, state *models.SnapshotPolicyResource) error {

	modifySnapshotPolicyParam := pmax.NewSnapshotPolicyModify()

	if plan.SnapshotPolicyName.ValueString() != state.SnapshotPolicyName.ValueString() {
		modifySnapshotPolicyParam.SetSnapshotPolicyName(plan.SnapshotPolicyName.ValueString())
	}

	if plan.Interval.ValueString() != "" && plan.Interval.ValueString() != state.Interval.ValueString() {
		mins, err := ConvertTimeStringToMinutes(plan.Interval.ValueString())
		if err != nil {
			tflog.Info(ctx, fmt.Sprintf("Error during converting time interval for Snapshot Policy Update: %s", err))
			return err
		}
		modifySnapshotPolicyParam.SetIntervalMins(mins)
	}

	if plan.OffsetMinutes.ValueInt64() != 0 && plan.OffsetMinutes != state.OffsetMinutes {
		modifySnapshotPolicyParam.SetOffsetMins(int32(plan.OffsetMinutes.ValueInt64()))
	}
	if plan.ComplianceCountCritical.ValueInt64() != state.ComplianceCountCritical.ValueInt64() {
		modifySnapshotPolicyParam.SetComplianceCountCritical(plan.ComplianceCountCritical.ValueInt64())
	}
	if plan.ComplianceCountWarning.ValueInt64() != state.ComplianceCountWarning.ValueInt64() {
		modifySnapshotPolicyParam.SetComplianceCountWarning(plan.ComplianceCountWarning.ValueInt64())
	}
	if plan.SnapshotCount.ValueInt64() != state.SnapshotCount.ValueInt64() {
		modifySnapshotPolicyParam.SetSnapshotCount(int32(plan.SnapshotCount.ValueInt64()))
	}
	snapshotPolicyUpdate := pmax.SnapshotPolicyUpdate{
		Action: "Modify",
		Modify: modifySnapshotPolicyParam,
	}

	updateReq := client.PmaxOpenapiClient.ReplicationApi.UpdateSnapshotPolicy(ctx, client.SymmetrixID, state.SnapshotPolicyName.ValueString())
	updateReq = updateReq.SnapshotPolicyUpdate(snapshotPolicyUpdate)
	_, _, err := updateReq.Execute()

	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error in update snapshot policy: %s", err))
		return err
	}

	return nil
}
