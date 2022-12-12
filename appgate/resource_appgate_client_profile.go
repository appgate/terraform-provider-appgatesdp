package appgate

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v18/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateClientProfile() *schema.Resource {
	return &schema.Resource{

		Create: resourceAppgateClientProfileCreate,
		Read:   resourceAppgateClientProfileRead,
		Delete: resourceAppgateClientProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

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
		},
	}
}

func getIdFromProfile(profile openapi.ClientConnectionsProfilesInner) string {
	// names are case sensitive and the controller only allows 1 per case type.
	// so it will be suitable as the identitfer.
	return profile.GetName()
}

func resourceAppgateClientProfileCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Create Client Profile %s", d.Get("name"))
	ctx := context.Background()
	api := meta.(*Client).API.ClientProfilesApi

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		rand.Seed(time.Now().UnixNano())
		duration := rand.Intn(10)
		time.Sleep(time.Duration(duration) * time.Second)

		// before we create a profile, make sure the controllers are in a healthy state
		if err := applianceStatsRetryable(ctx, meta); err != nil {
			return err
		}
		token, err := meta.(*Client).GetToken()
		if err != nil {
			return resource.RetryableError(err)
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
		return fmt.Errorf("Error create: %w", err)
	}
	return resourceAppgateClientProfileRead(d, meta)
}

func resourceAppgateClientProfileRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Client Profile id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.ClientProfilesApi
	ctx := context.Background()
	clientConnections, _, err := api.ClientConnectionsGet(ctx).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not read Client Connections %w", prettyPrintAPIError(err))
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

func resourceAppgateClientProfileDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete client profile %+v", d.Id())
	ctx := context.Background()
	api := meta.(*Client).API.ClientProfilesApi

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		rand.Seed(time.Now().UnixNano())
		duration := rand.Intn(10) // n will be between 0 and 20
		log.Printf("[DEBUG] Create Client Profile %s Sleep %d", d.Get("name"), duration)
		time.Sleep(time.Duration(duration) * time.Second)
		if err := applianceStatsRetryable(ctx, meta); err != nil {
			return err
		}
		token, err := meta.(*Client).GetToken()
		if err != nil {
			return resource.RetryableError(err)
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
		return fmt.Errorf("Could not delete Client profile %s after retry %w", d.Id(), err)
	}

	d.SetId("")
	return nil
}
