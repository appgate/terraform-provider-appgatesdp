package appgate

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v19/openapi"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateSiteCreate,
		Read:   resourceAppgateSiteRead,
		Update: resourceAppgateSiteUpdate,
		Delete: resourceAppgateSiteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{

			"site_id": resourceUUID(),

			"name": {
				Type:        schema.TypeString,
				Description: "Name of the object.",
				Required:    true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"notes": {
				Type:        schema.TypeString,
				Description: "Notes for the object. Used for documentation purposes.",
				Default:     DefaultDescription,
				Optional:    true,
			},

			"tags": tagsSchema(),

			"short_name": {
				Type:        schema.TypeString,
				Description: "A short 4 letter name for the site",
				Optional:    true,
			},

			"network_subnets": {
				Type:        schema.TypeSet,
				Description: "Network subnets in CIDR format to define the Site's boundaries. They are added as routes by the Client.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"ip_pool_mappings": {
				Type:       schema.TypeSet,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"from": {
							Type:     schema.TypeString,
							Required: true,
						},

						"to": {
							Type:     schema.TypeString,
							Required: true,
						},

						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"default_gateway": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"enabled_v4": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"enabled_v6": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"excluded_subnets": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"entitlement_based_routing": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"vpn": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"state_sharing": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"snat": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"tls": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},

						"dtls": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},

						"route_via": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ipv6": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"ip_access_log_interval_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"name_resolution": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"use_hosts_file": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},

						"dns_resolvers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"update_interval": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"query_aaaa": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"default_ttl_seconds": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"servers": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"search_domains": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},

						"aws_resolvers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"update_interval": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"vpcs": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"vpc_auto_discovery": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"regions": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"use_iam_role": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"access_key_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"secret_access_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"https_proxy": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"resolve_with_master_credentials": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"assumed_roles": {
										Type:     schema.TypeList,
										Optional: true,

										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{

												"account_id": {
													Type:     schema.TypeString,
													Optional: true,
												},

												"role_name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"external_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"regions": {
													Type:     schema.TypeList,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
								},
							},
						},

						"azure_resolvers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"update_interval": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"use_managed_identities": {
										Type:     schema.TypeBool,
										Default:  false,
										Optional: true,
									},
									"subscription_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"tenant_id": {
										Type:     schema.TypeString,
										Required: true,
									},

									"client_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"secret": {
										Type:      schema.TypeString,
										Optional:  true,
										Sensitive: true,
									},
								},
							},
						},

						"esx_resolvers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"update_interval": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"hostname": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"username": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
									"password": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
								},
							},
						},

						"gcp_resolvers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"update_interval": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"project_filter": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"instance_filter": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"dns_forwarding": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"site_ipv4": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"site_ipv6": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"dns_servers": {
										Type:     schema.TypeSet,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"allow_destinations": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:     schema.TypeString,
													Required: true,
												},
												"netmask": {
													Type:     schema.TypeInt,
													Required: true,
												},
											},
										},
									},
									"default_ttl_seconds": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},

						"illumio_resolvers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"update_interval": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"hostname": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"username": {
										Type:     schema.TypeString,
										Required: true,
									},
									"password": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
									"org_id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAppgateSiteCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating Site: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.SitesApi
	currentVersion := meta.(*Client).ApplianceVersion
	args := openapi.Site{}
	if v, ok := d.GetOk("site_id"); ok {
		args.SetId(v.(string))
	}
	args.SetName(d.Get("name").(string))
	args.SetShortName(d.Get("short_name").(string))
	args.SetDescription(d.Get("description").(string))
	args.SetNotes(d.Get("notes").(string))
	args.SetTags(schemaExtractTags(d))

	if v, ok := d.GetOk("network_subnets"); ok {
		networkSubnets, err := readArrayOfStringsFromConfig(v.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetNetworkSubnets(networkSubnets)
	}

	if v, ok := d.GetOk("ip_pool_mappings"); ok {
		ipPoolMappings, err := readIPPoolMappingsFromConfig(v.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetIpPoolMappings(ipPoolMappings)
	}

	if v, ok := d.GetOk("default_gateway"); ok {
		DefaultGateway, err := readSiteDefaultGatewayFromConfig(v.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetDefaultGateway(DefaultGateway)
	}

	if v, ok := d.GetOk("entitlement_based_routing"); ok {
		args.SetEntitlementBasedRouting(v.(bool))
	}

	if v, ok := d.GetOk("vpn"); ok {
		vpn, err := readSiteVPNFromConfig(currentVersion, v.([]interface{}))
		if err != nil {
			return err
		}
		args.SetVpn(vpn)
	}

	if v, ok := d.GetOk("name_resolution"); ok {
		nameResolution, err := readSiteNameResolutionFromConfig(currentVersion, v.([]interface{}))
		if err != nil {
			return err
		}
		args.SetNameResolution(nameResolution)
	}

	site, _, err := api.SitesPost(context.Background()).Site(args).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create site %w", prettyPrintAPIError(err))
	}

	d.SetId(site.GetId())

	return resourceAppgateSiteRead(d, meta)
}

func resourceAppgateSiteRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Site Name: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.SitesApi
	currentVersion := meta.(*Client).ApplianceVersion

	request := api.SitesIdGet(context.Background(), d.Id())
	site, res, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Failed to read Site, %w", err)
	}

	d.SetId(site.GetId())
	d.Set("site_id", site.GetId())
	d.Set("name", site.GetName())
	d.Set("description", site.GetDescription())
	d.Set("notes", site.GetNotes())
	d.Set("tags", site.GetTags())
	d.Set("network_subnets", site.NetworkSubnets)
	if site.IpPoolMappings != nil {
		if err = d.Set("ip_pool_mappings", flattenSiteIPpoolmappning(site.GetIpPoolMappings())); err != nil {
			return err
		}
	}
	if site.DefaultGateway != nil {
		if err = d.Set("default_gateway", flattenSiteDefaultGateway(*site.DefaultGateway)); err != nil {
			return err
		}
	}
	d.Set("short_name", site.ShortName)
	d.Set("entitlement_based_routing", site.EntitlementBasedRouting)

	if site.Vpn != nil {
		if err = d.Set("vpn", flattenSiteVPN(currentVersion, *site.Vpn)); err != nil {
			return err
		}
	}

	if site.NameResolution != nil {
		var localNameResolution map[string]interface{}
		localNameResolutionList := d.Get("name_resolution").([]interface{})
		for _, l := range localNameResolutionList {
			localNameResolution = l.(map[string]interface{})
		}
		ns, err := flattenNameResolution(currentVersion, localNameResolution, *site.NameResolution)
		if err != nil {
			return err
		}
		if err = d.Set("name_resolution", ns); err != nil {
			return err
		}
	}

	return nil
}

func flattenSiteIPpoolmappning(in []openapi.SiteAllOfIpPoolMappings) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["from"] = v.GetFrom()
		m["to"] = v.GetTo()
		m["type"] = v.GetType()

		out[i] = m
	}
	return out
}

func flattenSiteDefaultGateway(in openapi.SiteAllOfDefaultGateway) []interface{} {
	m := make(map[string]interface{})
	m["enabled_v4"] = in.EnabledV4
	m["enabled_v6"] = in.EnabledV6
	exsub := make([]interface{}, 0, 0)
	for _, sub := range in.GetExcludedSubnets() {
		exsub = append(exsub, sub)
	}
	m["excluded_subnets"] = exsub
	return []interface{}{m}
}

func flattenSiteVPN(currentVersion *version.Version, in openapi.SiteAllOfVpn) []interface{} {
	m := make(map[string]interface{})
	if v, ok := in.GetStateSharingOk(); ok {
		m["state_sharing"] = *v
	}
	if v, ok := in.GetSnatOk(); ok {
		m["snat"] = *v
	}

	if in.HasDtls() {
		dtls := make(map[string]interface{})
		if _, ok := in.Dtls.GetEnabledOk(); ok {
			dtls["enabled"] = in.Dtls.GetEnabled()
		}
		m["dtls"] = []interface{}{dtls}
	}

	if in.HasTls() {
		tls := make(map[string]interface{})
		if _, ok := in.Tls.GetEnabledOk(); ok {
			tls["enabled"] = in.Tls.GetEnabled()
		}
		m["tls"] = []interface{}{tls}
	}
	if in.HasRouteVia() && (in.RouteVia.Ipv4 != nil || in.RouteVia.Ipv6 != nil) {
		routeVia := make(map[string]interface{})
		if v, o := in.RouteVia.GetIpv4Ok(); o && len(*v) > 0 {
			routeVia["ipv4"] = v
		}
		if v, o := in.RouteVia.GetIpv6Ok(); o && len(*v) > 0 {
			routeVia["ipv6"] = v
		}
		m["route_via"] = []interface{}{routeVia}
	}

	m["ip_access_log_interval_seconds"] = in.IpAccessLogIntervalSeconds

	return []interface{}{m}
}

func flattenNameResolution(currentVersion *version.Version, local map[string]interface{}, in openapi.SiteAllOfNameResolution) ([]interface{}, error) {
	m := make(map[string]interface{})
	if v, ok := in.GetUseHostsFileOk(); ok {
		m["use_hosts_file"] = *v
	}

	if v, ok := in.GetDnsResolversOk(); ok {
		m["dns_resolvers"] = flattenSiteDNSResolver(v)
	}
	if v, ok := in.GetAwsResolversOk(); ok {
		l := getNSLocalChanges(local, "aws_resolvers")
		m["aws_resolvers"] = flattenSiteAWSResolver(v, l)
	}
	if v, ok := in.GetAzureResolversOk(); ok {
		l := getNSLocalChanges(local, "azure_resolvers")
		m["azure_resolvers"] = flattenSiteAzureResolver(currentVersion, v, l)
	}
	if v, ok := in.GetEsxResolversOk(); ok {
		l := getNSLocalChanges(local, "esx_resolvers")
		m["esx_resolvers"] = flattenSiteESXResolvers(v, l)
	}
	if v, ok := in.GetGcpResolversOk(); ok {
		m["gcp_resolvers"] = flattenSiteGCPResolvers(v)
	}
	if currentVersion.GreaterThanOrEqual(Appliance55Version) {
		if v, ok := in.GetDnsForwardingOk(); ok {
			dnsfwd, err := flattenSiteDnsForwading(*v)
			if err != nil {
				return nil, err
			}
			m["dns_forwarding"] = dnsfwd
		}
	}
	if currentVersion.GreaterThanOrEqual(Appliance61Version) {
		if v, ok := in.GetIllumioResolversOk(); ok {
			m["illumio_resolvers"] = flattenSiteIllumioResolvers(v, getNSLocalChanges(local, "illumio_resolvers"))
		}
	}
	return []interface{}{m}, nil
}

func getNSLocalChanges(local map[string]interface{}, name string) map[string]interface{} {
	if d, ok := local[name]; ok {
		for _, p := range (d.(*schema.Set)).List() {
			return p.(map[string]interface{})
		}
	}
	return make(map[string]interface{})
}

func flattenSiteGCPResolvers(in []openapi.SiteAllOfNameResolutionGcpResolvers) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["name"] = v.GetName()
		m["update_interval"] = v.GetUpdateInterval()
		m["project_filter"] = v.GetProjectFilter()
		m["instance_filter"] = v.GetInstanceFilter()
		out[i] = m
	}
	return out
}

func flattenSiteIllumioResolvers(in []openapi.SiteAllOfNameResolutionIllumioResolvers, local map[string]interface{}) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["name"] = v.GetName()
		m["update_interval"] = v.GetUpdateInterval()
		m["hostname"] = v.GetHostname()
		m["port"] = v.GetPort()
		m["username"] = v.GetUsername()
		m["org_id"] = v.GetOrgId()
		if val, ok := local["password"]; ok {
			m["password"] = val
		} else {
			m["password"] = v.GetPassword()
		}
		out[i] = m
	}
	return out
}

func flattenSiteESXResolvers(in []openapi.SiteAllOfNameResolutionEsxResolvers, local map[string]interface{}) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["name"] = v.GetName()
		m["update_interval"] = v.GetUpdateInterval()
		m["hostname"] = v.GetHostname()
		m["username"] = v.GetUsername()
		if val, ok := local["password"]; ok {
			m["password"] = val
		} else {
			m["password"] = v.GetPassword()
		}
		out[i] = m
	}
	return out
}

func flattenSiteAzureResolver(currentVersion *version.Version, in []openapi.SiteAllOfNameResolutionAzureResolvers, local map[string]interface{}) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["name"] = v.GetName()
		m["update_interval"] = v.GetUpdateInterval()
		m["subscription_id"] = v.GetSubscriptionId()
		m["tenant_id"] = v.GetTenantId()
		m["client_id"] = v.GetClientId()
		if val, ok := local["secret"]; ok {
			m["secret"] = val
		} else {
			m["secret"] = v.GetSecret()
		}
		if currentVersion.GreaterThanOrEqual(Appliance55Version) {
			m["use_managed_identities"] = v.GetUseManagedIdentities()
		}

		out[i] = m
	}
	return out
}

func flattenSiteAWSResolver(in []openapi.SiteAllOfNameResolutionAwsResolvers, local map[string]interface{}) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["name"] = v.GetName()
		m["update_interval"] = v.GetUpdateInterval()
		m["vpcs"] = v.GetVpcs()
		m["vpc_auto_discovery"] = v.GetVpcAutoDiscovery()
		m["regions"] = v.GetRegions()
		m["use_iam_role"] = v.GetUseIAMRole()
		m["access_key_id"] = v.GetAccessKeyId()
		if val, ok := local["secret_access_key"]; ok {
			m["secret_access_key"] = val
		} else {
			m["secret_access_key"] = v.GetSecretAccessKey()
		}
		m["https_proxy"] = v.GetHttpsProxy()
		m["resolve_with_master_credentials"] = v.GetResolveWithMasterCredentials()
		if vv, o := v.GetAssumedRolesOk(); o != false {
			m["assumed_roles"] = flattenSiteAwsAssumedRoles(vv)
		}
		out[i] = m
	}
	return out
}

func flattenSiteAwsAssumedRoles(in []openapi.SiteAllOfNameResolutionAssumedRoles) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["account_id"] = v.GetAccountId()
		m["role_name"] = v.GetRoleName()
		m["external_id"] = v.GetExternalId()
		m["regions"] = v.GetRegions()
		out[i] = m
	}
	return out
}

func flattenSiteDNSResolver(in []openapi.SiteAllOfNameResolutionDnsResolvers) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["name"] = v.GetName()
		m["update_interval"] = v.GetUpdateInterval()
		m["query_aaaa"] = v.GetQueryAAAA()
		m["default_ttl_seconds"] = v.GetDefaultTtlSeconds()
		m["servers"] = v.GetServers()
		m["search_domains"] = v.GetSearchDomains()

		out[i] = m
	}
	return out
}

func flattenSiteDnsForwading(in openapi.SiteAllOfNameResolutionDnsForwarding) ([]map[string]interface{}, error) {
	m := make(map[string]interface{})
	m["site_ipv4"] = in.GetSiteIpv4()
	m["site_ipv6"] = in.GetSiteIpv6()
	m["dns_servers"] = in.GetDnsServers()
	ad, err := flattenAllowDestinations(in.GetAllowDestinations())
	if err != nil {
		return nil, err
	}
	m["allow_destinations"] = ad
	m["default_ttl_seconds"] = in.GetDefaultTtlSeconds()

	return []map[string]interface{}{m}, nil
}

func flattenAllowDestinations(input []openapi.AllowResourcesInner) ([]map[string]interface{}, error) {
	r := make([]map[string]interface{}, 0)
	for _, raw := range input {
		row := make(map[string]interface{}, 0)
		if v, ok := raw.GetAddressOk(); ok {
			row["address"] = *v
		}
		if v, ok := raw.GetNetmaskOk(); ok {
			row["netmask"] = *v
		}
		r = append(r, row)
	}
	return r, nil
}

func resourceAppgateSiteUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Site: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.SitesApi
	currentVersion := meta.(*Client).ApplianceVersion
	request := api.SitesIdGet(context.Background(), d.Id())

	orginalSite, res, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Failed to read Site, %w", err)
	}

	if d.HasChange("name") {
		orginalSite.SetName(d.Get("name").(string))
	}
	if d.HasChange("short_name") {
		orginalSite.SetShortName(d.Get("short_name").(string))
	}
	if d.HasChange("notes") {
		orginalSite.SetNotes(d.Get("notes").(string))
	}
	if d.HasChange("tags") {
		orginalSite.SetTags(schemaExtractTags(d))
	}
	if d.HasChange("entitlement_based_routing") {
		_, v := d.GetChange("entitlement_based_routing")
		orginalSite.SetEntitlementBasedRouting(v.(bool))
	}
	if d.HasChange("network_subnets") {
		_, n := d.GetChange("network_subnets")
		networkSubnets, err := readArrayOfStringsFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		orginalSite.SetNetworkSubnets(networkSubnets)
	}

	if d.HasChange("network_subnets") {
		_, n := d.GetChange("network_subnets")
		networkSubnets, err := readArrayOfStringsFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		orginalSite.SetNetworkSubnets(networkSubnets)
	}

	if d.HasChange("ip_pool_mappings") {
		_, n := d.GetChange("ip_pool_mappings")
		ipPoolMappings, err := readIPPoolMappingsFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		orginalSite.SetIpPoolMappings(ipPoolMappings)
	}

	if d.HasChange("default_gateway") {
		_, n := d.GetChange("default_gateway")
		DefaultGateway, err := readSiteDefaultGatewayFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		orginalSite.SetDefaultGateway(DefaultGateway)
	}

	if d.HasChange("vpn") {
		_, v := d.GetChange("vpn")
		vpn, err := readSiteVPNFromConfig(currentVersion, v.([]interface{}))
		if err != nil {
			return err
		}
		orginalSite.SetVpn(vpn)
	}

	if d.HasChange("name_resolution") {
		_, v := d.GetChange("name_resolution")
		nameResolution, err := readSiteNameResolutionFromConfig(currentVersion, v.([]interface{}))
		if err != nil {
			return err
		}
		orginalSite.SetNameResolution(nameResolution)
	}

	putRequest := api.SitesIdPut(context.Background(), d.Id())
	_, _, err = putRequest.Site(*orginalSite).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update site %w", prettyPrintAPIError(err))
	}
	return resourceAppgateSiteRead(d, meta)
}

func resourceAppgateSiteDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete Site: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.SitesApi

	if _, err := api.SitesIdDelete(context.Background(), d.Id()).Authorization(token).Execute(); err != nil {
		return fmt.Errorf("Failed to delete Site, %w", err)
	}
	d.SetId("")
	return nil
}

func readIPPoolMappingsFromConfig(maps []interface{}) ([]openapi.SiteAllOfIpPoolMappings, error) {
	result := make([]openapi.SiteAllOfIpPoolMappings, 0)
	for _, ipPool := range maps {
		if ipPool == nil {
			continue
		}
		r := openapi.SiteAllOfIpPoolMappings{}
		raw := ipPool.(map[string]interface{})
		if v, ok := raw["from"]; ok {
			r.SetFrom(v.(string))
		}
		if v, ok := raw["to"]; ok {
			r.SetTo(v.(string))
		}
		if v, ok := raw["type"]; ok {
			r.SetType(v.(string))
		}

		result = append(result, r)
	}
	return result, nil
}

func readSiteDefaultGatewayFromConfig(defaultGateways []interface{}) (openapi.SiteAllOfDefaultGateway, error) {
	result := openapi.SiteAllOfDefaultGateway{}
	for _, defaultGateway := range defaultGateways {
		if defaultGateway == nil {
			continue
		}
		raw := defaultGateway.(map[string]interface{})
		if v, ok := raw["enabled_v4"]; ok {
			result.SetEnabledV4(v.(bool))
		}
		if v, ok := raw["enabled_v6"]; ok {
			result.SetEnabledV6(v.(bool))
		}

		if v := raw["excluded_subnets"]; len(v.([]interface{})) > 0 {
			nets, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, fmt.Errorf("Failed to resolve default gateway excluded subnets: %w", err)
			}
			result.SetExcludedSubnets(nets)
		}

	}
	return result, nil
}

func readSiteVPNFromConfig(currentVersion *version.Version, vpns []interface{}) (openapi.SiteAllOfVpn, error) {
	result := openapi.SiteAllOfVpn{}
	for _, vpn := range vpns {
		if vpn == nil {
			continue
		}
		raw := vpn.(map[string]interface{})

		if v, ok := raw["state_sharing"]; ok {
			result.SetStateSharing(v.(bool))
		}
		if v, ok := raw["snat"]; ok {
			result.SetSnat(v.(bool))
		}
		if v, ok := raw["tls"]; ok {
			tls := openapi.NewSiteAllOfVpnTlsWithDefaults()
			rawTLS := v.(*schema.Set).List()
			for _, d := range rawTLS {
				raw := d.(map[string]interface{})
				if v, ok := raw["enabled"]; ok {
					tls.SetEnabled(v.(bool))
				}
			}
			result.SetTls(*tls)
		}

		if v, ok := raw["dtls"]; ok {
			dtls := openapi.NewSiteAllOfVpnDtlsWithDefaults()
			rawDTLS := v.(*schema.Set).List()
			for _, d := range rawDTLS {
				raw := d.(map[string]interface{})

				if v, ok := raw["enabled"]; ok {
					dtls.SetEnabled(v.(bool))
				}
			}

			result.SetDtls(*dtls)
		}

		if v, ok := raw["route_via"]; ok {
			routeVia, err := readSiteVPNRouteViaFromConfig(v.(*schema.Set).List())
			if err != nil {
				return result, fmt.Errorf("Failed to parse/retrieve route_via config %w", err)
			}
			result.SetRouteVia(routeVia)
		}

		if v, ok := raw["ip_access_log_interval_seconds"]; ok {
			result.SetIpAccessLogIntervalSeconds(float32(v.(int)))
		}
	}
	return result, nil
}

func readSiteVPNRouteViaFromConfig(routeViaConf []interface{}) (openapi.SiteAllOfVpnRouteVia, error) {
	result := openapi.SiteAllOfVpnRouteVia{}
	for _, r := range routeViaConf {
		if r == nil {
			continue
		}

		raw := r.(map[string]interface{})
		if v, ok := raw["ipv4"].(string); ok && len(v) > 0 {
			result.SetIpv4(v)
		}
		if v, ok := raw["ipv6"].(string); ok && len(v) > 0 {
			result.SetIpv6(v)
		}
	}
	return result, nil
}

func readSiteNameResolutionFromConfig(currentVersion *version.Version, nameresolutions []interface{}) (openapi.SiteAllOfNameResolution, error) {
	result := openapi.SiteAllOfNameResolution{}
	for _, nr := range nameresolutions {
		if nr == nil {
			continue
		}
		raw := nr.(map[string]interface{})
		if v, ok := raw["use_hosts_file"]; ok {
			result.SetUseHostsFile(v.(bool))
		}
		if v, ok := raw["dns_resolvers"]; ok {
			dnss := v.(*schema.Set).List()
			dnsResolvers, err := readDNSResolversFromConfig(currentVersion, dnss)
			if err != nil {
				return result, err
			}
			result.SetDnsResolvers(dnsResolvers)
		}
		if v, ok := raw["aws_resolvers"]; ok {
			awsResolvers, err := readAWSResolversFromConfig(v.(*schema.Set).List())
			if err != nil {
				return result, err
			}
			result.SetAwsResolvers(awsResolvers)
		}
		if v, ok := raw["azure_resolvers"]; ok {
			azureResolvers, err := readAzureResolversFromConfig(currentVersion, v.(*schema.Set).List())
			if err != nil {
				return result, err
			}
			result.SetAzureResolvers(azureResolvers)
		}
		if v, ok := raw["esx_resolvers"]; ok {
			esxResolvers, err := readESXResolversFromConfig(v.(*schema.Set).List())
			if err != nil {
				return result, err
			}
			result.SetEsxResolvers(esxResolvers)
		}
		if v, ok := raw["gcp_resolvers"]; ok {
			gcpResolvers, err := readGCPResolversFromConfig(v.(*schema.Set).List())
			if err != nil {
				return result, err
			}
			result.SetGcpResolvers(gcpResolvers)
		}
		if v, ok := raw["dns_forwarding"]; ok {
			dnsForwardingResolvers, err := readDNSForwardingResolversFromConfig(currentVersion, v.(*schema.Set).List())
			if err != nil {
				return result, err
			}
			if len(dnsForwardingResolvers.GetDnsServers()) > 0 {
				if currentVersion.LessThan(Appliance55Version) {
					return result, errors.New("dns_forwarding is only available in 5.5 or above")
				}
				result.SetDnsForwarding(dnsForwardingResolvers)
			}
		}
		if currentVersion.GreaterThanOrEqual(Appliance61Version) {
			if v, ok := raw["illumio_resolvers"]; ok {
				resolvers, err := readIllumioResolversFromConfig(v.(*schema.Set).List())
				if err != nil {
					return result, err
				}
				result.SetIllumioResolvers(resolvers)
			}
		}
	}
	return result, nil
}

func readDNSResolversFromConfig(currentVersion *version.Version, dnsConfigs []interface{}) ([]openapi.SiteAllOfNameResolutionDnsResolvers, error) {
	result := make([]openapi.SiteAllOfNameResolutionDnsResolvers, 0)
	for _, dns := range dnsConfigs {
		raw := dns.(map[string]interface{})
		row := openapi.SiteAllOfNameResolutionDnsResolvers{}
		if v, ok := raw["name"]; ok {
			row.SetName(v.(string))
		}
		if v, ok := raw["update_interval"]; ok {
			row.SetUpdateInterval(int32(v.(int)))
		}
		if currentVersion.GreaterThanOrEqual(Appliance60Version) {
			if v, ok := raw["query_aaaa"]; ok {
				row.SetQueryAAAA(v.(bool))
			}
			if v, ok := raw["default_ttl_seconds"].(int); ok && v > 0 {
				row.SetDefaultTtlSeconds(int32(v))
			}
		}
		if v := raw["servers"]; len(v.([]interface{})) > 0 {
			servers, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, fmt.Errorf("Failed to resolve dns servers: %w", err)
			}
			if len(servers) > 0 {
				row.SetServers(servers)
			}
		}
		if v := raw["search_domains"]; len(v.([]interface{})) > 0 {
			domains, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, fmt.Errorf("Failed to resolve dns search domains: %w", err)
			}
			if len(domains) > 0 {
				row.SetSearchDomains(domains)
			}
		}
		result = append(result, row)
	}
	return result, nil
}

func readAWSResolversFromConfig(awsConfigs []interface{}) ([]openapi.SiteAllOfNameResolutionAwsResolvers, error) {
	result := make([]openapi.SiteAllOfNameResolutionAwsResolvers, 0)
	for _, resolver := range awsConfigs {
		raw := resolver.(map[string]interface{})
		row := openapi.SiteAllOfNameResolutionAwsResolvers{}

		if v, ok := raw["name"]; ok {
			row.SetName(v.(string))
		}
		if v, ok := raw["update_interval"]; ok && v != nil {
			row.SetUpdateInterval(int32(v.(int)))
		}
		if v := raw["vpcs"]; len(v.([]interface{})) > 0 {
			vpcs, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, fmt.Errorf("Failed to resolve vpcs from aws config: %w", err)
			}
			row.SetVpcs(vpcs)
		}
		if v, ok := raw["vpc_auto_discovery"]; ok && v.(bool) {
			row.SetVpcAutoDiscovery(v.(bool))
		}
		if v := raw["regions"]; len(v.([]interface{})) > 0 && v.([]interface{}) != nil {
			regions, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, fmt.Errorf("Failed to resolve regions from aws config: %w", err)
			}
			row.SetRegions(regions)
		}
		if v, ok := raw["use_iam_role"]; ok && v.(bool) {
			row.SetUseIAMRole(v.(bool))
		}
		if v, ok := raw["access_key_id"].(string); ok && len(v) > 0 {
			row.SetAccessKeyId(v)
		}
		if v, ok := raw["secret_access_key"].(string); ok && len(v) > 0 {
			row.SetSecretAccessKey(v)
		}
		if v, ok := raw["https_proxy"]; ok && len(v.(string)) > 0 {
			row.SetHttpsProxy(v.(string))
		}
		if v, ok := raw["resolve_with_master_credentials"]; ok {
			row.SetResolveWithMasterCredentials(v.(bool))
		}
		if v, ok := raw["assumed_roles"]; ok {
			assumedRoles, err := readAwsAssumedRolesFromConfig(v.([]interface{}))
			if err != nil {
				return result, err
			}
			if len(assumedRoles) > 0 {
				row.SetAssumedRoles(assumedRoles)
			}
		}
		result = append(result, row)
	}
	return result, nil
}

func readAwsAssumedRolesFromConfig(roles []interface{}) ([]openapi.SiteAllOfNameResolutionAssumedRoles, error) {
	result := make([]openapi.SiteAllOfNameResolutionAssumedRoles, 0)
	for _, role := range roles {
		raw := role.(map[string]interface{})
		row := openapi.SiteAllOfNameResolutionAssumedRoles{}
		if v, ok := raw["account_id"]; ok {
			row.SetAccountId(v.(string))
		}
		if v, ok := raw["role_name"]; ok {
			row.SetRoleName(v.(string))
		}
		if v, ok := raw["external_id"]; ok {
			row.SetExternalId(v.(string))
		}
		if v, ok := raw["regions"]; ok {
			regions, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, err
			}
			row.SetRegions(regions)
		}
		result = append(result, row)
	}
	return result, nil
}

func readAzureResolversFromConfig(currentVersion *version.Version, azureConfigs []interface{}) ([]openapi.SiteAllOfNameResolutionAzureResolvers, error) {
	result := make([]openapi.SiteAllOfNameResolutionAzureResolvers, 0)
	for _, azure := range azureConfigs {
		raw := azure.(map[string]interface{})
		row := openapi.SiteAllOfNameResolutionAzureResolvers{}
		if v, ok := raw["name"]; ok {
			row.SetName(v.(string))
		}
		if v, ok := raw["update_interval"]; ok {
			row.SetUpdateInterval(int32(v.(int)))
		}
		if v, ok := raw["subscription_id"]; ok {
			row.SetSubscriptionId(v.(string))
		}
		if v, ok := raw["tenant_id"]; ok {
			row.SetTenantId(v.(string))
		}
		if v, ok := raw["client_id"]; ok {
			row.SetClientId(v.(string))
		}
		if v, ok := raw["secret"]; ok {
			row.SetSecret(v.(string))
		}
		if currentVersion.GreaterThanOrEqual(Appliance55Version) {
			if v, ok := raw["use_managed_identities"]; ok {
				row.SetUseManagedIdentities(v.(bool))
			}
		}
		result = append(result, row)
	}
	return result, nil
}

func readESXResolversFromConfig(esxConfigs []interface{}) ([]openapi.SiteAllOfNameResolutionEsxResolvers, error) {
	result := make([]openapi.SiteAllOfNameResolutionEsxResolvers, 0)
	for _, esxConfig := range esxConfigs {
		raw := esxConfig.(map[string]interface{})
		row := openapi.SiteAllOfNameResolutionEsxResolvers{}
		if v, ok := raw["name"]; ok {
			row.SetName(v.(string))
		}
		if v, ok := raw["update_interval"]; ok {
			row.SetUpdateInterval(int32(v.(int)))
		}
		if v, ok := raw["hostname"]; ok {
			row.SetHostname(v.(string))
		}
		if v, ok := raw["username"]; ok {
			row.SetUsername(v.(string))
		}
		if v, ok := raw["password"]; ok {
			row.SetPassword(v.(string))
		}
		result = append(result, row)
	}
	return result, nil
}

func readGCPResolversFromConfig(gcpConfigs []interface{}) ([]openapi.SiteAllOfNameResolutionGcpResolvers, error) {
	result := make([]openapi.SiteAllOfNameResolutionGcpResolvers, 0)
	for _, gcpConfig := range gcpConfigs {
		raw := gcpConfig.(map[string]interface{})
		row := openapi.SiteAllOfNameResolutionGcpResolvers{}
		if v, ok := raw["name"]; ok {
			row.SetName(v.(string))
		}
		if v, ok := raw["update_interval"]; ok {
			row.SetUpdateInterval(int32(v.(int)))
		}
		if v, ok := raw["project_filter"]; ok {
			row.SetProjectFilter(v.(string))
		}
		if v, ok := raw["instance_filter"]; ok {
			row.SetInstanceFilter(v.(string))
		}
		result = append(result, row)
	}
	return result, nil
}

func readDNSForwardingResolversFromConfig(currentVersion *version.Version, dnsForwardingConfig []interface{}) (openapi.SiteAllOfNameResolutionDnsForwarding, error) {
	result := openapi.SiteAllOfNameResolutionDnsForwarding{}
	for _, dnsForwarding := range dnsForwardingConfig {
		raw := dnsForwarding.(map[string]interface{})
		if v, ok := raw["site_ipv4"]; ok {
			result.SetSiteIpv4(v.(string))
		}
		if v, ok := raw["site_ipv6"]; ok {
			result.SetSiteIpv6(v.(string))
		}
		if v, ok := raw["dns_servers"]; ok {
			servers, err := readArrayOfStringsFromConfig(v.(*schema.Set).List())
			if err != nil {
				return result, err
			}
			result.SetDnsServers(servers)
		}
		if v, ok := raw["allow_destinations"]; ok {
			destinations, err := readAllowDestinationsFromConfig(v.(*schema.Set).List())
			if err != nil {
				return result, err
			}
			result.SetAllowDestinations(destinations)
		}
		if currentVersion.GreaterThanOrEqual(Appliance60Version) {
			if v, ok := raw["default_ttl_seconds"].(int); ok && v > 0 {
				result.SetDefaultTtlSeconds(int32(v))
			}
		}
	}
	return result, nil
}

func readIllumioResolversFromConfig(resolvers []interface{}) ([]openapi.SiteAllOfNameResolutionIllumioResolvers, error) {
	result := make([]openapi.SiteAllOfNameResolutionIllumioResolvers, 0)
	for _, illumio := range resolvers {
		raw := illumio.(map[string]interface{})
		row := openapi.SiteAllOfNameResolutionIllumioResolvers{}
		if v, ok := raw["name"]; ok {
			row.SetName(v.(string))
		}
		if v, ok := raw["update_interval"]; ok {
			row.SetUpdateInterval(int32(v.(int)))
		}
		if v, ok := raw["hostname"]; ok {
			row.SetHostname(v.(string))
		}
		if v, ok := raw["port"]; ok {
			row.SetPort(int32(v.(int)))
		}
		if v, ok := raw["username"]; ok {
			row.SetUsername(v.(string))
		}
		if v, ok := raw["password"]; ok {
			row.SetPassword(v.(string))
		}
		if currentVersion.GreaterThanOrEqual(Appliance62Version) {
			if v, ok := raw["org_id"]; ok {
				row.SetOrgId(v.(string))
			}
		}
		result = append(result, row)
	}
	return result, nil
}

func readAllowDestinationsFromConfig(input []interface{}) ([]openapi.AllowResourcesInner, error) {
	result := make([]openapi.AllowResourcesInner, 0)
	for _, r := range input {
		raw := r.(map[string]interface{})
		row := openapi.AllowResourcesInner{}
		if v, ok := raw["address"].(string); ok {
			row.SetAddress(v)
		}
		if v, ok := raw["netmask"].(int); ok {
			row.SetNetmask(int32(v))
		}
		result = append(result, row)
	}
	return result, nil
}
