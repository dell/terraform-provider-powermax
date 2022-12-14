package powermax

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/models"

	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func buildHostLimits(planHostLimitDetails map[string]string) *pmaxTypes.SetHostIOLimitsParam {
	if len(planHostLimitDetails) > 0 {
		pmaxHostIOLimitParam := &pmaxTypes.SetHostIOLimitsParam{
			HostIOLimitMBSec:    planHostLimitDetails["host_io_limit_mb_sec"],
			HostIOLimitIOSec:    planHostLimitDetails["host_io_limit_io_sec"],
			DynamicDistribution: planHostLimitDetails["dynamicdistribution"],
		}
		return pmaxHostIOLimitParam
	}
	return nil
}

func buildSnapshotPolicy(ctx context.Context, plan models.StorageGroup, resp *tfsdk.CreateResourceResponse) []string {
	planSnapShotPolicies := &[]types.Object{}
	pmaxSnapShotPolicies := []string{}

	diags := plan.SnapshotPolicies.ElementsAs(ctx, planSnapShotPolicies, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return nil
	}

	for _, planSnapshotPolicy := range *planSnapShotPolicies {
		snapshotpolicy := &models.SnapshotPolicy{}
		diags = planSnapshotPolicy.As(ctx, snapshotpolicy, types.ObjectAsOptions{UnhandledNullAsEmpty: true})
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return nil
		}
		if !snapshotpolicy.IsActive.Value {
			resp.Diagnostics.AddError(
				"Error creating storage group",
				CreateSGDetailErrorMsg+plan.Name.Value+" with error: cannot create storage group with suspended snapshot policy",
			)
			return nil
		}
		pmaxSnapShotPolicies = append(pmaxSnapShotPolicies, snapshotpolicy.PolicyName.Value)
	}

	if len(pmaxSnapShotPolicies) > 0 {
		return pmaxSnapShotPolicies
	}
	return nil
}

func buildVolumeIDs(ctx context.Context, plan models.StorageGroup, resp *tfsdk.CreateResourceResponse) []string {
	volumeIDs := []string{}
	err := plan.VolumeIDs.ElementsAs(ctx, &volumeIDs, true)
	if err.HasError() {
		resp.Diagnostics.Append(err...)
	}
	return volumeIDs
}

func updateState(sgState, sgPlan *models.StorageGroup, sgWithData models.StorageGroupWithData, operation string) {
	sgResponse := sgWithData.PmaxStorageGroup
	sgState.ID.Value = sgResponse.StorageGroupID
	sgState.Name.Value = sgResponse.StorageGroupID

	sgState.SRPID.Value = sgResponse.SRP
	if sgState.SRPID.Value == "" {
		sgState.SRPID.Value = "none"
	}

	if sgPlan != nil {
		sgState.ServiceLevel.Value = sgPlan.ServiceLevel.Value

	} else {
		if !strings.EqualFold(sgState.ServiceLevel.Value, sgResponse.ServiceLevel) {
			sgState.ServiceLevel.Value = sgResponse.ServiceLevel
		}
	}
	if sgState.ServiceLevel.Value == "" {
		sgState.ServiceLevel.Value = "none"
	}

	sgState.CapGB = types.Number{Value: big.NewFloat(sgResponse.CapacityGB)}
	sgState.CompressionRatio = types.String{Value: sgResponse.CompressionRatio}
	sgState.CompressionRatioToOne = types.Number{Value: big.NewFloat(sgResponse.CompressionRatioToOne)}
	sgState.DeviceEmulation.Value = sgResponse.DeviceEmulation
	sgState.EnableCompression.Value = sgResponse.Compression
	sgState.NumOfChildSgs.Value = int64(sgResponse.NumOfChildSGs)
	sgState.NumOfMaskingViews.Value = int64(sgResponse.NumOfMaskingViews)
	sgState.NumOfParentSgs.Value = int64(sgResponse.NumOfParentSGs)
	sgState.NumOfSnapshots.Value = int64(sgResponse.NumOfSnapshots)
	sgState.NumOfVols.Value = int64(sgResponse.NumOfVolumes)
	sgState.Type.Value = sgResponse.Type
	sgState.UnreducibleDataGB = types.Number{Value: big.NewFloat(sgResponse.UnreducibleDataGB)}
	sgState.VPsavedPercent = types.Number{Value: big.NewFloat(sgResponse.VPSavedPercent)}
	sgState.Workload.Value = sgResponse.Workload
	sgState.Unprotected.Value = sgResponse.Unprotected
	sgState.UUID.Value = sgResponse.UUID
	sgState.HostIOLimits.Elems = make(map[string]attr.Value)
	sgState.HostIOLimits.ElemType = types.StringType
	if sgResponse.HostIOLimit != nil {
		sgState.HostIOLimits.Elems["host_io_limit_mb_sec"] = types.String{Value: sgResponse.HostIOLimit.HostIOLimitMBSec}
		sgState.HostIOLimits.Elems["host_io_limit_io_sec"] = types.String{Value: sgResponse.HostIOLimit.HostIOLimitIOSec}
		sgState.HostIOLimits.Elems["dynamicdistribution"] = types.String{Value: sgResponse.HostIOLimit.DynamicDistribution}
	}

	maskingViews := types.List{
		ElemType: types.StringType,
	}
	tfsdkmaskingViews := []attr.Value{}
	for _, maskingView := range sgWithData.PmaxStorageGroup.MaskingView {
		tfsdkmaskingViews = append(tfsdkmaskingViews, types.String{Value: maskingView})
	}
	maskingViews.Elems = tfsdkmaskingViews
	sgState.MaskingView = maskingViews

	snapshotPoliciesTfsdk := types.Set{
		ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"policy_name": types.StringType,
				"is_active":   types.BoolType,
			},
		},
	}
	snapshotPolicyObjects := []attr.Value{}

	if len(sgResponse.SnapshotPolicies) > 0 {
		if operation == "read" {
			snapshotPoliciesMap := make(map[string]bool)
			for _, stateSnapshotPolicy := range sgState.SnapshotPolicies.Elems {
				for _, snapshotPolicyResp := range sgWithData.SnapshotPolicies {
					snapshotPolicyDetails := make(map[string]attr.Value)
					stateSnapshotPolicyObject := stateSnapshotPolicy.(types.Object)
					if (stateSnapshotPolicyObject.Attrs["policy_name"] == types.String{Value: snapshotPolicyResp.SnapshotPolicyID}) {
						snapshotPolicyDetails["policy_name"] = types.String{Value: snapshotPolicyResp.SnapshotPolicyID}
						snapshotPolicyDetails["is_active"] = types.Bool{Value: !snapshotPolicyResp.Suspended}
						snapshotPolicyObject := types.Object{
							Attrs: snapshotPolicyDetails,
							AttrTypes: map[string]attr.Type{
								"policy_name": types.StringType,
								"is_active":   types.BoolType,
							},
						}
						snapshotPolicyObjects = append(snapshotPolicyObjects, snapshotPolicyObject)
						snapshotPoliciesMap[snapshotPolicyResp.SnapshotPolicyID] = true
						break
					}
				}
			}
			// Loop to identify the policies added outside terraform
			for _, snapshotPolicyResp := range sgWithData.SnapshotPolicies {
				snapshotPolicyDetails := make(map[string]attr.Value)
				if _, ok := snapshotPoliciesMap[snapshotPolicyResp.SnapshotPolicyID]; !ok {
					snapshotPolicyDetails["policy_name"] = types.String{Value: snapshotPolicyResp.SnapshotPolicyID}
					snapshotPolicyDetails["is_active"] = types.Bool{Value: !snapshotPolicyResp.Suspended}
					snapshotPolicyObject := types.Object{
						Attrs: snapshotPolicyDetails,
						AttrTypes: map[string]attr.Type{
							"policy_name": types.StringType,
							"is_active":   types.BoolType,
						},
					}
					snapshotPolicyObjects = append(snapshotPolicyObjects, snapshotPolicyObject)
					snapshotPoliciesMap[snapshotPolicyResp.SnapshotPolicyID] = true
				}
			}

		} else if operation == "import" {
			for _, snapshotPolicyResp := range sgWithData.SnapshotPolicies {
				snapshotPolicyDetails := make(map[string]attr.Value)
				snapshotPolicyDetails["policy_name"] = types.String{Value: snapshotPolicyResp.SnapshotPolicyID}
				snapshotPolicyDetails["is_active"] = types.Bool{Value: !snapshotPolicyResp.Suspended}
				snapshotPolicyObject := types.Object{
					Attrs: snapshotPolicyDetails,
					AttrTypes: map[string]attr.Type{
						"policy_name": types.StringType,
						"is_active":   types.BoolType,
					},
				}
				snapshotPolicyObjects = append(snapshotPolicyObjects, snapshotPolicyObject)
			}
		} else {
			for _, planSnapshotPolicy := range sgPlan.SnapshotPolicies.Elems {
				for _, snapshotPolicyResp := range sgWithData.SnapshotPolicies {
					snapshotPolicyDetails := make(map[string]attr.Value)
					planSnapshotPolicyObject := planSnapshotPolicy.(types.Object)
					if (planSnapshotPolicyObject.Attrs["policy_name"] == types.String{Value: snapshotPolicyResp.SnapshotPolicyID}) {
						snapshotPolicyDetails["policy_name"] = types.String{Value: snapshotPolicyResp.SnapshotPolicyID}
						snapshotPolicyDetails["is_active"] = types.Bool{Value: !snapshotPolicyResp.Suspended}
						snapshotPolicyObject := types.Object{
							Attrs: snapshotPolicyDetails,
							AttrTypes: map[string]attr.Type{
								"policy_name": types.StringType,
								"is_active":   types.BoolType,
							},
						}
						snapshotPolicyObjects = append(snapshotPolicyObjects, snapshotPolicyObject)
						break
					}
				}
			}
		}
	}
	snapshotPoliciesTfsdk.Elems = snapshotPolicyObjects
	sgState.SnapshotPolicies = snapshotPoliciesTfsdk

	sgVolumes := types.Set{
		ElemType: types.StringType,
	}
	tfsdkSgVols := []attr.Value{}
	if operation == "read" {
		stateVols := []string{}
		diags := sgState.VolumeIDs.ElementsAs(context.Background(), &stateVols, true)
		if !diags.HasError() {
			for _, vol := range stateVols {
				if containsString(sgWithData.VolumeIDs, vol) {
					tfsdkSgVols = append(tfsdkSgVols, types.String{Value: vol})
				}
			}
			// This loop is to add any new volumes added outside terraform to detect drift
			for _, vol := range sgWithData.VolumeIDs {
				if !containsString(stateVols, vol) {
					tfsdkSgVols = append(tfsdkSgVols, types.String{Value: vol})
				}
			}
		}
	} else if operation == "import" {
		for _, vol := range sgWithData.VolumeIDs {
			tfsdkSgVols = append(tfsdkSgVols, types.String{Value: vol})
		}
	} else {
		planVols := []string{}
		diags := sgPlan.VolumeIDs.ElementsAs(context.Background(), &planVols, true)
		if !diags.HasError() {
			for _, vol := range planVols {
				if containsString(sgWithData.VolumeIDs, vol) {
					tfsdkSgVols = append(tfsdkSgVols, types.String{Value: vol})
				}
			}
		}
	}

	sgVolumes.Elems = tfsdkSgVols
	sgState.VolumeIDs = sgVolumes
}

// UpdateSg updates the storageGroup with changed parameters
func UpdateSg(ctx context.Context, client client.Client, sgID string, planSg, stateSg models.StorageGroup) ([]string, []string, []string) {
	updatedParameters := []string{}
	updateFailedParameters := []string{}
	updateErrorMsgs := []string{}

	if planSg.SRPID.Value != stateSg.SRPID.Value {
		err := updateSRP(ctx, client, sgID, planSg.SRPID.Value)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "srp")
			updateErrorMsgs = append(updateErrorMsgs, fmt.Sprintf("Failed to modify srp_id: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "srp")
		}
	}

	if planSg.EnableCompression.Value != stateSg.EnableCompression.Value {
		err := updateSgCompression(ctx, client, sgID, planSg.EnableCompression.Value)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "enable_compression")
			updateErrorMsgs = append(updateErrorMsgs, fmt.Sprintf("Failed to modify compression: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "enable_compression")
		}

	}

	if !strings.EqualFold(planSg.ServiceLevel.Value, stateSg.ServiceLevel.Value) {
		err := updateServiceLevel(ctx, client, sgID, planSg.ServiceLevel.Value)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "service_level")
			updateErrorMsgs = append(updateErrorMsgs, fmt.Sprintf("Failed to modify service_level: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "service_level")
		}
	}

	if !(len(planSg.HostIOLimits.Elems) == 0 && len(stateSg.HostIOLimits.Elems) == 0) && !reflect.DeepEqual(planSg.HostIOLimits.Elems, stateSg.HostIOLimits.Elems) {
		err := updateHostIOLimits(ctx, client, sgID, planSg, stateSg)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "host_io_limits")
			updateErrorMsgs = append(updateErrorMsgs, fmt.Sprintf("Failed to modify host_io_limits: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "host_io_limits")
		}

	}

	if !(len(planSg.VolumeIDs.Elems) == 0 && len(stateSg.VolumeIDs.Elems) == 0) && (!reflect.DeepEqual(planSg.VolumeIDs.Elems, stateSg.VolumeIDs.Elems)) {
		err := updateVolumes(ctx, client, sgID, planSg, stateSg)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "volume_ids")
			updateErrorMsgs = append(updateErrorMsgs, fmt.Sprintf("Failed to modify volume_ids: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "volume_ids")
		}

	}

	if !reflect.DeepEqual(planSg.SnapshotPolicies.Elems, stateSg.SnapshotPolicies.Elems) {
		err := updateSnapshotPolicies(ctx, client, sgID, planSg, stateSg)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "snapshot_policies")
			updateErrorMsgs = append(updateErrorMsgs, fmt.Sprintf("Failed to modify snapshot_policies: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "snapshot_policies")
		}

	}

	if planSg.Name.Value != stateSg.Name.Value {
		err := updateSgName(ctx, client, sgID, planSg.Name.Value)
		if err != nil {
			updateFailedParameters = append(updateFailedParameters, "name")
			updateErrorMsgs = append(updateErrorMsgs, fmt.Sprintf("Failed to rename storageGroup: %s", err.Error()))
		} else {
			updatedParameters = append(updatedParameters, "name")
		}
	}
	return updatedParameters, updateFailedParameters, updateErrorMsgs
}

func updateSgName(ctx context.Context, client client.Client, sgID string, name string) error {
	payload := pmaxTypes.UpdateStorageGroupPayload{
		ExecutionOption: pmaxTypes.ExecutionOptionSynchronous,
	}
	payload.EditStorageGroupActionParam.RenameStorageGroupParam = &pmaxTypes.RenameStorageGroupParam{
		NewStorageGroupName: name,
	}
	err := client.PmaxClient.UpdateStorageGroupS(ctx, client.SymmetrixID, sgID, payload)
	if err != nil {
		return err
	}
	return nil
}

func updateSgCompression(ctx context.Context, client client.Client, sgID string, compression bool) error {
	payload := pmaxTypes.UpdateStorageGroupPayload{
		ExecutionOption: pmaxTypes.ExecutionOptionSynchronous,
	}
	payload.EditStorageGroupActionParam.EditCompressionParam = &pmaxTypes.EditCompressionParam{
		Compression: &compression,
	}
	err := client.PmaxClient.UpdateStorageGroupS(ctx, client.SymmetrixID, sgID, &payload)
	if err != nil {
		return err
	}
	return nil
}

func updateServiceLevel(ctx context.Context, client client.Client, sgID, serviceLevel string) error {
	payload := pmaxTypes.UpdateStorageGroupPayload{
		ExecutionOption: pmaxTypes.ExecutionOptionSynchronous,
	}
	payload.EditStorageGroupActionParam.EditStorageGroupSLOParam = &pmaxTypes.EditStorageGroupSLOParam{
		SLOID: serviceLevel,
	}
	err := client.PmaxClient.UpdateStorageGroupS(ctx, client.SymmetrixID, sgID, payload)
	if err != nil {
		return err
	}
	return nil
}

func updateSRP(ctx context.Context, client client.Client, sgID, srp string) error {
	payload := pmaxTypes.UpdateStorageGroupPayload{
		ExecutionOption: pmaxTypes.ExecutionOptionSynchronous,
	}
	payload.EditStorageGroupActionParam.EditStorageGroupSRPParam = &pmaxTypes.EditStorageGroupSRPParam{
		SRPID: srp,
	}
	err := client.PmaxClient.UpdateStorageGroupS(ctx, client.SymmetrixID, sgID, payload)
	if err != nil {
		return err
	}
	return nil
}

func updateHostIOLimits(ctx context.Context, client client.Client, sgID string, planSg, stateSg models.StorageGroup) error {
	payload := pmaxTypes.UpdateStorageGroupPayload{
		ExecutionOption: pmaxTypes.ExecutionOptionSynchronous,
	}
	planHostIOlimimtsMap := make(map[string]string)
	planHostIOlimits := pmaxTypes.SetHostIOLimitsParam{}
	diags := planSg.HostIOLimits.ElementsAs(ctx, &planHostIOlimimtsMap, true)
	if diags.HasError() {
		return fmt.Errorf("unable to parse hostIOLimits from plan")
	}
	hostIolimitBytes, marshaleErr := json.Marshal(planHostIOlimimtsMap)
	if marshaleErr != nil {
		return fmt.Errorf("unable to parse hostIOLimits from plan %s", marshaleErr.Error())
	}
	unmarshalErr := json.Unmarshal(hostIolimitBytes, &planHostIOlimits)
	if unmarshalErr != nil {
		return fmt.Errorf("unable to parse hostIOLimits from plan %s", unmarshalErr.Error())
	}

	stateHostIOlimitsMap := make(map[string]string)
	stateHostIOlimits := pmaxTypes.SetHostIOLimitsParam{}
	diags = stateSg.HostIOLimits.ElementsAs(ctx, &stateHostIOlimitsMap, true)
	if diags.HasError() {
		return fmt.Errorf("unable to parse hostIOLimits from state")
	}
	statehostIolimitBytes, marshaleErr := json.Marshal(stateHostIOlimitsMap)
	if marshaleErr != nil {
		return fmt.Errorf("unable to parse hostIOLimits from plan %s", marshaleErr.Error())
	}
	unmarshalErr = json.Unmarshal(statehostIolimitBytes, &stateHostIOlimits)
	if unmarshalErr != nil {
		return fmt.Errorf("unable to parse hostIOLimits from plan %s", unmarshalErr.Error())
	}

	setHostIOlimitsParam := pmaxTypes.SetHostIOLimitsParam{}
	if planHostIOlimits.DynamicDistribution != stateHostIOlimits.DynamicDistribution {
		setHostIOlimitsParam.DynamicDistribution = planHostIOlimits.DynamicDistribution
	}
	if planHostIOlimits.HostIOLimitIOSec != stateHostIOlimits.HostIOLimitIOSec {
		setHostIOlimitsParam.HostIOLimitIOSec = planHostIOlimits.HostIOLimitIOSec
	}
	if planHostIOlimits.HostIOLimitMBSec != stateHostIOlimits.HostIOLimitMBSec {
		setHostIOlimitsParam.HostIOLimitMBSec = planHostIOlimits.HostIOLimitMBSec
	}
	payload.EditStorageGroupActionParam.SetHostIOLimitsParam = &setHostIOlimitsParam
	apiErr := client.PmaxClient.UpdateStorageGroupS(ctx, client.SymmetrixID, sgID, payload)
	if apiErr != nil {
		return apiErr
	}
	return nil
}

func updateVolumes(ctx context.Context, client client.Client, sgID string, planSg, stateSg models.StorageGroup) error {

	var inVolumes, outVolumes []string
	planVolumes := []string{}
	diags := planSg.VolumeIDs.ElementsAs(ctx, &planVolumes, true)
	if diags.HasError() {
		return fmt.Errorf("unable to parse volumeIDs from plan")
	}
	stateVolumes := []string{}
	diags = stateSg.VolumeIDs.ElementsAs(ctx, &stateVolumes, true)
	if diags.HasError() {
		return fmt.Errorf("unable to parse volumeIDs from state")
	}
	for _, volume := range stateVolumes {
		if !containsString(planVolumes, volume) {
			outVolumes = append(outVolumes, volume)
		}

	}
	for _, volume := range planVolumes {
		if !containsString(stateVolumes, volume) {
			inVolumes = append(inVolumes, volume)
		}
	}
	symmID := client.SymmetrixID
	if len(outVolumes) > 0 {
		_, err := client.PmaxClient.RemoveVolumesFromStorageGroup(ctx, symmID, sgID, false, outVolumes...)
		if err != nil {
			return err
		}
	}
	if len(inVolumes) > 0 {
		err := client.PmaxClient.AddVolumesToStorageGroupS(ctx, symmID, sgID, false, inVolumes...)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateSnapshotPolicies(ctx context.Context, client client.Client, sgID string, planSg, stateSg models.StorageGroup) error {
	var inSnapshotPolicies, outSnapshotPolicies, resumeSnapshotPolicies, suspendSnapshotPolicies []string
	planSnapshotPolcies := []models.SnapshotPolicy{}
	diags := planSg.SnapshotPolicies.ElementsAs(ctx, &planSnapshotPolcies, true)
	if diags.HasError() {
		return fmt.Errorf("unable to parse snapshot polciies from plan")
	}
	stateSnapshotPolicies := []models.SnapshotPolicy{}
	diags = stateSg.SnapshotPolicies.ElementsAs(ctx, &stateSnapshotPolicies, true)
	if diags.HasError() {
		return fmt.Errorf("unable to parse snapshot polciies from state")
	}
	statetSnapshotPolicyNames := []string{}
	for _, statesnapshotPolicy := range stateSnapshotPolicies {
		statetSnapshotPolicyNames = append(statetSnapshotPolicyNames, statesnapshotPolicy.PolicyName.Value)
	}

	planSnapshotPolicyNames := []string{}
	for _, planSnapshotPolicy := range planSnapshotPolcies {
		planSnapshotPolicyNames = append(planSnapshotPolicyNames, planSnapshotPolicy.PolicyName.Value)
	}

	for _, planSnapshotPolicy := range planSnapshotPolcies {
		for _, stateSnapshotPolicy := range stateSnapshotPolicies {
			if planSnapshotPolicy.PolicyName.Value == stateSnapshotPolicy.PolicyName.Value {
				if !stateSnapshotPolicy.IsActive.Value && planSnapshotPolicy.IsActive.Value {
					resumeSnapshotPolicies = append(resumeSnapshotPolicies, planSnapshotPolicy.PolicyName.Value)
				} else if stateSnapshotPolicy.IsActive.Value && !planSnapshotPolicy.IsActive.Value {
					suspendSnapshotPolicies = append(suspendSnapshotPolicies, planSnapshotPolicy.PolicyName.Value)
				}
			}
		}
	}

	for _, snapshotPolicy := range planSnapshotPolcies {
		if !containsString(statetSnapshotPolicyNames, snapshotPolicy.PolicyName.Value) {
			if snapshotPolicy.IsActive.Value {
				inSnapshotPolicies = append(inSnapshotPolicies, snapshotPolicy.PolicyName.Value)
			} else {
				suspendSnapshotPolicies = append(suspendSnapshotPolicies, snapshotPolicy.PolicyName.Value)
			}

		}
	}

	for _, snapshotPolicy := range stateSnapshotPolicies {
		if !containsString(planSnapshotPolicyNames, snapshotPolicy.PolicyName.Value) {
			outSnapshotPolicies = append(outSnapshotPolicies, snapshotPolicy.PolicyName.Value)
		}
	}

	payload := pmaxTypes.UpdateStorageGroupPayload{
		ExecutionOption:             pmaxTypes.ExecutionOptionSynchronous,
		EditStorageGroupActionParam: pmaxTypes.EditStorageGroupActionParam{},
	}

	if len(suspendSnapshotPolicies) > 0 {
		payload.EditStorageGroupActionParam.EditSnapshotPoliciesParam = &pmaxTypes.EditSnapshotPoliciesParam{
			SuspendSnapshotPolicyParam: &pmaxTypes.SnapshotPolicies{
				SnapshotPolicies: suspendSnapshotPolicies,
			},
		}
		apiErr := client.PmaxClient.UpdateStorageGroupS(ctx, client.SymmetrixID, sgID, payload)
		if apiErr != nil {
			return apiErr
		}
	}

	if len(resumeSnapshotPolicies) > 0 {
		payload.EditStorageGroupActionParam.EditSnapshotPoliciesParam = &pmaxTypes.EditSnapshotPoliciesParam{
			ResumeSnapshotPolicyParam: &pmaxTypes.SnapshotPolicies{
				SnapshotPolicies: resumeSnapshotPolicies,
			},
		}
		apiErr := client.PmaxClient.UpdateStorageGroupS(ctx, client.SymmetrixID, sgID, payload)
		if apiErr != nil {
			return apiErr
		}
	}

	if len(outSnapshotPolicies) > 0 {
		payload.EditStorageGroupActionParam.EditSnapshotPoliciesParam = &pmaxTypes.EditSnapshotPoliciesParam{
			DisassociateSnapshotPolicyParam: &pmaxTypes.SnapshotPolicies{
				SnapshotPolicies: outSnapshotPolicies,
			},
		}
		apiErr := client.PmaxClient.UpdateStorageGroupS(ctx, client.SymmetrixID, sgID, payload)
		if apiErr != nil {
			return apiErr
		}
	}
	if len(inSnapshotPolicies) > 0 {
		payload.EditStorageGroupActionParam.EditSnapshotPoliciesParam = &pmaxTypes.EditSnapshotPoliciesParam{
			AssociateSnapshotPolicyParam: &pmaxTypes.SnapshotPolicies{
				SnapshotPolicies: inSnapshotPolicies,
			},
		}
		apiErr := client.PmaxClient.UpdateStorageGroupS(ctx, client.SymmetrixID, sgID, payload)
		if apiErr != nil {
			return apiErr
		}
	}
	return nil
}

func containsString(checkList []string, checkElement string) bool {
	for _, elem := range checkList {
		if elem == checkElement {
			return true
		}
	}
	return false
}

func isParamUpdated(updatedParams []string, paramName string) bool {
	isParamUpdate := false
	for _, updatedParam := range updatedParams {
		if updatedParam == paramName {
			isParamUpdate = true
			break
		}
	}
	return isParamUpdate
}
