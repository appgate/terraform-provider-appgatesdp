package appgate

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v17/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

			"mfa_provider_id": resourceUUID(),
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

			"input_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Masked",
				ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
					s := v.(string)
					list := []string{"Masked", "Numeric", "Text"}
					for _, x := range list {
						if s == x {
							return
						}
					}
					errs = append(errs, fmt.Errorf("input_type must be on of %v, got %s", list, s))
					return
				},
			},

			"shared_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"authentication_protocol": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"timeout": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},

			"mode": {
				Type:     schema.TypeString,
				Computed: true,
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
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"challenge_shared_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func resourceAppgateMfaProviderCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating MFA provider: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.MFAProvidersApi
	args := openapi.NewMfaProviderWithDefaults()
	if v, ok := d.GetOk("mfa_provider_id"); ok {
		args.SetId(v.(string))
	}
	args.SetName(d.Get("name").(string))
	args.SetNotes(d.Get("notes").(string))
	args.SetTags(schemaExtractTags(d))
	if v, ok := d.GetOk("type"); ok {
		args.SetType(v.(string))
	}
	if v, ok := d.GetOk("hostnames"); ok {
		hostnames, err := readArrayOfStringsFromConfig(v.(*schema.Set).List())
		if err != nil {
			return fmt.Errorf("Could not read hostnames %w", err)
		}
		args.SetHostnames(hostnames)
	}
	if v, ok := d.GetOk("port"); ok {
		args.SetPort(float32(v.(int)))
	}
	if v, ok := d.GetOk("shared_secret"); ok {
		args.SetSharedSecret(v.(string))
	}
	if v, ok := d.GetOk("input_type"); ok {
		args.SetInputType(v.(string))
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
		return fmt.Errorf("Could not create MFA provider %w", prettyPrintAPIError(err))
	}

	d.SetId(mfaProvider.Id)
	d.Set("mfa_provider_id", mfaProvider.Id)

	return resourceAppgateMfaProviderRead(d, meta)
}

func resourceAppgateMfaProviderRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading MFA provider id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.MFAProvidersApi
	ctx := context.TODO()
	request := api.MfaProvidersIdGet(ctx, d.Id())
	mfaProvider, res, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Failed to read MFA provider, %w", err)
	}
	d.SetId(mfaProvider.Id)
	d.Set("mfa_provider_id", mfaProvider.Id)
	d.Set("name", mfaProvider.GetName())
	d.Set("notes", mfaProvider.GetNotes())
	d.Set("tags", mfaProvider.GetTags())
	d.Set("hostnames", mfaProvider.GetHostnames())
	d.Set("port", mfaProvider.GetPort())
	d.Set("input_type", mfaProvider.GetInputType())
	d.Set("authentication_protocol", mfaProvider.GetAuthenticationProtocol())
	d.Set("timeout", mfaProvider.GetTimeout())
	d.Set("mode", mfaProvider.GetMode())
	d.Set("use_user_password", mfaProvider.GetUseUserPassword())
	d.Set("authentication_protocol", mfaProvider.GetAuthenticationProtocol())

	return nil
}

func resourceAppgateMfaProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating MFA provider: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.MFAProvidersApi
	ctx := context.TODO()
	request := api.MfaProvidersIdGet(ctx, d.Id())
	originalMfaProvider, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read MFA provider while updating, %w", err)
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
	if d.HasChange("type") {
		originalMfaProvider.SetType(d.Get("notes").(string))
	}
	if d.HasChange("hostnames") {
		_, v := d.GetChange("hostnames")
		hostnames, err := readArrayOfStringsFromConfig(v.(*schema.Set).List())
		if err != nil {
			return fmt.Errorf("Failed to read hostnames %w", err)
		}
		originalMfaProvider.SetHostnames(hostnames)
	}
	if d.HasChange("port") {
		originalMfaProvider.SetPort(float32(d.Get("port").(int)))
	}
	if d.HasChange("shared_secret") {
		originalMfaProvider.SetSharedSecret(d.Get("shared_secret").(string))
	}
	if d.HasChange("input_type") {
		originalMfaProvider.SetInputType(d.Get("input_type").(string))
	}
	if d.HasChange("authentication_protocol") {
		originalMfaProvider.SetAuthenticationProtocol(d.Get("authentication_protocol").(string))
	}
	if d.HasChange("timeout") {
		originalMfaProvider.SetTimeout(float32(d.Get("timeout").(int)))
	}
	if d.HasChange("mode") {
		originalMfaProvider.SetMode(d.Get("mode").(string))
	}
	if d.HasChange("use_user_password") {
		originalMfaProvider.SetUseUserPassword(d.Get("use_user_password").(bool))
	}
	if d.HasChange("challenge_shared_secret") {
		originalMfaProvider.SetChallengeSharedSecret(d.Get("challenge_shared_secret").(string))
	}
	req := api.MfaProvidersIdPut(ctx, d.Id())
	req = req.MfaProvider(originalMfaProvider)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update MFA provider %w", prettyPrintAPIError(err))
	}
	return resourceAppgateMfaProviderRead(d, meta)
}

func resourceAppgateMfaProviderDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete MFA provider: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.MFAProvidersApi

	if _, err := api.MfaProvidersIdDelete(context.TODO(), d.Id()).Authorization(token).Execute(); err != nil {
		return fmt.Errorf("Could not delete MFA provider %w", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
