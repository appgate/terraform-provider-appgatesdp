package appgate

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v13/openapi"

	"github.com/google/uuid"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAppgateEntitlement() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateEntitlementRuleCreate,
		Read:   resourceAppgateEntitlementRuleRead,
		Update: resourceAppgateEntitlementRuleUpdate,
		Delete: resourceAppgateEntitlementRuleDelete,
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
				Type:       schema.TypeSet,
				Required:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
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

						"monitor": {
							Type:       schema.TypeSet,
							Optional:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
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
		},
	}
}

func resourceAppgateEntitlementRuleCreate(d *schema.ResourceData, meta interface{}) error {
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

	if c, ok := d.GetOk("conditions"); ok {
		conditions, err := readArrayOfStringsFromConfig(c.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetConditions(conditions)
	}

	if c, ok := d.GetOk("actions"); ok {
		actions, err := readConditionActionsFromConfig(c.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetActions(actions)
	}

	if v, ok := d.GetOk("app_shortcuts"); ok {
		appShortcuts, err := readAppShortcutFromConfig(v.([]interface{}))
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Creating Entitlement SHORTCUT: %+v", appShortcuts)
		args.SetAppShortcuts(appShortcuts)
	}

	request := api.EntitlementsPost(context.Background())
	request = request.Entitlement(*args)
	ent, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create entitlement %+v", prettyPrintAPIError(err))
	}

	d.SetId(ent.Id)
	d.Set("entitlement_id", ent.Id)
	return resourceAppgateEntitlementRuleRead(d, meta)
}

func resourceAppgateEntitlementRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Entitlement Name: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Reading Entitlement id: %+v", d.Id())
	token := meta.(*Client).Token
	api := meta.(*Client).API.EntitlementsApi
	ctx := context.TODO()
	request := api.EntitlementsIdGet(ctx, d.Id())
	entitlement, _, err := request.Authorization(token).Execute()
	if err != nil {
		// TODO check if 404
		d.SetId("")
		return fmt.Errorf("Failed to read Entitlement, %+v", err)
	}
	d.SetId(entitlement.Id)
	d.Set("entitlement_id", entitlement.Id)
	d.Set("name", entitlement.Name)
	d.Set("disabled", entitlement.Disabled)
	d.Set("notes", entitlement.Notes)
	d.Set("actions", entitlement.Actions)
	d.Set("conditions", entitlement.Conditions)
	d.Set("condition_logic", entitlement.ConditionLogic)
	d.Set("tags", entitlement.Tags)
	d.Set("site", entitlement.Site)
	if entitlement.AppShortcuts != nil {
		if err = d.Set("app_shortcuts", flattenEntitlementAppShortcut(*entitlement.AppShortcuts)); err != nil {
			return err
		}
	}
	if entitlement.Actions != nil {
		if err = d.Set("actions", flattenEntitlementActions(entitlement.Actions)); err != nil {
			return err
		}
	}

	return nil
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

func flattenEntitlementActions(in []openapi.EntitlementAllOfActions) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["subtype"] = v.Subtype
		m["action"] = v.Action
		m["hosts"] = v.Hosts
		m["ports"] = v.Ports
		m["types"] = v.Types
		if v.Monitor != nil {
			m["monitor"] = flattenEntitlementActionMonitor(v.Monitor)
		}
		out[i] = m
	}
	return out
}

func flattenEntitlementActionMonitor(in *openapi.EntitlementAllOfMonitor) []interface{} {
	log.Printf("[DEBUG] flattenEntitlementActionMonitor %+v", in)
	m := make(map[string]interface{})
	m["enabled"] = in.Enabled
	m["timeout"] = in.Timeout

	return []interface{}{m}
}

func resourceAppgateEntitlementRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Entitlement: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Updating Entitlement id: %+v", d.Id())
	token := meta.(*Client).Token
	api := meta.(*Client).API.EntitlementsApi
	ctx := context.TODO()
	request := api.EntitlementsIdGet(ctx, d.Id())
	orginalEntitlment, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read Entitlement while updating, %+v", err)
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
			return err
		}
		orginalEntitlment.SetConditions(conditions)
	}

	if d.HasChange("actions") {
		_, n := d.GetChange("actions")
		actions, err := readConditionActionsFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		orginalEntitlment.SetActions(actions)
	}

	if d.HasChange("app_shortcut") {
		_, n := d.GetChange("app_shortcut")
		appShortcut, err := readAppShortcutFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		orginalEntitlment.SetAppShortcuts(appShortcut)
	}

	req := api.EntitlementsIdPut(ctx, d.Id())
	req = req.Entitlement(orginalEntitlment)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update Entitlement %+v", prettyPrintAPIError(err))
	}

	return resourceAppgateEntitlementRuleRead(d, meta)
}

func resourceAppgateEntitlementRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete Entitlement: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Reading Entitlement id: %+v", d.Id())
	token := meta.(*Client).Token
	api := meta.(*Client).API.EntitlementsApi

	request := api.EntitlementsIdDelete(context.Background(), d.Id())

	_, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not delete Entitlement %+v", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}

func readConditionActionsFromConfig(actions []interface{}) ([]openapi.EntitlementAllOfActions, error) {
	result := make([]openapi.EntitlementAllOfActions, 0)
	for _, action := range actions {
		if action == nil {
			continue
		}
		a := openapi.NewEntitlementAllOfActionsWithDefaults()
		raw := action.(map[string]interface{})
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

		if v, ok := raw["monitor"]; ok {
			monitor := openapi.NewEntitlementAllOfMonitorWithDefaults()
			rawMonitors := v.(*schema.Set).List()
			for _, v := range rawMonitors {
				rawMonitor := v.(map[string]interface{})
				if v, ok := rawMonitor["enabled"]; ok {
					monitor.SetEnabled(v.(bool))
				}
				if v, ok := rawMonitor["timeout"]; ok {
					monitor.SetTimeout(int32(v.(int)))
				}
				a.SetMonitor(*monitor)
			}
		}
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
