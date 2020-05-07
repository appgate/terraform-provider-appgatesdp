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
					},
				},
			},

			"app_shortcut": {
				Type:       schema.TypeSet,
				Required:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"url": {
							Type:     schema.TypeString,
							Required: true,
						},

						"color_code": {
							Type:     schema.TypeInt,
							Required: true,
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
	args.SetConditionLogic(d.Get("condition_logic").(string))

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

	if c, ok := d.GetOk("app_shortcut"); ok {
		appShortcut, err := readAppShortcutFromConfig(c.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetAppShortcut(appShortcut)
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
	request := api.EntitlementsIdGet(context.Background(), d.Id())
	ent, _, err := request.Authorization(token).Execute()
	if err != nil {
		// TODO check if 404
		d.SetId("")
		return fmt.Errorf("Failed to read Entitlement, %+v", err)
	}
	d.SetId(ent.Id)
	d.Set("entitlement_id", ent.Id)
	d.Set("displayName", ent.DisplayName) // Deprecated in 5.1
	d.Set("disabled", ent.Disabled)
	d.Set("notes", ent.Notes)
	d.Set("actions", ent.Actions)
	d.Set("conditions", ent.Conditions)

	return nil
}

func resourceAppgateEntitlementRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Entitlement: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Updating Entitlement id: %+v", d.Id())

	token := meta.(*Client).Token
	api := meta.(*Client).API.EntitlementsApi

	request := api.EntitlementsIdGet(context.Background(), d.Id())
	request.Authorization(token)
	orginalEntitlment, _, err := request.Execute()
	if err != nil {
		return fmt.Errorf("Failed to read Entitlement, %+v", err)
	}

	if d.HasChange("name") {
		orginalEntitlment.SetName(d.Get("name").(string))
	}

	if d.HasChange("site") {
		orginalEntitlment.SetSite(d.Get("site").(string))
	}

	if d.HasChange("tags") {
		orginalEntitlment.SetTags(schemaExtractTags(d))
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

	req := api.EntitlementsIdPut(context.Background(), d.Id())
	req = req.Entitlement(orginalEntitlment)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to delete Entitlement, %+v", err)
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
		return fmt.Errorf("Failed to delete Entitlement, %+v", err)
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
			rawMonitor := v.(map[string]interface{})

			if v, ok := rawMonitor["enabled"]; ok {
				monitor.SetEnabled(v.(bool))
			}
			if v, ok := rawMonitor["timeout"]; ok {
				monitor.SetTimeout(int32(v.(int)))
			}
			a.SetMonitor(*monitor)
		}
		result = append(result, *a)
	}
	return result, nil
}

func readAppShortcutFromConfig(shortcuts []interface{}) (openapi.EntitlementAllOfAppShortcut, error) {
	result := openapi.EntitlementAllOfAppShortcut{}
	for _, shortcut := range shortcuts {
		if shortcut == nil {
			continue
		}
		raw := shortcut.(map[string]interface{})
		if v, ok := raw["name"]; ok {
			result.SetName(v.(string))
		}
		if v, ok := raw["url"]; ok {
			result.SetUrl(v.(string))
		}
		if v, ok := raw["color_code"]; ok {
			result.SetColorCode(int32(v.(int)))
		}
	}
	return result, nil
}
