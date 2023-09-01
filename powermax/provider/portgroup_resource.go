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

package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/constants"
	"terraform-provider-powermax/powermax/helper"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
		MarkdownDescription: "Resource for managing PortGroups in PowerMax array. Supported Update (name, ports). PowerMax port groups contain director and port identification and belong to a masking view. Ports can be added to and removed from the port group. Port groups that are no longer associated with a masking view can be deleted. Note the following recommendations: Port groups should contain four or more ports. Each port in a port group should be on a different director. A port can belong to more than one port group. However, for storage systems running HYPERMAX OS 5977 or higher, you cannot mix different types of ports (physical FC ports, virtual ports, and iSCSI virtual ports) within a single port group",
		Description:         "Resource for managing PortGroups in PowerMax array. Supported Update (name, ports). PowerMax port groups contain director and port identification and belong to a masking view. Ports can be added to and removed from the port group. Port groups that are no longer associated with a masking view can be deleted. Note the following recommendations: Port groups should contain four or more ports. Each port in a port group should be on a different director. A port can belong to more than one port group. However, for storage systems running HYPERMAX OS 5977 or higher, you cannot mix different types of ports (physical FC ports, virtual ports, and iSCSI virtual ports) within a single port group",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the portgroup.",
				MarkdownDescription: "The ID of the portgroup.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the portgroup. Only alphanumeric characters, underscores ( _ ), and hyphens (-) are allowed. (Update Supported)",
				MarkdownDescription: "The name of the portgroup. Only alphanumeric characters, underscores ( _ ), and hyphens (-) are allowed. (Update Supported)",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(64),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9_-]*$`),
						"must contain only alphanumeric characters and _-",
					),
				},
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
				Description:         "The list of ports associated with the portgroup. (Update Supported)",
				MarkdownDescription: "The list of ports associated with the portgroup. (Update Supported)",
			},
			"protocol": schema.StringAttribute{
				Required:            true,
				Description:         "The portgroup protocol. Protocols: SCSI_FC, iSCSI, NVMe_FC, NVMe_TCP",
				MarkdownDescription: "The portgroup protocol. Protocols: SCSI_FC, iSCSI, NVMe_FC, NVMe_TCP",
				Validators: []validator.String{
					stringvalidator.OneOf("SCSI_FC", "iSCSI", "NVMe_FC", "NVMe_TCP"),
				},
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

	pgResponse, _, err := helper.CreatePortGroup(ctx, *r.client, plan)

	if err != nil {
		errStr := constants.CreatePGDetailErrorMsg + plan.Name.ValueString() + " with error: "
		msgStr := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating port group", msgStr,
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
	pgResponse, _, err := helper.ReadPortgroupByID(ctx, *r.client, pgID)
	if err != nil {
		errStr := constants.ReadPGDetailsErrorMsg + pgID + " with error: "
		msgStr := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading port group", msgStr,
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
		return
	}

	portGroupID := pgState.ID.ValueString()

	if helper.IsParamUpdated(updatedParams, "name") {
		portGroupID = pgPlan.Name.ValueString()
	}

	pgResponse, _, err := helper.ReadPortgroupByID(ctx, *r.client, portGroupID)
	if err != nil {
		errStr := constants.UpdatePGDetailsErrMsg + pgPlan.Name.ValueString() + " with error: "
		msgStr := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading port group", msgStr,
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
	_, err := r.client.PmaxOpenapiClient.SLOProvisioningApi.DeletePortGroup(ctx, r.client.SymmetrixID, pgID).Execute()

	if err != nil {
		errStr := constants.DeletePGDetailsErrorMsg + pgID + " with error: "
		msgStr := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error deleting port group", msgStr,
		)
	}
	tflog.Info(ctx, "delete portgroup completed")
}

// ImportState import resource.
func (r *PortGroup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "importing port group state")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
