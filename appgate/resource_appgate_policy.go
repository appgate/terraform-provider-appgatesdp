package appgate

import (
	"time"

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
				Required: true,
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
	return resourceAppgatePolicyRead(d, meta)
}

func resourceAppgatePolicyRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAppgatePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAppgatePolicyRead(d, meta)
}

func resourceAppgatePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
