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
	"dell/powermax-go-client"
	"fmt"
	"strings"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FilterPortIds Based on state either use the filtered list of ports or get all ports.
func FilterPortIds(ctx context.Context, state *models.PortDataSourceModel, plan *models.PortDataSourceModel, client client.Client) ([]powermax.SymmetrixPortKey, error) {
	var portIds []powermax.SymmetrixPortKey
	// Return all Ports available
	if plan.PortFilter == nil || len(plan.PortFilter.IDs) == 0 {
		getParam := client.PmaxOpenapiClient.SLOProvisioningApi.GetAllPorts(ctx, client.SymmetrixID)
		portResponse, _, err := getParam.Execute()
		if err != nil {
			return portIds, err
		}
		portIds = portResponse.GetSymmetrixPortKey()
	} else {
		state.PortFilter = plan.PortFilter
		// Loop through the list of port ids filter and create the symmetricxPortKey Objects
		for _, port := range plan.PortFilter.IDs {
			split := strings.Split(port.ValueString(), ":")
			if len(split) != 2 {
				return portIds, fmt.Errorf("invalid format for port filter, should be 'directorId:portId'")
			}
			portIds = append(portIds, powermax.SymmetrixPortKey{
				DirectorId: split[0],
				PortId:     split[1],
			})

		}
	}
	return portIds, nil
}

// PortDetailMapper maps the openApi port object to the port terraform model.
func PortDetailMapper(ctx context.Context, port *powermax.DirectorPort) (models.PortDetailModal, error) {
	model := models.PortDetailModal{}
	err := CopyFields(ctx, port.SymmetrixPort, &model)
	model.DirectorID = types.StringValue(port.SymmetrixPort.SymmetrixPortKey.DirectorId)
	model.PortID = types.StringValue(port.SymmetrixPort.SymmetrixPortKey.PortId)
	if networkID, ok := port.SymmetrixPort.GetNetworkIdOk(); ok {
		model.NetworkID = types.Int64Value(*networkID)
	}
	if tpcID, ok := port.SymmetrixPort.GetTcpPortOk(); ok {
		tpc := int64(*tpcID)
		model.TPCPort = types.Int64Value(tpc)
	}
	ipAttributeList := []attr.Value{}
	for _, ini := range port.GetSymmetrixPort().IpAddresses {
		ipAttributeList = append(ipAttributeList, types.StringValue(ini))
	}
	model.IPAddresses, _ = types.ListValue(types.StringType, ipAttributeList)

	if err != nil {
		return model, err
	}
	return model, nil
}
