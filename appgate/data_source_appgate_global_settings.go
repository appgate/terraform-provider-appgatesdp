package appgate

import (
	"context"
	"fmt"

	"github.com/appgate/sdp-api-client-go/api/v13/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGlobalSettings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgateGlobalSettingsRead,
		Schema: map[string]*schema.Schema{
			"claims_token_expiration": {
				Type:        schema.TypeInt,
				Description: "Number of minutes the Claims Token is valid both for administrators and clients.",
				Computed:    true,
			},
			"entitlement_token_expiration": {
				Type:        schema.TypeInt,
				Description: "Number of minutes the Entitlement Token is valid for clients.",
				Computed:    true,
			},
			"administration_token_expiration": {
				Type:        schema.TypeInt,
				Description: "Number of minutes the administration Token is valid for administrators.",
				Computed:    true,
			},
			"vpn_certificate_expiration": {
				Type:        schema.TypeInt,
				Description: "Number of minutes the VPN certificates is valid for clients.",
				Computed:    true,
			},
			"login_banner_message": {
				Type:        schema.TypeString,
				Description: "The configured message will be displayed on the login UI.",
				Computed:    true,
			},
			"message_of_the_day": {
				Type:        schema.TypeString,
				Description: "The configured message will be displayed after a successful loging.",
				Computed:    true,
			},
			"backup_api_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the backup API is enabled or not.",
				Computed:    true,
			},
			"has_backup_passphrase": {
				Type:        schema.TypeBool,
				Description: "Whether there is a backup passphrase set or not. Deprecated as of 5.0. Use backupApiEnabled instead.",
				Computed:    true,
			},
			"fips": {
				Type:        schema.TypeBool,
				Description: "FIPS 140-2 Compliant Tunneling.",
				Computed:    true,
			},
			"geo_ip_updates": {
				Type:        schema.TypeBool,
				Description: "Whether the automatic GeoIp updates are enabled or not.",
				Computed:    true,
			},
			"audit_log_persistence_mode": {
				Type:        schema.TypeString,
				Description: "Audit Log persistence mode.",
				Computed:    true,
			},
			"app_discovery_domains": {
				Type:        schema.TypeSet,
				Description: "Domains to monitor for for App Discovery feature.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"collective_id": {
				Type:        schema.TypeString,
				Description: "A randomly generated ID during first installation to identify the Collective.",
				Computed:    true,
			},
		},
	}
}

func dataSourceAppgateGlobalSettingsRead(d *schema.ResourceData, meta interface{}) error {
	token := meta.(*Client).Token
	api := meta.(*Client).API.GlobalSettingsApi

	settings, err := getGlobalSettings(api, token)
	if err != nil {
		return fmt.Errorf("Could not read global settings %s", err)
	}
	d.SetId(settings.GetCollectiveId())
	d.Set("claims_token_expiration", settings.GetClaimsTokenExpiration())
	d.Set("entitlement_token_expiration", settings.GetEntitlementTokenExpiration())
	d.Set("administration_token_expiration", settings.GetAdministrationTokenExpiration())
	d.Set("vpn_certificate_expiration", settings.GetVpnCertificateExpiration())
	d.Set("login_banner_message", settings.GetLoginBannerMessage())
	d.Set("message_of_the_day", settings.GetMessageOfTheDay())
	d.Set("backup_api_enabled", settings.GetBackupApiEnabled())
	d.Set("has_backup_passphrase", settings.GetHasBackupPassphrase())
	d.Set("fips", settings.GetFips())
	d.Set("geo_ip_updates", settings.GetGeoIpUpdates())
	d.Set("audit_log_persistence_mode", settings.GetAuditLogPersistenceMode())
	d.Set("app_discovery_domains", settings.GetAppDiscoveryDomains())
	d.Set("collective_id", settings.GetCollectiveId())
	return nil
}

func getGlobalSettings(api *openapi.GlobalSettingsApiService, token string) (*openapi.GlobalSettings, error) {
	globalSettings, _, err := api.GlobalSettingsGet(context.Background()).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return &globalSettings, nil
}
