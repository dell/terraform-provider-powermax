package powermax

import (
	"strings"
	"terraform-provider-powermax/models"

	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func updateMaskingViewState(mvState *models.MaskingView, mvResponse *pmaxTypes.MaskingView, mvPlan *models.MaskingView, operation string) {
	mvState.ID = types.String{Value: mvResponse.MaskingViewID}
	if operation == "read" || operation == "import" {
		if !strings.EqualFold(mvState.HostID.Value, mvResponse.HostID) {
			mvState.HostID = types.String{Value: mvResponse.HostID}
		}

		if !strings.EqualFold(mvState.StorageGroupID.Value, mvResponse.StorageGroupID) {
			mvState.StorageGroupID = types.String{Value: mvResponse.StorageGroupID}
		}

		if !strings.EqualFold(mvState.PortGroupID.Value, mvResponse.PortGroupID) {
			mvState.PortGroupID = types.String{Value: mvResponse.PortGroupID}
		}

		if !strings.EqualFold(mvState.HostGroupID.Value, mvResponse.HostGroupID) {
			mvState.HostGroupID = types.String{Value: mvResponse.HostGroupID}
		}

		mvState.Name = types.String{Value: mvResponse.MaskingViewID}
	} else {
		mvState.Name = mvPlan.Name
		mvState.StorageGroupID = mvPlan.StorageGroupID
		mvState.PortGroupID = mvPlan.PortGroupID
		mvState.HostID = mvPlan.HostID
		mvState.HostGroupID = mvPlan.HostGroupID
	}
}
