package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
)

// validateForVersion reports configuration the selected API version cannot
// accept. The Terraform schema is the superset of every API version, so these
// checks run at plan time once the provider (and thus the version) is known:
//   - attributes set in config that the version does not have, or has as
//     read-only or with a different type;
//   - attributes the version requires that are null in the plan.
func validateForVersion(config, plan tftypes.Value, fields map[string]*apimeta.Field, version string) diag.Diagnostics {
	var diags diag.Diagnostics
	checkConfigured(&diags, path.Empty(), config, fields, version)
	checkRequired(&diags, path.Empty(), plan, fields, version)
	return diags
}

func checkConfigured(diags *diag.Diagnostics, p path.Path, v tftypes.Value, fields map[string]*apimeta.Field, version string) {
	if v.IsNull() || !v.IsKnown() {
		return
	}
	var attrs map[string]tftypes.Value
	if v.As(&attrs) != nil {
		return
	}
	for _, name := range sortedNames(attrs) {
		av := attrs[name]
		if av.IsNull() {
			continue
		}
		ap := p.AtName(name)
		f := fields[name]
		switch {
		case f == nil || !typeMatches(av.Type(), f):
			if name == "id" && p.Equal(path.Empty()) {
				continue
			}
			diags.AddAttributeError(ap, "Attribute not supported by API version",
				fmt.Sprintf("%s is not available in SDP Admin API %s. Remove it from the configuration or select an api_version that supports it (see the provider's API versions guide).", ap, version))
		case f.ReadOnly:
			diags.AddAttributeError(ap, "Attribute is read-only in API version",
				fmt.Sprintf("%s is read-only in SDP Admin API %s and cannot be configured.", ap, version))
		default:
			forEachObject(ap, av, f, func(op path.Path, ov tftypes.Value) {
				checkConfigured(diags, op, ov, f.Fields, version)
			})
		}
	}
}

func checkRequired(diags *diag.Diagnostics, p path.Path, v tftypes.Value, fields map[string]*apimeta.Field, version string) {
	if v.IsNull() || !v.IsKnown() {
		return
	}
	var attrs map[string]tftypes.Value
	if v.As(&attrs) != nil {
		return
	}
	for _, name := range sortedNames(attrs) {
		av := attrs[name]
		f := fields[name]
		if f == nil || !typeMatches(av.Type(), f) {
			continue
		}
		ap := p.AtName(name)
		if f.Required && !f.ReadOnly && av.IsNull() {
			diags.AddAttributeError(ap, "Missing attribute required by API version",
				fmt.Sprintf("%s is required by SDP Admin API %s.", ap, version))
			continue
		}
		forEachObject(ap, av, f, func(op path.Path, ov tftypes.Value) {
			checkRequired(diags, op, ov, f.Fields, version)
		})
	}
}

// forEachObject calls fn for v if it is an object, or for each object element
// of a list, set or map.
func forEachObject(p path.Path, v tftypes.Value, f *apimeta.Field, fn func(path.Path, tftypes.Value)) {
	if f.Fields == nil || v.IsNull() || !v.IsKnown() {
		return
	}
	t := v.Type()
	switch {
	case t.Is(tftypes.Object{}):
		fn(p, v)
	case t.Is(tftypes.List{}):
		var elems []tftypes.Value
		_ = v.As(&elems)
		for i, e := range elems {
			fn(p.AtListIndex(i), e)
		}
	case t.Is(tftypes.Set{}):
		var elems []tftypes.Value
		_ = v.As(&elems)
		for _, e := range elems {
			// Set elements have no stable address; report on the set.
			fn(p, e)
		}
	case t.Is(tftypes.Map{}):
		var elems map[string]tftypes.Value
		_ = v.As(&elems)
		for _, k := range sortedNames(elems) {
			fn(p.AtMapKey(k), elems[k])
		}
	}
}
