package powermax

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-powermax/models"

	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type resourceHostType struct{}

// Host Resource schema
func (r resourceHostType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	hostFlagDefaultPlanModifier := tfsdk.AttributePlanModifiers{
		DefaultAttribute(types.Object{
			Attrs: map[string]attr.Value{
				"override": types.Bool{Value: false},
				"enabled":  types.Bool{Value: false},
			},
			AttrTypes: map[string]attr.Type{
				"override": types.BoolType,
				"enabled":  types.BoolType,
			},
		}),
	}
	boolDefaultPlanModifier := tfsdk.AttributePlanModifiers{
		DefaultAttribute(types.Bool{Value: false}),
	}
	hostFlagNestedAttr := tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
		"override": {
			Type:          types.BoolType,
			Optional:      true,
			Computed:      true,
			PlanModifiers: boolDefaultPlanModifier,
		},
		"enabled": {
			Type:          types.BoolType,
			Optional:      true,
			Computed:      true,
			PlanModifiers: boolDefaultPlanModifier,
		},
	})
	return tfsdk.Schema{
		MarkdownDescription: "Resource to manage hosts in PowerMax array. Updates are supported for the following parameters: `name`, `initiators`, `host_flags`.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The ID of the host.",
				MarkdownDescription: "The ID of the host.",
			},
			"name": {
				Type:                types.StringType,
				Required:            true,
				Description:         "The name of the host.",
				MarkdownDescription: "The name of the host.",
			},

			"host_flags": {
				Required: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"volume_set_addressing": {
						Optional:      true,
						Computed:      true,
						Attributes:    hostFlagNestedAttr,
						PlanModifiers: hostFlagDefaultPlanModifier,
					},
					"disable_q_reset_on_ua": {
						Optional:      true,
						Computed:      true,
						Attributes:    hostFlagNestedAttr,
						PlanModifiers: hostFlagDefaultPlanModifier,
					},
					"environ_set": {
						Optional:      true,
						Computed:      true,
						Attributes:    hostFlagNestedAttr,
						PlanModifiers: hostFlagDefaultPlanModifier,
					},
					"openvms": {
						Optional:      true,
						Computed:      true,
						Attributes:    hostFlagNestedAttr,
						PlanModifiers: hostFlagDefaultPlanModifier,
					},
					"avoid_reset_broadcast": {
						Optional:      true,
						Computed:      true,
						Attributes:    hostFlagNestedAttr,
						PlanModifiers: hostFlagDefaultPlanModifier,
					},
					"scsi_3": {
						Optional:      true,
						Computed:      true,
						Attributes:    hostFlagNestedAttr,
						PlanModifiers: hostFlagDefaultPlanModifier,
					},
					"spc2_protocol_version": {
						Optional:      true,
						Computed:      true,
						Attributes:    hostFlagNestedAttr,
						PlanModifiers: hostFlagDefaultPlanModifier,
					},
					"scsi_support1": {
						Optional:      true,
						Computed:      true,
						Attributes:    hostFlagNestedAttr,
						PlanModifiers: hostFlagDefaultPlanModifier,
					},
					"consistent_lun": {
						Type:          types.BoolType,
						Optional:      true,
						Computed:      true,
						PlanModifiers: boolDefaultPlanModifier,
					},
				}),
				Description:         "Flags set for the host.",
				MarkdownDescription: "Flags set for the host.",
			},
			"initiators": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Required:            true,
				Description:         "The initiators associated with the host.",
				MarkdownDescription: "The initiators associated with the host.",
			},
			"numofmaskingviews": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "The number of masking views associated with the host.",
				MarkdownDescription: "The number of masking views associated with the host.",
			},
			"numofinitiators": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "The number of initiators associated with the host.",
				MarkdownDescription: "The number of initiators associated with the host.",
			},
			"numofhostgroups": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "The number of hostgroups associated with the host.",
				MarkdownDescription: "The number of hostgroups associated with the host.",
			},
			"type": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "Specifies the type of host.",
				MarkdownDescription: "Specifies the type of host.",
			},
			"maskingview": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Computed:            true,
				Description:         "The masking views associated with the host.",
				MarkdownDescription: "The masking views associated with the host.",
			},
			"powerpath_hosts": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Computed:            true,
				Description:         "The powerpath hosts associated with the host.",
				MarkdownDescription: "The powerpath hosts associated with the host.",
			},
			"numofpowerpathhosts": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "The number of powerpath hosts associated with the host.",
				MarkdownDescription: "The number of powerpath hosts associated with the host.",
			},
			"bw_limit": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "Specifies the bandwidth limit for a host.",
				MarkdownDescription: "Specifies the bandwidth limit for a host.",
			},
			"port_flags_override": {
				Type:                types.BoolType,
				Computed:            true,
				Description:         "States whether port flags override is enabled on the host.",
				MarkdownDescription: "States whether port flags override is enabled on the host.",
			},
		},
	}, nil
}

// NewResource is a wrapper around provider
func (r resourceHostType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceHost{
		p: *(p.(*provider)),
	}, nil
}

type resourceHost struct {
	p provider
}

// Create Host
func (r resourceHost) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Debug(ctx, "creating host")
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	tflog.Debug(ctx, "setting host plan")
	var planHost models.Host
	diags := req.Plan.Get(ctx, &planHost)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	initiators := make([]string, len(planHost.Initiators.Elems))

	if len(planHost.Initiators.Elems) > 0 {
		for index, initiator := range planHost.Initiators.Elems {
			initiators[index] = strings.Trim(initiator.String(), "\"")
		}
	}

	tflog.Debug(ctx, "preparing host flags")
	hostFlags := pmaxTypes.HostFlags{
		VolumeSetAddressing: &pmaxTypes.HostFlag{
			Enabled:  planHost.HostFlags.VolumeSetAddressing.Enabled.Value,
			Override: planHost.HostFlags.VolumeSetAddressing.Override.Value,
		},
		DisableQResetOnUA: &pmaxTypes.HostFlag{
			Enabled:  planHost.HostFlags.DisableQResetOnUa.Enabled.Value,
			Override: planHost.HostFlags.DisableQResetOnUa.Override.Value,
		},
		EnvironSet: &pmaxTypes.HostFlag{
			Enabled:  planHost.HostFlags.EnvironSet.Enabled.Value,
			Override: planHost.HostFlags.EnvironSet.Override.Value,
		},
		AvoidResetBroadcast: &pmaxTypes.HostFlag{
			Enabled:  planHost.HostFlags.AvoidResetBroadcast.Enabled.Value,
			Override: planHost.HostFlags.AvoidResetBroadcast.Override.Value,
		},
		OpenVMS: &pmaxTypes.HostFlag{
			Enabled:  planHost.HostFlags.Openvms.Enabled.Value,
			Override: planHost.HostFlags.Openvms.Override.Value,
		},
		SCSI3: &pmaxTypes.HostFlag{
			Enabled:  planHost.HostFlags.Scsi3.Enabled.Value,
			Override: planHost.HostFlags.Scsi3.Override.Value,
		},
		Spc2ProtocolVersion: &pmaxTypes.HostFlag{
			Enabled:  planHost.HostFlags.Spc2ProtocolVersion.Enabled.Value,
			Override: planHost.HostFlags.Spc2ProtocolVersion.Override.Value,
		},
		SCSISupport1: &pmaxTypes.HostFlag{
			Enabled:  planHost.HostFlags.ScsiSupport1.Enabled.Value,
			Override: planHost.HostFlags.ScsiSupport1.Override.Value,
		},
		ConsistentLUN: planHost.HostFlags.ConsistentLun.Value,
	}

	tflog.Debug(ctx, "calling create host pmax client", map[string]interface{}{
		"symmetrixID": r.p.client.SymmetrixID,
		"host":        planHost.Name.Value,
		"initiators":  initiators,
		"hostFlags":   hostFlags,
	})
	hostResponse, err := r.p.client.PmaxClient.CreateHost(ctx, r.p.client.SymmetrixID, planHost.Name.Value, initiators, &hostFlags)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating host",
			CreateHostDetailErrorMsg+planHost.Name.Value+"with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "create host response", map[string]interface{}{
		"Create Host Response": hostResponse,
	})

	hostState := models.Host{}
	tflog.Debug(ctx, "updating host state", map[string]interface{}{
		"host state":   hostState,
		"initiators":   initiators,
		"hostResponse": hostResponse,
	})
	updateHostState(&hostState, initiators, hostResponse)
	diags = resp.State.Set(ctx, hostState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "completed create host")
}

// Read Host
func (r resourceHost) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	tflog.Debug(ctx, "reading host")
	var hostState models.Host
	diags := req.State.Get(ctx, &hostState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get order from API and then update what is in state from what the API returns
	hostID := hostState.ID.Value
	tflog.Debug(ctx, "fetching host by ID", map[string]interface{}{
		"symmetrixID": r.p.client.SymmetrixID,
		"hostID":      hostID,
	})
	hostResponse, err := r.p.client.PmaxClient.GetHostByID(ctx, r.p.client.SymmetrixID, hostID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading host",
			ReadHostDetailsErrorMsg+hostID+" with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Get Host By ID response", map[string]interface{}{
		"Host Response": hostResponse,
	})

	initiators := make([]string, len(hostState.Initiators.Elems))

	if len(hostState.Initiators.Elems) > 0 {
		for index, initiator := range hostState.Initiators.Elems {
			initiators[index] = strings.Trim(initiator.String(), "\"")
		}
	}

	tflog.Debug(ctx, "updating host state")
	updateHostState(&hostState, initiators, hostResponse)
	diags = resp.State.Set(ctx, hostState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "completed read host")
}

// Update Host
// Supported updates: name, initiators, host flags
func (r resourceHost) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Debug(ctx, "updating host")
	var plan models.Host
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "fetched host details from plan")

	var state models.Host
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update host on pmax client", map[string]interface{}{
		"plan":  plan,
		"state": state,
	})
	updatedParams, updateFailedParameters, errMessages := updateHost(ctx, r.p.client, plan, state)
	if len(errMessages) > 0 || len(updateFailedParameters) > 0 {
		errMessage := strings.Join(errMessages, ",\n")
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to update all parameters of Host, updated parameters are %v and parameters failed to update are %v", updatedParams, updateFailedParameters),
			errMessage)
	}
	tflog.Debug(ctx, "update host response", map[string]interface{}{
		"updatedParams":          updatedParams,
		"updateFailedParameters": updateFailedParameters,
		"error messages":         errMessages,
	})

	hostID := plan.Name.Value
	tflog.Debug(ctx, "calling get host by ID on pmax client", map[string]interface{}{
		"SymmetrixID": r.p.client.SymmetrixID,
		"hostID":      hostID,
	})
	hostResponse, err := r.p.client.PmaxClient.GetHostByID(ctx, r.p.client.SymmetrixID, hostID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading host",
			ReadHostDetailsErrorMsg+hostID+" with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "get host by ID response", map[string]interface{}{
		"Host Response": hostResponse,
	})

	initiators := make([]string, len(plan.Initiators.Elems))
	if len(plan.Initiators.Elems) > 0 {
		for index, initiator := range plan.Initiators.Elems {
			initiators[index] = strings.Trim(initiator.String(), "\"")
		}
	}

	tflog.Debug(ctx, "updating update host state", map[string]interface{}{
		"state":        state,
		"initiators":   initiators,
		"hostResponse": hostResponse,
	})
	updateHostState(&state, initiators, hostResponse)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "update host completed")
}

// Delete Host
func (r resourceHost) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Debug(ctx, "deleting host")
	var hostState models.Host
	diags := req.State.Get(ctx, &hostState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	hostID := hostState.ID.Value
	tflog.Debug(ctx, "deleting host by host ID", map[string]interface{}{
		"symmetrixID": r.p.client.SymmetrixID,
		"hostID":      hostID,
	})
	err := r.p.client.PmaxClient.DeleteHost(ctx, r.p.client.SymmetrixID, hostID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting host",
			DeleteHostDetailsErrorMsg+hostID+" with error: "+err.Error(),
		)
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
	tflog.Debug(ctx, "deleting host complete")
}

func (r resourceHost) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tflog.Debug(ctx, "importing host state")
	var hostState models.Host
	hostID := req.ID
	tflog.Debug(ctx, "fetching host by ID", map[string]interface{}{
		"symmetrixID": r.p.client.SymmetrixID,
		"hostID":      hostID,
	})
	hostResponse, err := r.p.client.PmaxClient.GetHostByID(ctx, r.p.client.SymmetrixID, hostID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading host",
			ImportHostDetailsErrorMsg+hostID+" with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Get Host By ID response", map[string]interface{}{
		"Host Response": hostResponse,
	})

	tflog.Debug(ctx, "updating host state after import")
	updateHostState(&hostState, hostResponse.Initiators, hostResponse)
	diags := resp.State.Set(ctx, hostState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "completed import host")

}
