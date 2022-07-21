package appgate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/imdario/mergo"
)

const (
	Version12 = 12
	Version13 = 13
	Version14 = 14
	Version15 = 15
	Version16 = 16
	Version17 = 17
	// DefaultClientVersion is the latest support version of appgate sdp client that is supported.
	// its not recommended to change this value.
	DefaultClientVersion = Version16
)

var (
	// ApplianceVersionMap match appliance version to go client version.
	ApplianceVersionMap = map[int]string{
		Version12: "5.1.0",
		Version13: "5.2.0",
		Version14: "5.3.0",
		Version15: "5.4.0",
		Version16: "5.5.0",
		Version17: "6.0.0",
	}

	Appliance53Version, _ = version.NewVersion(ApplianceVersionMap[Version14])
	Appliance54Version, _ = version.NewVersion(ApplianceVersionMap[Version15])
	Appliance55Version, _ = version.NewVersion(ApplianceVersionMap[Version16])
)

// Provider function returns the object that implements the terraform.ResourceProvider interface, specifically a schema.Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPGATE_ADDRESS", nil),
			},
			"username": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("APPGATE_USERNAME", nil),
				ConflictsWith: []string{"bearer_token"},
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("APPGATE_PASSWORD", nil),
				ConflictsWith: []string{"bearer_token"},
			},
			"provider": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPGATE_PROVIDER", "local"),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPGATE_INSECURE", false),
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
			},
			"pem_filepath": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPGATE_PEM_FILEPATH", nil),
				Description: "Path to the controller's CA cert file in PEM format",
			},
			"bearer_token": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("APPGATE_BEARER_TOKEN", nil),
				ConflictsWith: []string{"username", "password"},
				Description:   "The Token from the LoginResponse, provided from outside terraform.",
			},
			"device_id": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("APPGATE_DEVICE_ID", nil),
				ValidateFunc: validation.IsUUID,
				Description:  "UUID to distinguish the Client device making the request. It is supposed to be same for every login request from the same server.",
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
			"appgatesdp_user_claim_script":       dataSourceUserClaimScript(),
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
			"appgatesdp_appliance_controller_activation":    resourceAppgateApplianceControllerActivation(),
			"appgatesdp_entitlement":                        resourceAppgateEntitlement(),
			"appgatesdp_site":                               resourceAppgateSite(),
			"appgatesdp_ringfence_rule":                     resourceAppgateRingfenceRule(),
			"appgatesdp_condition":                          resourceAppgateCondition(),
			"appgatesdp_policy":                             resourceAppgatePolicy(),
			"appgatesdp_criteria_script":                    resourceAppgateCriteriaScript(),
			"appgatesdp_entitlement_script":                 resourceAppgateEntitlementScript(),
			"appgatesdp_device_script":                      resourceAppgateDeviceScript(),
			"appgatesdp_user_claim_script":                  resourceAppgateUserClaimScript(),
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

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	config := Config{}
	config.Timeout = 20
	configFile := Config{}
	usingFile := false

	if path, ok := d.GetOk("config_path"); ok {
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
		if err := json.NewDecoder(file).Decode(&configFile); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Missing Appgate SDP credentials",
				Detail:   fmt.Sprintf("appgate config_path invalid json format %s", err),
			})
			return nil, diags
		}
		usingFile = true
	}

	if v, ok := d.GetOk("bearer_token"); ok {
		config.BearerToken = v.(string)
	}
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
	if v, ok := d.GetOk("pem_filepath"); ok {
		config.PemFilePath = v.(string)
	}
	if v, ok := d.GetOk("device_id"); ok {
		config.DeviceID = v.(string)
	}

	if usingFile {
		// we do not allow bool config keys from the config file, since they will always default to false if omitted
		// for the boolean config attributes, we will fallback to the default values defined in the Schema.
		// (yes this can be solved by pointers, however we are not interesting in doing that now)
		// https://play.golang.org/p/QNkWPEjPlcD
		configFile.Insecure = config.Insecure
		configFile.Debug = config.Debug
		if err := mergo.Merge(&config, configFile, mergo.WithOverride); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error merging config_file",
				Detail:   fmt.Sprintf("Error merging config_file with computed values %s", err),
			})
			return nil, diags
		}
	}
	// if no device_id is set by the user, we will set
	// the value based on the machine id, fallback to random UUID
	_, errs := validation.IsUUID(config.DeviceID, "device_id")
	if len(errs) > 0 {
		config.DeviceID = defaultDeviceID()
	}
	if err := config.Validate(usingFile); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Missing Appgate SDP credentials",
			Detail:   fmt.Sprintf("Appgate client is unauthenticated. %s", err.Error()),
		})
		return nil, diags
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

func defaultDeviceID() string {
	readAndParseUUID := func() (uuid.UUID, error) {
		// machine.ID() tries to read
		// /etc/machine-id on Linux
		// /etc/hostid on BSD
		// ioreg -rd1 -c IOPlatformExpertDevice | grep IOPlatformUUID on OSX
		// reg query HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Cryptography /v MachineGuid on Windows
		// and tries to parse the value as a UUID
		// https://github.com/denisbrodbeck/machineid
		var id uuid.UUID
		mid, err := machineid.ID()
		if err != nil {
			return id, err
		}
		return uuid.Parse(mid)
	}
	// if we can't get a valid UUID based on the machine ID, we will fallback to a random UUID value.
	v, err := readAndParseUUID()
	if err != nil {
		return uuid.New().String()
	}
	return v.String()
}
