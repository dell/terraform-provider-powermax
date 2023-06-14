/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"errors"
	"fmt"
	"net/http"
	"strings"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/constants"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// UpdateHostState update host state.
func UpdateHostState(hostState *models.HostModel, planInitiators []string, hostResponse *pmax.Host) {
	hostState.HostID = types.StringValue(hostResponse.HostId)
	hostState.Name = types.StringValue(hostResponse.HostId)

	if hostNumMaskingview, ok := hostResponse.GetNumOfMaskingViewsOk(); ok {
		hostState.NumberMaskingViews = types.Int64Value(*hostNumMaskingview)
	}
	if hostNumOfHostGroups, ok := hostResponse.GetNumOfHostGroupsOk(); ok {
		hostState.NumberHostGroups = types.Int64Value(int64(*hostNumOfHostGroups))
	}
	if hostNumOfInitiators, ok := hostResponse.GetNumOfInitiatorsOk(); ok {
		hostState.NumberInitiators = types.Int64Value(int64(*hostNumOfInitiators))
	}
	if hostNumOfPowerPathHosts, ok := hostResponse.GetNumOfPowerpathHostsOk(); ok {
		hostState.NumPowerPathHosts = types.Int64Value(*hostNumOfPowerPathHosts)
	}
	hostBwLimit := hostResponse.GetBwLimit()
	hostState.BWLimit = types.Int64Value(hostBwLimit)

	if hostType, ok := hostResponse.GetTypeOk(); ok {
		hostState.HostType = types.StringValue(*hostType)
	}
	if hostPortFlagsOverride, ok := hostResponse.GetPortFlagsOverrideOk(); ok {
		hostState.PortFlagsOverride = types.BoolValue(*hostPortFlagsOverride)
	}
	if hostConsistentLun, ok := hostResponse.GetConsistentLunOk(); ok {
		hostState.ConsistentLun = types.BoolValue(*hostConsistentLun)
	}

	iniAttributeList := []attr.Value{}
	for _, ini := range hostResponse.Initiator {
		iniAttributeList = append(iniAttributeList, types.StringValue(ini))
	}
	hostState.Initiators, _ = types.ListValue(types.StringType, iniAttributeList)

	maskViewAttributeList := []attr.Value{}
	for _, id := range hostResponse.Maskingview {
		maskViewAttributeList = append(maskViewAttributeList, types.StringValue(id))
	}
	hostState.MaskingviewIDs, _ = types.ListValue(types.StringType, maskViewAttributeList)

	powerPathAttributeList := []attr.Value{}
	for _, id := range hostResponse.Powerpathhosts {
		powerPathAttributeList = append(powerPathAttributeList, types.StringValue(id))
	}
	hostState.PowerPathHosts, _ = types.ListValue(types.StringType, powerPathAttributeList)

	hostGroupList := []attr.Value{}
	for _, hg := range hostResponse.Hostgroup {
		hostGroupList = append(hostGroupList, types.StringValue(hg))
	}
	hostState.HostGroup, _ = types.ListValue(types.StringType, hostGroupList)

	setDefaultHostFlags(hostState)
	if hostEnabledFlags, ok := hostResponse.GetEnabledFlagsOk(); ok {
		setHostFlags(*hostEnabledFlags, true, hostState)
	}
	if hostDisabledFlags, ok := hostResponse.GetDisabledFlagsOk(); ok {
		setHostFlags(*hostDisabledFlags, true, hostState)
	}

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
		getReq := client.PmaxOpenapiClient.SLOProvisioningApi.GetHost(ctx, client.SymmetrixID, state.HostID.ValueString())
		hostResponse, _, err := getReq.Execute()
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "initiators")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify initiators: %s for %s %v", "couldn't get the host data", state.HostID.ValueString(), hostResponse))
		}

		var planInitiatorsLowerCase []string
		for _, planInitiator := range planInitiators {
			planInitiatorsLowerCase = append(planInitiatorsLowerCase, strings.ToLower(planInitiator))
		}
		initRemove := []string{}
		initAdd := []string{}

		// check for initiators that are being added
		for _, init := range planInitiatorsLowerCase {
			// if this initiator is not in the list of current initiators, add it
			if !StringInSlice(init, hostResponse.Initiator) {
				initAdd = append(initAdd, init)
			}
		}

		// check for initiators to be removed
		for _, init := range hostResponse.Initiator {
			if !StringInSlice(init, planInitiatorsLowerCase) {
				initRemove = append(initRemove, init)
			}
		}
		// add initiators if needed
		if len(initAdd) > 0 {
			addInitiatorParam := pmax.NewAddInitiatorParam(initAdd)
			edit := &pmax.EditHostActionParam{
				AddInitiatorParam: addInitiatorParam,
			}
			_, err := ModifyHost(client, ctx, state.HostID.ValueString(), *edit)
			if err != nil {
				err1, ok := err.(*pmax.GenericOpenAPIError)
				message := ""
				if ok {
					message, _ = ParseBody(err1.Body())
				}
				updateFailedParameters = append(updateFailedParameters, "add_initiators")
				errorMessages = append(errorMessages, fmt.Sprintf("Failed to add initiators to host: %s", message))
			} else {
				updatedParameters = append(updatedParameters, "add_initiators")
			}
		}

		// remove initiators if needed
		if len(initRemove) > 0 {
			removeInitiatorParam := pmax.NewRemoveInitiatorParam(initRemove)
			edit := &pmax.EditHostActionParam{
				RemoveInitiatorParam: removeInitiatorParam,
			}
			_, err := ModifyHost(client, ctx, state.HostID.ValueString(), *edit)
			if err != nil {
				err1, ok := err.(*pmax.GenericOpenAPIError)
				message := ""
				if ok {
					message, _ = ParseBody(err1.Body())
				}
				updateFailedParameters = append(updateFailedParameters, "remove_initiators")
				errorMessages = append(errorMessages, fmt.Sprintf("Failed to remove initiators from host: %s", message))
			} else {
				updatedParameters = append(updatedParameters, "remove_initiators")
			}
		}

	}
	// Update host flags
	if plan.HostFlags != state.HostFlags || plan.ConsistentLun.ValueBool() != state.ConsistentLun.ValueBool() {
		hostFlags := pmax.NewHostFlags(
			*pmax.NewVolumeSetAddressing(plan.HostFlags.VolumeSetAddressing.Enabled.ValueBool(), plan.HostFlags.VolumeSetAddressing.Override.ValueBool()),
			*pmax.NewDisableQResetOnUa(plan.HostFlags.DisableQResetOnUA.Enabled.ValueBool(), plan.HostFlags.DisableQResetOnUA.Override.ValueBool()),
			*pmax.NewEnvironSet(plan.HostFlags.EnvironSet.Enabled.ValueBool(), plan.HostFlags.EnvironSet.Override.ValueBool()),
			*pmax.NewAvoidResetBroadcast(plan.HostFlags.AvoidResetBroadcast.Enabled.ValueBool(), plan.HostFlags.AvoidResetBroadcast.Override.ValueBool()),
			*pmax.NewOpenvms(plan.HostFlags.OpenVMS.Enabled.ValueBool(), plan.HostFlags.OpenVMS.Override.ValueBool()),
			*pmax.NewScsi3(plan.HostFlags.SCSI3.Enabled.ValueBool(), plan.HostFlags.SCSI3.Override.ValueBool()),
			*pmax.NewSpc2ProtocolVersion(plan.HostFlags.Spc2ProtocolVersion.Enabled.ValueBool(), plan.HostFlags.Spc2ProtocolVersion.Override.ValueBool()),
			*pmax.NewScsiSupport1(plan.HostFlags.SCSISupport1.Enabled.ValueBool(), plan.HostFlags.SCSISupport1.Override.ValueBool()),
			plan.ConsistentLun.ValueBool(),
		)
		flagsParam := pmax.NewSetHostFlagsParam(*hostFlags)
		edit := &pmax.EditHostActionParam{
			SetHostFlagsParam: flagsParam,
		}
		_, err := ModifyHost(client, ctx, state.HostID.ValueString(), *edit)

		if err != nil {
			err1, ok := err.(*pmax.GenericOpenAPIError)
			message := ""
			if ok {
				message, _ = ParseBody(err1.Body())
			}
			updateFailedParameters = append(updateFailedParameters, "host_flags")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify the host flags: %s", message))
		} else {
			updatedParameters = append(updatedParameters, "host_flags")
		}
	}

	// Update host name
	if plan.Name.ValueString() != state.Name.ValueString() {

		renameHostParam := pmax.RenameHostParam{
			NewHostName: plan.Name.ValueStringPointer(),
		}
		edit := pmax.EditHostActionParam{
			RenameHostParam: &renameHostParam,
		}
		_, err := ModifyHost(client, ctx, state.Name.ValueString(), edit)
		if err != nil {
			err1, ok := err.(*pmax.GenericOpenAPIError)
			message := ""
			if ok {
				message, _ = ParseBody(err1.Body())
			}
			updateFailedParameters = append(updateFailedParameters, "name")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to rename host: %s", message))
		} else {
			updatedParameters = append(updatedParameters, "name")
		}
	}

	return updatedParameters, updateFailedParameters, errorMessages
}

func ModifyHost(client client.Client, ctx context.Context, hostId string, edit pmax.EditHostActionParam) (*pmax.Host, error) {
	modifyParam := client.PmaxOpenapiClient.SLOProvisioningApi.ModifyHost(ctx, client.SymmetrixID, hostId)
	editParam := pmax.NewEditHostParam(edit)
	modifyParam = modifyParam.EditHostParam(*editParam)
	hostResp, resp1, err := client.PmaxOpenapiClient.SLOProvisioningApi.ModifyHostExecute(modifyParam)
	if err != nil {
		return hostResp, err
	}
	if resp1.StatusCode != http.StatusOK {
		err1 := errors.New(
			"Unable to Update PowerMax Host . Got http error - " +
				resp1.Status,
		)
		return hostResp, err1
	}
	tflog.Debug(ctx, "get host  by ID response", map[string]interface{}{
		"hostResponse": hostResp,
	})
	return hostResp, nil
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
