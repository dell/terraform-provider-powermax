// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package provider

import (
	"context"
	"fmt"
	"reflect"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/helper"
	"terraform-provider-powermax/powermax/models"

	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &StorageGroup{}
var _ resource.ResourceWithConfigure = &StorageGroup{}
var _ resource.ResourceWithImportState = &StorageGroup{}

func NewStorageGroup() resource.Resource {
	return &StorageGroup{}
}

// StorageGroup defines the resource implementation.
type StorageGroup struct {
	client *client.Client
}

func (r *StorageGroup) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storagegroup"
}

func (r *StorageGroup) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Resource for managing StorageGroups in PowerMax array. Updates are supported for the following parameters: `name`, `srp`, `enable_compression`, `service_level`, `host_io_limits`, `volume_ids`, `snapshot_policies`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the storage group",
				MarkdownDescription: "The ID of the storage group",
			},
			"storage_group_id": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the storage group",
				MarkdownDescription: "The name of the storage group",
			},
			"slo": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The service level associated with the storage group",
				MarkdownDescription: "The service level associated with the storage group",
			},
			"srp_id": schema.StringAttribute{
				Required:            true,
				Description:         "The SRP to be associated with the Storage Group. An existing SRP or 'none' must be specified",
				MarkdownDescription: "The SRP to be associated with the Storage Group. An existing SRP or 'none' must be specified",
			},
			"service_level": schema.StringAttribute{
				Computed:            true,
				Description:         "The service level associated with the storage group",
				MarkdownDescription: "The service level associated with the storage group",
			},
			"workload": schema.StringAttribute{
				Computed:            true,
				Description:         "The workload associated with the storage group",
				MarkdownDescription: "The workload associated with the storage group",
			},
			"slo_compliance": schema.StringAttribute{
				Computed:            true,
				Description:         "The service level compliance status of the storage group",
				MarkdownDescription: "The service level compliance status of the storage group",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
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
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
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
			"host_io_limit": schema.MapAttribute{
				Computed:            true,
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Host IO limit of the storage group",
				MarkdownDescription: "Host IO limit of the storage group",
				PlanModifiers:       []planmodifier.Map{mapplanmodifier.UseStateForUnknown()},
			},
			"compression": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "States whether compression is enabled on storage group",
				MarkdownDescription: "States whether compression is enabled on storage group",
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"compression_ratio": schema.StringAttribute{
				Computed:            true,
				Description:         "States whether compression is enabled on storage group",
				MarkdownDescription: "States whether compression is enabled on storage group",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"compression_ratio_to_one": schema.NumberAttribute{
				Computed:            true,
				Description:         "Compression ratio numeric value of the storage group",
				MarkdownDescription: "Compression ratio numeric value of the storage group",
				PlanModifiers:       []planmodifier.Number{numberplanmodifier.UseStateForUnknown()},
			},
			"vp_saved_percent": schema.NumberAttribute{
				Computed:            true,
				Description:         "VP saved percentage figure",
				MarkdownDescription: "VP saved percentage figure",
				PlanModifiers:       []planmodifier.Number{numberplanmodifier.UseStateForUnknown()},
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
				MarkdownDescription: "SThe amount of unreducible data in Gb.",
				PlanModifiers:       []planmodifier.Number{numberplanmodifier.UseStateForUnknown()},
			},
			//"volume_ids": schema.SetAttribute{
			//	ElementType:         types.StringType,
			//	Optional:            true,
			//	Computed:            true,
			//	Description:         "The IDs of the volume associated with the storage group. Only pre-existing volumes are considered here.",
			//	MarkdownDescription: "The IDs of the volume associated with the storage group. Only pre-existing volumes are considered here.",
			//	PlanModifiers:       []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
			//},
			//"volume_size": schema.StringAttribute{
			//	Optional:            true,
			//	Computed:            true,
			//	Description:         "The desired volume size of storage group",
			//	MarkdownDescription: "The desired volume size of storage group",
			//	PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			//},
			//"capacity_unit": schema.StringAttribute{
			//	Optional:            true,
			//	Computed:            true,
			//	Description:         "The unit of capacity enumeration values: CYL,MB,GB,TB",
			//	MarkdownDescription: "The unit of capacity enumeration values: CYL,MB,GB,TB",
			//	PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			//},
			//"volume_identifier_name": schema.StringAttribute{
			//	Optional:            true,
			//	Computed:            true,
			//	Description:         "Volume identifier name",
			//	MarkdownDescription: "Volume identifier name",
			//	PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			//},
		},
	}
}

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

func (r *StorageGroup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Storage Group...")
	var plan models.StorageGroupResourceModel
	var state models.StorageGroupResourceModel

	// Read Terraform plan into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	optionalPayload := make(map[string]interface{})

	// HostIOLimit
	planHostIOLimit := make(map[string]string)
	diag := plan.HostIOLimit.ElementsAs(ctx, &planHostIOLimit, true)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(planHostIOLimit) > 0 {
		hostIOLimitParam := &pmaxTypes.SetHostIOLimitsParam{
			HostIOLimitIOSec:    planHostIOLimit["host_io_limit_io_sec"],
			HostIOLimitMBSec:    planHostIOLimit["host_io_limit_mb_sec"],
			DynamicDistribution: planHostIOLimit["dynamicDistribution"],
		}

		optionalPayload["hostLimits"] = hostIOLimitParam
	}

	sg, err := r.client.PmaxClient.CreateStorageGroup(ctx, r.client.SymmetrixID, plan.StorageGroupID.ValueString(), plan.SRP.ValueString(), plan.SLO.ValueString(), false, optionalPayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create storage group, got error: %s", err.Error()))
		return
	}

	tflog.Debug(ctx, "created a resource", map[string]interface{}{
		"storage group": sg,
	})

	// Add new or existing volumes to the storage group based on volume attributes
	//err = UpdateVolume(ctx, &plan, &state, r)
	//if err != nil {
	//	resp.Diagnostics.AddError("Failed to update volume", err.Error())
	//	return
	//}

	// Update all fields of state
	storageGroup, err := r.client.PmaxClient.GetStorageGroup(ctx, r.client.SymmetrixID, plan.StorageGroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading storage group", err.Error())
		return
	}

	err = helper.CopyFields(ctx, storageGroup, &state)
	state.ID = types.StringValue(storageGroup.StorageGroupID)
	if err != nil {
		resp.Diagnostics.AddError("Error copying storage group fields", err.Error())
		return
	}

	//state.CapacityUnit = plan.CapacityUnit
	//state.VolumeIdentifierName = plan.VolumeIdentifierName
	//state.VolumeSize = plan.VolumeSize

	// Save plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *StorageGroup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Storage Group...")
	var state models.StorageGroupResourceModel

	// Read Terraform prior state into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	storageGroup, err := r.client.PmaxClient.GetStorageGroup(ctx, r.client.SymmetrixID, state.StorageGroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading storage group", err.Error())
		return
	}

	err = helper.CopyFields(ctx, storageGroup, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error copying storage group fields", err.Error())
		return
	}
	state.ID = types.StringValue(storageGroup.StorageGroupID)

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

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
	payload := pmaxTypes.UpdateStorageGroupPayload{ExecutionOption: pmaxTypes.ExecutionOptionSynchronous}

	// Rename
	planID := plan.StorageGroupID.ValueString()
	if stateID != planID {
		payload.EditStorageGroupActionParam = pmaxTypes.EditStorageGroupActionParam{
			RenameStorageGroupParam: &pmaxTypes.RenameStorageGroupParam{
				NewStorageGroupName: planID,
			},
		}
		err := r.client.PmaxClient.UpdateStorageGroupS(ctx, r.client.SymmetrixID, sgID, payload)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to update Storage Group ID(name): %s", err.Error()))
			resp.Diagnostics.AddError("Failed to update Storage Group ID(name)", err.Error())
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Update Storage Group ID(name): %s", planID))
			sgID = planID
			state.StorageGroupID = types.StringValue(planID)
		}
	}

	// Edit Compression
	planCompression := plan.Compression.ValueBool()
	stateCompression := state.Compression.ValueBool()
	if planCompression != stateCompression {
		payload.EditStorageGroupActionParam = pmaxTypes.EditStorageGroupActionParam{
			EditCompressionParam: &pmaxTypes.EditCompressionParam{
				Compression: &planCompression,
			},
		}
		err := r.client.PmaxClient.UpdateStorageGroupS(ctx, r.client.SymmetrixID, sgID, payload)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to update compression: %s", err.Error()))
			resp.Diagnostics.AddError("Failed to update compression", err.Error())
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Update compression: %t", planCompression))
			state.Compression = types.BoolValue(planCompression)
		}
	}

	// SetHostIOLimit
	planHostIOLimit := make(map[string]string)
	stateHostIOLimit := make(map[string]string)
	diag := plan.HostIOLimit.ElementsAs(ctx, &planHostIOLimit, true)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	diag = state.HostIOLimit.ElementsAs(ctx, &stateHostIOLimit, true)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(planHostIOLimit) > 0 && !reflect.DeepEqual(planHostIOLimit, stateHostIOLimit) {
		payload.EditStorageGroupActionParam = pmaxTypes.EditStorageGroupActionParam{
			SetHostIOLimitsParam: &pmaxTypes.SetHostIOLimitsParam{
				HostIOLimitIOSec:    planHostIOLimit["host_io_limit_io_sec"],
				HostIOLimitMBSec:    planHostIOLimit["host_io_limit_mb_sec"],
				DynamicDistribution: planHostIOLimit["dynamicDistribution"],
			},
		}
		err := r.client.PmaxClient.UpdateStorageGroupS(ctx, r.client.SymmetrixID, sgID, payload)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to update hostIOLimit: %s", err.Error()))
			resp.Diagnostics.AddError("Failed to update hostIOLimit", err.Error())
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Update hostIOLimit: %v", planHostIOLimit))
		}
	}

	// Edit Workload
	planWorkload := plan.Workload.ValueString()
	stateWorkload := state.Workload.ValueString()
	if planWorkload != stateWorkload {
		payload.EditStorageGroupActionParam = pmaxTypes.EditStorageGroupActionParam{
			EditStorageGroupWorkloadParam: &pmaxTypes.EditStorageGroupWorkloadParam{
				WorkloadSelection: planWorkload,
			},
		}
		err := r.client.PmaxClient.UpdateStorageGroupS(ctx, r.client.SymmetrixID, sgID, payload)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to update workload: %s", err.Error()))
			resp.Diagnostics.AddError("Failed to update workload", err.Error())
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Update workload: %s", planWorkload))
			state.Workload = types.StringValue(planWorkload)
		}
	}

	// Edit SLO
	planSLO := plan.SLO.ValueString()
	stateSLO := state.SLO.ValueString()
	if planSLO != stateSLO {
		payload.EditStorageGroupActionParam = pmaxTypes.EditStorageGroupActionParam{
			EditStorageGroupSLOParam: &pmaxTypes.EditStorageGroupSLOParam{
				SLOID: planSLO,
			},
		}
		err := r.client.PmaxClient.UpdateStorageGroupS(ctx, r.client.SymmetrixID, sgID, payload)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to update SLO: %s", err.Error()))
			resp.Diagnostics.AddError("Failed to update SLO", err.Error())
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Update SLO: %s", planSLO))
			state.SLO = types.StringValue(planSLO)
		}
	}

	// Edit SRP
	planSRP := plan.SRP.ValueString()
	stateSRP := state.SRP.ValueString()
	if planSRP != stateSRP {
		payload.EditStorageGroupActionParam = pmaxTypes.EditStorageGroupActionParam{
			EditStorageGroupSRPParam: &pmaxTypes.EditStorageGroupSRPParam{
				SRPID: planSRP,
			},
		}
		err := r.client.PmaxClient.UpdateStorageGroupS(ctx, r.client.SymmetrixID, sgID, payload)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to update SRP: %s", err.Error()))
			resp.Diagnostics.AddError("Failed to update SRP", err.Error())
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Update SRP: %s", planSLO))
			state.SRP = types.StringValue(planSRP)
		}
	}

	// Update Volume
	//err := UpdateVolume(ctx, &plan, &state, r)
	//if err != nil {
	//	resp.Diagnostics.AddError("Failed to update volume", err.Error())
	//	return
	//}

	// Update all fields of state
	storageGroup, err := r.client.PmaxClient.GetStorageGroup(ctx, r.client.SymmetrixID, state.StorageGroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading storage group", err.Error())
		return
	}

	err = helper.CopyFields(ctx, storageGroup, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error copying storage group fields", err.Error())
		return
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StorageGroup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Storage Group...")
	var data models.StorageGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.PmaxClient.DeleteStorageGroup(ctx, r.client.SymmetrixID, data.StorageGroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete storage group, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *StorageGroup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("storage_group_id"), req, resp)
}