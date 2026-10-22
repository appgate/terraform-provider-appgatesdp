package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion"
	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion/generated"
	"github.com/appgate/terraform-provider-appgatesdp/internal/sdp"
)

var (
	_ resource.ResourceWithConfigure   = &genericResource{}
	_ resource.ResourceWithImportState = &genericResource{}
	_ resource.ResourceWithModifyPlan  = &genericResource{}
)

// genericResource implements every generated resource. Its behavior is
// driven by the generated schema and the API metadata of the configured
// version; SDP specific behavior that the specification cannot express lives
// in hooks.go.
type genericResource struct {
	name     string
	schemaFn func(context.Context) rschema.Schema
	data     *providerData

	shapeOnce sync.Once
	shape     shape
}

func (r *genericResource) schemaShape() shape {
	r.shapeOnce.Do(func() { r.shape = resourceShape(r.schemaFn(context.Background()).Attributes) })
	return r.shape
}

func (r *genericResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.name
}

func (r *genericResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schemaFn(ctx)
}

func (r *genericResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	data, ok := req.ProviderData.(*providerData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("got %T", req.ProviderData))
		return
	}
	r.data = data
}

// meta returns the API metadata of this resource for the configured version.
func (r *genericResource) meta(diags *diag.Diagnostics) *apimeta.Resource {
	if r.data == nil {
		diags.AddError("Provider not configured", "The appgatesdp provider must be configured before managing "+TypePrefix+r.name+".")
		return nil
	}
	m, ok := r.data.api.Resources[r.name]
	if !ok {
		diags.AddError("Resource not supported by API version", unsupportedDetail(TypePrefix+r.name, r.data.version, func(a *apimeta.API) bool {
			_, ok := a.Resources[r.name]
			return ok
		}))
		return nil
	}
	return m
}

// unsupportedDetail explains which API versions offer a resource.
func unsupportedDetail(name string, v apiversion.Version, has func(*apimeta.API) bool) string {
	var supported []string
	for _, rv := range apiversion.Default().Versions {
		if api, err := generated.API(rv.APIVersion); err == nil && has(api) {
			supported = append(supported, rv.APIVersion)
		}
	}
	return fmt.Sprintf("%s is not available in SDP Admin API %s (SDP %s). Supported API versions: %s.",
		name, v.APIVersion, v.SDPVersion, strings.Join(supported, ", "))
}

func (r *genericResource) ModifyPlan(_ context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() || r.data == nil {
		return
	}
	m := r.meta(&resp.Diagnostics)
	if m == nil {
		return
	}
	resp.Diagnostics.Append(validateForVersion(req.Config.Raw, req.Plan.Raw, m.Fields, r.data.version.APIVersion)...)
	if m.Immutable && !req.State.Raw.IsNull() {
		resp.RequiresReplace = append(resp.RequiresReplace, changedInputs(req.State.Raw, req.Plan.Raw, m.Fields)...)
	}
}

// changedInputs returns the configurable top-level attributes whose planned
// value differs from state.
func changedInputs(state, plan tftypes.Value, fields map[string]*apimeta.Field) path.Paths {
	var before, after map[string]tftypes.Value
	if state.As(&before) != nil || plan.As(&after) != nil {
		return nil
	}
	var out path.Paths
	for _, name := range sortedNames(after) {
		if f := fields[name]; f != nil && !f.ReadOnly && after[name].IsKnown() && !after[name].Equal(before[name]) {
			out = append(out, path.Root(name))
		}
	}
	return out
}

func (r *genericResource) objectPath(m *apimeta.Resource, id string) string {
	return m.ObjectPath(url.PathEscape(id))
}

// reread returns the object as GET returns it. Write responses are not
// canonical: controllers return timestamps with a different precision and
// omit server-set fields (e.g. blacklistedAt) that later reads include, which
// would surface as changes on the next refresh.
func (r *genericResource) reread(ctx context.Context, m *apimeta.Resource, id string, written map[string]any) (map[string]any, error) {
	if id == "" && !m.Singleton {
		return written, nil
	}
	return r.getObject(ctx, m, id)
}

// getObject reads one object in the shape of the field metadata. APIs
// without an object endpoint (ReadList) return all objects in a list
// property of the collection.
func (r *genericResource) getObject(ctx context.Context, m *apimeta.Resource, id string) (map[string]any, error) {
	if m.ReadList == "" {
		var out map[string]any
		err := r.data.client.Do(ctx, http.MethodGet, r.objectPath(m, id), nil, nil, &out)
		return m.Wrap(out), err
	}
	var list map[string]any
	if err := r.data.client.Do(ctx, http.MethodGet, m.Path, nil, nil, &list); err != nil {
		return nil, err
	}
	items, _ := list[m.ListProperty()].([]any)
	for _, item := range items {
		if obj, ok := item.(map[string]any); ok && fmt.Sprint(obj[m.ID()]) == id {
			return obj, nil
		}
	}
	return nil, &sdp.Error{Status: http.StatusNotFound, Method: http.MethodGet, Path: m.Path, Message: fmt.Sprintf("%s %q not found in %s", m.ID(), id, m.ReadList)}
}

// put replaces an existing object with the plan. PUT replaces the whole
// object, so it starts from the current object: fields Terraform does not
// manage, or that are unknown in the plan, keep their values.
func (r *genericResource) put(ctx context.Context, m *apimeta.Resource, id string, plan tftypes.Value) (map[string]any, error) {
	current, err := r.getObject(ctx, m, id)
	if err != nil {
		return nil, fmt.Errorf("read before update: %w", err)
	}
	if current == nil {
		current = map[string]any{}
	}
	if err := overlayJSON(current, plan, m.Fields); err != nil {
		return nil, fmt.Errorf("convert plan to API request: %w", err)
	}
	body := m.Unwrap(current)
	setDiscriminator(body, m)
	if !m.Singleton {
		body["id"] = id
	}
	if err := r.data.client.Do(ctx, http.MethodPut, r.objectPath(m, id), nil, body, nil); err != nil {
		return nil, err
	}
	return r.reread(ctx, m, id, nil)
}

func (r *genericResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	m := r.meta(&resp.Diagnostics)
	if m == nil {
		return
	}
	c := r.data.client
	var out map[string]any
	var err error
	switch {
	case m.Singleton && m.CreateMethod == "":
		// Singletons always exist: configure them like an update so settings
		// the configuration does not mention keep their values.
		out, err = r.put(ctx, m, "", req.Plan.Raw)
	case m.Adopt:
		name := stringAttr(req.Plan.Raw, "name")
		var existing map[string]any
		if existing, err = findByName(ctx, c, m, name); err == nil {
			out, err = r.put(ctx, m, fmt.Sprint(existing["id"]), req.Plan.Raw)
		}
	default:
		var body map[string]any
		if body, err = objectToJSON(req.Plan.Raw, m.Fields); err != nil {
			resp.Diagnostics.AddError("Cannot convert plan to API request", err.Error())
			return
		}
		body = m.Unwrap(body)
		setDiscriminator(body, m)
		method := m.CreateMethod
		if method == "" {
			method = http.MethodPost
		}
		if err = c.Do(ctx, method, m.Path, nil, body, &out); err == nil {
			id, _ := out[m.ID()].(string)
			if id == "" {
				id = stringAttr(req.Plan.Raw, "id")
			}
			out, err = r.reread(ctx, m, id, m.Wrap(out))
		}
	}
	if err != nil {
		resp.Diagnostics.AddError("Cannot create "+TypePrefix+r.name, err.Error())
		return
	}
	r.setState(&resp.State, &resp.Diagnostics, out, m, req.Plan.Raw, afterApply)
}

func (r *genericResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	m := r.meta(&resp.Diagnostics)
	if m == nil {
		return
	}
	id := stringAttr(req.State.Raw, "id")
	out, err := r.getObject(ctx, m, id)
	if sdp.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Cannot read "+TypePrefix+r.name, err.Error())
		return
	}
	if !matchesDiscriminator(out, m) {
		resp.Diagnostics.AddWarning("Object type changed",
			fmt.Sprintf("%s %s is no longer of %s %q; removing it from state.", TypePrefix+r.name, id, m.Discriminator.Property, m.Discriminator.Value))
		resp.State.RemoveResource(ctx)
		return
	}
	r.setState(&resp.State, &resp.Diagnostics, out, m, req.State.Raw, refresh)
}

func (r *genericResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	m := r.meta(&resp.Diagnostics)
	if m == nil {
		return
	}
	id := stringAttr(req.State.Raw, "id")
	var out map[string]any
	var err error
	if m.Immutable {
		// Every input change replaces the object (ModifyPlan); only computed
		// values can differ here, so the current object is the new state.
		out, err = r.getObject(ctx, m, id)
	} else {
		out, err = r.put(ctx, m, id, req.Plan.Raw)
	}
	if err != nil {
		resp.Diagnostics.AddError("Cannot update "+TypePrefix+r.name, err.Error())
		return
	}
	r.setState(&resp.State, &resp.Diagnostics, out, m, req.Plan.Raw, afterApply)
}

func (r *genericResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	m := r.meta(&resp.Diagnostics)
	if m == nil {
		return
	}
	switch {
	case m.Adopt:
		resp.Diagnostics.AddWarning("Built-in object kept", TypePrefix+r.name+" is built into SDP; it was removed from Terraform state but not deleted.")
		return
	case m.Singleton && !m.Delete:
		return
	}
	id := stringAttr(req.State.Raw, "id")
	if h := resourceHooks[r.name].beforeDelete; h != nil {
		if err := h(ctx, r.data.client, m, id); err != nil {
			resp.Diagnostics.AddError("Cannot delete "+TypePrefix+r.name, err.Error())
			return
		}
	}
	// DELETE on a singleton resets or disables it. A singleton that is
	// already disabled answers 412 (e.g. replication source not enabled).
	err := r.data.client.Do(ctx, http.MethodDelete, m.DeleteObjectPath(url.PathEscape(id)), nil, nil, nil)
	var apiErr *sdp.Error
	if errors.As(err, &apiErr) && apiErr.Status == http.StatusPreconditionFailed && m.Singleton {
		err = nil
	}
	if err != nil && !sdp.IsNotFound(err) {
		resp.Diagnostics.AddError("Cannot delete "+TypePrefix+r.name, err.Error())
	}
}

func (r *genericResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// setState converts an API response into state. Singletons have no API id;
// they use the resource name as a fixed id.
func (r *genericResource) setState(state *tfsdk.State, diags *diag.Diagnostics, out map[string]any, m *apimeta.Resource, prior tftypes.Value, md mode) {
	typ, ok := prior.Type().(tftypes.Object)
	if !ok {
		diags.AddError("Unexpected state type", prior.Type().String())
		return
	}
	v, err := objectFromJSON(out, typ, m.Fields, r.schemaShape(), prior, md)
	if err != nil {
		diags.AddError("Unexpected API response for "+TypePrefix+r.name, err.Error())
		return
	}
	if stringAttr(v, "id") == "" {
		fallback, _ := out[m.ID()].(string)
		if fallback == "" {
			fallback = stringAttr(prior, "id")
		}
		if fallback == "" {
			if !m.Singleton {
				diags.AddError("Unexpected API response for "+TypePrefix+r.name, "the response contains no id")
				return
			}
			fallback = r.name
		}
		v = withAttr(v, "id", tftypes.NewValue(tftypes.String, fallback))
	}
	state.Raw = v
}

// overlayJSON writes the planned values onto an API object in place.
func overlayJSON(current map[string]any, plan tftypes.Value, fields map[string]*apimeta.Field) error {
	var attrs map[string]tftypes.Value
	if err := plan.As(&attrs); err != nil {
		return err
	}
	for _, name := range sortedNames(attrs) {
		av := attrs[name]
		f := fields[name]
		if f == nil || f.ReadOnly || !av.IsKnown() || !typeMatches(av.Type(), f) {
			continue
		}
		if av.IsNull() {
			delete(current, f.JSON)
			continue
		}
		if av.Type().Is(tftypes.Object{}) {
			sub, _ := current[f.JSON].(map[string]any)
			if sub == nil {
				sub = map[string]any{}
			}
			if err := overlayJSON(sub, av, f.Fields); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
			current[f.JSON] = sub
			continue
		}
		j, err := valueToJSON(av, f)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		current[f.JSON] = j
	}
	return nil
}

func setDiscriminator(body map[string]any, m *apimeta.Resource) {
	if m.Discriminator != nil {
		body[m.Discriminator.Property] = m.Discriminator.Value
	}
}

func matchesDiscriminator(obj map[string]any, m *apimeta.Resource) bool {
	if m.Discriminator == nil {
		return true
	}
	v, ok := obj[m.Discriminator.Property]
	return !ok || v == m.Discriminator.Value
}

// findByName returns the single object of the collection with exactly the
// given name (and the resource's discriminator value).
func findByName(ctx context.Context, c *sdp.Client, m *apimeta.Resource, name string) (map[string]any, error) {
	const pageSize = 100
	var matches []map[string]any
	for start := 0; ; start += pageSize {
		q := url.Values{"query": {name}, "orderBy": {"name"}, "range": {fmt.Sprintf("%d-%d", start, start+pageSize)}}
		var page struct {
			Data []map[string]any `json:"data"`
		}
		if err := c.Do(ctx, http.MethodGet, m.Path, q, nil, &page); err != nil {
			return nil, err
		}
		for _, obj := range page.Data {
			if obj["name"] == name && matchesDiscriminator(obj, m) {
				matches = append(matches, obj)
			}
		}
		if len(page.Data) < pageSize {
			break
		}
	}
	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("no object named %q found in %s", name, m.Path)
	case 1:
		return matches[0], nil
	default:
		return nil, fmt.Errorf("%d objects named %q found in %s; look it up by id instead", len(matches), name, m.Path)
	}
}

// stringAttr returns a known string attribute of an object value or "".
func stringAttr(v tftypes.Value, name string) string {
	var attrs map[string]tftypes.Value
	if !v.IsKnown() || v.IsNull() || v.As(&attrs) != nil {
		return ""
	}
	a, ok := attrs[name]
	if !ok || !a.IsKnown() || a.IsNull() {
		return ""
	}
	var s string
	_ = a.As(&s)
	return s
}

// withAttr returns a copy of object v with attribute name set to a.
func withAttr(v tftypes.Value, name string, a tftypes.Value) tftypes.Value {
	var attrs map[string]tftypes.Value
	_ = v.As(&attrs)
	attrs[name] = a
	return tftypes.NewValue(v.Type(), attrs)
}
