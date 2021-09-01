package appgate

import (
	"context"
	"fmt"
	"log"

	"github.com/appgate/sdp-api-client-go/api/v15/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateBlacklistUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateBlacklistUserCreate,
		Read:   resourceAppgateBlacklistUserRead,
		Delete: resourceAppgateBlacklistUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{

			"user_distinguished_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"username": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"reason": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
		},
	}
}

func resourceAppgateBlacklistUserCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating blacklisted user")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.BlacklistedUsersApi
	args := openapi.NewBlacklistEntryWithDefaults()

	if v, ok := d.GetOk("user_distinguished_name"); ok {
		args.SetUserDistinguishedName(v.(string))
	}
	if v, ok := d.GetOk("username"); ok {
		args.SetUsername(v.(string))
	}
	if v, ok := d.GetOk("provider_name"); ok {
		args.SetProviderName(v.(string))
	}
	if v, ok := d.GetOk("reason"); ok {
		args.SetReason(v.(string))
	}
	request := api.BlacklistPost(context.TODO())
	request = request.BlacklistEntry(*args)

	entry, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create blacklisted user %+v", prettyPrintAPIError(err))
	}

	d.SetId(entry.GetUserDistinguishedName())

	return resourceAppgateBlacklistUserRead(d, meta)
}

func queryEntry(ctx context.Context, api *openapi.BlacklistedUsersApiService, token, distinguishedName string) (*openapi.BlacklistEntry, error) {
	request := api.BlacklistGet(ctx)

	list, _, err := request.Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	for _, s := range list.GetData() {
		if s.GetUserDistinguishedName() == distinguishedName {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("Failed to find blacklist user %s", distinguishedName)
}

func resourceAppgateBlacklistUserRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading blacklisted user id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.BlacklistedUsersApi
	entry, err := queryEntry(context.TODO(), api, token, d.Id())
	if err != nil {
		return err
	}

	d.Set("user_distinguished_name", entry.GetUserDistinguishedName())
	d.Set("username", entry.GetUsername())
	d.Set("provider_name", entry.GetProviderName())
	d.Set("reason", entry.GetReason())
	return nil
}

func resourceAppgateBlacklistUserDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading blacklisted user id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.BlacklistedUsersApi

	if _, err := api.BlacklistDistinguishedNameDelete(context.TODO(), d.Id()).Authorization(token).Execute(); err != nil {
		return fmt.Errorf("Could not delete blacklisted user %+v", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
