package appgate

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v21/openapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// errRemedyLogicUnsupportedVersion is used when trying to use remedy_logic on an older unsupported version.
var errRemedyLogicUnsupportedVersion = fmt.Errorf("remedy_logic is only supported in %s or higher", ApplianceVersionMap[Version14])

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

			"condition_id": resourceUUID(),

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

			"tags": tagsSchema(),

			"expression": {
				Type:        schema.TypeString,
				Description: "Boolean expression in JavaScript.",
				Required:    true,
			},

			"repeat_schedules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"remedy_logic": {
				Type:        schema.TypeString,
				Description: "Whether all the Remedy Methods must succeed to pass this Condition or just one.",
				Optional:    true,
				Computed:    true,
				ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
					s := v.(string)
					list := []string{"and", "or"}
					for _, x := range list {
						if s == x {
							return
						}
					}
					errs = append(errs, fmt.Errorf("remedy_logic must be on of %v, got %s", list, s))
					return
				},
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
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.ConditionsApi
	currentVersion := meta.(*Client).ApplianceVersion

	args := openapi.Condition{}
	if v, ok := d.GetOk("condition_id"); ok {
		args.SetId(v.(string))
	}
	args.SetName(d.Get("name").(string))

	if c, ok := d.GetOk("notes"); ok {
		args.SetNotes(c.(string))
	}

	args.SetTags(schemaExtractTags(d))

	if v, ok := d.GetOk("expression"); ok {
		args.SetExpression(v.(string))
	}

	if v, ok := d.GetOk("remedy_logic"); ok {
		if currentVersion.LessThan(Appliance53Version) {
			return fmt.Errorf("%w, you are using %q client v%d", errRemedyLogicUnsupportedVersion, currentVersion, meta.(*Client).ClientVersion)
		}
		args.SetRemedyLogic(v.(string))
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
	request = request.Condition(args)
	condition, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create condition %w", prettyPrintAPIError(err))
	}

	d.SetId(condition.GetId())

	return resourceAppgateConditionRead(d, meta)
}

func resourceAppgateConditionRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Condition Name: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.ConditionsApi
	currentVersion := meta.(*Client).ApplianceVersion
	ctx := context.Background()
	request := api.ConditionsIdGet(ctx, d.Id())
	remoteCondition, res, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Failed to read Condition, %w", err)
	}
	d.SetId(remoteCondition.GetId())
	d.Set("condition_id", remoteCondition.Id)
	d.Set("name", remoteCondition.Name)
	d.Set("notes", remoteCondition.Notes)
	d.Set("tags", remoteCondition.Tags)
	d.Set("expression", remoteCondition.Expression)

	if currentVersion.GreaterThanOrEqual(Appliance53Version) {
		d.Set("remedy_logic", remoteCondition.GetRemedyLogic())
	}

	d.Set("repeat_schedules", remoteCondition.RepeatSchedules)
	if remoteCondition.RemedyMethods != nil {
		if err = d.Set("remedy_methods", flattenConditionRemedyMethods(remoteCondition.RemedyMethods)); err != nil {
			return err
		}
	}
	return nil
}

func flattenConditionRemedyMethods(in []openapi.RemedyMethod) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["type"] = v.Type
		m["message"] = v.Message
		m["claim_suffix"] = v.ClaimSuffix
		m["provider_id"] = v.ProviderId

		out[i] = m
	}
	return out
}

func resourceAppgateConditionUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating condition: %s", d.Get("name").(string))
	ctx := context.Background()
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.ConditionsApi
	currentVersion := meta.(*Client).ApplianceVersion
	request := api.ConditionsIdGet(ctx, d.Id())
	orginalCondition, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read condition, %w", err)
	}
	if d.HasChange("name") {
		orginalCondition.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		orginalCondition.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		orginalCondition.SetTags(schemaExtractTags(d))
	}

	if d.HasChange("expression") {
		orginalCondition.SetExpression(d.Get("expression").(string))
	}
	if d.HasChange("remedy_logic") {
		if currentVersion.LessThan(Appliance53Version) {
			return fmt.Errorf("%w, you are using %q client v%d", errRemedyLogicUnsupportedVersion, currentVersion, meta.(*Client).ClientVersion)
		}
		orginalCondition.SetRemedyLogic(d.Get("remedy_logic").(string))
	}

	if d.HasChange("repeat_schedules") {
		_, n := d.GetChange("repeat_schedules")
		repeatSchedules, err := readArrayOfStringsFromConfig(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		orginalCondition.SetRepeatSchedules(repeatSchedules)
	}

	if d.HasChange("remedy_methods") {
		_, n := d.GetChange("remedy_methods")
		remedyMethods, err := readRemedyMethodsFromConfig(n.([]interface{}))
		if err != nil {
			return err
		}
		orginalCondition.SetRemedyMethods(remedyMethods)
	}

	req := api.ConditionsIdPut(ctx, d.Id())

	_, _, err = req.Condition(*orginalCondition).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update condition %w", prettyPrintAPIError(err))
	}

	return resourceAppgateConditionRead(d, meta)
}

func resourceAppgateConditionDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete condition with name: %s", d.Get("name").(string))
	ctx := context.Background()
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.ConditionsApi

	// Get condition
	request := api.ConditionsIdGet(ctx, d.Id())
	condition, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to delete condition while GET, %w", err)
	}

	deleteRequest := api.ConditionsIdDelete(ctx, condition.GetId())
	_, err = deleteRequest.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to delete condition, %w", err)
	}
	d.SetId("")
	return nil
}

func readRemedyMethodsFromConfig(methods []interface{}) ([]openapi.RemedyMethod, error) {
	result := make([]openapi.RemedyMethod, 0)
	for _, method := range methods {
		if method == nil {
			continue
		}
		r := openapi.RemedyMethod{}
		raw := method.(map[string]interface{})
		if v, ok := raw["type"]; ok {
			r.SetType(v.(string))
		}

		if v, ok := raw["message"]; ok {
			r.SetMessage(v.(string))
		}

		if v, ok := raw["claim_suffix"]; ok && len(v.(string)) > 0 {
			r.SetClaimSuffix(v.(string))
		}

		if v, ok := raw["provider_id"]; ok && len(v.(string)) > 0 {
			r.SetProviderId(v.(string))
		}
		result = append(result, r)
	}
	return result, nil
}
