package appgate

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/appgate/sdp-api-client-go/api/v18/openapi"
	"github.com/appgate/terraform-provider-appgatesdp/appgate/hashcode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgatePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgatePolicyCreate,
		Read:   resourceAppgatePolicyRead,
		Update: resourceAppgatePolicyUpdate,
		Delete: resourceAppgatePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"policy_id": resourceUUID(),

			"name": {
				Type:        schema.TypeString,
				Description: "Name of the object.",
				Required:    true,
			},

			"notes": {
				Type:        schema.TypeString,
				Description: "Name of the object.",
				Default:     DefaultDescription,
				Optional:    true,
			},

			"tags": tagsSchema(),

			"disabled": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},

			"expression": {
				Type:     schema.TypeString,
				Required: true,
			},

			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Type of the Policy. It is informational and not enforced.",
			},

			"entitlements": {
				Type:             schema.TypeSet,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				Set:              schema.HashString,
				Elem:             &schema.Schema{Type: schema.TypeString},
			},

			"entitlement_links": {
				Type:             schema.TypeSet,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				Set:              schema.HashString,
				Elem:             &schema.Schema{Type: schema.TypeString},
			},

			"ringfence_rules": {
				Type:             schema.TypeSet,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				Set:              schema.HashString,
				Elem:             &schema.Schema{Type: schema.TypeString},
			},

			"ringfence_rule_links": {
				Type:             schema.TypeSet,
				Computed:         true,
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				Set:              schema.HashString,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
			},

			"tamper_proofing": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},

			"override_site": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"override_site_claim": {
				Type:        schema.TypeString,
				Description: "The path of a claim that contains the UUID of an override site. It should be defined as 'claims.xxx.xxx' or 'claims.xxx.xxx.xxx'1.",
				Optional:    true,
			},

			"proxy_auto_config": {
				Type:             schema.TypeList,
				MaxItems:         1,
				Optional:         true,
				Description:      "Client configures PAC URL on the client OS.",
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"persist": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"trusted_network_check": {
				Type:             schema.TypeList,
				MaxItems:         1,
				Optional:         true,
				Description:      "Client suspends operations when it's in a trusted network.",
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"dns_suffix": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"dns_settings": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of domain names with DNS server IPs that the Client should be using.",
				Set:         resourcePolicyDnsSettingsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"servers": {
							Type:     schema.TypeSet,
							Required: true,
							Set:      schema.HashString,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"client_settings": {
				Type:             schema.TypeList,
				MaxItems:         1,
				Optional:         true,
				Description:      "Settings that admins can apply to the Client.",
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"entitlements_list": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"attention_level": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"auto_start": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"add_remove_profiles": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"keep_me_signed_in": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"saml_auto_sign_in": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"quit": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"sign_out": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"suspend": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"client_profile_settings": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"profiles": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"custom_client_help_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"administrative_roles": {
				Type:             schema.TypeSet,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				Set:              schema.HashString,
				Elem:             &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourcePolicyDnsSettingsHash(v interface{}) int {
	var buf bytes.Buffer
	if v == nil {
		return hashcode.String(buf.String())
	}
	raw := v.(map[string]interface{})
	// modifying raw actually modifies the values passed to the provider.
	// Use a copy to avoid that.
	copy := make((map[string]interface{}))
	for key, value := range raw {
		copy[key] = value
	}

	buf.WriteString(fmt.Sprintf("%s-", copy["domain"].(string)))
	if v, ok := copy["servers"]; ok {
		buf.WriteString(fmt.Sprintf("%v-", v.(*schema.Set).List()))
	}

	return hashcode.String(buf.String())
}

func resourceAppgatePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating Policy with name: %s", d.Get("name").(string))
	ctx := context.Background()
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.PoliciesApi
	currentVersion := meta.(*Client).ApplianceVersion
	args := openapi.Policy{}

	if v, ok := d.GetOk("policy_id"); ok {
		args.SetId(v.(string))
	}

	// Type is only available in >= 5.5
	if currentVersion.LessThan(Appliance55Version) {
		args.Type = nil
	}

	args.SetName(d.Get("name").(string))

	if c, ok := d.GetOk("notes"); ok {
		args.SetNotes(c.(string))
	}

	args.SetTags(schemaExtractTags(d))

	if c, ok := d.GetOk("disabled"); ok {
		args.SetDisabled(c.(bool))
	}

	if c, ok := d.GetOk("expression"); ok {
		args.SetExpression(c.(string))
	}

	if currentVersion.GreaterThanOrEqual(Appliance54Version) {
		if v, ok := d.GetOk("client_settings"); ok {
			settings, err := readPolicyClientSettingsFromConfig(v.([]interface{}))
			if err != nil {
				return err
			}
			args.SetClientSettings(settings)
		}
	}
	if currentVersion.GreaterThanOrEqual(Appliance61Version) {
		if v, ok := d.GetOk("client_profile_settings"); ok {
			settings, err := readPolicyClientProfileSettingsFromConfig(v.([]interface{}))
			if err != nil {
				return err
			}
			args.SetClientProfileSettings(settings)
		}
		if v, ok := d.GetOk("custom_client_help_url"); ok {
			args.SetCustomClientHelpUrl(v.(string))
		}
	}
	if currentVersion.GreaterThanOrEqual(Appliance55Version) {
		if v, ok := d.GetOk("type"); ok {
			args.SetType(v.(string))
		}
		if args.GetType() == "Dns" {
			args.SetTamperProofing(false)
		}
		if v, ok := d.GetOk("override_site_claim"); ok {
			args.SetOverrideSiteClaim(v.(string))
		}
		if v, ok := d.GetOk("dns_settings"); ok {
			if args.GetType() != "Dns" {
				return fmt.Errorf("appgatesdp_policy.dns_settings is only allowed on policy Type 'Dns', got %q", args.GetType())
			}
			servers, err := readPolicyDnsSettingsFromConfig(v.(*schema.Set).List())
			if err != nil {
				return err
			}
			args.SetDnsSettings(servers)
		}
	}

	if c, ok := d.GetOk("entitlements"); ok {
		entitlements, err := readArrayOfStringsFromConfig(c.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetEntitlements(entitlements)
	}

	if c, ok := d.GetOk("entitlement_links"); ok {
		entitlementLinks, err := readArrayOfStringsFromConfig(c.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetEntitlementLinks(entitlementLinks)
	}

	if c, ok := d.GetOk("ringfence_rules"); ok {
		ringfenceRules, err := readArrayOfStringsFromConfig(c.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetRingfenceRules(ringfenceRules)
	}

	if c, ok := d.GetOk("ringfence_rule_links"); ok {
		ringfenceRuleLinks, err := readArrayOfStringsFromConfig(c.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetRingfenceRuleLinks(ringfenceRuleLinks)
	}
	if v, ok := d.GetOkExists("tamper_proofing"); ok {
		args.SetTamperProofing(v.(bool))
	}

	if c, ok := d.GetOk("override_site"); ok {
		args.SetOverrideSite(c.(string))
	}
	if v, ok := d.GetOk("proxy_auto_config"); ok {
		if currentVersion.LessThan(Appliance53Version) {
			return fmt.Errorf("proxy_auto_config not supported on %q client v%d", currentVersion, meta.(*Client).ClientVersion)
		}
		args.SetProxyAutoConfig(readProxyAutoConfigFromConfig(v.([]interface{})))
	}

	if v, ok := d.GetOk("trusted_network_check"); ok {
		args.SetTrustedNetworkCheck(readTrustedNetworkCheckFromConfig(v.([]interface{})))
	}

	if c, ok := d.GetOk("administrative_roles"); ok {
		administrativeRoles, err := readArrayOfStringsFromConfig(c.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetAdministrativeRoles(administrativeRoles)
	}

	request := api.PoliciesPost(ctx)
	request = request.Policy(args)
	policy, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create policy %w", prettyPrintAPIError(err))
	}

	d.SetId(policy.GetId())
	return resourceAppgatePolicyRead(d, meta)
}

func readTrustedNetworkCheckFromConfig(trustedNetworks []interface{}) openapi.PolicyAllOfTrustedNetworkCheck {
	result := openapi.PolicyAllOfTrustedNetworkCheck{}
	for _, r := range trustedNetworks {
		if r == nil {
			continue
		}
		raw := r.(map[string]interface{})
		if v, ok := raw["enabled"]; ok {
			result.SetEnabled(v.(bool))
		}
		if v, ok := raw["dns_suffix"]; ok {
			result.SetDnsSuffix(v.(string))
		}
	}
	return result
}

func readPolicyClientProfileSettingsFromConfig(settings []interface{}) (openapi.PolicyAllOfClientProfileSettings, error) {
	result := openapi.PolicyAllOfClientProfileSettings{}
	for _, r := range settings {
		if r == nil {
			continue
		}
		raw := r.(map[string]interface{})
		if v, ok := raw["enabled"]; ok {
			result.SetEnabled(v.(bool))
		}
		if v, ok := raw["profiles"]; ok {
			profiles, err := readArrayOfStringsFromConfig(v.(*schema.Set).List())
			if err != nil {
				return result, err
			}
			result.SetProfiles(profiles)
		}

	}
	return result, nil
}

func readPolicyClientSettingsFromConfig(settings []interface{}) (openapi.PolicyAllOfClientSettings, error) {
	result := openapi.PolicyAllOfClientSettings{}
	for _, r := range settings {
		if r == nil {
			continue
		}
		raw := r.(map[string]interface{})
		if v, ok := raw["enabled"]; ok {
			result.SetEnabled(v.(bool))
		}
		if v, ok := raw["entitlements_list"].(string); ok && len(v) > 0 {
			result.SetEntitlementsList(v)
		}
		if v, ok := raw["attention_level"].(string); ok && len(v) > 0 {
			result.SetAttentionLevel(v)
		}
		if v, ok := raw["auto_start"].(string); ok && len(v) > 0 {
			result.SetAutoStart(v)
		}
		if v, ok := raw["add_remove_profiles"].(string); ok && len(v) > 0 {
			result.SetAddRemoveProfiles(v)
		}
		if v, ok := raw["keep_me_signed_in"].(string); ok && len(v) > 0 {
			result.SetKeepMeSignedIn(v)
		}
		if v, ok := raw["saml_auto_sign_in"].(string); ok && len(v) > 0 {
			result.SetSamlAutoSignIn(v)
		}
		if v, ok := raw["quit"].(string); ok && len(v) > 0 {
			result.SetQuit(v)
		}
		if v, ok := raw["sign_out"].(string); ok && len(v) > 0 {
			result.SetSignOut(v)
		}
		if v, ok := raw["suspend"].(string); ok && len(v) > 0 {
			result.SetSuspend(v)
		}
	}
	return result, nil
}

func readPolicyDnsSettingsFromConfig(dnsSettings []interface{}) ([]openapi.PolicyAllOfDnsSettings, error) {
	list := make([]openapi.PolicyAllOfDnsSettings, 0, 0)
	for _, r := range dnsSettings {
		if r == nil {
			continue
		}
		result := openapi.PolicyAllOfDnsSettings{}
		raw := r.(map[string]interface{})
		if v, ok := raw["domain"].(string); ok && len(v) > 0 {
			result.SetDomain(v)
		}
		if v, ok := raw["servers"]; ok {
			servers, err := readArrayOfStringsFromConfig(v.(*schema.Set).List())
			if err != nil {
				return list, fmt.Errorf("Failed to resolve dns_settings.servers: %w", err)
			}
			if len(servers) > 0 {
				result.SetServers(servers)
			}
		}
		list = append(list, result)
	}
	log.Printf("[DEBUG] readPolicyDnsSettingsFromConfig Result %+v", list)
	return list, nil
}

func readProxyAutoConfigFromConfig(proxyAutoConfigs []interface{}) openapi.PolicyAllOfProxyAutoConfig {
	pac := openapi.PolicyAllOfProxyAutoConfig{}
	for _, r := range proxyAutoConfigs {
		if r == nil {
			continue
		}
		raw := r.(map[string]interface{})
		if v, ok := raw["enabled"]; ok {
			pac.SetEnabled(v.(bool))
		}
		if v, ok := raw["url"]; ok {
			pac.SetUrl(v.(string))
		}
		if v, ok := raw["persist"]; ok {
			pac.SetPersist(v.(bool))
		}
	}
	return pac
}

func resourceAppgatePolicyRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Policy with name: %s", d.Get("name").(string))
	ctx := context.Background()
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.PoliciesApi
	currentVersion := meta.(*Client).ApplianceVersion
	request := api.PoliciesIdGet(ctx, d.Id())
	policy, response, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if response != nil && response.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Failed to read policy, %w", err)
	}
	d.Set("policy_id", policy.GetId())
	d.Set("name", policy.GetName())
	d.Set("notes", policy.GetNotes())
	d.Set("disabled", policy.GetDisabled())
	d.Set("expression", policy.GetExpression())
	d.Set("entitlements", policy.GetEntitlements())
	d.Set("entitlement_links", policy.GetEntitlementLinks())
	d.Set("ringfence_rule_links", policy.GetRingfenceRuleLinks())
	d.Set("ringfence_rules", policy.GetRingfenceRules())
	d.Set("tags", policy.GetTags())
	d.Set("tamper_proofing", policy.GetTamperProofing())
	d.Set("administrative_roles", policy.GetAdministrativeRoles())

	if v, o := policy.GetProxyAutoConfigOk(); o != false {
		pac, err := flattenProxyAutoConfig(*v)
		if err != nil {
			return err
		}
		if currentVersion.GreaterThanOrEqual(Appliance53Version) {
			d.Set("proxy_auto_config", pac)
		}
	}
	if v, o := policy.GetTrustedNetworkCheckOk(); o != false {
		t, err := flattenTrustedNetworkCheck(*v)
		if err != nil {
			return err
		}
		d.Set("trusted_network_check", t)
	}
	if currentVersion.GreaterThanOrEqual(Appliance54Version) {
		clientSettings, err := flattenPolicyClientSettings(policy.GetClientSettings())
		if err != nil {
			return err
		}
		d.Set("client_settings", clientSettings)
	}
	if currentVersion.GreaterThanOrEqual(Appliance55Version) {
		d.Set("type", policy.GetType())
		d.Set("override_site_claim", policy.GetOverrideSiteClaim())
		dnsSettings, err := flattenPolicyDnsSettings(policy.GetDnsSettings())
		if err != nil {
			return err
		}
		d.Set("dns_settings", dnsSettings)
	}
	if currentVersion.GreaterThanOrEqual(Appliance61Version) {
		clientProfileSettings, err := flattenPolicyClientProfileSettings(policy.GetClientProfileSettings())
		if err != nil {
			return err
		}
		d.Set("client_profile_settings", clientProfileSettings)
		d.Set("custom_client_help_url", policy.GetCustomClientHelpUrl())
	}
	return nil
}

func flattenProxyAutoConfig(in openapi.PolicyAllOfProxyAutoConfig) ([]interface{}, error) {
	m := make(map[string]interface{})
	if v, o := in.GetEnabledOk(); o != false {
		m["enabled"] = *v
	}
	if v, o := in.GetUrlOk(); o != false {
		m["url"] = *v
	}
	if v, o := in.GetPersistOk(); o != false {
		m["persist"] = *v
	}

	return []interface{}{m}, nil
}

func flattenTrustedNetworkCheck(in openapi.PolicyAllOfTrustedNetworkCheck) ([]interface{}, error) {
	m := make(map[string]interface{})
	if v, o := in.GetEnabledOk(); o != false {
		m["enabled"] = v
	}
	if v, o := in.GetDnsSuffixOk(); o != false {
		m["dns_suffix"] = v
	}
	return []interface{}{m}, nil
}

func flattenPolicyClientProfileSettings(clientSettings openapi.PolicyAllOfClientProfileSettings) ([]interface{}, error) {
	m := make(map[string]interface{})
	if v, ok := clientSettings.GetEnabledOk(); ok {
		m["enabled"] = *v
	}
	if v, ok := clientSettings.GetProfilesOk(); ok {
		m["profiles"] = v
	}
	return []interface{}{m}, nil
}

func flattenPolicyClientSettings(clientSettings openapi.PolicyAllOfClientSettings) ([]interface{}, error) {
	m := make(map[string]interface{})
	if v, ok := clientSettings.GetEnabledOk(); ok {
		m["enabled"] = *v
	}
	if v, ok := clientSettings.GetEntitlementsListOk(); ok {
		m["entitlements_list"] = *v
	}
	if v, ok := clientSettings.GetAttentionLevelOk(); ok {
		m["attention_level"] = *v
	}
	if v, ok := clientSettings.GetAutoStartOk(); ok {
		m["auto_start"] = *v
	}
	if v, ok := clientSettings.GetAddRemoveProfilesOk(); ok {
		m["add_remove_profiles"] = *v
	}
	if v, ok := clientSettings.GetKeepMeSignedInOk(); ok {
		m["keep_me_signed_in"] = *v
	}
	if v, ok := clientSettings.GetSamlAutoSignInOk(); ok {
		m["saml_auto_sign_in"] = *v
	}
	if v, ok := clientSettings.GetQuitOk(); ok {
		m["quit"] = *v
	}
	if v, ok := clientSettings.GetSignOutOk(); ok {
		m["sign_out"] = *v
	}
	if v, ok := clientSettings.GetSuspendOk(); ok {
		m["suspend"] = *v
	}
	return []interface{}{m}, nil
}

func flattenPolicyDnsSettings(dnsSettings []openapi.PolicyAllOfDnsSettings) (*schema.Set, error) {
	out := make([]interface{}, 0)
	for _, dnsSetting := range dnsSettings {
		m := make(map[string]interface{})
		if v, ok := dnsSetting.GetDomainOk(); ok {
			m["domain"] = *v
		}
		if v, ok := dnsSetting.GetServersOk(); ok {
			m["servers"] = schema.NewSet(schema.HashString, convertStringArrToInterface(v))
		}
		out = append(out, m)
	}
	return schema.NewSet(resourcePolicyDnsSettingsHash, out), nil
}

func resourceAppgatePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating policy: %s", d.Get("name").(string))
	ctx := context.Background()
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.PoliciesApi
	currentVersion := meta.(*Client).ApplianceVersion
	request := api.PoliciesIdGet(ctx, d.Id())
	orginalPolicy, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read policy, %w", err)
	}

	if d.HasChange("name") {
		orginalPolicy.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		orginalPolicy.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		orginalPolicy.SetTags(schemaExtractTags(d))
	}

	if d.HasChange("disabled") {
		orginalPolicy.SetDisabled(d.Get("disabled").(bool))
	}

	if d.HasChange("expression") {
		orginalPolicy.SetExpression(d.Get("expression").(string))
	}

	if d.HasChange("entitlements") {
		_, n := d.GetChange("entitlements")
		entitlements, err := readArrayOfStringsFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		orginalPolicy.SetEntitlements(entitlements)
	}

	if d.HasChange("entitlement_links") {
		_, n := d.GetChange("entitlement_links")
		entitlements, err := readArrayOfStringsFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		orginalPolicy.SetEntitlementLinks(entitlements)
	}

	if d.HasChange("ringfence_rules") {
		_, n := d.GetChange("ringfence_rules")
		entitlements, err := readArrayOfStringsFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		orginalPolicy.SetRingfenceRules(entitlements)
	}

	if d.HasChange("ringfence_rule_links") {
		_, n := d.GetChange("ringfence_rule_links")
		entitlements, err := readArrayOfStringsFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		orginalPolicy.SetRingfenceRuleLinks(entitlements)
	}

	if d.HasChange("tamper_proofing") {
		orginalPolicy.SetTamperProofing(d.Get("tamper_proofing").(bool))
	}

	if d.HasChange("override_site") {
		orginalPolicy.SetOverrideSite(d.Get("override_site").(string))
	}

	if d.HasChange("proxy_auto_config") {
		_, v := d.GetChange("proxy_auto_config")
		orginalPolicy.SetProxyAutoConfig(readProxyAutoConfigFromConfig(v.([]interface{})))
	}

	if d.HasChange("trusted_network_check") {
		_, v := d.GetChange("trusted_network_check")
		orginalPolicy.SetTrustedNetworkCheck(readTrustedNetworkCheckFromConfig(v.([]interface{})))
	}

	if d.HasChange("administrative_roles") {
		_, n := d.GetChange("administrative_roles")
		entitlements, err := readArrayOfStringsFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		orginalPolicy.SetAdministrativeRoles(entitlements)
	}
	if currentVersion.GreaterThanOrEqual(Appliance54Version) {
		if d.HasChange("client_settings") {
			_, v := d.GetChange("client_settings")
			clientSettings, err := readPolicyClientSettingsFromConfig(v.([]interface{}))
			if err != nil {
				return err
			}
			orginalPolicy.SetClientSettings(clientSettings)
		}
	}
	if currentVersion.GreaterThanOrEqual(Appliance55Version) {
		if d.HasChange("type") {
			orginalPolicy.SetType(d.Get("type").(string))
		}
		if d.HasChange("override_site_claim") {
			orginalPolicy.SetOverrideSiteClaim(d.Get("override_site_claim").(string))
		}
		if d.HasChange("dns_settings") {
			if orginalPolicy.GetType() != "Dns" {
				return fmt.Errorf("appgatesdp_policy.dns_settings is only allowed on policy Type 'Dns', got %q", orginalPolicy.GetType())
			}
			_, v := d.GetChange("dns_settings")
			dnsSettings, err := readPolicyDnsSettingsFromConfig(v.(*schema.Set).List())
			if err != nil {
				return err
			}
			orginalPolicy.SetDnsSettings(dnsSettings)
		}
	}
	if currentVersion.GreaterThanOrEqual(Appliance61Version) {
		if d.HasChange("client_profile_settings") {
			_, v := d.GetChange("client_profile_settings")
			clientProfileSettings, err := readPolicyClientProfileSettingsFromConfig(v.([]interface{}))
			if err != nil {
				return err
			}
			orginalPolicy.SetClientProfileSettings(clientProfileSettings)
		}
		if d.HasChange("custom_client_help_url") {
			orginalPolicy.SetCustomClientHelpUrl(d.Get("custom_client_help_url").(string))
		}
	}
	req := api.PoliciesIdPut(ctx, d.Id())

	_, _, err = req.Policy(*orginalPolicy).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update policy %w", prettyPrintAPIError(err))
	}

	return resourceAppgatePolicyRead(d, meta)
}

func resourceAppgatePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete Policy with name: %s", d.Get("name").(string))
	ctx := context.Background()
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.PoliciesApi

	// Get policy
	request := api.PoliciesIdGet(ctx, d.Id())
	policy, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to delete policy while GET, %w", err)
	}

	deleteRequest := api.PoliciesIdDelete(ctx, policy.GetId())
	_, err = deleteRequest.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to delete policy, %w", err)
	}
	d.SetId("")
	return nil
}
