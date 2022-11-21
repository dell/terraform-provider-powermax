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

func updateHostState(hostState *models.Host, planInitiators []string, hostResponse *pmaxTypes.Host) {
	hostState.ID = types.String{Value: hostResponse.HostID}
	hostState.Name = types.String{Value: hostResponse.HostID}
	hostState.NumOfHostGroups = types.Int64{Value: hostResponse.NumberHostGroups}
	hostState.NumOfInitiators = types.Int64{Value: hostResponse.NumberInitiators}
	hostState.NumOfMaskingViews = types.Int64{Value: hostResponse.NumberMaskingViews}
	hostState.NumOfPowerpathHosts = types.Int64{Value: hostResponse.NumPowerPathHosts}
	hostState.BWLimit = types.Int64{Value: int64(hostResponse.BWLimit)}
	hostState.Type = types.String{Value: hostResponse.HostType}
	hostState.PortFlagsOverride = types.Bool{Value: hostResponse.PortFlagsOverride}
	hostState.HostFlags.ConsistentLun = types.Bool{Value: hostResponse.ConsistentLun}

	initiators := matchPlanAndResponseInitiators(planInitiators, hostResponse.Initiators)
	saveListAttribute(hostState, initiators, "initiator")
	saveListAttribute(hostState, hostResponse.MaskingviewIDs, "maskingView")
	saveListAttribute(hostState, hostResponse.PowerPathHosts, "powerpathHost")
	setDefaultHostFlags(hostState)
	setHostFlags(hostResponse.EnabledFlags, true, hostState)
	setHostFlags(hostResponse.DisabledFlags, false, hostState)
}

func setHostFlags(flags string, isEnabled bool, hostState *models.Host) {
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

func setDefaultHostFlags(hostState *models.Host) {
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

func saveListAttribute(hostState *models.Host, listAttribute []string, attributeName string) {
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
	if attributeName == "initiator" {
		hostState.Initiators = attributeListType
	} else if attributeName == "maskingView" {
		hostState.Maskingview = attributeListType
	} else if attributeName == "powerpathHost" {
		hostState.PowerpathHosts = attributeListType
	}
}

func matchPlanAndResponseInitiators(planInitiators, responseInitiators []string) []string {
	stateInitiators := []string{}
	respInitiatorsMap := make(map[string]int)
	for _, initiator := range responseInitiators {
		if _, ok := respInitiatorsMap[initiator]; !ok {
			respInitiatorsMap[initiator] = 1
		}
	}
	for _, initiator := range planInitiators {
		if containsIgnoreCase(initiator, responseInitiators) {
			stateInitiators = append(stateInitiators, initiator)
			delete(respInitiatorsMap, strings.ToLower(initiator))
		}
	}

	for initiator := range respInitiatorsMap {
		stateInitiators = append(stateInitiators, initiator)
	}

	return stateInitiators
}

func containsIgnoreCase(elemToFind string, elems []string) bool {
	for _, elem := range elems {
		if strings.EqualFold(elemToFind, elem) {
			return true
		}
	}
	return false
}

func updateHost(ctx context.Context, client client.Client, plan, state models.Host) ([]string, []string, []string) {
	updatedParameters := []string{}
	updateFailedParameters := []string{}
	errorMessages := []string{}

	if plan.Name.Value != state.Name.Value {
		_, err := client.PmaxClient.UpdateHostName(ctx, client.SymmetrixID, state.ID.Value, plan.Name.Value)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "name")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to rename host: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "name")
		}
	}

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

	if !compareStringSlice(planInitiators, stateInitiators) {
		hostResponse, err := client.PmaxClient.GetHostByID(ctx, client.SymmetrixID, plan.Name.Value)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "initiators")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify initiators: %s", "couldn't get the host data"))
		}
		// confirm the lower case logic of initiators
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

	if plan.HostFlags != state.HostFlags {
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
			ConsistentLUN: plan.HostFlags.ConsistentLun.Value,
		}
		_, err := client.PmaxClient.UpdateHostFlags(ctx, client.SymmetrixID, plan.Name.Value, &hostFlags)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "host flags")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify the host flags: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "host flags")
		}
	}
	return updatedParameters, updateFailedParameters, errorMessages
}

func compareStringSlice(planInitiators, stateInitiators []string) bool {
	if len(planInitiators) != len(stateInitiators) {
		return false
	}

	itemAppearsTimes := make(map[string]int, len(planInitiators))
	for _, i := range planInitiators {
		itemAppearsTimes[i]++
	}

	for _, i := range stateInitiators {
		if _, ok := itemAppearsTimes[i]; !ok {
			return false
		}

		itemAppearsTimes[i]--
		if itemAppearsTimes[i] == 0 {
			delete(itemAppearsTimes, i)
		}
	}
	return len(itemAppearsTimes) == 0
}
