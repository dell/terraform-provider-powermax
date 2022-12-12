package powermax

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type resourcePortGroupType struct{}

// PortGroup Resource schema
func (r resourcePortGroupType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Resource to manage PortGroups in PowerMax array. Updates are supported for the following parameters: `name`, `ports`.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The ID of the portgroup.",
				MarkdownDescription: "The ID of the portgroup.",
			},
			"name": {
				Type:                types.StringType,
				Required:            true,
				Description:         "The name of the portgroup.",
				MarkdownDescription: "The name of the portgroup.",
			},
			"ports": {
				Required: true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"director_id": {
						Type:     types.StringType,
						Required: true,
					},
					"port_id": {
						Type:     types.StringType,
						Required: true,
					},
				}),
				Description:         "The ports associated with the portgroup.",
				MarkdownDescription: "The ports associated with the portgroup.",
			},
			"protocol": {
				Type:                types.StringType,
				Required:            true,
				Description:         "The portgroup protocol.",
				MarkdownDescription: "The portgroup protocol.",
			},
			"numofports": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "The number of ports associated with the portgroup.",
				MarkdownDescription: "The number of ports associated with the portgroup.",
			},
			"numofmaskingviews": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "The number of masking views associated with the portgroup.",
				MarkdownDescription: "The number of masking views associated with the portgroup.",
			},
			"type": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The type of the portgroup.",
				MarkdownDescription: "The type of the portgroup.",
			},
			"maskingview": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Computed:            true,
				Description:         "The masking views associated with the portgroup.",
				MarkdownDescription: "The masking views associated with the portgroup.",
			},
			"test_id": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The test ID of the portgroup.",
				MarkdownDescription: "The test ID of the portgroup.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
		},
	}, nil
}

// NewResource is a wrapper around provider
func (r resourcePortGroupType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourcePortGroup{
		p: *(p.(*provider)),
	}, nil
}

type resourcePortGroup struct {
	p provider
}

// Create PortGroup
func (r resourcePortGroup) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Info(ctx, "creating port group")
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan models.PortGroup
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "building ports", map[string]interface{}{
		"plan": plan,
		"resp": resp,
	})
	pmaxPorts := getPmaxPortsFromTfsdkPG(plan)

	tflog.Debug(ctx, "calling create port group on pmax client", map[string]interface{}{
		"symmetrixID": r.p.client.SymmetrixID,
		"name":        plan.Name.Value,
		"ports":       pmaxPorts,
	})
	pgResponse, err := r.p.client.PmaxClient.CreatePortGroup(ctx, r.p.client.SymmetrixID, plan.Name.Value, pmaxPorts, plan.Protocol.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating port group",
			CreatePGDetailErrorMsg+plan.Name.Value+" with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "create port group response", map[string]interface{}{
		"pgResponse": pgResponse,
	})

	pgState := models.PortGroup{}
	tflog.Debug(ctx, "updating port group state", map[string]interface{}{
		"pgResponse": pgResponse,
		"pgState":    pgState,
	})
	updatePGState(&pgState, &plan, pgResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating portgroup to terraform state, this can lead to a mismatch between infra and state ",
			CreatePGDetailErrorMsg+plan.Name.Value+" with error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, pgState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create portgroup completed")
}

// Read PortGroup
func (r resourcePortGroup) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	tflog.Info(ctx, "reading portgroup")
	var pgState models.PortGroup
	diags := req.State.Get(ctx, &pgState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get portgroup ID from API and then update what is in state from what the API returns
	pgID := pgState.ID.Value
	tflog.Debug(ctx, "getting portgroup by ID", map[string]interface{}{
		"symmetrixID": r.p.client.SymmetrixID,
		"portGroupID": pgID,
	})
	pgResponse, err := r.p.client.PmaxClient.GetPortGroupByID(ctx, r.p.client.SymmetrixID, pgID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading portgroup",
			ReadPGDetailsErrorMsg+pgID+" with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "get port group by ID response", map[string]interface{}{
		"pgResponse": pgResponse,
	})

	tflog.Debug(ctx, "updating portgroup state", map[string]interface{}{
		"pgState":    pgState,
		"pgResponse": pgResponse,
	})
	updatePGState(&pgState, &pgState, pgResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating portgroup to terraform state, this can lead to a mismatch between infra and state ",
			ReadPGDetailsErrorMsg+pgState.Name.Value+" with error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, pgState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read portgroup completed")
}

// Update PortGroup
// Supported updates: name, ports
func (r resourcePortGroup) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Info(ctx, "updating portgroup")
	var pgPlan, pgState models.PortGroup
	diags := req.State.Get(ctx, &pgState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = req.Plan.Get(ctx, &pgPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedParams, updateFailedParameters, errorMessages := updatePortGroup(ctx, r.p.client, pgPlan, pgState)
	if len(errorMessages) > 0 || len(updateFailedParameters) > 0 {
		errMessage := strings.Join(errorMessages, ",\n")
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s, updated parameters are %v and parameters failed to update are %v", UpdatePGDetailsErrMsg, updatedParams, updateFailedParameters),
			errMessage)
	}

	portGroupID := pgState.ID.Value

	if isParamUpdated(updatedParams, "name") {
		portGroupID = pgPlan.Name.Value
	}

	pgResponse, err := r.p.client.PmaxClient.GetPortGroupByID(ctx, r.p.client.SymmetrixID, portGroupID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading portgroup",
			UpdatePGDetailsErrMsg+pgPlan.Name.Value+" with error: "+err.Error(),
		)
		return
	}

	updatePGState(&pgState, &pgPlan, pgResponse)

	diags = resp.State.Set(ctx, pgState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "update portgroup completed")
}

// Delete PortGroup
func (r resourcePortGroup) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Info(ctx, "deleting portgroup")
	var pgState models.PortGroup
	diags := req.State.Get(ctx, &pgState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	pgID := pgState.ID.Value
	tflog.Debug(ctx, "calling delete port group on pmax client", map[string]interface{}{
		"symmetrixID": r.p.client.SymmetrixID,
		"portGroupID": pgID,
	})
	err := r.p.client.PmaxClient.DeletePortGroup(ctx, r.p.client.SymmetrixID, pgID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting portgroup",
			DeletePGDetailsErrorMsg+pgID+" with error: "+err.Error(),
		)
	}
	tflog.Info(ctx, "delete portgroup completed")
}

// Import resource
func (r resourcePortGroup) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tflog.Info(ctx, "importing port group state")
	var pgState models.PortGroup
	pgState.ID = types.String{Value: req.ID}

	// Get portgroup ID from API and then update what is in state from what the API returns
	pgID := pgState.ID.Value
	tflog.Debug(ctx, "getting portgroup by ID", map[string]interface{}{
		"symmetrixID": r.p.client.SymmetrixID,
		"portGroupID": pgID,
	})
	pgResponse, err := r.p.client.PmaxClient.GetPortGroupByID(ctx, r.p.client.SymmetrixID, pgID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing portgroup",
			ImportPGDetailsErrorMsg+pgID+" with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "get port group by ID response", map[string]interface{}{
		"pgResponse": pgResponse,
	})

	tflog.Debug(ctx, "updating portgroup state", map[string]interface{}{
		"pgState":    pgState,
		"pgResponse": pgResponse,
	})
	updatePGState(&pgState, &pgState, pgResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing portgroup to terraform state, this can lead to a mismatch between infra and state ",
			ImportPGDetailsErrorMsg+pgState.Name.Value+" with error: "+err.Error(),
		)
		return
	}

	diags := resp.State.Set(ctx, pgState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "import port group state completed")
}
