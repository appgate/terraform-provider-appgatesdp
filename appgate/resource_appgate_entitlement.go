package appgate

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v16/openapi"
	"github.com/appgate/terraform-provider-appgatesdp/appgate/hashcode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateEntitlement() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppgateEntitlementRuleCreate,
		ReadContext:   resourceAppgateEntitlementRuleRead,
		UpdateContext: resourceAppgateEntitlementRuleUpdate,
		DeleteContext: resourceAppgateEntitlementRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"entitlement_id": resourceUUID(),

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

			"disabled": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},

			"site": {
				Type:     schema.TypeString,
				Required: true,
			},

			"condition_logic": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
					s := v.(string)
					list := []string{"and", "or"}
					for _, x := range list {
						if s == x {
							return
						}
					}
					errs = append(errs, fmt.Errorf("condition_logic must be on of %v, got %s", list, s))
					return
				},
			},

			"conditions": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "List of Condition IDs applies to this Entitlement.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"actions": {
				Type:             schema.TypeSet,
				Required:         true,
				Set:              resourceAppgateEntitlementActionHash,
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"subtype": {
							Type:     schema.TypeString,
							Required: true,
						},

						"action": {
							Type:     schema.TypeString,
							Required: true,
						},

						"hosts": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      schema.HashString,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"ports": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      schema.HashString,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"types": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"monitor": {
							Type:             schema.TypeList,
							MaxItems:         1,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:             schema.TypeBool,
										Optional:         true,
										Computed:         true,
										DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
									},
									"timeout": {
										Type:             schema.TypeInt,
										Optional:         true,
										Computed:         true,
										DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
									},
								},
							},
						},
					},
				},
			},

			"app_shortcuts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"color_code": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},

			"app_shortcut_scripts": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAppgateEntitlementActionHash(v interface{}) int {
	raw := v.(map[string]interface{})
	// modifying raw actually modifies the values passed to the provider.
	// Use a copy to avoid that.
	copy := make((map[string]interface{}))
	for key, value := range raw {
		copy[key] = value
	}
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", copy["subtype"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", copy["action"].(string)))
	if v, ok := copy["hosts"]; ok {
		buf.WriteString(fmt.Sprintf("%v-", v.(*schema.Set).List()))
	}
	if v, ok := copy["ports"]; ok {
		buf.WriteString(fmt.Sprintf("%v-", v.(*schema.Set).List()))
	}
	if v, ok := copy["types"]; ok {
		vs := v.([]interface{})
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		sort.Strings(s)

		for _, v := range s {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}

	// monitor is only valid if subtype is tcp_up
	if copy["subtype"].(string) == "tcp_up" {
		if v, ok := copy["monitor"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			mHash := resourceAppgateEntitlementActionMonitorHash(v[0])
			buf.WriteString(fmt.Sprintf("%d-", mHash))
		} else {
			// hack, write default values here to correct the hash value.
			v := map[string]interface{}{
				"enabled": false,
				"timeout": 30,
			}
			mHash := resourceAppgateEntitlementActionMonitorHash(v)
			buf.WriteString(fmt.Sprintf("%d-", mHash))
		}
	}
	return hashcode.String(buf.String())
}

func resourceAppgateEntitlementActionMonitorHash(v interface{}) int {
	var buf bytes.Buffer
	m, ok := v.(map[string]interface{})

	if !ok {
		return 0
	}

	if v, ok := m["enabled"]; ok {
		buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
	}
	if v, ok := m["timeout"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}

	return hashcode.String(buf.String())
}

func resourceAppgateEntitlementRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	log.Printf("[DEBUG] Creating Entitlement: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.EntitlementsApi

	args := openapi.NewEntitlementWithDefaults()
	if v, ok := d.GetOk("entitlement_id"); ok {
		args.SetId(v.(string))
	}
	args.SetName(d.Get("name").(string))
	args.SetSite(d.Get("site").(string))
	args.SetNotes(d.Get("notes").(string))
	args.SetTags(schemaExtractTags(d))
	args.SetDisabled(d.Get("disabled").(bool))

	if v, ok := d.GetOk("condition_logic"); ok {
		args.SetConditionLogic(v.(string))
	}

	if v, ok := d.GetOk("conditions"); ok {
		conditions, err := readArrayOfStringsFromConfig(v.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetConditions(conditions)
	}

	if v, ok := d.GetOk("actions"); ok {
		actions, _, err := readEntitlmentActionsFromConfig(v.(*schema.Set).List(), diags)
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetActions(actions)
	}

	if v, ok := d.GetOk("app_shortcuts"); ok {
		appShortcuts, err := readAppShortcutFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetAppShortcuts(appShortcuts)
	}

	if v, ok := d.GetOk("app_shortcut_scripts"); ok {
		scripts, err := readArrayOfStringsFromConfig(v.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetAppShortcutScripts(scripts)
	}

	request := api.EntitlementsPost(context.Background())
	request = request.Entitlement(*args)
	ent, _, err := request.Authorization(token).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Could not create entitlement %w", prettyPrintAPIError(err)))
	}

	d.SetId(ent.Id)
	d.Set("entitlement_id", ent.Id)
	resourceAppgateEntitlementRuleRead(ctx, d, meta)

	return diags
}

func resourceAppgateEntitlementRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	log.Printf("[DEBUG] Reading Entitlement Name: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Reading Entitlement id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.EntitlementsApi

	request := api.EntitlementsIdGet(ctx, d.Id())
	entitlement, res, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.FromErr(fmt.Errorf("Failed to read Entitlement, %w", err))
	}
	d.SetId(entitlement.Id)
	d.Set("entitlement_id", entitlement.Id)
	d.Set("name", entitlement.Name)
	d.Set("disabled", entitlement.Disabled)
	d.Set("notes", entitlement.Notes)
	d.Set("conditions", entitlement.Conditions)
	d.Set("condition_logic", entitlement.ConditionLogic)
	d.Set("tags", entitlement.Tags)
	d.Set("site", entitlement.Site)
	if entitlement.AppShortcuts != nil {
		if err = d.Set("app_shortcuts", flattenEntitlementAppShortcut(*entitlement.AppShortcuts)); err != nil {
			return diag.FromErr(err)
		}
	}

	actions := flattenEntitlementActions(entitlement.Actions, d)
	if err = d.Set("actions", actions); err != nil {
		return diag.FromErr(err)
	}

	if v, ok := entitlement.GetAppShortcutScriptsOk(); ok {
		d.Set("app_shortcut_scripts", *v)
	}

	return diags
}

func flattenEntitlementAppShortcut(in []openapi.AppShortcut) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["name"] = v.GetName()
		m["description"] = v.GetDescription()
		m["url"] = v.GetUrl()
		m["color_code"] = v.GetColorCode()

		out[i] = m
	}

	return out
}

func icmpTypes() []string {
	return []string{"icmp_up", "icmp_down", "icmpv6_up", "icmpv6_down"}
}

func flattenEntitlementActions(actions []openapi.EntitlementAllOfActions, d *schema.ResourceData) *schema.Set {
	out := []interface{}{}
	for _, act := range actions {
		action := make(map[string]interface{})
		action["subtype"] = act.Subtype
		action["action"] = act.Action
		action["hosts"] = schema.NewSet(schema.HashString, convertStringArrToInterface(act.GetHosts()))
		action["ports"] = schema.NewSet(schema.HashString, convertStringArrToInterface(act.GetPorts()))
		types := act.GetTypes()
		if types != nil && inArray(act.Subtype, icmpTypes()) {
			action["types"] = convertStringArrToInterface(act.GetTypes())
		}
		if act.Monitor != nil && act.Subtype == "tcp_up" {
			action["monitor"] = flattenEntitlementActionMonitor(act.GetMonitor())
			dataActions := d.Get("actions")
			hash := resourceAppgateEntitlementActionHash(action)

			for _, k := range dataActions.(*schema.Set).List() {
				log.Printf("[DEBUG] flattenEntitlementActions action list OLD %+v", k)
				oldHash := resourceAppgateEntitlementActionHash(k)
				if oldHash == hash {
					oldV := k.(map[string]interface{})
					if v, ok := oldV["monitor"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
						log.Printf("[DEBUG] flattenEntitlementActions OLD MONITOR %+v", v)
						action["monitor"] = v
					}
				}
			}
		}
		out = append(out, action)
	}
	return schema.NewSet(resourceAppgateEntitlementActionHash, out)
}

func flattenEntitlementActionMonitor(monitor openapi.EntitlementAllOfMonitor) []interface{} {
	log.Printf("[DEBUG] flattenEntitlementActionMonitor %+v", monitor)
	m := make(map[string]interface{})
	m["enabled"] = *monitor.Enabled
	m["timeout"] = int(*monitor.Timeout)

	return []interface{}{m}
}

func resourceAppgateEntitlementRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating Entitlement: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Updating Entitlement id: %+v", d.Id())
	var diags diag.Diagnostics
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.EntitlementsApi

	request := api.EntitlementsIdGet(ctx, d.Id())
	orginalEntitlment, _, err := request.Authorization(token).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to read Entitlement while updating, %w", err))
	}

	if d.HasChange("name") {
		orginalEntitlment.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		orginalEntitlment.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		orginalEntitlment.SetTags(schemaExtractTags(d))
	}

	if d.HasChange("disabled") {
		orginalEntitlment.SetDisabled(d.Get("disabled").(bool))
	}

	if d.HasChange("site") {
		orginalEntitlment.SetSite(d.Get("site").(string))
	}

	if d.HasChange("condition_logic") {
		orginalEntitlment.SetConditionLogic(d.Get("condition_logic").(string))
	}

	if d.HasChange("conditions") {
		_, n := d.GetChange("conditions")
		conditions, err := readArrayOfStringsFromConfig(n.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		orginalEntitlment.SetConditions(conditions)
	}

	if d.HasChange("actions") {
		_, v := d.GetChange("actions")
		actions, _, err := readEntitlmentActionsFromConfig(v.(*schema.Set).List(), diags)
		if err != nil {
			return diag.FromErr(err)
		}
		orginalEntitlment.SetActions(actions)
	}

	if d.HasChange("app_shortcut") {
		_, n := d.GetChange("app_shortcut")
		appShortcut, err := readAppShortcutFromConfig(n.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		orginalEntitlment.SetAppShortcuts(appShortcut)
	}

	if d.HasChange("app_shortcut_scripts") {
		_, v := d.GetChange("app_shortcut_scripts")
		scripts, err := readArrayOfStringsFromConfig(v.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		orginalEntitlment.SetAppShortcutScripts(scripts)
	}

	req := api.EntitlementsIdPut(ctx, d.Id())
	req = req.Entitlement(orginalEntitlment)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Could not update Entitlement %w", prettyPrintAPIError(err)))
	}

	return resourceAppgateEntitlementRuleRead(ctx, d, meta)
}

func resourceAppgateEntitlementRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Delete Entitlement: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Reading Entitlement id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.EntitlementsApi

	if _, err := api.EntitlementsIdDelete(ctx, d.Id()).Authorization(token).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("Could not delete Entitlement %w", prettyPrintAPIError(err)))
	}
	d.SetId("")
	return diags
}

func readEntitlmentActionsFromConfig(actions []interface{}, diags diag.Diagnostics) ([]openapi.EntitlementAllOfActions, diag.Diagnostics, error) {
	result := make([]openapi.EntitlementAllOfActions, 0)
	for _, action := range actions {
		if action == nil {
			continue
		}
		a := openapi.EntitlementAllOfActions{}
		raw := action.(map[string]interface{})
		if v, ok := raw["subtype"]; ok {
			a.SetSubtype(v.(string))
		}
		if v, ok := raw["action"]; ok {
			a.SetAction(v.(string))
		}
		if v, ok := raw["hosts"]; ok {
			hosts, err := readArrayOfStringsFromConfig(v.(*schema.Set).List())
			if err != nil {
				return result, diags, fmt.Errorf("Failed to resolve entitlement action hosts: %w", err)
			}
			a.SetHosts(hosts)
		}
		if v, ok := raw["ports"]; ok {
			ports, err := readArrayOfStringsFromConfig(v.(*schema.Set).List())
			if err != nil {
				return result, diags, fmt.Errorf("Failed to resolve entitlement action ports: %w", err)
			}
			a.SetPorts(ports)
		}
		if v := raw["types"]; len(v.([]interface{})) > 0 {
			if !inArray(a.GetSubtype(), icmpTypes()) {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Invalid usage of types.",
					Detail:   "ICMP type. Only valid for icmp subtypes.",
				})
			}
			types, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, diags, fmt.Errorf("Failed to resolve entitlment action types: %w", err)
			}
			a.SetTypes(types)
		}

		if v, ok := raw["monitor"].([]interface{}); ok && len(v) > 0 {
			monitor := openapi.NewEntitlementAllOfMonitorWithDefaults()
			rawMonitors := v
			for _, v := range rawMonitors {
				rawMonitor := v.(map[string]interface{})
				if v, ok := rawMonitor["enabled"]; ok {
					monitor.SetEnabled(v.(bool))
				}
				if v, ok := rawMonitor["timeout"]; ok {
					monitor.SetTimeout(int32(v.(int)))
				}
			}
			a.SetMonitor(*monitor)
		}
		result = append(result, a)
	}
	return result, diags, nil
}

func readAppShortcutFromConfig(shortcuts []interface{}) ([]openapi.AppShortcut, error) {
	result := make([]openapi.AppShortcut, 0)
	for _, shortcut := range shortcuts {
		if shortcut == nil {
			continue
		}
		row := openapi.AppShortcut{}
		raw := shortcut.(map[string]interface{})
		if v, ok := raw["name"]; ok {
			row.SetName(v.(string))
		}
		if v, ok := raw["url"]; ok {
			row.SetUrl(v.(string))
		}
		if v, ok := raw["description"]; ok {
			row.SetDescription(v.(string))
		}
		if v, ok := raw["color_code"]; ok {
			row.SetColorCode(int32(v.(int)))
		}
		result = append(result, row)
	}
	log.Printf("[DEBUG] readAppShortcutFromConfig: %+v", result)
	return result, nil
}
