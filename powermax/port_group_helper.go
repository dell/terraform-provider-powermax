package powermax

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/models"

	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getPmaxPortsFromTfsdkPG(tfsdkPg models.PortGroup) []pmaxTypes.PortKey {

	if len(tfsdkPg.Ports) > 0 {
		pmaxPorts := make([]pmaxTypes.PortKey, 0)

		for _, port := range tfsdkPg.Ports {
			portmap := pmaxTypes.PortKey{
				DirectorID: port.DirectorID.Value,
				PortID:     port.PortID.Value,
			}
			pmaxPorts = append(pmaxPorts, portmap)
		}

		return pmaxPorts
	}
	return nil
}

func updatePGState(pgState, pgPlan *models.PortGroup, pgResponse *pmaxTypes.PortGroup) {
	pgState.ID.Value = pgResponse.PortGroupID
	pgState.Name.Value = pgResponse.PortGroupID
	pgState.Protocol.Value = pgResponse.PortGroupProtocol
	pgState.NumOfMaskingViews.Value = pgResponse.NumberMaskingViews
	pgState.NumOfPorts.Value = pgResponse.NumberPorts
	pgState.TestID.Value = pgResponse.TestID
	pgState.Type.Value = pgResponse.PortGroupType

	var attributeListType types.List
	if len(pgResponse.MaskingView) > 0 {
		var attributeList []attr.Value
		for _, attribute := range pgResponse.MaskingView {
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

	pgState.Maskingview = attributeListType

	symmetrixPortkeyDetails := make([]models.PortKey, 0)
	pmaxPortKeysToTfsdkPorts := make(map[string]models.PortKey)
	for _, symmetrixPortkey := range pgResponse.SymmetrixPortKey {
		director := strings.ToUpper(symmetrixPortkey.DirectorID)
		port := strings.ToLower(symmetrixPortkey.PortID)
		key := fmt.Sprintf("%s/%s", director, port)
		pmaxPortKeysToTfsdkPorts[key] = models.PortKey{
			DirectorID: types.String{Value: symmetrixPortkey.DirectorID},
			PortID:     types.String{Value: symmetrixPortkey.PortID},
		}
	}

	for _, statePort := range pgPlan.Ports {
		key := fmt.Sprintf("%s/%s", strings.ToUpper(statePort.DirectorID.Value), strings.ToUpper(statePort.PortID.Value))
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

func updatePortGroup(ctx context.Context, client client.Client, planPg, statePg models.PortGroup) (updatedParams []string, updateFailedParams []string, errorMessages []string) {
	if planPg.Name.Value != statePg.Name.Value {
		_, err := client.PmaxClient.RenamePortGroup(ctx, client.SymmetrixID, statePg.ID.Value, planPg.Name.Value)
		if err != nil {
			updateFailedParams = append(updateFailedParams, "name")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to rename PortGroup: %s", err.Error()))
		} else {
			updatedParams = append(updatedParams, "name")
		}
	}
	planPorts := getPmaxPortsFromTfsdkPG(planPg)
	statePorts := getPmaxPortsFromTfsdkPG(statePg)
	if !(len(planPorts) == 0 && len(statePorts) == 0) && !reflect.DeepEqual(planPorts, statePorts) {
		_, err := client.PmaxClient.UpdatePortGroup(ctx, client.SymmetrixID, planPg.Name.Value, planPorts)
		if err != nil {
			updateFailedParams = append(updateFailedParams, "ports")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to update ports: %s", err.Error()))
		} else {
			updatedParams = append(updatedParams, "ports")
		}

	}
	return updatedParams, updateFailedParams, errorMessages
}
