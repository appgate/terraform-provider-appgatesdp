// Package provider implements the appgatesdp Terraform provider on top of
// the generated schemas (internal/generated) and per-version API metadata
// (internal/apiversion/generated).
package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion"
	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion/generated"
	"github.com/appgate/terraform-provider-appgatesdp/internal/sdp"
)

// TypePrefix is prepended to every resource and data source name.
const TypePrefix = "appgatesdp_"

// New returns the provider factory.
func New(version string) func() provider.Provider {
	return func() provider.Provider { return &appgateProvider{version: version} }
}

type appgateProvider struct {
	version string
	// newClient is replaced in tests.
	newClient func(sdp.Config, int) (*sdp.Client, error)
}

// providerData is handed to every resource and data source. It binds the
// configured API version to its client and generated metadata, so version
// selection happens exactly once.
type providerData struct {
	client  *sdp.Client
	version apiversion.Version
	api     *apimeta.API
}

type providerModel struct {
	URL          types.String `tfsdk:"url"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
	Provider     types.String `tfsdk:"provider"`
	Insecure     types.Bool   `tfsdk:"insecure"`
	Debug        types.Bool   `tfsdk:"debug"`
	APIVersion   types.String `tfsdk:"api_version"`
	ClientVer    types.Int64  `tfsdk:"client_version"`
	ConfigPath   types.String `tfsdk:"config_path"`
	PEMFilePath  types.String `tfsdk:"pem_filepath"`
	BearerToken  types.String `tfsdk:"bearer_token"`
	DeviceID     types.String `tfsdk:"device_id"`
	LoginTimeout types.String `tfsdk:"login_timeout"`
}

func (p *appgateProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "appgatesdp"
	resp.Version = p.version
}

func (p *appgateProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	reg := apiversion.Default()
	resp.Schema = schema.Schema{
		Description: "Manages Appgate SDP through the Admin API. Supported API versions: " + strings.Join(reg.Names(), ", ") + ".",
		Attributes: map[string]schema.Attribute{
			"url":      schema.StringAttribute{Optional: true, Description: "Controller admin URL, e.g. https://controller.example.com:8443/admin. Env: APPGATE_ADDRESS."},
			"username": schema.StringAttribute{Optional: true, Description: "Admin username. Env: APPGATE_USERNAME."},
			"password": schema.StringAttribute{Optional: true, Sensitive: true, Description: "Admin password. Env: APPGATE_PASSWORD."},
			"provider": schema.StringAttribute{Optional: true, Description: "Identity provider used to log in. Defaults to `local`. Env: APPGATE_PROVIDER."},
			"insecure": schema.BoolAttribute{Optional: true, Description: "Skip TLS certificate verification. Defaults to `true` for compatibility with earlier releases; set `false` together with `pem_filepath` in production. Env: APPGATE_INSECURE."},
			"debug":    schema.BoolAttribute{Optional: true, DeprecationMessage: "HTTP debugging is controlled by TF_LOG=DEBUG; this attribute has no effect.", Description: "Deprecated. Env: APPGATE_HTTP_DEBUG."},
			"api_version": schema.StringAttribute{Optional: true, Description: fmt.Sprintf(
				"SDP Admin API version, one of %s. When unset, the newest version supported by both the controller and the provider is negotiated at login. Env: APPGATE_API_VERSION.",
				strings.Join(reg.Names(), ", "))},
			"client_version": schema.Int64Attribute{Optional: true, DeprecationMessage: "Use api_version instead, e.g. api_version = \"v25\".", Description: "Deprecated numeric API version. Env: APPGATE_CLIENT_VERSION."},
			"config_path":    schema.StringAttribute{Optional: true, Description: "Path to a JSON file with provider settings (keys appgate_url, appgate_username, ...). Env: APPGATE_CONFIG_PATH."},
			"pem_filepath":   schema.StringAttribute{Optional: true, Description: "Path to the controller CA certificate in PEM format. Env: APPGATE_PEM_FILEPATH."},
			"bearer_token":   schema.StringAttribute{Optional: true, Sensitive: true, Description: "Admin API token obtained outside Terraform; replaces username and password. Env: APPGATE_BEARER_TOKEN."},
			"device_id":      schema.StringAttribute{Optional: true, Description: "UUID identifying this client in login requests. Defaults to a UUID derived from the machine ID. Env: APPGATE_DEVICE_ID."},
			"login_timeout":  schema.StringAttribute{Optional: true, Description: "How long to retry login while the controller is starting, as a duration (default `10m`). Env: APPGATE_LOGIN_TIMEOUT."},
		},
	}
}

// fileConfig is the format of config_path.
type fileConfig struct {
	URL          string `json:"appgate_url"`
	Username     string `json:"appgate_username"`
	Password     string `json:"appgate_password"`
	Provider     string `json:"appgate_provider"`
	APIVersion   string `json:"appgate_api_version"`
	Version      int    `json:"appgate_client_version"`
	BearerToken  string `json:"appgate_bearer_token"`
	PEMFilePath  string `json:"appgate_pem_filepath"`
	DeviceID     string `json:"appgate_device_id"`
	LoginTimeout string `json:"appgate_login_timeout"`
}

// settings are the resolved provider settings; precedence is config file,
// then provider block, then environment.
type settings struct {
	sdp.Config
	apiVersion string
}

func resolveSettings(m providerModel) (settings, diag.Diagnostics) {
	var diags diag.Diagnostics
	str := func(v types.String, env string) string {
		if !v.IsNull() && !v.IsUnknown() {
			return v.ValueString()
		}
		return os.Getenv(env)
	}
	boolean := func(v types.Bool, env string, def bool) bool {
		if !v.IsNull() && !v.IsUnknown() {
			return v.ValueBool()
		}
		if b, err := strconv.ParseBool(os.Getenv(env)); err == nil {
			return b
		}
		return def
	}
	s := settings{Config: sdp.Config{
		URL:         str(m.URL, "APPGATE_ADDRESS"),
		Username:    str(m.Username, "APPGATE_USERNAME"),
		Password:    str(m.Password, "APPGATE_PASSWORD"),
		Provider:    str(m.Provider, "APPGATE_PROVIDER"),
		BearerToken: str(m.BearerToken, "APPGATE_BEARER_TOKEN"),
		PEMFile:     str(m.PEMFilePath, "APPGATE_PEM_FILEPATH"),
		DeviceID:    str(m.DeviceID, "APPGATE_DEVICE_ID"),
		Insecure:    boolean(m.Insecure, "APPGATE_INSECURE", true),
	}}
	s.apiVersion = str(m.APIVersion, "APPGATE_API_VERSION")
	clientVersion := os.Getenv("APPGATE_CLIENT_VERSION")
	if !m.ClientVer.IsNull() && !m.ClientVer.IsUnknown() {
		clientVersion = strconv.FormatInt(m.ClientVer.ValueInt64(), 10)
	}
	loginTimeout := str(m.LoginTimeout, "APPGATE_LOGIN_TIMEOUT")

	if cfgPath := str(m.ConfigPath, "APPGATE_CONFIG_PATH"); cfgPath != "" {
		b, err := os.ReadFile(cfgPath)
		if err != nil {
			diags.AddAttributeError(path.Root("config_path"), "Cannot read config_path", err.Error())
			return s, diags
		}
		var fc fileConfig
		if err := json.Unmarshal(b, &fc); err != nil {
			diags.AddAttributeError(path.Root("config_path"), "Invalid config_path file", err.Error())
			return s, diags
		}
		override := func(dst *string, v string) {
			if v != "" {
				*dst = v
			}
		}
		override(&s.URL, fc.URL)
		override(&s.Username, fc.Username)
		override(&s.Password, fc.Password)
		override(&s.Provider, fc.Provider)
		override(&s.BearerToken, fc.BearerToken)
		override(&s.PEMFile, fc.PEMFilePath)
		override(&s.DeviceID, fc.DeviceID)
		override(&s.apiVersion, fc.APIVersion)
		override(&loginTimeout, fc.LoginTimeout)
		if fc.Version > 0 {
			clientVersion = strconv.Itoa(fc.Version)
		}
	}

	if s.apiVersion == "" && clientVersion != "" {
		s.apiVersion = "v" + clientVersion
	}
	s.LoginTimeout = 10 * time.Minute
	if loginTimeout != "" {
		d, err := time.ParseDuration(loginTimeout)
		if err != nil {
			diags.AddAttributeError(path.Root("login_timeout"), "Invalid login_timeout", fmt.Sprintf("%q is not a duration: %s", loginTimeout, err))
		}
		s.LoginTimeout = d
	}
	if s.DeviceID == "" {
		s.DeviceID = defaultDeviceID()
	} else if _, err := uuid.Parse(s.DeviceID); err != nil {
		diags.AddAttributeError(path.Root("device_id"), "Invalid device_id", "device_id must be a UUID")
	}
	return s, diags
}

func (p *appgateProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var m providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	s, diags := resolveSettings(m)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	s.UserAgent = fmt.Sprintf("terraform-provider-appgatesdp/%s (+https://registry.terraform.io/providers/appgate/appgatesdp) Terraform/%s", p.version, req.TerraformVersion)

	reg := apiversion.Default()
	var selected apiversion.Version
	if s.apiVersion != "" {
		v, err := reg.Lookup(s.apiVersion)
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("api_version"), "Unsupported API version", err.Error())
			return
		}
		selected = v
	}
	newClient := p.newClient
	if newClient == nil {
		newClient = sdp.New
	}
	client, err := newClient(s.Config, selected.Number())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Appgate SDP provider configuration", err.Error())
		return
	}
	if s.apiVersion == "" {
		selected, err = negotiate(ctx, client, reg)
		if err != nil {
			resp.Diagnostics.AddError("Cannot select SDP Admin API version", err.Error())
			return
		}
		client.SetVersion(selected.Number())
	}
	api, err := generated.API(selected.APIVersion)
	if err != nil {
		resp.Diagnostics.AddError("Missing generated API implementation", err.Error())
		return
	}
	data := &providerData{client: client, version: selected, api: api}
	resp.ResourceData = data
	resp.DataSourceData = data
}

// negotiate picks the newest API version supported by both the controller and
// the provider.
func negotiate(ctx context.Context, c *sdp.Client, reg apiversion.Registry) (apiversion.Version, error) {
	lo, hi, err := c.SupportedVersions(ctx)
	if err != nil {
		return apiversion.Version{}, fmt.Errorf("%w; set api_version explicitly", err)
	}
	for i := len(reg.Versions) - 1; i >= 0; i-- {
		if n := reg.Versions[i].Number(); n >= lo && n <= hi {
			return reg.Versions[i], nil
		}
	}
	return apiversion.Version{}, fmt.Errorf("controller supports API v%d to v%d, this provider supports %s", lo, hi, strings.Join(reg.Names(), ", "))
}

func (p *appgateProvider) Resources(context.Context) []func() resource.Resource {
	// Hand-written resources for API behavior the generator cannot map.
	out := []func() resource.Resource{
		func() resource.Resource { return &controllerActivationResource{} },
	}
	for _, name := range sortedNames(generatedResourceSchemas) {
		out = append(out, func() resource.Resource {
			return &genericResource{name: name, schemaFn: generatedResourceSchemas[name]}
		})
	}
	return out
}

func (p *appgateProvider) DataSources(context.Context) []func() datasource.DataSource {
	// Hand-written data sources for API actions the generator cannot map.
	out := []func() datasource.DataSource{
		func() datasource.DataSource { return &applianceSeedDataSource{} },
	}
	for _, name := range sortedNames(generatedDataSourceSchemas) {
		out = append(out, func() datasource.DataSource {
			return &genericDataSource{name: name, schemaFn: generatedDataSourceSchemas[name]}
		})
	}
	return out
}
