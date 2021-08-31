package appgate

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v15/openapi"

	"github.com/google/uuid"
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"policy_id": {
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

			"entitlements": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"entitlement_links": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"ringfence_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"ringfence_rule_links": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"tamper_proofing": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"override_site": {
				Type:     schema.TypeString,
				Optional: true,
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

			"administrative_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
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
	args := openapi.NewPolicyWithDefaults()
	args.Id = uuid.New().String()

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

	if c, ok := d.GetOk("tamper_proofing"); ok {
		args.SetTamperProofing(c.(bool))
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
	request = request.Policy(*args)
	policy, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create policy %+v", prettyPrintAPIError(err))
	}

	d.SetId(policy.Id)
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
	policy, _, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Failed to read policy, %+v", err)
	}
	d.Set("policy_id", policy.Id)
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
	log.Printf("[DEBUG] flattenTrustedNetworkCheck: %+v", m)
	return []interface{}{m}, nil
}

func resourceAppgatePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating policy: %s", d.Get("name").(string))
	ctx := context.Background()
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.PoliciesApi
	request := api.PoliciesIdGet(ctx, d.Id())
	orginalPolicy, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read policy, %+v", err)
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

	req := api.PoliciesIdPut(ctx, d.Id())

	_, _, err = req.Policy(orginalPolicy).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update policy %+v", prettyPrintAPIError(err))
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
		return fmt.Errorf("Failed to delete policy while GET, %+v", err)
	}

	deleteRequest := api.PoliciesIdDelete(ctx, policy.GetId())
	_, err = deleteRequest.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to delete policy, %+v", err)
	}
	d.SetId("")
	return nil
}
