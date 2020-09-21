package appgate

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/appgate/terraform-provider-appgate/client/v12/openapi"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAppgateMfaProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateMfaProviderCreate,
		Read:   resourceAppgateMfaProviderRead,
		Update: resourceAppgateMfaProviderUpdate,
		Delete: resourceAppgateMfaProviderDelete,
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

			"mfa_provider_id": {
				Type:        schema.TypeString,
				Description: "ID of the object.",
				Computed:    true,
			},

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

			"tags": {
				Type:        schema.TypeSet,
				Description: "Array of tags.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
					s := v.(string)
					list := []string{"Radius", "DefaultTimeBased"}
					for _, x := range list {
						if s == x {
							return
						}
					}
					errs = append(errs, fmt.Errorf("type must be on of %v, got %s", list, s))
					return
				},
			},

			"hostnames": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"shared_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"authentication_protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
					s := v.(string)
					list := []string{"OneFactor", "Challenge", "Push"}
					for _, x := range list {
						if s == x {
							return
						}
					}
					errs = append(errs, fmt.Errorf("mode must be on of %v, got %s", list, s))
					return
				},
			},

			"use_user_password": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"challenge_shared_secret": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAppgateMfaProviderCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating MFA provider: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.MFAProvidersApi
	args := openapi.NewMfaProviderWithDefaults()
	args.Id = uuid.New().String()
	args.SetName(d.Get("name").(string))
	args.SetNotes(d.Get("notes").(string))
	args.SetTags(schemaExtractTags(d))
	if v, ok := d.GetOk("type"); ok {
		args.SetType(v.(string))
	}
	if v, ok := d.GetOk("hostnames"); ok {
		hostnames, err := readArrayOfStringsFromConfig(v.(*schema.Set).List())
		if err != nil {
			return fmt.Errorf("Could not read hostnames %s", err)
		}
		args.SetHostnames(hostnames)
	}
	if v, ok := d.GetOk("port"); ok {
		args.SetPort(float32(v.(int)))
	}
	if v, ok := d.GetOk("shared_secret"); ok {
		args.SetSharedSecret(v.(string))
	}
	if v, ok := d.GetOk("authentication_protocol"); ok {
		args.SetAuthenticationProtocol(v.(string))
	}
	if v, ok := d.GetOk("timeout"); ok {
		args.SetTimeout(float32(v.(int)))
	}
	if v, ok := d.GetOk("mode"); ok {
		args.SetMode(v.(string))
	}
	if v, ok := d.GetOk("use_user_password"); ok {
		args.SetUseUserPassword(v.(bool))
	}
	if v, ok := d.GetOk("challenge_shared_secret"); ok {
		args.SetChallengeSharedSecret(v.(string))
	}

	request := api.MfaProvidersPost(context.TODO())
	request = request.MfaProvider(*args)

	mfaProvider, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create MFA provider %+v", prettyPrintAPIError(err))
	}

	d.SetId(mfaProvider.Id)
	d.Set("mfa_provider_id", mfaProvider.Id)

	return resourceAppgateMfaProviderRead(d, meta)
}

func resourceAppgateMfaProviderRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading MFA provider id: %+v", d.Id())
	token := meta.(*Client).Token
	api := meta.(*Client).API.MFAProvidersApi
	ctx := context.TODO()
	request := api.MfaProvidersIdGet(ctx, d.Id())
	mfaProvider, _, err := request.Authorization(token).Execute()
	if err != nil {
		// TODO check if 404
		d.SetId("")
		return fmt.Errorf("Failed to read MFA provider, %+v", err)
	}
	d.SetId(mfaProvider.Id)
	d.Set("mfa_provider_id", mfaProvider.Id)
	d.Set("name", mfaProvider.GetName())
	d.Set("notes", mfaProvider.GetNotes())
	d.Set("tags", mfaProvider.GetTags())
	d.Set("hostnames", mfaProvider.GetHostnames())
	d.Set("port", mfaProvider.GetPort())
	d.Set("authentication_protocol", mfaProvider.GetAuthenticationProtocol())
	d.Set("timeout", mfaProvider.GetTimeout())
	d.Set("mode", mfaProvider.GetMode())
	d.Set("use_user_password", mfaProvider.GetUseUserPassword())

	return nil
}

func resourceAppgateMfaProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating MFA provider: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.MFAProvidersApi
	ctx := context.TODO()
	request := api.MfaProvidersIdGet(ctx, d.Id())
	originalMfaProvider, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read MFA provider while updating, %+v", err)
	}

	if d.HasChange("name") {
		originalMfaProvider.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalMfaProvider.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalMfaProvider.SetTags(schemaExtractTags(d))
	}

	req := api.MfaProvidersIdPut(ctx, d.Id())
	req = req.MfaProvider(originalMfaProvider)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update MFA provider %+v", prettyPrintAPIError(err))
	}
	return resourceAppgateMfaProviderRead(d, meta)
}

func resourceAppgateMfaProviderDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete MFA provider: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.MFAProvidersApi

	request := api.MfaProvidersIdDelete(context.TODO(), d.Id())

	_, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not delete MFA provider %+v", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
