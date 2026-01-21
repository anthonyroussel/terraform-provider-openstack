package provider

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-provider-openstack/terraform-provider-openstack/v3/openstack/internal/consts"
	"github.com/terraform-provider-openstack/terraform-provider-openstack/v3/openstack/version"
	"github.com/terraform-provider-openstack/utils/v2/auth"
	"github.com/terraform-provider-openstack/utils/v2/mutexkv"
)

var _ provider.Provider = &OpenStackProvider{}

type OpenStackProvider struct{}

type OpenStackProviderModel struct {
	AuthURL                     types.String `tfsdk:"auth_url"`
	Region                      types.String `tfsdk:"region"`
	UserName                    types.String `tfsdk:"user_name"`
	UserID                      types.String `tfsdk:"user_id"`
	ApplicationCredentialID     types.String `tfsdk:"application_credential_id"`
	ApplicationCredentialName   types.String `tfsdk:"application_credential_name"`
	ApplicationCredentialSecret types.String `tfsdk:"application_credential_secret"`
	TenantID                    types.String `tfsdk:"tenant_id"`
	TenantName                  types.String `tfsdk:"tenant_name"`
	Password                    types.String `tfsdk:"password"`
	Token                       types.String `tfsdk:"token"`
	UserDomainName              types.String `tfsdk:"user_domain_name"`
	UserDomainID                types.String `tfsdk:"user_domain_id"`
	ProjectDomainName           types.String `tfsdk:"project_domain_name"`
	ProjectDomainID             types.String `tfsdk:"project_domain_id"`
	DomainID                    types.String `tfsdk:"domain_id"`
	DomainName                  types.String `tfsdk:"domain_name"`
	DefaultDomain               types.String `tfsdk:"default_domain"`
	SystemScope                 types.Bool   `tfsdk:"system_scope"`
	Insecure                    types.Bool   `tfsdk:"insecure"`
	EndpointType                types.String `tfsdk:"endpoint_type"`
	CACertFile                  types.String `tfsdk:"cacert_file"`
	Cert                        types.String `tfsdk:"cert"`
	Key                         types.String `tfsdk:"key"`
	Swauth                      types.Bool   `tfsdk:"swauth"`
	DelayedAuth                 types.Bool   `tfsdk:"delayed_auth"`
	AllowReauth                 types.Bool   `tfsdk:"allow_reauth"`
	Cloud                       types.String `tfsdk:"cloud"`
	MaxRetries                  types.Int64  `tfsdk:"max_retries"`
	EndpointOverrides           types.Map    `tfsdk:"endpoint_overrides"`
	DisableNoCacheHeader        types.Bool   `tfsdk:"disable_no_cache_header"`
	EnableLogging               types.Bool   `tfsdk:"enable_logging"`
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
			"auth_url": schema.StringAttribute{
				MarkdownDescription: consts.AuthURL,
				Optional:            true,
			},

			"region": schema.StringAttribute{
				MarkdownDescription: consts.Region,
				Optional:            true,
			},

			"user_name": schema.StringAttribute{
				MarkdownDescription: consts.UserName,
				Optional:            true,
			},

			"user_id": schema.StringAttribute{
				MarkdownDescription: consts.UserID,
				Optional:            true,
			},

			"application_credential_id": schema.StringAttribute{
				MarkdownDescription: consts.ApplicationCredentialID,
				Optional:            true,
			},

			"application_credential_name": schema.StringAttribute{
				MarkdownDescription: consts.ApplicationCredentialName,
				Optional:            true,
			},

			"application_credential_secret": schema.StringAttribute{
				MarkdownDescription: consts.ApplicationCredentialSecret,
				Optional:            true,
				Sensitive:           true,
			},

			"tenant_id": schema.StringAttribute{
				MarkdownDescription: consts.TenantID,
				Optional:            true,
			},

			"tenant_name": schema.StringAttribute{
				MarkdownDescription: consts.TenantName,
				Optional:            true,
			},

			"password": schema.StringAttribute{
				MarkdownDescription: consts.Password,
				Optional:            true,
				Sensitive:           true,
			},

			"token": schema.StringAttribute{
				MarkdownDescription: consts.Token,
				Optional:            true,
			},

			"user_domain_name": schema.StringAttribute{
				MarkdownDescription: consts.UserDomainName,
				Optional:            true,
			},

			"user_domain_id": schema.StringAttribute{
				MarkdownDescription: consts.UserDomainID,
				Optional:            true,
			},

			"project_domain_name": schema.StringAttribute{
				MarkdownDescription: consts.ProjectDomainName,
				Optional:            true,
			},

			"project_domain_id": schema.StringAttribute{
				MarkdownDescription: consts.ProjectDomainID,
				Optional:            true,
			},

			"domain_id": schema.StringAttribute{
				MarkdownDescription: consts.DomainID,
				Optional:            true,
			},

			"domain_name": schema.StringAttribute{
				MarkdownDescription: consts.DomainName,
				Optional:            true,
			},

			"default_domain": schema.StringAttribute{
				MarkdownDescription: consts.DefaultDomain,
				Optional:            true,
			},

			"system_scope": schema.BoolAttribute{
				MarkdownDescription: consts.SystemScope,
				Optional:            true,
			},

			"insecure": schema.BoolAttribute{
				MarkdownDescription: consts.Insecure,
				Optional:            true,
			},

			"endpoint_type": schema.StringAttribute{
				MarkdownDescription: consts.EndpointType,
				Optional:            true,
			},

			"cacert_file": schema.StringAttribute{
				MarkdownDescription: consts.CACertFile,
				Optional:            true,
			},

			"cert": schema.StringAttribute{
				MarkdownDescription: consts.Cert,
				Optional:            true,
			},

			"key": schema.StringAttribute{
				MarkdownDescription: consts.Key,
				Optional:            true,
			},

			"swauth": schema.BoolAttribute{
				MarkdownDescription: consts.Swauth,
				Optional:            true,
			},

			"delayed_auth": schema.BoolAttribute{
				MarkdownDescription: consts.DelayedAuth,
				Optional:            true,
			},

			"allow_reauth": schema.BoolAttribute{
				MarkdownDescription: consts.AllowReauth,
				Optional:            true,
			},

			"cloud": schema.StringAttribute{
				MarkdownDescription: consts.Cloud,
				Optional:            true,
			},

			"max_retries": schema.Int64Attribute{
				MarkdownDescription: consts.MaxRetries,
				Optional:            true,
			},

			"endpoint_overrides": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: consts.EndpointOverrides,
				Optional:            true,
			},

			"disable_no_cache_header": schema.BoolAttribute{
				MarkdownDescription: consts.DisableNoCacheHeader,
				Optional:            true,
			},

			"enable_logging": schema.BoolAttribute{
				MarkdownDescription: consts.EnableLogging,
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

	authOpts := &gophercloud.AuthOptions{
		Scope: &gophercloud.AuthScope{System: getBoolValue(data.SystemScope, []string{"OS_SYSTEM_SCOPE"}, false)},
	}

	config := auth.Config{
		CACertFile:                  getStringValue(data.CACertFile, []string{"OS_CACERT"}, ""),
		ClientCertFile:              getStringValue(data.Cert, []string{"OS_CERT"}, ""),
		ClientKeyFile:               getStringValue(data.Key, []string{"OS_KEY"}, ""),
		Cloud:                       getStringValue(data.Cloud, []string{"OS_CLOUD"}, ""),
		DefaultDomain:               getStringValue(data.DefaultDomain, []string{"OS_DEFAULT_DOMAIN"}, "default"),
		DomainID:                    getStringValue(data.DomainID, []string{"OS_DOMAIN_ID"}, ""),
		DomainName:                  getStringValue(data.DomainName, []string{"OS_DOMAIN_NAME"}, ""),
		EndpointType:                getStringValue(data.EndpointType, []string{"OS_ENDPOINT_TYPE"}, ""),
		IdentityEndpoint:            getStringValue(data.AuthURL, []string{"OS_AUTH_URL"}, ""),
		Password:                    getStringValue(data.Password, []string{"OS_PASSWORD"}, ""),
		ProjectDomainID:             getStringValue(data.ProjectDomainID, []string{"OS_PROJECT_DOMAIN_ID"}, ""),
		ProjectDomainName:           getStringValue(data.ProjectDomainName, []string{"OS_PROJECT_DOMAIN_NAME"}, ""),
		Region:                      getStringValue(data.Region, []string{"OS_REGION_NAME"}, ""),
		Swauth:                      getBoolValue(data.Swauth, []string{"OS_SWAUTH"}, false),
		Token:                       getStringValue(data.Token, []string{"OS_TOKEN"}, ""),
		TenantID:                    getStringValue(data.TenantID, []string{"OS_PROJECT_ID"}, ""),
		TenantName:                  getStringValue(data.TenantName, []string{"OS_PROJECT_NAME"}, ""),
		UserDomainID:                getStringValue(data.UserDomainID, []string{"OS_USER_DOMAIN_ID"}, ""),
		UserDomainName:              getStringValue(data.UserDomainName, []string{"OS_USER_DOMAIN_NAME"}, ""),
		Username:                    getStringValue(data.UserName, []string{"OS_USERNAME"}, ""),
		UserID:                      getStringValue(data.UserID, []string{"OS_USER_ID"}, ""),
		UseOctavia:                  true,
		ApplicationCredentialID:     getStringValue(data.ApplicationCredentialID, []string{"OS_APPLICATION_CREDENTIAL_ID"}, ""),
		ApplicationCredentialName:   getStringValue(data.ApplicationCredentialName, []string{"OS_APPLICATION_CREDENTIAL_NAME"}, ""),
		ApplicationCredentialSecret: getStringValue(data.ApplicationCredentialSecret, []string{"OS_APPLICATION_CREDENTIAL_SECRET"}, ""),
		DelayedAuth:                 getBoolValue(data.DelayedAuth, []string{"OS_DELAYED_AUTH"}, true),
		AllowReauth:                 getBoolValue(data.AllowReauth, []string{"OS_ALLOW_REAUTH"}, false),
		AuthOpts:                    authOpts,
		MaxRetries:                  getIntValue(data.MaxRetries, []string{}, 0),
		DisableNoCacheHeader:        getBoolValue(data.DisableNoCacheHeader, []string{}, false),
		TerraformVersion:            req.TerraformVersion,
		SDKVersion:                  getSDKVersion() + " Terraform Provider OpenStack/" + version.Version,
		MutexKV:                     mutexkv.NewMutexKV(),
		EnableLogger:                getBoolValue(data.EnableLogging, []string{}, false),
	}

	// FIXME: handle endpoint_overrides + insecure + enable_logging

	if err := config.LoadAndValidate(ctx); err != nil {
		resp.Diagnostics.AddError("Unable to authenticate with OpenStack", err.Error())

		return
	}

	resp.DataSourceData = config
	resp.ResourceData = config
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
