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
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/constants"
	"terraform-provider-powermax/powermax/helper"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SnapshotPolicy{}
var _ resource.ResourceWithImportState = &SnapshotPolicy{}
var _ resource.ResourceWithConfigure = &SnapshotPolicy{}

// NewSnapshotPolicy creates a new Snapshot Policy resource.
func NewSnapshotPolicy() resource.Resource {
	return &SnapshotPolicy{}
}

// SnapshotPolicy defines the resource implementation.
type SnapshotPolicy struct {
	client *client.Client
}

// Metadata returns the metadata for the resource.
func (r *SnapshotPolicy) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshotpolicy"
}

// Schema returns the schema for the resource.
func (r *SnapshotPolicy) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Resource for a specific Snapshot Policy in PowerMax array.",
		Description:         "Resource for a specific Snapshot Policy in PowerMax array.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier",
				Computed:    true,
			},
			"snapshot_policy_name": schema.StringAttribute{
				Description:         "Name of the snapshot policy. Only alphanumeric characters, underscores ( _ ), and hyphens (-) are allowed and max length can be 32 characters",
				MarkdownDescription: "Name of the snapshot policy. Only alphanumeric characters, underscores ( _ ), and hyphens (-) are allowed and max length can be 32 characters",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(32),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9_-]*$`),
						"must contain only alphanumeric characters and _-",
					),
				},
			},
			"snapshot_count": schema.Int64Attribute{
				Description:         "Number of snapshots that will be taken before the oldest ones are no longer required",
				MarkdownDescription: "Number of snapshots that will be taken before the oldest ones are no longer required",
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(48),
			},
			"interval_minutes": schema.Int64Attribute{
				Description:         "Number of minutes between each policy execution",
				MarkdownDescription: "Number of minutes between each policy execution",
				Computed:            true,
				Optional:            true,
			},
			"offset_minutes": schema.Int64Attribute{
				Description:         "Number of minutes after 00:00 on Monday morning that the policy will execute",
				MarkdownDescription: "Number of minutes after 00:00 on Monday morning that the policy will execute",
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(420),
			},
			"provider_name": schema.StringAttribute{
				Description:         "The name of the cloud provider associated with this policy. Only applies to cloud policies",
				MarkdownDescription: "The name of the cloud provider associated with this policy. Only applies to cloud policies",
				Computed:            true,
				Optional:            true,
			},
			"retention_days": schema.Int64Attribute{
				Description:         "The number of days that snapshots will be retained in the cloud for. Only applies to cloud policies",
				MarkdownDescription: "The number of days that snapshots will be retained in the cloud for. Only applies to cloud policies",
				Computed:            true,
				Optional:            true,
			},
			"suspended": schema.BoolAttribute{
				Description:         "Set if the snapshot policy has been suspended",
				MarkdownDescription: "Set if the snapshot policy has been suspended",
				Computed:            true,
				Optional:            true,
			},
			"secure": schema.BoolAttribute{
				Description:         "Set if the snapshot policy creates secure snapshots",
				MarkdownDescription: "Set if the snapshot policy creates secure snapshots",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"last_time_used": schema.StringAttribute{
				Description:         "The last time that the snapshot policy was run",
				MarkdownDescription: "The last time that the snapshot policy was run",
				Computed:            true,
				Optional:            true,
			},
			"storage_group_count": schema.Int64Attribute{
				Description:         "The total number of storage groups that this snapshot policy is associated with",
				MarkdownDescription: "The total number of storage groups that this snapshot policy is associated with",
				Computed:            true,
			},
			"compliance_count_warning": schema.Int64Attribute{
				Description:         "The threshold of good snapshots which are not failed/bad for compliance to change from normal to warning.",
				MarkdownDescription: "The threshold of good snapshots which are not failed/bad for compliance to change from normal to warning.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(47),
			},
			"compliance_count_critical": schema.Int64Attribute{
				Description:         "The threshold of good snapshots which are not failed/bad for compliance to change from warning to critical",
				MarkdownDescription: "The threshold of good snapshots which are not failed/bad for compliance to change from warning to critical",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(46),
			},
			"type": schema.StringAttribute{
				Description:         "The type of Snapshots that are created with the policy, local or cloud",
				MarkdownDescription: "The type of Snapshots that are created with the policy, local or cloud",
				Computed:            true,
			},
			"interval": schema.StringAttribute{
				Description:         "The interval between snapshots Enumeration values: 10 Minutes, 12 Minutes, 15 Minutes, 20 Minutes, 30 Minutes, 1 Hour, 2 Hours, 3 Hours, 4 Hours, 6 Hours, 8 Hours, 12 Hours, 1 Day, 7 Days",
				MarkdownDescription: "The interval between snapshots Enumeration values: 10 Minutes, 12 Minutes, 15 Minutes, 20 Minutes, 30 Minutes, 1 Hour, 2 Hours, 3 Hours, 4 Hours, 6 Hours, 8 Hours, 12 Hours, 1 Day, 7 Days",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("1 Hour"),
				Validators: []validator.String{
					stringvalidator.OneOf("10 Minutes", "12 Minutes", "15 Minutes", "20 Minutes", "30 Minutes", "1 Hour", "2 Hours", "3 Hours", "4 Hours", "6 Hours", "8 Hours", "12 Hours", "1 Day", "7 Days"),
				},
			},
			"storage_groups": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "The storage groups associated with the snapshot policy.This field cannot be set during create and is only valid for Edit/Update.If user wants to delete the snapshot policy all associated storage groups will also be unlinked from the Snapshot Policy.",
				MarkdownDescription: "The storage groups associated with the snapshot policy..This field cannot be set during create and is only valid for Edit/Update.If user wants to delete the snapshot policy all associated storage groups will also be unlinked from the Snapshot Policy.",
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
		},
	}
}

// Configure configure client for Snapshot policy resource.
func (r *SnapshotPolicy) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a snapshot policy and refresh state.
func (r *SnapshotPolicy) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Snapshot Policy...")
	var planSnapPolicy models.SnapshotPolicyResource
	diags := req.Plan.Get(ctx, &planSnapPolicy)
	// Read Terraform plan into the model
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !planSnapPolicy.StorageGroups.IsNull() && len(planSnapPolicy.StorageGroups.Elements()) > 0 {
		resp.Diagnostics.AddError(
			"Unable to create snapshot policy",
			"Create snapshot policy does not support adding storage groups, only after the policy is created can you add storage groups",
		)
		return
	}

	snapPolicyCreateResp, _, err := helper.CreateSnapshotPolicy(ctx, *r.client, planSnapPolicy)
	if err != nil {
		snapPolicyID := planSnapPolicy.SnapshotPolicyName.ValueString()
		errStr := constants.CreateSnapPolicyDetailErrorMsg + snapPolicyID + ": "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating snapshot policy",
			message,
		)

		req := r.client.PmaxOpenapiClient.ReplicationApi.GetSnapshotPolicy(ctx, r.client.SymmetrixID, snapPolicyID)
		snapPolicyGetResp, _, getSnapPolicyErr := req.Execute()
		if snapPolicyGetResp != nil || getSnapPolicyErr == nil {
			_, err := helper.DeleteSnapshotPolicy(ctx, *r.client, snapPolicyID)
			if err != nil {
				errStr := constants.CreateSnapPolicyDetailErrorMsg + snapPolicyID + "with error: "
				message := helper.GetErrorString(err, errStr)
				resp.Diagnostics.AddError(
					"Error deleting the invalid snapshot policy, This may be a dangling resource and needs to be deleted manually",
					message,
				)
			}
		}
		return
	}
	tflog.Debug(ctx, "create snapshot policy response", map[string]interface{}{
		"Create Snapshot Policy Response": snapPolicyCreateResp,
	})
	//Get Storage Groups associated with the snapshot policy
	storageGroups, _, errStorageGroup := helper.GetSnapshotPolicyStorageGroups(ctx, *r.client, planSnapPolicy.SnapshotPolicyName.ValueString())
	if errStorageGroup != nil {
		errStr := ""
		msgStr := helper.GetErrorString(errStorageGroup, errStr)
		resp.Diagnostics.AddError("Error getting Snapshot Policy storage groups", msgStr)
		// Attempt to cleanup after failure
		_, err := helper.DeleteSnapshotPolicy(ctx, *r.client, planSnapPolicy.SnapshotPolicyName.ValueString())
		if err != nil {
			errStr := constants.CreateSnapPolicyDetailErrorMsg + planSnapPolicy.SnapshotPolicyName.ValueString() + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError(
				"Error deleting the invalid snapshot policy, This may be a dangling resource and needs to be deleted manually",
				message,
			)
		}
		return
	}

	var result models.SnapshotPolicyResource
	// Copy values with the same fields
	errCpy := helper.UpdateSnapshotPolicyResourceState(ctx, snapPolicyCreateResp, &result, storageGroups)

	if errCpy != nil {
		resp.Diagnostics.AddError("Error copying Snapshot Policy", errCpy.Error())
		// Attempt to cleanup after failure
		_, err := helper.DeleteSnapshotPolicy(ctx, *r.client, planSnapPolicy.SnapshotPolicyName.ValueString())
		if err != nil {
			errStr := constants.CreateSnapPolicyDetailErrorMsg + planSnapPolicy.SnapshotPolicyName.ValueString() + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError(
				"Error deleting the invalid snapshot policy, This may be a dangling resource and needs to be deleted manually",
				message,
			)
		}
		return
	}
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create snapshot policy completed")

}

// Delete SnapshotPolicy.
func (r *SnapshotPolicy) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting SnapshotPolicy")
	var snapPolicyState models.SnapshotPolicyResource
	diags := req.State.Get(ctx, &snapPolicyState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	snapPolicyID := snapPolicyState.SnapshotPolicyName.ValueString()

	// Remove any associated storage groups from snapshot policy before deleting the snapshot policy
	if !snapPolicyState.StorageGroups.IsNull() && len(snapPolicyState.StorageGroups.Elements()) > 0 {
		remSGs := make([]string, 0, len(snapPolicyState.StorageGroups.Elements()))
		for _, sg := range snapPolicyState.StorageGroups.Elements() {
			remSGs = append(remSGs, sg.String()[1:len(sg.String())-1])
		}
		removeSnapshotPolicyParam := pmax.NewSnapshotPolicyStorageGroupAddRemove()
		removeSnapshotPolicyParam.SetStorageGroupName(remSGs)
		snapshotPolicyUpdate := pmax.SnapshotPolicyUpdate{
			Action:                       "DisassociateFromStorageGroups",
			DisassociateFromStorageGroup: removeSnapshotPolicyParam,
		}

		updateReq := r.client.PmaxOpenapiClient.ReplicationApi.UpdateSnapshotPolicy(ctx, r.client.SymmetrixID, snapPolicyID)
		updateReq = updateReq.SnapshotPolicyUpdate(snapshotPolicyUpdate)
		_, _, err := updateReq.Execute()

		if err != nil {
			errStr := constants.DeleteSnapPolicyDetailErrorMsg + snapPolicyID + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Could not remove associated storage groups from Snapshot Policy", message)
			return
		}
	}
	tflog.Debug(ctx, "deleting snapshot policy by snapPolicyId", map[string]interface{}{
		"symmetrixID":  r.client.SymmetrixID,
		"snapPolicyID": snapPolicyID,
	})
	delReq := r.client.PmaxOpenapiClient.ReplicationApi.DeleteSnapshotPolicy(ctx, r.client.SymmetrixID, snapPolicyID)
	_, err := delReq.Execute()
	if err != nil {
		errStr := constants.DeleteSnapPolicyDetailErrorMsg + snapPolicyID + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error deleting snapshot policy",
			message,
		)
	}

	tflog.Info(ctx, "Delete snapshot policy complete")
}

// Update Snapshot Policy.
func (r *SnapshotPolicy) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "updating Snapshot Policy")
	var plan models.SnapshotPolicyResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "fetched snapshot policy details from plan")

	var state models.SnapshotPolicyResource
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update host on pmax client", map[string]interface{}{
		"plan":  plan,
		"state": state,
	})

	err := helper.ModifySnapshotPolicy(ctx, *r.client, &plan, &state)
	if err != nil {
		errStr := constants.UpdateSnapshotPolicy + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating snapshot Policy",
			message,
		)
		return
	}
	// Read and update state after the modification
	getReq := r.client.PmaxOpenapiClient.ReplicationApi.GetSnapshotPolicy(ctx, r.client.SymmetrixID, plan.SnapshotPolicyName.ValueString())
	snapPolicyDetail, _, err := getReq.Execute()
	if err != nil {
		errStr := fmt.Sprintf("Error reading snapshot policy %s after update with error:", state.SnapshotPolicyName.ValueString())
		msgStr := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading snapshot policy",
			msgStr,
		)
		return
	}
	// Get Storage Groups associated with the snapshot policy
	storageGroupReq := r.client.PmaxOpenapiClient.ReplicationApi.GetSnapshotPolicyStorageGroups(ctx, r.client.SymmetrixID, snapPolicyDetail.SnapshotPolicyName)
	storageGroups, _, errStorageGroup := storageGroupReq.Execute()
	if errStorageGroup != nil {
		errStr := ""
		msgStr := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error getting Snapshot Policy storage groups", msgStr)
	}

	errState := helper.UpdateSnapshotPolicyResourceState(ctx, snapPolicyDetail, &state, storageGroups)
	if errState != nil {
		resp.Diagnostics.AddError(
			"Error updating snapshot policy state",
			errState.Error(),
		)
		return
	}

	// Save plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "update snapshot policy completed")
}

// Read SnapshotPolicy.
func (r *SnapshotPolicy) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading SnapshotPolicy...")
	var snapPolicyState models.SnapshotPolicyResource
	diags := req.State.Get(ctx, &snapPolicyState)
	// Read Terraform prior state into the model
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	snapshotPolicyID := snapPolicyState.SnapshotPolicyName.ValueString()
	snapshotPolicy, _, err := helper.GetSnapshotPolicy(ctx, *r.client, snapshotPolicyID)
	if err != nil {
		errStr := constants.ReadSnapPolicyDetailsErrorMsg + snapshotPolicyID + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading snapshot policy",
			message,
		)
		return
	}
	// Get Storage Groups associated with the snapshot policy
	storageGroups, _, errStorageGroup := helper.GetSnapshotPolicyStorageGroups(ctx, *r.client, snapshotPolicyID)
	if errStorageGroup != nil {
		errStr := ""
		msgStr := helper.GetErrorString(errStorageGroup, errStr)
		resp.Diagnostics.AddError("Error getting snapshot policy storage groups", msgStr)
	}

	tflog.Debug(ctx, "Updating snapshot policy state")
	// Copy values with the same fields
	errCpy := helper.UpdateSnapshotPolicyResourceState(ctx, snapshotPolicy, &snapPolicyState, storageGroups)

	if errCpy != nil {
		resp.Diagnostics.AddError("Error reading snapshot policy", errCpy.Error())
		return
	}
	diags = resp.State.Set(ctx, snapPolicyState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Read SnapshotPolicy completed")

}

// ImportState imports the state of the resource from the req.
func (r *SnapshotPolicy) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "importing Snapshot Policy state")
	var snapPolicyState models.SnapshotPolicyResource
	snapshotPolicyID := req.ID
	tflog.Debug(ctx, "fetching snapshot policy by ID", map[string]interface{}{
		"symmetrixID":      r.client.SymmetrixID,
		"snapshotPolicyID": snapshotPolicyID,
	})

	getReq := r.client.PmaxOpenapiClient.ReplicationApi.GetSnapshotPolicy(ctx, r.client.SymmetrixID, snapshotPolicyID)
	snapshotPolicyResponse, _, err := getReq.Execute()

	if err != nil {
		errStr := constants.ImportHostDetailsErrorMsg + snapshotPolicyID + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading snapshot policy",
			message,
		)
		return
	}
	tflog.Debug(ctx, "Get snapshot policy By ID response", map[string]interface{}{
		"Snapshot Policy Response": snapshotPolicyResponse,
	})
	// Get Storage Groups associated with the snapshot policy
	storageGroupReq := r.client.PmaxOpenapiClient.ReplicationApi.GetSnapshotPolicyStorageGroups(ctx, r.client.SymmetrixID, snapshotPolicyID)
	storageGroups, _, errStorageGroup := storageGroupReq.Execute()
	if errStorageGroup != nil {
		errStr := ""
		msgStr := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error getting Snapshot Policy storage groups", msgStr)
	}

	tflog.Debug(ctx, "updating snapshot policy state after import")
	errCpy := helper.UpdateSnapshotPolicyResourceState(ctx, snapshotPolicyResponse, &snapPolicyState, storageGroups)
	if errCpy != nil {
		resp.Diagnostics.AddError("Error copying Snapshot Policy", errCpy.Error())
		return
	}

	diags := resp.State.Set(ctx, snapPolicyState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "import Snapshot Policy state completed")
}
