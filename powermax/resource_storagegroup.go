package powermax

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type resourceStorageGroupType struct{}

// StorageGroup Resource schema
func (r resourceStorageGroupType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Resource for managing StorageGroups in PowerMax array. Updates are supported for the following parameters: `name`, `srp`, `enable_compression`, `service_level`, `host_io_limits`, `volume_ids`, `snapshot_policies`.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The ID of the storage group.",
				MarkdownDescription: "The ID of the storage group.",
			},
			"name": {
				Type:                types.StringType,
				Required:            true,
				Description:         "The name of the storage group.",
				MarkdownDescription: "The name of the storage group.",
			},
			"srpid": {
				Type:                types.StringType,
				Required:            true,
				Description:         "The ID of the SRP associated with the storage group.",
				MarkdownDescription: "The ID of the SRP associated with the storage group.",
			},
			"service_level": {
				Type:                types.StringType,
				Required:            true,
				Description:         "The service level associated with the storage group. It can be one of Optimized, Diamond, Platinum, Gold, Bronze, Silver or None.",
				MarkdownDescription: "The service level associated with the storage group. It can be one of Optimized, Diamond, Platinum, Gold, Bronze, Silver or None.",
			},
			"enable_compression": {
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				Description:         "Enable compression on the storage group. By default, value is set to true.",
				MarkdownDescription: "Enable compression on the storage group. By default, value is set to true.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.Bool{Value: true}),
				},
			},
			"workload": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The workload associated with the storage group.",
				MarkdownDescription: "The workload associated with the storage group.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"volume_ids": {
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Optional:            true,
				Computed:            true,
				Description:         "The IDs of the volume associated with the storage group. Only pre-existing volumes are considered here.",
				MarkdownDescription: "The IDs of the volume associated with the storage group. Only pre-existing volumes are considered here.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"numofvols": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "The number of volumes associated with the storage group.",
				MarkdownDescription: "The number of volumes associated with the storage group.",
			},
			"numofchildsgs": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "The number of child storage groups associated with the storage group.",
				MarkdownDescription: "The number of child storage groups associated with the storage group.",
			},
			"numofparentsgs": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "The number of parent storage groups associated with the storage group.",
				MarkdownDescription: "The number of parent storage groups associated with the storage group.",
			},
			"numofmaskingviews": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "The number of masking views associated with the storage group.",
				MarkdownDescription: "The number of masking views associated with the storage group.",
			},
			"numofsnapshots": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "The number of snapshots associated with the storage group.",
				MarkdownDescription: "The number of snapshots associated with the storage group.",
			},
			"cap_gb": {
				Type:                types.NumberType,
				Computed:            true,
				Description:         "The capacity of the storage group in GB.",
				MarkdownDescription: "The capacity of the storage group in GB.",
			},
			"slo_compliance": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The service level compliance status of the storage group.",
				MarkdownDescription: "The service level compliance status of the storage group.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"device_emulation": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The emulation of the volumes in the storage group.",
				MarkdownDescription: "The emulation of the volumes in the storage group.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"type": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The storage group type.",
				MarkdownDescription: "The storage group type.",
			},
			"unprotected": {
				Type:                types.BoolType,
				Computed:            true,
				Description:         "This flag states whether the storage group is protected.",
				MarkdownDescription: "This flag states whether the storage group is protected.",
			},
			"compression_ratio": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "Compression ratio of the storage group.",
				MarkdownDescription: "Compression ratio of the storage group.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"compression_ratio_to_one": {
				Type:                types.NumberType,
				Computed:            true,
				Description:         "Compression ratio numeric value of the storage group.",
				MarkdownDescription: "Compression ratio numeric value of the storage group.",
			},
			"vp_saved_percent": {
				Type:                types.NumberType,
				Computed:            true,
				Description:         "VP saved percentage figure.",
				MarkdownDescription: "VP saved percentage figure.",
			},
			"maskingview": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Computed:            true,
				Description:         "The masking views associated with the storage group.",
				MarkdownDescription: "The masking views associated with the storage group.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			// "childstoragegroup": {
			// 	Type: types.ListType{
			// 		ElemType: types.StringType,
			// 	},
			// 	Computed:            true,
			// 	Description:         "The child storage group(s) associated with the storage group.",
			// 	MarkdownDescription: "The child storage group(s) associated with the storage group.",
			// },
			"uuid": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "Storage Group UUID.",
				MarkdownDescription: "Storage Group UUID.",
			},
			"unreducible_data_gb": {
				Type:                types.NumberType,
				Computed:            true,
				Description:         "The amount of unreducible data in GB.",
				MarkdownDescription: "The amount of unreducible data in GB.",
			},
			// "parentstoragegroup": {
			// 	Type: types.ListType{
			// 		ElemType: types.StringType,
			// 	},
			// 	Computed:            true,
			// 	Description:         "The parent storage group(s) associated with the storage group.",
			// 	MarkdownDescription: "The parent storage group(s) associated with the storage group.",
			// },
			"host_io_limits": {
				Computed: true,
				Optional: true,
				Type: types.MapType{
					ElemType: types.StringType,
				},
				Description:         "The host Limits for the storageGroup. Currently the supported host limit parameters are ['host_io_limit_mb_sec','host_io_limit_io_sec','dynamicdistribution'].",
				MarkdownDescription: "The host Limits for the storageGroup. Currently the supported host limit parameters are ['host_io_limit_mb_sec','host_io_limit_io_sec','dynamicdistribution'].",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"snapshot_policies": {
				Computed: true,
				Optional: true,
				Type: types.SetType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"policy_name": types.StringType,
							"is_active":   types.BoolType,
						},
					},
				},
				Description:         "The snapshot policies to be associated with the storageGroup. The 'is_active' field in the nested schema indicates whether the snapshot policy is in resumed/suspended state",
				MarkdownDescription: "The snapshot policies to be associated with the storageGroup. The 'is_active' field in the nested schema indicates whether the snapshot policy is in resumed/suspended state",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
		},
	}, nil
}

// NewResource is a wrapper around provider
func (r resourceStorageGroupType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceStorageGroup{
		p: *(p.(*provider)),
	}, nil
}

type resourceStorageGroup struct {
	p provider
}

// Create StorageGroup
func (r resourceStorageGroup) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Info(ctx, "creating storage group")
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var sgPlan models.StorageGroup
	diags := req.Plan.Get(ctx, &sgPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	planHostLimitDetails := make(map[string]string)
	optionalPayloadParams := make(map[string]interface{})
	diags = sgPlan.HostIOLimits.ElementsAs(ctx, &planHostLimitDetails, true)

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Debug(ctx, "building host limits", map[string]interface{}{
		"planHostLimitDetails": planHostLimitDetails,
	})
	hostLimitPayload := buildHostLimits(planHostLimitDetails)
	if hostLimitPayload != nil {
		optionalPayloadParams["hostLimits"] = hostLimitPayload
	}

	tflog.Debug(ctx, "building snapshot policy", map[string]interface{}{
		"sgPlan": sgPlan,
		"resp":   resp,
	})
	snapshotPolicyPayload := buildSnapshotPolicy(ctx, sgPlan, resp)
	if resp.Diagnostics.HasError() {
		return
	}
	if snapshotPolicyPayload != nil {
		optionalPayloadParams["snapshotPolicies"] = snapshotPolicyPayload
	}

	tflog.Debug(ctx, "building volume IDs", map[string]interface{}{
		"sgPlan": sgPlan,
		"resp":   resp,
	})
	volumeIDs := buildVolumeIDs(ctx, sgPlan, resp)

	symmID := r.p.client.SymmetrixID
	tflog.Debug(ctx, "calling create storage group on pmax client", map[string]interface{}{
		"symmetrixID":           symmID,
		"storageGroup":          sgPlan.Name.Value,
		"SRP ID":                sgPlan.SRPID.Value,
		"service level":         sgPlan.ServiceLevel.Value,
		"optionalPayloadParams": optionalPayloadParams,
	})
	pmaxSgResponse, err := r.p.client.PmaxClient.CreateStorageGroup(ctx, symmID, sgPlan.Name.Value, sgPlan.SRPID.Value, sgPlan.ServiceLevel.Value, false, optionalPayloadParams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating storage group",
			CreateSGDetailErrorMsg+sgPlan.Name.Value+"with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "create storage group response", map[string]interface{}{
		"pmaxSgResponse": pmaxSgResponse,
	})

	if len(volumeIDs) > 0 {
		tflog.Debug(ctx, "calling add volumes to storage group on pmax client", map[string]interface{}{
			"symmID":         symmID,
			"storageGroupID": pmaxSgResponse.StorageGroupID,
			"volumeIDs":      volumeIDs,
		})
		err = r.p.client.PmaxClient.AddVolumesToStorageGroupS(ctx, symmID, pmaxSgResponse.StorageGroupID, false, volumeIDs...)
		if err != nil {
			resp.Diagnostics.AddError(
				"could not add volumes to storageGroup",
				fmt.Sprintf("%s: %s due to %s", CreateSGDetailErrorMsg, sgPlan.Name.Value, err.Error()),
			)
		}
	}
	tflog.Debug(ctx, "calling get volume ID list in  storage group on pmax client", map[string]interface{}{
		"symmID":         symmID,
		"storageGroupID": pmaxSgResponse.StorageGroupID,
	})
	sgVolIDs, err := r.p.client.PmaxClient.GetVolumeIDListInStorageGroup(ctx, symmID, pmaxSgResponse.StorageGroupID)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading volumes of storagegroup",
			err.Error(),
		)
	}
	tflog.Debug(ctx, "get volume ID list in storage group response", map[string]interface{}{
		"sgVolIDs": sgVolIDs,
	})

	sgState := models.StorageGroup{}
	sgResponse := models.StorageGroupWithData{
		PmaxStorageGroup: pmaxSgResponse,
		VolumeIDs:        sgVolIDs,
	}

	for _, snapshotPolicyID := range pmaxSgResponse.SnapshotPolicies {
		tflog.Debug(ctx, "calling get storage group snapshot policy on pmax client", map[string]interface{}{
			"symmID":           symmID,
			"storageGroupID":   pmaxSgResponse.StorageGroupID,
			"snapshotPolicyID": snapshotPolicyID,
		})
		sgSnapshotPolicy, err := r.p.client.PmaxClient.GetStorageGroupSnapshotPolicy(ctx, symmID, snapshotPolicyID, pmaxSgResponse.StorageGroupID)
		if err != nil {
			continue
		}
		tflog.Debug(ctx, "get storage group snapshot policy response", map[string]interface{}{
			"sgSnapshotPolicy": sgSnapshotPolicy,
		})
		sgResponse.SnapshotPolicies = append(sgResponse.SnapshotPolicies, *sgSnapshotPolicy)
	}

	tflog.Debug(ctx, "updating create storage group state", map[string]interface{}{
		"sgState":    sgState,
		"sgPlan":     sgPlan,
		"sgResponse": sgResponse,
	})
	updateState(&sgState, &sgPlan, sgResponse, "create")

	diags = resp.State.Set(ctx, sgState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create storage group completed")
}

// Read StorageGroup
func (r resourceStorageGroup) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	tflog.Info(ctx, "reading storage group")
	var sgState models.StorageGroup
	diags := req.State.Get(ctx, &sgState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sgID := sgState.ID.Value
	tflog.Debug(ctx, "calling get storage group on pmax client", map[string]interface{}{
		"symmetrixID":    r.p.client.SymmetrixID,
		"storageGroupID": sgID,
	})
	pmaxSgResponse, err := r.p.client.PmaxClient.GetStorageGroup(ctx, r.p.client.SymmetrixID, sgID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading storagegroup",
			ReadSGDetailsErrorMsg+sgID+" with error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "calling get volume ID list in storage group on pmax client", map[string]interface{}{
		"symmetrixID":    r.p.client.SymmetrixID,
		"storageGroupID": sgID,
	})
	sgVolIDs, err := r.p.client.PmaxClient.GetVolumeIDListInStorageGroup(ctx, r.p.client.SymmetrixID, pmaxSgResponse.StorageGroupID)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading volumes of storagegroup",
			err.Error(),
		)
	}
	tflog.Debug(ctx, "get volume ID list in storage group response", map[string]interface{}{
		"sgVolIDs": sgVolIDs,
	})

	sgResponse := models.StorageGroupWithData{
		PmaxStorageGroup: pmaxSgResponse,
		VolumeIDs:        sgVolIDs,
	}

	for _, snapshotPolicyID := range pmaxSgResponse.SnapshotPolicies {
		tflog.Debug(ctx, "calling get storage group snapshot policy on pmax client", map[string]interface{}{
			"symmetrixID":      r.p.client.SymmetrixID,
			"storageGroupID":   sgID,
			"snapshotPolicyID": snapshotPolicyID,
		})
		sgSnapshotPolicy, err := r.p.client.PmaxClient.GetStorageGroupSnapshotPolicy(ctx, r.p.client.SymmetrixID, snapshotPolicyID, pmaxSgResponse.StorageGroupID)
		if err != nil {
			continue
		}
		tflog.Debug(ctx, "get storage group snapshot policy response", map[string]interface{}{
			"sgSnapshotPolicy": sgSnapshotPolicy,
		})
		sgResponse.SnapshotPolicies = append(sgResponse.SnapshotPolicies, *sgSnapshotPolicy)
	}

	tflog.Debug(ctx, "updating storage group state", map[string]interface{}{
		"sgState":    sgState,
		"sgResponse": sgResponse,
	})
	updateState(&sgState, nil, sgResponse, "read")

	// Set state
	diags = resp.State.Set(ctx, &sgState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read storage group completed")
}

// Update StorageGroup
// Supported updates: name, service_level, SRP, snapshot policies, compression, volume IDs, host IO limits
func (r resourceStorageGroup) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Info(ctx, "updating storage group")
	var planStorageGroup models.StorageGroup
	diags := req.Plan.Get(ctx, &planStorageGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateStorageGroup models.StorageGroup
	diags = req.State.Get(ctx, &stateStorageGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	storageGroupID := stateStorageGroup.ID.Value
	symmID := r.p.client.SymmetrixID
	tflog.Debug(ctx, "calling update storage group on pmax client", map[string]interface{}{
		"storageGroupID":    storageGroupID,
		"planStorageGroup":  planStorageGroup,
		"stateStorageGroup": stateStorageGroup,
	})
	updatedParameters, updateFailedParameters, errMsgs := UpdateSg(ctx, r.p.client, storageGroupID, planStorageGroup, stateStorageGroup)
	if len(updateFailedParameters) > 0 {
		errorMessage := strings.Join(errMsgs, ",\n ")
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s, The updated parameters are %v and the parameters failed to update are %v", UpdateSGDetailsErrorMsg, updatedParameters, updateFailedParameters),
			errorMessage)
	}
	tflog.Debug(ctx, "update storage group response", map[string]interface{}{
		"updatedParameters":      updatedParameters,
		"updateFailedParameters": updateFailedParameters,
		"errMsgs":                errMsgs,
	})

	tflog.Debug(ctx, "calling get storage group on pmax client", map[string]interface{}{
		"symmetrixID":  symmID,
		"storageGroup": planStorageGroup.Name.Value,
	})

	if isParamUpdated(updatedParameters, "name") {
		storageGroupID = planStorageGroup.Name.Value
	}

	pmaxSgResponse, err := r.p.client.PmaxClient.GetStorageGroup(ctx, symmID, storageGroupID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading storagegroup: %s post update", storageGroupID), err.Error())
		return
	}

	tflog.Debug(ctx, "get storage group response", map[string]interface{}{
		"pmaxSgResponse": pmaxSgResponse,
	})

	tflog.Debug(ctx, "calling get volume ID list in storage group on pmax client", map[string]interface{}{
		"symmetrixID":  symmID,
		"storageGroup": pmaxSgResponse.StorageGroupID,
	})
	sgVolIDs, err := r.p.client.PmaxClient.GetVolumeIDListInStorageGroup(ctx, symmID, pmaxSgResponse.StorageGroupID)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading volumes of storagegroup",
			err.Error(),
		)
	}
	tflog.Debug(ctx, "get volume ID list in storage group response", map[string]interface{}{
		"sgVolIDs": sgVolIDs,
	})

	sgResponse := models.StorageGroupWithData{
		PmaxStorageGroup: pmaxSgResponse,
		VolumeIDs:        sgVolIDs,
	}

	for _, snapshotPolicyID := range pmaxSgResponse.SnapshotPolicies {
		tflog.Debug(ctx, "calling get storage group snapshot policy on pmax client", map[string]interface{}{
			"symmetrixID":      symmID,
			"storageGroup":     pmaxSgResponse.StorageGroupID,
			"snapshotPolicyID": snapshotPolicyID,
		})
		sgSnapshotPolicy, err := r.p.client.PmaxClient.GetStorageGroupSnapshotPolicy(ctx, symmID, snapshotPolicyID, pmaxSgResponse.StorageGroupID)
		if err != nil {
			continue
		}
		tflog.Debug(ctx, "get storage group snapshot policy response", map[string]interface{}{
			"sgSnapshotPolicy": sgSnapshotPolicy,
		})
		sgResponse.SnapshotPolicies = append(sgResponse.SnapshotPolicies, *sgSnapshotPolicy)
	}

	tflog.Debug(ctx, "updating update storage group state", map[string]interface{}{
		"stateStorageGroup": stateStorageGroup,
		"planStorageGroup":  planStorageGroup,
		"sgResponse":        sgResponse,
	})
	updateState(&stateStorageGroup, &planStorageGroup, sgResponse, "update")

	// Set state
	diags = resp.State.Set(ctx, &stateStorageGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update storage group completed")
}

// Delete StorageGroup
func (r resourceStorageGroup) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Info(ctx, "deleting storage group")
	var sgState models.StorageGroup
	diags := req.State.Get(ctx, &sgState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	sgID := sgState.ID.Value
	tflog.Debug(ctx, "calling delete storage group on pmax client", map[string]interface{}{
		"symmetrixID":    r.p.client.SymmetrixID,
		"storageGroupID": sgID,
	})
	err := r.p.client.PmaxClient.DeleteStorageGroup(ctx, r.p.client.SymmetrixID, sgID)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting storagegroup", DeleteSGDetailsErrorMsg+sgID+" with error: "+err.Error())
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "delete storage group completed")
}

// Import resource
func (r resourceStorageGroup) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tflog.Info(ctx, "importing storage group state")
	var sgState models.StorageGroup
	sgState.ID = types.String{Value: req.ID}
	sgID := sgState.ID.Value
	tflog.Debug(ctx, "calling get storage group on pmax client", map[string]interface{}{
		"symmetrixID":    r.p.client.SymmetrixID,
		"storageGroupID": sgID,
	})
	pmaxSgResponse, err := r.p.client.PmaxClient.GetStorageGroup(ctx, r.p.client.SymmetrixID, sgID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing storagegroup",
			ImportSGDetailsErrorMsg+sgID+" with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "get storage group response", map[string]interface{}{
		"pmaxSgResponse": pmaxSgResponse,
	})

	tflog.Debug(ctx, "calling get volume ID list in storage group on pmax client", map[string]interface{}{
		"symmetrixID":    r.p.client.SymmetrixID,
		"storageGroupID": pmaxSgResponse.StorageGroupID,
	})
	sgVolIDs, err := r.p.client.PmaxClient.GetVolumeIDListInStorageGroup(ctx, r.p.client.SymmetrixID, pmaxSgResponse.StorageGroupID)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading volumes of storagegroup",
			err.Error(),
		)
	}
	tflog.Debug(ctx, "get volume ID list in storage group response", map[string]interface{}{
		"sgVolIDs": sgVolIDs,
	})

	sgResponse := models.StorageGroupWithData{
		PmaxStorageGroup: pmaxSgResponse,
		VolumeIDs:        sgVolIDs,
	}

	for _, snapshotPolicyID := range pmaxSgResponse.SnapshotPolicies {
		tflog.Debug(ctx, "calling get storage group snapshot policy on pmax client", map[string]interface{}{
			"symmetrixID":      r.p.client.SymmetrixID,
			"storageGroupID":   pmaxSgResponse.StorageGroupID,
			"snapshotPolicyID": snapshotPolicyID,
		})
		sgSnapshotPolicy, err := r.p.client.PmaxClient.GetStorageGroupSnapshotPolicy(ctx, r.p.client.SymmetrixID, snapshotPolicyID, pmaxSgResponse.StorageGroupID)
		if err != nil {
			continue
		}
		tflog.Debug(ctx, "get storage group snapshot policy response", map[string]interface{}{
			"sgSnapshotPolicy": sgSnapshotPolicy,
		})
		sgResponse.SnapshotPolicies = append(sgResponse.SnapshotPolicies, *sgSnapshotPolicy)
	}

	tflog.Debug(ctx, "updating import storage group state", map[string]interface{}{
		"sgState":    sgState,
		"sgResponse": sgResponse,
	})
	updateState(&sgState, nil, sgResponse, "import")

	// Set state
	diags := resp.State.Set(ctx, &sgState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "import storage group state completed")
}
