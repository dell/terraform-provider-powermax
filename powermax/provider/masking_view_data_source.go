// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.

package provider

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"terraform-provider-powermax/client"
	"terraform-provider-powermax/powermax/models"

	pmaxTypes "github.com/dell/gopowermax/v2/types/v100"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &maskingViewDataSource{}
	_ datasource.DataSourceWithConfigure = &maskingViewDataSource{}
)

var lockMutex sync.Mutex

// defaultMaxPowerMaxConnections is the number of workers that can query powermax at a time.
const defaultMaxPowerMaxConnections = 10

// NewMaskingViewDataSource returns the masking view data source object.
func NewMaskingViewDataSource() datasource.DataSource {
	return &maskingViewDataSource{}
}

// maskingViewDataSource defines the data source implementation.
type maskingViewDataSource struct {
	client *client.Client
}

func (d *maskingViewDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_maskingview"
}

func (d *maskingViewDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Masking View data source",
		Description:         "Masking View data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Unique identifier of the masking view instance.",
				MarkdownDescription: "Unique identifier of the masking view instance.",
			},
			"masking_views": schema.ListNestedAttribute{
				Description:         "List of masking views.",
				MarkdownDescription: "List of masking views.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"masking_view_name": schema.StringAttribute{
							Description:         "Unique identifier of the masking view.",
							MarkdownDescription: "Unique identifier of the masking view.",
							Computed:            true,
						},
						"host_id": schema.StringAttribute{
							Description:         "The host id of the masking view.",
							MarkdownDescription: "The host id of the masking view.",
							Computed:            true,
						},
						"host_group_id": schema.StringAttribute{
							Description:         "The host group id of the masking view.",
							MarkdownDescription: "The host group id of the masking view.",
							Computed:            true,
						},
						"port_group_id": schema.StringAttribute{
							Description:         "The port group id of the masking view.",
							MarkdownDescription: "The port group id of the masking view.",
							Computed:            true,
						},
						"storage_group_id": schema.StringAttribute{
							Description:         "The storage group id of the masking view.",
							MarkdownDescription: "The storage group id of the masking view.",
							Computed:            true,
						},
						"capacity_gb": schema.Float64Attribute{
							Computed:            true,
							Description:         "The capacity of the storage group in the masking view.",
							MarkdownDescription: "The capacity of the storage group in the masking view.",
						},
						"volumes": schema.ListAttribute{
							Description:         "List of Volumes.",
							MarkdownDescription: "List of Volumes.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"ports": schema.ListAttribute{
							Description:         "List of ports.",
							MarkdownDescription: "List of ports.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"initiators": schema.ListAttribute{
							Description:         "List of initiators.",
							MarkdownDescription: "List of initiators.",
							ElementType:         types.StringType,
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

func (d *maskingViewDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (d *maskingViewDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Masking View data source ")

	var state models.MaskingViewDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var maskingViewIds []string
	// Get masking view IDs from config or query all if not specified
	if state.MaskingViewFilter == nil || len(state.MaskingViewFilter.Names) == 0 {
		// Read all the masking views
		tflog.Debug(ctx, fmt.Sprintf("Calling api to get MaskingViewList for Symmetrix - %s", d.client.SymmetrixID))
		maskingViewList, err := d.client.PmaxClient.GetMaskingViewList(ctx, d.client.SymmetrixID)

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Get PowerMax Masking View List",
				err.Error(),
			)
			return
		}
		maskingViewIds = maskingViewList.MaskingViewIDs
	} else {
		tflog.Debug(ctx, fmt.Sprintf("Get masking view Ids from filter for Symmetrix - %s", d.client.SymmetrixID))
		// get ids from filter and assign to maskingViewIds
		for _, name := range state.MaskingViewFilter.Names {
			maskingViewIds = append(maskingViewIds, name.ValueString())
		}
	}

	var models []models.MaskingViewModel
	for model := range d.getMaskingViewToConnections(ctx, resp, d.getMaskingViews(ctx, resp, maskingViewIds)) {
		models = append(models, model)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	state.MaskingViews = models
	state.ID = types.StringValue("placeholder")

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Done with Read Masking View data source ")
}

// updateMaskingViewState parse masking view and masking view connections to the state model.
func (d *maskingViewDataSource) updateMaskingViewState(maskingView *pmaxTypes.MaskingView, connections []*pmaxTypes.MaskingViewConnection) (model models.MaskingViewModel, err error) {

	model = models.MaskingViewModel{
		MaskingViewName: types.StringValue(maskingView.MaskingViewID),
		HostID:          types.StringValue(maskingView.HostID),
		HostGroupID:     types.StringValue(maskingView.HostGroupID),
		PortGroupID:     types.StringValue(maskingView.PortGroupID),
		StorageGroupID:  types.StringValue(maskingView.StorageGroupID),
	}

	var totalCapacity float64
	var volumes []attr.Value
	var ports []attr.Value
	var initiators []attr.Value
	for _, conn := range connections {
		capacity, err := strconv.ParseFloat(conn.CapacityGB, 64)
		if err != nil {
			return model, err
		}
		if !contains(volumes, conn.VolumeID) {
			volumes = append(volumes, types.StringValue(conn.VolumeID))
			totalCapacity += capacity
		}
		if !contains(initiators, conn.InitiatorID) {
			initiators = append(initiators, types.StringValue(conn.InitiatorID))
		}
		if !contains(ports, conn.DirectorPort) {
			ports = append(ports, types.StringValue(conn.DirectorPort))
		}
	}

	model.CapacityGB = types.Float64Value(totalCapacity)
	model.Volumes, _ = types.ListValue(types.StringType, volumes)
	model.Initiators, _ = types.ListValue(types.StringType, initiators)
	model.Ports, _ = types.ListValue(types.StringType, ports)

	return
}

func (d *maskingViewDataSource) getMaskingViewToConnections(ctx context.Context, resp *datasource.ReadResponse, maskingView <-chan *pmaxTypes.MaskingView) <-chan models.MaskingViewModel {

	var wg sync.WaitGroup
	ch := make(chan models.MaskingViewModel)
	go func() {
		for mv := range maskingView {
			wg.Add(1)
			go func(mv *pmaxTypes.MaskingView) {
				defer wg.Done()
				maskingViewConnection, err := d.client.PmaxClient.GetMaskingViewConnections(ctx, d.client.SymmetrixID, mv.MaskingViewID, "")
				if err != nil {
					lockMutex.Lock()
					defer lockMutex.Unlock()
					resp.Diagnostics.AddError(
						fmt.Sprintf("Failed to get MaskingViewConnections - %s.", mv.MaskingViewID),
						err.Error(),
					)
					return
				}

				model, err := d.updateMaskingViewState(mv, maskingViewConnection)
				if err != nil {
					lockMutex.Lock()
					defer lockMutex.Unlock()
					resp.Diagnostics.AddError(
						fmt.Sprintf("Failed to update masking view state - %s.", mv.MaskingViewID),
						err.Error(),
					)
					return
				}
				ch <- model
			}(mv)
		}
		wg.Wait()
		close(ch)
	}()

	return ch
}

func (d *maskingViewDataSource) getMaskingViews(ctx context.Context, resp *datasource.ReadResponse, maskingViewNames []string) <-chan *pmaxTypes.MaskingView {

	ch := make(chan *pmaxTypes.MaskingView)
	var wg sync.WaitGroup
	sem := make(chan struct{}, defaultMaxPowerMaxConnections)

	go func() {
		for _, maskingViewID := range maskingViewNames {
			sem <- struct{}{}
			wg.Add(1)
			go func(id string) {
				defer func() {
					wg.Done()
					<-sem
				}()
				tflog.Debug(ctx, fmt.Sprintf("Calling api to get MaskingView - %s", id))
				maskingView, err := d.client.PmaxClient.GetMaskingViewByID(ctx, d.client.SymmetrixID, id)
				if err != nil {
					lockMutex.Lock()
					defer lockMutex.Unlock()
					resp.Diagnostics.AddError(
						fmt.Sprintf("Failed to get MaskingView - %s.", id),
						err.Error(),
					)
					return
				}

				ch <- maskingView

			}(maskingViewID)
		}
		wg.Wait()
		close(ch)
		close(sem)
	}()
	return ch
}

// contains will return true if the slice contains the given value.
func contains(slice []attr.Value, value string) bool {
	for _, element := range slice {
		if element.Equal(types.StringValue(value)) {
			return true
		}
	}
	return false
}
