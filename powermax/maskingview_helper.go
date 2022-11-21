package powermax

import (
	"terraform-provider-powermax/models"

	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func updateMaskingViewState(mvState *models.MaskingView, mvResponse *pmaxTypes.MaskingView) {
	mvState.ID = types.String{Value: mvResponse.MaskingViewID}
	mvState.Name = types.String{Value: mvResponse.MaskingViewID}
	mvState.StorageGroupID = types.String{Value: mvResponse.StorageGroupID}
	mvState.PortGroupID = types.String{Value: mvResponse.PortGroupID}
	mvState.HostID = types.String{Value: mvResponse.HostID}
	mvState.HostGroupID = types.String{Value: mvResponse.HostGroupID}
}
