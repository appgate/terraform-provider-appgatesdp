package appgate

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAppgateCondition() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateConditionCreate,
		Read:   resourceAppgateConditionRead,
		Update: resourceAppgateConditionUpdate,
		Delete: resourceAppgateConditionDelete,
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

			"condition_id": {
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

			"expression": {
				Type:        schema.TypeString,
				Description: "Name of the object.",
				Required:    true,
			},

			"repeat_schedules": {
				Type:        schema.TypeSet,
				Description: "Array of tags.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"remedy_methods": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"type": {
							Type:     schema.TypeString,
							Required: true,
						},

						"message": {
							Type:     schema.TypeString,
							Required: true,
						},

						"claimn_suffix": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"provider_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAppgateConditionCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceAppgateConditionRead(d, meta)
}
func resourceAppgateConditionRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceAppgateConditionUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAppgateConditionRead(d, meta)
}
func resourceAppgateConditionDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
