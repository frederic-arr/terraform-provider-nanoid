// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure NanoidProvider satisfies various provider interfaces.
var _ provider.Provider = &NanoidProvider{}
var _ provider.ProviderWithFunctions = &NanoidProvider{}

// NanoidProvider defines the provider implementation.
type NanoidProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// NanoidProviderModel describes the provider data model.
type NanoidProviderModel struct{}

type NanoidProviderData struct{}

func (p *NanoidProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "nanoid"
	resp.Version = p.version
}

func (p *NanoidProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Nanoid provider provides an interface to the go-nanoid library to generate unique resource identifiers.",
	}
}

func (p *NanoidProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data NanoidProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	providerData := NanoidProviderData{}
	resp.DataSourceData = &providerData
	resp.ResourceData = &providerData
}

func (p *NanoidProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewIdResource,
		NewDnsResource,
	}
}

func (p *NanoidProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *NanoidProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &NanoidProvider{
			version: version,
		}
	}
}
