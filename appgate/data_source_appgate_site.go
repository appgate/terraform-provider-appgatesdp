package appgate

import (
	"context"
	"fmt"
	"log"

	"github.com/appgate/sdp-api-client-go/api/v13/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateSite() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgateSiteRead,
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

func dataSourceAppgateSiteRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Data source Site")
	token := meta.(*Client).Token
	api := meta.(*Client).API.SitesApi

	siteID, iok := d.GetOk("site_id")
	siteName, nok := d.GetOk("site_name")

	if !iok && !nok {
		return fmt.Errorf("please provide one of site_id or site_name attributes")
	}
	var reqErr error
	var site *openapi.Site
	if iok {
		site, reqErr = findSiteByUUID(api, siteID.(string), token)
	} else {
		site, reqErr = findSiteByName(api, siteName.(string), token)
	}
	if reqErr != nil {
		return reqErr
	}
	log.Printf("[DEBUG] Got Site: %+v", site)

	d.SetId(site.Id)
	d.Set("name", site.Name)
	d.Set("short_name", site.ShortName)
	d.Set("site_id", site.Id)
	d.Set("notes", site.Notes)
	d.Set("tags", site.Tags)

	return nil
}

func findSiteByUUID(api *openapi.SitesApiService, id string, token string) (*openapi.Site, error) {
	log.Printf("[DEBUG] Data source Site get by UUID %s", id)
	site, _, err := api.SitesIdGet(context.Background(), id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return &site, nil
}

func findSiteByName(api *openapi.SitesApiService, name string, token string) (*openapi.Site, error) {
	log.Printf("[DEBUG] Data source Site get by name %s", name)
	request := api.SitesGet(context.Background())

	site, _, err := request.Query(name).OrderBy("name").Range_("0-1").Authorization(token).Execute()
	if err != nil {
		return nil, err
	}

	for _, s := range site.GetData() {
		return &s, nil
	}
	return nil, fmt.Errorf("Failed to find site %s", name)
}
