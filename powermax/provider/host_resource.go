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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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

// Metadata returns the metadata for the resource
func (r *Host) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

// Schema returns the schema for the resource
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
	resp.Schema = schema.Schema{

		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Resource for managing Host in PowerMax array.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the host.",
				MarkdownDescription: "The ID of the host.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the host.",
				MarkdownDescription: "The name of the host.",
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
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
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
				Required: true,
				Attributes: map[string]schema.Attribute{
					"volume_set_addressing": schema.SingleNestedAttribute{
						Optional:            true,
						Computed:            true,
						Attributes:          hostFlagNestedAttr,
						Description:         "It enables the volume set addressing mode.",
						MarkdownDescription: "It enables the volume set addressing mode.",
						PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
						Default:             objectdefault.StaticValue(objd),
					},
					"disable_q_reset_on_ua": schema.SingleNestedAttribute{
						Optional:            true,
						Computed:            true,
						Attributes:          hostFlagNestedAttr,
						Description:         "It is used for hosts that do not expect the queue to be flushed on a 0629 sense.",
						MarkdownDescription: "It is used for hosts that do not expect the queue to be flushed on a 0629 sense.",
						PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
						Default:             objectdefault.StaticValue(objd),
					},
					"environ_set": schema.SingleNestedAttribute{
						Optional:            true,
						Computed:            true,
						Attributes:          hostFlagNestedAttr,
						Description:         "It enables the environmental error reporting by the storage system to the host on the specific port.",
						MarkdownDescription: "It enables the environmental error reporting by the storage system to the host on the specific port.",
						PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
						Default:             objectdefault.StaticValue(objd),
					},
					"openvms": schema.SingleNestedAttribute{
						Optional:            true,
						Computed:            true,
						Attributes:          hostFlagNestedAttr,
						Description:         "This attribute enables an Open VMS fibre connection.",
						MarkdownDescription: "This attribute enables an Open VMS fibre connection.",
						PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
						Default:             objectdefault.StaticValue(objd),
					},
					"avoid_reset_broadcast": schema.SingleNestedAttribute{
						Optional:            true,
						Computed:            true,
						Attributes:          hostFlagNestedAttr,
						Description:         "It enables a SCSI bus reset to only occur to the port that received the reset.",
						MarkdownDescription: "It enables a SCSI bus reset to only occur to the port that received the reset.",
						PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
						Default:             objectdefault.StaticValue(objd),
					},
					"scsi_3": schema.SingleNestedAttribute{
						Optional:            true,
						Computed:            true,
						Attributes:          hostFlagNestedAttr,
						Description:         "Alters the inquiry data to report that the storage system supports the SCSI-3 protocol.",
						MarkdownDescription: "Alters the inquiry data to report that the storage system supports the SCSI-3 protocol.",
						PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
						Default:             objectdefault.StaticValue(objd),
					},
					"spc2_protocol_version": schema.SingleNestedAttribute{
						Optional:            true,
						Computed:            true,
						Attributes:          hostFlagNestedAttr,
						Description:         "When setting this flag, the port must be offline.",
						MarkdownDescription: "When setting this flag, the port must be offline.",
						PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
						Default:             objectdefault.StaticValue(objd),
					},
					"scsi_support1": schema.SingleNestedAttribute{
						Optional:            true,
						Computed:            true,
						Attributes:          hostFlagNestedAttr,
						Description:         "This attribute provides a stricter compliance with SCSI standards.",
						MarkdownDescription: "This attribute provides a stricter compliance with SCSI standards.",
						PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
						Default:             objectdefault.StaticValue(objd),
					},
				},
				Description:         "Flags set for the host. When host_flags = {} then default flags will be considered.",
				MarkdownDescription: "Flags set for the host. When host_flags = {} then default flags will be considered.",
			},
		},
	}
}

// Configure configure client for host resource
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

// Create creates a host and refresh state
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
	hostFlags := pmaxTypes.HostFlags{
		VolumeSetAddressing: &pmaxTypes.HostFlag{
			Enabled:  planHost.HostFlags.VolumeSetAddressing.Enabled.ValueBool(),
			Override: planHost.HostFlags.VolumeSetAddressing.Override.ValueBool(),
		},
		DisableQResetOnUA: &pmaxTypes.HostFlag{
			Enabled:  planHost.HostFlags.DisableQResetOnUA.Enabled.ValueBool(),
			Override: planHost.HostFlags.DisableQResetOnUA.Override.ValueBool(),
		},
		EnvironSet: &pmaxTypes.HostFlag{
			Enabled:  planHost.HostFlags.EnvironSet.Enabled.ValueBool(),
			Override: planHost.HostFlags.EnvironSet.Override.ValueBool(),
		},
		AvoidResetBroadcast: &pmaxTypes.HostFlag{
			Enabled:  planHost.HostFlags.AvoidResetBroadcast.Enabled.ValueBool(),
			Override: planHost.HostFlags.AvoidResetBroadcast.Override.ValueBool(),
		},
		OpenVMS: &pmaxTypes.HostFlag{
			Enabled:  planHost.HostFlags.OpenVMS.Enabled.ValueBool(),
			Override: planHost.HostFlags.OpenVMS.Override.ValueBool(),
		},
		SCSI3: &pmaxTypes.HostFlag{
			Enabled:  planHost.HostFlags.SCSI3.Enabled.ValueBool(),
			Override: planHost.HostFlags.SCSI3.Override.ValueBool(),
		},
		Spc2ProtocolVersion: &pmaxTypes.HostFlag{
			Enabled:  planHost.HostFlags.Spc2ProtocolVersion.Enabled.ValueBool(),
			Override: planHost.HostFlags.Spc2ProtocolVersion.Override.ValueBool(),
		},
		SCSISupport1: &pmaxTypes.HostFlag{
			Enabled:  planHost.HostFlags.SCSISupport1.Enabled.ValueBool(),
			Override: planHost.HostFlags.SCSISupport1.Override.ValueBool(),
		},
	}

	hostCreateResp, err := r.client.PmaxClient.CreateHost(ctx, r.client.SymmetrixID, planHost.Name.ValueString(), initiators, &hostFlags)
	if err != nil {
		hostID := planHost.Name.ValueString()
		resp.Diagnostics.AddError(
			"Error creating host",
			constants.CreateHostDetailErrorMsg+hostID+"with error: "+err.Error(),
		)
		hostCreateResp, getHostErr := r.client.PmaxClient.GetHostByID(ctx, r.client.SymmetrixID, hostID)
		if hostCreateResp != nil || getHostErr == nil {
			err := r.client.PmaxClient.DeleteHost(ctx, r.client.SymmetrixID, hostID)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error deleting the invalid host, This may be a dangling resource and needs to be deleted manually",
					constants.CreateHostDetailErrorMsg+hostID+"with error: "+err.Error(),
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
	err := r.client.PmaxClient.DeleteHost(ctx, r.client.SymmetrixID, hostID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting host",
			constants.DeleteHostDetailsErrorMsg+hostID+" with error: "+err.Error(),
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
	hostResponse, err := r.client.PmaxClient.GetHostByID(ctx, r.client.SymmetrixID, hostID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading host",
			constants.ReadHostDetailsErrorMsg+hostID+" with error: "+err.Error(),
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
	host, err := r.client.PmaxClient.GetHostByID(ctx, r.client.SymmetrixID, hostID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading host",
			constants.ReadHostDetailsErrorMsg+hostID+" with error: "+err.Error(),
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
	hostResponse, err := r.client.PmaxClient.GetHostByID(ctx, r.client.SymmetrixID, hostID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading host",
			constants.ImportHostDetailsErrorMsg+hostID+" with error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Get Host By ID response", map[string]interface{}{
		"Host Response": hostResponse,
	})

	tflog.Debug(ctx, "updating host state after import")
	helper.UpdateHostState(&hostState, hostResponse.Initiators, hostResponse)
	diags := resp.State.Set(ctx, hostState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "import host state completed")
}
