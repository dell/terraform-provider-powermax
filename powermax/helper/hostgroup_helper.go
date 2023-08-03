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
	"dell/powermax-go-client"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/constants"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// UpdateHostGroupState update host group state.
func UpdateHostGroupState(hostGroupState *models.HostGroupModel, hostGroupResponse *powermax.HostGroup) {
	hostGroupState.ID = types.StringValue(hostGroupResponse.HostGroupId)
	hostGroupState.Name = types.StringValue(hostGroupResponse.HostGroupId)
	hostGroupState.NumOfHosts = types.Int64Value(int64(*hostGroupResponse.NumOfHosts))
	hostGroupState.NumOfInitiators = types.Int64Value(int64(*hostGroupResponse.NumOfInitiators))
	hostGroupState.NumOfMaskingViews = types.Int64Value(*hostGroupResponse.NumOfMaskingViews)
	hostGroupState.Type = types.StringValue(*hostGroupResponse.Type)
	hostGroupState.PortFlagsOverride = types.BoolValue(*hostGroupResponse.PortFlagsOverride)
	hostGroupState.ConsistentLun = types.BoolValue(*hostGroupResponse.ConsistentLun)

	responseHostIDs := []string{}
	for _, hostSummary := range hostGroupResponse.Host {
		responseHostIDs = append(responseHostIDs, hostSummary.HostId)
	}

	saveHgListAttribute(hostGroupState, responseHostIDs, "hostIDs")
	saveHgListAttribute(hostGroupState, hostGroupResponse.Maskingview, "maskingViews")
	setDefaultHostFlagsForHostGroup(hostGroupState)
	if hostGroupResponse.EnabledFlags != nil {
		setHostFlagsInHg(*hostGroupResponse.EnabledFlags, true, hostGroupState)
	}
	if hostGroupResponse.DisabledFlags != nil {
		setHostFlagsInHg(*hostGroupResponse.DisabledFlags, false, hostGroupState)
	}
}

func setHostFlagsInHg(flags string, isEnabled bool, hostState *models.HostGroupModel) {
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

func setDefaultHostFlagsForHostGroup(hostState *models.HostGroupModel) {
	hostState.HostFlags = &models.HostFlags{
		VolumeSetAddressing: models.HostFlag{
			Enabled:  basetypes.NewBoolValue(false),
			Override: basetypes.NewBoolValue(false),
		},
		DisableQResetOnUA: models.HostFlag{
			Enabled:  basetypes.NewBoolValue(false),
			Override: basetypes.NewBoolValue(false),
		},
		EnvironSet: models.HostFlag{
			Enabled:  basetypes.NewBoolValue(false),
			Override: basetypes.NewBoolValue(false),
		},
		AvoidResetBroadcast: models.HostFlag{
			Enabled:  basetypes.NewBoolValue(false),
			Override: basetypes.NewBoolValue(false),
		},
		OpenVMS: models.HostFlag{
			Enabled:  basetypes.NewBoolValue(false),
			Override: basetypes.NewBoolValue(false),
		},
		SCSI3: models.HostFlag{
			Enabled:  basetypes.NewBoolValue(false),
			Override: basetypes.NewBoolValue(false),
		},
		Spc2ProtocolVersion: models.HostFlag{
			Enabled:  basetypes.NewBoolValue(false),
			Override: basetypes.NewBoolValue(false),
		},
		SCSISupport1: models.HostFlag{
			Enabled:  basetypes.NewBoolValue(false),
			Override: basetypes.NewBoolValue(false),
		},
	}
}

func saveHgListAttribute(hostGroupState *models.HostGroupModel, listAttribute []string, attributeName string) {
	var attributeListType types.List

	if len(listAttribute) > 0 {
		var attributeList []attr.Value
		for _, attribute := range listAttribute {
			attributeList = append(attributeList, types.StringValue(attribute))
		}
		attributeListType, _ = types.ListValue(types.StringType, attributeList)
	} else {
		// Empty List
		attributeListType, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if attributeName == "maskingViews" {
		hostGroupState.Maskingviews = attributeListType
	} else if attributeName == "hostIDs" {
		hostGroupState.HostIDs = types.Set(attributeListType)
	}
}

// UpdateHostGroup update host group and return updated parameters, failed updated parameters and errors.
func UpdateHostGroup(ctx context.Context, client client.Client, plan, state models.HostGroupModel) ([]string, []string, []string) {
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

		var add []string
		var remove []string

		// Host to Add
		for _, hostID := range planHostIDs {
			// if this host is not in the list of current hosts, add it
			if !stringInSlice(hostID, stateHostIDs) {
				tflog.Debug(ctx, fmt.Sprintf("Appending hosts to host group : %s", hostID))
				add = append(add, hostID)
			}
		}

		edit := &powermax.EditHostGroupActionParam{
			AddHostParam: &powermax.AddHostParam{
				Host: add,
			},
		}

		_, doReturn, err := ModifyHostGroup(ctx, client, state.ID.ValueString(), *edit)

		if doReturn {
			updateFailedParameters = append(updateFailedParameters, "host_ids")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to Add host_ids: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "host_ids")
		}

		// Hosts to remove
		for _, hostID := range stateHostIDs {
			// if this host is not in the list of current hosts, remove it
			if !stringInSlice(hostID, planHostIDs) {
				tflog.Debug(ctx, fmt.Sprintf("Removing hosts from host group : %s", hostID))
				remove = append(remove, hostID)
			}
		}

		editRemove := &powermax.EditHostGroupActionParam{
			RemoveHostParam: &powermax.RemoveHostParam{
				Host: remove,
			},
		}

		_, doReturnRemove, errRemove := ModifyHostGroup(ctx, client, state.ID.ValueString(), *editRemove)

		if doReturnRemove {
			updateFailedParameters = append(updateFailedParameters, "host_ids")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to Remove host_ids: %s", errRemove.Error()))
		} else {
			updatedParameters = append(updatedParameters, "host_ids")
		}
	}

	if *plan.HostFlags != *state.HostFlags || plan.ConsistentLun.ValueBool() != state.ConsistentLun.ValueBool() {
		flags := powermax.NewHostFlags(
			*powermax.NewVolumeSetAddressing(plan.HostFlags.VolumeSetAddressing.Enabled.ValueBool(), plan.HostFlags.VolumeSetAddressing.Override.ValueBool()),
			*powermax.NewDisableQResetOnUa(plan.HostFlags.DisableQResetOnUA.Enabled.ValueBool(), plan.HostFlags.DisableQResetOnUA.Override.ValueBool()),
			*powermax.NewEnvironSet(plan.HostFlags.EnvironSet.Enabled.ValueBool(), plan.HostFlags.EnvironSet.Override.ValueBool()),
			*powermax.NewAvoidResetBroadcast(plan.HostFlags.AvoidResetBroadcast.Enabled.ValueBool(), plan.HostFlags.AvoidResetBroadcast.Override.ValueBool()),
			*powermax.NewOpenvms(plan.HostFlags.OpenVMS.Enabled.ValueBool(), plan.HostFlags.OpenVMS.Override.ValueBool()),
			*powermax.NewScsi3(plan.HostFlags.SCSI3.Enabled.ValueBool(), plan.HostFlags.SCSI3.Override.ValueBool()),
			*powermax.NewSpc2ProtocolVersion(plan.HostFlags.Spc2ProtocolVersion.Enabled.ValueBool(), plan.HostFlags.Spc2ProtocolVersion.Override.ValueBool()),
			*powermax.NewScsiSupport1(plan.HostFlags.SCSISupport1.Enabled.ValueBool(), plan.HostFlags.SCSISupport1.Override.ValueBool()),
			plan.ConsistentLun.ValueBool(),
		)
		flagsParam := powermax.NewSetHostGroupFlagsParam(*flags)
		edit := &powermax.EditHostGroupActionParam{
			SetHostGroupFlagsParam: flagsParam,
		}
		_, doReturn, err := ModifyHostGroup(ctx, client, state.ID.ValueString(), *edit)

		if doReturn {
			errStr := ""
			message := GetErrorString(err, errStr)
			updateFailedParameters = append(updateFailedParameters, "host_flags")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify the host flags: %s", message))
		} else {
			updatedParameters = append(updatedParameters, "host_flags")
		}
	}

	if plan.Name.ValueString() != state.Name.ValueString() {

		RenameHostGroupParam := powermax.RenameHostGroupParam{
			NewHostGroupName: plan.Name.ValueStringPointer(),
		}

		edit := powermax.EditHostGroupActionParam{
			RenameHostGroupParam: &RenameHostGroupParam,
		}
		_, shouldReturn, err := ModifyHostGroup(ctx, client, state.Name.ValueString(), edit)

		if shouldReturn {
			updateFailedParameters = append(updateFailedParameters, "name")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to rename hostGroup: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "name")
		}
	}

	return updatedParameters, updateFailedParameters, errorMessages
}

// ModifyHostGroup modify host group.
func ModifyHostGroup(ctx context.Context, client client.Client, hostGroupID string, edit powermax.EditHostGroupActionParam) (*powermax.HostGroup, bool, error) {
	modifyParam := client.PmaxOpenapiClient.SLOProvisioningApi.ModifyHostGroup(ctx, client.SymmetrixID, hostGroupID)
	editParam := powermax.NewEditHostGroupParam(edit)
	modifyParam = modifyParam.EditHostGroupParam(*editParam)
	hgResponse, resp1, err := client.PmaxOpenapiClient.SLOProvisioningApi.ModifyHostGroupExecute(modifyParam)
	if err != nil {
		return hgResponse, true, err
	}
	if resp1.StatusCode != http.StatusOK {
		err1 := errors.New(
			"Unable to Read PowerMax Host Groups. Got http error - " +
				resp1.Status,
		)
		return hgResponse, true, err1
	}
	tflog.Debug(ctx, "get host group by ID response", map[string]interface{}{
		"hgResponse": hgResponse,
	})
	return hgResponse, false, nil
}

// FilterHostGroupIds Based on state either use the filtered list of host groups or get all host groups.
func FilterHostGroupIds(ctx context.Context, state *models.HostGroupDataSourceModel, plan *models.HostGroupDataSourceModel, client client.Client) ([]string, error) {
	var hostgroupIds []string
	if plan.HostGroupFilter == nil || len(plan.HostGroupFilter.IDs) == 0 {
		hostGroupReq := client.PmaxOpenapiClient.SLOProvisioningApi.ListHostGroups(ctx, client.SymmetrixID)
		hostGroupResponse, resp1, err := hostGroupReq.Execute()
		if err != nil {
			return hostgroupIds, err
		}
		if resp1.StatusCode != http.StatusOK {
			return hostgroupIds, fmt.Errorf("Unable to read PowerMax Hostgroups - %s", resp1.Status)
		}
		hostgroupIds = hostGroupResponse.GetHostGroupId()
	} else {
		for _, hg := range plan.HostGroupFilter.IDs {
			hostgroupIds = append(hostgroupIds, hg.ValueString())
		}
	}
	return hostgroupIds, nil
}

// HostGroupDetailMapper convert pmaxTypes.HostGroup to models.HostGroupDetailModal.
func HostGroupDetailMapper(hg *powermax.HostGroup) (models.HostGroupDetailModal, diag.Diagnostics) {
	model := models.HostGroupDetailModal{
		HostGroupID:       types.StringValue(hg.HostGroupId),
		Name:              types.StringValue(hg.HostGroupId),
		ConsistentLun:     types.BoolValue(*hg.ConsistentLun),
		PortFlagsOverride: types.BoolValue(*hg.PortFlagsOverride),
		NumOfMaskingViews: types.Int64Value(*hg.NumOfMaskingViews),
		NumOfHosts:        types.Int64Value(int64(*hg.NumOfHosts)),
		NumOfInitiators:   types.Int64Value(int64(*hg.NumOfInitiators)),
		Type:              types.StringValue(*hg.Type),
	}
	var hosts []models.HostGroupHostDetailModal
	var err diag.Diagnostics
	for _, host := range hg.Host {
		var intiators types.List

		tempHost := models.HostGroupHostDetailModal{
			HostID: types.StringValue(host.HostId),
		}
		if len(host.Initiator) > 0 {
			var attributeList []attr.Value
			for _, attribute := range host.Initiator {
				attributeList = append(attributeList, types.StringValue(attribute))
			}
			intiators, err = types.ListValue(types.StringType, attributeList)
			if err != nil {
				return model, err
			}
		} else {
			// Empty List
			intiators, err = types.ListValue(types.StringType, []attr.Value{})
			if err != nil {
				return model, err
			}
		}
		tempHost.Initiator = intiators
		hosts = append(hosts, tempHost)
		model.Host = hosts
	}

	if len(hg.Maskingview) > 0 {
		var attributeList []attr.Value
		for _, attribute := range hg.Maskingview {
			attributeList = append(attributeList, types.StringValue(attribute))
		}
		model.Maskingview, err = types.ListValue(types.StringType, attributeList)
		if err != nil {
			return model, err
		}
	} else {
		// Empty List
		model.Maskingview, err = types.ListValue(types.StringType, []attr.Value{})
		if err != nil {
			return model, err
		}
	}
	return model, err
}

// HandleHostFlag Sets the hostflag state in the plan if there is not one there by default, otherwise use what is in the plan.
func HandleHostFlag(plan models.HostGroupModel) powermax.HostFlags {
	if plan.HostFlags == nil {
		return *powermax.NewHostFlags(
			*powermax.NewVolumeSetAddressing(false, false),
			*powermax.NewDisableQResetOnUa(false, false),
			*powermax.NewEnvironSet(false, false),
			*powermax.NewAvoidResetBroadcast(false, false),
			*powermax.NewOpenvms(false, false),
			*powermax.NewScsi3(false, false),
			*powermax.NewSpc2ProtocolVersion(false, false),
			*powermax.NewScsiSupport1(false, false),
			false,
		)
	}

	return *powermax.NewHostFlags(
		*powermax.NewVolumeSetAddressing(plan.HostFlags.VolumeSetAddressing.Enabled.ValueBool(), plan.HostFlags.VolumeSetAddressing.Override.ValueBool()),
		*powermax.NewDisableQResetOnUa(plan.HostFlags.DisableQResetOnUA.Enabled.ValueBool(), plan.HostFlags.DisableQResetOnUA.Override.ValueBool()),
		*powermax.NewEnvironSet(plan.HostFlags.EnvironSet.Enabled.ValueBool(), plan.HostFlags.EnvironSet.Override.ValueBool()),
		*powermax.NewAvoidResetBroadcast(plan.HostFlags.AvoidResetBroadcast.Enabled.ValueBool(), plan.HostFlags.AvoidResetBroadcast.Override.ValueBool()),
		*powermax.NewOpenvms(plan.HostFlags.OpenVMS.Enabled.ValueBool(), plan.HostFlags.OpenVMS.Override.ValueBool()),
		*powermax.NewScsi3(plan.HostFlags.SCSI3.Enabled.ValueBool(), plan.HostFlags.SCSI3.Override.ValueBool()),
		*powermax.NewSpc2ProtocolVersion(plan.HostFlags.Spc2ProtocolVersion.Enabled.ValueBool(), plan.HostFlags.Spc2ProtocolVersion.Override.ValueBool()),
		*powermax.NewScsiSupport1(plan.HostFlags.SCSISupport1.Enabled.ValueBool(), plan.HostFlags.SCSISupport1.Override.ValueBool()),
		plan.ConsistentLun.ValueBool(),
	)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
