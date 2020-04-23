package appgate

import (
	"context"
	"fmt"
	"log"

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

			"created": {
				Type:        schema.TypeString,
				Description: "Create date.",
				Computed:    true,
			},

			"updated": {
				Type:        schema.TypeString,
				Description: "Create date.",
				Computed:    true,
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
				Type:     schema.TypeSet,
				Optional: true,
				// ConfigMode: schema.SchemaConfigModeAttr,
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
					},
				},
			}, // vpn

			"name_resolution": {
				Type:       schema.TypeSet,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"use_hosts_file": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"dns_resolvers": {
							Type:       schema.TypeSet,
							Optional:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"update_interval": {
										Type:     schema.TypeString,
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
							Type:       schema.TypeSet,
							Optional:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"update_interval": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"vpcs": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"vpc_auto_discovery": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"regions": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"use_iam_role": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"access_key_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"secret_access_key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"https_proxy": {
										Type:     schema.TypeString,
										Required: true,
									},
									"resolve_with_master_credentials": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"assumed_roles": {
										Type:       schema.TypeSet,
										Required:   true,
										ConfigMode: schema.SchemaConfigModeAttr,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{

												"account_id": {
													Type:     schema.TypeString,
													Required: true,
												},

												"role_name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"external_id": {
													Type:     schema.TypeString,
													Required: true,
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
							Type:       schema.TypeSet,
							Optional:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"update_interval": {
										Type:     schema.TypeString,
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
									"secret_id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},

						"esx_resolvers": {
							Type:       schema.TypeSet,
							Optional:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"update_interval": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"hostname": {
										Type:     schema.TypeString,
										Required: true,
									},
									"username": {
										Type:     schema.TypeString,
										Required: true,
									},
									"password": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"gcp_resolvers": {
							Type:       schema.TypeSet,
							Optional:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"update_interval": {
										Type:     schema.TypeString,
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
	rawsubnets := d.Get("network_subnets").(*schema.Set).List()
	subnets := make([]string, 0)
	for _, raw := range rawsubnets {
		subnets = append(subnets, raw.(string))
	}
	var defaultGateway []openapi.SiteAllOfDefaultGateway
	if g, ok := d.GetOk("default_gateway"); ok {
		gw := g.(*schema.Set).List()
		for _, r := range gw {
			l := r.(map[string]interface{})
			gwo := openapi.SiteAllOfDefaultGateway{
				EnabledV4: openapi.PtrBool(l["enabled_v4"].(bool)),
				EnabledV6: openapi.PtrBool(l["enabled_v6"].(bool)),
			}
			excludedSubnets := make([]string, 0)
			for _, t := range l["excluded_subnets"].([]interface{}) {
				excludedSubnets = append(excludedSubnets, t.(string))
			}
			gwo.ExcludedSubnets = &excludedSubnets
			defaultGateway = append(defaultGateway, gwo)
		}
	}
	args := openapi.NewSiteWithDefaults()
	args.SetName(d.Get("name").(string))
	args.SetId(uuid.New().String())
	args.SetNotes(d.Get("notes").(string))
	args.SetTags(schemaExtractTags(d))
	args.SetEntitlementBasedRouting(d.Get("entitlement_based_routing").(bool))
	args.SetNetworkSubnets(subnets)
	args.SetShortName(d.Get("short_name").(string))

	if len(defaultGateway) > 0 {
		args.DefaultGateway = &defaultGateway[0]
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
	gw := make(map[string]interface{})

	gw["enabled_v4"] = site.DefaultGateway.EnabledV4
	gw["enabled_v6"] = site.DefaultGateway.EnabledV6
	exsub := make([]interface{}, 0, 0)
	for _, sub := range site.DefaultGateway.GetExcludedSubnets() {
		exsub = append(exsub, sub)
	}
	gw["excluded_subnets"] = exsub

	d.SetId(site.Id)
	d.Set("site_id", site.Id)
	d.Set("notes", site.Notes)
	d.Set("created", site.Created.String())
	d.Set("updated", site.Updated.String())
	d.Set("tags", site.Tags)
	d.Set("network_subnets", site.NetworkSubnets)
	d.Set("short_name", site.ShortName)
	d.Set("entitlement_based_routing", site.EntitlementBasedRouting)

	if err := d.Set("default_gateway", []interface{}{gw}); err != nil {
		return fmt.Errorf("Failed to read default gateway on %s: %+v", d.Id(), err)
	}
	return nil
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

	rawsubnets := d.Get("network_subnets").(*schema.Set).List()
	subnets := make([]string, 0)
	for _, raw := range rawsubnets {
		subnets = append(subnets, raw.(string))
	}

	orginalSite.SetName(d.Get("name").(string))
	orginalSite.SetNotes(d.Get("notes").(string))
	orginalSite.SetShortName(d.Get("short_name").(string))
	orginalSite.SetEntitlementBasedRouting(d.Get("entitlement_based_routing").(bool))
	orginalSite.SetNetworkSubnets(subnets)
	orginalSite.SetTags(schemaExtractTags(d))

	putRequest := api.SitesIdPut(context.Background(), d.Id())
	_, _, err = putRequest.Site(orginalSite).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to update Site, %+v", err)
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
