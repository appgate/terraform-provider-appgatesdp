package appgate

import (
	"time"

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
		},
	}
}

func resourceAppgateRingfenceRuleCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceAppgateRingfenceRuleRead(d, meta)
}
func resourceAppgateRingfenceRuleRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceAppgateRingfenceRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAppgateRingfenceRuleRead(d, meta)
}
func resourceAppgateRingfenceRuleDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
