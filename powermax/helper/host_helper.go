// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.

package helper

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/constants"
	"terraform-provider-powermax/powermax/models"

	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateHostState update host state.
func UpdateHostState(hostState *models.HostModel, planInitiators []string, hostResponse *pmaxTypes.Host) {
	hostState.HostID = types.StringValue(hostResponse.HostID)
	hostState.Name = types.StringValue(hostResponse.HostID)
	hostState.NumberHostGroups = types.Int64Value(hostResponse.NumberHostGroups)
	hostState.NumberInitiators = types.Int64Value(hostResponse.NumberInitiators)
	hostState.NumberMaskingViews = types.Int64Value(hostResponse.NumberMaskingViews)
	hostState.NumPowerPathHosts = types.Int64Value(hostResponse.NumPowerPathHosts)
	hostState.BWLimit = types.Int64Value(int64(hostResponse.BWLimit))
	hostState.HostType = types.StringValue(hostResponse.HostType)
	hostState.PortFlagsOverride = types.BoolValue(hostResponse.PortFlagsOverride)
	hostState.ConsistentLun = types.BoolValue(hostResponse.PortFlagsOverride)

	iniAttributeList := []attr.Value{}
	for _, ini := range hostResponse.Initiators {
		iniAttributeList = append(iniAttributeList, types.StringValue(ini))
	}
	hostState.Initiators, _ = types.ListValue(types.StringType, iniAttributeList)

	maskViewAttributeList := []attr.Value{}
	for _, id := range hostResponse.MaskingviewIDs {
		maskViewAttributeList = append(maskViewAttributeList, types.StringValue(id))
	}
	hostState.MaskingviewIDs, _ = types.ListValue(types.StringType, maskViewAttributeList)

	powerPathAttributeList := []attr.Value{}
	for _, id := range hostResponse.PowerPathHosts {
		powerPathAttributeList = append(maskViewAttributeList, types.StringValue(id))
	}
	hostState.PowerPathHosts, _ = types.ListValue(types.StringType, powerPathAttributeList)
	setDefaultHostFlags(hostState)
	setHostFlags(hostResponse.EnabledFlags, true, hostState)
	setHostFlags(hostResponse.DisabledFlags, false, hostState)

}

// UpdateHost update host and return updated parameters, failed updated parameters and errors.
func UpdateHost(ctx context.Context, client client.Client, plan, state models.HostModel) ([]string, []string, []string) {
	updatedParameters := []string{}
	updateFailedParameters := []string{}
	errorMessages := []string{}

	var planInitiators []string
	diags := plan.Initiators.ElementsAs(ctx, &planInitiators, true)
	if diags.HasError() {
		updateFailedParameters = append(updateFailedParameters, "initiators")
		errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify initiators: %s", "couldn't get the plan initiator data"))
	}

	var stateInitiators []string
	diags = state.Initiators.ElementsAs(ctx, &stateInitiators, true)
	if diags.HasError() {
		updateFailedParameters = append(updateFailedParameters, "initiators")
		errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify initiators: %s", "couldn't get the state initiator data"))
	}

	if !CompareStringSlice(planInitiators, stateInitiators) {
		hostResponse, err := client.PmaxClient.GetHostByID(ctx, client.SymmetrixID, state.HostID.ValueString())
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "initiators")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify initiators: %s for %s", "couldn't get the host data", state.HostID.ValueString()))
		}

		var planInitiatorsLowerCase []string
		for _, planInitiator := range planInitiators {
			planInitiatorsLowerCase = append(planInitiatorsLowerCase, strings.ToLower(planInitiator))
		}
		_, err = client.PmaxClient.UpdateHostInitiators(ctx, client.SymmetrixID, hostResponse, planInitiatorsLowerCase)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "initiators")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify initiators: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "initiators")
		}
	}

	if plan.HostFlags != state.HostFlags || plan.ConsistentLun.ValueBool() != state.ConsistentLun.ValueBool() {
		hostFlags := pmaxTypes.HostFlags{
			VolumeSetAddressing: &pmaxTypes.HostFlag{
				Enabled:  plan.HostFlags.VolumeSetAddressing.Enabled.ValueBool(),
				Override: plan.HostFlags.VolumeSetAddressing.Override.ValueBool(),
			},
			DisableQResetOnUA: &pmaxTypes.HostFlag{
				Enabled:  plan.HostFlags.DisableQResetOnUA.Enabled.ValueBool(),
				Override: plan.HostFlags.DisableQResetOnUA.Override.ValueBool(),
			},
			EnvironSet: &pmaxTypes.HostFlag{
				Enabled:  plan.HostFlags.EnvironSet.Enabled.ValueBool(),
				Override: plan.HostFlags.EnvironSet.Override.ValueBool(),
			},
			AvoidResetBroadcast: &pmaxTypes.HostFlag{
				Enabled:  plan.HostFlags.AvoidResetBroadcast.Enabled.ValueBool(),
				Override: plan.HostFlags.AvoidResetBroadcast.Override.ValueBool(),
			},
			OpenVMS: &pmaxTypes.HostFlag{
				Enabled:  plan.HostFlags.OpenVMS.Enabled.ValueBool(),
				Override: plan.HostFlags.OpenVMS.Override.ValueBool(),
			},
			SCSI3: &pmaxTypes.HostFlag{
				Enabled:  plan.HostFlags.SCSI3.Enabled.ValueBool(),
				Override: plan.HostFlags.SCSI3.Override.ValueBool(),
			},
			Spc2ProtocolVersion: &pmaxTypes.HostFlag{
				Enabled:  plan.HostFlags.Spc2ProtocolVersion.Enabled.ValueBool(),
				Override: plan.HostFlags.Spc2ProtocolVersion.Override.ValueBool(),
			},
			SCSISupport1: &pmaxTypes.HostFlag{
				Enabled:  plan.HostFlags.SCSISupport1.Enabled.ValueBool(),
				Override: plan.HostFlags.SCSISupport1.Override.ValueBool(),
			},
			ConsistentLUN: plan.ConsistentLun.ValueBool(),
		}
		_, err := client.PmaxClient.UpdateHostFlags(ctx, client.SymmetrixID, state.HostID.ValueString(), &hostFlags)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "host_flags")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify the host flags: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "host_flags")
		}
	}
	if plan.Name.ValueString() != state.Name.ValueString() {
		_, err := client.PmaxClient.UpdateHostName(ctx, client.SymmetrixID, state.HostID.ValueString(), plan.Name.ValueString())
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "name")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to rename host: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "name")
		}
	}

	return updatedParameters, updateFailedParameters, errorMessages
}

func setDefaultHostFlags(hostState *models.HostModel) {
	hostState.HostFlags.VolumeSetAddressing.Enabled = types.BoolValue(false)
	hostState.HostFlags.VolumeSetAddressing.Override = types.BoolValue(false)
	hostState.HostFlags.DisableQResetOnUA.Enabled = types.BoolValue(false)
	hostState.HostFlags.DisableQResetOnUA.Override = types.BoolValue(false)
	hostState.HostFlags.AvoidResetBroadcast.Enabled = types.BoolValue(false)
	hostState.HostFlags.AvoidResetBroadcast.Override = types.BoolValue(false)
	hostState.HostFlags.EnvironSet.Enabled = types.BoolValue(false)
	hostState.HostFlags.EnvironSet.Override = types.BoolValue(false)
	hostState.HostFlags.OpenVMS.Enabled = types.BoolValue(false)
	hostState.HostFlags.OpenVMS.Override = types.BoolValue(false)
	hostState.HostFlags.SCSISupport1.Enabled = types.BoolValue(false)
	hostState.HostFlags.SCSISupport1.Override = types.BoolValue(false)
	hostState.HostFlags.SCSI3.Enabled = types.BoolValue(false)
	hostState.HostFlags.SCSI3.Override = types.BoolValue(false)
	hostState.HostFlags.Spc2ProtocolVersion.Enabled = types.BoolValue(false)
	hostState.HostFlags.Spc2ProtocolVersion.Override = types.BoolValue(false)
}
func setHostFlags(flags string, isEnabled bool, hostState *models.HostModel) {
	if flags != "" {
		flagsArr := strings.Split(flags, ",")
		for _, flag := range flagsArr {
			switch flag {
			case constants.VolumeSetAddressing:
				hostState.HostFlags.VolumeSetAddressing.Enabled = types.BoolValue(isEnabled)
				hostState.HostFlags.VolumeSetAddressing.Override = types.BoolValue(true)
			case constants.DisableQResetOnUa:
				hostState.HostFlags.DisableQResetOnUA.Enabled = types.BoolValue(isEnabled)
				hostState.HostFlags.DisableQResetOnUA.Override = types.BoolValue(true)
			case constants.AvoidResetBroadcast:
				hostState.HostFlags.AvoidResetBroadcast.Enabled = types.BoolValue(isEnabled)
				hostState.HostFlags.AvoidResetBroadcast.Override = types.BoolValue(true)
			case constants.EnvironSet:
				hostState.HostFlags.EnvironSet.Enabled = types.BoolValue(isEnabled)
				hostState.HostFlags.EnvironSet.Override = types.BoolValue(true)
			case constants.OpenVMS:
				hostState.HostFlags.OpenVMS.Enabled = types.BoolValue(isEnabled)
				hostState.HostFlags.OpenVMS.Override = types.BoolValue(true)
			case constants.SCSISupport1:
				hostState.HostFlags.SCSISupport1.Enabled = types.BoolValue(isEnabled)
				hostState.HostFlags.SCSISupport1.Override = types.BoolValue(true)
			case constants.SCSI3:
				hostState.HostFlags.SCSI3.Enabled = types.BoolValue(isEnabled)
				hostState.HostFlags.SCSI3.Override = types.BoolValue(true)
			case constants.SPC2ProtocolVersion:
				hostState.HostFlags.Spc2ProtocolVersion.Enabled = types.BoolValue(isEnabled)
				hostState.HostFlags.Spc2ProtocolVersion.Override = types.BoolValue(true)
			}
		}
	}
}
