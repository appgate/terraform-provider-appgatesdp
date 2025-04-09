package appgate

import (
	"context"
	"fmt"
	"github.com/appgate/sdp-api-client-go/api/v22/openapi"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateClientProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppgateClientProfileCreate,
		ReadContext:   resourceAppgateClientProfileRead,
		UpdateContext: resourceAppgateClientProfileUpdate,
		DeleteContext: resourceAppgateClientProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"notes": {
				Type:        schema.TypeString,
				Description: "Notes for the object. Used for documentation purposes.",
				Optional:    true,
			},

			"tags": tagsSchema(),

			"spa_key_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"identity_provider_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"exported": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAppgateClientProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Create Client Profile %s", d.Get("name"))
	api := meta.(*Client).API.ClientProfilesApi
	args := make(map[string]interface{}, 0)
	args["name"] = d.Get("name").(string)
	args["notes"] = d.Get("notes").(string)
	args["tags"] = schemaExtractTags(d)
	if v, ok := d.GetOk("spa_key_name"); ok {
		args["spaKeyName"] = v.(string)
	}
	if v, ok := d.GetOk("identity_provider_name"); ok {
		args["identityProviderName"] = v.(string)
	}
	if _, ok := d.GetOk("url"); ok {
		return diag.Errorf("url is not supported on your appliance version, use hostname instead")
	}
	if v, ok := d.GetOk("hostname"); ok {
		args["hostname"] = v.(string)
	}
	profile, _, err := api.ClientProfilesPost(ctx).Authorization(token).Body(args).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Could not create client profile %s", prettyPrintAPIError(err)))
	}
	d.SetId(profile.GetId())
	return resourceAppgateClientProfileRead(ctx, d, meta)
}

func resourceAppgateClientProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Reading Client Profile id: %+v", d.Id())
	api := meta.(*Client).API.ClientProfilesApi
	profile, res, err := api.ClientProfilesIdGet(ctx, d.Id()).Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "profile has been removed remote",
			})
			return diags
		}
		return diag.FromErr(fmt.Errorf("Failed to read client profile, %s", err))
	}

	id, ok := profile["id"].(string)
	if ok {
		d.Set("id", id)
	}
	if name, ok := profile["name"].(string); ok {
		d.Set("name", name)
	}
	if notes, ok := profile["notes"].(string); ok {
		d.Set("notes", notes)
	}
	if spaKeyName, ok := profile["spa_key_name"].(string); ok {
		d.Set("spa_key_name", spaKeyName)
	}
	if identityProviderName, ok := profile["identity_provider_name"].(string); ok {
		d.Set("identity_provider_name", identityProviderName)
	}
	if hostname, ok := profile["hostname"].(string); ok {
		d.Set("hostname", hostname)
	}
	if exported, ok := profile["exported"].(string); ok {
		d.Set("exported", exported)
	}

	url, _, err := api.ClientProfilesIdUrlGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		diags = AppendFromErr(diags, err)
		return diags
	}
	d.Set("url", url.GetUrl())

	return nil
}

func resourceAppgateClientProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	currentVersion := meta.(*Client).ApplianceVersion
	if currentVersion.LessThan(Appliance61Version) {
		return nil
	}
	log.Printf("[DEBUG] Updating client profile id: %+v", d.Id())
	var diags diag.Diagnostics

	api := meta.(*Client).API.ClientProfilesApi
	originalProfile, _, err := api.ClientProfilesIdGet(ctx, d.Id()).Authorization(token).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to read profile while updating, %w", err))
	}

	if d.HasChange("name") {
		originalProfile["name"] = d.Get("name").(string)
	}
	if d.HasChange("notes") {
		originalProfile["notes"] = d.Get("notes").(string)
	}

	if d.HasChange("tags") {
		originalProfile["tags"] = schemaExtractTags(d)
	}
	if d.HasChange("spa_key_name") {
		originalProfile["spa_key_name"] = d.Get("spa_key_name").(string)
	}
	if d.HasChange("identity_provider_name") {
		originalProfile["identity_provider_name"] = d.Get("identity_provider_name").(string)
	}
	if d.HasChange("hostname") {
		originalProfile["hostname"] = d.Get("hostname").(string)
	}
	if _, _, err := api.ClientProfilesIdPut(ctx, d.Id()).Authorization(token).Body(originalProfile).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("Could not update client profile %w", prettyPrintAPIError(err)))

	}
	return diags
}

func resourceAppgateClientProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Delete client profile %+v", d.Id())
	api := meta.(*Client).API.ClientProfilesApi
	if _, err := api.ClientProfilesIdDelete(ctx, d.Id()).Authorization(token).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("Could not delete client profile %w", prettyPrintAPIError(err)))
	}
	d.SetId("")
	return nil
}
