package appgate

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateSite() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateSiteRead,
		Schema: map[string]*schema.Schema{
			"site_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"site_name"},
			},
			"site_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"site_id"},
			},

			"notes": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"short_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAppgateSiteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Data source Site")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.SitesApi
	site, diags := ResolveSiteFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}

	d.SetId(site.GetId())
	d.Set("site_name", site.GetName())
	d.Set("short_name", site.GetShortName())
	d.Set("site_id", site.GetId())
	d.Set("notes", site.GetNotes())
	d.Set("tags", site.GetTags())

	return nil
}
