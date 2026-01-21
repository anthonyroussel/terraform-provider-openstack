package openstack

import (
	"context"
	"os"
	"runtime/debug"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-provider-openstack/terraform-provider-openstack/v3/openstack/internal/consts"
	"github.com/terraform-provider-openstack/terraform-provider-openstack/v3/openstack/version"
	"github.com/terraform-provider-openstack/utils/v2/auth"
	"github.com/terraform-provider-openstack/utils/v2/mutexkv"
)

// Use openstackbase.Config as the base/foundation of this provider's
// Config struct.
type Config struct {
	auth.Config
}

// Provider returns a schema.Provider for OpenStack.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_URL", ""),
				Description: consts.AuthURL,
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: consts.Region,
				DefaultFunc: schema.EnvDefaultFunc("OS_REGION_NAME", ""),
			},

			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USERNAME", ""),
				Description: consts.UserName,
			},

			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_ID", ""),
				Description: consts.UserID,
			},

			"application_credential_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_APPLICATION_CREDENTIAL_ID", ""),
				Description: consts.ApplicationCredentialID,
			},

			"application_credential_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_APPLICATION_CREDENTIAL_NAME", ""),
				Description: consts.ApplicationCredentialName,
			},

			"application_credential_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_APPLICATION_CREDENTIAL_SECRET", ""),
				Description: consts.ApplicationCredentialSecret,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_ID",
					"OS_PROJECT_ID",
				}, ""),
				Description: consts.TenantID,
			},

			"tenant_name": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_NAME",
					"OS_PROJECT_NAME",
				}, ""),
				Description: consts.TenantName,
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PASSWORD", ""),
				Description: consts.Password,
			},

			"token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TOKEN",
					"OS_AUTH_TOKEN",
				}, ""),
				Description: consts.Token,
			},

			"user_domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_DOMAIN_NAME", ""),
				Description: consts.UserDomainName,
			},

			"user_domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_DOMAIN_ID", ""),
				Description: consts.UserDomainID,
			},

			"project_domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PROJECT_DOMAIN_NAME", ""),
				Description: consts.ProjectDomainName,
			},

			"project_domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PROJECT_DOMAIN_ID", ""),
				Description: consts.ProjectDomainID,
			},

			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DOMAIN_ID", ""),
				Description: consts.DomainID,
			},

			"domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DOMAIN_NAME", ""),
				Description: consts.DomainName,
			},

			"default_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DEFAULT_DOMAIN", "default"),
				Description: consts.DefaultDomain,
			},

			"system_scope": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_SYSTEM_SCOPE", false),
				Description: consts.SystemScope,
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_INSECURE", nil),
				Description: consts.Insecure,
			},

			"endpoint_type": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_ENDPOINT_TYPE", ""),
				Description: consts.EndpointType,
			},

			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CACERT", ""),
				Description: consts.CACertFile,
			},

			"cert": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CERT", ""),
				Description: consts.Cert,
			},

			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_KEY", ""),
				Description: consts.Key,
			},

			"swauth": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_SWAUTH", false),
				Description: consts.Swauth,
			},

			"delayed_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DELAYED_AUTH", true),
				Description: consts.DelayedAuth,
			},

			"allow_reauth": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_ALLOW_REAUTH", true),
				Description: consts.AllowReauth,
			},

			"cloud": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CLOUD", ""),
				Description: consts.Cloud,
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: consts.MaxRetries,
			},

			"endpoint_overrides": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: consts.EndpointOverrides,
			},

			"disable_no_cache_header": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: consts.DisableNoCacheHeader,
			},

			"enable_logging": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: consts.EnableLogging,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"openstack_blockstorage_availability_zones_v3":       dataSourceBlockStorageAvailabilityZonesV3(),
			"openstack_blockstorage_snapshot_v3":                 dataSourceBlockStorageSnapshotV3(),
			"openstack_blockstorage_volume_v3":                   dataSourceBlockStorageVolumeV3(),
			"openstack_blockstorage_quotaset_v3":                 dataSourceBlockStorageQuotasetV3(),
			"openstack_compute_aggregate_v2":                     dataSourceComputeAggregateV2(),
			"openstack_compute_availability_zones_v2":            dataSourceComputeAvailabilityZonesV2(),
			"openstack_compute_instance_v2":                      dataSourceComputeInstanceV2(),
			"openstack_compute_flavor_v2":                        dataSourceComputeFlavorV2(),
			"openstack_compute_hypervisor_v2":                    dataSourceComputeHypervisorV2(),
			"openstack_compute_servergroup_v2":                   dataSourceComputeServerGroupV2(),
			"openstack_compute_keypair_v2":                       dataSourceComputeKeypairV2(),
			"openstack_compute_quotaset_v2":                      dataSourceComputeQuotasetV2(),
			"openstack_compute_limits_v2":                        dataSourceComputeLimitsV2(),
			"openstack_containerinfra_nodegroup_v1":              dataSourceContainerInfraNodeGroupV1(),
			"openstack_containerinfra_clustertemplate_v1":        dataSourceContainerInfraClusterTemplateV1(),
			"openstack_containerinfra_cluster_v1":                dataSourceContainerInfraCluster(),
			"openstack_dns_zone_v2":                              dataSourceDNSZoneV2(),
			"openstack_dns_zone_share_v2":                        dataSourceDNSZoneShareV2(),
			"openstack_fw_group_v2":                              dataSourceFWGroupV2(),
			"openstack_fw_policy_v2":                             dataSourceFWPolicyV2(),
			"openstack_fw_rule_v2":                               dataSourceFWRuleV2(),
			"openstack_identity_role_v3":                         dataSourceIdentityRoleV3(),
			"openstack_identity_project_v3":                      dataSourceIdentityProjectV3(),
			"openstack_identity_project_ids_v3":                  dataSourceIdentityProjectIDsV3(),
			"openstack_identity_user_v3":                         dataSourceIdentityUserV3(),
			"openstack_identity_auth_scope_v3":                   dataSourceIdentityAuthScopeV3(),
			"openstack_identity_endpoint_v3":                     dataSourceIdentityEndpointV3(),
			"openstack_identity_service_v3":                      dataSourceIdentityServiceV3(),
			"openstack_identity_group_v3":                        dataSourceIdentityGroupV3(),
			"openstack_images_image_v2":                          dataSourceImagesImageV2(),
			"openstack_images_image_ids_v2":                      dataSourceImagesImageIDsV2(),
			"openstack_networking_addressscope_v2":               dataSourceNetworkingAddressScopeV2(),
			"openstack_networking_network_v2":                    dataSourceNetworkingNetworkV2(),
			"openstack_networking_qos_bandwidth_limit_rule_v2":   dataSourceNetworkingQoSBandwidthLimitRuleV2(),
			"openstack_networking_qos_dscp_marking_rule_v2":      dataSourceNetworkingQoSDSCPMarkingRuleV2(),
			"openstack_networking_qos_minimum_bandwidth_rule_v2": dataSourceNetworkingQoSMinimumBandwidthRuleV2(),
			"openstack_networking_qos_policy_v2":                 dataSourceNetworkingQoSPolicyV2(),
			"openstack_networking_quota_v2":                      dataSourceNetworkingQuotaV2(),
			"openstack_networking_subnet_v2":                     dataSourceNetworkingSubnetV2(),
			"openstack_networking_subnet_ids_v2":                 dataSourceNetworkingSubnetIDsV2(),
			"openstack_networking_secgroup_v2":                   dataSourceNetworkingSecGroupV2(),
			"openstack_networking_subnetpool_v2":                 dataSourceNetworkingSubnetPoolV2(),
			"openstack_networking_floatingip_v2":                 dataSourceNetworkingFloatingIPV2(),
			"openstack_networking_router_v2":                     dataSourceNetworkingRouterV2(),
			"openstack_networking_port_v2":                       dataSourceNetworkingPortV2(),
			"openstack_networking_port_ids_v2":                   dataSourceNetworkingPortIDsV2(),
			"openstack_networking_trunk_v2":                      dataSourceNetworkingTrunkV2(),
			"openstack_networking_segment_v2":                    dataSourceNetworkingSegmentV2(),
			"openstack_sharedfilesystem_availability_zones_v2":   dataSourceSharedFilesystemAvailabilityZonesV2(),
			"openstack_sharedfilesystem_sharenetwork_v2":         dataSourceSharedFilesystemShareNetworkV2(),
			"openstack_sharedfilesystem_share_v2":                dataSourceSharedFilesystemShareV2(),
			"openstack_sharedfilesystem_snapshot_v2":             dataSourceSharedFilesystemSnapshotV2(),
			"openstack_keymanager_secret_v1":                     dataSourceKeyManagerSecretV1(),
			"openstack_keymanager_container_v1":                  dataSourceKeyManagerContainerV1(),
			"openstack_loadbalancer_flavor_v2":                   dataSourceLoadBalancerFlavorV2(),
			"openstack_lb_flavor_v2":                             dataSourceLBFlavorV2(),
			"openstack_lb_flavorprofile_v2":                      dataSourceLBFlavorProfileV2(),
			"openstack_lb_loadbalancer_v2":                       dataSourceLBLoadbalancerV2(),
			"openstack_lb_listener_v2":                           dataSourceLBListenerV2(),
			"openstack_lb_member_v2":                             dataSourceLBMemberV2(),
			"openstack_lb_monitor_v2":                            dataSourceLBMonitorV2(),
			"openstack_lb_pool_v2":                               dataSourceLBPoolV2(),
			"openstack_workflow_cron_trigger_v2":                 dataSourceWorkflowCronTriggerV2(),
			"openstack_workflow_workflow_v2":                     dataSourceWorkflowWorkflowV2(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"openstack_blockstorage_qos_association_v3":          resourceBlockStorageQosAssociationV3(),
			"openstack_blockstorage_qos_v3":                      resourceBlockStorageQosV3(),
			"openstack_blockstorage_quotaset_v3":                 resourceBlockStorageQuotasetV3(),
			"openstack_blockstorage_volume_v3":                   resourceBlockStorageVolumeV3(),
			"openstack_blockstorage_volume_attach_v3":            resourceBlockStorageVolumeAttachV3(),
			"openstack_blockstorage_volume_type_access_v3":       resourceBlockstorageVolumeTypeAccessV3(),
			"openstack_blockstorage_volume_type_v3":              resourceBlockStorageVolumeTypeV3(),
			"openstack_compute_aggregate_v2":                     resourceComputeAggregateV2(),
			"openstack_compute_flavor_v2":                        resourceComputeFlavorV2(),
			"openstack_compute_flavor_access_v2":                 resourceComputeFlavorAccessV2(),
			"openstack_compute_instance_v2":                      resourceComputeInstanceV2(),
			"openstack_compute_interface_attach_v2":              resourceComputeInterfaceAttachV2(),
			"openstack_compute_keypair_v2":                       resourceComputeKeypairV2(),
			"openstack_compute_servergroup_v2":                   resourceComputeServerGroupV2(),
			"openstack_compute_quotaset_v2":                      resourceComputeQuotasetV2(),
			"openstack_compute_volume_attach_v2":                 resourceComputeVolumeAttachV2(),
			"openstack_containerinfra_nodegroup_v1":              resourceContainerInfraNodeGroupV1(),
			"openstack_containerinfra_clustertemplate_v1":        resourceContainerInfraClusterTemplateV1(),
			"openstack_containerinfra_cluster_v1":                resourceContainerInfraClusterV1(),
			"openstack_db_instance_v1":                           resourceDatabaseInstanceV1(),
			"openstack_db_user_v1":                               resourceDatabaseUserV1(),
			"openstack_db_configuration_v1":                      resourceDatabaseConfigurationV1(),
			"openstack_db_database_v1":                           resourceDatabaseDatabaseV1(),
			"openstack_dns_recordset_v2":                         resourceDNSRecordSetV2(),
			"openstack_dns_zone_v2":                              resourceDNSZoneV2(),
			"openstack_dns_zone_share_v2":                        resourceDNSZoneShareV2(),
			"openstack_dns_transfer_request_v2":                  resourceDNSTransferRequestV2(),
			"openstack_dns_transfer_accept_v2":                   resourceDNSTransferAcceptV2(),
			"openstack_dns_quota_v2":                             resourceDNSQuotaV2(),
			"openstack_fw_group_v2":                              resourceFWGroupV2(),
			"openstack_fw_policy_v2":                             resourceFWPolicyV2(),
			"openstack_fw_rule_v2":                               resourceFWRuleV2(),
			"openstack_identity_endpoint_v3":                     resourceIdentityEndpointV3(),
			"openstack_identity_project_v3":                      resourceIdentityProjectV3(),
			"openstack_identity_role_v3":                         resourceIdentityRoleV3(),
			"openstack_identity_role_assignment_v3":              resourceIdentityRoleAssignmentV3(),
			"openstack_identity_inherit_role_assignment_v3":      resourceIdentityInheritRoleAssignmentV3(),
			"openstack_identity_service_v3":                      resourceIdentityServiceV3(),
			"openstack_identity_user_v3":                         resourceIdentityUserV3(),
			"openstack_identity_user_membership_v3":              resourceIdentityUserMembershipV3(),
			"openstack_identity_group_v3":                        resourceIdentityGroupV3(),
			"openstack_identity_application_credential_v3":       resourceIdentityApplicationCredentialV3(),
			"openstack_identity_ec2_credential_v3":               resourceIdentityEc2CredentialV3(),
			"openstack_identity_registered_limit_v3":             resourceIdentityRegisteredLimitV3(),
			"openstack_identity_limit_v3":                        resourceIdentityLimitV3(),
			"openstack_images_image_v2":                          resourceImagesImageV2(),
			"openstack_images_image_access_v2":                   resourceImagesImageAccessV2(),
			"openstack_images_image_access_accept_v2":            resourceImagesImageAccessAcceptV2(),
			"openstack_lb_flavor_v2":                             resourceLoadBalancerFlavorV2(),
			"openstack_lb_flavorprofile_v2":                      resourceLoadBalancerFlavorProfileV2(),
			"openstack_lb_loadbalancer_v2":                       resourceLoadBalancerV2(),
			"openstack_lb_listener_v2":                           resourceListenerV2(),
			"openstack_lb_pool_v2":                               resourcePoolV2(),
			"openstack_lb_member_v2":                             resourceMemberV2(),
			"openstack_lb_members_v2":                            resourceMembersV2(),
			"openstack_lb_monitor_v2":                            resourceMonitorV2(),
			"openstack_lb_l7policy_v2":                           resourceL7PolicyV2(),
			"openstack_lb_l7rule_v2":                             resourceL7RuleV2(),
			"openstack_lb_quota_v2":                              resourceLoadBalancerQuotaV2(),
			"openstack_networking_bgp_speaker_v2":                resourceNetworkingBGPSpeakerV2(),
			"openstack_networking_bgp_peer_v2":                   resourceNetworkingBGPPeerV2(),
			"openstack_networking_floatingip_v2":                 resourceNetworkingFloatingIPV2(),
			"openstack_networking_floatingip_associate_v2":       resourceNetworkingFloatingIPAssociateV2(),
			"openstack_networking_network_v2":                    resourceNetworkingNetworkV2(),
			"openstack_networking_port_v2":                       resourceNetworkingPortV2(),
			"openstack_networking_rbac_policy_v2":                resourceNetworkingRBACPolicyV2(),
			"openstack_networking_port_secgroup_associate_v2":    resourceNetworkingPortSecGroupAssociateV2(),
			"openstack_networking_qos_bandwidth_limit_rule_v2":   resourceNetworkingQoSBandwidthLimitRuleV2(),
			"openstack_networking_qos_dscp_marking_rule_v2":      resourceNetworkingQoSDSCPMarkingRuleV2(),
			"openstack_networking_qos_minimum_bandwidth_rule_v2": resourceNetworkingQoSMinimumBandwidthRuleV2(),
			"openstack_networking_qos_policy_v2":                 resourceNetworkingQoSPolicyV2(),
			"openstack_networking_quota_v2":                      resourceNetworkingQuotaV2(),
			"openstack_networking_router_v2":                     resourceNetworkingRouterV2(),
			"openstack_networking_router_interface_v2":           resourceNetworkingRouterInterfaceV2(),
			"openstack_networking_router_route_v2":               resourceNetworkingRouterRouteV2(),
			"openstack_networking_router_routes_v2":              resourceNetworkingRouterRoutesV2(),
			"openstack_networking_secgroup_v2":                   resourceNetworkingSecGroupV2(),
			"openstack_networking_secgroup_rule_v2":              resourceNetworkingSecGroupRuleV2(),
			"openstack_networking_address_group_v2":              resourceNetworkingAddressGroupV2(),
			"openstack_networking_subnet_v2":                     resourceNetworkingSubnetV2(),
			"openstack_networking_subnet_route_v2":               resourceNetworkingSubnetRouteV2(),
			"openstack_networking_subnetpool_v2":                 resourceNetworkingSubnetPoolV2(),
			"openstack_networking_addressscope_v2":               resourceNetworkingAddressScopeV2(),
			"openstack_networking_trunk_v2":                      resourceNetworkingTrunkV2(),
			"openstack_networking_portforwarding_v2":             resourceNetworkingPortForwardingV2(),
			"openstack_networking_segment_v2":                    resourceNetworkingSegmentV2(),
			"openstack_objectstorage_account_v1":                 resourceObjectStorageAccountV1(),
			"openstack_objectstorage_container_v1":               resourceObjectStorageContainerV1(),
			"openstack_objectstorage_object_v1":                  resourceObjectStorageObjectV1(),
			"openstack_objectstorage_tempurl_v1":                 resourceObjectstorageTempurlV1(),
			"openstack_orchestration_stack_v1":                   resourceOrchestrationStackV1(),
			"openstack_taas_tap_mirror_v2":                       resourceTapMirrorV2(),
			"openstack_vpnaas_ipsec_policy_v2":                   resourceIPSecPolicyV2(),
			"openstack_vpnaas_service_v2":                        resourceServiceV2(),
			"openstack_vpnaas_ike_policy_v2":                     resourceIKEPolicyV2(),
			"openstack_vpnaas_endpoint_group_v2":                 resourceEndpointGroupV2(),
			"openstack_vpnaas_site_connection_v2":                resourceSiteConnectionV2(),
			"openstack_sharedfilesystem_securityservice_v2":      resourceSharedFilesystemSecurityServiceV2(),
			"openstack_sharedfilesystem_sharenetwork_v2":         resourceSharedFilesystemShareNetworkV2(),
			"openstack_sharedfilesystem_share_v2":                resourceSharedFilesystemShareV2(),
			"openstack_sharedfilesystem_share_access_v2":         resourceSharedFilesystemShareAccessV2(),
			"openstack_keymanager_secret_v1":                     resourceKeyManagerSecretV1(),
			"openstack_keymanager_container_v1":                  resourceKeyManagerContainerV1(),
			"openstack_keymanager_order_v1":                      resourceKeyManagerOrderV1(),
			"openstack_bgpvpn_v2":                                resourceBGPVPNV2(),
			"openstack_bgpvpn_network_associate_v2":              resourceBGPVPNNetworkAssociateV2(),
			"openstack_bgpvpn_router_associate_v2":               resourceBGPVPNRouterAssociateV2(),
			"openstack_bgpvpn_port_associate_v2":                 resourceBGPVPNPortAssociateV2(),
			"openstack_workflow_cron_trigger_v2":                 resourceWorkflowCronTriggerV2(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		return configureProvider(ctx, provider, d)
	}

	return provider
}

func getSDKVersion() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return version.Version
	}

	for _, v := range buildInfo.Deps {
		if v.Path == "github.com/hashicorp/terraform-plugin-sdk/v2" {
			return v.Version
		}
	}

	return version.Version
}

func configureProvider(ctx context.Context, provider *schema.Provider, d *schema.ResourceData) (any, diag.Diagnostics) {
	// If the provider has already been configured, return the existing config.
	if v, ok := provider.Meta().(*Config); ok {
		return v, nil
	}

	terraformVersion := provider.TerraformVersion
	if terraformVersion == "" {
		// Terraform 0.12 introduced this field to the protocol
		// We can therefore assume that if it's missing it's 0.10 or 0.11
		terraformVersion = "0.11+compatible"
	}

	enableLogging := d.Get("enable_logging").(bool)
	if !enableLogging {
		// enforce logging (similar to OS_DEBUG) when TF_LOG is 'DEBUG' or 'TRACE'
		if logLevel := logging.LogLevel(); logLevel != "" && os.Getenv("OS_DEBUG") == "" {
			if logLevel == "DEBUG" || logLevel == "TRACE" {
				enableLogging = true
			}
		}
	}

	authOpts := &gophercloud.AuthOptions{
		Scope: &gophercloud.AuthScope{System: d.Get("system_scope").(bool)},
	}

	config := Config{
		auth.Config{
			CACertFile:                  d.Get("cacert_file").(string),
			ClientCertFile:              d.Get("cert").(string),
			ClientKeyFile:               d.Get("key").(string),
			Cloud:                       d.Get("cloud").(string),
			DefaultDomain:               d.Get("default_domain").(string),
			DomainID:                    d.Get("domain_id").(string),
			DomainName:                  d.Get("domain_name").(string),
			EndpointOverrides:           d.Get("endpoint_overrides").(map[string]any),
			EndpointType:                d.Get("endpoint_type").(string),
			IdentityEndpoint:            d.Get("auth_url").(string),
			Password:                    d.Get("password").(string),
			ProjectDomainID:             d.Get("project_domain_id").(string),
			ProjectDomainName:           d.Get("project_domain_name").(string),
			Region:                      d.Get("region").(string),
			Swauth:                      d.Get("swauth").(bool),
			Token:                       d.Get("token").(string),
			TenantID:                    d.Get("tenant_id").(string),
			TenantName:                  d.Get("tenant_name").(string),
			UserDomainID:                d.Get("user_domain_id").(string),
			UserDomainName:              d.Get("user_domain_name").(string),
			Username:                    d.Get("user_name").(string),
			UserID:                      d.Get("user_id").(string),
			UseOctavia:                  true,
			ApplicationCredentialID:     d.Get("application_credential_id").(string),
			ApplicationCredentialName:   d.Get("application_credential_name").(string),
			ApplicationCredentialSecret: d.Get("application_credential_secret").(string),
			DelayedAuth:                 d.Get("delayed_auth").(bool),
			AllowReauth:                 d.Get("allow_reauth").(bool),
			AuthOpts:                    authOpts,
			MaxRetries:                  d.Get("max_retries").(int),
			DisableNoCacheHeader:        d.Get("disable_no_cache_header").(bool),
			TerraformVersion:            terraformVersion,
			SDKVersion:                  getSDKVersion() + " Terraform Provider OpenStack/" + version.Version,
			MutexKV:                     mutexkv.NewMutexKV(),
			EnableLogger:                enableLogging,
		},
	}

	v, ok := getOkExists(d, "insecure")
	if ok {
		insecure := v.(bool)
		config.Insecure = &insecure
	}

	if err := config.LoadAndValidate(ctx); err != nil {
		return nil, diag.FromErr(err)
	}

	return &config, nil
}
