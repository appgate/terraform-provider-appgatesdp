package appgate

import (
	"context"
	"fmt"
	"log"

	"github.com/appgate/sdp-api-client-go/api/v21/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateLicense() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppgateLicenseCreate,
		ReadContext:   resourceAppgateLicenseRead,
		DeleteContext: resourceAppgateLicenseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: ImportLicenseState,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{

			"license": {
				Type:        schema.TypeString,
				Description: "The license file contents for this Controller (with the matching request code).",
				ForceNew:    true,
				Sensitive:   true,
				Optional:    true,
				Computed:    true,
			},
			"request_code": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
				Optional:  true,
			},
			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"max_users": {
				Type:     schema.TypeFloat,
				Computed: true,
				Optional: true,
			},
			"max_portal_users": {
				Type:     schema.TypeFloat,
				Computed: true,
				Optional: true,
			},
			"max_service_users": {
				Type:     schema.TypeFloat,
				Computed: true,
				Optional: true,
			},
			"users_lease_time_hours": {
				Type:     schema.TypeFloat,
				Computed: true,
				Optional: true,
			},
			"portal_users_lease_time_hours": {
				Type:     schema.TypeFloat,
				Computed: true,
				Optional: true,
			},
			"service_users_lease_time_hours": {
				Type:     schema.TypeFloat,
				Computed: true,
				Optional: true,
			},
			"max_sites": {
				Type:     schema.TypeFloat,
				Computed: true,
				Optional: true,
			},
			"max_connector_groups": {
				Type:     schema.TypeFloat,
				Computed: true,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeFloat,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func ImportLicenseState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func resourceAppgateLicenseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating license")
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	token, err := meta.(*Client).GetToken()
	diags = AppendFromErr(diags, err)

	api := meta.(*Client).API.LicenseApi
	args := openapi.LicenseImport{}
	if _, ok := d.GetOk("license"); !ok {
		return AppendErrorf(diags, "license attribute is required during creation")
	}
	args.SetLicense(d.Get("license").(string))

	license, _, err := api.LicensePost(ctx).LicenseImport(args).Authorization(token).Execute()
	if err != nil {
		return AppendFromErr(diags, fmt.Errorf("Could not create license %w", prettyPrintAPIError(err)))
	}
	d.SetId(license.GetId())

	return resourceAppgateLicenseRead(ctx, d, meta)
}

func resourceAppgateLicenseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Reading license")
	token, err := meta.(*Client).GetToken()
	diags = AppendFromErr(diags, err)
	api := meta.(*Client).API.LicenseApi

	licenses, _, err := api.LicenseGet(ctx).Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		return AppendFromErr(diags, fmt.Errorf("Failed to read license, %w", err))
	}
	for _, license := range licenses.GetLicenses() {
		if license.GetId() == d.Id() {
			d.SetId(license.GetId())
			if err := d.Set("request_code", licenses.GetRequestCode()); err != nil {
				return AppendFromErr(diags, err)
			}
			if err := d.Set("expiration", license.GetExpiration().String()); err != nil {
				return AppendFromErr(diags, err)
			}
			if err := d.Set("max_users", license.GetMaxUsers()); err != nil {
				return AppendFromErr(diags, err)
			}
			if err := d.Set("max_portal_users", license.GetMaxPortalUsers()); err != nil {
				return AppendFromErr(diags, err)
			}
			if err := d.Set("max_service_users", license.GetMaxServiceUsers()); err != nil {
				return AppendFromErr(diags, err)
			}
			if err := d.Set("users_lease_time_hours", license.GetUsersLeaseTimeHours()); err != nil {
				return AppendFromErr(diags, err)
			}
			if err := d.Set("portal_users_lease_time_hours", license.GetPortalUsersLeaseTimeHours()); err != nil {
				return AppendFromErr(diags, err)
			}
			if err := d.Set("service_users_lease_time_hours", license.GetServiceUsersLeaseTimeHours()); err != nil {
				return AppendFromErr(diags, err)
			}
			if err := d.Set("max_sites", license.GetMaxSites()); err != nil {
				return AppendFromErr(diags, err)
			}
			if err := d.Set("max_connector_groups", license.GetMaxConnectorGroups()); err != nil {
				return AppendFromErr(diags, err)
			}
			if err := d.Set("type", license.GetType()); err != nil {
				return AppendFromErr(diags, err)
			}
			return nil
		}
	}
	return AppendFromErr(diags, fmt.Errorf("license %s not found", d.Id()))
}

func resourceAppgateLicenseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete license")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.LicenseApi

	if _, err := api.LicenseDelete(ctx).Authorization(token).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("Could not delete license %w", prettyPrintAPIError(err)))
	}
	d.SetId("")
	return nil
}
