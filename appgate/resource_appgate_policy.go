package appgate

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/appgate/terraform-provider-appgate/client/v12/openapi"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/schema"
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

			"tags": {
				Type:        schema.TypeSet,
				Description: "Array of tags.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

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
	token := meta.(*Client).Token
	api := meta.(*Client).API.PoliciesApi

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

func resourceAppgatePolicyRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Policy with name: %s", d.Get("name").(string))
	ctx := context.Background()
	token := meta.(*Client).Token
	api := meta.(*Client).API.PoliciesApi

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

	return nil
}

func resourceAppgatePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating policy: %s", d.Get("name").(string))
	ctx := context.Background()
	token := meta.(*Client).Token
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
	token := meta.(*Client).Token
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
