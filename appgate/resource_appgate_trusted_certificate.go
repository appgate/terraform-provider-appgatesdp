package appgate

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v15/openapi"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateTrustedCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateTrustedCertificateCreate,
		Read:   resourceAppgateTrustedCertificateRead,
		Update: resourceAppgateTrustedCertificateUpdate,
		Delete: resourceAppgateTrustedCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		SchemaVersion: 1,
		Schema: func() map[string]*schema.Schema {
			return mergeSchemaMaps(baseEntitySchema(), map[string]*schema.Schema{
				"trusted_certificate_id": {
					Type:        schema.TypeString,
					Description: "ID of the object.",
					Computed:    true,
				},
				"pem": {
					Type:        schema.TypeString,
					Description: "A certificate in PEM format.",
					Required:    true,
				},
			})
		}(),
	}
}

func resourceAppgateTrustedCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating trusted certificate: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.TrustedCertificatesApi
	args := openapi.NewTrustedCertificateWithDefaults()
	args.Id = uuid.New().String()
	args.SetName(d.Get("name").(string))
	args.SetNotes(d.Get("notes").(string))
	args.SetTags(schemaExtractTags(d))

	if v, ok := d.GetOk("pem"); ok {
		args.SetPem(v.(string))
	}

	request := api.TrustedCertificatesPost(context.TODO())
	request = request.TrustedCertificate(*args)

	trustedCertificate, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create trusted certificate %+v", prettyPrintAPIError(err))
	}

	d.SetId(trustedCertificate.Id)
	d.Set("trusted_certificate_id", trustedCertificate.Id)

	return resourceAppgateTrustedCertificateRead(d, meta)
}

func resourceAppgateTrustedCertificateRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading trusted certificate id: %+v", d.Id())
	token := meta.(*Client).Token
	api := meta.(*Client).API.TrustedCertificatesApi
	ctx := context.TODO()
	request := api.TrustedCertificatesIdGet(ctx, d.Id())
	trustedCertificate, _, err := request.Authorization(token).Execute()
	if err != nil {
		// TODO check if 404
		d.SetId("")
		return fmt.Errorf("Failed to read trusted certificate, %+v", err)
	}
	d.SetId(trustedCertificate.Id)
	d.Set("trusted_certificate_id", trustedCertificate.Id)
	d.Set("name", trustedCertificate.GetName())
	d.Set("notes", trustedCertificate.GetNotes())
	d.Set("tags", trustedCertificate.GetTags())
	d.Set("pem", trustedCertificate.GetPem())

	return nil
}

func resourceAppgateTrustedCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating trusted certificate: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.TrustedCertificatesApi
	ctx := context.TODO()
	request := api.TrustedCertificatesIdGet(ctx, d.Id())
	originalTrustedCertificate, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read trusted certificate while updating, %+v", err)
	}

	if d.HasChange("name") {
		originalTrustedCertificate.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalTrustedCertificate.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalTrustedCertificate.SetTags(schemaExtractTags(d))
	}

	if d.HasChange("pem") {
		originalTrustedCertificate.SetPem(d.Get("pem").(string))
	}

	req := api.TrustedCertificatesIdPut(ctx, d.Id())
	req = req.TrustedCertificate(originalTrustedCertificate)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update trusted certificate %+v", prettyPrintAPIError(err))
	}
	return resourceAppgateTrustedCertificateRead(d, meta)
}

func resourceAppgateTrustedCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete trusted certificate: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.TrustedCertificatesApi

	request := api.TrustedCertificatesIdDelete(context.TODO(), d.Id())

	_, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not delete trusted certificate %+v", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
