// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"terraform-provider-tasks/internal/sdk"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure RightbrainProvider satisfies various provider interfaces.
var _ provider.Provider = &RightbrainProvider{}
var _ provider.ProviderWithFunctions = &RightbrainProvider{}
var _ provider.ProviderWithEphemeralResources = &RightbrainProvider{}

const (
	ProviderName     = "rightbrain"
	DefaultOAuthHost = "https://oauth.rightbrain.ai"
	DefaultAPIHost   = "https://app.rightbrain.ai"
)

// RightbrainProvider defines the provider implementation.
type RightbrainProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// RightbrainProviderModel describes the provider data model.
type RightbrainProviderModel struct {
	RightbrainAPIHost      types.String `tfsdk:"api_host"`
	RightbrainOAuthHost    types.String `tfsdk:"oauth_host"`
	RightbrainClientID     types.String `tfsdk:"client_id"`
	RightbrainClientSecret types.String `tfsdk:"client_secret"`
	RightbrainOrgID        types.String `tfsdk:"org_id"`
	RightbrainProjectID    types.String `tfsdk:"project_id"`
}

func (p *RightbrainProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = ProviderName
	resp.Version = p.version
}

func (p *RightbrainProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_host": schema.StringAttribute{
				MarkdownDescription: "The hostname for the Rightbrain API server",
				Optional:            true,
			},
			"oauth_host": schema.StringAttribute{
				MarkdownDescription: "The hostname for the Rightbrain OAuth server",
				Optional:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The OAuth Client ID",
				Required:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "The OAuth Client Secret",
				Required:            true,
			},
			"org_id": schema.StringAttribute{
				MarkdownDescription: "The Org ID",
				Required:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The Project ID",
				Required:            true,
			},
		},
	}
}

func (p *RightbrainProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data RightbrainProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client, err := p.newRightbrainClient(data)
	if err != nil {
		resp.Diagnostics.AddError("cannot create rightbrain client", err.Error())
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *RightbrainProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTaskResource,
	}
}

func (p *RightbrainProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *RightbrainProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewLLMModelDataSource,
	}
}

func (p *RightbrainProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RightbrainProvider{
			version: version,
		}
	}
}

func (p *RightbrainProvider) newRightbrainClient(data RightbrainProviderModel) (*sdk.TasksClient, error) {
	oauthURL := fmt.Sprintf("%s/oauth2/token", p.isEmptyValueElseDefault(data.RightbrainOAuthHost, DefaultOAuthHost))
	tokenStore, err := sdk.NewDefaultTokenStore(oauthURL)
	if err != nil {
		return nil, err
	}
	return sdk.NewTasksClient(TerraformLog{}, http.DefaultClient, tokenStore, sdk.Config{
		RightbrainAPIHost:      p.isEmptyValueElseDefault(data.RightbrainAPIHost, DefaultAPIHost),
		RightbrainClientID:     data.RightbrainClientID.ValueString(),
		RightbrainClientSecret: data.RightbrainClientSecret.ValueString(),
		RightbrainOrgID:        data.RightbrainOrgID.ValueString(),
		RightbrainProjectID:    data.RightbrainProjectID.ValueString(),
	}), nil
}

func (p *RightbrainProvider) isEmptyValueElseDefault(field types.String, defaultValue string) string {
	if field.IsNull() {
		return defaultValue
	}
	return field.ValueString()
}
