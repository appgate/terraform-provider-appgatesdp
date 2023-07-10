package appgate

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Global settings is as the name suggest global settings, and should be treated as a singleton object.
// we can never delete global_settings, but we can UPDATE and DELETE (reset) it.
// https://discuss.hashicorp.com/t/singleton-resource/9869
func resourceGlobalSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlobalSettingsCreate,
		ReadContext:   resourceGlobalSettingsRead,
		UpdateContext: resourceGlobalSettingsUpdate,
		DeleteContext: resourceGlobalSettingsDelete,
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
			"profile_hostname": {
				Type:        schema.TypeString,
				Description: "Client Connections, The hostname to use for generating profile URLs.",
				Optional:    true,
				Computed:    true,
			},
			"claims_token_expiration": {
				Type:        schema.TypeInt,
				Description: "Number of minutes the Claims Token is valid both for administrators and clients.",
				Optional:    true,
				Computed:    true,
			},
			"entitlement_token_expiration": {
				Type:        schema.TypeInt,
				Description: "Number of minutes the Entitlement Token is valid for clients.",
				Optional:    true,
				Computed:    true,
			},
			"administration_token_expiration": {
				Type:        schema.TypeInt,
				Description: "Number of minutes the administration Token is valid for administrators.",
				Optional:    true,
				Computed:    true,
			},
			"vpn_certificate_expiration": {
				Type:        schema.TypeInt,
				Description: "Number of minutes the VPN certificates is valid for clients.",
				Optional:    true,
				Computed:    true,
			},
			"login_banner_message": {
				Type:        schema.TypeString,
				Description: "The configured message will be displayed on the login UI.",
				Optional:    true,
				Computed:    true,
			},
			"message_of_the_day": {
				Type:        schema.TypeString,
				Description: "The configured message will be displayed after a successful logging.",
				Optional:    true,
				Computed:    true,
			},
			"backup_api_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the backup API is enabled or not.",
				Optional:    true,
				Computed:    true,
			},
			"has_backup_passphrase": {
				Type:        schema.TypeBool,
				Description: "Whether there is a backup passphrase set or not. Deprecated as of 5.0. Use backupApiEnabled instead.",
				Optional:    true,
				Computed:    true,
				Deprecated:  "Deprecated as of 5.0. Use backupApiEnabled instead.",
			},
			"backup_passphrase": {
				Type:        schema.TypeString,
				Description: "The passphrase to encrypt Appliance Backups when backup API is used.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
			},
			"fips": {
				Type:        schema.TypeBool,
				Description: "FIPS 140-2 Compliant Tunneling.",
				Optional:    true,
				Computed:    true,
			},
			"geo_ip_updates": {
				Type:        schema.TypeBool,
				Description: "Whether the automatic GeoIp updates are enabled or not.",
				Optional:    true,
				Computed:    true,
			},
			"audit_log_persistence_mode": {
				Type:        schema.TypeString,
				Description: "Audit Log persistence mode.",
				Optional:    true,
				Computed:    true,
			},
			"app_discovery_domains": {
				Type:        schema.TypeSet,
				Description: "Domains to monitor for for App Discovery feature.",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"registered_device_expiration_days": {
				Type:        schema.TypeInt,
				Description: "Number of days registered devices are kept in storage before being deleted",
				Optional:    true,
			},
			"spa_mode": {
				Type:        schema.TypeString,
				Description: "SPA Mode.",
				Optional:    true,
				Computed:    true,
				ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
					// controller validation error message is wrong:
					// spaMode FOO was not one of [TCP, UDP_TCP]
					// should be
					// spaMode FOO was not one of [TCP, UDP-TCP]
					// so will validate here first.
					s := v.(string)
					list := []string{
						"TCP",
						"UDP-TCP",
					}
					for _, x := range list {
						if s == x {
							return
						}
					}
					errs = append(errs, fmt.Errorf("type must be on of %v, got %s", list, s))
					return
				},
			},
			"spa_time_window_seconds": {
				Type:        schema.TypeInt,
				Description: "Number of seconds the time skew SPA will allow.",
				Optional:    true,
				Computed:    true,
			},
			"collective_name": {
				Type:        schema.TypeString,
				Description: "Friendly name for the Collective.",
				Optional:    true,
				Computed:    true,
			},
			"collective_id": {
				Type:        schema.TypeString,
				Description: "A randomly generated ID during first installation to identify the Collective.",
				Computed:    true,
			},
		},
	}
}

func resourceGlobalSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceGlobalSettingsUpdate(ctx, d, meta)
}

func resourceGlobalSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Reading Global settings id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.GlobalSettingsApi
	currentVersion := meta.(*Client).ApplianceVersion
	request := api.GlobalSettingsGet(ctx)
	settings, res, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.FromErr(fmt.Errorf("Failed to read Global settings, %w", err))
	}
	d.SetId(settings.GetCollectiveId())
	d.Set("claims_token_expiration", settings.GetClaimsTokenExpiration())
	d.Set("entitlement_token_expiration", settings.GetEntitlementTokenExpiration())
	d.Set("administration_token_expiration", settings.GetAdministrationTokenExpiration())
	d.Set("vpn_certificate_expiration", settings.GetVpnCertificateExpiration())
	d.Set("login_banner_message", settings.GetLoginBannerMessage())
	d.Set("message_of_the_day", settings.GetMessageOfTheDay())
	d.Set("backup_api_enabled", settings.GetBackupApiEnabled())
	if val, ok := d.GetOk("backup_passphrase"); ok {
		d.Set("backup_passphrase", val)
	} else {
		d.Set("backup_passphrase", settings.GetBackupPassphrase())
	}
	d.Set("fips", settings.GetFips())
	d.Set("geo_ip_updates", settings.GetGeoIpUpdates())
	d.Set("audit_log_persistence_mode", settings.GetAuditLogPersistenceMode())
	d.Set("app_discovery_domains", settings.GetAppDiscoveryDomains())
	d.Set("collective_id", settings.GetCollectiveId())

	d.Set("collective_name", settings.GetCollectiveName())
	d.Set("spa_time_window_seconds", settings.GetSpaTimeWindowSeconds())
	d.Set("spa_mode", settings.GetSpaMode())

	if currentVersion.GreaterThanOrEqual(Appliance54Version) {
		ccAPI := meta.(*Client).API.ClientProfilesApi
		request := ccAPI.ClientConnectionsGet(ctx)
		clientConnections, _, err := request.Authorization(token).Execute()
		if err != nil {
			d.SetId("")
			return diag.FromErr(fmt.Errorf("Failed to read Client Connections, %w", err))
		}
		d.Set("profile_hostname", clientConnections.GetProfileHostname())
	}
	if currentVersion.GreaterThanOrEqual(Appliance62Version) {
		d.Set("registered_device_expiration_days", settings.GetRegisteredDeviceExpirationDays())
	}
	return diags
}

func resourceGlobalSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating Global settings")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.GlobalSettingsApi
	currentVersion := meta.(*Client).ApplianceVersion

	request := api.GlobalSettingsGet(ctx)
	originalsettings, _, err := request.Authorization(token).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to read Global settings while updating, %w", err))
	}

	if d.HasChange("claims_token_expiration") {
		originalsettings.SetClaimsTokenExpiration(float32(d.Get("claims_token_expiration").(int)))
	}
	if d.HasChange("entitlement_token_expiration") {
		originalsettings.SetEntitlementTokenExpiration(float32(d.Get("entitlement_token_expiration").(int)))
	}
	if d.HasChange("administration_token_expiration") {
		originalsettings.SetAdministrationTokenExpiration(float32(d.Get("administration_token_expiration").(int)))
	}
	if d.HasChange("vpn_certificate_expiration") {
		originalsettings.SetVpnCertificateExpiration(float32(d.Get("vpn_certificate_expiration").(int)))
	}
	if d.HasChange("login_banner_message") {
		originalsettings.SetLoginBannerMessage(d.Get("login_banner_message").(string))
	}
	if d.HasChange("message_of_the_day") {
		originalsettings.SetMessageOfTheDay(d.Get("message_of_the_day").(string))
	}
	if d.HasChange("backup_api_enabled") {
		originalsettings.SetBackupApiEnabled(d.Get("backup_api_enabled").(bool))
	}
	if d.HasChange("backup_passphrase") {
		originalsettings.SetBackupPassphrase(d.Get("backup_passphrase").(string))
	}
	if d.HasChange("fips") {
		originalsettings.SetFips(d.Get("fips").(bool))
	}
	if d.HasChange("geo_ip_updates") {
		originalsettings.SetGeoIpUpdates(d.Get("geo_ip_updates").(bool))
	}
	if d.HasChange("audit_log_persistence_mode") {
		originalsettings.SetAuditLogPersistenceMode(d.Get("audit_log_persistence_mode").(string))
	}
	if d.HasChange("app_discovery_domains") {
		_, n := d.GetChange("app_discovery_domains")
		domains, err := readArrayOfStringsFromConfig(n.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		originalsettings.SetAppDiscoveryDomains(domains)
	}
	if d.HasChange("registered_device_expiration_days") {
		if currentVersion.LessThan(Appliance62Version) {
			return diag.Errorf("registered_device_expiration_days is not supported on %s", currentVersion.String())
		} else if currentVersion.GreaterThanOrEqual(Appliance62Version) {
			originalsettings.SetRegisteredDeviceExpirationDays(float32(d.Get("registered_device_expiration_days").(int)))
		}
	}
	if d.HasChange("spa_mode") {
		if currentVersion.LessThan(Appliance55Version) {
			return diag.Errorf("spa_mode is not supported on %s", currentVersion.String())
		} else if currentVersion.GreaterThanOrEqual(Appliance55Version) {
			originalsettings.SetSpaMode(d.Get("spa_mode").(string))
		}
	}
	if d.HasChange("spa_time_window_seconds") {
		if currentVersion.LessThan(Appliance60Version) {
			return diag.Errorf("spa_time_window_seconds is not supported on %s", currentVersion.String())
		} else if currentVersion.GreaterThanOrEqual(Appliance60Version) {
			originalsettings.SetSpaTimeWindowSeconds(float32(d.Get("spa_time_window_seconds").(int)))
		}
	}
	if d.HasChange("collective_name") {
		if currentVersion.LessThan(Appliance60Version) {
			return diag.Errorf("collective_name is not supported on %s", currentVersion.String())
		} else if currentVersion.GreaterThanOrEqual(Appliance60Version) {
			originalsettings.SetCollectiveName(d.Get("collective_name").(string))
		}
	}
	log.Printf("[DEBUG] Updating Global settings %+v", originalsettings)
	req := api.GlobalSettingsPut(ctx)
	_, err = req.GlobalSettings(*originalsettings).Authorization(token).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Could not update Global settings %w", prettyPrintAPIError(err)))
	}

	if currentVersion.GreaterThanOrEqual(Appliance54Version) && d.HasChange("profile_hostname") {
		ccAPI := meta.(*Client).API.ClientProfilesApi
		request := ccAPI.ClientConnectionsGet(ctx)
		originalclientConnections, _, err := request.Authorization(token).Execute()
		if err != nil {
			d.SetId("")
			return diag.FromErr(fmt.Errorf("Failed to read Client Connections, %w", err))
		}
		_, v := d.GetChange("profile_hostname")
		originalclientConnections.SetProfileHostname(v.(string))
		req := ccAPI.ClientConnectionsPut(ctx)
		_, _, err = req.ClientConnections(*originalclientConnections).Authorization(token).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf("Could not update Client Connections %w", prettyPrintAPIError(err)))
		}
	}
	return resourceGlobalSettingsRead(ctx, d, meta)
}

func resourceGlobalSettingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Delete Global settings")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.GlobalSettingsApi

	if _, err := api.GlobalSettingsDelete(context.Background()).Authorization(token).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("Could not reset Global settings %w", prettyPrintAPIError(err)))
	}
	d.SetId("")
	// The API won't allow us to delete/remove/reset profile_hostname from client_connections so we will just leave it be.
	return diags
}
