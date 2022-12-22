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

type resourceHostGroupType struct{}

// HostGroup Resource schema
func (r resourceHostGroupType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	hostFlagDefaultPlanModifier := tfsdk.AttributePlanModifiers{
		tfsdk.UseStateForUnknown(),
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
		tfsdk.UseStateForUnknown(),
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
		MarkdownDescription: "Resource to manage hostgroup in PowerMax array. Updates are supported for the following parameters: `name`, `host_ids`, `host_flags`, `consistent_lun`.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "The ID of the hostgroup.",
				MarkdownDescription: "The ID of the hostgroup.",
			},
			"name": {
				Type:                types.StringType,
				Required:            true,
				Description:         "The name of the hostgroup.",
				MarkdownDescription: "The name of the hostgroup.",
			},
			"host_ids": {
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Required:            true,
				Description:         "A list of host ids associated with the hostgroup.",
				MarkdownDescription: "The masking views associated with the hostgroup.",
				Validators: []tfsdk.AttributeValidator{
					SizeAtLeast(1),
				},
			},
			"host_flags": {
				Required: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"volume_set_addressing": {
						Optional:            true,
						Computed:            true,
						Attributes:          hostFlagNestedAttr,
						Description:         "It enables the volume set addressing mode.",
						MarkdownDescription: "It enables the volume set addressing mode.",
						PlanModifiers:       hostFlagDefaultPlanModifier,
					},
					"disable_q_reset_on_ua": {
						Optional:            true,
						Computed:            true,
						Attributes:          hostFlagNestedAttr,
						Description:         "It is used for hosts that do not expect the queue to be flushed on a 0629 sense.",
						MarkdownDescription: "It is used for hosts that do not expect the queue to be flushed on a 0629 sense.",
						PlanModifiers:       hostFlagDefaultPlanModifier,
					},
					"environ_set": {
						Optional:            true,
						Computed:            true,
						Attributes:          hostFlagNestedAttr,
						Description:         "It enables the environmental error reporting by the storage system to the host on the specific port.",
						MarkdownDescription: "It enables the environmental error reporting by the storage system to the host on the specific port.",
						PlanModifiers:       hostFlagDefaultPlanModifier,
					},
					"openvms": {
						Optional:            true,
						Computed:            true,
						Attributes:          hostFlagNestedAttr,
						Description:         "This attribute enables an Open VMS fibre connection.",
						MarkdownDescription: "This attribute enables an Open VMS fibre connection.",
						PlanModifiers:       hostFlagDefaultPlanModifier,
					},
					"avoid_reset_broadcast": {
						Optional:            true,
						Computed:            true,
						Attributes:          hostFlagNestedAttr,
						Description:         "It enables a SCSI bus reset to only occur to the port that received the reset.",
						MarkdownDescription: "It enables a SCSI bus reset to only occur to the port that received the reset.",
						PlanModifiers:       hostFlagDefaultPlanModifier,
					},
					"scsi_3": {
						Optional:            true,
						Computed:            true,
						Attributes:          hostFlagNestedAttr,
						Description:         "Alters the inquiry data to report that the storage system supports the SCSI-3 protocol.",
						MarkdownDescription: "Alters the inquiry data to report that the storage system supports the SCSI-3 protocol.",
						PlanModifiers:       hostFlagDefaultPlanModifier,
					},
					"spc2_protocol_version": {
						Optional:            true,
						Computed:            true,
						Attributes:          hostFlagNestedAttr,
						Description:         "When setting this flag, the port must be offline.",
						MarkdownDescription: "When setting this flag, the port must be offline.",
						PlanModifiers:       hostFlagDefaultPlanModifier,
					},
					"scsi_support1": {
						Optional:            true,
						Computed:            true,
						Attributes:          hostFlagNestedAttr,
						Description:         "This attribute provides a stricter compliance with SCSI standards.",
						MarkdownDescription: "This attribute provides a stricter compliance with SCSI standards.",
						PlanModifiers:       hostFlagDefaultPlanModifier,
					},
				}),
				Description:         "Host Flags set for the hostgroup. When host_flags = {} then default flags will be considered.",
				MarkdownDescription: "Host Flags set for the hostgroup. When host_flags = {} then default flags will be considered.",
			},
			"consistent_lun": {
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				Description:         "It enables the rejection of any masking operation involving this hostgroup that would result in inconsistent LUN values.",
				MarkdownDescription: "It enables the rejection of any masking operation involving this hostgroup that would result in inconsistent LUN values.",
				PlanModifiers:       boolDefaultPlanModifier,
			},
			"numofmaskingviews": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "The number of masking views associated with the hostgroup.",
				MarkdownDescription: "The number of masking views associated with the hostgroup.",
			},
			"numofinitiators": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "The number of initiators associated with the hostgroup.",
				MarkdownDescription: "The number of initiators associated with the hostgroup.",
			},
			"numofhosts": {
				Type:                types.Int64Type,
				Computed:            true,
				Description:         "The number of hosts associated with the hostgroup.",
				MarkdownDescription: "The number of hosts associated with the hostgroup.",
			},
			"type": {
				Type:                types.StringType,
				Computed:            true,
				Description:         "Specifies the type of hostgroup.",
				MarkdownDescription: "Specifies the type of hostgroup.",
			},
			"maskingviews": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Computed:            true,
				Description:         "The masking views associated with the hostgroup.",
				MarkdownDescription: "The masking views associated with the hostgroup.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"port_flags_override": {
				Type:                types.BoolType,
				Computed:            true,
				Description:         "States whether port flags override is enabled on the hostgroup.",
				MarkdownDescription: "States whether port flags override is enabled on the hostgroup.",
			},
		},
	}, nil
}

// NewResource is a wrapper around provider
func (r resourceHostGroupType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceHostGroup{
		p: *(p.(*provider)),
	}, nil
}

type resourceHostGroup struct {
	p provider
}

// Create Host
func (r resourceHostGroup) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Info(ctx, "creating hostGroup")
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	tflog.Debug(ctx, "setting hostgroup plan")
	var planHostGroup models.HostGroup
	diags := req.Plan.Get(ctx, &planHostGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hostIDs := make([]string, len(planHostGroup.HostIDs.Elems))

	diags = planHostGroup.HostIDs.ElementsAs(ctx, &hostIDs, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "preparing hostgroup host flags")
	hostFlags := pmaxTypes.HostFlags{
		VolumeSetAddressing: &pmaxTypes.HostFlag{
			Enabled:  planHostGroup.HostFlags.VolumeSetAddressing.Enabled.Value,
			Override: planHostGroup.HostFlags.VolumeSetAddressing.Override.Value,
		},
		DisableQResetOnUA: &pmaxTypes.HostFlag{
			Enabled:  planHostGroup.HostFlags.DisableQResetOnUa.Enabled.Value,
			Override: planHostGroup.HostFlags.DisableQResetOnUa.Override.Value,
		},
		EnvironSet: &pmaxTypes.HostFlag{
			Enabled:  planHostGroup.HostFlags.EnvironSet.Enabled.Value,
			Override: planHostGroup.HostFlags.EnvironSet.Override.Value,
		},
		AvoidResetBroadcast: &pmaxTypes.HostFlag{
			Enabled:  planHostGroup.HostFlags.AvoidResetBroadcast.Enabled.Value,
			Override: planHostGroup.HostFlags.AvoidResetBroadcast.Override.Value,
		},
		OpenVMS: &pmaxTypes.HostFlag{
			Enabled:  planHostGroup.HostFlags.Openvms.Enabled.Value,
			Override: planHostGroup.HostFlags.Openvms.Override.Value,
		},
		SCSI3: &pmaxTypes.HostFlag{
			Enabled:  planHostGroup.HostFlags.Scsi3.Enabled.Value,
			Override: planHostGroup.HostFlags.Scsi3.Override.Value,
		},
		Spc2ProtocolVersion: &pmaxTypes.HostFlag{
			Enabled:  planHostGroup.HostFlags.Spc2ProtocolVersion.Enabled.Value,
			Override: planHostGroup.HostFlags.Spc2ProtocolVersion.Override.Value,
		},
		SCSISupport1: &pmaxTypes.HostFlag{
			Enabled:  planHostGroup.HostFlags.ScsiSupport1.Enabled.Value,
			Override: planHostGroup.HostFlags.ScsiSupport1.Override.Value,
		},
		ConsistentLUN: planHostGroup.ConsistentLun.Value,
	}

	tflog.Info(ctx, "calling create hostgroup pmax client", map[string]interface{}{
		"symmetrixID": r.p.client.SymmetrixID,
		"host":        planHostGroup.Name.Value,
		"hostIDS":     hostIDs,
		"hostFlags":   hostFlags,
	})

	hostGroupResponse, err := r.p.client.PmaxClient.CreateHostGroup(ctx, r.p.client.SymmetrixID, planHostGroup.Name.Value, hostIDs, &hostFlags)
	if err != nil {
		hostgroupID := planHostGroup.Name.Value
		resp.Diagnostics.AddError(
			"Error creating hostgroup",
			CreateHostGroupDetailErrorMsg+hostgroupID+"with error: "+err.Error(),
		)
		hostGroupResponse, getHostGroupErr := r.p.client.PmaxClient.GetHostGroupByID(ctx, r.p.client.SymmetrixID, hostgroupID)
		if hostGroupResponse != nil || getHostGroupErr == nil {
			err := r.p.client.PmaxClient.DeleteHost(ctx, r.p.client.SymmetrixID, hostgroupID)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error deleting the invalid hostGroup, This may be a dangling resource and needs to be deleted manually",
					CreateHostGroupDetailErrorMsg+hostgroupID+"with error: "+err.Error(),
				)
			}
		}
		return
	}
	tflog.Debug(ctx, "create hostgroup response", map[string]interface{}{
		"Create HostGroup Response": hostGroupResponse,
	})

	hostGroupState := models.HostGroup{}
	tflog.Debug(ctx, "updating hostgroup state", map[string]interface{}{
		"host state":        hostGroupState,
		"hostIDs":           hostIDs,
		"hostGroupResponse": hostGroupResponse,
	})
	updateHostGroupState(&hostGroupState, hostGroupResponse)
	diags = resp.State.Set(ctx, hostGroupState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create hostGroup completed")
}

// Read Host
func (r resourceHostGroup) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	tflog.Info(ctx, "reading hostgroup")
	var hostGroupState models.HostGroup
	diags := req.State.Get(ctx, &hostGroupState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hostGroupID := hostGroupState.ID.Value
	tflog.Debug(ctx, "fetching hostgroup by ID", map[string]interface{}{
		"symmetrixID": r.p.client.SymmetrixID,
		"hostGroupID": hostGroupID,
	})
	hostGroupResponse, err := r.p.client.PmaxClient.GetHostGroupByID(ctx, r.p.client.SymmetrixID, hostGroupID)
	fmt.Println("Hostgroup response ", hostGroupResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading hostGroup",
			ReadHostGroupDetailsErrorMsg+hostGroupID+" with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Get HostGroup By ID response", map[string]interface{}{
		"HostGroup Response": hostGroupResponse,
	})

	tflog.Debug(ctx, "updating hostgroup state")
	updateHostGroupState(&hostGroupState, hostGroupResponse)
	diags = resp.State.Set(ctx, hostGroupState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read hostgroup completed")
}

// Update HostGroup
// Supported updates: name, host_ids, host_flags
func (r resourceHostGroup) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Info(ctx, "updating host")
	var planHostGroup models.HostGroup
	diags := req.Plan.Get(ctx, &planHostGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "fetched hostgroup details from plan")

	var stateHostGroup models.HostGroup
	diags = req.State.Get(ctx, &stateHostGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update hostgroup on pmax client", map[string]interface{}{
		"plan":  planHostGroup,
		"state": stateHostGroup,
	})
	updatedParams, updateFailedParameters, errMessages := updateHostGroup(ctx, r.p.client, planHostGroup, stateHostGroup)
	if len(errMessages) > 0 || len(updateFailedParameters) > 0 {
		errMessage := strings.Join(errMessages, ",\n")
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s, updated parameters are %v and parameters failed to update are %v", UpdateHostGroupDetailsErrorMsg, updatedParams, updateFailedParameters),
			errMessage)
	}
	tflog.Debug(ctx, "update hostgroup response", map[string]interface{}{
		"updatedParams":          updatedParams,
		"updateFailedParameters": updateFailedParameters,
		"error messages":         errMessages,
	})

	hostGroupID := stateHostGroup.ID.Value
	if isParamUpdated(updatedParams, "name") {
		hostGroupID = planHostGroup.Name.Value
	}

	tflog.Debug(ctx, "calling get hostgroup by ID on pmax client", map[string]interface{}{
		"SymmetrixID": r.p.client.SymmetrixID,
		"hostgroupID": hostGroupID,
	})
	hostGroupResponse, err := r.p.client.PmaxClient.GetHostGroupByID(ctx, r.p.client.SymmetrixID, hostGroupID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading hostgroup",
			ReadHostGroupDetailsErrorMsg+hostGroupID+" with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "get hostgroup by ID response", map[string]interface{}{
		"HostGroup Response": hostGroupResponse,
	})

	tflog.Debug(ctx, "updating hostgroup state", map[string]interface{}{
		"state":             stateHostGroup,
		"hostGroupResponse": hostGroupResponse,
	})
	updateHostGroupState(&stateHostGroup, hostGroupResponse)
	diags = resp.State.Set(ctx, stateHostGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update hostgroup completed")
}

// Delete Host
func (r resourceHostGroup) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Info(ctx, "deleting hostgroup")
	var hostGroupState models.HostGroup
	diags := req.State.Get(ctx, &hostGroupState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	hostGroupID := hostGroupState.ID.Value
	tflog.Debug(ctx, "deleting hostgroup by hostgroup ID", map[string]interface{}{
		"symmetrixID": r.p.client.SymmetrixID,
		"hostGroupID": hostGroupID,
	})
	err := r.p.client.PmaxClient.DeleteHostGroup(ctx, r.p.client.SymmetrixID, hostGroupID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting hostGroup",
			DeleteHostGroupDetailsErrorMsg+hostGroupID+" with error: "+err.Error(),
		)
	}

	tflog.Info(ctx, "delete hostgroup complete")
}

func (r resourceHostGroup) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tflog.Info(ctx, "importing hostgroup state")
	var hostGroupState models.HostGroup
	hostGroupID := req.ID
	tflog.Debug(ctx, "fetching hostgroup by ID", map[string]interface{}{
		"symmetrixID": r.p.client.SymmetrixID,
		"hostID":      hostGroupID,
	})
	hostGroupResponse, err := r.p.client.PmaxClient.GetHostGroupByID(ctx, r.p.client.SymmetrixID, hostGroupID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading hostgroup",
			ImportHostGroupDetailsErrorMsg+hostGroupID+" with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Get HostGroup By ID response", map[string]interface{}{
		"HostGroup Response": hostGroupResponse,
	})

	tflog.Debug(ctx, "updating hostgroup state after import")
	updateHostGroupState(&hostGroupState, hostGroupResponse)
	diags := resp.State.Set(ctx, hostGroupState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "import hostgroup state completed")

}
