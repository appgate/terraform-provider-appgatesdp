package appgate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/appgate/sdp-api-client-go/api/v19/openapi"
	"github.com/appgate/terraform-provider-appgatesdp/appgate/hashcode"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	identityProviderLocalDatabase   = "LocalDatabase"
	identityProviderRadius          = "Radius"
	identityProviderLdap            = "Ldap"
	identityProviderSaml            = "Saml"
	identityProviderLdapCertificate = "LdapCertificate"
	identityProviderConnector       = "Connector"
	identityProviderOidc            = "Oidc"
	builtinProviderLocal            = "local"
	builtinProviderConnector        = "Connector"
)

var (
	ErrNetworkInactivityTimeoutEnabled = errors.New("network_inactivity_timeout_enabled is only available in 6.1 or higher")
)

func identityProviderSchema() map[string]*schema.Schema {
	return mergeSchemaMaps(baseEntitySchema(), identityProviderIPPoolSchema(), identityProviderClaimsSchema(), func() map[string]*schema.Schema {
		ip := map[string]*schema.Schema{
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
						identityProviderConnector,
						identityProviderOidc,
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

			"admin_provider": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"device_limit_per_user": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},

			"on_boarding_two_factor": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mfa_provider_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"message": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"device_limit_per_user": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"claim_suffix": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "onBoarding",
						},
						"always_required": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"inactivity_timeout_minutes": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"network_inactivity_timeout_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"user_scripts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dns_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dns_search_domains": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"block_local_dns_requests": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
		}
		return ip
	}())
}

func identityProviderIPPoolSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ip_pool_v4": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ip_pool_v6": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func identityProviderClaimsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"claim_mappings": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Set:      resourceIdentityProviderClaimMappingsHash,
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
						Type:             schema.TypeBool,
						Computed:         true,
						Optional:         true,
						DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
					},
					"encrypt": {
						Type:             schema.TypeBool,
						Computed:         true,
						Optional:         true,
						DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
					},
				},
			},
		},
		"on_demand_claim_mappings": {
			Type:     schema.TypeSet,
			Optional: true,
			Set:      resourceIdentityProviderOnDemandClaimMappingsHash,
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
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     schema.TypeString,
									Computed: true,
									Optional: true,
								},
								"path": {
									Type:     schema.TypeString,
									Computed: true,
									Optional: true,
								},
								"args": {
									Type:     schema.TypeString,
									Computed: true,
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
	}
}

func resourceIdentityProviderClaimMappingsHash(v interface{}) int {
	raw, ok := v.(map[string]interface{})
	if !ok {
		return 0
	}
	// modifying raw actually modifies the values passed to the provider.
	// Use a copy to avoid that.
	copy := make((map[string]interface{}))
	for key, value := range raw {
		copy[key] = value
	}
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", copy["attribute_name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", copy["claim_name"].(string)))
	buf.WriteString(fmt.Sprintf("%t-", copy["list"].(bool)))
	buf.WriteString(fmt.Sprintf("%t-", copy["encrypt"].(bool)))

	h := hashcode.String(buf.String())

	return h
}

func resourceIdentityProviderOnDemandClaimMappingsHash(v interface{}) int {
	raw := v.(map[string]interface{})
	// modifying raw actually modifies the values passed to the provider.
	// Use a copy to avoid that.
	copy := make((map[string]interface{}))
	for key, value := range raw {
		copy[key] = value
	}
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", copy["command"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", copy["claim_name"].(string)))

	if v, ok := copy["parameters"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		phash := func(i interface{}) int {
			var buf bytes.Buffer
			m, ok := i.(map[string]interface{})

			if !ok {
				return 0
			}
			if v, ok := m["name"]; ok {
				buf.WriteString(fmt.Sprintf("%s-", v.(string)))
			}
			if v, ok := m["path"]; ok {
				buf.WriteString(fmt.Sprintf("%s-", v.(string)))
			}
			if v, ok := m["args"]; ok {
				buf.WriteString(fmt.Sprintf("%s-", v.(string)))
			}

			return hashcode.String(buf.String())
		}
		mHash := phash(v[0])
		buf.WriteString(fmt.Sprintf("%d-", mHash))
	}
	buf.WriteString(fmt.Sprintf("%s-", copy["platform"].(string)))

	return hashcode.String(buf.String())
}

// ldapProviderSchema return the default base schema for
// LDAP and LDAP Certificate provider.
func ldapProviderSchema() map[string]*schema.Schema {
	s := identityProviderSchema()

	s["hostnames"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
	s["port"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}
	s["ssl_enabled"] = &schema.Schema{
		Type:     schema.TypeBool,
		Computed: true,
		Optional: true,
	}
	s["admin_distinguished_name"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	s["admin_password"] = &schema.Schema{
		Type:      schema.TypeString,
		Sensitive: true,
		Optional:  true,
	}
	s["base_dn"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	s["object_class"] = &schema.Schema{
		Type:       schema.TypeString,
		Deprecated: "Deprecated as of 6.2. Use userFilter instead",
		Computed:   true,
		Optional:   true,
	}
	s["user_filter"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	s["username_attribute"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Optional: true,
	}
	s["membership_filter"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Optional: true,
	}
	s["membership_base_dn"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	s["password_warning"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:     schema.TypeBool,
					Computed: true,
					Optional: true,
				},
				"threshold_days": {
					Type:     schema.TypeInt,
					Computed: true,
					Optional: true,
				},
				"message": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
	return s
}

// readProviderFromConfig reads all the common attributes for the IdentityProviders.
func readProviderFromConfig(d *schema.ResourceData, provider openapi.ConfigurableIdentityProvider, currentVersion *version.Version) (*openapi.ConfigurableIdentityProvider, error) {
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

	if v, ok := d.GetOk("admin_provider"); ok {
		provider.SetAdminProvider(v.(bool))
	}

	// device_limit_per_user is only available on 5.5 or higher on root level,
	// previous version has this on on_boarding_two_factor.device_limit_per_user
	if v, ok := d.GetOk("device_limit_per_user"); ok {
		if currentVersion.LessThan(Appliance55Version) {
			return &provider, fmt.Errorf(
				"device_limit_per_user is only available on 5.5, your current version is %s, Use on_boarding_two_factor.device_limit_per_user for appliances less then 5.5",
				currentVersion.String(),
			)
		}
		provider.SetDeviceLimitPerUser(int32(v.(int)))
	}

	if v, ok := d.GetOk("on_boarding_two_factor"); ok {
		onboarding, err := readOnBoardingTwoFactorFromConfig(v.([]interface{}), currentVersion)
		if err != nil {
			return &provider, err
		}
		provider.SetOnBoarding2FA(onboarding)
	}

	if v, ok := d.GetOk("inactivity_timeout_minutes"); ok {
		provider.SetInactivityTimeoutMinutes(int32(v.(int)))
	}
	if v, ok := d.GetOk("network_inactivity_timeout_enabled"); ok {
		provider.SetNetworkInactivityTimeoutEnabled(v.(bool))
	}
	if v, ok := d.GetOk("ip_pool_v4"); ok {
		provider.SetIpPoolV4(v.(string))
	}
	if v, ok := d.GetOk("ip_pool_v6"); ok {
		provider.SetIpPoolV6(v.(string))
	}
	if v, ok := d.GetOk("user_scripts"); ok {
		us, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return &provider, err
		}
		provider.SetUserScripts(us)
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
		claims := readIdentityProviderClaimMappingFromConfig(v.(*schema.Set).List())
		if len(claims) > 0 {
			provider.SetClaimMappings(claims)
		}
	}
	if v, ok := d.GetOk("on_demand_claim_mappings"); ok {
		claims := readIdentityProviderOnDemandClaimMappingFromConfig(v.(*schema.Set).List())
		if len(claims) > 0 {
			provider.SetOnDemandClaimMappings(claims)
		}
	}
	return &provider, nil
}

func readOnBoardingTwoFactorFromConfig(input []interface{}, currentVersion *version.Version) (openapi.ConfigurableIdentityProviderAllOfOnBoarding2FA, error) {
	onboarding := openapi.ConfigurableIdentityProviderAllOfOnBoarding2FA{}
	for _, r := range input {
		raw := r.(map[string]interface{})
		if v, ok := raw["mfa_provider_id"]; ok {
			onboarding.SetMfaProviderId(v.(string))
		}
		if v, ok := raw["message"]; ok {
			onboarding.SetMessage(v.(string))
		}
		if v, ok := raw["device_limit_per_user"]; ok {
			val := int32(v.(int))
			log.Printf("[DEBUG] on_boarding_two_factor.device_limit_per_user only available in version 5.4 or less got %v - %v", val, currentVersion.LessThan(Appliance55Version))
			if currentVersion.LessThan(Appliance55Version) && val > 0 {
				onboarding.SetDeviceLimitPerUser(val)
			} else if val > 0 {
				// device_limit_per_user is not allowed in 5.5
				return onboarding, fmt.Errorf(
					"on_boarding_two_factor.device_limit_per_user is deprecated in %s. Use root level field instead. Got %d",
					currentVersion.String(),
					val,
				)
			} else {
				// else omit devicelmit per user from the request.
				log.Printf("[DEBUG] on_boarding_two_factor.device_limit_per_user is not allowed on %s, omitted it from request, use root level instead", currentVersion.String())
				onboarding.DeviceLimitPerUser = nil
			}
		}

		if v, ok := raw["claim_suffix"]; ok {
			onboarding.SetClaimSuffix(v.(string))
		}
		if v, ok := raw["always_required"]; ok {
			onboarding.SetAlwaysRequired(v.(bool))
		}
	}
	return onboarding, nil
}

// func readIdentityProviderClaimMappingFromConfig(input []interface{}) []map[string]interface{} {
func readIdentityProviderClaimMappingFromConfig(input []interface{}) []openapi.ClaimMappingsInner {
	claims := make([]openapi.ClaimMappingsInner, 0)
	for _, raw := range input {
		claim := raw.(map[string]interface{})
		c := openapi.ClaimMappingsInner{}
		if v, ok := claim["attribute_name"]; ok {
			c.SetAttributeName(v.(string))
		}
		if v, ok := claim["claim_name"]; ok {
			c.SetClaimName(v.(string))
		}
		c.SetList(claim["list"].(bool))
		c.SetEncrypt(claim["encrypt"].(bool))
		claims = append(claims, c)
	}
	return claims
}

func readIdentityProviderOnDemandClaimMappingFromConfig(input []interface{}) []openapi.OnDemandClaimMappingsInner {
	claims := make([]openapi.OnDemandClaimMappingsInner, 0)
	for _, raw := range input {
		claim := raw.(map[string]interface{})
		c := openapi.NewOnDemandClaimMappingsInnerWithDefaults()
		if v, ok := claim["command"]; ok {
			c.SetCommand(v.(string))
		}
		if v, ok := claim["claim_name"]; ok {
			c.SetClaimName(v.(string))
		}
		if v, ok := claim["platform"]; ok {
			c.SetPlatform(v.(string))
		}
		if v, ok := claim["parameters"]; ok {
			p := openapi.NewOnDemandClaimMappingsInnerParametersWithDefaults()
			for _, para := range v.([]interface{}) {
				parameters := para.(map[string]interface{})
				if v, ok := parameters["name"]; ok && len(v.(string)) > 0 {
					p.SetName(v.(string))
				}
				if v, ok := parameters["path"]; ok && len(v.(string)) > 0 {
					p.SetPath(v.(string))
				}
				if v, ok := parameters["args"]; ok && len(v.(string)) > 0 {
					p.SetArgs(v.(string))
				}
				c.SetParameters(*p)
			}
		}

		claims = append(claims, *c)
	}
	return claims
}

func identityProviderDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete LdapProvider: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.IdentityProvidersApi

	if _, err := api.IdentityProvidersIdDelete(context.Background(), d.Id()).Authorization(token).Execute(); err != nil {
		return fmt.Errorf("Could not delete LdapProvider %w", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}

func flattenIdentityProviderClaimsMappning(claims []openapi.ClaimMappingsInner) *schema.Set {
	out := make([]interface{}, 0)
	for _, claim := range claims {
		row := make(map[string]interface{})
		row["attribute_name"] = claim.GetAttributeName()
		row["claim_name"] = claim.GetClaimName()
		row["list"] = claim.GetList()
		row["encrypt"] = claim.GetEncrypt()

		out = append(out, row)
	}
	return schema.NewSet(resourceIdentityProviderClaimMappingsHash, out)
}

func flattenIdentityProviderOnDemandClaimsMappning(claims []openapi.OnDemandClaimMappingsInner) *schema.Set {
	out := []interface{}{}
	for _, claim := range claims {
		row := make(map[string]interface{})
		if v, ok := claim.GetCommandOk(); ok {
			row["command"] = *v
		}
		if v, ok := claim.GetClaimNameOk(); ok {
			row["claim_name"] = *v
		}
		if v, ok := claim.GetParametersOk(); ok {
			parameters := make([]map[string]interface{}, 0)
			parameter := make(map[string]interface{})
			if v, ok := v.GetNameOk(); ok && len(*v) > 0 {
				parameter["name"] = v
			}
			if v, ok := v.GetPathOk(); ok && len(*v) > 0 {
				parameter["path"] = v
			}
			if v, ok := v.GetArgsOk(); ok && len(*v) > 0 {
				parameter["args"] = v
			}
			parameters = append(parameters, parameter)
			row["parameters"] = parameters
		}
		if v, ok := claim.GetPlatformOk(); ok {
			row["platform"] = *v
		}
		out = append(out, row)
	}
	return schema.NewSet(resourceIdentityProviderOnDemandClaimMappingsHash, out)
}

func flattenIdentityProviderOnboarding2fa(input openapi.ConfigurableIdentityProviderAllOfOnBoarding2FA, currentVersion *version.Version) []interface{} {
	o := make(map[string]interface{})
	if v, ok := input.GetMfaProviderIdOk(); ok {
		o["mfa_provider_id"] = v
	}
	if v, ok := input.GetMessageOk(); ok {
		o["message"] = v
	}
	// we will only save device_limit_per_user in the statefile if the currentversion still supports it.
	if currentVersion.LessThan(Appliance55Version) {
		if v, ok := input.GetDeviceLimitPerUserOk(); ok {
			o["device_limit_per_user"] = int(*v)
		}
	}
	if v, ok := input.GetClaimSuffixOk(); ok {
		o["claim_suffix"] = v
	}
	if v, ok := input.GetAlwaysRequiredOk(); ok {
		o["always_required"] = v
	}

	return []interface{}{o}
}
