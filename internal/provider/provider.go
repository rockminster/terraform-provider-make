package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure MakeProvider satisfies various provider interfaces.
var _ provider.Provider = &MakeProvider{}
var _ provider.ProviderWithFunctions = &MakeProvider{}

// MakeProvider defines the provider implementation.
type MakeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// MakeProviderModel describes the provider data model.
type MakeProviderModel struct {
	ApiToken types.String `tfsdk:"api_token"`
	BaseUrl  types.String `tfsdk:"base_url"`
}

func (p *MakeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "make"
	resp.Version = p.version
}

func (p *MakeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				MarkdownDescription: "API token for Make.com authentication. Can also be set via the MAKE_API_TOKEN environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL for Make.com API. Defaults to https://api.make.com/. Can also be set via the MAKE_BASE_URL environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *MakeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MakeProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Default configuration values
	apiToken := os.Getenv("MAKE_API_TOKEN")
	baseUrl := os.Getenv("MAKE_BASE_URL")

	if baseUrl == "" {
		baseUrl = "https://api.make.com/"
	}

	// Override with provider configuration if specified
	if !data.ApiToken.IsNull() {
		apiToken = data.ApiToken.ValueString()
	}

	if !data.BaseUrl.IsNull() {
		baseUrl = data.BaseUrl.ValueString()
	}

	// Validation
	if apiToken == "" {
		resp.Diagnostics.AddError(
			"Missing API Token Configuration",
			"While configuring the provider, the API token was not found in "+
				"the MAKE_API_TOKEN environment variable or provider "+
				"configuration block api_token attribute.",
		)
		return
	}

	// Create API client
	client := &MakeAPIClient{
		ApiToken: apiToken,
		BaseUrl:  baseUrl,
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *MakeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewScenarioResource,
	}
}

func (p *MakeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewScenarioDataSource,
	}
}

func (p *MakeProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		// Example function
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MakeProvider{
			version: version,
		}
	}
}

// MakeAPIClient represents the Make.com API client
type MakeAPIClient struct {
	ApiToken string
	BaseUrl  string
}
