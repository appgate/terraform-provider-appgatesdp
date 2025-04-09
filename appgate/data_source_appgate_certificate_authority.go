package appgate

import (
	"context"
	"log"

	"github.com/appgate/sdp-api-client-go/api/v22/openapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateCertificateAuthority() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateCertificateAuthorityRead,
		Schema: map[string]*schema.Schema{
			"pem": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Get the current CA Certificate in PEM format.",
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"serial": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issuer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_from": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_to": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject_public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAppgateCertificateAuthorityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Data source CA Certificate.")
	var diags diag.Diagnostics
	_, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.CertificateAuthorityApi

	if pem, ok := d.GetOk("pem"); pem.(bool) && ok {
		cert, _, err := api.CertificateAuthorityCaPemGet(ctx).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		setCertificate(d, cert)
	} else {
		cert, _, err := api.CertificateAuthorityCaGet(ctx).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		setCertificate(d, cert)
	}

	return diags
}

func setCertificate(d *schema.ResourceData, cert *openapi.CaConfig) {
	d.Set("version", int(cert.GetVersion()))
	d.Set("serial", cert.GetSerial())
	d.Set("issuer", cert.GetIssuer())
	d.Set("subject", cert.GetSubject())
	d.Set("valid_from", cert.GetValidFrom().String())
	d.Set("valid_to", cert.GetValidTo().String())
	d.Set("fingerprint", cert.GetFingerprint())
	d.Set("certificate", cert.GetCertificate())
	d.Set("subject_public_key", cert.GetSubjectPublicKey())
	d.SetId(cert.GetFingerprint())
}
