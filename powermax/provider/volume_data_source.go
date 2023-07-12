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
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/helper"
	"terraform-provider-powermax/powermax/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &volumeDataSource{}
	_ datasource.DataSourceWithConfigure = &volumeDataSource{}
)

// NewVolumeDataSource returns the volume data source object.
func NewVolumeDataSource() datasource.DataSource {
	return &volumeDataSource{}
}

type volumeDataSource struct {
	client *client.Client
}

func (d *volumeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

func (d *volumeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for reading Volumes in PowerMax array.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder for acc testing",
				Computed:    true,
			},
			"volumes": schema.ListNestedAttribute{
				Description:         "List of volumes.",
				MarkdownDescription: "List of volumes.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "The ID of the volume.",
							MarkdownDescription: "The ID of the volume.",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							Description:         "The type of the volume.",
							MarkdownDescription: "The type of the volume.",
						},
						"volume_identifier": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							Description:         "The identifier of the volume.",
							MarkdownDescription: "The identifier of the volume.",
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
						"physical_name": schema.StringAttribute{
							Computed:            true,
							Description:         "The physical name of the volume.",
							MarkdownDescription: "The physical name of the volume.",
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
						},
						"mobility_id_enabled": schema.BoolAttribute{
							Computed:            true,
							Optional:            true,
							Description:         "States whether mobility ID is enabled on the volume.",
							MarkdownDescription: "States whether mobility ID is enabled on the volume.",
						},
						"unreducible_data_gb": schema.Float64Attribute{
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
						},
						"cap_gb": schema.Float64Attribute{
							Computed:            true,
							Description:         "The capability of volume in the unit of GB.",
							MarkdownDescription: "The capability of volume in the unit of GB.",
						},
						"cap_mb": schema.Float64Attribute{
							Computed:            true,
							Description:         "The capability of volume in the unit of MB.",
							MarkdownDescription: "The capability of volume in the unit of MB.",
						},
						"cap_cyl": schema.Int64Attribute{
							Computed:            true,
							Description:         "The capability of volume in the unit of CYL.",
							MarkdownDescription: "The capability of volume in the unit of CYL.",
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"storage_group_name": schema.StringAttribute{
						Description:         "The name of the storage group.",
						MarkdownDescription: "The name of the storage group.",
						Optional:            true,
					},
					"encapsulated_wwn": schema.StringAttribute{
						Description:         "The specified volume encapsulated_wwn.",
						MarkdownDescription: "The specified volume encapsulated_wwn.",
						Optional:            true,
					},
					"wwn": schema.StringAttribute{
						Description:         "The specified volume wwn.",
						MarkdownDescription: "The specified volume wwn.",
						Optional:            true,
					},
					"symmlun": schema.StringAttribute{
						Description:         "Greater than, Less than or equal to the specified symmlun.",
						MarkdownDescription: "Greater than, Less than or equal to the specified symmlun.",
						Optional:            true,
					},
					"status": schema.StringAttribute{
						Description:         "The specified volume status.",
						MarkdownDescription: "The specified volume status.",
						Optional:            true,
					},
					"physical_name": schema.StringAttribute{
						Description:         "The specified volume physical name.",
						MarkdownDescription: "The specified volume physical name.",
						Optional:            true,
					},
					"volume_identifier": schema.StringAttribute{
						Description:         "The specified volume volume identifier.",
						MarkdownDescription: "The specified volume volume identifier.",
						Optional:            true,
					},
					"allocated_percent": schema.StringAttribute{
						Description:         "Greater than, Less than or equal to the allocated percent.",
						MarkdownDescription: "Greater than, Less than or equal to the allocated percent.",
						Optional:            true,
					},
					"cap_tb": schema.StringAttribute{
						Description:         "Greater than, Less than or equal to the cap tb.",
						MarkdownDescription: "Greater than, Less than or equal to the cap tb.",
						Optional:            true,
					},
					"cap_gb": schema.StringAttribute{
						Description:         "Greater than, Less than or equal to the cap gb.",
						MarkdownDescription: "Greater than, Less than or equal to the cap gb.",
						Optional:            true,
					},
					"cap_mb": schema.StringAttribute{
						Description:         "Greater than, Less than or equal to the cap mb.",
						MarkdownDescription: "Greater than, Less than or equal to the cap mb.",
						Optional:            true,
					},
					"cap_cyl": schema.StringAttribute{
						Description:         "Greater than, Less than or equal to the cap CYL.",
						MarkdownDescription: "Greater than, Less than or equal to the cap CYL.",
						Optional:            true,
					},
					"num_of_storage_groups": schema.StringAttribute{
						Description:         "Greater than, Less than or equal to the number of storage groups.",
						MarkdownDescription: "Greater than, Less than or equal to the number of storage groups.",
						Optional:            true,
					},
					"num_of_masking_views": schema.StringAttribute{
						Description:         "Greater than, Less than or equal to the number of masking views.",
						MarkdownDescription: "Greater than, Less than or equal to the number of masking views.",
						Optional:            true,
					},
					"num_of_front_end_paths": schema.StringAttribute{
						Description:         "Greater than, Less than or equal to the number of front end paths.",
						MarkdownDescription: "Greater than, Less than or equal to the number of front end paths.",
						Optional:            true,
					},
					"virtual_volumes": schema.BoolAttribute{
						Description:         "Volumes that are virtual volumes (true/false).",
						MarkdownDescription: "Volumes that are virtual volumes (true/false).",
						Optional:            true,
					},
					"private_volumes": schema.BoolAttribute{
						Description:         "Volumes that are private (true/false).",
						MarkdownDescription: "Volumes that are private (true/false).",
						Optional:            true,
					},
					"available_thin_volumes": schema.BoolAttribute{
						Description:         "Volumes that are available thin volumes (true/false).",
						MarkdownDescription: "Volumes that are available thin volumes (true/false).",
						Optional:            true,
					},
					"tdev": schema.BoolAttribute{
						Description:         "Volumes that are tdev (true/false).",
						MarkdownDescription: "Volumes that are tdev (true/false).",
						Optional:            true,
					},
					"thin_bcv": schema.BoolAttribute{
						Description:         "Volumes that are thin bcv (true/false).",
						MarkdownDescription: "Volumes that are thin bcv (true/false).",
						Optional:            true,
					},
					"vdev": schema.BoolAttribute{
						Description:         "Volumes that are vdev (true/false).",
						MarkdownDescription: "Volumes that are vdev (true/false).",
						Optional:            true,
					},
					"gatekeeper": schema.BoolAttribute{
						Description:         "Volumes that are gatekeeper (true/false).",
						MarkdownDescription: "Volumes that are gatekeeper (true/false).",
						Optional:            true,
					},
					"data_volume": schema.BoolAttribute{
						Description:         "Volumes that are data volume (true/false).",
						MarkdownDescription: "Volumes that are data volume (true/false).",
						Optional:            true,
					},
					"dld": schema.BoolAttribute{
						Description:         "Volumes that are dld (true/false).",
						MarkdownDescription: "Volumes that are dld (true/false).",
						Optional:            true,
					},
					"drv": schema.BoolAttribute{
						Description:         "Volumes that are drv (true/false).",
						MarkdownDescription: "Volumes that are drv (true/false).",
						Optional:            true,
					},
					"mapped": schema.BoolAttribute{
						Description:         "Volumes that are mapped (true/false).",
						MarkdownDescription: "Volumes that are mapped (true/false).",
						Optional:            true,
					},
					"bound_tdev": schema.BoolAttribute{
						Description:         "Volumes that are bound tdev (true/false).",
						MarkdownDescription: "Volumes that are bound tdev (true/false).",
						Optional:            true,
					},
					"reserved": schema.BoolAttribute{
						Description:         "Volumes that are reserved (true/false).",
						MarkdownDescription: "Volumes that are reserved (true/false).",
						Optional:            true,
					},
					"pinned": schema.BoolAttribute{
						Description:         "Volumes that are pinned (true/false).",
						MarkdownDescription: "Volumes that are pinned (true/false).",
						Optional:            true,
					},
					"encapsulated": schema.BoolAttribute{
						Description:         "Volumes that are encapsulated (true/false).",
						MarkdownDescription: "Volumes that are encapsulated (true/false).",
						Optional:            true,
					},
					"associated": schema.BoolAttribute{
						Description:         "Volumes that are associated (true/false).",
						MarkdownDescription: "Volumes that are associated (true/false).",
						Optional:            true,
					},
					"emulation": schema.StringAttribute{
						Description:         "Volumes that are of the specified emulation.",
						MarkdownDescription: "Volumes that are of the specified emulation.",
						Optional:            true,
					},
					"split_name": schema.StringAttribute{
						Description:         "Volumes that are mapped to CU images associated to the specified FICON split.",
						MarkdownDescription: "Volumes that are mapped to CU images associated to the specified FICON split.",
						Optional:            true,
					},
					"cu_image_num": schema.StringAttribute{
						Description:         "Volumes that are mapped to a CU image with the specified CU image number.",
						MarkdownDescription: "Volumes that are mapped to a CU image with the specified CU image number.",
						Optional:            true,
					},
					"cu_image_ssid": schema.StringAttribute{
						Description:         "Volumes that are mapped to a CU image with the specified CU SSID.",
						MarkdownDescription: "Volumes that are mapped to a CU image with the specified CU SSID.",
						Optional:            true,
					},
					"rdf_group_number": schema.StringAttribute{
						Description:         "Volumes that are part of the specified rdf group.",
						MarkdownDescription: "Volumes that are part of the specified rdf group.",
						Optional:            true,
					},
					"has_effective_wwn": schema.BoolAttribute{
						Description:         "Volumes that have effective wwns (true/false)",
						MarkdownDescription: "Volumes that have effective wwns (true/false)",
						Optional:            true,
					},
					"effective_wwn": schema.StringAttribute{
						Description:         "Volumes that contain the specified effective_wwn.",
						MarkdownDescription: "Volumes that contain the specified effective_wwn.",
						Optional:            true,
					},
					"type": schema.StringAttribute{
						Description:         "Volumes that contain the specified volume type.",
						MarkdownDescription: "Volumes that contain the specified volume type.",
						Optional:            true,
					},
					"oracle_instance_name": schema.StringAttribute{
						Description:         "Volumes that contain the specified Oracle Instance Name.",
						MarkdownDescription: "Volumes that contain the specified Oracle Instance Name.",
						Optional:            true,
					},
					"mobility_id_enabled": schema.BoolAttribute{
						Description:         "Volumes that are mobility ID enabled (true/false).",
						MarkdownDescription: "Volumes that are mobility ID enabled (true/false).",
						Optional:            true,
					},
					"unreducible_data_gb": schema.StringAttribute{
						Description:         "Greater than,Less than or equal to the unreducible data gb.",
						MarkdownDescription: "Greater than,Less than or equal to the unreducible data gb.",
						Optional:            true,
					},
					"nguid": schema.StringAttribute{
						Description:         "Volumes that correspond to Namespace Globally Unique Identifier that uses the EUI64 16-byte designator format.",
						MarkdownDescription: "Volumes that correspond to Namespace Globally Unique Identifier that uses the EUI64 16-byte designator format.",
						Optional:            true,
					},
				},
			},
		},
	}
}

func (d *volumeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	pmaxClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = pmaxClient
}

func (d *volumeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.VolumeDatasource
	var err error

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	param, err := helper.GetVolumeFilterParam(ctx, d.client, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get volume filter param",
			err.Error(),
		)
		return
	}
	state.Volumes, err = updateVolumeState(ctx, d.client, param)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update volume state",
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue("place_holder")
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// updateVolumeState iterates over the volume list and update the state.
func updateVolumeState(ctx context.Context, p *client.Client, params powermax.ApiListVolumesRequest) (response []models.VolumeDatasourceEntity, err error) {
	volIDs, _, err := params.Execute()
	if err != nil {
		errStr := ""
		message := helper.GetErrorString(err, errStr)
		return nil, fmt.Errorf(message)
	}

	for _, vol := range volIDs.ResultList.GetResult() {
		for _, volumeID := range vol {
			volumeModel := p.PmaxOpenapiClient.SLOProvisioningApi.GetVolume(ctx, p.SymmetrixID, fmt.Sprint(volumeID))
			volResponse, _, err := volumeModel.Execute()
			if err != nil {
				errStr := ""
				message := helper.GetErrorString(err, errStr)
				return nil, fmt.Errorf(message)

			}
			volState := models.VolumeDatasourceEntity{}
			err = helper.CopyFields(ctx, volResponse, &volState)
			volState.SymmetrixPortKey, _ = helper.GetSymmetrixPortKeyObjects(volResponse)
			volState.StorageGroups, _ = helper.GetStorageGroupObjects(volResponse)
			volState.RfdGroupIDList, _ = helper.GetRfdGroupIdsObjects(volResponse)
			if id, ok := volResponse.GetVolumeIdOk(); ok {
				volState.VolumeID = types.StringValue(*id)
			}
			if mobid, ok := volResponse.GetMobilityIdEnabledOk(); ok {
				volState.MobilityIDEnabled = types.BoolValue(*mobid)
			}
			if err != nil {
				return nil, err
			}
			response = append(response, volState)
		}
	}
	return response, nil
}
