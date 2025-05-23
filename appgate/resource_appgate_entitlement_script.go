package appgate

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v22/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateEntitlementScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateEntitlementScriptCreate,
		Read:   resourceAppgateEntitlementScriptRead,
		Update: resourceAppgateEntitlementScriptUpdate,
		Delete: resourceAppgateEntitlementScriptDelete,
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

			"entitlement_script_id": resourceUUID(),
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

			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"expression": {
				Type:        schema.TypeString,
				Description: "A JavaScript expression that returns a list of IPs and names.",
				Required:    true,
			},
		},
	}
}

func resourceAppgateEntitlementScriptCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating Entitlement script: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.EntitlementScriptsApi
	args := openapi.NewEntitlementScriptWithDefaults()
	if v, ok := d.GetOk("entitlement_script_id"); ok {
		args.SetId(v.(string))
	}
	args.SetName(d.Get("name").(string))
	args.SetNotes(d.Get("notes").(string))
	args.SetTags(schemaExtractTags(d))

	if v, ok := d.GetOk("expression"); ok {
		args.SetExpression(v.(string))
	}
	if v, ok := d.GetOk("type"); ok {
		args.SetType(v.(string))
	}
	ctx := BaseAuthContext(token)
	request := api.EntitlementScriptsPost(ctx)
	request = request.EntitlementScript(*args)
	EntitlementScript, _, err := request.Execute()
	if err != nil {
		return fmt.Errorf("Could not create Entitlement script %w", prettyPrintAPIError(err))
	}

	d.SetId(EntitlementScript.GetId())
	d.Set("entitlement_script_id", EntitlementScript.GetId())

	return resourceAppgateEntitlementScriptRead(d, meta)
}

func resourceAppgateEntitlementScriptRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Entitlement script id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.EntitlementScriptsApi
	ctx := BaseAuthContext(token)
	request := api.EntitlementScriptsIdGet(ctx, d.Id())
	EntitlementScript, res, err := request.Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Failed to read Entitlement script, %w", err)
	}
	d.SetId(EntitlementScript.GetId())
	d.Set("entitlement_script_id", EntitlementScript.GetId())
	d.Set("name", EntitlementScript.GetName())
	d.Set("notes", EntitlementScript.GetNotes())
	d.Set("tags", EntitlementScript.GetTags())
	d.Set("expression", EntitlementScript.GetExpression())
	d.Set("type", EntitlementScript.GetType())

	return nil
}

func resourceAppgateEntitlementScriptUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Entitlement script: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Updating Entitlement script id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.EntitlementScriptsApi
	ctx := BaseAuthContext(token)
	request := api.EntitlementScriptsIdGet(ctx, d.Id())
	originalEntitlementScript, _, err := request.Execute()
	if err != nil {
		return fmt.Errorf("Failed to read Entitlement script while updating, %w", err)
	}

	if d.HasChange("name") {
		originalEntitlementScript.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalEntitlementScript.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalEntitlementScript.SetTags(schemaExtractTags(d))
	}

	if d.HasChange("type") {
		originalEntitlementScript.SetType(d.Get("type").(string))
	}

	if d.HasChange("expression") {
		originalEntitlementScript.SetExpression(d.Get("expression").(string))
	}

	req := api.EntitlementScriptsIdPut(ctx, d.Id())
	req = req.EntitlementScript(*originalEntitlementScript)
	_, _, err = req.Execute()
	if err != nil {
		return fmt.Errorf("Could not update Entitlement script %w", prettyPrintAPIError(err))
	}
	return resourceAppgateEntitlementScriptRead(d, meta)
}

func resourceAppgateEntitlementScriptDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete Entitlement script: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Reading Entitlement script id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.EntitlementScriptsApi
	if _, err := api.EntitlementScriptsIdDelete(BaseAuthContext(token), d.Id()).Execute(); err != nil {
		return fmt.Errorf("Could not delete Entitlement script %w", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
