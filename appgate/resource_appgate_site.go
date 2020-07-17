package appgate

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/appgate/terraform-provider-appgate/client/v12/openapi"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/schema"
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

			"site_id": {
				Type:        schema.TypeString,
				Description: "ID of the object.",
				Computed:    true,
			},

			"name": {
				Type:        schema.TypeString,
				Description: "Name of the object.",
				Required:    true,
			},

			"notes": {
				Type:        schema.TypeString,
				Description: "Notes for the object. Used for documentation purposes.",
				Default:     DefaultDescription,
				Optional:    true,
			},

			"tags": {
				Type:        schema.TypeSet,
				Description: "Array of tags.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

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
				Computed: true,
				// Default:  true,
			},

			"vpn": {
				// Due to the limitation of tf-11115 it is not possible to nest maps.
				// https://github.com/hashicorp/terraform/issues/11115
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
							Type:     schema.TypeMap,
							Optional: true,
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
							Type:     schema.TypeMap,
							Optional: true,
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
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": {
										Type:     schema.TypeString,
										Required: true,
									},
									"ipv6": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"web_proxy_enabled": {
							Type:        schema.TypeBool,
							Description: "Flag for manipulating web proxy p12 file. Setting this false will delete the existing p12 file from database.",
							Optional:    true,
						},
						"web_proxy_key_store": {
							Type:        schema.TypeString,
							Description: "The PKCS12 package to be used for web proxy. The file must be with no password and must include the full certificate chain and a private key. In Base64 format.",
							Optional:    true,
						},
						"ip_access_log_interval_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
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
							Type:     schema.TypeList,
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
							Type:     schema.TypeList,
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
							Type:     schema.TypeList,
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
					},
				},
			},
		},
	}
}

func resourceAppgateSiteCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating Site: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.SitesApi

	args := openapi.NewSiteWithDefaults()
	args.Id = uuid.New().String()
	args.SetName(d.Get("name").(string))
	args.SetShortName(d.Get("short_name").(string))
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
		vpn, err := readSiteVPNFromConfig(v.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetVpn(vpn)
	}

	if v, ok := d.GetOk("name_resolution"); ok {
		nameResolution, err := readSiteNameResolutionFromConfig(v.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetNameResolution(nameResolution)
	}

	request := api.SitesPost(context.Background())
	request = request.Site(*args)
	site, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create site %+v", prettyPrintAPIError(err))
	}

	d.SetId(site.Id)

	return resourceAppgateSiteRead(d, meta)
}

func resourceAppgateSiteRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Site Name: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.SitesApi
	request := api.SitesIdGet(context.Background(), d.Id())
	site, _, err := request.Authorization(token).Execute()
	if err != nil {
		// TODO check if 404
		d.SetId("")
		return fmt.Errorf("Failed to read Site, %+v", err)
	}

	d.SetId(site.Id)
	d.Set("site_id", site.Id)
	d.Set("name", site.Name)
	d.Set("notes", site.Notes)
	d.Set("tags", site.Tags)
	d.Set("network_subnets", site.NetworkSubnets)
	if site.IpPoolMappings != nil {
		if err = d.Set("ip_pool_mappings", flattenSiteIPpoolmappning(*site.IpPoolMappings)); err != nil {
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
		if err = d.Set("vpn", flattenSiteVPN(*site.Vpn)); err != nil {
			return err
		}
	}
	if site.NameResolution != nil {
		if err = d.Set("name_resolution", flattenNameResolution(*site.NameResolution)); err != nil {
			return err
		}
	}
	return nil
}

func flattenSiteIPpoolmappning(in []openapi.SiteAllOfIpPoolMappings) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["from"] = v.From
		m["to"] = v.To

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

func flattenSiteVPN(in openapi.SiteAllOfVpn) []interface{} {
	m := make(map[string]interface{})

	if v, o := in.GetStateSharingOk(); o != false {
		m["state_sharing"] = v
	}
	if v, o := in.GetSnatOk(); o != false {
		m["snat"] = v
	}
	if in.HasTls() {
		tls := make(map[string]interface{})
		tls["enabled"] = strconv.FormatBool(in.Tls.GetEnabled())
		m["tls"] = tls
	}
	if in.HasDtls() {
		dtls := make(map[string]interface{})
		dtls["enabled"] = strconv.FormatBool(in.Dtls.GetEnabled())
		m["dtls"] = dtls
	}
	if in.HasRouteVia() && in.RouteVia.Ipv4 != nil {
		routeVia := make(map[string]interface{})
		routeVia["ipv4"] = in.RouteVia.GetIpv4()
		routeVia["ipv6"] = in.RouteVia.GetIpv6()
		m["route_via"] = routeVia
	}

	if v, o := in.GetWebProxyEnabledOk(); o != false {
		m["web_proxy_enabled"] = v
	}
	if v, o := in.GetWebProxyKeyStoreOk(); o != false {
		m["web_proxy_key_store"] = v
	}
	m["ip_access_log_interval_seconds"] = in.IpAccessLogIntervalSeconds
	return []interface{}{m}
}

func flattenNameResolution(in openapi.SiteAllOfNameResolution) []interface{} {
	m := make(map[string]interface{})
	if v, o := in.GetUseHostsFileOk(); o != false {
		m["use_hosts_file"] = v
	}
	if v, o := in.GetDnsResolversOk(); o != false {
		m["dns_resolvers"] = flattenSiteDNSResolver(*v)
	}
	if v, o := in.GetAwsResolversOk(); o != false {
		m["aws_resolvers"] = flattenSiteAWSResolver(*v)
	}
	if v, o := in.GetAzureResolversOk(); o != false {
		m["azure_resolvers"] = flattenSiteAzureResolver(*v)
	}
	if v, o := in.GetEsxResolversOk(); o != false {
		m["esx_resolvers"] = flattenSiteESXResolvers(*v)
	}
	if v, o := in.GetGcpResolversOk(); o != false {
		m["gcp_resolvers"] = flattenSiteGCPResolvers(*v)
	}
	return []interface{}{m}
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

func flattenSiteESXResolvers(in []openapi.SiteAllOfNameResolutionEsxResolvers) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["name"] = v.GetName()
		m["update_interval"] = v.GetUpdateInterval()
		m["hostname"] = v.GetHostname()
		m["username"] = v.GetUsername()
		m["password"] = v.GetPassword()

		out[i] = m
	}
	return out
}

func flattenSiteAzureResolver(in []openapi.SiteAllOfNameResolutionAzureResolvers) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["name"] = v.GetName()
		m["update_interval"] = v.GetUpdateInterval()
		m["subscription_id"] = v.GetSubscriptionId()
		m["tenant_id"] = v.GetTenantId()
		m["client_id"] = v.GetClientId()
		m["secret"] = v.GetSecret()

		out[i] = m
	}
	return out
}

func flattenSiteAWSResolver(in []openapi.SiteAllOfNameResolutionAwsResolvers) []map[string]interface{} {
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
		m["secret_access_key"] = v.GetSecretAccessKey()
		m["https_proxy"] = v.GetHttpsProxy()
		m["resolve_with_master_credentials"] = v.GetResolveWithMasterCredentials()
		if vv, o := v.GetAssumedRolesOk(); o != false {
			m["assumed_roles"] = flattenSiteAwsAssumedRoles(*vv)
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
		m["servers"] = v.GetServers()
		m["search_domains"] = v.GetSearchDomains()

		out[i] = m
	}
	return out
}

func resourceAppgateSiteUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Site: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.SitesApi
	request := api.SitesIdGet(context.Background(), d.Id())

	orginalSite, _, err := request.Authorization(token).Execute()
	if err != nil {
		// TODO check if 404
		d.SetId("")
		return fmt.Errorf("Failed to read Site, %+v", err)
	}

	if d.HasChange("name") {
		orginalSite.SetName(d.Get("name").(string))
	}

	if d.HasChange("tags") {
		orginalSite.SetTags(schemaExtractTags(d))
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
		_, n := d.GetChange("vpn")
		vpn, err := readSiteVPNFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		orginalSite.SetVpn(vpn)
	}

	if d.HasChange("name_resolution") {
		_, n := d.GetChange("name_resolution")
		nameResolution, err := readSiteNameResolutionFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		orginalSite.SetNameResolution(nameResolution)
	}

	putRequest := api.SitesIdPut(context.Background(), d.Id())
	_, _, err = putRequest.Site(orginalSite).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update site %+v", prettyPrintAPIError(err))
	}
	return resourceAppgateSiteRead(d, meta)
}

func resourceAppgateSiteDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete Site: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.SitesApi
	request := api.SitesIdDelete(context.Background(), d.Id())
	_, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to delete Site, %+v", err)
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
				return result, fmt.Errorf("Failed to resolve default gateway excluded subnets: %+v", err)
			}
			result.SetExcludedSubnets(nets)
		}

	}
	return result, nil
}

func readSiteVPNFromConfig(vpns []interface{}) (openapi.SiteAllOfVpn, error) {
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
			rawTLS := v.(map[string]interface{})

			if v, ok := rawTLS["enabled"]; ok {
				tls.SetEnabled(v.(bool))
			}
			result.SetTls(*tls)
		}

		if v, ok := raw["dtls"]; ok {
			dtls := openapi.NewSiteAllOfVpnDtlsWithDefaults()
			rawDTLS := v.(map[string]interface{})

			if v, ok := rawDTLS["enabled"]; ok {
				dtls.SetEnabled(v.(bool))
			}
			result.SetDtls(*dtls)
		}

		if v, ok := raw["route_via"]; ok {
			routeVia := openapi.NewSiteAllOfVpnRouteViaWithDefaults()
			rawRouteVia := v.(map[string]interface{})

			if v, ok := rawRouteVia["ipv4"]; ok {
				routeVia.SetIpv4(v.(string))
			}
			if v, ok := rawRouteVia["ipv6"]; ok {
				routeVia.SetIpv6(v.(string))
			}
			result.SetRouteVia(*routeVia)
		}

		if v, ok := raw["web_proxy_enabled"]; ok {
			result.SetWebProxyEnabled(v.(bool))
		}
		if v, ok := raw["web_proxy_key_store"]; ok {
			result.SetWebProxyKeyStore(v.(string))
		}
		if v, ok := raw["ip_access_log_interval_seconds"]; ok {
			result.SetIpAccessLogIntervalSeconds(float32(v.(int)))
		}
	}
	return result, nil
}

func readSiteNameResolutionFromConfig(nameresolutions []interface{}) (openapi.SiteAllOfNameResolution, error) {
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
			dnsResolvers, err := readDNSResolversFromConfig(v.([]interface{}))
			if err != nil {
				return result, err
			}
			result.SetDnsResolvers(dnsResolvers)
		}
		if v, ok := raw["aws_resolvers"]; ok {
			awsResolvers, err := readAWSResolversFromConfig(v.([]interface{}))
			if err != nil {
				return result, err
			}
			result.SetAwsResolvers(awsResolvers)
		}
		if v, ok := raw["azure_resolvers"]; ok {
			azureResolvers, err := readAzureResolversFromConfig(v.([]interface{}))
			if err != nil {
				return result, err
			}
			result.SetAzureResolvers(azureResolvers)
		}
		if v, ok := raw["esx_resolvers"]; ok {
			esxResolvers, err := readESXResolversFromConfig(v.([]interface{}))
			if err != nil {
				return result, err
			}
			result.SetEsxResolvers(esxResolvers)
		}
		if v, ok := raw["gcp_resolvers"]; ok {
			gcpResolvers, err := readGCPResolversFromConfig(v.([]interface{}))
			if err != nil {
				return result, err
			}
			result.SetGcpResolvers(gcpResolvers)
		}
	}
	return result, nil
}

func readDNSResolversFromConfig(dnsConfigs []interface{}) ([]openapi.SiteAllOfNameResolutionDnsResolvers, error) {
	result := make([]openapi.SiteAllOfNameResolutionDnsResolvers, 0)
	for _, dns := range dnsConfigs {
		raw := dns.(map[string]interface{})
		row := openapi.SiteAllOfNameResolutionDnsResolvers{}
		log.Printf("[DEBUG] readDNSResolversFromConfig RAW IS: %+v", raw)
		if v, ok := raw["name"]; ok {
			row.SetName(v.(string))
		}
		if v, ok := raw["update_interval"]; ok {
			row.SetUpdateInterval(int32(v.(int)))
		}
		if v := raw["servers"]; len(v.([]interface{})) > 0 {
			servers, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, fmt.Errorf("Failed to resolve dns serers: %+v", err)
			}
			if len(servers) > 0 {
				row.SetServers(servers)
			}
		}
		if v := raw["search_domains"]; len(v.([]interface{})) > 0 {
			domains, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, fmt.Errorf("Failed to resolve dns search domains: %+v", err)
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
		if v, ok := raw["update_interval"]; ok {
			row.SetUpdateInterval(int32(v.(int)))
		}
		if v := raw["vpcs"]; len(v.([]interface{})) > 0 {
			vpcs, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, fmt.Errorf("Failed to resolve vpcs from aws config: %+v", err)
			}
			row.SetVpcs(vpcs)
		}
		if v, ok := raw["vpc_auto_discovery"]; ok && v.(bool) {
			row.SetVpcAutoDiscovery(v.(bool))
		}
		if v := raw["regions"]; len(v.([]interface{})) > 0 {
			regions, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, fmt.Errorf("Failed to resolve regions from aws config: %+v", err)
			}
			row.SetRegions(regions)
		}
		if v, ok := raw["use_iam_role"]; ok && v.(bool) {
			row.SetUseIAMRole(v.(bool))
		}
		if v, ok := raw["access_key_id"]; ok {
			row.SetAccessKeyId(v.(string))
		}
		if v, ok := raw["secret_access_key"]; ok {
			row.SetSecretAccessKey(v.(string))
		}
		if v, ok := raw["https_proxy"]; ok && len(v.(string)) > 0 {
			row.SetHttpsProxy(v.(string))
		}
		if v, ok := raw["resolve_with_master_credentials"]; ok && v.(bool) {
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

func readAzureResolversFromConfig(azureConfigs []interface{}) ([]openapi.SiteAllOfNameResolutionAzureResolvers, error) {
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
