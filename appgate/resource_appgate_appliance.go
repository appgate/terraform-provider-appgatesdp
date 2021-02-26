package appgate

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v14/openapi"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"activated": {
				Type:     schema.TypeBool,
				Computed: true,
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
				Computed:    true,
			},

			"client_interface": {
				Type:     schema.TypeList,
				MaxItems: 1,
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
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
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
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
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
							Type:        schema.TypeList,
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
				MaxItems: 1,
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
							Computed: true,
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
									"mtu": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"ipv4": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"dhcp": {
													Type:     schema.TypeSet,
													Optional: true,
													Computed: true,
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
													Type:     schema.TypeList,
													Optional: true,

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
																Type:     schema.TypeString,
																Optional: true,
															},
															"snat": {
																Type:     schema.TypeBool,
																Optional: true,
															},
														},
													},
												},
												"virtual_ip": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"ipv6": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"dhcp": {
													Type:     schema.TypeSet,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"enabled": {
																Type:     schema.TypeBool,
																Optional: true,
																Computed: true,
															},
															"dns": {
																Type:     schema.TypeBool,
																Optional: true,
																Computed: true,
															},
															"ntp": {
																Type:     schema.TypeBool,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
												"static": {
													Type:     schema.TypeList,
													Optional: true,
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
																Type:     schema.TypeString,
																Required: true,
															},
															"snat": {
																Type:     schema.TypeBool,
																Optional: true,
															},
														},
													},
												},
												"virtual_ip": {
													Type:     schema.TypeString,
													Optional: true,
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
			"ntp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
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
				MaxItems: 1,
				Optional: true,
				Computed: true,
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
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
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
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
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
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
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
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_sources": reUsableSchemas["allow_sources"],
					},
				},
			},

			"log_server": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
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
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
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
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"connector"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
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
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"elasticsearch": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
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
										Required: true,
									},
									"host": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"format": {
										Type:     schema.TypeString,
										Required: true,
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
									"filter": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"aws_kineses": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
									"type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
											s := v.(string)
											enums := []string{"Stream", "Firehose"}
											if inArray(s, enums) {
												return
											}
											errs = append(errs, fmt.Errorf(
												"%s: is invalid option, expected %+v", name, enums,
											))
											return
										},
									},
									"stream_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"batch_size": {
										Type:     schema.TypeInt,
										Computed: true,
										Optional: true,
									},
									"number_of_partition_keys": {
										Type:     schema.TypeInt,
										Computed: true,
										Optional: true,
									},
									"filter": {
										Type:     schema.TypeString,
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

			"connector": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"gateway"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"express_clients": {
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

									"allow_resources": {
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
											},
										},
									},

									"snat_to_resources": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"advanced_clients": {
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

									"allow_resources": reUsableSchemas["allow_sources"],

									"snat_to_tunnel": {
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

	args.SetTags(schemaExtractTags(d))

	if v, ok := d.GetOk("notes"); ok {
		args.SetNotes(v.(string))
	}

	if v, ok := d.GetOk("site"); ok {
		args.SetSite(v.(string))
	}

	if v, ok := d.GetOk("customization"); ok {
		args.SetCustomization(v.(string))
	}

	if v, ok := d.GetOk("connect_to_peers_using_client_port_with_spa"); ok {
		args.SetConnectToPeersUsingClientPortWithSpa(v.(bool))
	}

	if c, ok := d.GetOk("client_interface"); ok {
		cinterface, err := readClientInterfaceFromConfig(c.([]interface{}))
		if err != nil {
			return err
		}
		args.SetClientInterface(cinterface)
	}

	if p, ok := d.GetOk("peer_interface"); ok {
		pinterface, err := readPeerInterfaceFromConfig(p.([]interface{}))
		if err != nil {
			return err
		}
		args.SetPeerInterface(pinterface)
	}

	if a, ok := d.GetOk("admin_interface"); ok {
		ainterface, err := readAdminInterfaceFromConfig(a.([]interface{}))
		if err != nil {
			return err
		}
		args.SetAdminInterface(ainterface)
	}

	if n, ok := d.GetOk("networking"); ok {
		network, err := readNetworkingFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		args.SetNetworking(network)
	}

	if n, ok := d.GetOk("ssh_server"); ok {
		sshServer, err := readSSHServerFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		args.SetSshServer(sshServer)
	}

	if n, ok := d.GetOk("snmp_server"); ok {
		srv, err := readSNMPServerFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		args.SetSnmpServer(srv)
	}

	if n, ok := d.GetOk("healthcheck_server"); ok {
		srv, err := readHealthcheckServerFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		args.SetHealthcheckServer(srv)
	}

	if n, ok := d.GetOk("prometheus_exporter"); ok {
		exporter, err := readPrometheusExporterFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		args.SetPrometheusExporter(exporter)
	}

	if n, ok := d.GetOk("ping"); ok {
		p, err := readPingFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		args.SetPing(p)
	}

	if n, ok := d.GetOk("ntp"); ok {
		ntp, err := readNTPFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		args.SetNtp(ntp)
	}

	if v, ok := d.GetOk("log_server"); ok {
		logSrv, err := readLogServerFromConfig(v.([]interface{}))
		if err != nil {
			return err
		}
		args.SetLogServer(logSrv)
	}

	if v, ok := d.GetOk("controller"); ok {
		ctrl, err := readControllerFromConfig(v.([]interface{}))
		if err != nil {
			return err
		}
		args.SetController(ctrl)
	}

	if v, ok := d.GetOk("gateway"); ok {
		gw, err := readGatewayFromConfig(v.([]interface{}))
		if err != nil {
			return err
		}
		args.SetGateway(gw)
	}

	if v, ok := d.GetOk("log_forwarder"); ok {
		lf, err := readLogForwardFromConfig(v.([]interface{}))
		if err != nil {
			return err
		}
		args.SetLogForwarder(lf)
	}

	if v, ok := d.GetOk("connector"); ok {
		connector, err := readApplianceConnectorFromConfig(v.([]interface{}))
		if err != nil {
			return err
		}
		args.SetConnector(connector)
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
		if v, ok := raw["mtu"]; ok {
			mtu := openapi.PtrInt32(int32(v.(int)))
			if *mtu > int32(64) {
				nic.SetMtu(*mtu)
			}
		}

		if v := raw["ipv4"].([]interface{}); len(v) > 0 {
			ipv4networking := openapi.ApplianceAllOfNetworkingIpv4{}
			for _, item := range v {
				ipv4Data := item.(map[string]interface{})
				if item := ipv4Data["dhcp"].(*schema.Set); item.Len() > 0 {
					for _, i := range item.List() {
						ipv4networking.SetDhcp(readNetworkIpv4DhcpFromConfig(i.(map[string]interface{})))
					}
				}
				if item := ipv4Data["static"]; len(item.([]interface{})) > 0 {
					ipv4networking.SetStatic(readNetworkIpv4StaticFromConfig(item.([]interface{})))
				}
				if v, o := ipv4Data["virtual_ip"]; o && len(v.(string)) > 0 {
					ipv4networking.SetVirtualIp(v.(string))
				}
			}
			nic.SetIpv4(ipv4networking)
		}

		if v := raw["ipv6"].([]interface{}); len(v) > 0 {
			ipv6networking := openapi.ApplianceAllOfNetworkingIpv6{}
			for _, item := range v {
				ipv6Data := item.(map[string]interface{})
				if item := ipv6Data["dhcp"].(*schema.Set); item.Len() > 0 {
					for _, i := range item.List() {
						ipv6networking.SetDhcp(readNetworkIpv6DhcpFromConfig(i.(map[string]interface{})))
					}
				}
				if v := ipv6Data["static"]; len(v.([]interface{})) > 0 {
					ipv6networking.SetStatic(readNetworkIpv6StaticFromConfig(v.([]interface{})))
				}
				if v, o := ipv6Data["virtual_ip"]; o && len(v.(string)) > 0 {
					ipv6networking.SetVirtualIp(v.(string))
				}
			}
			nic.SetIpv6(ipv6networking)
		}
		apiNics = append(apiNics, nic)
	}
	return apiNics, nil
}

func readNetworkIpv4StaticFromConfig(ipv4staticraw []interface{}) []openapi.ApplianceAllOfNetworkingIpv4Static {
	var r []openapi.ApplianceAllOfNetworkingIpv4Static
	for _, s := range ipv4staticraw {
		raw := s.(map[string]interface{})
		row := openapi.ApplianceAllOfNetworkingIpv4Static{}
		if v, ok := raw["address"]; ok {
			row.SetAddress(v.(string))
		}
		if v, ok := raw["netmask"]; ok {
			row.SetNetmask(int32(v.(int)))
		}
		if v, ok := raw["hostname"]; ok {
			row.SetHostname(v.(string))
		}
		if v, ok := raw["snat"]; ok {
			row.SetSnat(v.(bool))
		}
		r = append(r, row)
	}
	return r
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

func readNetworkIpv6StaticFromConfig(ipv6staticraw []interface{}) []openapi.ApplianceAllOfNetworkingIpv6Static {
	var r []openapi.ApplianceAllOfNetworkingIpv6Static
	for _, s := range ipv6staticraw {
		raw := s.(map[string]interface{})
		row := openapi.ApplianceAllOfNetworkingIpv6Static{}
		if v, ok := raw["address"]; ok {
			row.SetAddress(v.(string))
		}
		if v, ok := raw["netmask"]; ok {
			row.SetNetmask(int32(v.(int)))
		}
		if v, ok := raw["hostname"]; ok {
			row.SetHostname(v.(string))
		}
		if v, ok := raw["snat"]; ok {
			row.SetSnat(v.(bool))
		}
		r = append(r, row)
	}
	return r
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

func readAllowSourcesFromConfig(input []map[string]interface{}) ([]map[string]interface{}, error) {
	r := make([]map[string]interface{}, 0)
	for _, raw := range input {
		row := make(map[string]interface{}, 0)
		// TODO, address can be both list and single string.
		if v, ok := raw["address"]; ok {
			row["address"] = v.(string)
		}
		if v, ok := raw["netmask"]; ok {
			row["netmask"] = v
		}
		if v, ok := raw["nic"]; ok && v != "" {
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
	d.Set("notes", appliance.GetNotes())
	d.Set("hostname", appliance.GetHostname())

	if v, ok := d.GetOkExists("site"); ok {
		d.Set("site", v)
	}

	if v, ok := d.GetOkExists("customization"); ok {
		d.Set("customization", v)
	}

	d.Set("connect_to_peers_using_client_port_with_spa", appliance.GetConnectToPeersUsingClientPortWithSpa())

	if v, o := appliance.GetClientInterfaceOk(); o != false {
		ci, err := flattenApplianceClientInterface(*v)
		if err != nil {
			return err
		}
		d.Set("client_interface", ci)
	}

	if v, o := appliance.GetPeerInterfaceOk(); o != false {
		peerInterface, err := flattenAppliancePeerInterface(*v)
		if err != nil {
			return err
		}
		d.Set("peer_interface", peerInterface)
	}

	if v, o := appliance.GetAdminInterfaceOk(); o != false {
		adminInterface, err := flattenApplianceAdminInterface(*v)
		if err != nil {
			return err
		}
		if err := d.Set("admin_interface", adminInterface); err != nil {
			return err
		}
	}

	if v, o := appliance.GetNetworkingOk(); o != false {
		networking, err := flattenApplianceNetworking(*v)
		if err != nil {
			return err
		}
		if err := d.Set("networking", networking); err != nil {
			return err
		}
	}

	if _, o := d.GetOkExists("ntp"); o {
		v := appliance.GetNtp()
		ntp := make(map[string]interface{})
		servers := make([]map[string]interface{}, 0)
		for _, n := range v.GetServers() {
			srv := make(map[string]interface{})
			srv["hostname"] = n.GetHostname()
			srv["key_type"] = n.GetKeyType()
			srv["key"] = n.GetKey()
			servers = append(servers, srv)
		}
		ntp["servers"] = servers
		if err := d.Set("ntp", []interface{}{ntp}); err != nil {
			return err
		}
	}

	if _, o := d.GetOkExists("ssh_server"); o {
		v := appliance.GetSshServer()
		sshServer := make(map[string]interface{})
		sshServer["enabled"] = v.GetEnabled()
		sshServer["port"] = v.GetPort()
		sshServer["allow_sources"] = v.GetAllowSources()
		sshServer["password_authentication"] = v.GetPasswordAuthentication()

		if err := d.Set("ssh_server", []interface{}{sshServer}); err != nil {
			return err
		}
	}

	if _, o := d.GetOkExists("snmp_server"); o {
		v := appliance.GetSnmpServer()
		snmpSrv := make(map[string]interface{})
		snmpSrv["enabled"] = v.GetEnabled()
		snmpSrv["tcp_port"] = v.GetTcpPort()
		snmpSrv["udp_port"] = v.GetUdpPort()
		snmpSrv["snmpd_conf"] = v.GetSnmpdConf()
		snmpSrv["allow_sources"] = v.GetAllowSources()

		if err := d.Set("snmp_server", []interface{}{snmpSrv}); err != nil {
			return err
		}
	}

	if _, o := d.GetOkExists("healthcheck_server"); o {
		v := appliance.GetHealthcheckServer()
		healthSrv := make(map[string]interface{})
		healthSrv["enabled"] = v.GetEnabled()
		healthSrv["port"] = v.GetPort()
		healthSrv["allow_sources"] = v.GetAllowSources()

		if err := d.Set("healthcheck_server", []interface{}{healthSrv}); err != nil {
			return err
		}
	}

	if _, o := d.GetOkExists("prometheus_exporter"); o {
		v := appliance.GetPrometheusExporter()
		exporter := make(map[string]interface{})
		exporter["enabled"] = v.GetEnabled()
		exporter["port"] = v.GetPort()
		exporter["allow_sources"] = v.GetAllowSources()

		if err := d.Set("prometheus_exporter", []interface{}{exporter}); err != nil {
			return err
		}
	}

	if _, o := d.GetOkExists("ping"); o {
		v := appliance.GetPing()
		ping := make(map[string]interface{})
		ping["allow_sources"] = v.GetAllowSources()

		if err := d.Set("ping", []interface{}{ping}); err != nil {
			return err
		}
	}

	if _, o := d.GetOkExists("log_server"); o {
		v := appliance.GetLogServer()
		logsrv := make(map[string]interface{})
		logsrv["enabled"] = v.GetEnabled()
		logsrv["retention_days"] = v.GetRetentionDays()

		if err := d.Set("log_server", []interface{}{logsrv}); err != nil {
			return err
		}
	}

	if _, o := d.GetOkExists("controller"); o {
		v := appliance.GetController()
		ctrl := make(map[string]interface{})
		ctrl["enabled"] = v.GetEnabled()

		if err := d.Set("controller", []interface{}{ctrl}); err != nil {
			return err
		}
	}

	if _, o := d.GetOkExists("gateway"); o {
		gateway, err := flatttenApplianceGateway(appliance.GetGateway())
		if err != nil {
			return err
		}
		if err := d.Set("gateway", gateway); err != nil {
			return err
		}
	}

	if _, o := d.GetOkExists("log_forwarder"); o {
		logforward, err := flatttenApplianceLogForwarder(appliance.GetLogForwarder())
		if err != nil {
			return err
		}
		if err := d.Set("log_forwarder", logforward); err != nil {
			return fmt.Errorf("Unable to read log fowarder %s", err)
		}
	}

	if _, o := d.GetOkExists("connector"); o {
		iot, err := flatttenApplianceConnector(appliance.GetConnector())
		if err != nil {
			return err
		}
		if err := d.Set("connector", iot); err != nil {
			return fmt.Errorf("Unable to read connectors %s", err)
		}
	}

	if v, o := appliance.GetRsyslogDestinationsOk(); o != false {
		rsyslogs := make([]map[string]interface{}, 0)
		for _, rsys := range *v {
			rsyslog := make(map[string]interface{})
			rsyslog["selector"] = rsys.GetSelector()
			rsyslog["template"] = rsys.GetTemplate()
			rsyslog["destination"] = rsys.GetDestination()
			rsyslogs = append(rsyslogs, rsyslog)
		}

		if err := d.Set("rsyslog_destinations", rsyslogs); err != nil {
			return err
		}
	}

	if v, o := appliance.GetHostnameAliasesOk(); o != false {

		if err := d.Set("hostname_aliases", v); err != nil {
			return err
		}
	}
	return nil
}

func flatttenApplianceGateway(in openapi.ApplianceAllOfGateway) ([]map[string]interface{}, error) {
	var gateways []map[string]interface{}
	gateway := make(map[string]interface{})
	if v, o := in.GetEnabledOk(); o != false {
		gateway["enabled"] = v
	}

	if v, o := in.GetVpnOk(); o != false {
		vpn := make(map[string]interface{})
		if v, o := v.GetWeightOk(); o {
			vpn["weight"] = *v
		}
		if v, o := v.GetAllowDestinationsOk(); o {
			destinations := make([]map[string]interface{}, 0)
			for _, d := range *v {
				destination := make(map[string]interface{})
				destination["address"] = d.GetAddress()
				destination["netmask"] = d.GetNetmask()
				destination["nic"] = d.GetNic()
				destinations = append(destinations, destination)
			}
			vpn["allow_destinations"] = destinations
		}
		gateway["vpn"] = []map[string]interface{}{vpn}
	}

	gateways = append(gateways, gateway)
	return gateways, nil
}

func flatttenApplianceLogForwarder(in openapi.ApplianceAllOfLogForwarder) ([]map[string]interface{}, error) {
	var logforwarders []map[string]interface{}
	logforward := make(map[string]interface{})
	if v, o := in.GetEnabledOk(); o != false {
		logforward["enabled"] = *v
	}

	if v, o := in.GetElasticsearchOk(); o != false {
		elasticsearch := make(map[string]interface{})
		if v, o := v.GetUrlOk(); o != false {
			elasticsearch["url"] = *v
		}
		if v, o := v.GetAwsIdOk(); o != false {
			elasticsearch["aws_id"] = *v
		}
		if v, o := v.GetAwsSecretOk(); o != false {
			elasticsearch["aws_secret"] = *v
		}
		if v, o := v.GetAwsRegionOk(); o != false {
			elasticsearch["aws_region"] = *v
		}
		if v, o := v.GetUseInstanceCredentialsOk(); o != false {
			elasticsearch["use_instance_credentials"] = *v
		}
		if v, o := v.GetRetentionDaysOk(); o != false {
			elasticsearch["retention_days"] = *v
		}
		logforward["elasticsearch"] = []map[string]interface{}{elasticsearch}
	}
	if v, o := in.GetTcpClientsOk(); o != false {
		tcpClientList := make([]map[string]interface{}, 0)
		for _, tcpClient := range *v {
			client := make(map[string]interface{})
			client["name"] = tcpClient.GetName()
			client["host"] = tcpClient.GetHost()
			client["port"] = tcpClient.GetPort()
			client["format"] = tcpClient.GetFormat()
			client["use_tls"] = tcpClient.GetUseTLS()
			client["filter"] = tcpClient.GetFilter()
			tcpClientList = append(tcpClientList, client)
		}
		logforward["tcp_clients"] = tcpClientList
	}
	if v, o := in.GetAwsKinesesOk(); o != false {
		kinesesList := make([]map[string]interface{}, 0)
		for _, kineses := range *v {
			k := make(map[string]interface{})
			k["aws_id"] = kineses.GetAwsId()
			k["aws_secret"] = kineses.GetAwsSecret()
			k["aws_region"] = kineses.GetAwsRegion()
			k["use_instance_credentials"] = kineses.GetUseInstanceCredentials()
			k["type"] = kineses.GetType()
			k["stream_name"] = kineses.GetStreamName()
			k["batch_size"] = kineses.GetBatchSize()
			k["number_of_partition_keys"] = kineses.GetNumberOfPartitionKeys()
			k["filter"] = kineses.GetFilter()
			kinesesList = append(kinesesList, k)
		}
		logforward["aws_kineses"] = kinesesList
	}
	logforward["sites"] = in.GetSites()

	logforwarders = append(logforwarders, logforward)
	return logforwarders, nil
}

func flatttenApplianceConnector(in openapi.ApplianceAllOfConnector) ([]map[string]interface{}, error) {
	var connectors []map[string]interface{}
	connector := make(map[string]interface{})
	if v, o := in.GetEnabledOk(); o != false {
		connector["enabled"] = *v
	}
	if v, o := in.GetExpressClientsOk(); o != false {
		clients := make([]map[string]interface{}, 0)
		for _, client := range *v {
			c := make(map[string]interface{})
			c["name"] = client.GetName()
			c["device_id"] = client.GetDeviceId()
			alloweResources := make([]map[string]interface{}, 0)
			for _, arRaw := range client.GetAllowResources() {
				ar := make(map[string]interface{})
				ar["address"] = arRaw.GetAddress()
				ar["netmask"] = arRaw.GetNetmask()
				alloweResources = append(alloweResources, ar)
			}
			c["allow_resources"] = alloweResources
			c["snat_to_resources"] = client.GetSnatToResources()

			clients = append(clients, c)
		}
		connector["express_clients"] = clients
	}
	if v, o := in.GetAdvancedClientsOk(); o != false {
		clients := make([]map[string]interface{}, 0)
		for _, client := range *v {
			c := make(map[string]interface{})
			c["name"] = client.GetName()
			c["device_id"] = client.GetDeviceId()
			alloweResources, err := readAllowSourcesFromConfig(client.GetAllowResources())
			if err != nil {
				return nil, err
			}
			c["allow_resources"] = alloweResources
			c["snat_to_tunnel"] = client.GetSnatToTunnel()

			clients = append(clients, c)
		}
		connector["advanced_clients"] = clients
	}
	connectors = append(connectors, connector)
	return connectors, nil
}

func flattenApplianceClientInterface(in openapi.ApplianceAllOfClientInterface) ([]interface{}, error) {
	m := make(map[string]interface{})
	if v, o := in.GetProxyProtocolOk(); o != false {
		m["proxy_protocol"] = *v
	}
	if v, o := in.GetHostnameOk(); o != false {
		m["hostname"] = *v
	}
	if v, o := in.GetHttpsPortOk(); o != false {
		m["https_port"] = *v
	}
	if v, o := in.GetDtlsPortOk(); o != false {
		m["dtls_port"] = *v
	}
	if _, o := in.GetAllowSourcesOk(); o != false {
		allowSources, err := readAllowSourcesFromConfig(in.GetAllowSources())
		if err != nil {
			return nil, err
		}
		m["allow_sources"] = allowSources
	}
	if v, o := in.GetOverrideSpaModeOk(); o != false {
		m["override_spa_mode"] = v
	}

	return []interface{}{m}, nil
}

func flattenAppliancePeerInterface(in openapi.ApplianceAllOfPeerInterface) ([]interface{}, error) {
	m := make(map[string]interface{})
	if v, o := in.GetHostnameOk(); o != false {
		m["hostname"] = v
	}
	if v, o := in.GetHttpsPortOk(); o != false {
		m["https_port"] = v
	}
	if _, o := in.GetAllowSourcesOk(); o != false {
		allowSources, err := readAllowSourcesFromConfig(in.GetAllowSources())
		if err != nil {
			return nil, err
		}
		m["allow_sources"] = allowSources
	}
	return []interface{}{m}, nil
}

func flattenApplianceAdminInterface(in openapi.ApplianceAllOfAdminInterface) ([]interface{}, error) {
	m := make(map[string]interface{})
	if v, o := in.GetHostnameOk(); o != false {
		m["hostname"] = *v
	}
	if v, o := in.GetHttpsPortOk(); o != false {
		m["https_port"] = *v
	}
	if v, o := in.GetHttpsCiphersOk(); o != false {
		m["https_ciphers"] = *v
	}

	if _, o := in.GetAllowSourcesOk(); o != false {
		allowSources, err := readAllowSourcesFromConfig(in.GetAllowSources())
		if err != nil {
			return nil, err
		}
		m["allow_sources"] = allowSources
	}
	return []interface{}{m}, nil
}

func flattenApplianceNetworking(in openapi.ApplianceAllOfNetworking) ([]map[string]interface{}, error) {
	var networkings []map[string]interface{}
	networking := make(map[string]interface{})

	if v, o := in.GetHostsOk(); o != false {
		hosts := make([]map[string]interface{}, 0)
		for _, h := range *v {
			host := make(map[string]interface{})
			if v, o := h.GetAddressOk(); o != false {
				host["address"] = *v
			}
			if v, o := h.GetHostnameOk(); o != false {
				host["hostname"] = *v
			}
			hosts = append(hosts, host)
		}
		networking["hosts"] = hosts
	}

	if v, o := in.GetNicsOk(); o != false {
		nics := make([]map[string]interface{}, 0)
		for _, h := range *v {
			nic := make(map[string]interface{})
			if v, o := h.GetEnabledOk(); o != false {
				nic["enabled"] = *v
			}
			if v, o := h.GetNameOk(); o != false {
				nic["name"] = *v
			}
			if v, o := h.GetMtuOk(); o != false {
				nic["mtu"] = *v
			}

			if v, o := h.GetIpv4Ok(); o != false {
				dhcp := make(map[string]interface{})
				staticList := make([]map[string]interface{}, 0)
				dhcpValue := v.GetDhcp()

				if v, o := dhcpValue.GetEnabledOk(); o != false {
					dhcp["enabled"] = *v
				}
				if v, o := dhcpValue.GetDnsOk(); o != false {
					dhcp["dns"] = *v
				}
				if v, o := dhcpValue.GetRoutersOk(); o != false {
					dhcp["routers"] = *v
				}
				if v, o := dhcpValue.GetNtpOk(); o != false {
					dhcp["ntp"] = *v
				}
				for _, s := range v.GetStatic() {
					static := make(map[string]interface{})
					if v, o := s.GetAddressOk(); o {
						static["address"] = *v
					}
					if v, o := s.GetNetmaskOk(); o {
						static["netmask"] = *v
					}
					if v, o := s.GetHostnameOk(); o {
						static["hostname"] = *v
					}
					if v, o := s.GetSnatOk(); o {
						static["snat"] = *v
					}
					staticList = append(staticList, static)
				}
				ipv4map := make(map[string]interface{})
				ipv4map["dhcp"] = []map[string]interface{}{dhcp}
				ipv4map["static"] = staticList
				if v, o := v.GetVirtualIpOk(); o && len(*v) > 0 {
					ipv4map["virtual_ip"] = *v
				}
				nic["ipv4"] = []map[string]interface{}{ipv4map}
			}
			if v, o := h.GetIpv6Ok(); o != false {
				dhcp := make(map[string]interface{})
				staticList := make([]map[string]interface{}, 0)
				dhcpValue := v.GetDhcp()
				if v, o := dhcpValue.GetEnabledOk(); o != false {
					dhcp["enabled"] = *v
				}
				if v, o := dhcpValue.GetDnsOk(); o != false {
					dhcp["dns"] = *v
				}

				if v, o := dhcpValue.GetNtpOk(); o != false {
					dhcp["ntp"] = *v
				}
				for _, s := range v.GetStatic() {
					static := make(map[string]interface{})
					if v, o := s.GetAddressOk(); o {
						static["address"] = *v
					}
					if v, o := s.GetNetmaskOk(); o {
						static["netmask"] = *v
					}
					if v, o := s.GetHostnameOk(); o {
						static["hostname"] = *v
					}
					if v, o := s.GetSnatOk(); o {
						static["snat"] = *v
					}
					staticList = append(staticList, static)
				}
				ipv6map := make(map[string]interface{})
				ipv6map["dhcp"] = []map[string]interface{}{dhcp}
				ipv6map["static"] = staticList
				if v, o := v.GetVirtualIpOk(); o && len(*v) > 0 {
					ipv6map["virtual_ip"] = *v
				}
				nic["ipv6"] = []map[string]interface{}{ipv6map}

			}
			nics = append(nics, nic)
		}
		networking["nics"] = nics
	}
	if v, o := in.GetDnsServersOk(); o {
		networking["dns_servers"] = *v
	}
	if v, o := in.GetDnsDomainsOk(); o {
		networking["dns_domains"] = *v
	}
	if v, o := in.GetRoutesOk(); o {
		routes := make([]map[string]interface{}, 0)
		for _, r := range *v {
			route := make(map[string]interface{})
			if v, o := r.GetAddressOk(); o {
				route["address"] = *v
			}
			if v, o := r.GetNetmaskOk(); o {
				route["netmask"] = *v
			}
			if v, o := r.GetGatewayOk(); o {
				route["gateway"] = *v
			}
			if v, o := r.GetNicOk(); o {
				route["nic"] = *v
			}
			routes = append(routes, route)
		}
		networking["routes"] = routes

	}
	networkings = append(networkings, networking)

	return networkings, nil
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

	if d.HasChange("connect_to_peers_using_client_port_with_spa") {
		originalAppliance.SetConnectToPeersUsingClientPortWithSpa(d.Get("connect_to_peers_using_client_port_with_spa").(bool))
	}

	if d.HasChange("client_interface") {
		_, n := d.GetChange("client_interface")
		cinterface, err := readClientInterfaceFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		originalAppliance.SetClientInterface(cinterface)
	}

	if d.HasChange("peer_interface") {
		_, n := d.GetChange("peer_interface")
		pinterface, err := readPeerInterfaceFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		originalAppliance.SetPeerInterface(pinterface)
	}

	if d.HasChange("admin_interface") {
		_, a := d.GetChange("admin_interface")
		ainterface, err := readAdminInterfaceFromConfig(a.([]interface{}))
		if err != nil {
			return err
		}
		originalAppliance.SetAdminInterface(ainterface)
	}

	if d.HasChange("networking") {
		_, n := d.GetChange("networking")
		networking, err := readNetworkingFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		originalAppliance.SetNetworking(networking)
	}

	if d.HasChange("ntp") {
		_, n := d.GetChange("ntp")
		ntp, err := readNTPFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		originalAppliance.SetNtp(ntp)
	}

	if d.HasChange("ssh_server") {
		_, n := d.GetChange("ssh_server")
		sshServer, err := readSSHServerFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		originalAppliance.SetSshServer(sshServer)
	}

	if d.HasChange("snmp_server") {
		_, n := d.GetChange("snmp_server")
		srv, err := readSNMPServerFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		originalAppliance.SetSnmpServer(srv)
	}

	if d.HasChange("healthcheck_server") {
		_, n := d.GetChange("healthcheck_server")
		srv, err := readHealthcheckServerFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		originalAppliance.SetHealthcheckServer(srv)
	}

	if d.HasChange("prometheus_exporter") {
		_, n := d.GetChange("prometheus_exporter")
		exporter, err := readPrometheusExporterFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		originalAppliance.SetPrometheusExporter(exporter)
	}

	if d.HasChange("ping") {
		_, n := d.GetChange("ping")
		p, err := readPingFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		originalAppliance.SetPing(p)
	}

	if d.HasChange("log_server") {
		_, n := d.GetChange("log_server")
		logSrv, err := readLogServerFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		originalAppliance.SetLogServer(logSrv)
	}

	if d.HasChange("controller") {
		_, n := d.GetChange("controller")
		ctrl, err := readControllerFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		originalAppliance.SetController(ctrl)
	}

	if d.HasChange("gateway") {
		_, n := d.GetChange("gateway")
		gw, err := readGatewayFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		originalAppliance.SetGateway(gw)
	}

	if d.HasChange("log_forwarder") {
		_, n := d.GetChange("log_forwarder")
		lf, err := readLogForwardFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		originalAppliance.SetLogForwarder(lf)
	}

	if d.HasChange("connector") {
		_, n := d.GetChange("connector")
		iot, err := readApplianceConnectorFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		originalAppliance.SetConnector(iot)
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
		if v, ok := raw["proxy_protocol"]; ok {
			cinterface.SetProxyProtocol(v.(bool))
		}
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
			as, err := listToMapList(v.([]interface{}))
			if err != nil {
				return cinterface, err
			}
			allowSources, err := readAllowSourcesFromConfig(as)
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
			as, err := listToMapList(v.([]interface{}))
			if err != nil {
				return pinterf, err
			}
			allowSources, err := readAllowSourcesFromConfig(as)
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
		if v := raw["https_ciphers"]; len(v.([]interface{})) > 0 {
			ciphers := make([]string, 0)
			for _, cipher := range v.([]interface{}) {
				ciphers = append(ciphers, cipher.(string))
			}
			aInterface.SetHttpsCiphers(ciphers)
		}

		if v := raw["allow_sources"]; len(v.([]interface{})) > 0 {
			as, err := listToMapList(v.([]interface{}))
			if err != nil {
				return aInterface, err
			}
			allowSources, err := readAllowSourcesFromConfig(as)
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
			as, err := listToMapList(v.([]interface{}))
			if err != nil {
				return sshServer, err
			}
			allowSources, err := readAllowSourcesFromConfig(as)
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
			as, err := listToMapList(v.([]interface{}))
			if err != nil {
				return server, err
			}
			allowSources, err := readAllowSourcesFromConfig(as)
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
			as, err := listToMapList(v.([]interface{}))
			if err != nil {
				return server, err
			}
			allowSources, err := readAllowSourcesFromConfig(as)
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
		if v := r["vpn"].([]interface{}); len(v) > 0 {
			vpn := openapi.ApplianceAllOfGatewayVpn{}
			for _, s := range v {
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
			as, err := listToMapList(v.([]interface{}))
			if err != nil {
				return val, err
			}
			allowSources, err := readAllowSourcesFromConfig(as)
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
			as, err := listToMapList(v.([]interface{}))
			if err != nil {
				return val, err
			}
			allowSources, err := readAllowSourcesFromConfig(as)
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

		if v := raw["elasticsearch"].([]interface{}); len(v) > 0 {
			elasticsearch := openapi.Elasticsearch{}
			for _, s := range v {
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
				if v, ok := r["filter"]; ok {
					tcpClient.SetFilter(v.(string))
				}
				tcpClients = append(tcpClients, tcpClient)
			}
			val.SetTcpClients(tcpClients)
		}
		if v := raw["aws_kineses"]; len(v.([]interface{})) > 0 {
			awsKineses := make([]openapi.AwsKinesis, 0)
			for _, awsk := range v.([]interface{}) {
				kinesis := openapi.AwsKinesis{}
				row := awsk.(map[string]interface{})
				if v, ok := row["aws_id"]; ok {
					kinesis.SetAwsId(v.(string))
				}
				if v, ok := row["aws_secret"]; ok {
					kinesis.SetAwsSecret(v.(string))
				}
				if v, ok := row["aws_region"]; ok {
					kinesis.SetAwsRegion(v.(string))
				}
				if v, ok := row["use_instance_credentials"]; ok {
					kinesis.SetUseInstanceCredentials(v.(bool))
				}
				if v, ok := row["type"]; ok {
					kinesis.SetType(v.(string))
				}
				if v, ok := row["stream_name"]; ok {
					kinesis.SetStreamName(v.(string))
				}
				if v, ok := row["batch_size"]; ok {
					kinesis.SetBatchSize(int32(v.(int)))
				}
				if v, ok := row["number_of_partition_keys"]; ok {
					kinesis.SetNumberOfPartitionKeys(int32(v.(int)))
				}
				if v, ok := row["filter"]; ok {
					kinesis.SetFilter(v.(string))
				}
			}
			val.SetAwsKineses(awsKineses)
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

func readApplianceConnectorFromConfig(connectors []interface{}) (openapi.ApplianceAllOfConnector, error) {
	val := openapi.ApplianceAllOfConnector{}
	for _, connector := range connectors {
		if connector == nil {
			continue
		}

		raw := connector.(map[string]interface{})

		if v, ok := raw["enabled"]; ok {
			val.SetEnabled(v.(bool))
		}
		if v := raw["express_clients"]; len(v.([]interface{})) > 0 {
			clients := make([]openapi.ApplianceAllOfConnectorExpressClients, 0)
			for _, c := range v.([]interface{}) {
				client := openapi.ApplianceAllOfConnectorExpressClients{}
				r := c.(map[string]interface{})
				if v, ok := r["name"]; ok {
					client.SetName(v.(string))
				}
				if v, ok := r["device_id"]; ok {
					client.SetDeviceId(v.(string))
				}
				// allowed sources
				if v := r["allow_resources"]; len(v.([]interface{})) > 0 {
					as, err := listToMapList(v.([]interface{}))
					if err != nil {
						return val, err
					}
					sources, err := readAllowSourcesFromConfig(as)
					if err != nil {
						return val, err
					}
					allowedResources := make([]openapi.ApplianceAllOfConnectorAllowResources, 0)
					for _, s := range sources {
						ar := openapi.ApplianceAllOfConnectorAllowResources{}
						if v, ok := s["address"]; ok {
							ar.SetAddress(v.(string))
						}
						if v, ok := s["netmask"]; ok {
							ar.SetNetmask(int32(v.(int)))
						}
						allowedResources = append(allowedResources, ar)
					}
					// TODO loop source convert to openapi.ApplianceAllOfConnectorAllowResources
					client.SetAllowResources(allowedResources)
				}
				if v, ok := r["snat_to_resources"]; ok {
					client.SetSnatToResources(v.(bool))
				}
				clients = append(clients, client)
			}
			val.SetExpressClients(clients)
		}
		if v := raw["advanced_clients"]; len(v.([]interface{})) > 0 {
			clients := make([]openapi.ApplianceAllOfConnectorAdvancedClients, 0)
			for _, c := range v.([]interface{}) {
				client := openapi.ApplianceAllOfConnectorAdvancedClients{}
				r := c.(map[string]interface{})
				if v, ok := r["name"]; ok {
					client.SetName(v.(string))
				}
				if v, ok := r["device_id"]; ok {
					client.SetDeviceId(v.(string))
				}
				if v := r["allow_resources"]; len(v.([]interface{})) > 0 {
					as, err := listToMapList(v.([]interface{}))
					if err != nil {
						return val, err
					}
					sources, err := readAllowSourcesFromConfig(as)
					if err != nil {
						return val, err
					}
					client.SetAllowResources(sources)
				}
				if v, ok := r["snat_to_tunnel"]; ok {
					client.SetSnatToTunnel(v.(bool))
				}
				clients = append(clients, client)
			}
			val.SetAdvancedClients(clients)
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
