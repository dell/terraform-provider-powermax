package powermax

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/models"

	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func updateHostGroupState(hostGroupState *models.HostGroup, hostGroupResponse *pmaxTypes.HostGroup) {
	hostGroupState.ID = types.String{Value: hostGroupResponse.HostGroupID}
	hostGroupState.Name = types.String{Value: hostGroupResponse.HostGroupID}
	hostGroupState.NumOfHosts = types.Int64{Value: hostGroupResponse.NumOfHosts}
	hostGroupState.NumOfInitiators = types.Int64{Value: hostGroupResponse.NumberInitiators}
	hostGroupState.NumOfMaskingViews = types.Int64{Value: hostGroupResponse.NumberMaskingViews}
	hostGroupState.Type = types.String{Value: hostGroupResponse.HostGroupType}
	hostGroupState.PortFlagsOverride = types.Bool{Value: hostGroupResponse.PortFlagsOverride}
	hostGroupState.ConsistentLun = types.Bool{Value: hostGroupResponse.ConsistentLun}

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

func setHostFlagsInHg(flags string, isEnabled bool, hostState *models.HostGroup) {
	if flags != "" {
		flagsArr := strings.Split(flags, ",")
		for _, flag := range flagsArr {
			switch flag {
			case VolumeSetAddressing:
				hostState.HostFlags.VolumeSetAddressing.Enabled.Value = isEnabled
				hostState.HostFlags.VolumeSetAddressing.Override.Value = true
			case DisableQResetOnUa:
				hostState.HostFlags.DisableQResetOnUa.Enabled.Value = isEnabled
				hostState.HostFlags.DisableQResetOnUa.Override.Value = true
			case AvoidResetBroadcast:
				hostState.HostFlags.AvoidResetBroadcast.Enabled.Value = isEnabled
				hostState.HostFlags.AvoidResetBroadcast.Override.Value = true
			case EnvironSet:
				hostState.HostFlags.EnvironSet.Enabled.Value = isEnabled
				hostState.HostFlags.EnvironSet.Override.Value = true
			case OpenVMS:
				hostState.HostFlags.Openvms.Enabled.Value = isEnabled
				hostState.HostFlags.Openvms.Override.Value = true
			case SCSISupport1:
				hostState.HostFlags.ScsiSupport1.Enabled.Value = isEnabled
				hostState.HostFlags.ScsiSupport1.Override.Value = true
			case SCSI3:
				hostState.HostFlags.Scsi3.Enabled.Value = isEnabled
				hostState.HostFlags.Scsi3.Override.Value = true
			case SPC2ProtocolVersion:
				hostState.HostFlags.Spc2ProtocolVersion.Enabled.Value = isEnabled
				hostState.HostFlags.Spc2ProtocolVersion.Override.Value = true
			}
		}
	}
}

func setDefaultHostFlagsForHostGroup(hostState *models.HostGroup) {
	hostState.HostFlags.VolumeSetAddressing.Enabled = types.Bool{Value: false}
	hostState.HostFlags.VolumeSetAddressing.Override = types.Bool{Value: false}
	hostState.HostFlags.DisableQResetOnUa.Enabled = types.Bool{Value: false}
	hostState.HostFlags.DisableQResetOnUa.Override = types.Bool{Value: false}
	hostState.HostFlags.AvoidResetBroadcast.Enabled = types.Bool{Value: false}
	hostState.HostFlags.AvoidResetBroadcast.Override = types.Bool{Value: false}
	hostState.HostFlags.EnvironSet.Enabled = types.Bool{Value: false}
	hostState.HostFlags.EnvironSet.Override = types.Bool{Value: false}
	hostState.HostFlags.Openvms.Enabled = types.Bool{Value: false}
	hostState.HostFlags.Openvms.Override = types.Bool{Value: false}
	hostState.HostFlags.ScsiSupport1.Enabled = types.Bool{Value: false}
	hostState.HostFlags.ScsiSupport1.Override = types.Bool{Value: false}
	hostState.HostFlags.Scsi3.Enabled = types.Bool{Value: false}
	hostState.HostFlags.Scsi3.Override = types.Bool{Value: false}
	hostState.HostFlags.Spc2ProtocolVersion.Enabled = types.Bool{Value: false}
	hostState.HostFlags.Spc2ProtocolVersion.Override = types.Bool{Value: false}
}

func saveHgListAttribute(hostGroupState *models.HostGroup, listAttribute []string, attributeName string) {
	var attributeListType types.List
	if len(listAttribute) > 0 {
		var attributeList []attr.Value
		for _, attribute := range listAttribute {
			attributeList = append(attributeList, types.String{Value: attribute})
		}
		attributeListType = types.List{
			ElemType: types.StringType,
			Elems:    attributeList,
		}
	} else {
		attributeListType = types.List{
			ElemType: types.StringType,
			Elems:    []attr.Value{},
		}
	}
	if attributeName == "maskingViews" {
		hostGroupState.Maskingviews = attributeListType
	} else if attributeName == "hostIDs" {
		hostGroupState.HostIDs = types.Set(attributeListType)
	}
}

func updateHostGroup(ctx context.Context, client client.Client, plan, state models.HostGroup) ([]string, []string, []string) {
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

	if !compareStringSlice(planHostIDs, stateHostIDs) {
		hostGroupID := state.ID.Value
		_, err := client.PmaxClient.UpdateHostGroupHosts(ctx, client.SymmetrixID, hostGroupID, planHostIDs)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "host_ids")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify host_ids: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "host_ids")
		}
	}

	if plan.HostFlags != state.HostFlags || plan.ConsistentLun.Value != state.ConsistentLun.Value {
		hostFlags := pmaxTypes.HostFlags{
			VolumeSetAddressing: &pmaxTypes.HostFlag{
				Enabled:  plan.HostFlags.VolumeSetAddressing.Enabled.Value,
				Override: plan.HostFlags.VolumeSetAddressing.Override.Value,
			},
			DisableQResetOnUA: &pmaxTypes.HostFlag{
				Enabled:  plan.HostFlags.DisableQResetOnUa.Enabled.Value,
				Override: plan.HostFlags.DisableQResetOnUa.Override.Value,
			},
			EnvironSet: &pmaxTypes.HostFlag{
				Enabled:  plan.HostFlags.EnvironSet.Enabled.Value,
				Override: plan.HostFlags.EnvironSet.Override.Value,
			},
			AvoidResetBroadcast: &pmaxTypes.HostFlag{
				Enabled:  plan.HostFlags.AvoidResetBroadcast.Enabled.Value,
				Override: plan.HostFlags.AvoidResetBroadcast.Override.Value,
			},
			OpenVMS: &pmaxTypes.HostFlag{
				Enabled:  plan.HostFlags.Openvms.Enabled.Value,
				Override: plan.HostFlags.Openvms.Override.Value,
			},
			SCSI3: &pmaxTypes.HostFlag{
				Enabled:  plan.HostFlags.Scsi3.Enabled.Value,
				Override: plan.HostFlags.Scsi3.Override.Value,
			},
			Spc2ProtocolVersion: &pmaxTypes.HostFlag{
				Enabled:  plan.HostFlags.Spc2ProtocolVersion.Enabled.Value,
				Override: plan.HostFlags.Spc2ProtocolVersion.Override.Value,
			},
			SCSISupport1: &pmaxTypes.HostFlag{
				Enabled:  plan.HostFlags.ScsiSupport1.Enabled.Value,
				Override: plan.HostFlags.ScsiSupport1.Override.Value,
			},
			ConsistentLUN: plan.ConsistentLun.Value,
		}
		_, err := client.PmaxClient.UpdateHostGroupFlags(ctx, client.SymmetrixID, state.Name.Value, &hostFlags)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "host_flags")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify the host flags: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "host_flags")
		}
	}

	if plan.Name.Value != state.Name.Value {
		_, err := client.PmaxClient.UpdateHostGroupName(ctx, client.SymmetrixID, state.ID.Value, plan.Name.Value)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "name")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to rename hostGroup: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "name")
		}
	}

	return updatedParameters, updateFailedParameters, errorMessages
}
