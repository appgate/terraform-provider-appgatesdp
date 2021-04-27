package appgate

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DefaultClientVersion is the latest support version of appgate sdp client that is supported.
// its not recommended to change this value.
const DefaultClientVersion = 14

// ApplianceVersionMap match appliance version to go client version.
var ApplianceVersionMap = map[int]string{
	12: "5.1.0",
	13: "5.2.0",
	14: "5.3.0",
	15: "5.4.0",
}

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
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if (username != "") && (password != "") {
		v := d.Get("client_version").(int)
		config := Config{
			URL:      d.Get("url").(string),
			Username: d.Get("username").(string),
			Password: d.Get("password").(string),
			Provider: d.Get("provider").(string),
			Insecure: d.Get("insecure").(bool),
			Timeout:  20,
			Debug:    d.Get("debug").(bool),
			Version:  v,
		}
		c, err := config.Client()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Unable to create Appgate SDP SDK client v%d", v),
				Detail:   fmt.Sprintf("Unable to authenticate user for authenticated Appgate SDP client %s", err),
			})

			return nil, diags
		}
		if c.ApplianceVersion.LessThan(c.LatestSupportedVersion) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("You are using old version of Appgate SDP SDK client v%d", v),
				Detail:   "You are using an outdated version of appgate appliances, you should consider updating to the latest version.",
			})
		}
		log.Printf("[INFO] Appgate SDP Appliance Version %q client version v%d", c.ApplianceVersion, c.ClientVersion)
		return c, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Using unauthenticated Appgate client",
		Detail:   "Appgate client is unauthenticated. Provide user credentials to access restricted resources.",
	})

	return nil, diags
}
