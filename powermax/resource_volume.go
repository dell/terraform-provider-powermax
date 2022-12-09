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

type resourceVolumeType struct{}

// Volume Resource schema
func (r resourceVolumeType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Resource to manage Volumes in PowerMax Array. Updates are supported for the following parameters: `name`, `enable_mobility_id`, `size`, `cap_unit`.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The ID of the volume.",
				MarkdownDescription: "The ID of the volume.",
			},
			"name": {
				Type:                types.StringType,
				Required:            true,
				Description:         "The name of the volume.",
				MarkdownDescription: "The name of the volume.",
			},
			"size": {
				Type:                types.NumberType,
				Required:            true,
				Description:         "The size of the volume.",
				MarkdownDescription: "The size of the volume.",
			},
			"cap_unit": {
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "Volume capacity unit.",
				MarkdownDescription: "Volume capacity unit.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.String{Value: "GB"}),
				},
				Validators: []tfsdk.AttributeValidator{
					validCapUnitValidator{},
				},
			},
			"sg_name": {
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "The name of the storage group. sg_name is required while creating the volume.",
				MarkdownDescription: "The name of the storage group. sg_name is required while creating the volume.",
			},
			"type": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The type of the volume.",
				MarkdownDescription: "The type of the volume.",
			},
			"emulation": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The emulation of the volume Enumeration values.",
				MarkdownDescription: "The emulation of the volume Enumeration values.",
			},
			"ssid": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The ssid of the volume.",
				MarkdownDescription: "The ssid of the volume.",
			},
			"allocated_percent": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "The allocated percentage of the volume.",
				MarkdownDescription: "The allocated percentage of the volume.",
			},
			"status": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The status of the volume.",
				MarkdownDescription: "The status of the volume.",
			},
			"reserved": {
				Type:                types.BoolType,
				Computed:            true,
				Description:         "States whether the volume is reserved.",
				MarkdownDescription: "States whether the volume is reserved.",
			},
			"pinned": {
				Type:                types.BoolType,
				Computed:            true,
				Description:         "States whether the volume is pinned.",
				MarkdownDescription: "States whether the volume is pinned.",
			},
			"wwn": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The WWN of the volume.",
				MarkdownDescription: "The WWN of the volume.",
			},
			"encapsulated": {
				Type:                types.BoolType,
				Computed:            true,
				Description:         "States whether the volume is encapsulated.",
				MarkdownDescription: "States whether the volume is encapsulated.",
			},
			"num_of_storage_groups": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "The number of storage groups associated with the volume.",
				MarkdownDescription: "The number of storage groups associated with the volume.",
			},
			"num_of_front_end_paths": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "The number of front end paths of the volume.",
				MarkdownDescription: "The number of front end paths of the volume.",
			},
			"storagegroup_ids": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Computed:            true,
				Description:         "List of storage groups which are associated with the volume.",
				MarkdownDescription: "List of storage groups which are associated with the volume.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"symmetrix_port_keys": {
				Computed: true,
				Type: types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"director_id": types.StringType,
							"port_id":     types.StringType,
						},
					},
				},
				Description:         "The symmetrix ports associated with the volume.",
				MarkdownDescription: "The symmetrix ports associated with the volume.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"rdf_group_ids": {
				Computed: true,
				Type: types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"rdf_group_number": types.Int64Type,
							"label":            types.StringType,
						},
					},
				},
				Description:         "The RDF groups associated with the volume.",
				MarkdownDescription: "The RDF groups associated with the volume.",
			},
			"snap_source": {
				Type:                types.BoolType,
				Computed:            true,
				Description:         "States whether the volume is a snapvx source.",
				MarkdownDescription: "States whether the volume is a snapvx source.",
			},
			"snap_target": {
				Type:                types.BoolType,
				Computed:            true,
				Description:         "States whether the volume is a snapvx target.",
				MarkdownDescription: "States whether the volume is a snapvx target.",
			},
			"has_effective_wwn": {
				Type:                types.BoolType,
				Computed:            true,
				Description:         "States whether volume has effective WWN.",
				MarkdownDescription: "States whether volume has effective WWN.",
			},
			"effective_wwn": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "Effective WWN of the volume.",
				MarkdownDescription: "Effective WWN of the volume.",
			},
			"encapsulated_wwn": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "Encapsulated  WWN of the volume.",
				MarkdownDescription: "Encapsulated  WWN of the volume.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"oracle_instance_name": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "Oracle instance name associated with the volume.",
				MarkdownDescription: "Oracle instance name associated with the volume.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"enable_mobility_id": {
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				Description:         "States whether mobility ID is enabled on the volume.",
				MarkdownDescription: "States whether mobility ID is enabled on the volume.",
			},
			"unreducible_data_gb": {
				Type:                types.Float64Type,
				Computed:            true,
				Description:         "The amount of unreducible data in Gb.",
				MarkdownDescription: "The amount of unreducible data in Gb.",
			},
			"nguid": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The NGUID of the volume.",
				MarkdownDescription: "The NGUID of the volume.",
			},
		},
	}, nil
}

// NewResource is a wrapper around provider
func (r resourceVolumeType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceVolume{
		p: *(p.(*provider)),
	}, nil
}

type resourceVolume struct {
	p provider
}

// Create Volume
func (r resourceVolume) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Info(ctx, "creating volume")
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan models.Volume
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.SGName.Value == "" {
		resp.Diagnostics.AddError(
			"Error creating volume",
			fmt.Sprintf(CreateVolDetailErrorMsg+" %s with error: %s", plan.Name.Value, "storage group name cannot be empty"),
		)
		return
	}

	size, err := getVolumeSize(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating volume",
			fmt.Sprintf(CreateVolDetailErrorMsg+" %s with error: %s", plan.Name.Value, err.Error()),
		)
		return
	}
	volumeOptions := make(map[string]interface{})
	volumeOptions["capacityUnit"] = plan.CapUnit.Value
	volumeOptions["enableMobility"] = plan.EnableMobilityID.Value
	tflog.Debug(ctx, "calling create volume in storage groups on pmax client", map[string]interface{}{
		"symmetrixID":      r.p.client.SymmetrixID,
		"storageGroupName": plan.SGName.Value,
		"name":             plan.Name.Value,
		"size":             size,
		"volumeOptions":    volumeOptions,
	})
	volResponse, err := r.p.client.PmaxClient.CreateVolumeInStorageGroupS(ctx, r.p.client.SymmetrixID, plan.SGName.Value, plan.Name.Value, size, volumeOptions)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating volume",
			CreateVolDetailErrorMsg+plan.Name.Value+"with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "create volume in storage groups response", map[string]interface{}{
		"volResponse": volResponse,
	})

	volState := models.Volume{}

	tflog.Debug(ctx, "updating create volume state", map[string]interface{}{
		"volResponse": volResponse,
		"plan":        plan,
		"volState":    volState,
	})
	updateVolState(&volState, volResponse, &plan, "create")

	diags = resp.State.Set(ctx, volState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create volume completed")
}

// Read Volume
func (r resourceVolume) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	tflog.Info(ctx, "reading volume")
	var volState models.Volume
	diags := req.State.Get(ctx, &volState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	volID := volState.ID.Value
	tflog.Debug(ctx, "calling get volume by ID", map[string]interface{}{
		"symmetrixID": r.p.client.SymmetrixID,
		"volumeID":    volID,
	})
	volResponse, err := r.p.client.PmaxClient.GetVolumeByID(ctx, r.p.client.SymmetrixID, volID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading volume",
			ReadVolDetailsErrorMsg+volID+" with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "get volume by ID response", map[string]interface{}{
		"volResponse": volResponse,
	})

	tflog.Debug(ctx, "updating read volume state", map[string]interface{}{
		"volResponse": volResponse,
		"volState":    volState,
	})
	updateVolState(&volState, volResponse, nil, "read")
	diags = resp.State.Set(ctx, volState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read volume completed")
}

// Update Volume
// Supported updates: name, mobilityID, size
func (r resourceVolume) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Info(ctx, "updating volume")
	var planVol models.Volume
	diags := req.Plan.Get(ctx, &planVol)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Fetched vol from plan")
	var stateVol models.Volume
	diags = req.State.Get(ctx, &stateVol)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update volume on pmax client", map[string]interface{}{
		"planVol":  planVol,
		"stateVol": stateVol,
	})
	updatedParams, updateFailedParameters, errMessages := updateVol(ctx, r.p.client, planVol, stateVol)
	if len(errMessages) > 0 || len(updateFailedParameters) > 0 {
		errMessage := strings.Join(errMessages, ",\n")
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to update all parameters of Volume, updated parameters are %v and parameters failed to update are %v", updatedParams, updateFailedParameters),
			errMessage)
	}

	volID := stateVol.ID.Value
	tflog.Debug(ctx, "calling get volume by ID on pmax client", map[string]interface{}{
		"symmetrixID": r.p.client.SymmetrixID,
		"volumeID":    volID,
	})
	volResponse, err := r.p.client.PmaxClient.GetVolumeByID(ctx, r.p.client.SymmetrixID, volID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading volume",
			ReadVolDetailsErrorMsg+volID+" with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "get volume by ID response", map[string]interface{}{
		"volResponse": volResponse,
	})

	tflog.Debug(ctx, "updating volume state", map[string]interface{}{
		"volResponse": volResponse,
		"planVol":     planVol,
	})
	updateVolState(&stateVol, volResponse, &planVol, "update")
	diags = resp.State.Set(ctx, stateVol)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update volume completed")
}

// Delete Volume
func (r resourceVolume) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Info(ctx, "deleting volume")
	var volumeState models.Volume
	diags := req.State.Get(ctx, &volumeState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	volumeID := volumeState.ID.Value
	sgAssociatedWithVolume := []string{}
	diags = volumeState.StorageGroupIDs.ElementsAs(ctx, &sgAssociatedWithVolume, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	for _, sgID := range sgAssociatedWithVolume {
		tflog.Debug(ctx, "calling get storage group on pmax client", map[string]interface{}{
			"symmetrixID":    r.p.client.SymmetrixID,
			"storageGroupID": sgID,
		})
		sg, _ := r.p.client.PmaxClient.GetStorageGroup(ctx, r.p.client.SymmetrixID, sgID)
		tflog.Debug(ctx, "get storage group response", map[string]interface{}{
			"sg": sg,
		})
		if sg != nil {
			tflog.Debug(ctx, "calling remove volumes from storage group on pmax client", map[string]interface{}{
				"symmetrixID":    r.p.client.SymmetrixID,
				"storageGroupID": sgID,
				"volumeID":       volumeID,
			})
			_, err := r.p.client.PmaxClient.RemoveVolumesFromStorageGroup(ctx, r.p.client.SymmetrixID, sgID, true, volumeID)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error removing volume from storage group",
					RemoveVolumeFromSGDetailsErrorMsg+"Volume ID: "+volumeID+", SG ID: "+volumeState.SGName.Value+" with error: "+err.Error(),
				)
				return
			}
		}
	}

	tflog.Debug(ctx, "calling delete volume on pmax client", map[string]interface{}{
		"symmetrixID": r.p.client.SymmetrixID,
		"volumeID":    volumeID,
	})
	err := r.p.client.PmaxClient.DeleteVolume(ctx, r.p.client.SymmetrixID, volumeID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting volume",
			DeleteVolumeDetailsErrorMsg+volumeID+" with error: "+err.Error(),
		)
	}
	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "delete volume completed")
}

// Import resource
func (r resourceVolume) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tflog.Info(ctx, "importing volume state")
	var stateVol models.Volume
	volID := req.ID
	tflog.Debug(ctx, "calling get volume by ID on pmax client", map[string]interface{}{
		"symmetrixID": r.p.client.SymmetrixID,
		"volumeID":    volID,
	})
	volResponse, err := r.p.client.PmaxClient.GetVolumeByID(ctx, r.p.client.SymmetrixID, volID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing volume",
			ImportVolDetailsErrorMsg+volID+" with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "get volume by ID response", map[string]interface{}{
		"volResponse": volResponse,
	})
	tflog.Debug(ctx, "updating import volume state", map[string]interface{}{
		"stateVol":    stateVol,
		"volResponse": volResponse,
	})
	updateVolState(&stateVol, volResponse, nil, "import")

	diags := resp.State.Set(ctx, stateVol)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "import volume state completed")
}
