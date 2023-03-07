package appgate

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateTrustedCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateTrustedCertificateRead,
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

func dataSourceAppgateTrustedCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Data source trusted certificate")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.TrustedCertificatesApi
	trustedCert, diags := ResolveTrustedCertificateFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}

	d.SetId(trustedCert.GetId())
	d.Set("trusted_certificate_name", trustedCert.GetName())
	d.Set("trusted_certificate_id", trustedCert.GetId())
	return nil
}
