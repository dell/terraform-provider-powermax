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
	"regexp"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/helper"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &StorageGroup{}
var _ resource.ResourceWithConfigure = &StorageGroup{}
var _ resource.ResourceWithImportState = &StorageGroup{}

// NewStorageGroup is a helper function to simplify the provider implementation.
func NewStorageGroup() resource.Resource {
	return &StorageGroup{}
}

// StorageGroup defines the resource implementation.
type StorageGroup struct {
	client *client.Client
}

// Metadata Resource metadata.
func (r *StorageGroup) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storagegroup"
}

// Schema Resource schema.
func (r *StorageGroup) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Resource for managing StorageGroups in PowerMax array. Updates are supported for the following parameters: `name`, `srp`, `enable_compression`, `service_level`, `host_io_limits`, `volume_ids`, `snapshot_policies`.",
		Description:         "Resource for managing StorageGroups in PowerMax array. Updates are supported for the following parameters: `name`, `srp`, `enable_compression`, `service_level`, `host_io_limits`, `volume_ids`, `snapshot_policies`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the storage group",
				MarkdownDescription: "The ID of the storage group",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the storage group. Only alphanumeric characters, underscores ( _ ), and hyphens (-) are allowed.",
				MarkdownDescription: "The name of the storage group. Only alphanumeric characters, underscores ( _ ), and hyphens (-) are allowed.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(64),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9_-]*$`),
						"must contain only alphanumeric characters and _-",
					),
				},
			},
			"slo": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The service level associated with the storage group",
				MarkdownDescription: "The service level associated with the storage group",
			},
			"srp_id": schema.StringAttribute{
				Required:            true,
				Description:         "The Srp to be associated with the Storage Group. An existing Srpor 'none' must be specified",
				MarkdownDescription: "The Srp to be associated with the Storage Group. An existing Srpor 'none' must be specified",
			},
			"service_level": schema.StringAttribute{
				Computed:            true,
				Description:         "The service level associated with the storage group",
				MarkdownDescription: "The service level associated with the storage group",
			},
			"workload": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The workload associated with the storage group",
				MarkdownDescription: "The workload associated with the storage group",
			},
			"slo_compliance": schema.StringAttribute{
				Computed:            true,
				Description:         "The service level compliance status of the storage group",
				MarkdownDescription: "The service level compliance status of the storage group",
			},
			"num_of_vols": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The number of volumes associated with the storage group",
				MarkdownDescription: "The number of volumes associated with the storage group",
			},
			"num_of_child_sgs": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of child storage groups associated with the storage group",
				MarkdownDescription: "The number of child storage groups associated with the storage group",
			},
			"num_of_parent_sgs": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of parent storage groups associated with the storage group",
				MarkdownDescription: "The number of parent storage groups associated with the storage group",
			},
			"num_of_masking_views": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of masking views associated with the storage group",
				MarkdownDescription: "The number of masking views associated with the storage group",
			},
			"num_of_snapshots": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of snapshots associated with the storage group",
				MarkdownDescription: "The number of snapshots associated with the storage group",
			},
			"num_of_snapshot_policies": schema.Int64Attribute{
				Computed:            true,
				Description:         "The number of snapshot policies associated with the storage group",
				MarkdownDescription: "The number of snapshot policies associated with the storage group",
			},
			"cap_gb": schema.NumberAttribute{
				Computed:            true,
				Description:         "The capacity of the storage group",
				MarkdownDescription: "The capacity of the storage group",
			},
			"device_emulation": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "The emulation of the volumes in the storage group",
				MarkdownDescription: "The emulation of the volumes in the storage group",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				Description:         "The storage group type",
				MarkdownDescription: "The storage group type",
			},
			"unprotected": schema.BoolAttribute{
				Computed:            true,
				Description:         "States whether the storage group is protected",
				MarkdownDescription: "States whether the storage group is protected",
			},
			"child_storage_group": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "The child storage group(s) associated with the storage group",
				MarkdownDescription: "The child storage group(s) associated with the storage group",
			},
			"parent_storage_group": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "The parent storage group(s) associated with the storage group",
				MarkdownDescription: "The parent storage group(s) associated with the storage group",
			},
			"maskingview": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "The masking views associated with the storage group",
				MarkdownDescription: "The masking views associated with the storage group",
			},
			"snapshot_policies": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "The snapshot policies associated with the storage group",
				MarkdownDescription: "The snapshot policies associated with the storage group",
			},
			"host_io_limit": schema.ObjectAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Host IO limit of the storage group",
				MarkdownDescription: "Host IO limit of the storage group",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				AttributeTypes: map[string]attr.Type{
					"host_io_limit_io_sec": types.StringType,
					"host_io_limit_mb_sec": types.StringType,
					"dynamic_distribution": types.StringType,
				},
			},
			"compression": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "States whether compression is enabled on storage group",
				MarkdownDescription: "States whether compression is enabled on storage group",
			},
			"compression_ratio": schema.StringAttribute{
				Computed:            true,
				Description:         "States whether compression is enabled on storage group",
				MarkdownDescription: "States whether compression is enabled on storage group",
			},
			"compression_ratio_to_one": schema.NumberAttribute{
				Computed:            true,
				Description:         "Compression ratio numeric value of the storage group",
				MarkdownDescription: "Compression ratio numeric value of the storage group",
			},
			"vp_saved_percent": schema.NumberAttribute{
				Computed:            true,
				Description:         "VP saved percentage figure",
				MarkdownDescription: "VP saved percentage figure",
			},
			"tags": schema.StringAttribute{
				Computed:            true,
				Description:         "The tags associated with the storage group",
				MarkdownDescription: "The tags associated with the storage group",
			},
			"uuid": schema.StringAttribute{
				Computed:            true,
				Description:         "Storage Group UUID",
				MarkdownDescription: "Storage Group UUID",
			},
			"unreducible_data_gb": schema.NumberAttribute{
				Computed:            true,
				Description:         "The amount of unreducible data in Gb.",
				MarkdownDescription: "The amount of unreducible data in Gb.",
			},
			"volume_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "The IDs of the volume associated with the storage group. Only pre-existing volumes are considered here.",
				MarkdownDescription: "The IDs of the volume associated with the storage group. Only pre-existing volumes are considered here.",
			},
		},
	}
}

// Configure the resource.
func (r *StorageGroup) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create a storage group.
func (r *StorageGroup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Storage Group...")
	var plan models.StorageGroupResourceModel
	var state models.StorageGroupResourceModel

	// Read Terraform plan into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	sgModel := r.client.PmaxOpenapiClient.SLOProvisioningApi.CreateStorageGroup(ctx, r.client.SymmetrixID)
	create := powermax.NewCreateStorageGroupParam(plan.StorageGroupID.ValueString())
	create.SetSrpId(plan.Srp.ValueString())
	create.SetSloBasedStorageGroupParam(helper.CreateSloParam(plan))
	sgModel = sgModel.CreateStorageGroupParam(*create)
	sg, _, err := sgModel.Execute()
	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create storage group, got error: %s", message))
		return
	}

	tflog.Debug(ctx, "created a resource", map[string]interface{}{
		"storage group": sg,
	})

	// Add or remove existing volumes to the storage group based on volume attributes
	err = helper.AddRemoveVolume(ctx, &plan, &state, r.client, plan.StorageGroupID.ValueString())
	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Failed to update volume", message)
		return
	}

	err = helper.UpdateSgState(ctx, r.client, plan.StorageGroupID.ValueString(), &state)
	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error updating state for storage group", message)
		return
	}

	// Save plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read a storage group.
func (r *StorageGroup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Storage Group...")
	var state models.StorageGroupResourceModel

	// Read Terraform prior state into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := helper.UpdateSgState(ctx, r.client, state.StorageGroupID.ValueString(), &state)
	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error updating state for storage group", message)
		return
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update a storage group.
func (r *StorageGroup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Storage group...")
	// Read Terraform plan into the model
	var plan models.StorageGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform state into the model
	var state models.StorageGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Read Storage Group ID from state in case of renaming
	stateID := state.StorageGroupID.ValueString()
	sgID := stateID
	payload := r.client.PmaxOpenapiClient.SLOProvisioningApi.ModifyStorageGroup(ctx, r.client.SymmetrixID, sgID)
	updateName := false
	// Storage Group update need to be done separately because only one payload is accepted by the REST API
	// Rename
	planID := plan.StorageGroupID.ValueString()
	if stateID != planID {
		payload = payload.EditStorageGroupParam(powermax.EditStorageGroupParam{
			EditStorageGroupActionParam: powermax.EditStorageGroupActionParam{
				RenameStorageGroupParam: &powermax.RenameStorageGroupParam{
					NewStorageGroupName: planID,
				},
			},
		})
		_, _, err := payload.Execute()
		if err != nil {
			errStr := ""
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Failed to update Storage Group ID(name)", message)
			tflog.Error(ctx, fmt.Sprintf("Failed to update Storage Group ID(name): %s", message))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Update Storage Group ID(name): %s", planID))
			sgID = planID
			state.StorageGroupID = types.StringValue(planID)
		}
		updateName = true
	}

	// Recreate the modify storage group param with the new name on a rename job
	if updateName {
		payload = r.client.PmaxOpenapiClient.SLOProvisioningApi.ModifyStorageGroup(ctx, r.client.SymmetrixID, planID)
	}

	// Edit Compression
	planCompression := plan.Compression.ValueBool()
	stateCompression := state.Compression.ValueBool()
	if planCompression != stateCompression {
		payload = payload.EditStorageGroupParam(powermax.EditStorageGroupParam{
			EditStorageGroupActionParam: powermax.EditStorageGroupActionParam{
				EditCompressionParam: &powermax.EditCompressionParam{
					Compression: &planCompression,
				},
			},
		})
		_, _, err := payload.Execute()
		if err != nil {
			errStr := ""
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Failed to update compression:", message)
			tflog.Error(ctx, fmt.Sprintf("Failed to update compression: %s", err.Error()))

		} else {
			tflog.Debug(ctx, fmt.Sprintf("Update compression: %t", planCompression))
			state.Compression = types.BoolValue(planCompression)
		}
	}

	// SetHostIOLimit
	if !plan.HostIOLimit.IsNull() && !plan.HostIOLimit.Equal(state.HostIOLimit) {
		hostIOLimit := helper.ConstructHostIOLimit(plan)
		payload = payload.EditStorageGroupParam(powermax.EditStorageGroupParam{
			EditStorageGroupActionParam: powermax.EditStorageGroupActionParam{
				SetHostIOLimitsParam: &powermax.SetHostIOLimitsParam{
					HostIoLimitMbSec:    &hostIOLimit.HostIOLimitMBSec,
					HostIoLimitIoSec:    &hostIOLimit.HostIOLimitIOSec,
					DynamicDistribution: &hostIOLimit.DynamicDistribution,
				},
			},
		})
		_, _, err := payload.Execute()
		if err != nil {
			errStr := ""
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Failed to update hostIOLimit:", message)
			tflog.Error(ctx, fmt.Sprintf("Failed to update hostIOLimit: %s", err.Error()))

		} else {
			tflog.Debug(ctx, fmt.Sprintf("Update hostIOLimit: %v", plan.HostIOLimit))
		}
	}

	// Edit Workload
	planWorkload := plan.Workload.ValueString()
	stateWorkload := state.Workload.ValueString()
	if planWorkload != stateWorkload {
		payload = payload.EditStorageGroupParam(powermax.EditStorageGroupParam{
			EditStorageGroupActionParam: powermax.EditStorageGroupActionParam{
				EditStorageGroupWorkloadParam: &powermax.EditStorageGroupWorkloadParam{
					WorkloadSelection: planWorkload,
				},
			},
		})
		_, _, err := payload.Execute()
		if err != nil {
			errStr := ""
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Failed to update workload:", message)
			tflog.Error(ctx, fmt.Sprintf("Failed to update workload: %s", err.Error()))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Update workload: %s", planWorkload))
			state.Workload = types.StringValue(planWorkload)
		}
	}

	// Edit Slo
	planSLO := plan.Slo.ValueString()
	stateSLO := state.Slo.ValueString()
	if planSLO != stateSLO {
		payload = payload.EditStorageGroupParam(powermax.EditStorageGroupParam{
			EditStorageGroupActionParam: powermax.EditStorageGroupActionParam{
				EditStorageGroupSLOParam: &powermax.EditStorageGroupSLOParam{
					SloId: planSLO,
				},
			},
		})
		_, _, err := payload.Execute()
		if err != nil {
			errStr := ""
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Failed to update Slo:", message)
			tflog.Error(ctx, fmt.Sprintf("Failed to update Slo: %s", err.Error()))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Update Slo: %s", planSLO))
			state.Slo = types.StringValue(planSLO)
		}
	}

	// Edit Srp
	planSRP := plan.Srp.ValueString()
	stateSRP := state.Srp.ValueString()
	if planSRP != stateSRP {
		payload = payload.EditStorageGroupParam(powermax.EditStorageGroupParam{
			EditStorageGroupActionParam: powermax.EditStorageGroupActionParam{
				EditStorageGroupSRPParam: &powermax.EditStorageGroupSRPParam{
					SrpId: planSRP,
				},
			},
		})
		_, _, err := payload.Execute()
		if err != nil {
			errStr := ""
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Failed to update Srp:", message)
			tflog.Error(ctx, fmt.Sprintf("Failed to update Srp: %s", message))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Update Srp: %s", planSLO))
			state.Srp = types.StringValue(planSRP)
		}
	}

	// Update Volume
	err := helper.AddRemoveVolume(ctx, &plan, &state, r.client, sgID)
	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Failed to update volume:", message)
		return
	}

	err = helper.UpdateSgState(ctx, r.client, sgID, &state)
	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error updating state for storage group:", message)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Applying this State!!! %v", state))
	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete deletes a Storage Group.
func (r *StorageGroup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Storage Group...")
	var data models.StorageGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	deletePayload := r.client.PmaxOpenapiClient.SLOProvisioningApi.DeleteStorageGroup(ctx, r.client.SymmetrixID, data.StorageGroupID.ValueString())
	_, err := deletePayload.Execute()
	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete storage group, got error: %s", message))
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports a Storage Group.
func (r *StorageGroup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
