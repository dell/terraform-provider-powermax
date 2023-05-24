// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.

package helper

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/models"

	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetPmaxPortsFromTfsdkPG returns a slice of pmaxTypes.PortKey from a models.PortGroup
func GetPmaxPortsFromTfsdkPG(tfsdkPg models.PortGroup) []pmaxTypes.PortKey {

	if len(tfsdkPg.Ports) > 0 {
		pmaxPorts := make([]pmaxTypes.PortKey, 0)

		for _, port := range tfsdkPg.Ports {
			portmap := pmaxTypes.PortKey{
				DirectorID: port.DirectorID.ValueString(),
				PortID:     port.PortID.ValueString(),
			}
			pmaxPorts = append(pmaxPorts, portmap)
		}

		return pmaxPorts
	}
	return nil
}

// UpdatePGState updates the state of a PortGroup
func UpdatePGState(pgState, pgPlan *models.PortGroup, pgResponse *pmaxTypes.PortGroup) {
	pgState.ID = types.StringValue(pgResponse.PortGroupID)
	pgState.Name = types.StringValue(pgResponse.PortGroupID)
	if pgResponse.PortGroupProtocol == "" {
		pgState.Protocol = pgPlan.Protocol
	} else {
		pgState.Protocol = types.StringValue(pgResponse.PortGroupProtocol)
	}

	pgState.NumOfMaskingViews = types.Int64Value(pgResponse.NumberMaskingViews)
	pgState.NumOfPorts = types.Int64Value(pgResponse.NumberPorts)
	pgState.Type = types.StringValue(pgResponse.PortGroupType)

	var attributeListType types.List
	if len(pgResponse.MaskingView) > 0 {
		var attributeList []attr.Value
		for _, attribute := range pgResponse.MaskingView {
			attributeList = append(attributeList, types.StringValue(attribute))
		}
		attributeListType, _ = types.ListValue(
			types.StringType,
			attributeList,
		)

	} else {
		attributeListType, _ = types.ListValue(
			types.StringType,
			[]attr.Value{},
		)
	}

	pgState.Maskingview = attributeListType

	symmetrixPortkeyDetails := make([]models.PortKey, 0)
	pmaxPortKeysToTfsdkPorts := make(map[string]models.PortKey)
	for _, symmetrixPortkey := range pgResponse.SymmetrixPortKey {
		director := strings.ToUpper(symmetrixPortkey.DirectorID)
		port := strings.ToLower(symmetrixPortkey.PortID)
		key := fmt.Sprintf("%s/%s", director, port)
		pmaxPortKeysToTfsdkPorts[key] = models.PortKey{
			DirectorID: types.StringValue(symmetrixPortkey.DirectorID),
			PortID:     types.StringValue(symmetrixPortkey.PortID),
		}
	}

	for _, statePort := range pgPlan.Ports {
		key := fmt.Sprintf("%s/%s", strings.ToUpper(statePort.DirectorID.ValueString()), strings.ToUpper(statePort.PortID.ValueString()))
		if val, ok := pmaxPortKeysToTfsdkPorts[key]; ok {
			symmetrixPortkeyDetails = append(symmetrixPortkeyDetails, val)
			delete(pmaxPortKeysToTfsdkPorts, key)
		}
	}
	// following code adds any ports which have been added outside terraform. This will be used for drift detection.
	if len(pmaxPortKeysToTfsdkPorts) > 0 {
		for _, value := range pmaxPortKeysToTfsdkPorts {
			symmetrixPortkeyDetails = append(symmetrixPortkeyDetails, value)
		}

	}
	pgState.Ports = symmetrixPortkeyDetails

}

// UpdatePortGroup updates a PortGroup and returns a slice of updated parameters, failed parameters and error messages
func UpdatePortGroup(ctx context.Context, client client.Client, planPg, statePg models.PortGroup) (updatedParams []string, updateFailedParams []string, errorMessages []string) {
	planPorts := GetPmaxPortsFromTfsdkPG(planPg)
	statePorts := GetPmaxPortsFromTfsdkPG(statePg)
	if !(len(planPorts) == 0 && len(statePorts) == 0) {
		_, err := client.PmaxClient.UpdatePortGroup(ctx, client.SymmetrixID, statePg.Name.ValueString(), planPorts)
		if err != nil {
			updateFailedParams = append(updateFailedParams, "ports")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to update ports: %s", err.Error()))
		} else {
			updatedParams = append(updatedParams, "ports")
		}

	}

	if planPg.Name.ValueString() != statePg.Name.ValueString() {
		_, err := client.PmaxClient.RenamePortGroup(ctx, client.SymmetrixID, statePg.ID.ValueString(), planPg.Name.ValueString())
		if err != nil {
			updateFailedParams = append(updateFailedParams, "name")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to rename PortGroup: %s", err.Error()))
		} else {
			updatedParams = append(updatedParams, "name")
		}
	}
	return updatedParams, updateFailedParams, errorMessages
}
