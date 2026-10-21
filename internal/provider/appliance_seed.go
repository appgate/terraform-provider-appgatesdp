package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// applianceSeedDataSource exports the activation seed of an appliance.
//
// Hand-written: the seed export is an action (POST /appliances/{id}/export)
// returning an opaque document, which the OpenAPI generator cannot map to a
// resource or data source. Its request body (SSHConfig) is identical in every
// supported API version.
type applianceSeedDataSource struct {
	data *providerData
}

type applianceSeedModel struct {
	ID                 types.String `tfsdk:"id"`
	ApplianceID        types.String `tfsdk:"appliance_id"`
	Password           types.String `tfsdk:"password"`
	SSHKey             types.String `tfsdk:"ssh_key"`
	ProvideCloudSSHKey types.Bool   `tfsdk:"provide_cloud_ssh_key"`
	AllowCustomization types.Bool   `tfsdk:"allow_customization"`
	ValidityDays       types.Int64  `tfsdk:"validity_days"`
	SeedFile           types.String `tfsdk:"seed_file"`
}

func (d *applianceSeedDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_appliance_seed"
}

func (d *applianceSeedDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Exports the JSON seed used to activate an inactive appliance.",
		Attributes: map[string]schema.Attribute{
			"id":                    schema.StringAttribute{Computed: true, Description: "Same as appliance_id."},
			"appliance_id":          schema.StringAttribute{Required: true, Description: "ID of the inactive appliance."},
			"password":              schema.StringAttribute{Optional: true, Sensitive: true, Description: "Password of the appliance's cz user."},
			"ssh_key":               schema.StringAttribute{Optional: true, Description: "SSH public key allowed to log in as cz."},
			"provide_cloud_ssh_key": schema.BoolAttribute{Optional: true, Description: "Use the SSH key provided by the cloud platform (AWS, Azure)."},
			"allow_customization":   schema.BoolAttribute{Optional: true, Description: "Whether the appliance allows customizations. Defaults to true."},
			"validity_days":         schema.Int64Attribute{Optional: true, Description: "How many days the seed is valid. Defaults to 1."},
			"seed_file": schema.StringAttribute{Computed: true, Sensitive: true,
				Description: "Base64 encoded seed JSON; empty when the appliance is already activated."},
		},
	}
}

func (d *applianceSeedDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		d.data, _ = req.ProviderData.(*providerData)
	}
}

func (d *applianceSeedDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.data == nil {
		resp.Diagnostics.AddError("Provider not configured", "The appgatesdp provider must be configured before reading appgatesdp_appliance_seed.")
		return
	}
	var m applianceSeedModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	path := "/appliances/" + url.PathEscape(m.ApplianceID.ValueString())
	var appliance map[string]any
	if err := d.data.client.Do(ctx, http.MethodGet, path, nil, nil, &appliance); err != nil {
		resp.Diagnostics.AddError("Cannot read appliance", err.Error())
		return
	}
	m.ID = m.ApplianceID
	m.SeedFile = types.StringValue("")
	if appliance["activated"] == true {
		resp.Diagnostics.AddWarning("Appliance already activated", fmt.Sprintf("Appliance %s is activated; no seed was exported.", m.ApplianceID.ValueString()))
		resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
		return
	}
	body := map[string]any{}
	if !m.Password.IsNull() {
		body["password"] = m.Password.ValueString()
	}
	if !m.SSHKey.IsNull() {
		body["sshKey"] = m.SSHKey.ValueString()
	}
	if !m.ProvideCloudSSHKey.IsNull() {
		body["provideCloudSSHKey"] = m.ProvideCloudSSHKey.ValueBool()
	}
	if !m.AllowCustomization.IsNull() {
		body["allowCustomization"] = m.AllowCustomization.ValueBool()
	}
	if !m.ValidityDays.IsNull() {
		body["validityDays"] = m.ValidityDays.ValueInt64()
	}
	var seed map[string]any
	if err := d.data.client.Do(ctx, http.MethodPost, path+"/export", nil, body, &seed); err != nil {
		resp.Diagnostics.AddError("Cannot export appliance seed", err.Error())
		return
	}
	raw, err := json.Marshal(seed)
	if err != nil {
		resp.Diagnostics.AddError("Cannot encode appliance seed", err.Error())
		return
	}
	m.SeedFile = types.StringValue(base64.StdEncoding.EncodeToString(raw))
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
