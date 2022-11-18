package powermax

import (
	"terraform-provider-powermax/models"

	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
)

func updateMaskingViewState(mvState *models.MaskingView, mvResponse *pmaxTypes.MaskingView) {
	mvState.ID.Value = mvResponse.MaskingViewID
	mvState.Name.Value = mvResponse.MaskingViewID
	mvState.StorageGroupID.Value = mvResponse.StorageGroupID
	mvState.PortGroupID.Value = mvResponse.PortGroupID
	mvState.HostID.Value = mvResponse.HostID
	mvState.HostGroupID.Value = mvResponse.HostGroupID
}
