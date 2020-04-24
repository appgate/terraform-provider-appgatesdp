package appgate

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/appgate/terraform-provider-appgate/client/v12/openapi"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAppgateAppliance() *schema.Resource {
	reUsableSchemas := make(map[string]*schema.Schema)

	reUsableSchemas["allow_sources"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"address": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validateIPaddress,
				},
				"netmask": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"nic": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
	return &schema.Resource{
		Create: resourceAppgateApplianceCreate,
		Read:   resourceAppgateApplianceRead,
		Update: resourceAppgateApplianceUpdate,
		Delete: resourceAppgateApplianceDelete,
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

			"appliance_id": {
				Type:        schema.TypeString,
				Description: "ID of the object.",
				Computed:    true,
			},
			"seed_file": {
				Type:        schema.TypeString,
				Description: "Seed file (json) generated from appliance used in remote-exec.",
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

			"hostname": {
				Type:        schema.TypeString,
				Description: "Name of the object.",
				Required:    true,
			},

			"site": {
				Type:        schema.TypeString,
				Description: "Site served by the Appliance. Entitlements on this Site will be included in the Entitlement Token for this Appliance. Not useful if Gateway role is not enabled.",
				Optional:    true,
			},

			"customization": {
				Type:        schema.TypeString,
				Description: "Customization assigned to this Appliance.",
				Optional:    true,
			},

			"connect_to_peers_using_client_port_with_spa": {
				Type:        schema.TypeBool,
				Description: "Makes the Appliance to connect to Controller/LogServer/LogForwarders using their clientInterface.httpsPort instead of peerInterface.httpsPort. The Appliance uses SPA to connect.",
				Optional:    true,
			},

			"client_interface": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"proxy_protocol": {
							Type:        schema.TypeBool,
							Description: "To enable/disable Proxy protocol on this Appliance.",
							Optional:    true,
							Default:     false,
						},

						"hostname": {
							Type:     schema.TypeString,
							Required: true,
						},

						"https_port": {
							Type:     schema.TypeInt,
							Default:  443,
							Optional: true,
						},
						"dtls_port": {
							Type:     schema.TypeInt,
							Default:  443,
							Optional: true,
						},

						"allow_sources": reUsableSchemas["allow_sources"],

						"override_spa_mode": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
								s := v.(string)
								list := []string{"Disabled", "TCP", "UDP-TCP"}
								for _, x := range list {
									if s == x {
										return
									}
								}
								errs = append(errs, fmt.Errorf("override_spa_mode must be on of %v, got %s", list, s))
								return
							},
						},
					},
				},
			},

			"peer_interface": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:     schema.TypeString,
							Required: true,
						},
						"https_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"allow_sources": reUsableSchemas["allow_sources"],
					},
				},
			},

			"admin_interface": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:     schema.TypeString,
							Required: true,
						},
						"https_port": {
							Type:     schema.TypeInt,
							Default:  443,
							Optional: true,
						},
						"https_ciphers": {
							Type:        schema.TypeSet,
							Description: "The type of TLS ciphers to allow. See: https://www.openssl.org/docs/man1.0.2/apps/ciphers.html for all supported ciphers.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"allow_sources": reUsableSchemas["allow_sources"],
					},
				},
			},

			"networking": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"hosts": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": {
										Type:     schema.TypeString,
										Required: true,
									},
									"address": {
										Type:         schema.TypeString,
										ValidateFunc: validateIPaddress,
										Required:     true,
									},
								},
							},
						},

						"nics": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"ipv4": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"dhcp": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"enabled": {
																Type:     schema.TypeBool,
																Optional: true,
															},
															"dns": {
																Type:     schema.TypeBool,
																Optional: true,
															},
															"routers": {
																Type:     schema.TypeBool,
																Optional: true,
															},
															"ntp": {
																Type:     schema.TypeBool,
																Optional: true,
															},
														},
													},
												},
												"static": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"address": {
																Type:         schema.TypeString,
																ValidateFunc: validateIPaddress,
																Optional:     true,
															},
															"netmask": {
																Type:     schema.TypeInt,
																Optional: true,
															},
															"hostname": {
																Type:     schema.TypeInt,
																Optional: true,
															},
															"snat": {
																Type:     schema.TypeBool,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},
									"ipv6": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"dhcp": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"enabled": {
																Type:     schema.TypeBool,
																Optional: true,
															},
															"dns": {
																Type:     schema.TypeBool,
																Optional: true,
															},
															"routers": {
																Type:     schema.TypeBool,
																Optional: true,
															},
															"ntp": {
																Type:     schema.TypeBool,
																Optional: true,
															},
														},
													},
												},
												"static": {
													Type:     schema.TypeSet,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"address": {
																Type:         schema.TypeString,
																ValidateFunc: validateIPaddress,
																Required:     true,
															},
															"netmask": {
																Type:     schema.TypeInt,
																Required: true,
															},
															"hostname": {
																Type:     schema.TypeInt,
																Required: true,
															},
															"snat": {
																Type:     schema.TypeBool,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},

						"dns_servers": {
							Type:        schema.TypeSet,
							Description: "DNS Server addresses.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},

						"dns_domains": {
							Type:        schema.TypeSet,
							Description: "DNS Search domains.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"routes": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateIPaddress,
									},
									"netmask": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"gateway": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"nic": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			// ntp_server Deprecated as of 4.3.0, use 'ntp' field instead. NTP servers to synchronize time.
			// "ntp_servers": {
			// 	Type:        schema.TypeSet,
			// 	Description: "Array of tags.",
			// 	Optional:    true,
			// 	Elem:        &schema.Schema{Type: schema.TypeString},
			// },
			"ntp": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"servers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"hostname": {
										Type:     schema.TypeString,
										Required: true,
									},

									"key_type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
											s := v.(string)
											enums := []string{"MD5", "SHA", "SHA1", "SHA256", "SHA512", "RMD160"}
											if inArray(s, enums) {
												return
											}
											errs = append(errs, fmt.Errorf(
												"%s: is invalid option, expected %+v", name, enums,
											))
											return
										},
									},

									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"ssh_server": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  22,
						},

						"allow_sources": reUsableSchemas["allow_sources"],

						"password_authentication": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			// TODO
			// "snmp_server": {},
			// "healthcheck_server": {},
			// "prometheus_exporter": {},
			// "ping": {},
			"log_server": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"retention_days": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  30,
						},
					},
				},
			},
			"controller": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"gateway": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"vpn": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"weight": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  100,
									},
									"allow_destinations": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:         schema.TypeString,
													ValidateFunc: validateIPaddress,
													Optional:     true,
												},
												"netmask": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"nic": {
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
				},
			},
			// "log_forwarder": {},
			// "iot_connector": {},
			// "rsyslog_destinations": {},
			// "hostname_aliases": {},
		},
	}
}

func resourceAppgateApplianceCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating Appliance with name: %s", d.Get("name").(string))
	ctx := context.Background()
	token := meta.(*Client).Token
	api := meta.(*Client).API.AppliancesApi

	args := openapi.NewApplianceWithDefaults()
	args.Id = uuid.New().String()
	args.SetName(d.Get("name").(string))
	args.SetHostname(d.Get("hostname").(string))
	args.SetNotes(d.Get("notes").(string))
	args.SetTags(schemaExtractTags(d))
	args.SetSite(d.Get("site").(string))

	var ci []openapi.ApplianceAllOfClientInterface
	if c, ok := d.GetOk("client_interface"); ok {
		cinterfaces := c.(*schema.Set).List()
		for _, r := range cinterfaces {
			cinterface := openapi.ApplianceAllOfClientInterface{}
			raw := r.(map[string]interface{})
			if v, ok := raw["hostname"]; ok {
				cinterface.SetHostname(v.(string))
			}
			if v, ok := raw["https_port"]; ok {
				cinterface.SetHttpsPort(int32(v.(int)))
			}
			if v, ok := raw["dtls_port"]; ok {
				cinterface.SetDtlsPort(int32(v.(int)))
			}
			if v := raw["allow_sources"]; len(v.([]interface{})) > 0 {
				allowSources, err := readAllowSourcesFromConfig(v.([]interface{}))
				if err != nil {
					return fmt.Errorf("Failed to resolve network hosts: %+v", err)
				}
				cinterface.SetAllowSources(allowSources)
			}
			if v, ok := raw["override_spa_mode"]; ok {
				cinterface.SetOverrideSpaMode(v.(string))
			}
			ci = append(ci, cinterface)
		}
	}
	if len(ci) > 0 {
		args.ClientInterface = ci[0]
	}
	var pi []openapi.ApplianceAllOfPeerInterface
	if p, ok := d.GetOk("peer_interface"); ok {
		pinterfaces := p.(*schema.Set).List()
		for _, r := range pinterfaces {
			pinterf := openapi.ApplianceAllOfPeerInterface{}
			raw := r.(map[string]interface{})
			if v, ok := raw["hostname"]; ok {
				pinterf.SetHostname(v.(string))
			}
			if v, ok := raw["https_port"]; ok {
				pinterf.SetHttpsPort(int32(v.(int)))
			}
			if v := raw["allow_sources"]; len(v.([]interface{})) > 0 {
				allowSources, err := readAllowSourcesFromConfig(v.([]interface{}))
				if err != nil {
					return fmt.Errorf("Failed to resolve network hosts: %+v", err)
				}
				pinterf.SetAllowSources(allowSources)
			}
			pi = append(pi, pinterf)
		}
	}
	if len(pi) > 0 {
		args.PeerInterface = pi[0]
	}

	var networkings []openapi.ApplianceAllOfNetworking
	if n, ok := d.GetOk("networking"); ok {
		networks := n.([]interface{})
		for _, netw := range networks {
			if netw == nil {
				continue
			}
			network := openapi.ApplianceAllOfNetworking{}
			rawNetwork := netw.(map[string]interface{})
			if v := rawNetwork["hosts"]; len(v.([]interface{})) > 0 {
				hosts, err := readNetworkHostFromConfig(v.([]interface{}))
				if err != nil {
					return fmt.Errorf("Failed to resolve network hosts: %+v", err)
				}
				network.SetHosts(hosts)
			}

			if v := rawNetwork["nics"]; len(v.([]interface{})) > 0 {
				nics, err := readNetworkNicsFromConfig(v.([]interface{}))
				if err != nil {
					return fmt.Errorf("Failed to resolve network nics: %+v", err)
				}
				network.SetNics(nics)
			}
			if v := rawNetwork["dns_servers"].(*schema.Set); v.Len() > 0 {
				d := make([]string, 0)
				for _, dns := range v.List() {
					d = append(d, dns.(string))
				}
				network.SetDnsServers(d)
			}
			if v := rawNetwork["dns_domains"].(*schema.Set); v.Len() > 0 {
				d := make([]string, 0)
				for _, dns := range v.List() {
					d = append(d, dns.(string))
				}
				network.SetDnsDomains(d)
			}
			if v := rawNetwork["routes"]; len(v.([]interface{})) > 0 {
				routes, err := readNetworkRoutesFromConfig(v.([]interface{}))
				if err != nil {
					return fmt.Errorf("Failed to resolve network routes: %+v", err)
				}
				network.SetRoutes(routes)
			}
			networkings = append(networkings, network)
		}
	}
	if len(networkings) > 0 {
		args.Networking = networkings[0]
	}

	var sshServers []openapi.ApplianceAllOfSshServer
	if s, ok := d.GetOk("ssh_server"); ok {
		servers := s.([]interface{})
		for _, srv := range servers {
			if srv == nil {
				continue
			}
			sshServer := openapi.ApplianceAllOfSshServer{}
			rawServer := srv.(map[string]interface{})
			if v, ok := rawServer["enabled"]; ok {
				sshServer.SetEnabled(v.(bool))
			}
			if v, ok := rawServer["port"]; ok {
				sshServer.SetPort(int32(v.(int)))
			}
			if v, ok := rawServer["password_authentication"]; ok {
				sshServer.SetPasswordAuthentication(v.(bool))
			}
			if v := rawServer["allow_sources"]; len(v.([]interface{})) > 0 {
				allowSources, err := readAllowSourcesFromConfig(v.([]interface{}))
				if err != nil {
					return fmt.Errorf("Failed to resolve network hosts: %+v", err)
				}
				sshServer.SetAllowSources(allowSources)
			}
			sshServers = append(sshServers, sshServer)
		}
	}

	if len(sshServers) > 0 {
		args.SetSshServer(sshServers[0])
	}

	var ntpList []openapi.ApplianceAllOfNtp
	if n, ok := d.GetOk("ntp"); ok {
		for _, ntp := range n.(*schema.Set).List() {
			if ntp == nil {
				continue
			}
			ntpCfg := openapi.ApplianceAllOfNtp{}
			raw := ntp.(map[string]interface{})
			if servers := raw["servers"]; len(servers.([]interface{})) > 0 {
				ntpServers, err := readNtpServersFromConfig(servers.([]interface{}))
				if err != nil {
					return fmt.Errorf("Failed to resolve ntp servers: %+v", err)
				}
				ntpCfg.SetServers(ntpServers)
			}
			ntpList = append(ntpList, ntpCfg)
		}
	}
	if len(ntpList) > 0 {
		args.SetNtp(ntpList[0])
	}

	if v, ok := d.GetOk("controller"); ok {
		for _, ctrl := range v.(*schema.Set).List() {
			r := ctrl.(map[string]interface{})
			if v, ok := r["enabled"]; ok {
				val := openapi.ApplianceAllOfController{}
				val.SetEnabled(v.(bool))
				args.SetController(val)
			}
		}
	}
	if v, ok := d.GetOk("log_server"); ok {
		for _, ctrl := range v.(*schema.Set).List() {
			r := ctrl.(map[string]interface{})
			val := openapi.ApplianceAllOfLogServer{}
			if v, ok := r["enabled"]; ok {
				val.SetEnabled(v.(bool))
			}
			if v, ok := r["retention_days"]; ok {
				val.SetRetentionDays(int32(v.(int)))
			}
			args.SetLogServer(val)
		}
	}
	if v, ok := d.GetOk("gateway"); ok {
		val := openapi.ApplianceAllOfGateway{}
		for _, ctrl := range v.(*schema.Set).List() {
			r := ctrl.(map[string]interface{})
			if v, ok := r["enabled"]; ok {
				val.SetEnabled(v.(bool))
			}
			if v := r["vpn"].(*schema.Set); v.Len() > 0 {
				vpn := openapi.ApplianceAllOfGatewayVpn{}
				for _, s := range v.List() {
					raw := s.(map[string]interface{})
					if v, ok := raw["weight"]; ok {
						vpn.SetWeight(int32(v.(int)))
					}
					if v := raw["allow_destinations"]; len(v.([]interface{})) > 0 {
						rawAllowedDestinations := v.([]interface{})
						allowDestinations := make([]openapi.ApplianceAllOfGatewayVpnAllowDestinations, 0)
						for _, r := range rawAllowedDestinations {
							raw := r.(map[string]interface{})
							ad := openapi.ApplianceAllOfGatewayVpnAllowDestinations{}
							if v := raw["address"].(string); v != "" {
								ad.SetAddress(v)
							}
							if v, ok := raw["netmask"]; ok {
								ad.SetNetmask(int32(v.(int)))
							}
							if v := raw["nic"].(string); v != "" {
								ad.SetNic(v)
							}
							allowDestinations = append(allowDestinations, ad)
						}
						vpn.SetAllowDestinations(allowDestinations)
					}
				}
				val.SetVpn(vpn)
			}
		}
		args.SetGateway(val)
	}
	log.Printf("\n appliance arguments: \n %+v \n", args)
	request := api.AppliancesPost(ctx)
	request = request.Appliance(*args)
	appliance, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create appliance %+v", prettyPrintAPIError(err))
	}

	d.SetId(appliance.Id)

	return resourceAppgateApplianceRead(d, meta)
}

func readNetworkNicsFromConfig(hosts []interface{}) ([]openapi.ApplianceAllOfNetworkingNics, error) {
	apiNics := make([]openapi.ApplianceAllOfNetworkingNics, 0)
	for _, h := range hosts {
		raw := h.(map[string]interface{})
		nic := openapi.ApplianceAllOfNetworkingNics{}
		if v, ok := raw["enabled"]; ok {
			nic.Enabled = openapi.PtrBool(v.(bool))
		}
		if v := raw["name"].(string); v != "" {
			nic.Name = v
		}
		if v := raw["ipv4"].(*schema.Set); v.Len() > 0 {
			ipv4networking := openapi.ApplianceAllOfNetworkingIpv4{}
			for _, v := range v.List() {
				ipv4Data := v.(map[string]interface{})
				if v := ipv4Data["dhcp"].(*schema.Set); v.Len() > 0 {
					for _, v := range v.List() {
						ipv4networking.SetDhcp(readNetworkIpv4DhcpFromConfig(v.(map[string]interface{})))
					}
				}
				// TODO do static.
			}
			nic.SetIpv4(ipv4networking)
		}
		apiNics = append(apiNics, nic)
	}
	return apiNics, nil
}

func readNetworkIpv4DhcpFromConfig(ipv4raw map[string]interface{}) openapi.ApplianceAllOfNetworkingIpv4Dhcp {
	ipv4dhcp := openapi.ApplianceAllOfNetworkingIpv4Dhcp{}
	if v, ok := ipv4raw["enabled"]; ok {
		ipv4dhcp.SetEnabled(v.(bool))
	}
	if v, ok := ipv4raw["dns"]; ok {
		ipv4dhcp.SetDns(v.(bool))
	}
	if v, ok := ipv4raw["routers"]; ok {
		ipv4dhcp.SetDns(v.(bool))
	}
	if v, ok := ipv4raw["ntp"]; ok {
		ipv4dhcp.SetDns(v.(bool))
	}
	return ipv4dhcp
}

func readNtpServersFromConfig(input []interface{}) ([]openapi.ApplianceAllOfNtpServers, error) {
	var r []openapi.ApplianceAllOfNtpServers
	for _, s := range input {
		raw := s.(map[string]interface{})
		row := openapi.ApplianceAllOfNtpServers{}
		if v, ok := raw["hostname"]; ok {
			row.SetHostname(v.(string))
		}
		if v, ok := raw["key_type"]; ok {
			row.SetKeyType(v.(string))
		}
		if v, ok := raw["key"]; ok {
			row.SetKey(v.(string))
		}
		r = append(r, row)
	}
	return r, nil
}

func readAllowSourcesFromConfig(input []interface{}) ([]map[string]interface{}, error) {
	r := make([]map[string]interface{}, 0)
	for _, s := range input {
		raw := s.(map[string]interface{})
		row := make(map[string]interface{}, 0)
		// TODO, address can be both list and single string.
		if v, ok := raw["address"]; ok {
			row["address"] = v.(string)
		}
		if v, ok := raw["netmask"]; ok {
			row["netmask"] = v.(int)
		}
		if v := raw["nic"].(string); v != "" {
			row["nic"] = v
		}
		r = append(r, row)
	}
	return r, nil
}

func readNetworkRoutesFromConfig(routes []interface{}) ([]openapi.ApplianceAllOfNetworkingRoutes, error) {
	apiRoutes := make([]openapi.ApplianceAllOfNetworkingRoutes, 0)
	for _, r := range routes {
		raw := r.(map[string]interface{})
		route := openapi.ApplianceAllOfNetworkingRoutes{}
		if v := raw["address"].(string); v != "" {
			route.SetAddress(v)
		}
		if v, ok := raw["netmask"]; ok {
			route.SetNetmask(int32(v.(int)))
		}
		if v := raw["gateway"].(string); v != "" {
			route.SetGateway(v)
		}
		if v := raw["nic"].(string); v != "" {
			route.SetNic(v)
		}
		apiRoutes = append(apiRoutes, route)
	}
	return apiRoutes, nil
}

func readNetworkHostFromConfig(hosts []interface{}) ([]openapi.ApplianceAllOfNetworkingHosts, error) {
	apiHosts := make([]openapi.ApplianceAllOfNetworkingHosts, 0)
	for _, h := range hosts {
		raw := h.(map[string]interface{})
		host := openapi.ApplianceAllOfNetworkingHosts{}
		if v := raw["hostname"].(string); v != "" {
			host.Hostname = v
		}
		if v := raw["address"].(string); v != "" {
			host.Address = v
		}
		apiHosts = append(apiHosts, host)
	}
	return apiHosts, nil
}

func resourceAppgateApplianceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Appliance Name: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.AppliancesApi
	ctx := context.Background()
	request := api.AppliancesIdGet(ctx, d.Id())
	appliance, _, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Failed to read Appliance, %+v", err)
	}
	d.Set("appliance_id", appliance.Id)
	d.Set("name", appliance.Name)
	d.Set("tags", appliance.Tags)
	d.Set("appliance_id", appliance.Id)

	localClientInterface := d.Get("client_interface").(*schema.Set).List()
	var saveClientInterface []map[string]interface{}
	for _, raw := range localClientInterface {
		l := raw.(map[string]interface{})
		o := make(map[string]interface{}, 0)
		log.Printf("[DEBUG] raw client interface: %+v", l)
		if v, ok := l["hostname"]; ok {
			o["hostname"] = v
		}
		if v, ok := l["dtls_port"]; ok {
			o["dtls_port"] = v.(int)
		}
		if v, ok := l["https_port"]; ok {
			o["https_port"] = v.(int)
		}
		if v, ok := l["proxy_protocol"]; ok {
			o["proxy_protocol"] = v.(bool)
		}
		if v, ok := l["allow_sources"]; ok {
			os := make(map[string]interface{}, 0)
			for _, y := range v.([]interface{}) {
				z := y.(map[string]interface{})
				if v, ok := z["address"]; ok {
					os["address"] = v
				}
				if v, ok := z["nic"]; ok {
					os["nic"] = v
				}
				if v, ok := z["netmask"]; ok {
					os["netmask"] = v
				}

			}
			o["allow_sources"] = os
		}
		saveClientInterface = append(saveClientInterface, o)
	}
	d.Set("client_interface", saveClientInterface)
	// d.Set("networking", appliance.GetNetworking())
	d.Set("notes", appliance.GetNotes())
	d.Set("hostname", appliance.GetHostname())
	// d.Set("peer_interface", appliance.PeerInterface)

	if ok, _ := appliance.GetActivatedOk(); *ok {
		d.Set("seed_file", "")
		return nil
	}
	exportRequest := api.AppliancesIdExportPost(ctx, appliance.Id)
	exportRequest = exportRequest.SSHConfig(openapi.SSHConfig{
		Password: openapi.PtrString("cz"),
	})
	seedmap, _, err := exportRequest.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not export appliance %+v", prettyPrintAPIError(err))
	}
	encodedSeed, err := json.Marshal(seedmap)
	if err != nil {
		return fmt.Errorf("Could not parse json seed file: %+v", err)
	}
	d.Set("seed_file", b64.StdEncoding.EncodeToString([]byte(encodedSeed)))

	return nil
}

func resourceAppgateApplianceUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Appliance: %s", d.Get("name").(string))
	ctx := context.Background()
	token := meta.(*Client).Token
	api := meta.(*Client).API.AppliancesApi
	request := api.AppliancesIdGet(ctx, d.Id())
	originalAppliance, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read Appliance, %+v", err)
	}
	// TODO check fields..
	originalAppliance.SetName(d.Get("name").(string))

	req := api.AppliancesIdPut(ctx, d.Id())

	_, _, err = req.Appliance(originalAppliance).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to update Appliance, %+v", err)
	}
	return resourceAppgateApplianceRead(d, meta)
}

func resourceAppgateApplianceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete Appliance: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.AppliancesApi
	ctx := context.Background()

	// Get appliance
	request := api.AppliancesIdGet(ctx, d.Id())
	appliance, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to delete Appliance while GET, %+v", err)
	}
	// Deactive
	if ok, _ := appliance.GetActivatedOk(); *ok {
		log.Printf("[DEBUG] Appliance is active, deactive and wiping before deleting")
		deactiveRequest := api.AppliancesIdDeactivatePost(ctx, appliance.GetId())
		_, err = deactiveRequest.Wipe(true).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("Failed to delete Appliance while deactiving, %+v", err)
		}
	}

	// Delete
	deleteRequest := api.AppliancesIdDelete(ctx, appliance.GetId())
	_, err = deleteRequest.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to delete Appliance, %+v", err)
	}
	d.SetId("")
	return nil
}
