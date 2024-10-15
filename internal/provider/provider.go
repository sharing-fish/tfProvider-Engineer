// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure DevOpsAPIProvider satisfies various provider interfaces.
var _ provider.Provider = &DevOpsAPIProvider{}

// DevOpsAPIProvider defines the provider implementation.
type DevOpsAPIProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// DevOpsAPIProviderModel describes the provider data model.
type DevOpsAPIProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *DevOpsAPIProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "devops-bootcamp"
	resp.Version = p.version
}

func (p *DevOpsAPIProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
		},
	}
}

func (p *DevOpsAPIProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring the provider")
	var data DevOpsAPIProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if data.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown DevOps API Host",
			"The provider cannot create the DevOps API client as there is an unknown configuration value for the DevOps API endpoint. "+
				"Value is stored in BOOTCAMP_API_ENDPOINT. It should be http://localhost:8080",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("BOOTCAMP_API_ENDPOINT")

	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing DevOps API Endpoint",
			"The provider cannot create the DevOps API client as there is a missing or empty value for the DevOps API endpoint. ",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "devops_api_endpoint", endpoint)
	tflog.Debug(ctx, "Creating DevOps API client")

	client, err := NewClient(&endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create DevOps API Client",
			"An unexpected error occurred when creating the DevOps API client: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured DevOps API client", map[string]any{"success": true})

}

func (p *DevOpsAPIProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewEngineerResource,
		NewDevResource,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *DevOpsAPIProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewEngineerDataSource,
		NewDevDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DevOpsAPIProvider{
			version: version,
		}
	}
}
