package appgate

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/appgate/sdp-api-client-go/api/v18/openapi"

	"github.com/google/uuid"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateAppliance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppgateApplianceCreate,
		ReadContext:   resourceAppgateApplianceRead,
		UpdateContext: resourceAppgateApplianceUpdate,
		DeleteContext: resourceAppgateApplianceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{

			"appliance_id": resourceUUID(),
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

			"tags": tagsSchema(),

			"hostname": {
				Type:        schema.TypeString,
				Deprecated:  "appliance hostname is deprecated as of 5.4.",
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
				Computed:    true,
			},

			"connect_to_peers_using_client_port_with_spa": {
				Type:        schema.TypeBool,
				Deprecated:  "connect_to_peers_using_client_port_with_spa is deprecated as of 5.4. It will always be enabled when the support for peerInterface is removed.",
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

						"allow_sources": allowSourcesSchema(),

						"override_spa_mode": {
							Type:     schema.TypeString,
							Optional: true,
							// We will have a default value here instead of omitting the attribute when its disabled.
							// https://github.com/appgate/terraform-provider-appgatesdp/issues/117#issuecomment-846381509
							Default: "Disabled",
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
				// TODO:
				// Temporary removed this warning, since its not scheduled to be removed until the version after 5.5
				// and since its still required for all existing supported versions, we will not show this error for the users.
				//
				// Deprecated: "peer_interface is deprecated as of 5.4. All connections will be handled by clientInterface and adminInterface in the future. The hostname field is used as identifier and will take over the hostname field in the root of Appliance when this interface is removed.",
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
							Default:  444,
						},
						"allow_sources": allowSourcesSchema(),
					},
				},
			},

			"admin_interface": adminInterfaceSchema(),

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
													Type:     schema.TypeList,
													MaxItems: 1,
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
																Type:       schema.TypeString,
																Deprecated: "Removed in >= 5.4",
																Optional:   true,
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
																Type:       schema.TypeString,
																Deprecated: "Removed in >= 5.4",
																Required:   true,
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
							Set:         schema.HashString,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},

						"dns_domains": {
							Type:        schema.TypeSet,
							Description: "DNS Search domains.",
							Set:         schema.HashString,
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

						"allow_sources": allowSourcesSchema(),

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

						"allow_sources": allowSourcesSchema(),
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

						"allow_sources": allowSourcesSchema(),
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

						"allow_sources": allowSourcesSchema(),
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
						"allow_sources": allowSourcesSchema(),
					},
				},
			},

			"log_server": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:             schema.TypeBool,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
						},
						"retention_days": {
							Type:             schema.TypeInt,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
						},
					},
				},
			},

			"controller": controllerSchema(),

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
										Required: true,
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
									"compatibility_mode": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"authentication": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Required: true,
												},
												"token": {
													Type:      schema.TypeString,
													Required:  true,
													Sensitive: true,
												},
											},
										},
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
							Computed:    true,
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
									"dnat_to_resource": {
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

									"allow_resources": allowSourcesSchema(),

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

			"portal": {
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

						"https_p12": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subject_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"content": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"password": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"proxy_p12s": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subject_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"content": {
										Type:        schema.TypeString,
										Description: "path to file with p12",
										Optional:    true,
									},
									"password": {
										Type:      schema.TypeString,
										Sensitive: true,
										Optional:  true,
									},
									"verify_upstream": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},

						"profiles": {
							Type:        schema.TypeList,
							Description: "Names of the profiles in this Collective to use in the Portal.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},

						"external_profiles": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"url": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"sign_in_customization": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"background_color": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Changes the background color on the sign-in page. In hexadecimal format.",
									},
									"background_image": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Changes the background image on the sign-in page. Must be in PNG, JPEG or GIF format.",
									},
									"background_image_checksum": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"logo": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Changes the logo on the sign-in page. Must be in PNG, JPEG or GIF format.",
									},
									"logo_checksum": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"text": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Adds a text to the sign-in page.",
									},
									"text_color": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Changes the text color on the sign-in page. In hexadecimal format.",
									},
									"auto_redirect": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
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
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAppgateApplianceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating Appliance with name: %s", d.Get("name").(string))
	var diags diag.Diagnostics
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.AppliancesApi
	currentVersion := meta.(*Client).ApplianceVersion
	args := openapi.NewApplianceWithDefaults()
	if v, ok := d.GetOk("appliance_id"); ok {
		args.SetId(v.(string))
	}
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
			return diag.FromErr(err)
		}
		args.SetClientInterface(cinterface)
	}

	if p, ok := d.GetOk("peer_interface"); ok {
		if currentVersion.GreaterThanOrEqual(Appliance60Version) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("peer_interface is not supported in %s", currentVersion.String()),
				Detail: `peer_interface is removed in >= 6.0.
All connections will be handled by client_interface and admin_interface in the future.
The hostname field is used as identifier and will take over the hostname field in
the root of Appliance when this interface is removed.`,
			})
			return diags
		}
		pinterface, err := readPeerInterfaceFromConfig(p.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetPeerInterface(pinterface)
	}

	if a, ok := d.GetOk("admin_interface"); ok {
		ainterface, err := readAdminInterfaceFromConfig(a.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetAdminInterface(ainterface)
	}

	if n, ok := d.GetOk("networking"); ok {
		network, err := readNetworkingFromConfig(n.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetNetworking(network)
	}

	if n, ok := d.GetOk("ssh_server"); ok {
		sshServer, err := readSSHServerFromConfig(n.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetSshServer(sshServer)
	}

	if n, ok := d.GetOk("snmp_server"); ok {
		srv, err := readSNMPServerFromConfig(n.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetSnmpServer(srv)
	}

	if n, ok := d.GetOk("healthcheck_server"); ok {
		srv, err := readHealthcheckServerFromConfig(n.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetHealthcheckServer(srv)
	}

	if n, ok := d.GetOk("prometheus_exporter"); ok {
		exporter, err := readPrometheusExporterFromConfig(n.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetPrometheusExporter(exporter)
	}

	if n, ok := d.GetOk("ping"); ok {
		p, err := readPingFromConfig(n.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetPing(p)
	}

	if n, ok := d.GetOk("ntp"); ok {
		ntp, err := readNTPFromConfig(n.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetNtp(ntp)
	}

	if v, ok := d.GetOk("log_server"); ok {
		logSrv, err := readLogServerFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		if logSrv.GetEnabled() {
			args.SetLogServer(logSrv)
		} else {
			args.LogServer = nil
		}
	}

	if v, ok := d.GetOk("controller"); ok {
		ctrl, err := readControllerFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetController(ctrl)
	}

	if v, ok := d.GetOk("gateway"); ok {
		gw, err := readGatewayFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetGateway(gw)
	}

	if v, ok := d.GetOk("log_forwarder"); ok {
		lf, err := readLogForwardFromConfig(v.([]interface{}), currentVersion)
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetLogForwarder(lf)
	}

	if v, ok := d.GetOk("connector"); ok {
		connector, err := readApplianceConnectorFromConfig(currentVersion, v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetConnector(connector)
	}

	if v, ok := d.GetOk("rsyslog_destinations"); ok {
		rsyslog, err := readRsyslogDestinationFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetRsyslogDestinations(rsyslog)
	}

	if v, ok := d.GetOk("portal"); ok {
		if !currentVersion.GreaterThanOrEqual(Appliance54Version) {
			return diag.Errorf("appliance.portal requires %s, you are using %q client v%d", Appliance54Version, currentVersion, meta.(*Client).ClientVersion)
		}
		portal, err := readAppliancePortalFromConfig(d, v.([]interface{}), currentVersion)
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetPortal(portal)
	}

	if v, ok := d.GetOk("hostname_aliases"); ok {
		hostnames, err := readHostnameAliasesFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetHostnameAliases(hostnames)
	}

	request := api.AppliancesPost(ctx)
	request = request.Appliance(*args)
	appliance, _, err := request.Authorization(token).Execute()
	if err != nil {
		return diag.Errorf("Could not create appliance %s", prettyPrintAPIError(err))
	}

	d.SetId(appliance.GetId())

	resourceAppgateApplianceRead(ctx, d, meta)
	return diags
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
				if item, ok := ipv4Data["dhcp"].([]interface{}); ok && len(item) > 0 && item[0] != nil {
					ipv4networking.SetDhcp(readNetworkIpv4DhcpFromConfig(item[0].(map[string]interface{})))
				}
				if item := ipv4Data["static"]; len(item.([]interface{})) > 0 {
					ipv4networking.SetStatic(readNetworkIpv4StaticFromConfig(item.([]interface{})))
				}
				if v, ok := ipv4Data["virtual_ip"]; ok && len(v.(string)) > 0 {
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
				if v, ok := ipv6Data["virtual_ip"]; ok && len(v.(string)) > 0 {
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

func flattenAllowSources(input []openapi.AllowSourcesInner) ([]map[string]interface{}, error) {
	r := make([]map[string]interface{}, 0)
	for _, raw := range input {
		row := make(map[string]interface{}, 0)
		if v, ok := raw.GetAddressOk(); ok {
			row["address"] = v
		}
		if v, ok := raw.GetNetmaskOk(); ok {
			row["netmask"] = v
		}
		if v, ok := raw.GetNicOk(); ok && *v != "" {
			row["nic"] = v
		}
		r = append(r, row)
	}
	return r, nil
}

func flattenAllowResources(input []openapi.AllowResourcesInner) ([]map[string]interface{}, error) {
	r := make([]map[string]interface{}, 0)
	for _, raw := range input {
		row := make(map[string]interface{}, 0)
		if v, ok := raw.GetAddressOk(); ok {
			row["address"] = v
		}
		if v, ok := raw.GetNetmaskOk(); ok {
			row["netmask"] = v
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

func resourceAppgateApplianceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading Appliance Name: %s", d.Get("name").(string))
	var diags diag.Diagnostics

	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.AppliancesApi
	currentVersion := meta.(*Client).ApplianceVersion

	request := api.AppliancesIdGet(ctx, d.Id())
	appliance, res, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.Errorf("Failed to read Appliance, %s", err)
	}
	d.Set("appliance_id", appliance.GetId())
	d.Set("name", appliance.GetName())
	d.Set("tags", appliance.GetTags())
	d.Set("notes", appliance.GetNotes())
	d.Set("hostname", appliance.GetHostname())

	if err := d.Set("site", appliance.GetSite()); err != nil {
		return diag.Errorf("Error setting appliance.site %s", err)
	}
	if err := d.Set("customization", appliance.GetCustomization()); err != nil {
		return diag.Errorf("Error setting appliance.customization %s", err)
	}

	if err := d.Set("connect_to_peers_using_client_port_with_spa", appliance.GetConnectToPeersUsingClientPortWithSpa()); err != nil {
		return diag.Errorf("Error setting appliance.connect_to_peers_using_client_port_with_spa %s", err)
	}

	if v, ok := appliance.GetClientInterfaceOk(); ok {
		ci, err := flattenApplianceClientInterface(*v)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("client_interface", ci)
	}

	if v, ok := appliance.GetPeerInterfaceOk(); ok {
		peerInterface, err := flattenAppliancePeerInterface(*v)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("peer_interface", peerInterface)
	}

	if v, ok := appliance.GetAdminInterfaceOk(); ok {
		adminInterface, err := flattenApplianceAdminInterface(*v)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("admin_interface", adminInterface); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := appliance.GetNetworkingOk(); ok {
		networking, err := flattenApplianceNetworking(*v)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("networking", networking); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := appliance.GetNtpOk(); ok {
		ntp := make(map[string]interface{})
		servers := make([]map[string]interface{}, 0)
		for _, v := range v.GetServers() {
			srv := make(map[string]interface{})
			srv["hostname"] = v.GetHostname()
			srv["key_type"] = v.GetKeyType()
			srv["key"] = v.GetKey()
			servers = append(servers, srv)
		}
		ntp["servers"] = servers
		if err := d.Set("ntp", []interface{}{ntp}); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := appliance.GetSshServerOk(); ok {
		sshServer := make(map[string]interface{})
		sshServer["enabled"] = v.GetEnabled()
		sshServer["port"] = v.GetPort()
		as, err := flattenAllowSources(v.GetAllowSources())
		if err != nil {
			return diag.FromErr(err)
		}
		sshServer["allow_sources"] = as
		sshServer["password_authentication"] = v.GetPasswordAuthentication()

		if err := d.Set("ssh_server", []interface{}{sshServer}); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := appliance.GetSnmpServerOk(); ok {
		snmpSrv := make(map[string]interface{})
		snmpSrv["enabled"] = v.GetEnabled()
		snmpSrv["tcp_port"] = v.GetTcpPort()
		snmpSrv["udp_port"] = v.GetUdpPort()
		snmpSrv["snmpd_conf"] = v.GetSnmpdConf()
		as, err := flattenAllowSources(v.GetAllowSources())
		if err != nil {
			return diag.FromErr(err)
		}
		snmpSrv["allow_sources"] = as

		if err := d.Set("snmp_server", []interface{}{snmpSrv}); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := appliance.GetHealthcheckServerOk(); ok {
		healthSrv := make(map[string]interface{})
		healthSrv["enabled"] = v.GetEnabled()
		healthSrv["port"] = v.GetPort()
		as, err := flattenAllowSources(v.GetAllowSources())
		if err != nil {
			return diag.FromErr(err)
		}
		healthSrv["allow_sources"] = as

		if err := d.Set("healthcheck_server", []interface{}{healthSrv}); err != nil {
			return diag.FromErr(err)
		}
	}
	if v, ok := appliance.GetPrometheusExporterOk(); ok {
		exporter := make(map[string]interface{})
		exporter["enabled"] = v.GetEnabled()
		exporter["port"] = v.GetPort()
		as, err := flattenAllowSources(v.GetAllowSources())
		if err != nil {
			return diag.FromErr(err)
		}
		exporter["allow_sources"] = as

		if err := d.Set("prometheus_exporter", []interface{}{exporter}); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := appliance.GetPingOk(); ok {
		ping := make(map[string]interface{})
		as, err := flattenAllowSources(v.GetAllowSources())
		if err != nil {
			return diag.FromErr(err)
		}
		ping["allow_sources"] = as

		if err := d.Set("ping", []interface{}{ping}); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := appliance.GetLogServerOk(); ok {
		logsrv := make(map[string]interface{})
		enabledLogServer := v.GetEnabled()
		// we will only save log_server to the state ifs enabled,
		// since all appliances include default log_server enabled: false in every response.
		if enabledLogServer {
			logsrv["enabled"] = enabledLogServer
			logsrv["retention_days"] = v.GetRetentionDays()
			if err := d.Set("log_server", []interface{}{logsrv}); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if v, ok := appliance.GetControllerOk(); ok {
		ctrl := make(map[string]interface{})
		ctrl["enabled"] = v.GetEnabled()

		if err := d.Set("controller", []interface{}{ctrl}); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := appliance.GetGatewayOk(); ok {
		gateway, err := flatttenApplianceGateway(*v)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("gateway", gateway); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := appliance.GetLogForwarderOk(); ok {
		logforward, err := flatttenApplianceLogForwarder(*v, currentVersion, d)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("log_forwarder", logforward); err != nil {
			return diag.Errorf("Unable to read log fowarder %s", err)
		}
	}

	if v, ok := appliance.GetConnectorOk(); ok {
		connector, err := flatttenApplianceConnector(currentVersion, *v)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("connector", connector); err != nil {
			return diag.Errorf("Unable to read connectors %s", err)
		}
	}

	if v, ok := appliance.GetRsyslogDestinationsOk(); ok {
		rsyslogs := make([]map[string]interface{}, 0)
		for _, rsys := range v {
			rsyslog := make(map[string]interface{})
			rsyslog["selector"] = rsys.GetSelector()
			rsyslog["template"] = rsys.GetTemplate()
			rsyslog["destination"] = rsys.GetDestination()
			rsyslogs = append(rsyslogs, rsyslog)
		}

		if err := d.Set("rsyslog_destinations", rsyslogs); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := appliance.GetPortalOk(); ok {
		portals := make([]map[string]interface{}, 0)
		portal := make(map[string]interface{})
		portal["enabled"] = v.GetEnabled()
		// get local state from the portal attribute, for values that are not included
		// in the response body
		var localPortal map[string]interface{}
		localPortalList := d.Get("portal").([]interface{})
		for _, l := range localPortalList {
			localPortal = l.(map[string]interface{})
		}
		if len(v.GetProxyP12s()) > 0 {
			proxyp12s, err := flattenAppliancePortalProxyp12s(localPortal, v.GetProxyP12s())
			if err != nil {
				return diag.FromErr(err)
			}
			portal["proxy_p12s"] = proxyp12s
		}
		https_p12, err := flattenApplianceProxyp12s(localPortal, v.GetHttpsP12())
		if err != nil {
			return diag.FromErr(err)
		}
		portal["https_p12"] = https_p12

		portal["profiles"] = v.GetProfiles()
		portal["external_profiles"] = v.GetExternalProfiles()
		if currentVersion.GreaterThanOrEqual(Appliance55Version) {
			signInCustomization, err := flattenAppliancePortalSignInCustomziation(d, localPortal, v.GetSignInCustomization())
			if err != nil {
				return diag.FromErr(err)
			}
			portal["sign_in_customization"] = signInCustomization
		}
		portals = append(portals, portal)
		if err := d.Set("portal", portals); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := appliance.GetHostnameAliasesOk(); ok {
		if err := d.Set("hostname_aliases", v); err != nil {
			return diag.FromErr(err)
		}
	}
	return diags
}

func flattenAppliancePortalProxyp12s(local map[string]interface{}, p12s []openapi.Portal12) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	for k, p12 := range p12s {
		raw := make(map[string]interface{})
		if len(p12.GetId()) < 1 {
			continue
		}
		raw["id"] = p12.GetId()
		raw["verify_upstream"] = p12.GetVerifyUpstream()
		raw["subject_name"] = p12.GetSubjectName()
		// content, and password not always known, not included in the response body
		if state, ok := local["proxy_p12s"].([]interface{}); ok && state[k] != nil {
			stateRow := state[k].(map[string]interface{})
			raw["content"] = stateRow["content"].(string)
			raw["password"] = stateRow["password"].(string)
		}
		result = append(result, raw)
	}

	return result, nil
}

func flattenApplianceProxyp12s(local map[string]interface{}, p12 openapi.P12) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	if len(p12.GetId()) < 1 {
		return result, nil
	}
	raw := make(map[string]interface{})
	raw["id"] = p12.GetId()
	raw["subject_name"] = p12.GetSubjectName()
	// content, and password not always known, not included in the response body
	if v, ok := local["https_p12"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		stateRow := v[0].(map[string]interface{})
		raw["content"] = stateRow["content"].(string)
		raw["password"] = stateRow["password"].(string)

	}
	result = append(result, raw)
	return result, nil
}

func flattenAppliancePortalSignInCustomziation(d *schema.ResourceData, local map[string]interface{}, customization openapi.PortalSignInCustomization) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	raw := make(map[string]interface{})

	raw["background_color"] = customization.GetBackgroundColor()
	raw["background_image"] = d.Get("portal.0.sign_in_customization.0.background_image").(string)
	raw["background_image_checksum"] = customization.GetBackgroundImage()
	raw["logo"] = d.Get("portal.0.sign_in_customization.0.logo").(string)
	raw["logo_checksum"] = customization.GetLogo()
	raw["text"] = customization.GetText()
	raw["text_color"] = customization.GetTextColor()
	raw["auto_redirect"] = customization.GetAutoRedirect()
	result = append(result, raw)
	return result, nil
}

func flatttenApplianceGateway(in openapi.ApplianceAllOfGateway) ([]map[string]interface{}, error) {
	var gateways []map[string]interface{}
	gateway := make(map[string]interface{})
	if v, ok := in.GetEnabledOk(); ok {
		gateway["enabled"] = v
	}

	if v, ok := in.GetVpnOk(); ok {
		vpn := make(map[string]interface{})
		if v, ok := v.GetWeightOk(); ok {
			vpn["weight"] = *v
		}
		if v, ok := v.GetAllowDestinationsOk(); ok {
			destinations := make([]map[string]interface{}, 0)
			for _, d := range v {
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

func flatttenApplianceLogForwarder(in openapi.ApplianceAllOfLogForwarder, currentVersion *version.Version, d *schema.ResourceData) ([]map[string]interface{}, error) {
	var logforwarders []map[string]interface{}
	logforward := make(map[string]interface{})
	if v, ok := in.GetEnabledOk(); ok {
		logforward["enabled"] = *v
	}

	if v, ok := in.GetElasticsearchOk(); ok {
		elasticsearch := make(map[string]interface{})
		if v, ok := v.GetUrlOk(); ok {
			elasticsearch["url"] = *v
		}
		if v, ok := v.GetAwsIdOk(); ok {
			elasticsearch["aws_id"] = *v
		}
		if v, ok := v.GetAwsSecretOk(); ok {
			elasticsearch["aws_secret"] = *v
		}
		if v, ok := v.GetAwsRegionOk(); ok {
			elasticsearch["aws_region"] = *v
		}
		if v, ok := v.GetUseInstanceCredentialsOk(); ok {
			elasticsearch["use_instance_credentials"] = *v
		}
		if v, ok := v.GetRetentionDaysOk(); ok {
			elasticsearch["retention_days"] = *v
		}
		if currentVersion.GreaterThanOrEqual(Appliance55Version) {
			if v, ok := v.GetCompatibilityModeOk(); ok {
				elasticsearch["compatibility_mode"] = *v
			}
			if authRaw, ok := v.GetAuthenticationOk(); ok {
				auth := make(map[string]interface{})
				if v, ok := authRaw.GetTypeOk(); ok {
					auth["type"] = v
				}

				// token is sensitive, so we won't get it in the response body, but we can lookup it from the state
				if state := d.Get("log_forwarder.0.elasticsearch.0.authentication").([]interface{}); len(state) > 0 && state[0] != nil {
					s := state[0].(map[string]interface{})
					if v, ok := s["token"]; ok {
						auth["token"] = v.(string)
					}
				} else if v, ok := authRaw.GetTokenOk(); ok {
					log.Printf("[DEBUG] Could not find log_forwarder.0.elasticsearch.0.authentication.token in state, fallback to API response")
					auth["token"] = v
				}
				elasticsearch["authentication"] = []map[string]interface{}{auth}
			}
		}
		logforward["elasticsearch"] = []map[string]interface{}{elasticsearch}
	}
	if v, ok := in.GetTcpClientsOk(); ok {
		tcpClientList := make([]map[string]interface{}, 0)
		for _, tcpClient := range v {
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
	if v, ok := in.GetAwsKinesesOk(); ok {
		kinesesList := make([]map[string]interface{}, 0)
		for _, kineses := range v {
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

func flatttenApplianceConnector(currentVersion *version.Version, in openapi.ApplianceAllOfConnector) ([]map[string]interface{}, error) {
	var connectors []map[string]interface{}
	connector := make(map[string]interface{})
	if v, ok := in.GetEnabledOk(); ok {
		connector["enabled"] = *v
	}
	if v, ok := in.GetExpressClientsOk(); ok {
		clients := make([]map[string]interface{}, 0)
		for _, client := range v {
			c := make(map[string]interface{})
			c["name"] = client.GetName()
			c["device_id"] = client.GetDeviceId()

			alloweResources, err := flattenAllowResources(client.GetAllowResources())
			if err != nil {
				return nil, err
			}
			c["allow_resources"] = alloweResources
			c["snat_to_resources"] = client.GetSnatToResources()
			if currentVersion.GreaterThanOrEqual(Appliance54Version) {
				c["dnat_to_resource"] = client.GetDnatToResource()
			}

			clients = append(clients, c)
		}
		connector["express_clients"] = clients
	}
	if v, ok := in.GetAdvancedClientsOk(); ok {
		clients := make([]map[string]interface{}, 0)
		for _, client := range v {
			c := make(map[string]interface{})
			c["name"] = client.GetName()
			c["device_id"] = client.GetDeviceId()
			alloweResources, err := flattenAllowSources(client.GetAllowResources())
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
	if v, ok := in.GetProxyProtocolOk(); ok {
		m["proxy_protocol"] = *v
	}
	if v, ok := in.GetHostnameOk(); ok {
		m["hostname"] = *v
	}
	if v, ok := in.GetHttpsPortOk(); ok {
		m["https_port"] = *v
	}
	if v, ok := in.GetDtlsPortOk(); ok {
		m["dtls_port"] = *v
	}
	if _, ok := in.GetAllowSourcesOk(); ok {
		allowSources, err := flattenAllowSources(in.GetAllowSources())
		if err != nil {
			return nil, err
		}
		m["allow_sources"] = allowSources
	}

	if v, ok := in.GetOverrideSpaModeOk(); ok {
		m["override_spa_mode"] = v
	} else {
		// If we dont get any from the response body, we will manually set it to Disabled to
		// make it explicit in the the tf plan.
		// https://github.com/appgate/terraform-provider-appgatesdp/issues/117#issuecomment-846381509
		m["override_spa_mode"] = "Disabled"
	}

	return []interface{}{m}, nil
}

func flattenAppliancePeerInterface(in openapi.ApplianceAllOfPeerInterface) ([]interface{}, error) {
	m := make(map[string]interface{})
	if v, ok := in.GetHostnameOk(); ok {
		m["hostname"] = v
	}
	if v, ok := in.GetHttpsPortOk(); ok {
		m["https_port"] = v
	}
	if _, ok := in.GetAllowSourcesOk(); ok {
		allowSources, err := flattenAllowSources(in.GetAllowSources())
		if err != nil {
			return nil, err
		}
		m["allow_sources"] = allowSources
	}
	return []interface{}{m}, nil
}

func flattenApplianceAdminInterface(in openapi.ApplianceAllOfAdminInterface) ([]interface{}, error) {
	m := make(map[string]interface{})
	if v, ok := in.GetHostnameOk(); ok {
		m["hostname"] = *v
	}
	if v, ok := in.GetHttpsPortOk(); ok {
		m["https_port"] = *v
	}
	if v, ok := in.GetHttpsCiphersOk(); ok {
		m["https_ciphers"] = v
	}

	if _, ok := in.GetAllowSourcesOk(); ok {
		allowSources, err := flattenAllowSources(in.GetAllowSources())
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

	if v, ok := in.GetHostsOk(); ok {
		hosts := make([]map[string]interface{}, 0)
		for _, h := range v {
			host := make(map[string]interface{})
			if v, ok := h.GetAddressOk(); ok {
				host["address"] = *v
			}
			if v, ok := h.GetHostnameOk(); ok {
				host["hostname"] = *v
			}
			hosts = append(hosts, host)
		}
		networking["hosts"] = hosts
	}

	if v, ok := in.GetNicsOk(); ok {
		nics := make([]map[string]interface{}, 0)
		for _, h := range v {
			nic := make(map[string]interface{})
			if v, ok := h.GetEnabledOk(); ok {
				nic["enabled"] = *v
			}
			if v, ok := h.GetNameOk(); ok {
				nic["name"] = *v
			}
			if v, ok := h.GetMtuOk(); ok {
				nic["mtu"] = *v
			}

			if v, ok := h.GetIpv4Ok(); ok {
				dhcp := make(map[string]interface{})
				staticList := make([]map[string]interface{}, 0)
				dhcpValue := v.GetDhcp()

				if v, ok := dhcpValue.GetEnabledOk(); ok {
					dhcp["enabled"] = *v
				}
				if v, ok := dhcpValue.GetDnsOk(); ok {
					dhcp["dns"] = *v
				}
				if v, ok := dhcpValue.GetRoutersOk(); ok {
					dhcp["routers"] = *v
				}
				if v, ok := dhcpValue.GetNtpOk(); ok {
					dhcp["ntp"] = *v
				}
				for _, s := range v.GetStatic() {
					static := make(map[string]interface{})
					if v, ok := s.GetAddressOk(); ok {
						static["address"] = *v
					}
					if v, ok := s.GetNetmaskOk(); ok {
						static["netmask"] = *v
					}
					if v, ok := s.GetSnatOk(); ok {
						static["snat"] = *v
					}
					staticList = append(staticList, static)
				}
				ipv4map := make(map[string]interface{})
				ipv4map["dhcp"] = []map[string]interface{}{dhcp}
				ipv4map["static"] = staticList
				if v, ok := v.GetVirtualIpOk(); ok && len(*v) > 0 {
					ipv4map["virtual_ip"] = *v
				}
				nic["ipv4"] = []map[string]interface{}{ipv4map}
			}
			if v, ok := h.GetIpv6Ok(); ok {
				dhcp := make(map[string]interface{})
				staticList := make([]map[string]interface{}, 0)
				dhcpValue := v.GetDhcp()
				if v, ok := dhcpValue.GetEnabledOk(); ok {
					dhcp["enabled"] = *v
				}
				if v, ok := dhcpValue.GetDnsOk(); ok {
					dhcp["dns"] = *v
				}

				if v, ok := dhcpValue.GetNtpOk(); ok {
					dhcp["ntp"] = *v
				}
				for _, s := range v.GetStatic() {
					static := make(map[string]interface{})
					if v, ok := s.GetAddressOk(); ok {
						static["address"] = *v
					}
					if v, ok := s.GetNetmaskOk(); ok {
						static["netmask"] = *v
					}
					if v, ok := s.GetSnatOk(); ok {
						static["snat"] = *v
					}
					staticList = append(staticList, static)
				}
				ipv6map := make(map[string]interface{})
				ipv6map["dhcp"] = []map[string]interface{}{dhcp}
				ipv6map["static"] = staticList
				if v, ok := v.GetVirtualIpOk(); ok && len(*v) > 0 {
					ipv6map["virtual_ip"] = *v
				}
				nic["ipv6"] = []map[string]interface{}{ipv6map}

			}
			nics = append(nics, nic)
		}
		networking["nics"] = nics
	}

	if v, ok := in.GetDnsServersOk(); ok {
		networking["dns_servers"] = schema.NewSet(schema.HashString, convertStringArrToInterface(v))
	}
	if _, ok := in.GetDnsDomainsOk(); ok {
		networking["dns_domains"] = schema.NewSet(schema.HashString, convertStringArrToInterface(in.GetDnsDomains()))
	}

	if v, ok := in.GetRoutesOk(); ok {
		routes := make([]map[string]interface{}, 0)
		for _, r := range v {
			route := make(map[string]interface{})
			if v, ok := r.GetAddressOk(); ok {
				route["address"] = *v
			}
			if v, ok := r.GetNetmaskOk(); ok {
				route["netmask"] = *v
			}
			if v, ok := r.GetGatewayOk(); ok {
				route["gateway"] = *v
			}
			if v, ok := r.GetNicOk(); ok && *v != "" {
				route["nic"] = *v
			}
			routes = append(routes, route)
		}
		networking["routes"] = routes

	}
	networkings = append(networkings, networking)

	return networkings, nil
}

func resourceAppgateApplianceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating Appliance: %s", d.Get("name").(string))
	// var diags diag.Diagnostics
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.AppliancesApi
	currentVersion := meta.(*Client).ApplianceVersion
	request := api.AppliancesIdGet(ctx, d.Id())
	originalAppliance, _, err := request.Authorization(token).Execute()
	if err != nil {
		return diag.Errorf("Failed to read Appliance, %s", err)
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
		_, v := d.GetChange("client_interface")
		cinterface, err := readClientInterfaceFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		originalAppliance.SetClientInterface(cinterface)
	}

	if d.HasChange("peer_interface") {
		_, v := d.GetChange("peer_interface")
		pinterface, err := readPeerInterfaceFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		originalAppliance.SetPeerInterface(pinterface)
	}

	if d.HasChange("admin_interface") {
		_, v := d.GetChange("admin_interface")
		ainterface, err := readAdminInterfaceFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		// since admin_interface is Optional, but admin_interface.hostname is required
		// if it set, we will make sure that hostname is not None before we send the request
		// to avoid sending empty fields in the request body.
		// otherwise, admin_interface has been removed, and we can set admin_interface to nil
		// and it will be omitted from the PUT request.
		if v, ok := ainterface.GetHostnameOk(); ok && v != nil && len(*v) > 0 {
			originalAppliance.SetAdminInterface(ainterface)
		} else {
			originalAppliance.AdminInterface = nil
		}
	}

	if d.HasChange("networking") {
		_, v := d.GetChange("networking")
		networking, err := readNetworkingFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		originalAppliance.SetNetworking(networking)
	}

	if d.HasChange("ntp") {
		_, v := d.GetChange("ntp")
		ntp, err := readNTPFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		originalAppliance.SetNtp(ntp)
	}

	if d.HasChange("ssh_server") {
		_, v := d.GetChange("ssh_server")
		sshServer, err := readSSHServerFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		originalAppliance.SetSshServer(sshServer)
	}

	if d.HasChange("snmp_server") {
		_, v := d.GetChange("snmp_server")
		srv, err := readSNMPServerFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		originalAppliance.SetSnmpServer(srv)
	}

	if d.HasChange("healthcheck_server") {
		_, v := d.GetChange("healthcheck_server")
		srv, err := readHealthcheckServerFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		originalAppliance.SetHealthcheckServer(srv)
	}

	if d.HasChange("prometheus_exporter") {
		_, v := d.GetChange("prometheus_exporter")
		exporter, err := readPrometheusExporterFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		originalAppliance.SetPrometheusExporter(exporter)
	}

	if d.HasChange("ping") {
		_, v := d.GetChange("ping")
		p, err := readPingFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		originalAppliance.SetPing(p)
	}

	if d.HasChange("log_server") {
		_, v := d.GetChange("log_server")
		logSrv, err := readLogServerFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		// we will only include log server in the request body if its enabled,
		// otherwise we will omit it in the request body and let the controller compute the rest.
		if logSrv.Enabled != nil && *logSrv.Enabled {
			originalAppliance.SetLogServer(logSrv)
		} else {
			originalAppliance.LogServer = nil
		}

	}

	if d.HasChange("controller") {
		_, v := d.GetChange("controller")
		ctrl, err := readControllerFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		originalAppliance.SetController(ctrl)
	}

	if d.HasChange("gateway") {
		_, v := d.GetChange("gateway")
		gw, err := readGatewayFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		originalAppliance.SetGateway(gw)
	}

	if d.HasChange("log_forwarder") {
		_, v := d.GetChange("log_forwarder")
		lf, err := readLogForwardFromConfig(v.([]interface{}), currentVersion)
		if err != nil {
			return diag.FromErr(err)
		}
		originalAppliance.SetLogForwarder(lf)
	}

	if d.HasChange("connector") {
		_, v := d.GetChange("connector")
		iot, err := readApplianceConnectorFromConfig(currentVersion, v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		originalAppliance.SetConnector(iot)
	}

	if d.HasChange("portal") {
		_, v := d.GetChange("portal")
		portal, err := readAppliancePortalFromConfig(d, v.([]interface{}), currentVersion)
		if err != nil {
			return diag.FromErr(err)
		}
		originalAppliance.SetPortal(portal)
	}

	if d.HasChange("rsyslog_destinations") {
		_, v := d.GetChange("rsyslog_destinations")
		rsys, err := readRsyslogDestinationFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		originalAppliance.SetRsyslogDestinations(rsys)
	}

	if d.HasChange("hostname_aliases") {
		_, v := d.GetChange("hostname_aliases")
		hostnames, err := readHostnameAliasesFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		originalAppliance.SetHostnameAliases(hostnames)
	}

	req := api.AppliancesIdPut(ctx, d.Id())

	_, _, err = req.Appliance(*originalAppliance).Authorization(token).Execute()
	if err != nil {
		return diag.Errorf("Could not update appliance %s", prettyPrintAPIError(err))
	}
	return resourceAppgateApplianceRead(ctx, d, meta)
}

func resourceAppgateApplianceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete Appliance: %s", d.Get("name").(string))
	var diags diag.Diagnostics

	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.AppliancesApi

	// Get appliance
	request := api.AppliancesIdGet(ctx, d.Id())
	appliance, _, err := request.Authorization(token).Execute()
	if err != nil {
		return diag.Errorf("Failed to delete Appliance while GET, %s", err)
	}
	// Deactivate
	if ok, _ := appliance.GetActivatedOk(); *ok {
		log.Printf("[DEBUG] Appliance is active, deactivate and wiping before deleting")
		deactiveRequest := api.AppliancesIdDeactivatePost(ctx, appliance.GetId())
		_, err = deactiveRequest.Wipe(true).Authorization(token).Execute()
		if err != nil {
			return diag.Errorf("Failed to delete Appliance while deactivating, %s", err)
		}
	}

	// Delete
	deleteRequest := api.AppliancesIdDelete(ctx, appliance.GetId())
	_, err = deleteRequest.Authorization(token).Execute()
	if err != nil {
		return diag.Errorf("Failed to delete Appliance, %s", err)
	}
	d.SetId("")
	return diags
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
		if v := raw["allow_sources"].([]interface{}); len(v) > 0 {
			allowSources, err := readAllowSources(v)
			if err != nil {
				return cinterface, fmt.Errorf("Failed to resolve client_interface.allow_sources: %w", err)
			}
			cinterface.SetAllowSources(allowSources)
		}
		if v, ok := raw["override_spa_mode"].(string); ok && len(v) > 0 {
			if v != "Disabled" {
				cinterface.SetOverrideSpaMode(v)
			}
		}
	}
	return cinterface, nil
}

func readPeerInterfaceFromConfig(pinterfaces []interface{}) (openapi.ApplianceAllOfPeerInterface, error) {
	pinterf := openapi.ApplianceAllOfPeerInterface{}
	for _, r := range pinterfaces {
		raw := r.(map[string]interface{})
		if v, ok := raw["hostname"].(string); ok && len(v) > 0 {
			pinterf.SetHostname(v)
		}
		if v, ok := raw["https_port"]; ok {
			pinterf.SetHttpsPort(int32(v.(int)))
		}
		if v := raw["allow_sources"].([]interface{}); len(v) > 0 {
			allowSources, err := readAllowSources(v)
			if err != nil {
				return pinterf, fmt.Errorf("Failed to resolve peer_interface.allow_sources: %w", err)
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
		if v, ok := raw["hostname"].(string); ok && len(v) > 0 {
			aInterface.SetHostname(v)
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

		if v := raw["allow_sources"].([]interface{}); len(v) > 0 {
			allowSources, err := readAllowSources(v)
			if err != nil {
				return aInterface, fmt.Errorf("Failed to admin interface allowed sources: %w", err)
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
				return network, fmt.Errorf("Failed to resolve network hosts: %w", err)
			}
			network.SetHosts(hosts)
		}

		if v := rawNetwork["nics"]; len(v.([]interface{})) > 0 {
			nics, err := readNetworkNicsFromConfig(v.([]interface{}))
			if err != nil {
				return network, fmt.Errorf("Failed to resolve network nics: %w", err)
			}
			network.SetNics(nics)
		}
		dnsServers := make([]string, 0)
		if v, ok := rawNetwork["dns_servers"]; ok {
			list := v.(*schema.Set).List()
			for _, dns := range list {
				dnsServers = append(dnsServers, dns.(string))
			}
		}
		if len(dnsServers) > 0 {
			network.SetDnsServers(dnsServers)
		}

		dnsDomains := make([]string, 0)
		if v, ok := rawNetwork["dns_domains"]; ok {
			list := v.(*schema.Set).List()
			for _, dns := range list {
				dnsDomains = append(dnsDomains, dns.(string))
			}
		}
		if len(dnsDomains) > 0 {
			network.SetDnsDomains(dnsDomains)
		}

		if v := rawNetwork["routes"]; len(v.([]interface{})) > 0 {
			routes, err := readNetworkRoutesFromConfig(v.([]interface{}))
			if err != nil {
				return network, fmt.Errorf("Failed to resolve network routes: %w", err)
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
		if v := rawServer["allow_sources"].([]interface{}); len(v) > 0 {
			allowSources, err := readAllowSources(v)
			if err != nil {
				return sshServer, err
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

		if v := rawServer["allow_sources"].([]interface{}); len(v) > 0 {
			allowSources, err := readAllowSources(v)
			if err != nil {
				return server, err
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
		if v := rawServer["allow_sources"].([]interface{}); len(v) > 0 {
			allowSources, err := readAllowSources(v)
			if err != nil {
				return server, err
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
				return ntpCfg, fmt.Errorf("Failed to resolve ntp servers: %w", err)
			}
			ntpCfg.SetServers(ntpServers)
		}
	}
	return ntpCfg, nil
}

func readLogServerFromConfig(logServers []interface{}) (openapi.ApplianceAllOfLogServer, error) {
	srv := openapi.ApplianceAllOfLogServer{}
	for _, raw := range logServers {
		r := raw.(map[string]interface{})
		if v, ok := r["enabled"].(bool); ok {
			srv.SetEnabled(v)
		}
		if v, ok := r["retention_days"]; ok {
			srv.SetRetentionDays(int32(v.(int)))
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
						if v, ok := raw["netmask"].(int); ok && v >= 0 {
							ad.SetNetmask(int32(v))
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
		if v := rawServer["allow_sources"].([]interface{}); len(v) > 0 {
			allowSources, err := readAllowSources(v)
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
		if v := rawServer["allow_sources"].([]interface{}); len(v) > 0 {
			allowSources, err := readAllowSources(v)
			if err != nil {
				return val, err
			}
			val.SetAllowSources(allowSources)
		}
	}
	return val, nil
}

func readLogForwardFromConfig(logforwards []interface{}, currentVersion *version.Version) (openapi.ApplianceAllOfLogForwarder, error) {
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
				if v, ok := r["compatibility_mode"]; ok {
					if currentVersion.LessThan(Appliance55Version) {
						return val, fmt.Errorf("elasticsearch.compatibility_mode is only available in 5.5 or greater, got %s", currentVersion)
					}
					elasticsearch.SetCompatibilityMode(int32(v.(int)))
				}

				if v, ok := r["authentication"].([]interface{}); ok {
					if currentVersion.LessThan(Appliance55Version) {
						return val, fmt.Errorf("elasticsearch.authentication is only available in 5.5 or greater, got %s", currentVersion)
					}
					val := v[0].(map[string]interface{})
					a := openapi.ElasticsearchAllOfAuthentication{}
					if v, ok := val["type"].(string); ok && len(v) > 0 {
						a.SetType(v)
					}
					if v, ok := val["token"].(string); ok && len(v) > 0 {
						a.SetToken(v)
					}
					if len(a.GetType()) > 0 {
						elasticsearch.SetAuthentication(a)
					}
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

func readApplianceConnectorFromConfig(currentVersion *version.Version, connectors []interface{}) (openapi.ApplianceAllOfConnector, error) {
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
				if v := r["allow_resources"].([]interface{}); len(v) > 0 {
					allowedSources, err := listToMapList(v)
					if err != nil {
						return val, err
					}
					sources := make([]openapi.AllowResourcesInner, 0)
					for _, as := range allowedSources {
						row := openapi.NewAllowResourcesInnerWithDefaults()
						if v, ok := as["address"].(string); ok {
							row.SetAddress(v)
						}
						if v, ok := as["netmask"].(int); ok {
							row.SetNetmask(int32(v))
						}
						sources = append(sources, *row)
					}
					client.SetAllowResources(sources)
				}
				if v, ok := r["snat_to_resources"]; ok {
					client.SetSnatToResources(v.(bool))
				}
				if currentVersion.GreaterThanOrEqual(Appliance54Version) {
					if v, ok := r["dnat_to_resource"]; ok {
						client.SetDnatToResource(v.(bool))
					}
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
					allowedSources, err := listToMapList(v.([]interface{}))
					if err != nil {
						return val, err
					}
					sources := make([]openapi.AllowSourcesInner, 0)
					for _, as := range allowedSources {
						row := openapi.NewAllowSourcesInnerWithDefaults()
						if v, ok := as["address"].(string); ok {
							row.SetAddress(v)
						}
						if v, ok := as["netmask"].(int); ok {
							row.SetNetmask(int32(v))
						}
						if v, ok := as["nic"].(string); ok {
							row.SetNic(v)
						}
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

func readAppliancePortalFromConfig(d *schema.ResourceData, portals []interface{}, currentVersion *version.Version) (openapi.Portal, error) {
	p := openapi.Portal{}
	for _, portal := range portals {
		if portal == nil {
			continue
		}

		raw := portal.(map[string]interface{})
		if v, ok := raw["enabled"]; ok {
			p.SetEnabled(v.(bool))
		}
		if v, ok := raw["profiles"]; ok {
			profiles, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return p, err
			}
			p.SetProfiles(profiles)
		}

		if v, ok := raw["https_p12"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			p12 := openapi.P12{}
			raw := v[0].(map[string]interface{})
			p12.SetId(uuid.New().String())
			if v, ok := raw["content"]; ok {
				content, err := appliancePortalReadp12Content(v.(string))
				if err != nil {
					return p, fmt.Errorf("unable to read https_p12 file content %w", err)
				}
				p12.SetContent(content)
			}
			if v, ok := raw["password"]; ok {
				p12.SetPassword(v.(string))
			}
			p.SetHttpsP12(p12)
		}

		if v, ok := raw["proxy_p12s"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			p12s := make([]openapi.Portal12, 0)
			for _, k := range v {
				raw := k.(map[string]interface{})
				proxyp12 := openapi.Portal12{}
				id := uuid.New().String()
				proxyp12.SetId(id)
				proxyp12.Id = &id
				if v, ok := raw["content"]; ok {
					content, err := appliancePortalReadp12Content(v.(string))
					if err != nil {
						return p, fmt.Errorf("unable to read p12 file content %w", err)
					}
					proxyp12.SetContent(content)
				}
				if v, ok := raw["password"]; ok {
					proxyp12.SetPassword(v.(string))
				}
				p12s = append(p12s, proxyp12)
			}

			p.SetProxyP12s(p12s)
		}
		if v, ok := raw["external_profiles"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			profiles := make([]openapi.PortalExternalProfilesInner, 0)
			for _, k := range v {
				raw := k.(map[string]interface{})
				profile := openapi.PortalExternalProfilesInner{}
				if v, ok := raw["id"]; ok {
					profile.SetId(v.(string))
				}
				if v, ok := raw["url"]; ok {
					profile.SetUrl(v.(string))
				}
				profiles = append(profiles, profile)
			}
			p.SetExternalProfiles(profiles)
		}
		if v, ok := raw["sign_in_customization"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			customization := openapi.PortalSignInCustomization{}
			raw := v[0].(map[string]interface{})

			if v, ok := raw["background_color"].(string); ok && len(v) > 0 {
				customization.SetBackgroundColor(v)
			}
			if _, ok := raw["background_image"]; ok {
				k := "portal.0.sign_in_customization.0.background_image"
				content, err := getResourceFileContent(d, k)
				if err != nil {
					return p, err
				}
				customization.SetBackgroundImage(base64.StdEncoding.EncodeToString(content))
			}
			if v, ok := raw["logo"].(string); ok && len(v) > 0 {
				k := "portal.0.sign_in_customization.0.logo"
				content, err := getResourceFileContent(d, k)
				if err != nil {
					return p, err
				}
				customization.SetLogo(base64.StdEncoding.EncodeToString(content))
			}
			if v, ok := raw["text"].(string); ok && len(v) > 0 {
				customization.SetText(v)
			}
			if v, ok := raw["text_color"].(string); ok && len(v) > 0 {
				customization.SetTextColor(v)
			}

			if v, ok := raw["auto_redirect"].(bool); ok {
				if currentVersion.LessThan(Appliance60Version) && v {
					return p, fmt.Errorf("portal.sign_in_customization.auto_redirect is not allowed in %s", currentVersion.String())
				} else if currentVersion.GreaterThanOrEqual(Appliance60Version) {
					customization.SetAutoRedirect(v)
				}
			}

			p.SetSignInCustomization(customization)
		}

	}
	return p, nil
}

func appliancePortalReadp12Content(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("Error opening p12 file (%s): %w", path, err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("[WARN] Error closing p12 file (%s): %s", path, err)
		}
	}()
	reader := bufio.NewReader(file)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("Error reading file (%s): %w", path, err)
	}
	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded, nil
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
