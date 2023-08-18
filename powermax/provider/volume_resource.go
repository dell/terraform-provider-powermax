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
	"strings"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/helper"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type volumeResource struct {
	client *client.Client
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &volumeResource{}
	_ resource.ResourceWithConfigure   = &volumeResource{}
	_ resource.ResourceWithImportState = &volumeResource{}
)

// NewVolumeResource is a helper function to simplify the provider implementation.
func NewVolumeResource() resource.Resource {
	return &volumeResource{}
}

func (r volumeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

func (r volumeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Resource for managing Volumes in PowerMax array. Updates are supported for the following parameters: 'vol_name', 'mobility_id_enabled', 'size', 'cap_unit'",
		Description:         "Resource for managing Volumes in PowerMax array. Updates are supported for the following parameters: 'vol_name', 'mobility_id_enabled', 'size', 'cap_unit'",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The ID of the volume.",
				MarkdownDescription: "The ID of the volume.",
				Computed:            true,
			},
			"vol_name": schema.StringAttribute{
				Description:         "The name of the volume.",
				MarkdownDescription: "The name of the volume.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"sg_name": schema.StringAttribute{
				Description:         "The name of the storage group. sg_name is required while creating the volume.",
				MarkdownDescription: "The name of the storage group. sg_name is required while creating the volume.",
				Computed:            true,
				Optional:            true,
			},
			"size": schema.NumberAttribute{
				Description:         "The size of the volume.",
				MarkdownDescription: "The size of the volume.",
				Required:            true,
			},
			"cap_unit": schema.StringAttribute{
				Description:         "The Capacity Unit corresponding to the size.",
				MarkdownDescription: "The Capacity Unit corresponding to the size.",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(helper.CapacityUnitGb),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						helper.CapacityUnitMb,
						helper.CapacityUnitGb,
						helper.CapacityUnitTb,
						helper.CapacityUnitCyl,
					}...),
				},
			},
			"type": schema.StringAttribute{
				Computed:            true,
				Description:         "The type of the volume.",
				MarkdownDescription: "The type of the volume.",
			},
			"emulation": schema.StringAttribute{
				Computed:            true,
				Description:         "The emulation of the volume Enumeration values.",
				MarkdownDescription: "The emulation of the volume Enumeration values.",
			},
			"ssid": schema.StringAttribute{
				Computed:            true,
				Description:         "The ssid of the volume.",
				MarkdownDescription: "The ssid of the volume.",
			},
			"allocated_percent": schema.Int64Attribute{
				Computed:            true,
				Description:         "The allocated percentage of the volume.",
				MarkdownDescription: "The allocated percentage of the volume.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				Description:         "The status of the volume.",
				MarkdownDescription: "The status of the volume.",
			},
			"reserved": schema.BoolAttribute{
				Computed:            true,
				Description:         "States whether the volume is reserved.",
				MarkdownDescription: "States whether the volume is reserved.",
			},
			"pinned": schema.BoolAttribute{
				Computed:            true,
				Description:         "States whether the volume is pinned.",
				MarkdownDescription: "States whether the volume is pinned.",
			},
			"wwn": schema.StringAttribute{
				Computed:            true,
				Description:         "The WWN of the volume.",
				MarkdownDescription: "The WWN of the volume.",
			},
			"encapsulated": schema.BoolAttribute{
				Computed:            true,
				Description:         "States whether the volume is encapsulated.",
				MarkdownDescription: "States whether the volume is encapsulated.",
			},
			"num_of_storage_groups": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of storage groups associated with the volume.",
				MarkdownDescription: "The number of storage groups associated with the volume.",
			},
			"num_of_front_end_paths": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of front end paths of the volume.",
				MarkdownDescription: "The number of front end paths of the volume.",
			},
			"rdf_group_ids": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "The RDF groups associated with the volume.",
				MarkdownDescription: "The RDF groups associated with the volume.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rdf_group_number": schema.Int64Attribute{
							Description:         "The number of rdf group.",
							MarkdownDescription: "The number of rdf group.",
							Computed:            true,
						},
						"label": schema.StringAttribute{
							Description:         "The label of the rdf group.",
							MarkdownDescription: "The label of the rdf group.",
							Computed:            true,
						},
					},
				},
			},
			"snapvx_source": schema.BoolAttribute{
				Computed:            true,
				Description:         "States whether the volume is a snapvx source.",
				MarkdownDescription: "States whether the volume is a snapvx source.",
			},
			"snapvx_target": schema.BoolAttribute{
				Computed:            true,
				Description:         "States whether the volume is a snapvx target.",
				MarkdownDescription: "States whether the volume is a snapvx target.",
			},
			"has_effective_wwn": schema.BoolAttribute{
				Computed:            true,
				Description:         "States whether volume has effective WWN.",
				MarkdownDescription: "States whether volume has effective WWN.",
			},
			"effective_wwn": schema.StringAttribute{
				Computed:            true,
				Description:         "Effective WWN of the volume.",
				MarkdownDescription: "Effective WWN of the volume.",
			},
			"encapsulated_wwn": schema.StringAttribute{
				Computed:            true,
				Description:         "Encapsulated  WWN of the volume.",
				MarkdownDescription: "Encapsulated  WWN of the volume.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mobility_id_enabled": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "States whether mobility ID is enabled on the volume.",
				MarkdownDescription: "States whether mobility ID is enabled on the volume.",
			},
			"unreducible_data_gb": schema.NumberAttribute{
				Computed:            true,
				Description:         "The amount of unreducible data in Gb.",
				MarkdownDescription: "The amount of unreducible data in Gb.",
			},
			"nguid": schema.StringAttribute{
				Computed:            true,
				Description:         "The NGUID of the volume.",
				MarkdownDescription: "The NGUID of the volume.",
			},
			"oracle_instance_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Oracle instance name associated with the volume.",
				MarkdownDescription: "Oracle instance name associated with the volume.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"storage_groups": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of storage groups which are associated with the volume.",
				MarkdownDescription: "List of storage groups which are associated with the volume.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"storage_group_name": schema.StringAttribute{
							Description:         "The ID of the storage group.",
							MarkdownDescription: "The ID of the storage group.",
							Computed:            true,
						},
						"parent_storage_group_name": schema.StringAttribute{
							Description:         "The ID of the storage group parents.",
							MarkdownDescription: "The ID of the storage group parents.",
							Computed:            true,
						},
					},
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"symmetrix_port_key": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "The symmetrix ports associated with the volume.",
				MarkdownDescription: "The symmetrix ports associated with the volume.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"director_id": schema.StringAttribute{
							Description:         "The ID of the director.",
							MarkdownDescription: "The ID of the director.",
							Computed:            true,
						},
						"port_id": schema.StringAttribute{
							Description:         "The ID of the symmetrix port.",
							MarkdownDescription: "The ID of the symmetrix port.",
							Computed:            true,
						},
					},
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure - defines configuration for volume resource.
func (r *volumeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *c.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

// Create - method to create volume resource.
func (r volumeResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Info(ctx, "creating volume")

	var plan models.VolumeResource
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	if !plan.Size.IsNull() {
		size, _ := plan.Size.ValueBigFloat().Float64()
		if plan.CapUnit.ValueString() == "CYL" && size != float64(int(size)) {
			response.Diagnostics.AddError(
				"Error creating volume",
				fmt.Sprintf("Could not create volume %s with error: %s", plan.VolumeIdentifier.ValueString(), "Invalid Config, size type 'CYL' must be integer"),
			)
			return
		}
	}

	if plan.StorageGroupName.ValueString() == "" {
		response.Diagnostics.AddError(
			"Error creating volume",
			fmt.Sprintf("Could not create volume %s with error: %s", plan.VolumeIdentifier.ValueString(), "storage group name cannot be empty"),
		)
		return
	}

	volResponse, _, err := helper.CreateVolume(ctx, *r.client, plan)
	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error creating volume",
			fmt.Sprintf("Could not create volume %s with error: %s", plan.VolumeIdentifier.ValueString(), message))
		return
	}
	tflog.Debug(ctx, "create volume in storage groups response", map[string]interface{}{
		"volResponse": volResponse,
	})
	// Extrct the new volume ID from the storage group
	volState := models.VolumeResource{}
	volumeIDListInStorageGroup, _, err := helper.ListVolumes(ctx, *r.client, plan)
	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error creating volume",
			fmt.Sprintf("Could not find volume %s after creating with error: %s", plan.VolumeIdentifier.ValueString(), message))
		return
	}
	volID := ""
	for _, v := range volumeIDListInStorageGroup.ResultList.Result {
		for _, v2 := range v {
			volID = fmt.Sprint(v2)
		}
	}
	if volID == "" {
		response.Diagnostics.AddError("Error creating volume",
			fmt.Sprintf("Could not find find volume id for %s after creating", plan.VolumeIdentifier.ValueString()),
		)
		return
	}
	// Now that we have the ID get the specific volume info
	vol, _, err := helper.GetVolume(ctx, *r.client, volID)

	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error creating volume",
			fmt.Sprintf("Could not find volume %s after creating with error: %s", plan.VolumeIdentifier.ValueString(), message))
		return
	}
	tflog.Debug(ctx, "updating create volume state", map[string]interface{}{
		"volResponse": volResponse,
		"vol":         vol,
		"plan":        plan,
		"volState":    volState,
		"volId":       volID,
	})
	err = helper.UpdateVolResourceState(ctx, &volState, vol, &plan)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating volume",
			fmt.Sprintf("Could not upda volume state %s with error: %s", plan.VolumeIdentifier.ValueString(), err.Error()),
		)
		return
	}

	diags = response.State.Set(ctx, volState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create volume completed")
}

func (r volumeResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Info(ctx, "reading volume")
	var volState models.VolumeResource
	diags := request.State.Get(ctx, &volState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	volID := volState.ID.ValueString()
	tflog.Debug(ctx, "calling get volume by ID", map[string]interface{}{
		"symmetrixID": r.client.SymmetrixID,
		"volumeID":    volID,
	})
	volResponse, _, err := helper.GetVolume(ctx, *r.client, volID)
	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError(
			"Error reading volume",
			fmt.Sprintf("Could not read volume %s with error: %s", volID, message),
		)

		return
	}
	tflog.Debug(ctx, "get volume by ID response", map[string]interface{}{
		"volResponse": volResponse,
	})

	tflog.Debug(ctx, "updating read volume state", map[string]interface{}{
		"volResponse": volResponse,
		"volState":    volState,
	})
	err = helper.UpdateVolResourceState(ctx, &volState, volResponse, nil)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating volume",
			fmt.Sprintf("Could not update volume %s with error: %s", volID, err.Error()),
		)
		return
	}
	diags = response.State.Set(ctx, volState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read volume completed")
}

// Update VolumeResource
// Supported updates: vol_name, mobility_id_enabled, size, cap_unit.
func (r volumeResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	tflog.Info(ctx, "updating volume")
	var planVol models.VolumeResource
	diags := request.Plan.Get(ctx, &planVol)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Fetched vol from plan")
	var stateVol models.VolumeResource
	diags = response.State.Get(ctx, &stateVol)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if !planVol.Size.IsNull() {
		size, _ := planVol.Size.ValueBigFloat().Float64()
		if planVol.CapUnit.ValueString() == "CYL" && size != float64(int(size)) {
			response.Diagnostics.AddError(
				"Error updating volume",
				fmt.Sprintf("Could not create volume %s with error: %s", planVol.VolumeIdentifier.ValueString(), "Invalid Config, size type 'CYL' must be integer"),
			)
			return
		}
	}

	tflog.Debug(ctx, "calling update volume on pmax client", map[string]interface{}{
		"planVol":  planVol,
		"stateVol": stateVol,
	})
	updatedParams, updateFailedParameters, errMessages := helper.UpdateVol(ctx, r.client, planVol, stateVol)
	if len(errMessages) > 0 || len(updateFailedParameters) > 0 {
		errMessage := strings.Join(errMessages, ",\n")
		response.Diagnostics.AddError(
			fmt.Sprintf("Failed to update all parameters of Volume, updated parameters are %v and parameters failed to update are %v", updatedParams, updateFailedParameters),
			errMessage)
	}

	volID := stateVol.ID.ValueString()
	tflog.Debug(ctx, "calling get volume by ID on pmax client", map[string]interface{}{
		"symmetrixID": r.client.SymmetrixID,
		"volumeID":    volID,
	})
	volResponse, _, err := helper.GetVolume(ctx, *r.client, volID)
	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError(
			"Error reading volume",
			fmt.Sprintf("Could not read volume %s with error: %s", volID, message),
		)
		return
	}
	tflog.Debug(ctx, "get volume by ID response", map[string]interface{}{
		"volResponse": volResponse,
	})

	tflog.Debug(ctx, "updating volume state", map[string]interface{}{
		"volResponse": volResponse,
		"planVol":     planVol,
	})
	err = helper.UpdateVolResourceState(ctx, &stateVol, volResponse, &planVol)
	if err != nil {
		response.Diagnostics.AddError(
			"Error reading volume",
			fmt.Sprintf("Could not read volume %s with error: %s", volID, err.Error()),
		)
		return
	}
	diags = response.State.Set(ctx, stateVol)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update volume completed")
}

func (r volumeResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting volume")
	var volumeState models.VolumeResource
	diags := request.State.Get(ctx, &volumeState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	volumeID := volumeState.ID.ValueString()
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
	}
	removeVol := make([]string, 0)
	removeVol = append(removeVol, volumeID)
	// Unbind volume
	var sgAssociatedWithVolume []models.StorageGroupName
	diags = volumeState.StorageGroups.ElementsAs(ctx, &sgAssociatedWithVolume, true)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
	}

	for _, associatedSG := range sgAssociatedWithVolume {
		tflog.Debug(ctx, "calling get storage group on pmax client", map[string]interface{}{
			"symmetrixID":    r.client.SymmetrixID,
			"storageGroupID": associatedSG.StorageGroupName.ValueString(),
		})
		sgModel := r.client.PmaxOpenapiClient.SLOProvisioningApi.GetStorageGroup2(ctx, r.client.SymmetrixID, associatedSG.StorageGroupName.ValueString())
		sg, _, _ := sgModel.Execute()
		tflog.Debug(ctx, "get storage group response", map[string]interface{}{
			"associatedSG": sg,
		})
		if sg != nil {
			tflog.Debug(ctx, "calling remove volumes from storage group on pmax client", map[string]interface{}{
				"symmetrixID":    r.client.SymmetrixID,
				"storageGroupID": sg,
				"volumeID":       volumeID,
			})
			deleteParam := r.client.PmaxOpenapiClient.SLOProvisioningApi.ModifyStorageGroup(ctx, r.client.SymmetrixID, associatedSG.StorageGroupName.ValueString())
			deleteParam = deleteParam.EditStorageGroupParam(
				powermax.EditStorageGroupParam{
					EditStorageGroupActionParam: powermax.EditStorageGroupActionParam{
						RemoveVolumeParam: &powermax.RemoveVolumeParam{
							VolumeId: removeVol,
						},
					},
				},
			)
			_, _, err := deleteParam.Execute()
			if err != nil {
				errStr := ""
				message := helper.GetErrorString(err, errStr)
				response.Diagnostics.AddError(
					"Error removing volume from storage group",
					fmt.Sprintf("Could not remove  Volume ID: %s from storage group: %s with error: %s",
						volumeID, volumeState.StorageGroupName.ValueString(), message),
				)

				return
			}
		}
	}
	tflog.Debug(ctx, "calling delete volume on pmax client", map[string]interface{}{
		"symmetrixID": r.client.SymmetrixID,
		"volumeID":    volumeID,
	})
	delParam := r.client.PmaxOpenapiClient.SLOProvisioningApi.DeleteVolume(ctx, r.client.SymmetrixID, volumeID)
	_, err := delParam.Execute()
	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError(
			"Error deleting volume",
			fmt.Sprintf("Could not remove Volume ID: %s with error: %s ",
				volumeID, message),
		)

	}
	response.State.RemoveResource(ctx)
	tflog.Info(ctx, "delete volume completed")
}

func (r volumeResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
	var stateVol models.VolumeResource
	response.State.Get(ctx, &stateVol)
	// For importing volume, storage group for creating should leave as empty
	stateVol.StorageGroupName = types.StringValue("")
	// Default cap unit
	stateVol.CapUnit = types.StringValue(helper.CapacityUnitGb)
	response.State.Set(ctx, stateVol)
}
