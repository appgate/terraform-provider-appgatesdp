package appgate

import (
	"context"
	"fmt"
	"log"

	"github.com/appgate/terraform-provider-appgate/client/v12/openapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	identityProviderLocalDatabase   = "LocalDatabase"
	identityProviderRadius          = "Radius"
	identityProviderLdap            = "Ldap"
	identityProviderSaml            = "Saml"
	identityProviderLdapCertificate = "LdapCertificate"
	identityProviderIotConnector    = "IotConnector"
)

func identityProviderSchema() map[string]*schema.Schema {
	return mergeSchemaMaps(baseEntitySchema(), map[string]*schema.Schema{
		"type": {
			Optional: true,
			Type:     schema.TypeString,
			ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
				s := v.(string)
				list := []string{
					identityProviderLocalDatabase,
					identityProviderRadius,
					identityProviderLdap,
					identityProviderSaml,
					identityProviderLdapCertificate,
					identityProviderIotConnector,
				}
				for _, x := range list {
					if s == x {
						return
					}
				}
				errs = append(errs, fmt.Errorf("type must be on of %v, got %s", list, s))
				return
			},
		},

		"default": {
			Type:     schema.TypeBool,
			Optional: true,
		},

		"admin_provider": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"on_boarding_2fa": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"mfa_provider_id": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"message": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"device_limit_per_user": {
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"inactivity_timeout_minutes": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"ip_pool_v4": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ip_pool_v6": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"dns_servers": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"dns_search_domains": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"block_local_dns_requests": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"claim_mappings": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"attribute_name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"claim_name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"list": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"encrypted": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"on_demand_claim_mappings": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"command": {
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
							s := v.(string)
							list := []string{
								"fileSize",
								"fileExists",
								"fileCreated",
								"fileUpdated",
								"fileVersion",
								"fileSha512",
								"processRunning",
								"processList",
								"serviceRunning",
								"serviceList",
								"regExists",
								"regQuery",
								"runScript",
							}
							for _, x := range list {
								if s == x {
									return
								}
							}
							errs = append(errs, fmt.Errorf("type must be on of %v, got %s", list, s))
							return
						},
					},
					"claim_name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"parameters": {
						Type:     schema.TypeMap,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"path": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"args": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
					"platform": {
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
							s := v.(string)
							list := []string{
								"desktop.windows.all",
								"desktop.macos.all",
								"desktop.linux.all",
								"desktop.all",
								"mobile.android.all",
								"mobile.ios.all",
								"mobile.all",
								"all",
							}
							for _, x := range list {
								if s == x {
									return
								}
							}
							errs = append(errs, fmt.Errorf("type must be on of %v, got %s", list, s))
							return
						},
					},
				},
			},
		},
	})
}

// readProviderFromConfig reads all the common attribudes for the IdentityProviders.
func readProviderFromConfig(d *schema.ResourceData, provider openapi.IdentityProvider) (*openapi.IdentityProvider, error) {
	base, err := readBaseEntityFromConfig(d)
	if err != nil {
		return &provider, err
	}
	if _, o := base.GetNameOk(); o {
		provider.SetName(base.GetName())
	}
	if _, o := base.GetTagsOk(); o {
		provider.SetTags(base.GetTags())
	}
	if _, o := base.GetNotesOk(); o {
		provider.SetNotes(base.GetNotes())
	}
	if v, ok := d.GetOk("display_name"); ok {
		provider.SetDisplayName(v.(string))
	}
	if v, ok := d.GetOk("default"); ok {
		provider.SetDefault(v.(bool))
	}
	if v, ok := d.GetOk("client_provider"); ok {
		provider.SetClientProvider(v.(bool))
	}
	if v, ok := d.GetOk("admin_provider"); ok {
		provider.SetAdminProvider(v.(bool))
	}
	if v, ok := d.GetOk("on_boarding_2fa"); ok {
		onboarding := openapi.IdentityProviderAllOfOnBoarding2FA{}
		raw := v.(map[string]interface{})
		if v, ok := raw["mfa_provider_id"]; ok {
			onboarding.SetMfaProviderId(v.(string))
		}
		if v, ok := raw["message"]; ok {
			onboarding.SetMessage(v.(string))
		}
		if v, ok := raw["device_limit_per_user"]; ok {
			onboarding.SetDeviceLimitPerUser(int32(v.(int)))
		}
		provider.SetOnBoarding2FA(onboarding)
	}
	if v, ok := d.GetOk("on_boarding_type"); ok {
		provider.SetOnBoardingType(v.(string))
	}
	if v, ok := d.GetOk("on_boarding_otp_provider"); ok {
		provider.SetOnBoardingOtpProvider(v.(string))
	}
	if v, ok := d.GetOk("on_boarding_otp_message"); ok {
		provider.SetOnBoardingOtpMessage(v.(string))
	}
	if v, ok := d.GetOk("inactivity_timeout_minutes"); ok {
		provider.SetInactivityTimeoutMinutes(int32(v.(int)))
	}
	if v, ok := d.GetOk("ip_pool_v4"); ok {
		provider.SetIpPoolV4(v.(string))
	}
	if v, ok := d.GetOk("ip_pool_v6"); ok {
		provider.SetIpPoolV6(v.(string))
	}
	if v, ok := d.GetOk("dns_servers"); ok {
		servers, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return &provider, err
		}
		provider.SetDnsServers(servers)
	}
	if v, ok := d.GetOk("dns_search_domains"); ok {
		servers, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return &provider, err
		}
		provider.SetDnsSearchDomains(servers)
	}
	if v, ok := d.GetOk("block_local_dns_requests"); ok {
		provider.SetBlockLocalDnsRequests(v.(bool))
	}
	if v, ok := d.GetOk("claim_mappings"); ok {
		claims := make([]map[string]interface{}, 0)
		for _, claim := range v.([]map[string]interface{}) {
			c := make(map[string]interface{})
			if v, ok := claim["attribute_name"]; ok {
				c["attribute_name"] = v.(string)
			}
			if v, ok := claim["claim_name"]; ok {
				c["claim_name"] = v.(string)
			}
			if v, ok := claim["list"]; ok {
				c["list"] = v.(bool)
			}
			if v, ok := claim["encrypt"]; ok {
				c["encrypt"] = v.(bool)
			}
			claims = append(claims, c)
		}
		if len(claims) > 0 {
			provider.SetClaimMappings(claims)
		}
	}
	if v, ok := d.GetOk("on_demand_claim_mappings"); ok {
		claims := make([]map[string]interface{}, 0)
		for _, claim := range v.([]map[string]interface{}) {
			c := make(map[string]interface{})
			if v, ok := claim["command"]; ok {
				c["command"] = v.(string)
			}
			if v, ok := claim["claim_name"]; ok {
				c["claim_name"] = v.(string)
			}
			if v, ok := claim["parameters"]; ok {
				parameters := v.(map[string]interface{})
				if v, ok := parameters["name"]; ok {
					c["name"] = v.(string)
				}
				if v, ok := parameters["path"]; ok {
					c["path"] = v.(string)
				}
				if v, ok := parameters["args"]; ok {
					c["args"] = v.(string)
				}
				c["parameters"] = parameters
			}
			if v, ok := claim["platform"]; ok {
				c["platform"] = v.(bool)
			}
			claims = append(claims, c)
		}
		if len(claims) > 0 {
			provider.SetOnDemandClaimMappings(claims)
		}
	}
	return &provider, nil
}

func identityProviderDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete LdapProvider: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.IdentityProvidersApi

	request := api.IdentityProvidersIdDelete(context.Background(), d.Id())

	_, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not delete LdapProvider %+v", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
