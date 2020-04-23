package appgate

import (
	"context"
	"fmt"
	"log"

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
		// TODO add optional fields.
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

			"site": {
				Type:        schema.TypeString,
				Description: "ID of the site for this Entitlement.",
				Required:    true,
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
						// Required for ICMP
						"types": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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

	rawconditions := d.Get("conditions").(*schema.Set).List()
	conditions := make([]string, 0)
	for _, raw := range rawconditions {
		conditions = append(conditions, raw.(string))
	}
	rawactions := d.Get("actions").(*schema.Set).List()

	var actions []openapi.EntitlementAllOfActions
	for _, raw := range rawactions {
		l := raw.(map[string]interface{})

		action := openapi.EntitlementAllOfActions{
			Subtype: l["subtype"].(string),
			Action:  l["action"].(string),
		}

		types := make([]string, 0)
		for _, t := range l["types"].([]interface{}) {
			types = append(types, t.(string))
		}
		action.Types = &types

		rawhosts := l["hosts"].([]interface{})
		hosts := make([]string, 0)
		for _, h := range rawhosts {
			hosts = append(hosts, h.(string))
		}
		action.Hosts = hosts

		actions = append(actions, action)
	}
	args := openapi.NewEntitlementWithDefaults()
	args.SetId(uuid.New().String())
	args.SetName(d.Get("name").(string))
	args.SetSite(d.Get("site").(string))
	args.SetConditions(conditions)
	args.SetActions(actions)

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

	rawconditions := d.Get("conditions").(*schema.Set).List()
	conditions := make([]string, 0)
	for _, raw := range rawconditions {
		conditions = append(conditions, raw.(string))
	}
	rawactions := d.Get("actions").(*schema.Set).List()

	var actions []openapi.EntitlementAllOfActions
	for _, raw := range rawactions {
		l := raw.(map[string]interface{})

		action := openapi.EntitlementAllOfActions{
			Subtype: l["subtype"].(string),
			Action:  l["action"].(string),
		}

		types := make([]string, 0)
		for _, t := range l["types"].([]interface{}) {
			types = append(types, t.(string))
		}
		action.Types = &types

		rawhosts := l["hosts"].([]interface{})
		hosts := make([]string, 0)
		for _, h := range rawhosts {
			hosts = append(hosts, h.(string))
		}
		action.Hosts = hosts

		actions = append(actions, action)
	}
	orginalEntitlment.SetName(d.Get("name").(string))
	orginalEntitlment.SetSite(d.Get("site").(string))
	orginalEntitlment.SetConditions(conditions)
	orginalEntitlment.SetActions(actions)

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
