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

	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure implementation
var (
	_ resource.Resource                = &HostGroup{}
	_ resource.ResourceWithConfigure   = &HostGroup{}
	_ resource.ResourceWithImportState = &HostGroup{}
)

func NewHostGroup() resource.Resource {
	return &HostGroup{}
}

type HostGroup struct {
	client *client.Client
}

func (r *HostGroup) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hostgroup"
}

func (r *HostGroup) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	hostFlagNestedAttr := map[string]schema.Attribute{
		"override": schema.BoolAttribute{
			Optional: true,
			Computed: true,
			Default:  booldefault.StaticBool(false),
		},
		"enabled": schema.BoolAttribute{
			Optional: true,
			Computed: true,
			Default:  booldefault.StaticBool(false),
		},
	}
	objd, err := basetypes.NewObjectValue(
		map[string]attr.Type{
			"override": types.BoolType,
			"enabled":  types.BoolType,
		},
		map[string]attr.Value{
			"override": types.BoolValue(false),
			"enabled":  types.BoolValue(false),
		},
	)

	if err != nil {
		return
	}

	hostFlagAttr := map[string]schema.Attribute{
		"volume_set_addressing": schema.SingleNestedAttribute{
			Optional:            true,
			Computed:            true,
			Attributes:          hostFlagNestedAttr,
			Description:         "It enables the volume set addressing mode.",
			MarkdownDescription: "It enables the volume set addressing mode.",
			Default:             objectdefault.StaticValue(objd),
		},
		"disable_q_reset_on_ua": schema.SingleNestedAttribute{
			Optional:            true,
			Computed:            true,
			Attributes:          hostFlagNestedAttr,
			Description:         "It is used for hosts that do not expect the queue to be flushed on a 0629 sense.",
			MarkdownDescription: "It is used for hosts that do not expect the queue to be flushed on a 0629 sense.",
			Default:             objectdefault.StaticValue(objd),
		},
		"environ_set": schema.SingleNestedAttribute{
			Optional:            true,
			Computed:            true,
			Attributes:          hostFlagNestedAttr,
			Description:         "It enables the environmental error reporting by the storage system to the host on the specific port.",
			MarkdownDescription: "It enables the environmental error reporting by the storage system to the host on the specific port.",
			Default:             objectdefault.StaticValue(objd),
		},
		"openvms": schema.SingleNestedAttribute{
			Optional:            true,
			Computed:            true,
			Attributes:          hostFlagNestedAttr,
			Description:         "This attribute enables an Open VMS fibre connection.",
			MarkdownDescription: "This attribute enables an Open VMS fibre connection.",
			Default:             objectdefault.StaticValue(objd),
		},
		"avoid_reset_broadcast": schema.SingleNestedAttribute{
			Optional:            true,
			Computed:            true,
			Attributes:          hostFlagNestedAttr,
			Description:         "It enables a SCSI bus reset to only occur to the port that received the reset.",
			MarkdownDescription: "It enables a SCSI bus reset to only occur to the port that received the reset.",
			Default:             objectdefault.StaticValue(objd),
		},
		"scsi_3": schema.SingleNestedAttribute{
			Optional:            true,
			Computed:            true,
			Attributes:          hostFlagNestedAttr,
			Description:         "Alters the inquiry data to report that the storage system supports the SCSI-3 protocol.",
			MarkdownDescription: "Alters the inquiry data to report that the storage system supports the SCSI-3 protocol.",
			Default:             objectdefault.StaticValue(objd),
		},
		"spc2_protocol_version": schema.SingleNestedAttribute{
			Optional:            true,
			Computed:            true,
			Attributes:          hostFlagNestedAttr,
			Description:         "When setting this flag, the port must be offline.",
			MarkdownDescription: "When setting this flag, the port must be offline.",
			Default:             objectdefault.StaticValue(objd),
		},
		"scsi_support1": schema.SingleNestedAttribute{
			Optional:            true,
			Computed:            true,
			Attributes:          hostFlagNestedAttr,
			Description:         "This attribute provides a stricter compliance with SCSI standards.",
			MarkdownDescription: "This attribute provides a stricter compliance with SCSI standards.",
			Default:             objectdefault.StaticValue(objd),
		},
	}

	resp.Schema = schema.Schema{
		// Description for Docs
		MarkdownDescription: "Resource for managing HostGroups for a PowerMax Array",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the hostgroup.",
				MarkdownDescription: "The ID of the hostgroup.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the hostgroup.",
				MarkdownDescription: "The name of the hostgroup.",
			},
			"host_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Required:            true,
				Description:         "A list of host ids associated with the hostgroup.",
				MarkdownDescription: "The masking views associated with the hostgroup.",
			},
			"host_flags": schema.SingleNestedAttribute{
				Description:         "Host Flags set for the hostgroup. When host_flags = {} then default flags will be considered.",
				MarkdownDescription: "Host Flags set for the hostgroup. When host_flags = {} then default flags will be considered.",
				Required:            true,
				Attributes:          hostFlagAttr,
			},
			"consistent_lun": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				Description:         "It enables the rejection of any masking operation involving this hostgroup that would result in inconsistent LUN values.",
				MarkdownDescription: "It enables the rejection of any masking operation involving this hostgroup that would result in inconsistent LUN values.",
			},
			"numofmaskingviews": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of masking views associated with the hostgroup.",
				MarkdownDescription: "The number of masking views associated with the hostgroup.",
			},
			"numofinitiators": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of initiators associated with the hostgroup.",
				MarkdownDescription: "The number of initiators associated with the hostgroup.",
			},
			"numofhosts": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of hosts associated with the hostgroup.",
				MarkdownDescription: "The number of hosts associated with the hostgroup.",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				Description:         "Specifies the type of hostgroup.",
				MarkdownDescription: "Specifies the type of hostgroup.",
			},
			"maskingviews": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Description:         "The masking views associated with the hostgroup.",
				MarkdownDescription: "The masking views associated with the hostgroup.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"port_flags_override": schema.BoolAttribute{
				Computed:            true,
				Description:         "States whether port flags override is enabled on the hostgroup.",
				MarkdownDescription: "States whether port flags override is enabled on the hostgroup.",
			},
		},
	}
}

func (r *HostGroup) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if provider is not config
	if req.ProviderData == nil {
		return
	}

	client, err := req.ProviderData.(*client.Client)

	if !err {
		resp.Diagnostics.AddError(
			"Unexpected Resource Config Failure",
			fmt.Sprintf("Expected client, %T. Please report this issue to the provider developers", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *HostGroup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Create Host Group")
	var plan models.HostGroupModal
	var state models.HostGroupModal
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	hostIds := make([]string, len(plan.HostIDs.Elements()))
	diags := plan.HostIDs.ElementsAs(ctx, &hostIds, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hostFlags := pmaxTypes.HostFlags{
		VolumeSetAddressing: &pmaxTypes.HostFlag{
			Enabled:  plan.HostFlags.VolumeSetAddressing.Enabled.ValueBool(),
			Override: plan.HostFlags.VolumeSetAddressing.Override.ValueBool(),
		},
		DisableQResetOnUA: &pmaxTypes.HostFlag{
			Enabled:  plan.HostFlags.DisableQResetOnUA.Enabled.ValueBool(),
			Override: plan.HostFlags.DisableQResetOnUA.Override.ValueBool(),
		},
		EnvironSet: &pmaxTypes.HostFlag{
			Enabled:  plan.HostFlags.EnvironSet.Enabled.ValueBool(),
			Override: plan.HostFlags.EnvironSet.Override.ValueBool(),
		},
		AvoidResetBroadcast: &pmaxTypes.HostFlag{
			Enabled:  plan.HostFlags.AvoidResetBroadcast.Enabled.ValueBool(),
			Override: plan.HostFlags.AvoidResetBroadcast.Override.ValueBool(),
		},
		OpenVMS: &pmaxTypes.HostFlag{
			Enabled:  plan.HostFlags.OpenVMS.Enabled.ValueBool(),
			Override: plan.HostFlags.OpenVMS.Override.ValueBool(),
		},
		SCSI3: &pmaxTypes.HostFlag{
			Enabled:  plan.HostFlags.SCSI3.Enabled.ValueBool(),
			Override: plan.HostFlags.SCSI3.Override.ValueBool(),
		},
		Spc2ProtocolVersion: &pmaxTypes.HostFlag{
			Enabled:  plan.HostFlags.Spc2ProtocolVersion.Enabled.ValueBool(),
			Override: plan.HostFlags.Spc2ProtocolVersion.Override.ValueBool(),
		},
		SCSISupport1: &pmaxTypes.HostFlag{
			Enabled:  plan.HostFlags.SCSISupport1.Enabled.ValueBool(),
			Override: plan.HostFlags.SCSISupport1.Override.ValueBool(),
		},
		ConsistentLUN: plan.ConsistentLun.ValueBool(),
	}
	tflog.Info(ctx, "calling create hostgroup with client", map[string]interface{}{
		"symmetrixID": r.client.SymmetrixID,
		"host":        plan.Name.ValueString(),
		"hostIds":     hostIds,
		"hostFlags":   hostFlags,
	})

	newHostGroup, err := r.client.PmaxClient.CreateHostGroup(ctx, r.client.SymmetrixID, plan.Name.ValueString(), hostIds, &hostFlags)
	if err != nil {
		hostgroupID := plan.Name.ValueString()
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create host group, got error: %s", err.Error()))
		//Attempt to remove any partially created obejcts if there are any
		hostGroupResponse, getHostGroupErr := r.client.PmaxClient.GetHostGroupByID(ctx, r.client.SymmetrixID, hostgroupID)
		if hostGroupResponse != nil || getHostGroupErr == nil {
			err := r.client.PmaxClient.DeleteHost(ctx, r.client.SymmetrixID, hostgroupID)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error deleting the invalid hostGroup, This may be a dangling resource and needs to be deleted manually",
					constants.CreateHostGroupDetailErrorMsg+hostgroupID+"with error: "+err.Error(),
				)
			}
		}
		return
	}
	tflog.Debug(ctx, "created a resource", map[string]interface{}{
		"host group": newHostGroup,
	})

	tflog.Debug(ctx, "updating hostgroup state", map[string]interface{}{
		"host state":   state,
		"hostIDs":      hostIds,
		"newHostGroup": newHostGroup,
	})
	helper.UpdateHostGroupState(&state, newHostGroup)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *HostGroup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Host Group")
	var state models.HostGroupModal
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hostGroupId := state.ID.ValueString()
	tflog.Debug(ctx, "fetching hostgroup by ID", map[string]interface{}{
		"symmetricxId": r.client.SymmetrixID,
		"hostGroupId":  hostGroupId,
	})

	hgResponse, err := r.client.PmaxClient.GetHostGroupByID(ctx, r.client.SymmetrixID, hostGroupId)
	tflog.Debug(ctx, "Get HostGroup By ID response", map[string]interface{}{
		"HostGroup Response": hgResponse,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading hostGroup",
			constants.ReadHostGroupDetailsErrorMsg+hostGroupId+" with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Updating Hostgroup State")
	helper.UpdateHostGroupState(&state, hgResponse)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Read HostGroup Completed")
}

// Update HostGroup
// Supported updates: name, host_ids, host_flags
func (r *HostGroup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating HostGroup")
	var planHostGroup models.HostGroupModal
	diags := req.Plan.Get(ctx, &planHostGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "fetched hostgroup details from plan")

	var stateHostGroup models.HostGroupModal
	diags = req.State.Get(ctx, &stateHostGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update hostgroup on pmax client", map[string]interface{}{
		"plan":  planHostGroup,
		"state": stateHostGroup,
	})
	updatedParams, updateFailedParameters, errMessages := helper.UpdateHostGroup(ctx, *r.client, planHostGroup, stateHostGroup)
	if len(errMessages) > 0 || len(updateFailedParameters) > 0 {
		errMessage := strings.Join(errMessages, ",\n")
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s, updated parameters are %v and parameters failed to update are %v", constants.UpdateHostGroupDetailsErrorMsg, updatedParams, updateFailedParameters),
			errMessage)
	}
	tflog.Debug(ctx, "update hostgroup response", map[string]interface{}{
		"updatedParams":          updatedParams,
		"updateFailedParameters": updateFailedParameters,
		"error messages":         errMessages,
	})

	hostGroupID := stateHostGroup.ID.ValueString()
	if helper.IsParamUpdated(updatedParams, "name") {
		hostGroupID = planHostGroup.Name.ValueString()
	}

	tflog.Debug(ctx, "calling get hostgroup by ID on pmax client", map[string]interface{}{
		"SymmetrixID": r.client.SymmetrixID,
		"hostgroupID": hostGroupID,
	})
	hostGroupResponse, err := r.client.PmaxClient.GetHostGroupByID(ctx, r.client.SymmetrixID, hostGroupID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading hostgroup",
			constants.ReadHostGroupDetailsErrorMsg+hostGroupID+" with error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "updating hostgroup state", map[string]interface{}{
		"state":             stateHostGroup,
		"hostGroupResponse": hostGroupResponse,
	})
	helper.UpdateHostGroupState(&stateHostGroup, hostGroupResponse)
	diags = resp.State.Set(ctx, stateHostGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Update Hostgroup Completed")
}

func (r *HostGroup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting hostgroup")
	var hostGroupState models.HostGroupModal
	diags := req.State.Get(ctx, &hostGroupState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	hostGroupID := hostGroupState.ID.ValueString()
	tflog.Debug(ctx, "deleting hostgroup by hostgroup ID", map[string]interface{}{
		"symmetrixID": r.client.SymmetrixID,
		"hostGroupID": hostGroupID,
	})
	err := r.client.PmaxClient.DeleteHostGroup(ctx, r.client.SymmetrixID, hostGroupID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting hostGroup",
			constants.DeleteHostGroupDetailsErrorMsg+hostGroupID+" with error: "+err.Error(),
		)
	}

	tflog.Info(ctx, "delete hostgroup complete")
}

func (r *HostGroup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing Hostgroup State")
	var hostGroupState models.HostGroupModal
	hostGroupID := req.ID
	tflog.Debug(ctx, "fetching Hostgroup by ID", map[string]interface{}{
		"symmetrixID": r.client.SymmetrixID,
		"hostID":      hostGroupID,
	})
	hostGroupResponse, err := r.client.PmaxClient.GetHostGroupByID(ctx, r.client.SymmetrixID, hostGroupID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading hostgroup",
			constants.ImportHostGroupDetailsErrorMsg+hostGroupID+" with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Get HostGroup By ID response", map[string]interface{}{
		"HostGroup Response": hostGroupResponse,
	})

	tflog.Debug(ctx, "updating hostgroup state after import")
	helper.UpdateHostGroupState(&hostGroupState, hostGroupResponse)
	diags := resp.State.Set(ctx, hostGroupState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Import Hostgroup State Completed")
}
