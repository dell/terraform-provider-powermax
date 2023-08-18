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
	"net/http"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/models"
)

// CreateMaskingView Creates a new masking view
func CreateMaskingView(ctx context.Context, client client.Client, plan models.MaskingViewResourceModel, hostOrHostGroupID string, isHost bool) (*pmax.MaskingView, *http.Response, error) {
	hostOrHostGroupSelection := *pmax.NewHostOrHostGroupSelection()
	if isHost {
		hostOrHostGroupSelection.UseExistingHostParam = pmax.NewUseExistingHostParam(hostOrHostGroupID)

	} else {
		hostOrHostGroupSelection.UseExistingHostGroupParam = pmax.NewUseExistingHostGroupParam(hostOrHostGroupID)
	}

	portGroupSelection := *pmax.NewPortGroupSelection()
	portGroupSelection.UseExistingPortGroupParam = pmax.NewUseExistingPortGroupParam(plan.PortGroupID.ValueString())

	storageGroupSelection := *pmax.NewStorageGroupSelection()
	storageGroupSelection.UseExistingStorageGroupParam = pmax.NewUseExistingStorageGroupParam(plan.StorageGroupID.ValueString())

	createMaskingViewParam := pmax.NewCreateMaskingViewParam(plan.Name.ValueString())
	createMaskingViewParam.SetHostOrHostGroupSelection(hostOrHostGroupSelection)
	createMaskingViewParam.SetPortGroupSelection(portGroupSelection)
	createMaskingViewParam.SetStorageGroupSelection(storageGroupSelection)

	maskingViewReq := client.PmaxOpenapiClient.SLOProvisioningApi.CreateMaskingView(ctx, client.SymmetrixID)
	maskingViewReq = maskingViewReq.CreateMaskingViewParam(*createMaskingViewParam)
	return maskingViewReq.Execute()
}

// GetMaskingView Gets a Masking View
func GetMaskingView(ctx context.Context, client client.Client, name string) (*pmax.MaskingView, *http.Response, error) {
	return client.PmaxOpenapiClient.SLOProvisioningApi.GetMaskingView(ctx, client.SymmetrixID, name).Execute()
}
