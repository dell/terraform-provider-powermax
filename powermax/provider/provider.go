/*
Copyright (c) 2025 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"os"
	"strconv"
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

// Data describes the provider data model.
type Data struct {
	Endpoint     types.String `tfsdk:"endpoint"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
	SerialNumber types.String `tfsdk:"serial_number"`
	PmaxVersion  types.String `tfsdk:"pmax_version"`
	Insecure     types.Bool   `tfsdk:"insecure"`
}

// Metadata returns the provider metadata.
func (p *PmaxProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "powermax"
	resp.Version = p.version
}

// Schema returns the provider schema.
func (p *PmaxProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Terraform provider for Dell PowerMax " +
			"can be used to interact with a Dell PowerMax array in order to manage the array resources.",
		MarkdownDescription: "The Terraform provider for Dell PowerMax " +
			"can be used to interact with a Dell PowerMax array in order to manage the array resources.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "IP or FQDN of the PowerMax host. This can also be set using the environment variable POWERMAX_ENDPOINT",
				Description:         "IP or FQDN of the PowerMax host. This can also be set using the environment variable POWERMAX_ENDPOINT",
				// This should remain optional so user can use environment variables if they choose.
				Optional: true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username of the PowerMax host. This can also be set using the environment variable POWERMAX_USERNAME",
				Description:         "The username of the PowerMax host. This can also be set using the environment variable POWERMAX_USERNAME",
				// This should remain optional so user can use environment variables if they choose.
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password of the PowerMax host. This can also be set using the environment variable POWERMAX_PASSWORD",
				Description:         "The password of the PowerMax host. This can also be set using the environment variable POWERMAX_PASSWORD",
				// This should remain optional so user can use environment variables if they choose.
				Optional:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"serial_number": schema.StringAttribute{
				MarkdownDescription: "The serial_number of the PowerMax host. This can also be set using the environment variable POWERMAX_SERIAL_NUMBER",
				Description:         "The serial_number of the PowerMax host. This can also be set using the environment variable POWERMAX_SERIAL_NUMBER",
				// This should remain optional so user can use environment variables if they choose.
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "Boolean variable to specify whether to validate SSL certificate or not. This can also be set using the environment variable POWERMAX_INSECURE",
				Description:         "Boolean variable to specify whether to validate SSL certificate or not. This can also be set using the environment variable POWERMAX_INSECURE",
				// This should remain optional so user can use environment variables if they choose.
				Optional: true,
			},
			"pmax_version": schema.StringAttribute{
				MarkdownDescription: "The version of the PowerMax host. This can also be set using the environment variable POWERMAX_POWERMAX_VERSION",
				Description:         "The version of the PowerMax host. This can also be set using the environment variable POWERMAX_POWERMAX_VERSION",
				// This should remain optional so user can use environment variables if they choose.
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

// Configure configures the provider.
func (p *PmaxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data Data

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	usernameEnv := os.Getenv("POWERMAX_USERNAME")
	if usernameEnv != "" {
		data.Username = types.StringValue(usernameEnv)
	}

	passEnv := os.Getenv("POWERMAX_PASSWORD")
	if passEnv != "" {
		data.Password = types.StringValue(passEnv)
	}

	endpointEnv := os.Getenv("POWERMAX_ENDPOINT")
	if endpointEnv != "" {
		data.Endpoint = types.StringValue(endpointEnv)
	}

	serialNumberEnv := os.Getenv("POWERMAX_SERIAL_NUMBER")
	if serialNumberEnv != "" {
		data.SerialNumber = types.StringValue(serialNumberEnv)
	}

	versionEnv := os.Getenv("POWERMAX_VERSION")
	if versionEnv != "" {
		data.PmaxVersion = types.StringValue(versionEnv)
	}

	insecureEnv, errInsecure := strconv.ParseBool(os.Getenv("POWERMAX_INSECURE"))
	if errInsecure == nil {
		data.Insecure = types.BoolValue(insecureEnv)
	}

	// Configuration values are now available.
	pmaxClient, err := client.NewClient(
		ctx,
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

	// Do a dummy call to validate the client configuration
	_, _, errClient := pmaxClient.PmaxOpenapiClient.SLOProvisioningApi.ListHosts(ctx, pmaxClient.SymmetrixID).Execute()
	if errClient != nil {
		resp.Diagnostics.AddError(
			"Unable to create powermax client",
			"Please validate that the endpoint, username, serial_number and password are correct.",
		)
		return
	}

	// client configuration for data sources and resources
	p.client = pmaxClient
	resp.DataSourceData = pmaxClient
	resp.ResourceData = pmaxClient
}

// Resources returns the provider resources.
func (p *PmaxProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewStorageGroup,
		NewHostGroup,
		NewHost,
		NewPortGroup,
		NewMaskingView,
		NewVolumeResource,
		NewSnapshotResource,
		NewSnapshotPolicy,
	}
}

// DataSources returns the provider data sources.
func (p *PmaxProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewHostDataSource,
		NewHostGroupDataSource,
		NewPortgroupDataSource,
		NewVolumeDataSource,
		NewMaskingViewDataSource,
		NewStorageGroupDataSource,
		NewSnapshotDataSource,
		NewPortDataSource,
		NewSnapshotPolicyDataSource,
	}
}

// New returns a new provider.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PmaxProvider{
			version: version,
		}
	}
}
