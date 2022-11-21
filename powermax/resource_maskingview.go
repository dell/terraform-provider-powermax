package powermax

import (
	"context"
	"fmt"
	"terraform-provider-powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type resourceMaskingViewType struct{}

// MaskingView Resource schema
func (r resourceMaskingViewType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Resource to manage Maskingview in PowerMax array.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The ID of the maskingview.",
				MarkdownDescription: "The ID of the maskingview.",
			},
			"name": {
				Type:                types.StringType,
				Required:            true,
				Description:         "The name of the maskingview.",
				MarkdownDescription: "The name of the maskingview.",
			},
			"storage_group_id": {
				Type:                types.StringType,
				Required:            true,
				Description:         "The ID of the storagegroup associated with maskingview.",
				MarkdownDescription: "The ID of the storagegroup associated with maskingview.",
			},
			"port_group_id": {
				Type:                types.StringType,
				Required:            true,
				Description:         "The ID of the portgroup associated with maskingview.",
				MarkdownDescription: "The ID of the portgroup associated with maskingview.",
			},
			"host_id": {
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "The ID of the host associated with maskingview.",
				MarkdownDescription: "The ID of the host associated with maskingview.",
			},
			"host_group_id": {
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "The ID of the hostgroup associated with maskingview. Either of host_id/host_group_id is expected but not both",
				MarkdownDescription: "The ID of the hostgroup associated with maskingview. Either of host_id/host_group_id is expected but not both",
			},
		},
	}, nil
}

// NewResource is a wrapper around provider
func (r resourceMaskingViewType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceMaskingView{
		p: *(p.(*provider)),
	}, nil
}

type resourceMaskingView struct {
	p provider
}

// Create Maskingview
func (r resourceMaskingView) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Debug(ctx, "creating masking view")
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var planMaskingView models.MaskingView
	diags := req.Plan.Get(ctx, &planMaskingView)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if planMaskingView.HostID.Value != "" && planMaskingView.HostGroupID.Value != "" {
		resp.Diagnostics.AddError(
			"Error creating maskingview",
			fmt.Sprintf("%s: %s because only one of host_id/host_group_id is expected as input", CreateMVDetailErrorMsg, planMaskingView.Name.Value),
		)
		return
	}

	tflog.Debug(ctx, "calling create masking view on pmax client", map[string]interface{}{
		"symmetrixID":     r.p.client.SymmetrixID,
		"maskingViewName": planMaskingView.Name.Value,
		"storageGroup":    planMaskingView.StorageGroupID.Value,
		"hostID":          planMaskingView.HostID.Value,
		"portGroupID":     planMaskingView.PortGroupID.Value,
	})
	hostOrHostGroupID := planMaskingView.HostID.Value
	isHost := planMaskingView.HostID.Value != ""
	if !isHost {
		hostOrHostGroupID = planMaskingView.HostGroupID.Value
	}

	mvResponse, err := r.p.client.PmaxClient.CreateMaskingView(ctx, r.p.client.SymmetrixID, planMaskingView.Name.Value,
		planMaskingView.StorageGroupID.Value, hostOrHostGroupID, isHost, planMaskingView.PortGroupID.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating maskingview",
			fmt.Sprintf("%s %s with error: %s", CreateMVDetailErrorMsg, planMaskingView.Name.Value, err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "create masking view response", map[string]interface{}{
		"mvResponse": mvResponse,
	})

	var stateMaskingView models.MaskingView

	tflog.Debug(ctx, "updating masking view state", map[string]interface{}{
		"stateMaskingView": stateMaskingView,
		"mvResponse":       mvResponse,
	})
	updateMaskingViewState(&stateMaskingView, mvResponse)
	diags = resp.State.Set(ctx, stateMaskingView)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "create masking view completed")
}

// Read Maskingview
func (r resourceMaskingView) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	tflog.Debug(ctx, "reading masking view")
	var stateMaskingView models.MaskingView
	diags := req.State.Get(ctx, &stateMaskingView)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get maskingview ID from API and then update what is in state from what the API returns
	mvID := stateMaskingView.ID.Value
	tflog.Debug(ctx, "calling get masking view by ID on pmax client", map[string]interface{}{
		"symmetrixID":   r.p.client.SymmetrixID,
		"maskingViewID": mvID,
	})
	mvResponse, err := r.p.client.PmaxClient.GetMaskingViewByID(ctx, r.p.client.SymmetrixID, mvID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading maskingview",
			fmt.Sprintf("%s %s with error: %s", ReadMVDetailsErrorMsg, stateMaskingView.Name.Value, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, "updating masking view state", map[string]interface{}{
		"stateMaskingView": stateMaskingView,
		"mvResponse":       mvResponse,
	})
	updateMaskingViewState(&stateMaskingView, mvResponse)
	diags = resp.State.Set(ctx, stateMaskingView)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "reading masking view completed")
}

// Update Maskingview
// Supported updates: name
func (r resourceMaskingView) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Debug(ctx, "updating masking view")
	var planMaskingView models.MaskingView
	diags := req.Plan.Get(ctx, &planMaskingView)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateMaskingView models.MaskingView
	diags = req.State.Get(ctx, &stateMaskingView)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling rename masking view on pmax client", map[string]interface{}{
		"symmetrixID":   r.p.client.SymmetrixID,
		"maskingViewID": stateMaskingView.ID.Value,
		"newName":       planMaskingView.Name.Value,
	})
	mvResponse, err := r.p.client.PmaxClient.RenameMaskingView(ctx, r.p.client.SymmetrixID, stateMaskingView.ID.Value, planMaskingView.Name.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error renaming maskingview",
			fmt.Sprintf("%s %s with error: %s", RenameMVDetailErrorMsg, stateMaskingView.ID.Value, err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "rename masking view response", map[string]interface{}{
		"mvResponse": mvResponse,
	})

	tflog.Debug(ctx, "updating masking view state", map[string]interface{}{
		"stateMaskingView": stateMaskingView,
		"mvResponse":       mvResponse,
	})
	updateMaskingViewState(&stateMaskingView, mvResponse)
	diags = resp.State.Set(ctx, stateMaskingView)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "update masking view completed")
}

// Delete Maskingview
func (r resourceMaskingView) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Debug(ctx, "deleting masking view")
	var stateMaskingView models.MaskingView
	diags := req.State.Get(ctx, &stateMaskingView)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	mvID := stateMaskingView.ID.Value
	tflog.Debug(ctx, "calling delete masking view on pmax client", map[string]interface{}{
		"symmetrixID":   r.p.client.SymmetrixID,
		"maskingViewID": mvID,
	})
	err := r.p.client.PmaxClient.DeleteMaskingView(ctx, r.p.client.SymmetrixID, mvID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting maskingview",
			fmt.Sprintf("%s %s with error: %s", DeleteMVDetailsErrorMsg, stateMaskingView.Name.Value, err.Error()),
		)
	}
	resp.State.RemoveResource(ctx)
	tflog.Debug(ctx, "deleting masking view completed")
}

func (r resourceMaskingView) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tflog.Debug(ctx, "importing masking view")
	var stateMaskingView models.MaskingView
	mvID := req.ID
	tflog.Debug(ctx, "calling get masking view by ID on pmax client", map[string]interface{}{
		"symmetrixID":   r.p.client.SymmetrixID,
		"maskingViewID": mvID,
	})
	mvResponse, err := r.p.client.PmaxClient.GetMaskingViewByID(ctx, r.p.client.SymmetrixID, mvID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing maskingview",
			fmt.Sprintf("%s with masking view id %s with error: %s", ImportMVDetailsErrorMsg, mvID, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, "Get masking view by ID response", map[string]interface{}{
		"mvResponse": mvResponse,
	})

	tflog.Debug(ctx, "updating masking view state after import")
	updateMaskingViewState(&stateMaskingView, mvResponse)
	diags := resp.State.Set(ctx, stateMaskingView)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "completed import masking view")
}
