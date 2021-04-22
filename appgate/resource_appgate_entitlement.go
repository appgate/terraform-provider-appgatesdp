package appgate

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v14/openapi"
	"github.com/appgate/terraform-provider-appgatesdp/appgate/hashcode"
	"github.com/google/uuid"

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

			"entitlement_id": {
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
				Type:     schema.TypeSet,
				Required: true,
				Set:      resourceAppgateEntitlementActionHash,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					log.Printf("[DEBUG] DiffSuppressFunc ACTION: k %s old %s new %s", k, old, new)
					return old == "1" && new == "0"
				},
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
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"types": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"monitor_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if old == "" && new == "" {
									return true
								}
								if old == "" && new == "true" {
									return true
								}
								return old == "" && new == "false"
							},
						},

						"monitor_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if old == "" && new == "" {
									return true
								}
								if old == "" && len(new) > 0 {
									return true
								}
								return false
							},
						},

						"monitor": {
							Type:             schema.TypeList,
							MaxItems:         1,
							Optional:         true,
							Deprecated:       "monitor {} has been replaced by actions.monitor_enabled and actions.monitor_timeout",
							DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  30,
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
	m := v.(map[string]interface{})
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", m["subtype"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["action"].(string)))
	if v, ok := m["hosts"]; ok {
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
	if v, ok := m["ports"]; ok {
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
	if v, ok := m["types"]; ok {
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
	if _, ok := m["monitor_enabled"]; ok {
		buf.WriteString(fmt.Sprintf("%t-", m["monitor_enabled"].(bool)))
	}
	if _, ok := m["monitor_timeout"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", m["monitor_timeout"].(int)))
	}

	r := hashcode.String(buf.String())

	log.Printf("[DEBUG] action SET HASH: %d", r)
	return r
}

func resourceAppgateEntitlementRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	log.Printf("[DEBUG] Creating Entitlement: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.EntitlementsApi

	args := openapi.NewEntitlementWithDefaults()
	args.Id = uuid.New().String()
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
		actions, err := readConditionActionsFromConfig(v.(*schema.Set).List())
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
		return diag.FromErr(fmt.Errorf("Could not create entitlement %+v", prettyPrintAPIError(err)))
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
	token := meta.(*Client).Token
	api := meta.(*Client).API.EntitlementsApi

	request := api.EntitlementsIdGet(ctx, d.Id())
	entitlement, _, err := request.Authorization(token).Execute()
	if err != nil {
		// TODO check if 404
		d.SetId("")
		return diag.FromErr(fmt.Errorf("Failed to read Entitlement, %+v", err))
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
	if entitlement.Actions != nil {
		actions := flattenEntitlementActions(entitlement.Actions, d)
		if err = d.Set("actions", actions); err != nil {
			return diag.FromErr(err)
		}
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

func flattenEntitlementActions(in []openapi.EntitlementAllOfActions, d *schema.ResourceData) []interface{} {
	// var out = make([]map[string]interface{}, len(in), len(in))
	var out = make([]interface{}, 0)
	for _, v := range in {
		action := make(map[string]interface{})
		action["subtype"] = v.Subtype
		action["action"] = v.Action
		action["hosts"] = convertStringArrToInterface(v.GetHosts())
		action["ports"] = convertStringArrToInterface(v.GetPorts())
		action["types"] = convertStringArrToInterface(v.GetTypes())
		if v.Monitor != nil && action["sybtype"] == "tcp_up" {
			action["monitor_enabled"] = v.Monitor.GetEnabled()
			action["monitor_timeout"] = int(v.Monitor.GetTimeout())
			// Deprecated
			action["monitor"] = flattenEntitlementActionMonitor(v.Monitor)
		}

		out = append(out, action)
	}
	return out
}

func flattenEntitlementActionMonitor(in *openapi.EntitlementAllOfMonitor) []interface{} {
	log.Printf("[DEBUG] flattenEntitlementActionMonitor %+v", in)
	m := make(map[string]interface{})
	m["enabled"] = in.GetEnabled()
	m["timeout"] = int(in.GetTimeout())

	return []interface{}{m}
}

func resourceAppgateEntitlementRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating Entitlement: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Updating Entitlement id: %+v", d.Id())
	token := meta.(*Client).Token
	api := meta.(*Client).API.EntitlementsApi

	request := api.EntitlementsIdGet(ctx, d.Id())
	orginalEntitlment, _, err := request.Authorization(token).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to read Entitlement while updating, %+v", err))
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
		actions, err := readConditionActionsFromConfig(v.(*schema.Set).List())
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
		return diag.FromErr(fmt.Errorf("Could not update Entitlement %+v", prettyPrintAPIError(err)))
	}

	return resourceAppgateEntitlementRuleRead(ctx, d, meta)
}

func resourceAppgateEntitlementRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Delete Entitlement: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Reading Entitlement id: %+v", d.Id())
	token := meta.(*Client).Token
	api := meta.(*Client).API.EntitlementsApi

	request := api.EntitlementsIdDelete(context.Background(), d.Id())

	_, err := request.Authorization(token).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Could not delete Entitlement %+v", prettyPrintAPIError(err)))
	}
	d.SetId("")
	return diags
}

func readConditionActionsFromConfig(actions []interface{}) ([]openapi.EntitlementAllOfActions, error) {
	result := make([]openapi.EntitlementAllOfActions, 0)
	for _, action := range actions {
		if action == nil {
			continue
		}
		a := openapi.NewEntitlementAllOfActionsWithDefaults()
		raw := action.(map[string]interface{})
		log.Printf("[DEBUG] readConditionActionsFromConfig RAWT: %+v", raw)
		if v, ok := raw["subtype"]; ok {
			a.SetSubtype(v.(string))
		}
		if v, ok := raw["action"]; ok {
			a.SetAction(v.(string))
		}
		if v := raw["hosts"]; len(v.([]interface{})) > 0 {
			hosts, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, fmt.Errorf("Failed to resolve condition action hosts: %+v", err)
			}
			a.SetHosts(hosts)
		}
		if v := raw["ports"]; len(v.([]interface{})) > 0 {
			ports, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, fmt.Errorf("Failed to resolve condition action ports: %+v", err)
			}
			a.SetPorts(ports)
		}
		if v := raw["types"]; len(v.([]interface{})) > 0 {
			types, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, fmt.Errorf("Failed to resolve condition action types: %+v", err)
			}
			a.SetTypes(types)
		}
		monitor := openapi.NewEntitlementAllOfMonitorWithDefaults()

		// Check monitor for backwards compatibility
		if v, ok := raw["monitor"]; ok {
			rawMonitors := v.([]interface{})
			for _, v := range rawMonitors {
				rawMonitor := v.(map[string]interface{})
				if v, ok := rawMonitor["enabled"]; ok {
					raw["monitor_enabled"] = v.(bool)
				}
				if v, ok := rawMonitor["timeout"].(int); ok {
					raw["monitor_timeout"] = v
				}
			}
		}
		if v, ok := raw["monitor_enabled"].(bool); ok {
			monitor.SetEnabled(v)
		}
		if v, ok := raw["monitor_timeout"].(int); ok && v > 0 {
			monitor.SetTimeout(int32(v))
		}
		a.SetMonitor(*monitor)
		result = append(result, *a)
	}
	return result, nil
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
