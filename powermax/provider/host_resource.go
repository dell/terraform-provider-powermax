/*
Copyright (c) 2022-2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"strings"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/constants"
	"terraform-provider-powermax/powermax/helper"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &Host{}
var _ resource.ResourceWithImportState = &Host{}
var _ resource.ResourceWithConfigure = &Host{}

// NewHost creates a new Host resource.
func NewHost() resource.Resource {
	return &Host{}
}

// Host defines the resource implementation.
type Host struct {
	client *client.Client
}

// Metadata returns the metadata for the resource.
func (r *Host) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

// Schema returns the schema for the resource.
func (r *Host) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	hostFlagNestedAttr := map[string]schema.Attribute{
		"override": schema.BoolAttribute{
			Default:       booldefault.StaticBool(false),
			Optional:      true,
			Computed:      true,
			PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
		},
		"enabled": schema.BoolAttribute{
			Default:       booldefault.StaticBool(false),
			Optional:      true,
			Computed:      true,
			PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
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

	hostDefaultObj, _ := basetypes.NewObjectValue(
		map[string]attr.Type{
			"volume_set_addressing": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"override": types.BoolType,
					"enabled":  types.BoolType,
				},
			},
			"disable_q_reset_on_ua": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"override": types.BoolType,
					"enabled":  types.BoolType,
				},
			},
			"environ_set": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"override": types.BoolType,
					"enabled":  types.BoolType,
				},
			},
			"openvms": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"override": types.BoolType,
					"enabled":  types.BoolType,
				},
			},
			"avoid_reset_broadcast": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"override": types.BoolType,
					"enabled":  types.BoolType,
				},
			},
			"scsi_3": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"override": types.BoolType,
					"enabled":  types.BoolType,
				},
			},
			"spc2_protocol_version": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"override": types.BoolType,
					"enabled":  types.BoolType,
				},
			},
			"scsi_support1": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"override": types.BoolType,
					"enabled":  types.BoolType,
				},
			},
		},
		map[string]attr.Value{
			"volume_set_addressing": objd,
			"disable_q_reset_on_ua": objd,
			"environ_set":           objd,
			"openvms":               objd,
			"avoid_reset_broadcast": objd,
			"scsi_3":                objd,
			"spc2_protocol_version": objd,
			"scsi_support1":         objd,
		},
	)

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

		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Resource for managing Host in PowerMax array.",
		Description:         "Resource for managing Host in PowerMax array.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the host.",
				MarkdownDescription: "The ID of the host.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the host. Only alphanumeric characters, underscores ( _ ), and hyphens (-) are allowed.",
				MarkdownDescription: "The name of the host. Only alphanumeric characters, underscores ( _ ), and hyphens (-) are allowed.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(64),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9_-]*$`),
						"must contain only alphanumeric characters and _-",
					),
				},
			},
			"num_of_masking_views": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of masking views associated with the host.",
				MarkdownDescription: "The number of masking views associated with the host.",
			},
			"num_of_initiators": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of initiators associated with the host.",
				MarkdownDescription: "The number of initiators associated with the host.",
			},
			"num_of_host_groups": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of hostgroups associated with the host.",
				MarkdownDescription: "The number of hostgroups associated with the host.",
			},
			"port_flags_override": schema.BoolAttribute{
				Computed:            true,
				Description:         "States whether port flags override is enabled on the host.",
				MarkdownDescription: "States whether port flags override is enabled on the host.",
			},
			"consistent_lun": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "It enables the rejection of any masking operation involving this host that would result in inconsistent LUN values.",
				MarkdownDescription: "It enables the rejection of any masking operation involving this host that would result in inconsistent LUN values.",
				Default:             booldefault.StaticBool(false),
			},
			"type": schema.StringAttribute{
				Computed:            true,
				Description:         "Specifies the type of host.",
				MarkdownDescription: "Specifies the type of host.",
			},
			"initiator": schema.ListAttribute{
				ElementType:         types.StringType,
				Required:            true,
				Description:         "The initiators associated with the host.",
				MarkdownDescription: "The initiators associated with the host.",
			},
			"hostgroup": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Description:         "The host group associated with the host.",
				MarkdownDescription: "The host group associated with the host.",
			},
			"maskingview": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Description:         "The masking views associated with the host.",
				MarkdownDescription: "The masking views associated with the host.",
				PlanModifiers:       []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
			"powerpathhosts": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Description:         "The powerpath hosts associated with the host.",
				MarkdownDescription: "The powerpath hosts associated with the host.",
			},
			"numofpowerpathhosts": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of powerpath hosts associated with the host.",
				MarkdownDescription: "The number of powerpath hosts associated with the host.",
			},
			"bw_limit": schema.Int64Attribute{
				Computed:            true,
				Description:         "Specifies the bandwidth limit for a host.",
				MarkdownDescription: "Specifies the bandwidth limit for a host.",
			},
			"host_flags": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				Default:             objectdefault.StaticValue(hostDefaultObj),
				Attributes:          hostFlagAttr,
				Description:         "Flags set for the host. When host_flags = {} then default flags will be considered.",
				MarkdownDescription: "Flags set for the host. When host_flags = {} then default flags will be considered.",
			},
		},
	}
}

// Configure configure client for host resource.
func (r *Host) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a host and refresh state.
func (r *Host) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Host...")
	var planHost models.HostModel
	diags := req.Plan.Get(ctx, &planHost)
	// Read Terraform plan into the model
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	initiators := make([]string, len(planHost.Initiators.Elements()))
	if len(planHost.Initiators.Elements()) > 0 {
		for index, initiator := range planHost.Initiators.Elements() {
			initiators[index] = strings.Trim(initiator.String(), "\"")
		}
	}
	tflog.Debug(ctx, "preparing host flags")
	hostFlags := *pmax.NewHostFlags(
		*pmax.NewVolumeSetAddressing(planHost.HostFlags.VolumeSetAddressing.Enabled.ValueBool(), planHost.HostFlags.VolumeSetAddressing.Override.ValueBool()),
		*pmax.NewDisableQResetOnUa(planHost.HostFlags.DisableQResetOnUA.Enabled.ValueBool(), planHost.HostFlags.DisableQResetOnUA.Override.ValueBool()),
		*pmax.NewEnvironSet(planHost.HostFlags.EnvironSet.Enabled.ValueBool(), planHost.HostFlags.EnvironSet.Override.ValueBool()),
		*pmax.NewAvoidResetBroadcast(planHost.HostFlags.AvoidResetBroadcast.Enabled.ValueBool(), planHost.HostFlags.AvoidResetBroadcast.Override.ValueBool()),
		*pmax.NewOpenvms(planHost.HostFlags.OpenVMS.Enabled.ValueBool(), planHost.HostFlags.OpenVMS.Override.ValueBool()),
		*pmax.NewScsi3(planHost.HostFlags.SCSI3.Enabled.ValueBool(), planHost.HostFlags.SCSI3.Override.ValueBool()),
		*pmax.NewSpc2ProtocolVersion(planHost.HostFlags.Spc2ProtocolVersion.Enabled.ValueBool(), planHost.HostFlags.Spc2ProtocolVersion.Override.ValueBool()),
		*pmax.NewScsiSupport1(planHost.HostFlags.SCSISupport1.Enabled.ValueBool(), planHost.HostFlags.SCSISupport1.Override.ValueBool()),
		planHost.ConsistentLun.ValueBool(),
	)

	hostCreateReq := r.client.PmaxOpenapiClient.SLOProvisioningApi.CreateHost(ctx, r.client.SymmetrixID)
	createHostParam := pmax.NewCreateHostParam(planHost.Name.ValueString())
	createHostParam.SetHostFlags(hostFlags)
	createHostParam.SetInitiatorId(initiators)
	hostCreateReq = hostCreateReq.CreateHostParam(*createHostParam)
	hostCreateResp, _, err := r.client.PmaxOpenapiClient.SLOProvisioningApi.CreateHostExecute(hostCreateReq)
	if err != nil {
		hostID := planHost.Name.ValueString()

		errStr := constants.CreateHostDetailErrorMsg + hostID + ": "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating host",
			message,
		)

		req := r.client.PmaxOpenapiClient.SLOProvisioningApi.GetHost(ctx, r.client.SymmetrixID, hostID)
		hostGetResp, _, getHostErr := req.Execute()
		if hostGetResp != nil || getHostErr == nil {
			delReq := r.client.PmaxOpenapiClient.SLOProvisioningApi.DeleteHost(ctx, r.client.SymmetrixID, hostID)
			_, err := delReq.Execute()
			if err != nil {
				errStr := constants.CreateHostDetailErrorMsg + hostID + "with error: "
				message := helper.GetErrorString(err, errStr)
				resp.Diagnostics.AddError(
					"Error deleting the invalid host, This may be a dangling resource and needs to be deleted manually",
					message,
				)
			}
		}
		return
	}
	tflog.Debug(ctx, "create host response", map[string]interface{}{
		"Create Host Response": hostCreateResp,
	})
	result := models.HostModel{}
	helper.UpdateHostState(&result, initiators, hostCreateResp)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create host completed")

}

// Delete Host.
func (r *Host) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting Host")
	var hostState models.HostModel
	diags := req.State.Get(ctx, &hostState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	hostID := hostState.HostID.ValueString()
	tflog.Debug(ctx, "deleting host by host ID", map[string]interface{}{
		"symmetrixID": r.client.SymmetrixID,
		"hostID":      hostID,
	})
	delReq := r.client.PmaxOpenapiClient.SLOProvisioningApi.DeleteHost(ctx, r.client.SymmetrixID, hostID)
	_, err := delReq.Execute()
	if err != nil {
		errStr := constants.DeleteHostDetailsErrorMsg + hostID + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error deleting host",
			message,
		)
	}

	tflog.Info(ctx, "Delete host complete")
}

// Update Host.
func (r *Host) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "updating host")
	var plan models.HostModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "fetched host details from plan")

	var state models.HostModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update host on pmax client", map[string]interface{}{
		"plan":  plan,
		"state": state,
	})
	updatedParams, updateFailedParameters, errMessages := helper.UpdateHost(ctx, *r.client, plan, state)
	if len(errMessages) > 0 || len(updateFailedParameters) > 0 {
		errMessage := strings.Join(errMessages, ",\n")
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s, updated parameters are %v and parameters failed to update are %v", constants.UpdateHostDetailsErrorMsg, updatedParams, updateFailedParameters),
			errMessage)
	}
	tflog.Debug(ctx, "update host response", map[string]interface{}{
		"updatedParams":          updatedParams,
		"updateFailedParameters": updateFailedParameters,
		"error messages":         errMessages,
	})

	hostID := state.HostID.ValueString()
	if helper.IsParamUpdated(updatedParams, "name") {
		hostID = plan.Name.ValueString()
	}
	getReq := r.client.PmaxOpenapiClient.SLOProvisioningApi.GetHost(ctx, r.client.SymmetrixID, hostID)
	hostResponse, _, err := getReq.Execute()
	if err != nil {
		errStr := constants.ReadHostDetailsErrorMsg + hostID + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading host",
			message,
		)
		return
	}
	tflog.Debug(ctx, "get host by ID response", map[string]interface{}{
		"Host Response": hostResponse,
	})

	initiators := make([]string, len(plan.Initiators.Elements()))
	if len(plan.Initiators.Elements()) > 0 {
		for index, initiator := range plan.Initiators.Elements() {
			initiatorVal := strings.TrimSpace(strings.Trim(initiator.String(), "\""))
			if initiatorVal == "" {
				resp.Diagnostics.AddError(
					"Error updating host",
					"Empty initiator values are not allowed",
				)
				return
			}
			initiators[index] = initiatorVal
		}
	}

	tflog.Debug(ctx, "updating update host state", map[string]interface{}{
		"state":        state,
		"initiators":   initiators,
		"hostResponse": hostResponse,
	})
	helper.UpdateHostState(&state, initiators, hostResponse)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update host completed")
}

// Read Host.
func (r *Host) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Host...")
	var hostState models.HostModel
	diags := req.State.Get(ctx, &hostState)
	// Read Terraform prior state into the model
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	hostID := hostState.HostID.ValueString()
	getReq := r.client.PmaxOpenapiClient.SLOProvisioningApi.GetHost(ctx, r.client.SymmetrixID, hostID)
	host, _, err := getReq.Execute()
	if err != nil {
		errStr := constants.ReadHostDetailsErrorMsg + hostID + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading host",
			message,
		)
		return
	}
	initiators := make([]string, len(hostState.Initiators.Elements()))

	if len(hostState.Initiators.Elements()) > 0 {
		for index, initiator := range hostState.Initiators.Elements() {
			initiators[index] = strings.Trim(initiator.String(), "\"")
		}
	}

	tflog.Debug(ctx, "Updating host state")
	helper.UpdateHostState(&hostState, initiators, host)
	diags = resp.State.Set(ctx, hostState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Read host completed")

}

// ImportState imports the state of the resource from the req.
func (r *Host) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "importing host state")
	var hostState models.HostModel
	hostID := req.ID
	tflog.Debug(ctx, "fetching host by ID", map[string]interface{}{
		"symmetrixID": r.client.SymmetrixID,
		"hostID":      hostID,
	})

	getReq := r.client.PmaxOpenapiClient.SLOProvisioningApi.GetHost(ctx, r.client.SymmetrixID, hostID)
	hostResponse, _, err := getReq.Execute()

	if err != nil {
		errStr := constants.ImportHostDetailsErrorMsg + hostID + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading host",
			message,
		)
		return
	}
	tflog.Debug(ctx, "Get Host By ID response", map[string]interface{}{
		"Host Response": hostResponse,
	})

	tflog.Debug(ctx, "updating host state after import")
	helper.UpdateHostState(&hostState, hostResponse.Initiator, hostResponse)
	diags := resp.State.Set(ctx, hostState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "import host state completed")
}
