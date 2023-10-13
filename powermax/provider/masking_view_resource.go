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
	pmax "dell/powermax-go-client"
	"fmt"
	"regexp"
	"terraform-provider-powermax/client"
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

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &maskingView{}
	_ resource.ResourceWithConfigure   = &maskingView{}
	_ resource.ResourceWithImportState = &maskingView{}
)

// NewMaskingView returns the masking view resource object.
func NewMaskingView() resource.Resource {
	return &maskingView{}
}

// maskingView defines the resource implementation.
type maskingView struct {
	client *client.Client
}

func (r *maskingView) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_maskingview"
}

func (r *maskingView) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Resource for managing MaskingViews in PowerMax array. PowerMax masking views are a container of a storage group, a port group, and an initiator group, and makes the storage group visible to the host. Devices are masked and mapped automatically. The groups must contain some device entries.",
		Description:         "Resource for managing MaskingViews in PowerMax array. PowerMax masking views are a container of a storage group, a port group, and an initiator group, and makes the storage group visible to the host. Devices are masked and mapped automatically. The groups must contain some device entries.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the masking view.",
				MarkdownDescription: "The ID of the masking view.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Unique identifier of the masking view. (Update Supported)",
				MarkdownDescription: "Unique identifier of the masking view. (Update Supported)",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(64),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9_-]*$`),
						"must contain only alphanumeric characters and _-",
					),
				},
			},
			"storage_group_id": schema.StringAttribute{
				Required:            true,
				Description:         "The storage group id of the masking view.",
				MarkdownDescription: "The storage group id of the masking view.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"host_id": schema.StringAttribute{
				Required:            true,
				Description:         "The host id of the masking view.",
				MarkdownDescription: "The host id of the masking view.",
			},
			"host_group_id": schema.StringAttribute{
				Required:            true,
				Description:         "The host group id of the masking view.",
				MarkdownDescription: "The host group id of the masking view.",
			},
			"port_group_id": schema.StringAttribute{
				Required:            true,
				Description:         "The port group id of the masking view.",
				MarkdownDescription: "The port group id of the masking view.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

func (r *maskingView) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *maskingView) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Masking View...")
	var plan models.MaskingViewResourceModel

	// Read Terraform plan into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var hostOrHostGroupID string
	var isHost = false
	if plan.HostID.ValueString() != "" && plan.HostGroupID.ValueString() == "" {
		hostOrHostGroupID = plan.HostID.ValueString()
		isHost = true
	} else if plan.HostID.ValueString() == "" && plan.HostGroupID.ValueString() != "" {
		hostOrHostGroupID = plan.HostGroupID.ValueString()
	} else {
		resp.Diagnostics.AddError(
			"Specify either host_id or host_group_id.",
			"unexpected error: Specify either host_id or host_group_id",
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Calling api to create MaskingView - %s", plan.Name.ValueString()))

	maskingView, _, err := helper.CreateMaskingView(ctx, *r.client, plan, hostOrHostGroupID, isHost)

	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error creating masking view", message)

		return
	}

	tflog.Debug(ctx, "created a resource", map[string]interface{}{
		"masking view": maskingView,
	})

	tflog.Debug(ctx, fmt.Sprintf("Calling api to get MaskingView - %s", plan.Name.ValueString()))
	maskingView, _, err = helper.GetMaskingView(ctx, *r.client, plan.Name.ValueString())

	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading masking view", message)
		// Attempt to clean up the errored masking view after the host/hostgroup mistake
		_, delErr := r.client.PmaxOpenapiClient.SLOProvisioningApi.DeleteMaskingView(ctx, r.client.SymmetrixID, plan.Name.ValueString()).Execute()
		if delErr != nil {
			errStr := ""
			message := helper.GetErrorString(delErr, errStr)
			tflog.Error(ctx, "Error deleting maskingview after host_group error"+message)
		}
		return
	}

	err = helper.CopyFields(ctx, maskingView, &plan)
	if err != nil {
		// Attempt to clean up the errored masking view after the host/hostgroup mistake
		_, delErr := r.client.PmaxOpenapiClient.SLOProvisioningApi.DeleteMaskingView(ctx, r.client.SymmetrixID, plan.Name.ValueString()).Execute()
		if delErr != nil {
			errStr := ""
			message := helper.GetErrorString(delErr, errStr)
			tflog.Error(ctx, "Error deleting maskingview after host_group error"+message)
		}
		resp.Diagnostics.AddError("Error copying masking view fields", err.Error())
		return
	}

	plan.StorageGroupID = types.StringValue(*maskingView.StorageGroupId)
	if plan.HostGroupID.ValueString() != "" && maskingView.HostId != nil {
		resp.Diagnostics.AddError("Error creating masking view", fmt.Sprintf("The host_group_id '%s' is actually a host_id, change '%s' to host_id to create a masking view with this host", plan.HostGroupID, plan.HostGroupID))
		// Attempt to clean up the errored masking view after the host/hostgroup mistake
		_, err := r.client.PmaxOpenapiClient.SLOProvisioningApi.DeleteMaskingView(ctx, r.client.SymmetrixID, plan.Name.ValueString()).Execute()
		if err != nil {
			errStr := ""
			message := helper.GetErrorString(err, errStr)
			tflog.Error(ctx, "Error deleting maskingview after host_group error"+message)
			return
		}
		return
	}
	if plan.HostID.ValueString() != "" && maskingView.HostGroupId != nil {
		resp.Diagnostics.AddError("Error creating masking view", fmt.Sprintf("The host_id '%s' is actually a host_group_id, change '%s' to host_group_id to create a masking view with this host_group", plan.HostID, plan.HostID))
		// Attempt to clean up the errored masking view after the host/hostgroup mistake
		_, err := r.client.PmaxOpenapiClient.SLOProvisioningApi.DeleteMaskingView(ctx, r.client.SymmetrixID, plan.Name.ValueString()).Execute()
		if err != nil {
			errStr := ""
			message := helper.GetErrorString(err, errStr)
			tflog.Error(ctx, "Error deleting maskingview after host error"+message)
			return
		}
		return
	}
	if maskingView.HostId != nil {
		plan.HostID = types.StringValue(*maskingView.HostId)
	}
	if maskingView.HostGroupId != nil {
		plan.HostGroupID = types.StringValue(*maskingView.HostGroupId)
	}
	plan.PortGroupID = types.StringValue(*maskingView.PortGroupId)
	plan.Name = types.StringValue(maskingView.MaskingViewId)
	plan.ID = types.StringValue(maskingView.MaskingViewId)
	// Save plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Done with Create Masking View resource")
}

func (r *maskingView) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Masing View...")
	var state models.MaskingViewResourceModel

	// Read Terraform prior state into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Calling api to get MaskingView - %s", state.Name.ValueString()))
	getMaskingViewReq := r.client.PmaxOpenapiClient.SLOProvisioningApi.GetMaskingView(ctx, r.client.SymmetrixID, state.Name.ValueString())
	maskingView, _, err := getMaskingViewReq.Execute()

	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading masking view", message)

		return
	}

	err = helper.CopyFields(ctx, maskingView, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error copying masking view fields", err.Error())
		return
	}

	state.StorageGroupID = types.StringValue(*maskingView.StorageGroupId)
	if maskingView.HostId != nil {
		state.HostID = types.StringValue(*maskingView.HostId)
	}
	if maskingView.HostGroupId != nil {
		state.HostGroupID = types.StringValue(*maskingView.HostGroupId)
	}
	state.PortGroupID = types.StringValue(*maskingView.PortGroupId)
	state.Name = types.StringValue(maskingView.MaskingViewId)
	state.ID = types.StringValue(maskingView.MaskingViewId)
	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Info(ctx, "Done with Read Masking View resource")
}

// Update: support rename.
func (r *maskingView) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Masking View...")
	// Read Terraform plan into the model
	var plan models.MaskingViewResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform state into the model
	var state models.MaskingViewResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// prompt error on change in maskingView's hostGroup, portGroup or storageGroup, as we can't update the them after the creation
	if !plan.StorageGroupID.Equal(state.StorageGroupID) || !plan.PortGroupID.Equal(state.PortGroupID) || !plan.HostID.Equal(state.HostID) ||
		!plan.HostGroupID.Equal(state.HostGroupID) {
		resp.Diagnostics.AddError(
			"maskingView's host, hostGroup, portGroup or storageGroup cannot be update after creation.",
			"unexpected error: maskingView's host, hostGroup, portGroup or storageGroup change is not supported",
		)
		return
	}

	// Rename masking view
	if !plan.Name.Equal(state.Name) {
		tflog.Debug(ctx, fmt.Sprintf("Calling api to rename MaskingView from %s to %s", state.Name.ValueString(), plan.Name.ValueString()))

		renameMaskingViewParam := pmax.NewRenameMaskingViewParam(plan.Name.ValueString())
		rename := pmax.EditMaskingViewActionParam{
			RenameMaskingViewParam: renameMaskingViewParam,
		}
		modifyReq := r.client.PmaxOpenapiClient.SLOProvisioningApi.ModifyMaskingView(ctx, r.client.SymmetrixID, state.Name.ValueString())
		editParam := pmax.NewEditMaskingViewParam(rename)
		modifyReq = modifyReq.EditMaskingViewParam(*editParam)
		_, _, err := modifyReq.Execute()
		if err != nil {
			errStr := ""
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Error renaming masking view", message)

			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Calling api to get MaskingView - %s", plan.Name.ValueString()))
	getMaskingViewReq := r.client.PmaxOpenapiClient.SLOProvisioningApi.GetMaskingView(ctx, r.client.SymmetrixID, plan.Name.ValueString())
	maskingView, _, err := getMaskingViewReq.Execute()
	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading masking view", message)
		return
	}

	err = helper.CopyFields(ctx, maskingView, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error copying masking view fields", err.Error())
		return
	}

	state.StorageGroupID = types.StringValue(*maskingView.StorageGroupId)
	if maskingView.HostId != nil {
		state.HostID = types.StringValue(*maskingView.HostId)
	}
	if maskingView.HostGroupId != nil {
		state.HostGroupID = types.StringValue(*maskingView.HostGroupId)
	}
	state.PortGroupID = types.StringValue(*maskingView.PortGroupId)
	state.Name = types.StringValue(maskingView.MaskingViewId)
	state.ID = types.StringValue(maskingView.MaskingViewId)
	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Info(ctx, "Done with Update Masking View resource")
}

func (r *maskingView) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Masking View...")
	var state models.MaskingViewResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Calling api to delete MaskingView - %s", state.Name.ValueString()))
	delReq := r.client.PmaxOpenapiClient.SLOProvisioningApi.DeleteMaskingView(ctx, r.client.SymmetrixID, state.Name.ValueString())
	_, err := delReq.Execute()
	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete masking view, got error: %s", message))
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Info(ctx, "Done with Delete Masking View resource")
}

func (r *maskingView) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
