package appgate

import (
	"context"
	"fmt"
	"log"

	"github.com/appgate/sdp-api-client-go/api/v16/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceClientConnections() *schema.Resource {
	return &schema.Resource{

		DeprecationMessage: "Deprecated resource, replaced by appgatesdp_client_profile",

		Create: resourceClientConnectionsCreate,
		Read:   resourceClientConnectionsRead,
		Update: resourceClientConnectionsUpdate,
		Delete: resourceClientConnectionsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"spa_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
					s := v.(string)
					list := []string{"Disabled", "TCP", "UDP-TCP"}
					for _, x := range list {
						if s == x {
							return
						}
					}
					errs = append(errs, fmt.Errorf("spa_mode must be on of %v, got %s", list, s))
					return
				},
			},
			"profiles": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"name": {
							Type:     schema.TypeString,
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
				},
			},
		},
	}
}

func resourceClientConnectionsCreate(d *schema.ResourceData, meta interface{}) error {
	//TODO: Fix function, as it causes overwrites each run & drops all the SPA Keys
	return resourceClientConnectionsUpdate(d, meta)
}

func resourceClientConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Client Connections id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.ClientConnectionsApi
	ctx := context.TODO()
	request := api.ClientConnectionsGet(ctx)
	clientConnections, _, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Failed to read Client Connections, %+v", err)
	}
	d.SetId("spa_mode")
	if v, o := clientConnections.GetSpaModeOk(); o {
		d.Set("spa_mode", v)
	}
	if profiles, o := clientConnections.GetProfilesOk(); o {
		flattenProfiles := make([]map[string]interface{}, 0)
		for _, p := range *profiles {
			profile := make(map[string]interface{})
			if v, o := p.GetNameOk(); o {
				profile["name"] = *v
			}
			if v, o := p.GetSpaKeyNameOk(); o {
				profile["spa_key_name"] = *v
			}
			if v, o := p.GetIdentityProviderNameOk(); o {
				profile["identity_provider_name"] = *v
			}
			if v, o := p.GetUrlOk(); o {
				profile["url"] = *v
			}
			flattenProfiles = append(flattenProfiles, profile)
		}
		d.Set("profiles", flattenProfiles)
	}
	return nil
}

func resourceClientConnectionsUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Client Connections")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.ClientConnectionsApi
	ctx := context.TODO()
	request := api.ClientConnectionsGet(ctx)
	originalclientConnections, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read Client Connections while updating, %+v", err)
	}
	d.SetId("client_connections")

	if d.HasChange("spa_mode") {
		originalclientConnections.SetSpaMode(d.Get("spa_mode").(string))
	}
	if d.HasChange("profiles") {
		_, v := d.GetChange("profiles")
		profiles := readClientConnectionProfilesFromConfig(v.([]interface{}))
		log.Printf("[DEBUG] Updating Client PROFILES SET %+v", profiles)
		if err != nil {
			return fmt.Errorf("Failed to read profiles %s", err)
		}
		originalclientConnections.SetProfiles(profiles)
	}

	log.Printf("[DEBUG] Updating Client Connections %+v", originalclientConnections)
	req := api.ClientConnectionsPut(ctx)
	_, _, err = req.ClientConnections(originalclientConnections).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update Client Connections %+v", prettyPrintAPIError(err))
	}

	return resourceClientConnectionsRead(d, meta)
}

func readClientConnectionProfilesFromConfig(input []interface{}) []openapi.ClientConnectionsProfiles {
	result := make([]openapi.ClientConnectionsProfiles, 0)
	for _, p := range input {
		rawProfile := p.(map[string]interface{})
		log.Printf("[DEBUG] Updating RAW PROFILE %+v", rawProfile)
		profile := openapi.ClientConnectionsProfiles{}
		if v, o := rawProfile["name"]; o {
			profile.SetName(v.(string))
		}
		if v, o := rawProfile["spa_key_name"]; o {
			profile.SetSpaKeyName(v.(string))
		}
		if v, o := rawProfile["identity_provider_name"]; o {
			profile.SetIdentityProviderName(v.(string))
		}
		log.Printf("[DEBUG] Updating SINGLE PROFILE %+v", profile)
		result = append(result, profile)
	}
	return result
}

func resourceClientConnectionsDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete/Resetting Client Connections")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.ClientConnectionsApi

	if _, err := api.ClientConnectionsDelete(context.Background()).Authorization(token).Execute(); err != nil {
		return fmt.Errorf("Could reset Client Connections %+v", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
