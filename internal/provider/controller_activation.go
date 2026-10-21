package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
	"github.com/appgate/terraform-provider-appgatesdp/internal/sdp"
)

// controllerActivationResource enables the controller function on an
// activated appliance and waits until the controller is ready.
//
// Hand-written because the waiting is behavior, not API shape: it updates
// the generated appliance object (controller and admin_interface use the
// generated appliance schema and per-version metadata) and then polls
// GET /appliances/status for the appliance state.
type controllerActivationResource struct {
	data *providerData
}

var (
	_ resource.ResourceWithConfigure   = &controllerActivationResource{}
	_ resource.ResourceWithImportState = &controllerActivationResource{}
	_ resource.ResourceWithModifyPlan  = &controllerActivationResource{}
)

// Appliance states reported by /appliances/status.
const (
	stateApplianceReady  = "appliance_ready"
	stateControllerReady = "controller_ready"
)

var (
	defaultActivationTimeout = 45 * time.Minute
	// statusPollInterval is the first poll interval; it grows to 30s.
	statusPollInterval = 5 * time.Second
)

// activationAttributes are the appliance attributes this resource manages.
var activationAttributes = []string{"controller", "admin_interface"}

func (r *controllerActivationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_appliance_controller_activation"
}

func (r *controllerActivationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	appliance := generatedResourceSchemas["appliance"](ctx)
	attrs := map[string]rschema.Attribute{
		"id": rschema.StringAttribute{Computed: true, Description: "ID of the appliance.",
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"appliance_id": rschema.StringAttribute{Required: true, Description: "ID of an activated appliance.",
			PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"wait_timeout": rschema.StringAttribute{Optional: true,
			Description: "How long to wait for the appliance to reach its target state, as a duration. Defaults to `45m`. If it is exceeded, the apply succeeds with a warning."},
	}
	for _, name := range activationAttributes {
		a := appliance.Attributes[name].(rschema.SingleNestedAttribute)
		a.Required, a.Optional, a.Computed = true, false, false
		attrs[name] = a
	}
	resp.Schema = rschema.Schema{
		Description: "Enables or disables the controller function of an activated appliance and waits until the appliance reaches the matching state (controller_ready or appliance_ready). The `controller` and `admin_interface` attributes are those of `appgatesdp_appliance`.",
		Attributes:  attrs,
	}
}

func (r *controllerActivationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		r.data, _ = req.ProviderData.(*providerData)
	}
}

// fields returns the appliance metadata of the managed attributes.
func (r *controllerActivationResource) fields(diags *diag.Diagnostics) map[string]*apimeta.Field {
	if r.data == nil {
		diags.AddError("Provider not configured", "The appgatesdp provider must be configured before managing appgatesdp_appliance_controller_activation.")
		return nil
	}
	appliance, ok := r.data.api.Resources["appliance"]
	if !ok {
		diags.AddError("Resource not supported by API version", "appliances are not available in SDP Admin API "+r.data.version.APIVersion)
		return nil
	}
	out := map[string]*apimeta.Field{}
	for _, name := range activationAttributes {
		out[name] = appliance.Fields[name]
	}
	return out
}

func (r *controllerActivationResource) ModifyPlan(_ context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() || r.data == nil {
		return
	}
	fields := r.fields(&resp.Diagnostics)
	if fields == nil {
		return
	}
	var config, plan map[string]tftypes.Value
	if req.Config.Raw.As(&config) != nil || req.Plan.Raw.As(&plan) != nil {
		return
	}
	for _, name := range activationAttributes {
		f := fields[name]
		checkConfigured(&resp.Diagnostics, path.Root(name), config[name], f.Fields, r.data.version.APIVersion)
		checkRequired(&resp.Diagnostics, path.Root(name), plan[name], f.Fields, r.data.version.APIVersion)
	}
	if d := stringAttr(req.Config.Raw, "wait_timeout"); d != "" {
		if _, err := time.ParseDuration(d); err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("wait_timeout"), "Invalid wait_timeout", err.Error())
		}
	}
}

func (r *controllerActivationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.apply(ctx, req.Plan.Raw, &resp.State, &resp.Diagnostics)
}

func (r *controllerActivationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.apply(ctx, req.Plan.Raw, &resp.State, &resp.Diagnostics)
}

// apply writes controller and admin_interface onto the appliance and waits
// for the resulting state.
func (r *controllerActivationResource) apply(ctx context.Context, plan tftypes.Value, state *tfsdk.State, diags *diag.Diagnostics) {
	fields := r.fields(diags)
	if fields == nil {
		return
	}
	id := stringAttr(plan, "appliance_id")
	appliance, err := r.getAppliance(ctx, id)
	if err != nil {
		diags.AddError("Cannot read appliance", err.Error())
		return
	}
	if appliance["activated"] != true {
		diags.AddError("Appliance not activated", fmt.Sprintf("The controller function can only be enabled on an activated appliance; seed appliance %q first.", appliance["name"]))
		return
	}
	if err := overlayJSON(appliance, plan, fields); err != nil {
		diags.AddError("Cannot convert plan to API request", err.Error())
		return
	}
	var out map[string]any
	if err := r.data.client.Do(ctx, http.MethodPut, "/appliances/"+url.PathEscape(id), nil, appliance, &out); err != nil {
		diags.AddError("Cannot update appliance", err.Error())
		return
	}
	target := stateApplianceReady
	if ctrl, _ := out["controller"].(map[string]any); ctrl["enabled"] == true {
		target = stateControllerReady
	}
	r.wait(ctx, id, target, stringAttr(plan, "wait_timeout"), diags)
	if diags.HasError() {
		return
	}
	r.setState(state, diags, out, fields, plan, afterApply)
}

func (r *controllerActivationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	fields := r.fields(&resp.Diagnostics)
	if fields == nil {
		return
	}
	appliance, err := r.getAppliance(ctx, stringAttr(req.State.Raw, "appliance_id"))
	if sdp.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Cannot read appliance", err.Error())
		return
	}
	r.setState(&resp.State, &resp.Diagnostics, appliance, fields, req.State.Raw, refresh)
}

// Delete disables the controller function and waits until the appliance is
// a plain appliance again.
func (r *controllerActivationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.fields(&resp.Diagnostics) == nil {
		return
	}
	id := stringAttr(req.State.Raw, "appliance_id")
	appliance, err := r.getAppliance(ctx, id)
	if sdp.IsNotFound(err) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Cannot read appliance", err.Error())
		return
	}
	ctrl, _ := appliance["controller"].(map[string]any)
	if ctrl == nil {
		ctrl = map[string]any{}
	}
	ctrl["enabled"] = false
	appliance["controller"] = ctrl
	if err := r.data.client.Do(ctx, http.MethodPut, "/appliances/"+url.PathEscape(id), nil, appliance, nil); err != nil {
		resp.Diagnostics.AddError("Cannot disable controller", err.Error())
		return
	}
	r.wait(ctx, id, stateApplianceReady, stringAttr(req.State.Raw, "wait_timeout"), &resp.Diagnostics)
}

func (r *controllerActivationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("appliance_id"), req.ID)...)
}

func (r *controllerActivationResource) getAppliance(ctx context.Context, id string) (map[string]any, error) {
	var appliance map[string]any
	err := r.data.client.Do(ctx, http.MethodGet, "/appliances/"+url.PathEscape(id), nil, nil, &appliance)
	return appliance, err
}

func (r *controllerActivationResource) setState(state *tfsdk.State, diags *diag.Diagnostics, appliance map[string]any, fields map[string]*apimeta.Field, prior tftypes.Value, md mode) {
	typ := prior.Type().(tftypes.Object)
	v, err := objectFromJSON(appliance, typ, fields, nil, prior, md)
	if err != nil {
		diags.AddError("Unexpected API response for appliance", err.Error())
		return
	}
	state.Raw = withAttr(v, "id", tftypes.NewValue(tftypes.String, stringAttr(v, "appliance_id")))
}

// wait polls /appliances/status until the appliance reports target. Running
// out of time is a warning, as the change itself has been applied.
func (r *controllerActivationResource) wait(ctx context.Context, id, target, timeout string, diags *diag.Diagnostics) {
	limit := defaultActivationTimeout
	if timeout != "" {
		limit, _ = time.ParseDuration(timeout)
	}
	deadline := time.Now().Add(limit)
	interval := statusPollInterval
	last := "unknown"
	for {
		var status struct {
			Data []map[string]any `json:"data"`
		}
		if err := r.data.client.Do(ctx, http.MethodGet, "/appliances/status", nil, nil, &status); err == nil {
			for _, a := range status.Data {
				if a["id"] == id {
					last = fmt.Sprint(a["state"])
				}
			}
			if last == target {
				return
			}
		}
		if time.Now().Add(interval).After(deadline) {
			diags.AddWarning("Appliance did not reach the expected state",
				fmt.Sprintf("Appliance %s is in state %s after %s; expected %s. The change was applied; check the appliance in the admin UI.", id, last, limit, target))
			return
		}
		select {
		case <-ctx.Done():
			diags.AddError("Interrupted while waiting for appliance", ctx.Err().Error())
			return
		case <-time.After(interval):
		}
		interval = min(interval*2, 30*time.Second)
	}
}
