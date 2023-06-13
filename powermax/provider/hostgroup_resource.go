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
	"dell/powermax-go-client"
	"fmt"
	"net/http"
	"strings"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/constants"
	"terraform-provider-powermax/powermax/helper"
	"terraform-provider-powermax/powermax/models"

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

// Ensure implementation.
var (
	_ resource.Resource                = &HostGroup{}
	_ resource.ResourceWithConfigure   = &HostGroup{}
	_ resource.ResourceWithImportState = &HostGroup{}
)

// NewHostGroup is a helper function to simplify the provider implementation.
func NewHostGroup() resource.Resource {
	return &HostGroup{}
}

// HostGroup is the resource implementation.
type HostGroup struct {
	client *client.Client
}

// Metadata returns the metadata for the resource.
func (r *HostGroup) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hostgroup"
}

// Schema returns the schema for the resource.
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
	objd, _ := basetypes.NewObjectValue(
		map[string]attr.Type{
			"override": types.BoolType,
			"enabled":  types.BoolType,
		},
		map[string]attr.Value{
			"override": types.BoolValue(false),
			"enabled":  types.BoolValue(false),
		},
	)

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
		// Description for Docs
		MarkdownDescription: "Resource for managing HostGroups for a PowerMax Array",
		Description:         "Resource for managing HostGroups for a PowerMax Array",
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
			"consistent_lun": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				Description:         "It enables the rejection of any masking operation involving this hostgroup that would result in inconsistent LUN values.",
				MarkdownDescription: "It enables the rejection of any masking operation involving this hostgroup that would result in inconsistent LUN values.",
			},
			"host_flags": schema.SingleNestedAttribute{
				Description:         "Host Flags set for the hostgroup. When host_flags = {} or not set then default flags will be considered.",
				MarkdownDescription: "Host Flags set for the hostgroup. When host_flags = {} or not set then default flags will be considered.",
				Optional:            true,
				Computed:            true,
				Default:             objectdefault.StaticValue(hostDefaultObj),
				Attributes:          hostFlagAttr,
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

// Configure the HostGroup resource.
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

// Create a HostGroup resource.
func (r *HostGroup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Create Host Group")
	var plan models.HostGroupModel
	var state models.HostGroupModel
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

	if len(hostIds) > 0 {
		for _, id := range hostIds {
			if id == "" {
				resp.Diagnostics.AddError("Invaid Configuration: ", "host_ids can not have an empty \"\" value")
				return
			}
		}
	}

	hostFlags := helper.HandleHostFlag(plan)

	tflog.Info(ctx, "calling create hostgroup with client", map[string]interface{}{
		"symmetrixID": r.client.SymmetrixID,
		"host":        plan.Name.ValueString(),
		"hostIds":     hostIds,
		"hostFlags":   hostFlags,
	})

	newHgModel := r.client.PmaxOpenapiClient.SLOProvisioningApi.CreateHostGroup(ctx, r.client.SymmetrixID)
	create := powermax.NewCreateHostGroupParam(plan.Name.ValueString())
	create.SetHostFlags(hostFlags)
	create.SetHostId(hostIds)
	newHgModel = newHgModel.CreateHostGroupParam(*create)
	newHostGroup, _, err := newHgModel.Execute()

	if err != nil {
		hostgroupID := plan.Name.ValueString()
		resp.Diagnostics.AddError("Client Error", "Unable to create host group, please make sure only existing host(s) are set in the host_id flag")
		if err != nil {
			tflog.Debug(ctx, err.Error())
		}
		//Attempt to remove any partially created obejcts if there are any
		hgModel := r.client.PmaxOpenapiClient.SLOProvisioningApi.GetHostGroup(ctx, r.client.SymmetrixID, hostgroupID)
		hostGroupResponse, _, getHostGroupErr := hgModel.Execute()
		if hostGroupResponse != nil || getHostGroupErr == nil {
			deleteModel := r.client.PmaxOpenapiClient.SLOProvisioningApi.DeleteHostGroup(ctx, r.client.SymmetrixID, hostgroupID)
			_, err := deleteModel.Execute()
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

// Read a HostGroup resource.
func (r *HostGroup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Host Group")
	var state models.HostGroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hostGroupID := state.ID.ValueString()
	tflog.Debug(ctx, "fetching hostgroup by ID", map[string]interface{}{
		"symmetricxId": r.client.SymmetrixID,
		"hostGroupID":  hostGroupID,
	})
	hgModel := r.client.PmaxOpenapiClient.SLOProvisioningApi.GetHostGroup(ctx, r.client.SymmetrixID, hostGroupID)
	hgResponse, resp1, err := hgModel.Execute()
	tflog.Debug(ctx, "Get HostGroup By ID response", map[string]interface{}{
		"HostGroup Response": hgResponse,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading hostGroup",
			constants.ReadHostGroupDetailsErrorMsg+hostGroupID+" with error: "+err.Error(),
		)
		return
	}
	if resp1.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to Read PowerMax Host Groups. Got http error:",
			resp1.Status,
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
// Supported updates: name, host_ids, host_flags.
func (r *HostGroup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating HostGroup")
	var planHostGroup models.HostGroupModel
	diags := req.Plan.Get(ctx, &planHostGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "fetched hostgroup details from plan")

	var stateHostGroup models.HostGroupModel
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
	hgModel := r.client.PmaxOpenapiClient.SLOProvisioningApi.GetHostGroup(ctx, r.client.SymmetrixID, hostGroupID)
	hostGroupResponse, resp1, err := hgModel.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading hostgroup",
			constants.ReadHostGroupDetailsErrorMsg+hostGroupID+" with error: "+err.Error(),
		)
		return
	}
	if resp1.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to Read PowerMax Host Groups. Got http error:",
			resp1.Status,
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

// Delete a HostGroup resource.
func (r *HostGroup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting hostgroup")
	var hostGroupState models.HostGroupModel
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
	deleteModel := r.client.PmaxOpenapiClient.SLOProvisioningApi.DeleteHostGroup(ctx, r.client.SymmetrixID, hostGroupID)
	_, err := deleteModel.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting hostGroup",
			constants.DeleteHostGroupDetailsErrorMsg+hostGroupID+" with error: "+err.Error(),
		)
	}

	tflog.Info(ctx, "delete hostgroup complete")
}

// ImportState method used to import hostgroup state.
func (r *HostGroup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing Hostgroup State")
	var hostGroupState models.HostGroupModel
	hostGroupID := req.ID
	tflog.Debug(ctx, "fetching Hostgroup by ID", map[string]interface{}{
		"symmetrixID": r.client.SymmetrixID,
		"hostID":      hostGroupID,
	})
	hgModel := r.client.PmaxOpenapiClient.SLOProvisioningApi.GetHostGroup(ctx, r.client.SymmetrixID, hostGroupID)
	hostGroupResponse, resp1, err := hgModel.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading hostgroup",
			constants.ImportHostGroupDetailsErrorMsg+hostGroupID+" with error: "+err.Error(),
		)
		return
	}
	if resp1.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to Read PowerMax Host Groups. Got http error:",
			resp1.Status,
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
