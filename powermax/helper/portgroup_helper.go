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
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/models"

	"dell/powermax-go-client"
	pmax "dell/powermax-go-client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// GetPmaxPortsFromTfsdkPG returns a slice of pmaxTypes.PortKey from a models.PortGroup.
func GetPmaxPortsFromTfsdkPG(tfsdkPg models.PortGroup) []pmax.SymmetrixPortKey {

	if len(tfsdkPg.Ports) > 0 {
		pmaxPorts := make([]pmax.SymmetrixPortKey, 0)

		for _, port := range tfsdkPg.Ports {
			portmap := pmax.SymmetrixPortKey{
				DirectorId: port.DirectorID.ValueString(),
				PortId:     port.PortID.ValueString(),
			}
			pmaxPorts = append(pmaxPorts, portmap)
		}

		return pmaxPorts
	}
	return nil
}

// GetPortGroupList get the list of portgroups.
func GetPortGroupList(ctx context.Context, client client.Client, pgPlan models.PortgroupsDataSourceModel) (*powermax.ListPortGroupResult, *http.Response, error) {
	// Read the portgroup based on portgroup type and if nothing is mentioned, then it returns all the port groups
	portGroupsParam := client.PmaxOpenapiClient.SLOProvisioningApi.ListPortGroups(ctx, client.SymmetrixID)
	var typeStr string = ""
	if pgPlan.PgFilter != nil {
		if !pgPlan.PgFilter.Type.IsNull() {
			typeStr = pgPlan.PgFilter.Type.ValueString()
		}
	}
	if typeStr == "iscsi" {
		portGroupsParam = portGroupsParam.Iscsi("true")
	} else { //default Fiber
		portGroupsParam = portGroupsParam.Fibre("true")
	}

	return client.PmaxOpenapiClient.SLOProvisioningApi.ListPortGroupsExecute(portGroupsParam)
}

// CreatePortGroup get the list of portgroups.
func CreatePortGroup(ctx context.Context, client client.Client, plan models.PortGroup) (*powermax.PortGroup, *http.Response, error) {
	pmaxPorts := GetPmaxPortsFromTfsdkPG(plan)
	// Read the portgroup based on portgroup type and if nothing is mentioned, then it returns all the port groups
	//Read the portgroup based on portgroup type and if nothing is mentioned, then it returns all the port groups
	portGroups := client.PmaxOpenapiClient.SLOProvisioningApi.CreatePortGroup(ctx, client.SymmetrixID) //(ctx, d.client.SymmetrixID, state.Type.ValueString())

	createParam := pmax.NewCreatePortGroupParam(plan.Name.ValueString())
	createParam.SetPortGroupProtocol(plan.Protocol.ValueString())
	createParam.SetSymmetrixPortKey(pmaxPorts)

	portGroups = portGroups.CreatePortGroupParam(*createParam)

	tflog.Debug(ctx, "calling create port group on pmax client", map[string]interface{}{
		"symmetrixID": client.SymmetrixID,
		"name":        plan.Name.ValueString(),
		"ports":       pmaxPorts,
	})

	return portGroups.Execute()
}

// UpdatePGState updates the state of a PortGroup.
func UpdatePGState(pgState, pgPlan *models.PortGroup, pgResponse *pmax.PortGroup) {
	pgState.ID = types.StringValue(pgResponse.PortGroupId)
	pgState.Name = types.StringValue(pgResponse.PortGroupId)
	// Always use the state protocol to avoid Response collisions
	pgState.Protocol = pgPlan.Protocol

	if portGroupNoMaskingview, ok := pgResponse.GetNumOfMaskingViewsOk(); ok {
		pgState.NumOfMaskingViews = types.Int64Value(*portGroupNoMaskingview)
	}
	if numOfPorts, ok := pgResponse.GetNumOfPortsOk(); ok {
		pgState.NumOfPorts = types.Int64Value(int64(*numOfPorts))
	}
	if pgtype, ok := pgResponse.GetTypeOk(); ok {
		pgState.Type = types.StringValue(*pgtype)
	}

	var attributeListType types.List
	if portGroupMaskingview, ok := pgResponse.GetMaskingviewOk(); ok {
		var attributeList []attr.Value
		for _, attribute := range portGroupMaskingview {
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
		director := strings.ToUpper(symmetrixPortkey.DirectorId)
		port := strings.ToLower(symmetrixPortkey.PortId)
		key := fmt.Sprintf("%s/%s", director, port)
		pmaxPortKeysToTfsdkPorts[key] = models.PortKey{
			DirectorID: types.StringValue(symmetrixPortkey.DirectorId),
			PortID:     types.StringValue(symmetrixPortkey.PortId),
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

// UpdatePortGroup updates a PortGroup and returns a slice of updated parameters, failed parameters and error messages.
func UpdatePortGroup(ctx context.Context, client client.Client, planPg, statePg models.PortGroup) (updatedParams []string, updateFailedParams []string, errorMessages []string) {
	planPorts := GetPmaxPortsFromTfsdkPG(planPg)
	statePorts := GetPmaxPortsFromTfsdkPG(statePg)
	if !(len(planPorts) == 0 && len(statePorts) == 0) {
		_, err := updatePortGroupParams(ctx, client, statePg.Name.ValueString(), planPorts)
		if err != nil {
			updateFailedParams = append(updateFailedParams, "ports")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to update ports: %s", err.Error()))
		} else {
			updatedParams = append(updatedParams, "ports")
		}

	}

	if planPg.Name.ValueString() != statePg.Name.ValueString() {
		_, err := RenamePortGroup(ctx, client, client.SymmetrixID, statePg.ID.ValueString(), planPg.Name.ValueString())
		if err != nil {
			updateFailedParams = append(updateFailedParams, "name")
			errorMessages = append(errorMessages, fmt.Sprintf("Failed to rename PortGroup: %s", err.Error()))
		} else {
			updatedParams = append(updatedParams, "name")
		}
	}
	return updatedParams, updateFailedParams, errorMessages
}

// UpdatePortGroup - Update the PortGroup based on the 'ports' slice. The slice represents the intended
// configuration of the PortGroup after successful completion of the request.
// based on the passed in 'ports' the implementation will determine how to update
// the PortGroup and make appropriate REST calls sequentially. Take this into
// consideration when making parallel calls.
func updatePortGroupParams(ctx context.Context, client client.Client, portGroupID string, ports []pmax.SymmetrixPortKey) (*pmax.PortGroup, error) {

	// Create map of string "<DIRECTOR ID>/<PORT ID>" to a SymmetrixPortKeyType object based on the passed in 'ports'
	inPorts := make(map[string]*pmax.SymmetrixPortKey)
	for _, port := range ports {
		director := strings.ToUpper(port.DirectorId)
		port := strings.ToLower(port.PortId)
		key := fmt.Sprintf("%s/%s", director, port)
		if inPorts[key] == nil {
			inPorts[key] = &pmax.SymmetrixPortKey{
				DirectorId: director,
				PortId:     port,
			}
		}
	}
	pg, shouldReturn, returnValue := ReadPortgroupByID(ctx, client, portGroupID)
	if shouldReturn {
		return pg, returnValue
	}

	portIDRegex, err := regexp.Compile(`\\w+:(\\d+)`)

	if err != nil {
		return nil, fmt.Errorf("unable to update port group error: %s", err.Error())
	}

	// Create map of string "<DIRECTOR ID>/<PORT ID>" to a SymmetrixPortKeyType object based on what's found
	// in the PortGroup
	pgPorts := make(map[string]*pmax.SymmetrixPortKey)
	for _, p := range pg.SymmetrixPortKey {
		director := strings.ToUpper(p.DirectorId)
		// PortID string may come as a combination of directory + port_number
		// Extract just the port_number part
		port := strings.ToLower(p.PortId)
		submatch := portIDRegex.FindAllStringSubmatch(port, -1)
		if len(submatch) > 0 {
			port = submatch[0][1]
		}
		key := fmt.Sprintf("%s/%s", director, port)
		pgPorts[key] = &pmax.SymmetrixPortKey{
			DirectorId: director,
			PortId:     port,
		}
	}

	// Diff ports in request with ones in PortGroup --> ports to add
	var added []pmax.SymmetrixPortKey
	for k, v := range inPorts {
		if pgPorts[k] == nil {
			added = append(added, *v)
		}
	}

	// Diff ports in PortGroup with ones in request --> ports to remove
	var removed []pmax.SymmetrixPortKey
	for k, v := range pgPorts {
		if inPorts[k] == nil {
			removed = append(removed, *v)
		}
	}

	if len(added) > 0 {
		tflog.Info(ctx, fmt.Sprintf("Adding ports %v", added))

		edit := &pmax.EditPortGroupActionParam{
			AddPortParam: &pmax.AddPortParam{
				Port: added,
			},
		}
		pgResponse, shouldReturn, err1 := modifyPortGroup(ctx, client, portGroupID, *edit)
		if shouldReturn {
			return pgResponse, err1
		}
	}

	if len(removed) > 0 {
		tflog.Info(ctx, fmt.Sprintf("Removing ports %v", removed))
		edit := &pmax.EditPortGroupActionParam{
			RemovePortParam: &pmax.RemovePortParam{
				Port: removed,
			},
		}
		pgResponse, shouldReturn, err1 := modifyPortGroup(ctx, client, portGroupID, *edit)
		if shouldReturn {
			return pgResponse, err1
		}
	}
	return pg, nil
}

func modifyPortGroup(ctx context.Context, client client.Client, portGroupID string, edit pmax.EditPortGroupActionParam) (*pmax.PortGroup, bool, error) {
	modifyParam := client.PmaxOpenapiClient.SLOProvisioningApi.ModifyPortGroup(ctx, client.SymmetrixID, portGroupID)
	editParam := pmax.NewEditPortGroupParam(edit)
	modifyParam = modifyParam.EditPortGroupParam(*editParam)
	pgResponse, resp1, err := client.PmaxOpenapiClient.SLOProvisioningApi.ModifyPortGroupExecute(modifyParam)
	if err != nil {
		return pgResponse, true, err
	}
	if resp1.StatusCode != http.StatusOK {
		err1 := errors.New(
			"Unable to Read PowerMax Port Groups. Got http error - " +
				resp1.Status,
		)
		return pgResponse, true, err1
	}
	tflog.Debug(ctx, "get port group by ID response", map[string]interface{}{
		"pgResponse": pgResponse,
	})
	return pgResponse, false, nil
}

// RenamePortGroup - Renames a port group.
func RenamePortGroup(ctx context.Context, client client.Client, symID string, portGroupID string, newName string) (*pmax.PortGroup, error) {

	edit := pmax.EditPortGroupActionParam{
		RenamePortGroupParam: &pmax.RenamePortGroupParam{
			NewPortGroupName: newName,
		},
	}
	pgResponse, shouldReturn, err1 := modifyPortGroup(ctx, client, portGroupID, edit)
	if shouldReturn {
		return pgResponse, err1
	}
	return pgResponse, nil
}

// ReadPortgroupByID Read PortGroup by ID.
func ReadPortgroupByID(ctx context.Context, client client.Client, portGroupID string) (*pmax.PortGroup, bool, error) {
	portGroups := client.PmaxOpenapiClient.SLOProvisioningApi.GetPortGroup(ctx, client.SymmetrixID, portGroupID)
	pgResponse, resp1, err := client.PmaxOpenapiClient.SLOProvisioningApi.GetPortGroupExecute(portGroups)

	if err != nil {
		return pgResponse, true, err
	}
	if resp1.StatusCode != http.StatusOK {
		err1 := errors.New(
			"Unable to Read PowerMax Port Groups. Got http error - " +
				resp1.Status,
		)
		return pgResponse, true, err1
	}
	tflog.Debug(ctx, "get port group by ID response", map[string]interface{}{
		"pgResponse": pgResponse,
	})
	return pgResponse, false, nil
}
