package appgate

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/appgate/terraform-provider-appgate/client/v12/openapi"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAppgateRingfenceRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateRingfenceRuleCreate,
		Read:   resourceAppgateRingfenceRuleRead,
		Update: resourceAppgateRingfenceRuleUpdate,
		Delete: resourceAppgateRingfenceRuleDelete,
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

			"ringfence_rule_id": {
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

			"actions": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
								s := v.(string)
								list := []string{"icmp", "icmpv6", "udp", "tcp"}
								for _, x := range list {
									if s == x {
										return
									}
								}
								errs = append(errs, fmt.Errorf("protocol must be on of %v, got %s", list, s))
								return
							},
						},

						"direction": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
								s := v.(string)
								list := []string{"up", "down"}
								for _, x := range list {
									if s == x {
										return
									}
								}
								errs = append(errs, fmt.Errorf("direction must be on of %v, got %s", list, s))
								return
							},
						},

						"action": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
								s := v.(string)
								list := []string{"allow", "block"}
								for _, x := range list {
									if s == x {
										return
									}
								}
								errs = append(errs, fmt.Errorf("action must be on of %v, got %s", list, s))
								return
							},
						},

						"hosts": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"ports": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"types": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceAppgateRingfenceRuleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating Ringfence rule with name: %s", d.Get("name").(string))
	ctx := context.Background()
	token := meta.(*Client).Token
	api := meta.(*Client).API.RingfenceRulesApi

	args := openapi.NewRingfenceRuleWithDefaults()
	args.SetName(d.Get("name").(string))

	if c, ok := d.GetOk("notes"); ok {
		args.SetNotes(c.(string))
	}
	args.SetTags(schemaExtractTags(d))

	if c, ok := d.GetOk("actions"); ok {
		action, err := readRingfencActionFromConfig(c.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetActions(action)
	}

	request := api.RingfenceRulesPost(ctx)
	request = request.RingfenceRule(*args)
	ringfenceRule, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create Ringfence rule  %+v", prettyPrintAPIError(err))
	}

	d.SetId(ringfenceRule.Id)
	return resourceAppgateRingfenceRuleRead(d, meta)
}

func resourceAppgateRingfenceRuleRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAppgateRingfenceRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Ringfence rule with name: %s", d.Get("name").(string))
	ctx := context.Background()
	token := meta.(*Client).Token
	api := meta.(*Client).API.RingfenceRulesApi
	request := api.RingfenceRulesIdGet(ctx, d.Id())
	originalRingfenceRule, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read Ringfence rule, %+v", err)
	}

	if d.HasChange("name") {
		originalRingfenceRule.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalRingfenceRule.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalRingfenceRule.SetTags(schemaExtractTags(d))
	}
	if d.HasChange("actions") {
		_, n := d.GetChange("actions")
		actions, err := readRingfencActionFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		originalRingfenceRule.SetActions(actions)
	}
	req := api.RingfenceRulesIdPut(ctx, d.Id())

	_, _, err = req.RingfenceRule(originalRingfenceRule).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update Ringfence rule %+v", prettyPrintAPIError(err))
	}

	return resourceAppgateRingfenceRuleRead(d, meta)
}

func resourceAppgateRingfenceRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete Ringfence rule: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.RingfenceRulesApi
	ctx := context.Background()

	request := api.RingfenceRulesIdGet(ctx, d.Id())
	ringfenceRule, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to delete Ringfence rule while GET, %+v", err)
	}

	deleteRequest := api.RingfenceRulesIdDelete(ctx, ringfenceRule.GetId())
	_, err = deleteRequest.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to delete Ringfence rule, %+v", err)
	}
	d.SetId("")
	return nil
}

func readRingfencActionFromConfig(actions []interface{}) (openapi.RingfenceRuleAllOfActions, error) {
	result := openapi.RingfenceRuleAllOfActions{}
	for _, action := range actions {
		if action == nil {
			continue
		}
		raw := action.(map[string]interface{})
		if v, ok := raw["protocol"]; ok {
			result.SetProtocol(v.(string))
		}
		if v, ok := raw["direction"]; ok {
			result.SetDirection(v.(string))
		}
		if v, ok := raw["action"]; ok {
			result.SetAction(v.(string))
		}
		if v := raw["hosts"]; len(v.([]interface{})) > 0 {
			hosts, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, err
			}
			result.SetHosts(hosts)
		}
		if v := raw["ports"]; len(v.([]interface{})) > 0 {
			ports, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, err
			}
			result.SetPorts(ports)
		}
		if v := raw["types"]; len(v.([]interface{})) > 0 {
			types, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, err
			}
			result.SetTypes(types)
		}

	}
	return result, nil
}
