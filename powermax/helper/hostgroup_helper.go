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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func UpdateHostGroupState(hostGroupState *models.HostGroupModal, hostGroupResponse *pmaxTypes.HostGroup) {
	hostGroupState.ID = types.StringValue(hostGroupResponse.HostGroupID)
	hostGroupState.Name = types.StringValue(hostGroupResponse.HostGroupID)
	hostGroupState.NumOfHosts = types.Int64Value(hostGroupResponse.NumOfHosts)
	hostGroupState.NumOfInitiators = types.Int64Value(hostGroupResponse.NumberInitiators)
	hostGroupState.NumOfMaskingViews = types.Int64Value(hostGroupResponse.NumberMaskingViews)
	hostGroupState.Type = types.StringValue(hostGroupResponse.HostGroupType)
	hostGroupState.PortFlagsOverride = types.BoolValue(hostGroupResponse.PortFlagsOverride)
	hostGroupState.ConsistentLun = types.BoolValue(hostGroupResponse.ConsistentLun)

	responseHostIDs := []string{}
	for _, hostSummary := range hostGroupResponse.Hosts {
		responseHostIDs = append(responseHostIDs, hostSummary.HostID)
	}

	saveHgListAttribute(hostGroupState, responseHostIDs, "hostIDs")
	saveHgListAttribute(hostGroupState, hostGroupResponse.MaskingviewIDs, "maskingViews")
	setDefaultHostFlagsForHostGroup(hostGroupState)
	setHostFlagsInHg(hostGroupResponse.EnabledFlags, true, hostGroupState)
	setHostFlagsInHg(hostGroupResponse.DisabledFlags, false, hostGroupState)
}

func setHostFlagsInHg(flags string, isEnabled bool, hostState *models.HostGroupModal) {
	if flags != "" {
		flagsArr := strings.Split(flags, ",")
		for _, flag := range flagsArr {
			switch flag {
			case constants.VolumeSetAddressing:
				hostState.HostFlags.VolumeSetAddressing.Enabled = basetypes.NewBoolValue(isEnabled)
				hostState.HostFlags.VolumeSetAddressing.Override = basetypes.NewBoolValue(true)
			case constants.DisableQResetOnUa:
				hostState.HostFlags.DisableQResetOnUA.Enabled = basetypes.NewBoolValue(isEnabled)
				hostState.HostFlags.DisableQResetOnUA.Override = basetypes.NewBoolValue(true)
			case constants.AvoidResetBroadcast:
				hostState.HostFlags.AvoidResetBroadcast.Enabled = basetypes.NewBoolValue(isEnabled)
				hostState.HostFlags.AvoidResetBroadcast.Override = basetypes.NewBoolValue(true)
			case constants.EnvironSet:
				hostState.HostFlags.EnvironSet.Enabled = basetypes.NewBoolValue(isEnabled)
				hostState.HostFlags.EnvironSet.Override = basetypes.NewBoolValue(true)
			case constants.OpenVMS:
				hostState.HostFlags.OpenVMS.Enabled = basetypes.NewBoolValue(isEnabled)
				hostState.HostFlags.OpenVMS.Override = basetypes.NewBoolValue(true)
			case constants.SCSISupport1:
				hostState.HostFlags.SCSISupport1.Enabled = basetypes.NewBoolValue(isEnabled)
				hostState.HostFlags.SCSISupport1.Override = basetypes.NewBoolValue(true)
			case constants.SCSI3:
				hostState.HostFlags.SCSI3.Enabled = basetypes.NewBoolValue(isEnabled)
				hostState.HostFlags.SCSI3.Override = basetypes.NewBoolValue(true)
			case constants.SPC2ProtocolVersion:
				hostState.HostFlags.Spc2ProtocolVersion.Enabled = basetypes.NewBoolValue(isEnabled)
				hostState.HostFlags.Spc2ProtocolVersion.Override = basetypes.NewBoolValue(true)
			}
		}
	}
}

func setDefaultHostFlagsForHostGroup(hostState *models.HostGroupModal) {
	hostState.HostFlags.VolumeSetAddressing.Enabled = basetypes.NewBoolValue(false)
	hostState.HostFlags.VolumeSetAddressing.Override = basetypes.NewBoolValue(false)
	hostState.HostFlags.DisableQResetOnUA.Enabled = basetypes.NewBoolValue(false)
	hostState.HostFlags.DisableQResetOnUA.Override = basetypes.NewBoolValue(false)
	hostState.HostFlags.AvoidResetBroadcast.Enabled = basetypes.NewBoolValue(false)
	hostState.HostFlags.AvoidResetBroadcast.Override = basetypes.NewBoolValue(false)
	hostState.HostFlags.EnvironSet.Enabled = basetypes.NewBoolValue(false)
	hostState.HostFlags.EnvironSet.Override = basetypes.NewBoolValue(false)
	hostState.HostFlags.OpenVMS.Enabled = basetypes.NewBoolValue(false)
	hostState.HostFlags.OpenVMS.Override = basetypes.NewBoolValue(false)
	hostState.HostFlags.SCSISupport1.Enabled = basetypes.NewBoolValue(false)
	hostState.HostFlags.SCSISupport1.Override = basetypes.NewBoolValue(false)
	hostState.HostFlags.SCSI3.Enabled = basetypes.NewBoolValue(false)
	hostState.HostFlags.SCSI3.Override = basetypes.NewBoolValue(false)
	hostState.HostFlags.Spc2ProtocolVersion.Enabled = basetypes.NewBoolValue(false)
	hostState.HostFlags.Spc2ProtocolVersion.Override = basetypes.NewBoolValue(false)
}

func saveHgListAttribute(hostGroupState *models.HostGroupModal, listAttribute []string, attributeName string) diag.Diagnostics {
	var attributeListType types.List
	var err diag.Diagnostics
	if len(listAttribute) > 0 {
		var attributeList []attr.Value
		for _, attribute := range listAttribute {
			attributeList = append(attributeList, types.StringValue(attribute))
		}
		attributeListType, err = types.ListValue(types.StringType, attributeList)
	} else {
		// Empty List
		attributeListType, err = types.ListValue(types.StringType, []attr.Value{})
	}
	if attributeName == "maskingViews" {
		hostGroupState.Maskingviews = attributeListType
	} else if attributeName == "hostIDs" {
		hostGroupState.HostIDs = types.Set(attributeListType)
	}

	return err
}

func UpdateHostGroup(ctx context.Context, client client.Client, plan, state models.HostGroupModal) ([]string, []string, []string) {
	updatedParameters := []string{}
	updateFailedParameters := []string{}
	errorMessages := []string{}

	var planHostIDs []string
	diags := plan.HostIDs.ElementsAs(ctx, &planHostIDs, true)
	if diags.HasError() {
		updateFailedParameters = append(updateFailedParameters, "host_ids")
		errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify host_ids: %s", "couldn't get the plan host_ids data"))
	}

	var stateHostIDs []string
	diags = state.HostIDs.ElementsAs(ctx, &stateHostIDs, true)
	if diags.HasError() {
		updateFailedParameters = append(updateFailedParameters, "host_ids")
		errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify host_ids: %s", "couldn't get the state host_ids data"))
	}

	if !CompareStringSlice(planHostIDs, stateHostIDs) {
		hostGroupID := state.ID.ValueString()
		_, err := client.PmaxClient.UpdateHostGroupHosts(ctx, client.SymmetrixID, hostGroupID, planHostIDs)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "host_ids")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify host_ids: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "host_ids")
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
		_, err := client.PmaxClient.UpdateHostGroupFlags(ctx, client.SymmetrixID, state.Name.ValueString(), &hostFlags)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "host_flags")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify the host flags: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "host_flags")
		}
	}

	if plan.Name.ValueString() != state.Name.ValueString() {
		_, err := client.PmaxClient.UpdateHostGroupName(ctx, client.SymmetrixID, state.ID.ValueString(), plan.Name.ValueString())
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "name")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to rename hostGroup: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "name")
		}
	}

	return updatedParameters, updateFailedParameters, errorMessages
}
