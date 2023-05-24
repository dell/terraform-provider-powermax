// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.

package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/constants"
	"terraform-provider-powermax/powermax/helper"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &PortGroup{}
	_ resource.ResourceWithConfigure   = &PortGroup{}
	_ resource.ResourceWithImportState = &PortGroup{}
)

// NewPortGroup is a helper function to simplify the provider implementation.
func NewPortGroup() resource.Resource {
	return &PortGroup{}
}

// PortGroup defines the resource implementation.
type PortGroup struct {
	client *client.Client
}

// Schema Resource schema.
func (r *PortGroup) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Resource for managing PortGroups in PowerMax array. Updates are supported for the following parameters: `name`, `ports`.",
		Description:         "Resource for managing PortGroups in PowerMax array. Updates are supported for the following parameters: `name`, `ports`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the portgroup.",
				MarkdownDescription: "The ID of the portgroup.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the portgroup.",
				MarkdownDescription: "The name of the portgroup.",
			},
			"ports": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"director_id": schema.StringAttribute{
							Required: true,
						},
						"port_id": schema.StringAttribute{
							Required: true,
						},
					},
				},
				Description:         "The list of ports associated with the portgroup.",
				MarkdownDescription: "The list of ports associated with the portgroup.",
			},
			"protocol": schema.StringAttribute{
				Required:            true,
				Description:         "The portgroup protocol.",
				MarkdownDescription: "The portgroup protocol.",
			},
			"numofports": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of ports associated with the portgroup.",
				MarkdownDescription: "The number of ports associated with the portgroup.",
			},
			"numofmaskingviews": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of masking views associated with the portgroup.",
				MarkdownDescription: "The number of masking views associated with the portgroup.",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				Description:         "The type of the portgroup.",
				MarkdownDescription: "The type of the portgroup.",
			},
			"maskingview": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Description:         "The masking views associated with the portgroup.",
				MarkdownDescription: "The masking views associated with the portgroup.",
			},
		},
	}
}

// Metadata Resource metadata.
func (r *PortGroup) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_portgroup"
}

// Configure PortGroup.
func (r *PortGroup) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pmaxClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = pmaxClient
}

// Create PortGroup.
func (r *PortGroup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Info(ctx, "creating port group")

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
	pmaxPorts := helper.GetPmaxPortsFromTfsdkPG(plan)

	tflog.Debug(ctx, "calling create port group on pmax client", map[string]interface{}{
		"symmetrixID": r.client.SymmetrixID,
		"name":        plan.Name.ValueString(),
		"ports":       pmaxPorts,
	})
	pgResponse, err := r.client.PmaxClient.CreatePortGroup(ctx, r.client.SymmetrixID, plan.Name.ValueString(), pmaxPorts, plan.Protocol.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating port group",
			constants.CreatePGDetailErrorMsg+plan.Name.ValueString()+" with error: "+err.Error(),
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
	helper.UpdatePGState(&pgState, &plan, pgResponse)

	diags = resp.State.Set(ctx, pgState)
	resp.Diagnostics.Append(diags...)
}

// Read PortGroup.
func (r *PortGroup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "reading portgroup")
	var pgState models.PortGroup
	diags := req.State.Get(ctx, &pgState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get portgroup ID from API and then update what is in state from what the API returns
	pgID := pgState.ID.ValueString()
	tflog.Debug(ctx, "getting portgroup by ID", map[string]interface{}{
		"symmetrixID": r.client.SymmetrixID,
		"portGroupID": pgID,
	})
	pgResponse, err := r.client.PmaxClient.GetPortGroupByID(ctx, r.client.SymmetrixID, pgID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading portgroup",
			constants.ReadPGDetailsErrorMsg+pgID+" with error: "+err.Error(),
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
	helper.UpdatePGState(&pgState, &pgState, pgResponse)

	diags = resp.State.Set(ctx, pgState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read portgroup completed")
}

// Update PortGroup
// Supported updates: name, ports.
func (r *PortGroup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	updatedParams, updateFailedParameters, errorMessages := helper.UpdatePortGroup(ctx, *r.client, pgPlan, pgState)
	if len(errorMessages) > 0 || len(updateFailedParameters) > 0 {
		errMessage := strings.Join(errorMessages, ",\n")
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s, updated parameters are %v and parameters failed to update are %v", constants.UpdatePGDetailsErrMsg, updatedParams, updateFailedParameters),
			errMessage)
	}

	portGroupID := pgState.ID.ValueString()

	if helper.IsParamUpdated(updatedParams, "name") {
		portGroupID = pgPlan.Name.ValueString()
	}

	pgResponse, err := r.client.PmaxClient.GetPortGroupByID(ctx, r.client.SymmetrixID, portGroupID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading portgroup",
			constants.UpdatePGDetailsErrMsg+pgPlan.Name.ValueString()+" with error: "+err.Error(),
		)
		return
	}

	helper.UpdatePGState(&pgState, &pgPlan, pgResponse)

	diags = resp.State.Set(ctx, pgState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "update portgroup completed")
}

// Delete PortGroup.
func (r *PortGroup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting portgroup")
	var pgState models.PortGroup
	diags := req.State.Get(ctx, &pgState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	pgID := pgState.ID.ValueString()
	tflog.Debug(ctx, "calling delete port group on pmax client", map[string]interface{}{
		"symmetrixID": r.client.SymmetrixID,
		"portGroupID": pgID,
	})
	err := r.client.PmaxClient.DeletePortGroup(ctx, r.client.SymmetrixID, pgID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting portgroup",
			constants.DeletePGDetailsErrorMsg+pgID+" with error: "+err.Error(),
		)
	}
	tflog.Info(ctx, "delete portgroup completed")
}

// ImportState import resource.
func (r *PortGroup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "importing port group state")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
