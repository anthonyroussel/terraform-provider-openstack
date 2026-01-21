package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-provider-openstack/terraform-provider-openstack/v3/openstack/version"
)

var _ provider.Provider = &OpenStackProvider{}

type OpenStackProvider struct{}

type OpenStackProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func New() func() provider.Provider {
	return func() provider.Provider {
		return &OpenStackProvider{}
	}
}

func (p *OpenStackProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "openstack"
	resp.Version = version.Version
}

func (p *OpenStackProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for managing OpenStack ressources",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
		},
	}
}

func (p *OpenStackProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data OpenStackProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *OpenStackProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *OpenStackProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *OpenStackProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *OpenStackProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *OpenStackProvider) Actions(_ context.Context) []func() action.Action {
	return []func() action.Action{}
}
