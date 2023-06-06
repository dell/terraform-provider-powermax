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

// UpdateHostGroupState update host group state.
func UpdateHostGroupState(hostGroupState *models.HostGroupModel, hostGroupResponse *pmaxTypes.HostGroup) {
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
		hostGroupID := state.ID.ValueString()
		_, err := client.PmaxClient.UpdateHostGroupHosts(ctx, client.SymmetrixID, hostGroupID, planHostIDs)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "host_ids")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to modify host_ids: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "host_ids")
		}
	}

	if *plan.HostFlags != *state.HostFlags || plan.ConsistentLun.ValueBool() != state.ConsistentLun.ValueBool() {
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

// FilterHostGroupIds Based on state either use the filtered list of host groups or get all host groups.
func FilterHostGroupIds(ctx context.Context, state *models.HostGroupDataSourceModel, plan *models.HostGroupDataSourceModel, client client.Client) ([]string, error) {
	var hostgroupIds []string
	if plan.HostGroupFilter == nil || len(plan.HostGroupFilter.IDs) == 0 {
		hostGroupResponse, err := client.PmaxClient.GetHostGroupList(ctx, client.SymmetrixID)
		if err != nil {
			return hostgroupIds, err
		}
		hostgroupIds = hostGroupResponse.HostGroupIDs
	} else {
		for _, hg := range plan.HostGroupFilter.IDs {
			hostgroupIds = append(hostgroupIds, hg.ValueString())
		}
	}
	return hostgroupIds, nil
}

// HostGroupDetailMapper convert pmaxTypes.HostGroup to models.HostGroupDetailModal.
func HostGroupDetailMapper(hg *pmaxTypes.HostGroup) (models.HostGroupDetailModal, diag.Diagnostics) {
	model := models.HostGroupDetailModal{
		HostGroupID:       types.StringValue(hg.HostGroupID),
		Name:              types.StringValue(hg.HostGroupID),
		ConsistentLun:     types.BoolValue(hg.ConsistentLun),
		PortFlagsOverride: types.BoolValue(hg.PortFlagsOverride),
		NumOfMaskingViews: types.Int64Value(hg.NumberMaskingViews),
		NumOfHosts:        types.Int64Value(hg.NumOfHosts),
		NumOfInitiators:   types.Int64Value(hg.NumberInitiators),
		Type:              types.StringValue(hg.HostGroupType),
	}
	var hosts []models.HostGroupHostDetailModal
	var err diag.Diagnostics
	for _, host := range hg.Hosts {
		var intiators types.List

		tempHost := models.HostGroupHostDetailModal{
			HostID: types.StringValue(host.HostID),
		}
		if len(host.Initiators) > 0 {
			var attributeList []attr.Value
			for _, attribute := range host.Initiators {
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

	if len(hg.MaskingviewIDs) > 0 {
		var attributeList []attr.Value
		for _, attribute := range hg.MaskingviewIDs {
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
func HandleHostFlag(plan models.HostGroupModel) pmaxTypes.HostFlags {
	if plan.HostFlags == nil {
		return pmaxTypes.HostFlags{
			VolumeSetAddressing: &pmaxTypes.HostFlag{
				Enabled:  false,
				Override: false,
			},
			DisableQResetOnUA: &pmaxTypes.HostFlag{
				Enabled:  false,
				Override: false,
			},
			EnvironSet: &pmaxTypes.HostFlag{
				Enabled:  false,
				Override: false,
			},
			AvoidResetBroadcast: &pmaxTypes.HostFlag{
				Enabled:  false,
				Override: false,
			},
			OpenVMS: &pmaxTypes.HostFlag{
				Enabled:  false,
				Override: false,
			},
			SCSI3: &pmaxTypes.HostFlag{
				Enabled:  false,
				Override: false,
			},
			Spc2ProtocolVersion: &pmaxTypes.HostFlag{
				Enabled:  false,
				Override: false,
			},
			SCSISupport1: &pmaxTypes.HostFlag{
				Enabled:  false,
				Override: false,
			},
			ConsistentLUN: plan.ConsistentLun.ValueBool(),
		}
	}
	return pmaxTypes.HostFlags{
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
}
