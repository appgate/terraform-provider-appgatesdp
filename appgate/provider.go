package appgate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	Version12 = 12
	Version13 = 13
	Version14 = 14
	Version15 = 15
	Version16 = 16
	// DefaultClientVersion is the latest support version of appgate sdp client that is supported.
	// its not recommended to change this value.
	DefaultClientVersion = Version15
)

var (
	// ApplianceVersionMap match appliance version to go client version.
	ApplianceVersionMap = map[int]string{
		Version12: "5.1.0",
		Version13: "5.2.0",
		Version14: "5.3.0",
		Version15: "5.4.0",
		Version16: "5.5.0",
	}

	Appliance53Version, _ = version.NewVersion(ApplianceVersionMap[Version14])
	Appliance54Version, _ = version.NewVersion(ApplianceVersionMap[Version15])
)

// Provider function returns the object that implements the terraform.ResourceProvider interface, specifically a schema.Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPGATE_ADDRESS", ""),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPGATE_USERNAME", ""),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPGATE_PASSWORD", ""),
			},
			"provider": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPGATE_PROVIDER", "local"),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPGATE_INSECURE", true),
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPGATE_HTTP_DEBUG", false),
			},
			"client_version": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPGATE_CLIENT_VERSION", DefaultClientVersion),
			},
			"config_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPGATE_CONFIG_PATH", nil),
				Description: "Path to the appgate config file. Can be set with APPGATE_CONFIG_PATH.",
				// ConflictsWith: []string{"url"},
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"appgatesdp_appliance":               dataSourceAppgateAppliance(),
			"appgatesdp_entitlement":             dataSourceAppgateEntitlement(),
			"appgatesdp_site":                    dataSourceAppgateSite(),
			"appgatesdp_condition":               dataSourceAppgateCondition(),
			"appgatesdp_policy":                  dataSourceAppgatePolicy(),
			"appgatesdp_ringfence_rule":          dataSourceAppgateRingfenceRule(),
			"appgatesdp_criteria_script":         dataSourceCriteriaScript(),
			"appgatesdp_entitlement_script":      dataSourceEntitlementScript(),
			"appgatesdp_device_script":           dataSourceDeviceScript(),
			"appgatesdp_appliance_customization": dataSourceAppgateApplianceCustomization(),
			"appgatesdp_ip_pool":                 dataSourceAppgateIPPool(),
			"appgatesdp_administrative_role":     dataSourceAppgateAdministrativeRole(),
			"appgatesdp_global_settings":         dataSourceGlobalSettings(),
			"appgatesdp_trusted_certificate":     dataSourceAppgateTrustedCertificate(),
			"appgatesdp_mfa_provider":            dataSourceAppgateMfaProvider(),
			"appgatesdp_local_user":              dataSourceAppgateLocalUser(),
			"appgatesdp_identity_provider":       dataSourceAppgateIdentityProvider(),
			"appgatesdp_appliance_seed":          dataSourceAppgateApplianceSeed(),
			"appgatesdp_certificate_authority":   dataSourceAppgateCertificateAuthority(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"appgatesdp_appliance":                          resourceAppgateAppliance(),
			"appgatesdp_entitlement":                        resourceAppgateEntitlement(),
			"appgatesdp_site":                               resourceAppgateSite(),
			"appgatesdp_ringfence_rule":                     resourceAppgateRingfenceRule(),
			"appgatesdp_condition":                          resourceAppgateCondition(),
			"appgatesdp_policy":                             resourceAppgatePolicy(),
			"appgatesdp_criteria_script":                    resourceAppgateCriteriaScript(),
			"appgatesdp_entitlement_script":                 resourceAppgateEntitlementScript(),
			"appgatesdp_device_script":                      resourceAppgateDeviceScript(),
			"appgatesdp_appliance_customization":            resourceAppgateApplianceCustomizations(),
			"appgatesdp_ip_pool":                            resourceAppgateIPPool(),
			"appgatesdp_administrative_role":                resourceAppgateAdministrativeRole(),
			"appgatesdp_global_settings":                    resourceGlobalSettings(),
			"appgatesdp_ldap_identity_provider":             resourceAppgateLdapProvider(),
			"appgatesdp_trusted_certificate":                resourceAppgateTrustedCertificate(),
			"appgatesdp_mfa_provider":                       resourceAppgateMfaProvider(),
			"appgatesdp_local_user":                         resourceAppgateLocalUser(),
			"appgatesdp_license":                            resourceAppgateLicense(),
			"appgatesdp_admin_mfa_settings":                 resourceAdminMfaSettings(),
			"appgatesdp_client_connections":                 resourceClientConnections(),
			"appgatesdp_blacklist_user":                     resourceAppgateBlacklistUser(),
			"appgatesdp_radius_identity_provider":           resourceAppgateRadiusProvider(),
			"appgatesdp_saml_identity_provider":             resourceAppgateSamlProvider(),
			"appgatesdp_local_database_identity_provider":   resourceAppgateLocalDatabaseProvider(),
			"appgatesdp_ldap_certificate_identity_provider": resourceAppgateLdapCertificateProvider(),
			"appgatesdp_connector_identity_provider":        resourceAppgateConnectorProvider(),
			"appgatesdp_client_profile":                     resourceAppgateClientProfile(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func requiredParameters(d *schema.ResourceData) bool {
	required := []string{
		"username",
		"password",
		"url",
	}
	for _, r := range required {
		if _, ok := d.GetOk(r); !ok {
			return false
		}
	}
	return true
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	config := Config{}
	config.Timeout = 20

	if path, ok := d.GetOkExists("config_path"); ok {
		p := path.(string)
		file, err := os.OpenFile(p, os.O_RDWR, 0644)
		if errors.Is(err, os.ErrNotExist) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Missing Appgate SDP credentials",
				Detail:   fmt.Sprintf("appgate config_path file not found %s", err),
			})
			return nil, diags
		} else if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Missing Appgate SDP credentials",
				Detail:   err.Error(),
			})
			return nil, diags
		}
		defer file.Close()
		configFile := AppgateConfigFile{}
		if err := json.NewDecoder(file).Decode(&configFile); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Missing Appgate SDP credentials",
				Detail:   fmt.Sprintf("appgate config_path invalid json format %s", err),
			})
			return nil, diags
		}
		if len(configFile.URL) > 0 {
			config.URL = configFile.URL
		}
		if len(config.Username) > 0 {
			config.Username = configFile.Username
		}
		if len(config.Password) > 0 {
			config.Password = configFile.Password
		}
		if len(config.Provider) > 0 {
			config.Provider = configFile.Provider
		}
		if config.Version > 0 {
			config.Version = configFile.ClientVersion
		}
		config.Insecure = configFile.Insecure

	} else if !requiredParameters(d) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Missing Appgate SDP credentials",
			Detail:   "Appgate client is unauthenticated. Provide user credentials and URL to access restricted resources.",
		})
		return nil, diags
	}
	// overwrite config file value with ENV Variables or set attributes in the provider block
	if v, ok := d.GetOk("username"); ok {
		config.Username = v.(string)
	}
	if v, ok := d.GetOk("password"); ok {
		config.Password = v.(string)
	}
	if v, ok := d.GetOk("url"); ok {
		config.URL = v.(string)
	}
	if v, ok := d.GetOk("provider"); ok {
		config.Provider = v.(string)
	}
	if v, ok := d.GetOk("insecure"); ok {
		config.Insecure = v.(bool)
	}
	if v, ok := d.GetOk("provider"); ok {
		config.Provider = v.(string)
	}
	if v, ok := d.GetOk("debug"); ok {
		config.Debug = v.(bool)
	}
	if v, ok := d.GetOk("client_version"); ok {
		config.Version = v.(int)
	}

	c, err := config.Client()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to create Appgate SDP SDK client v%d", config.Version),
			Detail:   fmt.Sprintf("Unable to authenticate user for authenticated Appgate SDP client %s", err),
		})

		return nil, diags
	}
	if c == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to create Appgate SDP SDK client v%d", config.Version),
			Detail:   "Appgate sdp client is nil, internal error",
		})
		return nil, diags
	}

	return c, diags
}
