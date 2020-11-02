package appgate

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v12/openapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	token := meta.(*Client).Token
	api := meta.(*Client).API.LicenseApi
	args := openapi.NewLicenseImportWithDefaults()
	args.SetLicense(d.Get("license").(string))
	request := api.LicensePost(context.TODO())
	license, _, err := request.LicenseImport(*args).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create license %+v", prettyPrintAPIError(err))
	}
	d.SetId(getLicenseIdentifier(license))
	return resourceAppgateLicenseRead(d, meta)
}

func getLicenseIdentifier(license openapi.License) string {
	return fmt.Sprintf(
		"%d-%s-%s-%d-%d",
		int(license.GetType()),
		license.GetRequest(),
		license.GetExpiration().Format("2006-01-02T15:04:05Z0700"),
		int(license.GetMaxSites()),
		int(license.GetMaxUsers()),
	)
}

func resourceAppgateLicenseRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading license")
	token := meta.(*Client).Token
	api := meta.(*Client).API.LicenseApi
	ctx := context.TODO()
	request := api.LicenseGet(ctx)
	license, _, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Failed to read license, %+v", err)
	}
	d.SetId(getLicenseIdentifier(license.GetEntitled()))

	return nil
}

func resourceAppgateLicenseDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete license")
	token := meta.(*Client).Token
	api := meta.(*Client).API.LicenseApi

	request := api.LicenseDelete(context.TODO())

	_, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not delete license %+v", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
