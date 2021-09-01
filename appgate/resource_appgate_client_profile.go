package appgate

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v15/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateClientProfile() *schema.Resource {
	return &schema.Resource{

		Create: resourceAppgateClientProfileCreate,
		Read:   resourceAppgateClientProfileRead,
		Update: resourceAppgateClientProfileUpdate,
		Delete: resourceAppgateClientProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
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
			},

			"identity_provider_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getIdFromProfile(profile openapi.ClientConnectionsProfiles) string {
	// names are case sensitive and the controller only allows 1 per case type.
	// so it will be suitable as the identitfer.
	return profile.GetName()
}

func resourceAppgateClientProfileCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Create Client Profile %s", d.Get("name"))
	ctx := context.Background()

	api := meta.(*Client).API.ClientConnectionsApi

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
			return resource.RetryableError(err)
		}
		existingProfiles := clientConnections.GetProfiles()
		log.Printf("[DEBUG] Create  %s -- existing profiles count %d", d.Get("name"), len(existingProfiles))
		profile := openapi.ClientConnectionsProfiles{}
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
		_, response, err := api.ClientConnectionsPut(ctx).ClientConnections(clientConnections).Authorization(token).Execute()
		if err != nil {
			if response.StatusCode == http.StatusConflict {
				return resource.RetryableError(fmt.Errorf("Expected client profile %q to be created but was in state %s", d.Get("name").(string), response.Status))
			}
			return resource.NonRetryableError(fmt.Errorf("Error updating client connection profiles: %s", err))
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error create: %s", err)
	}
	return resourceAppgateClientProfileRead(d, meta)
}

func resourceAppgateClientProfileRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Client Profile id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.ClientConnectionsApi
	ctx := context.Background()
	clientConnections, _, err := api.ClientConnectionsGet(ctx).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not read Client Connections %+v", prettyPrintAPIError(err))
	}
	existingProfiles := clientConnections.GetProfiles()
	var p *openapi.ClientConnectionsProfiles
	for _, profile := range existingProfiles {
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
		return fmt.Errorf("could not find client profile %s during read", d.Id())
	}
	return nil
}

func resourceAppgateClientProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Client Profile %+v", d.Id())
	ctx := context.Background()
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.ClientConnectionsApi
	clientConnections, _, err := api.ClientConnectionsGet(ctx).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not read Client Connections %+v", prettyPrintAPIError(err))
	}
	existingProfiles := clientConnections.GetProfiles()
	var p *openapi.ClientConnectionsProfiles
	for _, profile := range existingProfiles {
		if strings.EqualFold(profile.GetName(), d.Id()) && profile.GetName() == d.Id() {
			p = &profile
			break
		}
	}
	if p == nil {
		diag.FromErr(fmt.Errorf("could not find client profile %s during update", d.Id()))
	}

	if d.HasChange("name") {
		p.SetName(d.Get("name").(string))
	}
	if d.HasChange("spa_key_name") {
		p.SetSpaKeyName(d.Get("spa_key_name").(string))
	}
	if d.HasChange("identity_provider_name") {
		p.SetIdentityProviderName(d.Get("identity_provider_name").(string))
	}
	req := api.ClientConnectionsPut(ctx)
	_, _, err = req.ClientConnections(clientConnections).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update Client profile %s %+v", d.Id(), prettyPrintAPIError(err))
	}
	return resourceAppgateClientProfileRead(d, meta)
}

func resourceAppgateClientProfileDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete client profile %+v", d.Id())
	ctx := context.Background()
	api := meta.(*Client).API.ClientConnectionsApi

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
			return resource.RetryableError(fmt.Errorf("Could not read Client Connections during delete %+v", prettyPrintAPIError(err)))
		}
		existingProfiles := clientConnections.GetProfiles()
		var p *openapi.ClientConnectionsProfiles
		var newProfiles []openapi.ClientConnectionsProfiles
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

		_, _, err = api.ClientConnectionsPut(ctx).ClientConnections(clientConnections).Authorization(token).Execute()
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Could not delete Client Connections after retry %+v", err)
	}

	d.SetId("")
	return nil
}
