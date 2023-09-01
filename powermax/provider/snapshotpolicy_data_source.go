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
	"fmt"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/helper"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &snapshotPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &snapshotPolicyDataSource{}
)

// NewSnapshotPolicyDataSource is a helper function to simplify the provider implementation.
func NewSnapshotPolicyDataSource() datasource.DataSource {
	return &snapshotPolicyDataSource{}
}

// snapshotPolicyDataSource is the data source implementation.
type snapshotPolicyDataSource struct {
	client *client.Client
}

func (d *snapshotPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshotpolicy"
}

func (d *snapshotPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source for a specific Snapshot Policy in PowerMax array. PowerMax snapshot policy feature provides snapshot orchestration at scale (1,024 snaps per storage group). The resource simplifies snapshot management for standard and cloud snapshots.",
		Description:         "Data source for a specific Snapshot Policy in PowerMax array. PowerMax snapshot policy feature provides snapshot orchestration at scale (1,024 snaps per storage group). The resource simplifies snapshot management for standard and cloud snapshots.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier",
				Computed:    true,
			},
			"snapshot_policies": schema.ListNestedAttribute{
				Description:         "List of Snapshot Policies",
				MarkdownDescription: "List of Snapshot Policies",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"snapshot_policy_name": schema.StringAttribute{
							Description:         "Name of the snapshot policy",
							MarkdownDescription: "Name of the snapshot policy",
							Computed:            true,
						},
						"snapshot_count": schema.Int64Attribute{
							Description:         "Number of snapshots that will be taken before the oldest ones are no longer required",
							MarkdownDescription: "Number of snapshots that will be taken before the oldest ones are no longer required",
							Computed:            true,
						},
						"interval_minutes": schema.Int64Attribute{
							Description:         "Number of minutes between each policy execution",
							MarkdownDescription: "Number of minutes between each policy execution",
							Computed:            true,
						},
						"offset_minutes": schema.Int64Attribute{
							Description:         "Number of minutes after 00:00 on Monday morning that the policy will execute",
							MarkdownDescription: "Number of minutes after 00:00 on Monday morning that the policy will execute",
							Computed:            true,
						},
						"provider_name": schema.StringAttribute{
							Description:         "The name of the cloud provider associated with this policy. Only applies to cloud policies",
							MarkdownDescription: "The name of the cloud provider associated with this policy. Only applies to cloud policies",
							Computed:            true,
						},
						"retention_days": schema.Int64Attribute{
							Description:         "The number of days that snapshots will be retained in the cloud for. Only applies to cloud policies",
							MarkdownDescription: "The number of days that snapshots will be retained in the cloud for. Only applies to cloud policies",
							Computed:            true,
						},
						"suspended": schema.BoolAttribute{
							Description:         "Set if the snapshot policy has been suspended",
							MarkdownDescription: "Set if the snapshot policy has been suspended",
							Computed:            true,
						},
						"secure": schema.BoolAttribute{
							Description:         "Set if the snapshot policy creates secure snapshots",
							MarkdownDescription: "Set if the snapshot policy creates secure snapshots",
							Computed:            true,
						},
						"last_time_used": schema.StringAttribute{
							Description:         "The last time that the snapshot policy was run",
							MarkdownDescription: "The last time that the snapshot policy was run",
							Computed:            true,
						},
						"storage_group_count": schema.Int64Attribute{
							Description:         "The total number of storage groups that this snapshot policy is associated with",
							MarkdownDescription: "The total number of storage groups that this snapshot policy is associated with",
							Computed:            true,
						},
						"compliance_count_warning": schema.Int64Attribute{
							Description:         "The threshold of good snapshots which are not failed/bad for compliance to change from normal to warning.",
							MarkdownDescription: "The threshold of good snapshots which are not failed/bad for compliance to change from normal to warning.",
							Computed:            true,
						},
						"compliance_count_critical": schema.Int64Attribute{
							Description:         "The threshold of good snapshots which are not failed/bad for compliance to change from warning to critical",
							MarkdownDescription: "The threshold of good snapshots which are not failed/bad for compliance to change from warning to critical",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							Description:         "The type of Snapshots that are created with the policy, local or cloud",
							MarkdownDescription: "The type of Snapshots that are created with the policy, local or cloud",
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"names": schema.SetAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
		},
	}
}

func (d *snapshotPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = client
}

func (d *snapshotPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.SnapshotPolicyDataSourceModel

	tflog.Info(ctx, "Attempting to read snapshot policies")
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var snapshotPolicyIds []string
	// Get snapshot policy IDs from config or query all if not specified
	if state.SnapshotPolicyFilter == nil || len(state.SnapshotPolicyFilter.Names) == 0 {
		// Read all the snapshot policies
		snapshotPolicyList, _, err := helper.GetSnapshotPolicies(ctx, *d.client)
		if err != nil {
			errStr := ""
			msgStr := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Error reading Snapshot Policy ids", msgStr)
			return
		}
		snapshotPolicyIds = snapshotPolicyList.Name
	} else {
		// get ids from filter and assign to snapshotPolicyIds
		for _, ids := range state.SnapshotPolicyFilter.Names {
			snapshotPolicyIds = append(snapshotPolicyIds, ids.ValueString())
		}
	}
	for _, id := range snapshotPolicyIds {
		snapshotPolicyResponse, _, err := helper.GetSnapshotPolicy(ctx, *d.client, id)
		if err != nil || snapshotPolicyResponse == nil {
			errStr := ""
			msgStr := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Error reading snapshot policy with id", msgStr)
			continue
		}
		var snapshotPolicy models.SnapshotPolicyModel
		tflog.Debug(ctx, "Updating snapshot policy state")
		// Copy values with the same fields
		errCpy := helper.CopyFields(ctx, snapshotPolicyResponse, &snapshotPolicy)
		if errCpy != nil {
			resp.Diagnostics.AddError("Error copying Snapshot Policies", errCpy.Error())
			return
		}
		state.SnapshotPolicies = append(state.SnapshotPolicies, snapshotPolicy)
	}
	state.ID = types.StringValue("snapshot-policy-datasource")
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}
