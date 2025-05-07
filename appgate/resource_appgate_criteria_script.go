package appgate

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v22/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateCriteriaScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateCriteriaScriptCreate,
		Read:   resourceAppgateCriteriaScriptRead,
		Update: resourceAppgateCriteriaScriptUpdate,
		Delete: resourceAppgateCriteriaScriptDelete,
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

			"criteria_script_id": resourceUUID(),

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
				Description: "A JavaScript expression that returns boolean.",
				Required:    true,
			},
		},
	}
}

func resourceAppgateCriteriaScriptCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating Criteria script: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.CriteriaScriptsApi
	args := openapi.NewCriteriaScriptWithDefaults()
	if v, ok := d.GetOk("criteria_script_id"); ok {
		args.SetId(v.(string))
	}
	args.SetName(d.Get("name").(string))
	args.SetNotes(d.Get("notes").(string))
	args.SetTags(schemaExtractTags(d))

	if v, ok := d.GetOk("expression"); ok {
		args.SetExpression(v.(string))
	}

	ctx := BaseAuthContext(token)
	request := api.CriteriaScriptsPost(ctx)
	request = request.CriteriaScript(*args)
	criteraScript, _, err := request.Execute()
	if err != nil {
		return fmt.Errorf("Could not create Criteria script %w", prettyPrintAPIError(err))
	}

	d.SetId(criteraScript.GetId())
	d.Set("criteria_script_id", criteraScript.GetId())

	return resourceAppgateCriteriaScriptRead(d, meta)
}

func resourceAppgateCriteriaScriptRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Criteria script id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.CriteriaScriptsApi
	ctx := BaseAuthContext(token)
	request := api.CriteriaScriptsIdGet(ctx, d.Id())
	criteraScript, res, err := request.Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Failed to read Criteria script, %w", err)
	}
	d.SetId(criteraScript.GetId())
	d.Set("criteria_script_id", criteraScript.GetId())
	d.Set("name", criteraScript.GetName())
	d.Set("notes", criteraScript.GetNotes())
	d.Set("tags", criteraScript.GetTags())
	d.Set("expression", criteraScript.GetExpression())

	return nil
}

func resourceAppgateCriteriaScriptUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Criteria script: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Updating Criteria script id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.CriteriaScriptsApi
	ctx := BaseAuthContext(token)
	request := api.CriteriaScriptsIdGet(ctx, d.Id())
	originalCriteriaScript, _, err := request.Execute()
	if err != nil {
		return fmt.Errorf("Failed to read Criteria script while updating, %w", err)
	}

	if d.HasChange("name") {
		originalCriteriaScript.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalCriteriaScript.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalCriteriaScript.SetTags(schemaExtractTags(d))
	}

	if d.HasChange("expression") {
		originalCriteriaScript.SetExpression(d.Get("expression").(string))
	}

	req := api.CriteriaScriptsIdPut(ctx, d.Id())
	req = req.CriteriaScript(*originalCriteriaScript)
	_, _, err = req.Execute()
	if err != nil {
		return fmt.Errorf("Could not update Criteria script %w", prettyPrintAPIError(err))
	}
	return resourceAppgateCriteriaScriptRead(d, meta)
}

func resourceAppgateCriteriaScriptDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete Criteria script: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Reading Criteria script id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.CriteriaScriptsApi
	if _, err := api.CriteriaScriptsIdDelete(BaseAuthContext(token), d.Id()).Execute(); err != nil {
		return fmt.Errorf("Could not delete Criteria script %w", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
