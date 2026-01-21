package consts

const (
	AuthURL = "The Identity authentication URL."

	Cloud = "An entry in a `clouds.yaml` file to use."

	Region = "The OpenStack region to connect to."

	UserName = "Username to login with."

	UserID = "User ID to login with."

	ApplicationCredentialID = "Application Credential ID to login with."

	ApplicationCredentialName = "Application Credential name to login with."

	ApplicationCredentialSecret = "Application Credential secret to login with."

	TenantID = "The ID of the Tenant (Identity v2) or Project (Identity v3)\n" +
		"to login with."

	TenantName = "The name of the Tenant (Identity v2) or Project (Identity v3)\n" +
		"to login with."

	Password = "Password to login with."

	Token = "Authentication token to use as an alternative to username/password."

	UserDomainName = "The name of the domain where the user resides (Identity v3)."

	UserDomainID = "The ID of the domain where the user resides (Identity v3)."

	ProjectDomainName = "The name of the domain where the project resides (Identity v3)."

	ProjectDomainID = "The ID of the domain where the proejct resides (Identity v3)."

	DomainID = "The ID of the Domain to scope to (Identity v3)."

	DomainName = "The name of the Domain to scope to (Identity v3)."

	DefaultDomain = "The name of the Domain ID to scope to if no other domain is specified. Defaults to `default` (Identity v3)."

	SystemScope = "If set to `true`, system scoped authorization will be enabled. Defaults to `false` (Identity v3)."

	Insecure = "Trust self-signed certificates."

	CACertFile = "A Custom CA certificate."

	Cert = "A client certificate to authenticate with."

	Key = "A client private key to authenticate with."

	EndpointType = "The catalog endpoint type to use."

	EndpointOverrides = "A map of services with an endpoint to override what was\n" +
		"from the Keystone catalog"

	Swauth = "Use Swift's authentication system instead of Keystone. Only used for\n" +
		"interaction with Swift."

	DisableNoCacheHeader = "If set to `true`, the HTTP `Cache-Control: no-cache` header will not be added by default to all API requests."

	DelayedAuth = "If set to `false`, OpenStack authorization will be perfomed,\n" +
		"every time the service provider client is called. Defaults to `true`."

	AllowReauth = "If set to `false`, OpenStack authorization won't be perfomed\n" +
		"automatically, if the initial auth token get expired. Defaults to `true`"

	MaxRetries = "How many times HTTP connection should be retried until giving up."

	EnableLogging = "Outputs very verbose logs with all calls made to and responses from OpenStack"
)
