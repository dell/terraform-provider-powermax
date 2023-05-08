// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.

package provider

import (
	"context"
	"fmt"
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

func newMaskingView() resource.Resource {
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
		MarkdownDescription: "Resource for managing MaskingViews in PowerMax array. Updates are supported for the following parameters: `name`.",
		Description:         "Resource for managing MaskingViews in PowerMax array. Updates are supported for the following parameters: `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the masking view.",
				MarkdownDescription: "The ID of the masking view.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Unique identifier of the masking view.",
				MarkdownDescription: "Unique identifier of the masking view.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
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
	if plan.HostID.ValueString() != "" && plan.HostGroupID.ValueString() == "" {
		hostOrHostGroupID = plan.HostID.ValueString()
	} else if plan.HostID.ValueString() == "" && plan.HostGroupID.ValueString() != "" {
		hostOrHostGroupID = plan.HostGroupID.ValueString()
	} else {
		resp.Diagnostics.AddError(
			"The host_id or host_group_id only needs to be specified one.",
			"unexpected error: The host_id or host_group_id only needs to be specified one.",
		)
	}

	tflog.Debug(ctx, fmt.Sprintf("Calling api to create MaskingView - %s", plan.Name.ValueString()))

	maskingView, err := r.client.PmaxClient.CreateMaskingView(ctx, r.client.SymmetrixID, plan.Name.ValueString(), plan.StorageGroupID.ValueString(), hostOrHostGroupID, false, plan.PortGroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create masking view, got error: %s", err.Error()))
		return
	}

	tflog.Debug(ctx, "created a resource", map[string]interface{}{
		"masking view": maskingView,
	})

	tflog.Debug(ctx, fmt.Sprintf("Calling api to get MaskingView - %s", plan.Name.ValueString()))
	maskingView, err = r.client.PmaxClient.GetMaskingViewByID(ctx, r.client.SymmetrixID, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading masking view", err.Error())
		return
	}

	err = helper.CopyFields(ctx, maskingView, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Error copying masking view fields", err.Error())
		return
	}

	plan.Name = types.StringValue(maskingView.MaskingViewID)
	plan.ID = types.StringValue(maskingView.MaskingViewID)
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
	maskingView, err := r.client.PmaxClient.GetMaskingViewByID(ctx, r.client.SymmetrixID, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading masking view", err.Error())
		return
	}

	err = helper.CopyFields(ctx, maskingView, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error copying masking view fields", err.Error())
		return
	}

	state.Name = types.StringValue(maskingView.MaskingViewID)
	state.ID = types.StringValue(maskingView.MaskingViewID)
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
	}

	// Rename masking view
	if !plan.Name.Equal(state.Name) {
		tflog.Debug(ctx, fmt.Sprintf("Calling api to rename MaskingView from %s to %s", state.Name.ValueString(), plan.Name.ValueString()))

		_, err := r.client.PmaxClient.RenameMaskingView(ctx, r.client.SymmetrixID, state.Name.ValueString(), plan.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error renaming masking view", err.Error())
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Calling api to get MaskingView - %s", plan.Name.ValueString()))
	maskingView, err := r.client.PmaxClient.GetMaskingViewByID(ctx, r.client.SymmetrixID, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading masking view", err.Error())
		return
	}

	err = helper.CopyFields(ctx, maskingView, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error copying masking view fields", err.Error())
		return
	}

	state.Name = types.StringValue(maskingView.MaskingViewID)
	state.ID = types.StringValue(maskingView.MaskingViewID)
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
	err := r.client.PmaxClient.DeleteMaskingView(ctx, r.client.SymmetrixID, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete masking view, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Info(ctx, "Done with Delete Masking View resource")
}

func (r *maskingView) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
