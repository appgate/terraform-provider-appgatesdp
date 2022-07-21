package appgate

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/appgate/sdp-api-client-go/api/v16/openapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateUserClaimScript() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppgateUserClaimScriptCreate,
		ReadContext:   resourceAppgateUserClaimScriptRead,
		UpdateContext: resourceAppgateUserClaimScriptUpdate,
		DeleteContext: resourceAppgateUserClaimScriptDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{

			"user_claim_script_id": resourceUUID(),

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
				Description: "The User Claim Script content.",
				Optional:    true,
			},
		},
	}
}

func resourceAppgateUserClaimScriptCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating User Claim Script: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.UserClaimScriptsApi
	args := openapi.NewUserScriptWithDefaults()
	if v, ok := d.GetOk("user_claim_script_id"); ok {
		args.SetId(v.(string))
	}
	args.SetName(d.Get("name").(string))
	args.SetNotes(d.Get("notes").(string))
	args.SetTags(schemaExtractTags(d))
	if v, ok := d.GetOk("expression"); ok {
		args.SetExpression(v.(string))
	}

	UserClaimScript, _, err := api.UserScriptsPost(ctx).UserScript(*args).Authorization(token).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Could not create User Claim Script %w", prettyPrintAPIError(err)))
	}

	d.SetId(UserClaimScript.Id)
	d.Set("user_claim_script_id", UserClaimScript.Id)

	return resourceAppgateUserClaimScriptRead(ctx, d, meta)
}

func resourceAppgateUserClaimScriptRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading User Claim Script id: %s", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.UserClaimScriptsApi
	request := api.UserScriptsIdGet(ctx, d.Id())
	UserClaimScript, res, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.FromErr(fmt.Errorf("Failed to read User claim script, %w", err))
	}
	d.SetId(UserClaimScript.Id)
	d.Set("user_claim_script_id", UserClaimScript.GetId())
	d.Set("name", UserClaimScript.GetName())
	d.Set("notes", UserClaimScript.GetNotes())
	d.Set("tags", UserClaimScript.GetTags())
	d.Set("expression", UserClaimScript.GetExpression())

	return nil
}

func resourceAppgateUserClaimScriptUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating User Claim Script: %s", d.Get("name").(string))

	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to read user claim script %w", err))
	}
	api := meta.(*Client).API.UserClaimScriptsApi
	request := api.UserScriptsIdGet(ctx, d.Id())
	originalUserClaimScript, _, err := request.Authorization(token).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to read User Claim Script while updating, %w", err))
	}

	if d.HasChange("name") {
		originalUserClaimScript.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalUserClaimScript.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalUserClaimScript.SetTags(schemaExtractTags(d))
	}

	if d.HasChange("expression") {
		originalUserClaimScript.SetExpression(d.Get("expression").(string))
	}

	req := api.UserScriptsIdPut(ctx, d.Id())
	req = req.UserScript(originalUserClaimScript)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Could not update User Claim Script %w", prettyPrintAPIError(err)))
	}
	return resourceAppgateUserClaimScriptRead(ctx, d, meta)
}

func resourceAppgateUserClaimScriptDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete User Claim Script: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Reading User Claim Script id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.UserClaimScriptsApi

	if _, err := api.UserScriptsIdDelete(ctx, d.Id()).Authorization(token).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("Could not delete User Claim Script %w", prettyPrintAPIError(err)))
	}
	d.SetId("")
	return nil
}
