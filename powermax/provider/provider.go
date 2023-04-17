// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package provider

import (
	"context"
	"terraform-provider-powermax/client"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure PmaxProvider satisfies various provider interfaces.
var _ provider.Provider = &PmaxProvider{}

// PmaxProvider defines the provider implementation.
type PmaxProvider struct {
	// client can contain the upstream provider SDK or HTTP client used to
	// communicate with the upstream service. Resource and DataSource
	// implementations can then make calls using this client.
	client *client.Client
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ProviderData describes the provider data model.
type ProviderData struct {
	Endpoint     types.String `tfsdk:"endpoint"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
	SerialNumber types.String `tfsdk:"serial_number"`
	PmaxVersion  types.String `tfsdk:"pmax_version"`
	Insecure     types.Bool   `tfsdk:"insecure"`
}

func (p *PmaxProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "powermax"
	resp.Version = p.version
}

func (p *PmaxProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "IP or FQDN of the PowerMax host",
				Description:         "IP or FQDN of the PowerMax host",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username of the PowerMax host.",
				Description:         "The username of the PowerMax host.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password of the PowerMax host.",
				Description:         "The password of the PowerMax host.",
				Required:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"serial_number": schema.StringAttribute{
				MarkdownDescription: "The serial_number of the PowerMax host.",
				Description:         "The serial_number of the PowerMax host.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "Boolean variable to specify whether to validate SSL certificate or not.",
				Description:         "Boolean variable to specify whether to validate SSL certificate or not.",
				Optional:            true,
			},
			"pmax_version": schema.StringAttribute{
				MarkdownDescription: "The version of the PowerMax host.",
				Description:         "The version of the PowerMax host.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

func (p *PmaxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ProviderData

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	pmaxClient, err := client.NewClient(
		data.Endpoint.ValueString(),
		data.Username.ValueString(),
		data.Password.ValueString(),
		data.SerialNumber.ValueString(),
		data.PmaxVersion.ValueString(),
		data.Insecure.ValueBool(),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create powermax client",
			err.Error(),
		)
		return
	}

	// client configuration for data sources and resources
	p.client = pmaxClient
	resp.DataSourceData = pmaxClient
	resp.ResourceData = pmaxClient
}

func (p *PmaxProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewStorageGroup,
		NewHostGroup,
		NewHost,
		NewPortGroup,
		newMaskingView,
	}
}

func (p *PmaxProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewHostDataSource,
		NewHostGroupDataSource,
		NewPortgroupDataSource,
		newMaskingViewDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PmaxProvider{
			version: version,
		}
	}
}
