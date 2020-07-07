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
				Type:     schema.TypeSet,
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
				Type:     schema.TypeSet,
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

			"snmp_server": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"tcp_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"udp_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"snmpd_conf": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"allow_sources": reUsableSchemas["allow_sources"],
					},
				},
			},

			"healthcheck_server": {
				Type:     schema.TypeSet,
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
						},

						"allow_sources": reUsableSchemas["allow_sources"],
					},
				},
			},

			"prometheus_exporter": {
				Type:     schema.TypeSet,
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
						},

						"allow_sources": reUsableSchemas["allow_sources"],
					},
				},
			},

			"ping": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_sources": reUsableSchemas["allow_sources"],
					},
				},
			},

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
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"iot_connector"},
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
			"log_forwarder": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"elasticsearch": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"aws_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"aws_secret": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"aws_region": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"use_instance_credentials": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"retention_days": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},

						"tcp_clients": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"host": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"format": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
											s := v.(string)
											enums := []string{"json", "syslog"}
											if inArray(s, enums) {
												return
											}
											errs = append(errs, fmt.Errorf(
												"%s: is invalid option, expected %+v", name, enums,
											))
											return
										},
									},
									"use_tls": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},

						"sites": {
							Type:        schema.TypeSet,
							Description: "Array of sites.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"iot_connector": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"gateway"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"clients": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"device_id": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"sources": reUsableSchemas["allow_sources"],

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

			"rsyslog_destinations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"selector": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"template": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"destination": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"hostname_aliases": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
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
	args.SetCustomization(d.Get("customization").(string))

	if c, ok := d.GetOk("client_interface"); ok {
		cinterface, err := readClientInterfaceFromConfig(c.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetClientInterface(cinterface)
	}

	if p, ok := d.GetOk("peer_interface"); ok {
		pinterface, err := readPeerInterfaceFromConfig(p.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetPeerInterface(pinterface)
	}

	if a, ok := d.GetOk("admin_interface"); ok {
		ainterface, err := readAdminInterfaceFromConfig(a.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetAdminInterface(ainterface)
	}

	if n, ok := d.GetOk("networking"); ok {
		network, err := readNetworkingFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetNetworking(network)
	}

	if n, ok := d.GetOk("ssh_server"); ok {
		sshServer, err := readSSHServerFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetSshServer(sshServer)
	}

	if n, ok := d.GetOk("snmp_server"); ok {
		srv, err := readSNMPServerFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetSnmpServer(srv)
	}

	if n, ok := d.GetOk("healthcheck_server"); ok {
		srv, err := readHealthcheckServerFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetHealthcheckServer(srv)
	}

	if n, ok := d.GetOk("prometheus_exporter"); ok {
		exporter, err := readPrometheusExporterFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetPrometheusExporter(exporter)
	}

	if n, ok := d.GetOk("ping"); ok {
		p, err := readPingFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetPing(p)
	}

	if n, ok := d.GetOk("ntp"); ok {
		ntp, err := readNTPFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetNtp(ntp)
	}

	if v, ok := d.GetOk("log_server"); ok {
		logSrv, err := readLogServerFromConfig(v.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetLogServer(logSrv)
	}

	if v, ok := d.GetOk("controller"); ok {
		ctrl, err := readControllerFromConfig(v.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetController(ctrl)
	}

	if v, ok := d.GetOk("gateway"); ok {
		gw, err := readGatewayFromConfig(v.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetGateway(gw)
	}

	if v, ok := d.GetOk("log_forwarder"); ok {
		lf, err := readLogForwardFromConfig(v.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetLogForwarder(lf)
	}

	if v, ok := d.GetOk("iot_connector"); ok {
		iot, err := readIotConnectorFromConfig(v.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetIotConnector(iot)
	}

	if v, ok := d.GetOk("rsyslog_destinations"); ok {
		rsyslog, err := readRsyslogDestinationFromConfig(v.([]interface{}))
		if err != nil {
			return err
		}
		args.SetRsyslogDestinations(rsyslog)
	}

	if v, ok := d.GetOk("hostname_aliases"); ok {
		hostnames, err := readHostnameAliasesFromConfig(v.([]interface{}))
		if err != nil {
			return err
		}
		args.SetHostnameAliases(hostnames)
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

		if v := raw["ipv6"].(*schema.Set); v.Len() > 0 {
			ipv6networking := openapi.ApplianceAllOfNetworkingIpv6{}
			for _, v := range v.List() {
				ipv4Data := v.(map[string]interface{})
				if v := ipv4Data["dhcp"].(*schema.Set); v.Len() > 0 {
					for _, v := range v.List() {
						ipv6networking.SetDhcp(readNetworkIpv6DhcpFromConfig(v.(map[string]interface{})))
					}
				}
				// TODO do static.
			}
			nic.SetIpv6(ipv6networking)
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
		ipv4dhcp.SetRouters(v.(bool))
	}
	if v, ok := ipv4raw["ntp"]; ok {
		ipv4dhcp.SetNtp(v.(bool))
	}
	return ipv4dhcp
}

func readNetworkIpv6DhcpFromConfig(ipv6raw map[string]interface{}) openapi.ApplianceAllOfNetworkingIpv6Dhcp {
	ipv6dhcp := openapi.ApplianceAllOfNetworkingIpv6Dhcp{}
	if v, ok := ipv6raw["enabled"]; ok {
		ipv6dhcp.SetEnabled(v.(bool))
	}
	if v, ok := ipv6raw["dns"]; ok {
		ipv6dhcp.SetDns(v.(bool))
	}
	if v, ok := ipv6raw["ntp"]; ok {
		ipv6dhcp.SetNtp(v.(bool))
	}
	return ipv6dhcp
}

func readNtpServersFromConfig(input []interface{}) ([]openapi.ApplianceAllOfNtpServers, error) {
	var r []openapi.ApplianceAllOfNtpServers
	for _, s := range input {
		raw := s.(map[string]interface{})
		row := openapi.ApplianceAllOfNtpServers{}
		if v, ok := raw["hostname"]; ok {
			row.SetHostname(v.(string))
		}
		if v, ok := raw["key_type"]; ok && len(v.(string)) > 0 {
			row.SetKeyType(v.(string))
		}
		if v, ok := raw["key"]; ok && len(v.(string)) > 0 {
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

	// d.Set("client_interface", saveClientInterface)
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

	if d.HasChange("name") {
		originalAppliance.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalAppliance.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalAppliance.SetTags(schemaExtractTags(d))
	}

	if d.HasChange("hostname") {
		originalAppliance.SetHostname(d.Get("hostname").(string))
	}

	if d.HasChange("site") {
		originalAppliance.SetSite(d.Get("site").(string))
	}

	if d.HasChange("customization") {
		originalAppliance.SetCustomization(d.Get("customization").(string))
	}

	if d.HasChange("client_interface") {
		_, n := d.GetChange("client_interface")
		cinterface, err := readClientInterfaceFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		originalAppliance.SetClientInterface(cinterface)
	}

	if d.HasChange("peer_interface") {
		_, n := d.GetChange("peer_interface")
		pinterface, err := readPeerInterfaceFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		originalAppliance.SetPeerInterface(pinterface)
	}

	if d.HasChange("admin_interface") {
		_, a := d.GetChange("admin_interface")
		ainterface, err := readAdminInterfaceFromConfig(a.(*schema.Set).List())
		if err != nil {
			return err
		}
		originalAppliance.SetAdminInterface(ainterface)
	}

	if d.HasChange("networking") {
		_, n := d.GetChange("networking")
		networking, err := readNetworkingFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		originalAppliance.SetNetworking(networking)
	}

	if d.HasChange("ntp") {
		_, n := d.GetChange("ntp")
		ntp, err := readNTPFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		originalAppliance.SetNtp(ntp)
	}

	if d.HasChange("ssh_server") {
		_, n := d.GetChange("ssh_server")
		sshServer, err := readSSHServerFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		originalAppliance.SetSshServer(sshServer)
	}

	if d.HasChange("snmp_server") {
		_, n := d.GetChange("ssh_server")
		srv, err := readSNMPServerFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		originalAppliance.SetSnmpServer(srv)
	}

	if d.HasChange("healthcheck_server") {
		_, n := d.GetChange("healthcheck_server")
		srv, err := readHealthcheckServerFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		originalAppliance.SetHealthcheckServer(srv)
	}

	if d.HasChange("prometheus_exporter") {
		_, n := d.GetChange("prometheus_exporter")
		exporter, err := readPrometheusExporterFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		originalAppliance.SetPrometheusExporter(exporter)
	}

	if d.HasChange("ping") {
		_, n := d.GetChange("ping")
		p, err := readPingFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		originalAppliance.SetPing(p)
	}

	if d.HasChange("log_server") {
		_, n := d.GetChange("log_server")
		logSrv, err := readLogServerFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		originalAppliance.SetLogServer(logSrv)
	}

	if d.HasChange("controller") {
		_, n := d.GetChange("controller")
		ctrl, err := readControllerFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		originalAppliance.SetController(ctrl)
	}

	if d.HasChange("gateway") {
		_, n := d.GetChange("gateway")
		gw, err := readGatewayFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		originalAppliance.SetGateway(gw)
	}

	if d.HasChange("log_forwarder") {
		_, n := d.GetChange("log_forwarder")
		lf, err := readLogForwardFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		originalAppliance.SetLogForwarder(lf)
	}

	if d.HasChange("iot_connector") {
		_, n := d.GetChange("iot_connector")
		iot, err := readIotConnectorFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		originalAppliance.SetIotConnector(iot)
	}

	if d.HasChange("rsyslog_destinations") {
		_, n := d.GetChange("rsyslog_destinations")
		rsys, err := readRsyslogDestinationFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		originalAppliance.SetRsyslogDestinations(rsys)
	}

	if d.HasChange("hostname_aliases") {
		_, n := d.GetChange("hostname_aliases")
		hostnames, err := readHostnameAliasesFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		originalAppliance.SetHostnameAliases(hostnames)
	}

	req := api.AppliancesIdPut(ctx, d.Id())

	_, _, err = req.Appliance(originalAppliance).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update appliance %+v", prettyPrintAPIError(err))
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

func readClientInterfaceFromConfig(cinterfaces []interface{}) (openapi.ApplianceAllOfClientInterface, error) {
	cinterface := openapi.ApplianceAllOfClientInterface{}
	for _, r := range cinterfaces {
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
				return cinterface, fmt.Errorf("Failed to resolve network hosts: %+v", err)
			}
			cinterface.SetAllowSources(allowSources)
		}
		if v, ok := raw["override_spa_mode"]; ok {
			cinterface.SetOverrideSpaMode(v.(string))
		}
	}
	return cinterface, nil
}

func readPeerInterfaceFromConfig(pinterfaces []interface{}) (openapi.ApplianceAllOfPeerInterface, error) {
	pinterf := openapi.ApplianceAllOfPeerInterface{}
	for _, r := range pinterfaces {
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
				return pinterf, fmt.Errorf("Failed to resolve network hosts: %+v", err)
			}
			pinterf.SetAllowSources(allowSources)
		}
	}
	return pinterf, nil
}
func readAdminInterfaceFromConfig(adminInterfaces []interface{}) (openapi.ApplianceAllOfAdminInterface, error) {
	aInterface := openapi.ApplianceAllOfAdminInterface{}
	for _, admin := range adminInterfaces {
		if admin == nil {
			continue
		}

		raw := admin.(map[string]interface{})
		if v, ok := raw["hostname"]; ok {
			aInterface.SetHostname(v.(string))
		}
		if v, ok := raw["https_port"]; ok {
			aInterface.SetHttpsPort(int32(v.(int)))
		}
		if v := raw["https_ciphers"].(*schema.Set); v.Len() > 0 {
			ciphers := make([]string, 0)
			for _, c := range v.List() {
				ciphers = append(ciphers, c.(string))
			}
			aInterface.SetHttpsCiphers(ciphers)
		}
		if v := raw["allow_sources"]; len(v.([]interface{})) > 0 {
			allowSources, err := readAllowSourcesFromConfig(v.([]interface{}))
			if err != nil {
				return aInterface, fmt.Errorf("Failed to admin interface allowed sources: %+v", err)
			}
			aInterface.SetAllowSources(allowSources)
		}
	}
	return aInterface, nil
}

func readNetworkingFromConfig(networks []interface{}) (openapi.ApplianceAllOfNetworking, error) {
	network := openapi.ApplianceAllOfNetworking{}
	for _, netw := range networks {
		if netw == nil {
			continue
		}

		rawNetwork := netw.(map[string]interface{})
		if v := rawNetwork["hosts"]; len(v.([]interface{})) > 0 {
			hosts, err := readNetworkHostFromConfig(v.([]interface{}))
			if err != nil {
				return network, fmt.Errorf("Failed to resolve network hosts: %+v", err)
			}
			network.SetHosts(hosts)
		}

		if v := rawNetwork["nics"]; len(v.([]interface{})) > 0 {
			nics, err := readNetworkNicsFromConfig(v.([]interface{}))
			if err != nil {
				return network, fmt.Errorf("Failed to resolve network nics: %+v", err)
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
				return network, fmt.Errorf("Failed to resolve network routes: %+v", err)
			}
			network.SetRoutes(routes)
		}
	}
	return network, nil
}

func readSSHServerFromConfig(sshServers []interface{}) (openapi.ApplianceAllOfSshServer, error) {
	sshServer := openapi.ApplianceAllOfSshServer{}
	for _, srv := range sshServers {
		if srv == nil {
			continue
		}
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
				return sshServer, fmt.Errorf("Failed to resolve ssh server allowed sources: %+v", err)
			}
			sshServer.SetAllowSources(allowSources)
		}
	}
	return sshServer, nil
}

func readSNMPServerFromConfig(snmpServers []interface{}) (openapi.ApplianceAllOfSnmpServer, error) {
	server := openapi.ApplianceAllOfSnmpServer{}
	for _, srv := range snmpServers {
		if srv == nil {
			continue
		}

		rawServer := srv.(map[string]interface{})
		if v, ok := rawServer["enabled"]; ok {
			server.SetEnabled(v.(bool))
		}
		if v, ok := rawServer["tcp_port"]; ok {
			server.SetTcpPort(int32(v.(int)))
		}
		if v, ok := rawServer["udp_port"]; ok {
			server.SetUdpPort(int32(v.(int)))
		}
		if v, ok := rawServer["snmpd_conf"]; ok {
			server.SetSnmpdConf(v.(string))
		}
		if v := rawServer["allow_sources"]; len(v.([]interface{})) > 0 {
			allowSources, err := readAllowSourcesFromConfig(v.([]interface{}))
			if err != nil {
				return server, fmt.Errorf("Failed to resolve network hosts: %+v", err)
			}
			server.SetAllowSources(allowSources)
		}
	}
	return server, nil
}

func readHealthcheckServerFromConfig(healhCheckServers []interface{}) (openapi.ApplianceAllOfHealthcheckServer, error) {
	server := openapi.ApplianceAllOfHealthcheckServer{}
	for _, srv := range healhCheckServers {
		if srv == nil {
			continue
		}

		rawServer := srv.(map[string]interface{})
		if v, ok := rawServer["enabled"]; ok {
			server.SetEnabled(v.(bool))
		}
		if v, ok := rawServer["port"]; ok {
			server.SetPort(int32(v.(int)))
		}
		if v := rawServer["allow_sources"]; len(v.([]interface{})) > 0 {
			allowSources, err := readAllowSourcesFromConfig(v.([]interface{}))
			if err != nil {
				return server, fmt.Errorf("Failed to resolve network hosts: %+v", err)
			}
			server.SetAllowSources(allowSources)
		}
	}
	return server, nil
}

func readNTPFromConfig(ntps []interface{}) (openapi.ApplianceAllOfNtp, error) {
	ntpCfg := openapi.ApplianceAllOfNtp{}
	for _, ntp := range ntps {
		if ntp == nil {
			continue
		}
		raw := ntp.(map[string]interface{})
		if servers := raw["servers"]; len(servers.([]interface{})) > 0 {
			ntpServers, err := readNtpServersFromConfig(servers.([]interface{}))
			if err != nil {
				return ntpCfg, fmt.Errorf("Failed to resolve ntp servers: %+v", err)
			}
			ntpCfg.SetServers(ntpServers)
		}
	}
	return ntpCfg, nil
}

func readLogServerFromConfig(logServers []interface{}) (openapi.ApplianceAllOfLogServer, error) {
	srv := openapi.ApplianceAllOfLogServer{}
	for _, ctrl := range logServers {
		r := ctrl.(map[string]interface{})
		val := openapi.ApplianceAllOfLogServer{}
		if v, ok := r["enabled"]; ok {
			val.SetEnabled(v.(bool))
		}
		if v, ok := r["retention_days"]; ok {
			val.SetRetentionDays(int32(v.(int)))
		}
	}
	return srv, nil
}

func readControllerFromConfig(controllers []interface{}) (openapi.ApplianceAllOfController, error) {
	val := openapi.ApplianceAllOfController{}
	for _, ctrl := range controllers {
		r := ctrl.(map[string]interface{})
		if v, ok := r["enabled"]; ok {
			val.SetEnabled(v.(bool))
		}
	}
	return val, nil
}

func readGatewayFromConfig(gateways []interface{}) (openapi.ApplianceAllOfGateway, error) {
	val := openapi.ApplianceAllOfGateway{}
	for _, ctrl := range gateways {
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
	return val, nil
}

func readPrometheusExporterFromConfig(exporters []interface{}) (openapi.ApplianceAllOfPrometheusExporter, error) {
	val := openapi.ApplianceAllOfPrometheusExporter{}
	for _, srv := range exporters {
		if srv == nil {
			continue
		}

		rawServer := srv.(map[string]interface{})
		if v, ok := rawServer["enabled"]; ok {
			val.SetEnabled(v.(bool))
		}
		if v, ok := rawServer["port"]; ok {
			val.SetPort(int32(v.(int)))
		}

		if v := rawServer["allow_sources"]; len(v.([]interface{})) > 0 {
			allowSources, err := readAllowSourcesFromConfig(v.([]interface{}))
			if err != nil {
				return val, err
			}
			val.SetAllowSources(allowSources)
		}
	}
	return val, nil
}

func readPingFromConfig(pingers []interface{}) (openapi.ApplianceAllOfPing, error) {
	val := openapi.ApplianceAllOfPing{}
	for _, srv := range pingers {
		if srv == nil {
			continue
		}

		rawServer := srv.(map[string]interface{})

		if v := rawServer["allow_sources"]; len(v.([]interface{})) > 0 {
			allowSources, err := readAllowSourcesFromConfig(v.([]interface{}))
			if err != nil {
				return val, err
			}
			val.SetAllowSources(allowSources)
		}
	}
	return val, nil
}

func readLogForwardFromConfig(logforwards []interface{}) (openapi.ApplianceAllOfLogForwarder, error) {
	val := openapi.ApplianceAllOfLogForwarder{}
	for _, logforward := range logforwards {
		if logforward == nil {
			continue
		}

		raw := logforward.(map[string]interface{})

		if v, ok := raw["enabled"]; ok {
			val.SetEnabled(v.(bool))
		}

		if v := raw["elasticsearch"].(*schema.Set); v.Len() > 0 {
			elasticsearch := openapi.ApplianceAllOfLogForwarderElasticsearch{}
			for _, s := range v.List() {
				r := s.(map[string]interface{})
				if v, ok := r["url"]; ok {
					elasticsearch.SetUrl(v.(string))
				}
				if v, ok := r["aws_id"]; ok {
					elasticsearch.SetAwsId(v.(string))
				}
				if v, ok := r["aws_secret"]; ok {
					elasticsearch.SetAwsSecret(v.(string))
				}
				if v, ok := r["aws_region"]; ok {
					elasticsearch.SetAwsRegion(v.(string))
				}
				if v, ok := r["use_instance_credentials"]; ok {
					elasticsearch.SetUseInstanceCredentials(v.(bool))
				}
				if v, ok := r["retention_days"]; ok {
					elasticsearch.SetRetentionDays(int32(v.(int)))
				}
			}
			val.SetElasticsearch(elasticsearch)
		}

		if v := raw["tcp_clients"]; len(v.([]interface{})) > 0 {
			tcpClients := make([]openapi.TcpClient, 0)
			for _, s := range v.([]interface{}) {
				tcpClient := openapi.TcpClient{}
				r := s.(map[string]interface{})
				if v, ok := r["name"]; ok {
					tcpClient.SetName(v.(string))
				}
				if v, ok := r["host"]; ok {
					tcpClient.SetHost(v.(string))
				}
				if v, ok := r["port"]; ok {
					tcpClient.SetPort(int32(v.(int)))
				}
				if v, ok := r["format"]; ok {
					tcpClient.SetFormat(v.(string))
				}
				if v, ok := r["use_tls"]; ok {
					tcpClient.SetUseTLS(v.(bool))
				}
				tcpClients = append(tcpClients, tcpClient)
			}
			val.SetTcpClients(tcpClients)
		}
		sites := make([]string, 0)
		if v := raw["sites"].(*schema.Set); v.Len() > 0 {
			for _, s := range v.List() {
				site := s.(string)
				if len(site) > 0 {
					sites = append(sites, site)
				}
			}
		}
		val.SetSites(sites)
	}
	return val, nil
}
func readIotConnectorFromConfig(iots []interface{}) (openapi.ApplianceAllOfIotConnector, error) {
	val := openapi.ApplianceAllOfIotConnector{}
	for _, iot := range iots {
		if iot == nil {
			continue
		}

		raw := iot.(map[string]interface{})

		if v, ok := raw["enabled"]; ok {
			val.SetEnabled(v.(bool))
		}
		if v := raw["clients"]; len(v.([]interface{})) > 0 {
			clients := make([]openapi.ApplianceAllOfIotConnectorClients, 0)
			for _, c := range v.([]interface{}) {
				client := openapi.ApplianceAllOfIotConnectorClients{}
				r := c.(map[string]interface{})
				if v, ok := r["name"]; ok {
					client.SetName(v.(string))
				}
				if v, ok := r["device_id"]; ok {
					client.SetDeviceId(v.(string))
				}
				// allowed sources
				if v := r["sources"]; len(v.([]interface{})) > 0 {
					sources, err := readAllowSourcesFromConfig(v.([]interface{}))
					if err != nil {
						return val, err
					}
					client.SetSources(sources)
				}
				if v, ok := r["snat"]; ok {
					client.SetSnat(v.(bool))
				}
				clients = append(clients, client)
			}
			val.SetClients(clients)
		}
	}
	return val, nil
}

func readRsyslogDestinationFromConfig(rsyslogs []interface{}) ([]openapi.ApplianceAllOfRsyslogDestinations, error) {
	result := make([]openapi.ApplianceAllOfRsyslogDestinations, 0)
	for _, rsys := range rsyslogs {
		if rsys == nil {
			continue
		}
		r := openapi.ApplianceAllOfRsyslogDestinations{}
		raw := rsys.(map[string]interface{})
		if v, ok := raw["selector"]; ok {
			r.SetSelector(v.(string))
		}
		if v, ok := raw["template"]; ok {
			r.SetTemplate(v.(string))
		}
		if v, ok := raw["destination"]; ok {
			r.SetDestination(v.(string))
		}
		result = append(result, r)
	}
	return result, nil
}

func readHostnameAliasesFromConfig(hostnameAliases []interface{}) ([]string, error) {
	result := make([]string, 0)
	for _, hostname := range hostnameAliases {
		if hostname == nil {
			continue
		}
		result = append(result, hostname.(string))
	}
	return result, nil
}
