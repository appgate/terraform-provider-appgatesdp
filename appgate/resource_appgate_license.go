package appgate

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v17/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateLicense() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateLicenseCreate,
		Read:   resourceAppgateLicenseRead,
		Delete: resourceAppgateLicenseDelete,
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

			"license": {
				Type:        schema.TypeString,
				Description: "The license file contents for this Controller (with the matching request code).",
				ForceNew:    true,
				Sensitive:   true,
				Required:    true,
			},
		},
	}
}

func resourceAppgateLicenseCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating license")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.LicenseApi
	args := openapi.NewLicenseImportWithDefaults()
	args.SetLicense(d.Get("license").(string))

	if _, _, err := api.LicensePost(context.TODO()).LicenseImport(*args).Authorization(token).Execute(); err != nil {
		return fmt.Errorf("Could not create license %w", prettyPrintAPIError(err))
	}
	return resourceAppgateLicenseRead(d, meta)
}

func getLicenseIdentifier(license *openapi.LicenseDetails) string {
	return license.GetRequestCode()
}

func resourceAppgateLicenseRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading license")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.LicenseApi
	ctx := context.TODO()
	request := api.LicenseGet(ctx)
	license, _, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Failed to read license, %w", err)
	}
	d.SetId(getLicenseIdentifier(license))

	return nil
}

func resourceAppgateLicenseDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete license")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.LicenseApi

	if _, err := api.LicenseDelete(context.TODO()).Authorization(token).Execute(); err != nil {
		return fmt.Errorf("Could not delete license %w", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
