package appgate

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v19/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

func getIdFromProfile(profile openapi.ClientConnectionsProfilesInner) string {
	// names are case sensitive and the controller only allows 1 per case type.
	// so it will be suitable as the identitfer.
	return profile.GetName()
}

func resourceAppgateClientProfileCreateLegacy(ctx context.Context, d *schema.ResourceData, meta interface{}, token string) diag.Diagnostics {
	api := meta.(*Client).API.ClientProfilesApi

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		rand.Seed(time.Now().UnixNano())
		duration := rand.Intn(10)
		time.Sleep(time.Duration(duration) * time.Second)

		// before we create a profile, make sure the controllers are in a healthy state
		if err := applianceStatsRetryable(ctx, meta); err != nil {
			return err
		}
		clientConnections, _, err := api.ClientConnectionsGet(ctx).Authorization(token).Execute()
		if err != nil {
			return resource.RetryableError(err)
		}
		existingProfiles := clientConnections.GetProfiles()

		profile := openapi.ClientConnectionsProfilesInner{}
		if v, ok := d.GetOk("name"); ok {
			profile.SetName(v.(string))
		}
		if v, ok := d.GetOk("spa_key_name"); ok {
			profile.SetSpaKeyName(v.(string))
		}
		if v, ok := d.GetOk("identity_provider_name"); ok {
			profile.SetIdentityProviderName(v.(string))
		}

		d.SetId(getIdFromProfile(profile))

		existingProfiles = append(existingProfiles, profile)
		clientConnections.SetProfiles(existingProfiles)
		_, _, err = api.ClientConnectionsPut(ctx).ClientConnections(*clientConnections).Authorization(token).Execute()
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error updating client connection profile %s: %w", d.Id(), err))
		}
		// check number of client profiles again and verify that is existingProfiles+1
		newConnections, _, err := api.ClientConnectionsGet(ctx).Authorization(token).Execute()
		if err != nil {
			return resource.RetryableError(err)
		}
		newProfiles := newConnections.GetProfiles()
		beep := false
		for _, existingProfile := range newProfiles {
			log.Printf("[DEBUG] Creation Read Found profile %q - Looking for %s", profile.GetName(), d.Id())
			if strings.EqualFold(existingProfile.GetName(), profile.GetName()) {
				log.Printf("[DEBUG] Found profile %s after create, OK!", profile.GetName())
				beep = true
			}
		}
		if !beep {
			return resource.RetryableError(fmt.Errorf("Profile %q did not get created", profile.GetName()))
		}
		// give the controller a moment before we check the initial status
		time.Sleep(time.Duration(duration) * time.Second)
		if err := waitForControllers(ctx, meta); err != nil {
			return resource.NonRetryableError(fmt.Errorf("1 or more controller never reached a healthy state after creating a client_profile: %w", err))
		}
		return nil
	})
	if err != nil {
		return diag.Errorf("Error create: %s", err)
	}
	return resourceAppgateClientProfileRead(ctx, d, meta)
}

func resourceAppgateClientProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	currentVersion := meta.(*Client).ApplianceVersion
	// starting from 6.1 we will use
	// /admin/client-profiles
	// instead of
	// /admin/client-connections
	if currentVersion.LessThan(Appliance61Version) {
		log.Printf("[DEBUG] Create Client Profile Legacy %s", d.Get("name"))
		return resourceAppgateClientProfileCreateLegacy(ctx, d, meta, token)
	}
	log.Printf("[DEBUG] Create Client Profile %s", d.Get("name"))
	api := meta.(*Client).API.ClientProfilesApi
	args := openapi.ClientProfile{}
	args.SetName(d.Get("name").(string))
	args.SetNotes(d.Get("notes").(string))
	args.SetTags(schemaExtractTags(d))
	if v, ok := d.GetOk("spa_key_name"); ok {
		args.SetSpaKeyName(v.(string))
	}
	if v, ok := d.GetOk("identity_provider_name"); ok {
		args.SetIdentityProviderName(v.(string))
	}
	if _, ok := d.GetOk("url"); ok {
		return diag.Errorf("url is not supported on your appliance version, use hostname instead")
	}
	if v, ok := d.GetOk("hostname"); ok {
		args.SetHostname(v.(string))
	}
	profile, _, err := api.ClientProfilesPost(ctx).Authorization(token).ClientProfile(args).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Could not create client profile %s", prettyPrintAPIError(err)))
	}
	d.SetId(profile.GetId())
	return resourceAppgateClientProfileRead(ctx, d, meta)
}

func resourceAppgateClientProfileReadLegacy(ctx context.Context, d *schema.ResourceData, meta interface{}, token string) diag.Diagnostics {
	api := meta.(*Client).API.ClientProfilesApi
	clientConnections, _, err := api.ClientConnectionsGet(ctx).Authorization(token).Execute()
	if err != nil {
		return diag.Errorf("Could not read Client Connections %s", prettyPrintAPIError(err))
	}
	existingProfiles := clientConnections.GetProfiles()
	var p *openapi.ClientConnectionsProfilesInner
	for _, profile := range existingProfiles {
		log.Printf("[DEBUG] Reading Found profile %q - Looking for %s", profile.GetName(), d.Id())
		if strings.EqualFold(profile.GetName(), d.Id()) && profile.GetName() == d.Id() {
			p = &profile
			d.Set("name", p.GetName())
			d.Set("spa_key_name", p.GetSpaKeyName())
			d.Set("identity_provider_name", p.GetIdentityProviderName())
			d.Set("url", p.GetUrl())
			break
		}
	}
	if p == nil {
		log.Printf("[DEBUG] Client Profile id %q not found in client connections profiles", d.Id())
		d.SetId("")
	}
	return nil
}

func resourceAppgateClientProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	currentVersion := meta.(*Client).ApplianceVersion
	if currentVersion.LessThan(Appliance61Version) {
		log.Printf("[DEBUG] Reading Client Profile Legacy id: %+v", d.Id())
		return resourceAppgateClientProfileReadLegacy(ctx, d, meta, token)
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
	d.Set("id", profile.GetId())
	d.Set("name", profile.GetName())
	d.Set("notes", profile.GetNotes())
	d.Set("spa_key_name", profile.GetSpaKeyName())
	d.Set("identity_provider_name", profile.GetIdentityProviderName())
	d.Set("hostname", profile.GetHostname())
	d.Set("exported", profile.GetExported().String())

	url, _, err := api.ClientProfilesIdUrlGet(ctx, profile.GetId()).Authorization(token).Execute()
	if err != nil {
		diags = AppendFromErr(diags, err)
		return diags
	}
	d.Set("url", url.GetUrl())

	return nil
}

func resourceAppgateClientProfileDeleteLegacy(ctx context.Context, d *schema.ResourceData, meta interface{}, token string) diag.Diagnostics {
	api := meta.(*Client).API.ClientProfilesApi

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		rand.Seed(time.Now().UnixNano())
		duration := rand.Intn(10) // n will be between 0 and 20
		log.Printf("[DEBUG] Create Client Profile %s Sleep %d", d.Get("name"), duration)
		time.Sleep(time.Duration(duration) * time.Second)
		if err := applianceStatsRetryable(ctx, meta); err != nil {
			return err
		}

		clientConnections, _, err := api.ClientConnectionsGet(ctx).Authorization(token).Execute()
		if err != nil {
			return resource.RetryableError(fmt.Errorf("Could not read Client Connections during delete %w", prettyPrintAPIError(err)))
		}
		existingProfiles := clientConnections.GetProfiles()
		var p *openapi.ClientConnectionsProfilesInner
		var newProfiles []openapi.ClientConnectionsProfilesInner
		for i, profile := range existingProfiles {
			if strings.EqualFold(profile.GetName(), d.Id()) && profile.GetName() == d.Id() {
				p = &profile
				// remove the profile from the list and maintain order.
				newProfiles = append(existingProfiles[:i], existingProfiles[i+1:]...)
				break
			}
		}
		if p == nil {
			diag.FromErr(fmt.Errorf("could not find client profile %s during delete", d.Id()))
		}
		clientConnections.SetProfiles(newProfiles)
		_, _, err = api.ClientConnectionsPut(ctx).ClientConnections(*clientConnections).Authorization(token).Execute()
		if err != nil {
			return resource.NonRetryableError(err)
		}
		// give the controller a moment before we check the initial status
		time.Sleep(time.Duration(duration) * time.Second)
		if err := waitForControllers(ctx, meta); err != nil {
			return resource.NonRetryableError(fmt.Errorf("1 or more controller never reached a healthy state after deleting a client_profile: %w", err))
		}
		return nil
	})
	if err != nil {
		return diag.Errorf("Could not delete Client profile %s after retry %s", d.Id(), err)
	}

	d.SetId("")
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
	orginalProfile, _, err := api.ClientProfilesIdGet(ctx, d.Id()).Authorization(token).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to read profile while updating, %w", err))
	}

	if d.HasChange("name") {
		orginalProfile.SetName(d.Get("name").(string))
	}
	if d.HasChange("notes") {
		orginalProfile.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		orginalProfile.SetTags(schemaExtractTags(d))
	}
	if d.HasChange("spa_key_name") {
		orginalProfile.SetSpaKeyName(d.Get("spa_key_name").(string))
	}
	if d.HasChange("identity_provider_name") {
		orginalProfile.SetIdentityProviderName(d.Get("identity_provider_name").(string))
	}
	if d.HasChange("hostname") {
		orginalProfile.SetHostname(d.Get("hostname").(string))
	}
	if _, _, err := api.ClientProfilesIdPut(ctx, d.Id()).Authorization(token).ClientProfile(*orginalProfile).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("Could not update client profile %w", prettyPrintAPIError(err)))

	}
	return diags
}

func resourceAppgateClientProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	currentVersion := meta.(*Client).ApplianceVersion
	if currentVersion.LessThan(Appliance61Version) {
		log.Printf("[DEBUG] Delete client profile Legacy %+v", d.Id())
		return resourceAppgateClientProfileDeleteLegacy(ctx, d, meta, token)
	}
	log.Printf("[DEBUG] Delete client profile %+v", d.Id())
	api := meta.(*Client).API.ClientProfilesApi
	if _, err := api.ClientProfilesIdDelete(ctx, d.Id()).Authorization(token).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("Could not delete client profile %w", prettyPrintAPIError(err)))
	}
	d.SetId("")
	return nil
}
