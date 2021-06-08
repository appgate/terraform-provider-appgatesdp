package appgate

import (
	"context"
	"fmt"
	"log"

	"github.com/appgate/sdp-api-client-go/api/v15/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateTrustedCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgateTrustedCertificateRead,
		Schema: map[string]*schema.Schema{
			"trusted_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"trusted_certificate_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceAppgateTrustedCertificateRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Data source trusted certificate")
	token := meta.(*Client).Token
	api := meta.(*Client).API.TrustedCertificatesApi

	trustedCertID, iok := d.GetOk("trusted_certificate_id")
	trustedCertName, nok := d.GetOk("trusted_certificate_name")

	if !iok && !nok {
		return fmt.Errorf("please provide one of trusted_certificate_id or trusted_certificate_name attributes")
	}
	var reqErr error
	var trustedCert *openapi.TrustedCertificate
	if iok {
		trustedCert, reqErr = findTrustedCertificateByUUID(api, trustedCertID.(string), token)
	} else {
		trustedCert, reqErr = findTrustedCertificateByName(api, trustedCertName.(string), token)
	}
	if reqErr != nil {
		return reqErr
	}
	log.Printf("[DEBUG] Got trusted certificate: %+v", trustedCert.Id)

	d.SetId(trustedCert.Id)
	d.Set("trusted_certificate_name", trustedCert.Name)
	d.Set("trusted_certificate_id", trustedCert.Id)
	return nil
}

func findTrustedCertificateByUUID(api *openapi.TrustedCertificatesApiService, id string, token string) (*openapi.TrustedCertificate, error) {
	log.Printf("[DEBUG] Data source trusted_certificate get by UUID %s", id)
	trustedCert, _, err := api.TrustedCertificatesIdGet(context.Background(), id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return &trustedCert, nil
}

func findTrustedCertificateByName(api *openapi.TrustedCertificatesApiService, name string, token string) (*openapi.TrustedCertificate, error) {
	log.Printf("[DEBUG] Data trusted_certificate get by name %s", name)
	request := api.TrustedCertificatesGet(context.Background())

	trustedCert, _, err := request.Query(name).OrderBy("name").Range_("0-1").Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	for _, s := range trustedCert.GetData() {
		return &s, nil
	}
	return nil, fmt.Errorf("Failed to find trusted_certificate %s", name)
}
