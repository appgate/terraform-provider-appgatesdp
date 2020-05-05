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
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"remedy_methods": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
								s := v.(string)
								list := []string{"DisplayMessage", "OtpAuthentication", "PasswordAuthentication", "Reason"}
								for _, x := range list {
									if s == x {
										return
									}
								}
								errs = append(errs, fmt.Errorf("type must be on of %v, got %s", list, s))
								return
							},
						},

						"message": {
							Type:     schema.TypeString,
							Required: true,
						},

						"claim_suffix": {
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
	log.Printf("[DEBUG] Creating Condition with name: %s", d.Get("name").(string))
	ctx := context.Background()
	token := meta.(*Client).Token
	api := meta.(*Client).API.ConditionsApi

	args := openapi.NewConditionWithDefaults()
	args.Id = uuid.New().String()
	args.SetName(d.Get("name").(string))

	if c, ok := d.GetOk("notes"); ok {
		args.SetNotes(c.(string))
	}

	args.SetTags(schemaExtractTags(d))

	if c, ok := d.GetOk("expression"); ok {
		args.SetExpression(c.(string))
	}

	if c, ok := d.GetOk("repeat_schedules"); ok {
		repeatSchedules, err := readArrayOfStringsFromConfig(c.(*schema.Set).List())
		if err != nil {
			return err
		}
		args.SetRepeatSchedules(repeatSchedules)
	}

	if v, ok := d.GetOk("remedy_methods"); ok {
		remedyMethods, err := readRemedyMethodsFromConfig(v.([]interface{}))
		if err != nil {
			return err
		}
		args.SetRemedyMethods(remedyMethods)
	}

	request := api.ConditionsPost(ctx)
	request = request.Condition(*args)
	condition, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create condition %+v", prettyPrintAPIError(err))
	}

	d.SetId(condition.Id)

	return resourceAppgateConditionRead(d, meta)
}
func resourceAppgateConditionRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceAppgateConditionUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAppgateConditionRead(d, meta)
}
func resourceAppgateConditionDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete condition with name: %s", d.Get("name").(string))
	ctx := context.Background()
	token := meta.(*Client).Token
	api := meta.(*Client).API.ConditionsApi

	// Get condition
	request := api.ConditionsIdGet(ctx, d.Id())
	condition, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to delete condition while GET, %+v", err)
	}

	deleteRequest := api.ConditionsIdDelete(ctx, condition.GetId())
	_, err = deleteRequest.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to delete condition, %+v", err)
	}
	d.SetId("")
	return nil
}

func readRemedyMethodsFromConfig(methods []interface{}) ([]openapi.ConditionAllOfRemedyMethods, error) {
	result := make([]openapi.ConditionAllOfRemedyMethods, 0)
	for _, method := range methods {
		if method == nil {
			continue
		}
		r := openapi.ConditionAllOfRemedyMethods{}
		raw := method.(map[string]interface{})

		if v, ok := raw["type"]; ok {
			r.SetType(v.(string))
		}

		if v, ok := raw["message"]; ok {
			r.SetMessage(v.(string))
		}

		if v, ok := raw["claim_suffix"]; ok {
			r.SetClaimSuffix(v.(string))
		}

		if v, ok := raw["provider_id"]; ok {
			r.SetProviderId(v.(string))
		}
	}
	return result, nil
}
