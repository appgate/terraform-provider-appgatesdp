package appgate

import (
	"fmt"
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
