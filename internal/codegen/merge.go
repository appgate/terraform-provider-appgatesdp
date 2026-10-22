package codegen

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// The provider code specification (IR) is handled as generic JSON so the
// merge stays independent of the many typed attribute structs in
// terraform-plugin-codegen-spec.

// versionedIR is the IR generated for one API version.
type versionedIR struct {
	version string
	ir      map[string]any
}

// mergeIR combines the per-version IRs into the single superset schema the
// provider exposes. Terraform fixes schemas before the provider is configured,
// so every attribute of every version must be present; version-specific
// support is enforced at plan time using the per-version API metadata.
// readOnly names data sources that read an endpoint instead of looking an
// object up by id or name.
func mergeIR(irs []versionedIR, readOnly map[string]bool) (map[string]any, error) {
	out := map[string]any{
		"provider": map[string]any{"name": "appgatesdp"},
		"version":  "0.1",
	}
	for _, kind := range []string{"resources", "datasources"} {
		byName := map[string][]versionedAttrs{}
		for _, v := range irs {
			for _, r := range listOf(v.ir[kind]) {
				rm := r.(map[string]any)
				name := rm["name"].(string)
				schema, _ := rm["schema"].(map[string]any)
				byName[name] = append(byName[name], versionedAttrs{v.version, listOf(schema["attributes"])})
			}
		}
		var merged []any
		for _, name := range sortedKeys(byName) {
			versions := byName[name]
			var present []string
			for _, v := range versions {
				present = append(present, v.version)
			}
			attrs, err := mergeAttributes(versions, present)
			if err != nil {
				return nil, fmt.Errorf("%s %s: %w", strings.TrimSuffix(kind, "s"), name, err)
			}
			for _, a := range attrs {
				_, body := attrType(a.(map[string]any))
				if _, nested := nestedAttributes(body); nested != nil {
					stripNestedDefaults(nested)
					if c := body["computed_optional_required"]; c == computed || c == computedOptional {
						computeChildren(nested)
					}
				}
			}
			switch {
			case kind == "resources":
				attrs = resourceOverlay(attrs)
			case readOnly[name]:
				attrs = readOnlyOverlay(attrs)
			default:
				attrs = dataSourceOverlay(attrs)
			}
			if err := checkReservedNames(attrs); err != nil {
				return nil, fmt.Errorf("%s %s: %w", strings.TrimSuffix(kind, "s"), name, err)
			}
			desc := fmt.Sprintf("Supported SDP Admin API versions: %s.", strings.Join(present, ", "))
			merged = append(merged, map[string]any{
				"name": name,
				"schema": map[string]any{
					"attributes":  attrs,
					"description": desc,
				},
			})
		}
		out[kind] = merged
	}
	return out, nil
}

type versionedAttrs struct {
	version string
	attrs   []any
}

// computability values used by the IR.
const (
	required         = "required"
	optional         = "optional"
	computed         = "computed"
	computedOptional = "computed_optional"
)

// mergeAttributes merges attribute lists of the given versions. parents lists
// every version in which the enclosing object exists.
func mergeAttributes(inputs []versionedAttrs, parents []string) ([]any, error) {
	type entry struct {
		version string
		typ     string
		body    map[string]any
	}
	byName := map[string][]entry{}
	for _, in := range inputs {
		for _, a := range in.attrs {
			am := a.(map[string]any)
			name, _ := am["name"].(string)
			typ, body := attrType(am)
			if typ == "" {
				return nil, fmt.Errorf("attribute %s has no type", name)
			}
			byName[name] = append(byName[name], entry{in.version, typ, body})
		}
	}
	var out []any
	for _, name := range sortedKeys(byName) {
		entries := byName[name]
		newest := entries[len(entries)-1]
		var compatible []entry
		for _, e := range entries {
			if e.typ == newest.typ {
				compatible = append(compatible, e)
			}
		}
		body := deepCopy(newest.body)
		var supported []string
		values := map[string]bool{}
		for _, e := range compatible {
			supported = append(supported, e.version)
			c, _ := e.body["computed_optional_required"].(string)
			values[c] = true
		}
		body["computed_optional_required"] = mergeComputability(values, len(supported) == len(parents))

		for _, key := range []string{"default", "validators", "plan_modifiers"} {
			first, has := compatible[0].body[key]
			same := true
			for _, e := range compatible[1:] {
				if v, ok := e.body[key]; ok != has || !reflect.DeepEqual(v, first) {
					same = false
				}
			}
			// A default must hold for every version: versions without the
			// attribute cannot store the defaulted value.
			if !same || (key == "default" && (body["computed_optional_required"] != computedOptional || len(supported) != len(parents))) {
				delete(body, key)
			}
		}
		for _, e := range compatible {
			if e.body["sensitive"] == true {
				body["sensitive"] = true
			}
		}
		if len(supported) != len(parents) {
			note := fmt.Sprintf("Supported in API versions: %s.", strings.Join(supported, ", "))
			desc, _ := body["description"].(string)
			body["description"] = strings.TrimSpace(desc + " " + note)
		}

		nestedKey, _ := nestedAttributes(body)
		if nestedKey != "" {
			var nested []versionedAttrs
			for _, e := range compatible {
				_, attrs := nestedAttributes(e.body)
				nested = append(nested, versionedAttrs{e.version, attrs})
			}
			merged, err := mergeAttributes(nested, supported)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", name, err)
			}
			setNestedAttributes(body, merged)
		}
		out = append(out, map[string]any{"name": name, newest.typ: body})
	}
	return out, nil
}

// stripNestedDefaults removes defaults below the root level.
// terraform-plugin-framework applies nested defaults wherever the
// configuration is null, including elements that only exist in prior state
// because the controller created them (e.g. a site's built-in DNS resolver
// uses update_interval 300 while the specification default is 60). The
// result would be a permanent diff, so unset nested attributes take the
// API's value instead.
func stripNestedDefaults(attrs []any) {
	for _, a := range attrs {
		_, body := attrType(a.(map[string]any))
		delete(body, "default")
		if _, nested := nestedAttributes(body); nested != nil {
			stripNestedDefaults(nested)
		}
	}
}

// computeChildren makes every attribute below a computed parent computed.
// The controller fills computed objects and lists itself (e.g. a site's
// built-in DNS resolver). Terraform core treats a non-computed child value
// in prior state as proof the configuration once set the parent, and plans
// to remove it when the configuration does not. Required children are still
// enforced per API version at plan time from the version metadata.
func computeChildren(attrs []any) {
	for _, a := range attrs {
		_, body := attrType(a.(map[string]any))
		if c := body["computed_optional_required"]; c == required || c == optional {
			body["computed_optional_required"] = computedOptional
		}
		if _, nested := nestedAttributes(body); nested != nil {
			computeChildren(nested)
		}
	}
}

// mergeComputability picks the least restrictive computability that is valid
// for every version: an attribute is only required if every version requires
// it, and only computed if no version accepts it as input.
func mergeComputability(values map[string]bool, everywhere bool) string {
	switch {
	case len(values) == 1 && values[required] && everywhere:
		return required
	case len(values) == 1 && values[computed]:
		return computed
	case !values[computed] && !values[computedOptional]:
		return optional
	default:
		return computedOptional
	}
}

func attrType(a map[string]any) (string, map[string]any) {
	for k, v := range a {
		if k != "name" {
			body, _ := v.(map[string]any)
			return k, body
		}
	}
	return "", nil
}

// nestedAttributes returns the nested attribute list of nested attribute types.
func nestedAttributes(body map[string]any) (string, []any) {
	if attrs, ok := body["attributes"].([]any); ok {
		return "attributes", attrs
	}
	if obj, ok := body["nested_object"].(map[string]any); ok {
		return "nested_object", listOf(obj["attributes"])
	}
	return "", nil
}

func setNestedAttributes(body map[string]any, attrs []any) {
	if _, ok := body["attributes"]; ok {
		body["attributes"] = attrs
		return
	}
	body["nested_object"].(map[string]any)["attributes"] = attrs
}

var (
	useStateForUnknown = map[string]any{"custom": map[string]any{
		"imports":           []any{map[string]any{"path": "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"}},
		"schema_definition": "stringplanmodifier.UseStateForUnknown()",
	}}
	requiresReplaceIfConfigured = map[string]any{"custom": map[string]any{
		"imports":           []any{map[string]any{"path": "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"}},
		"schema_definition": "stringplanmodifier.RequiresReplaceIfConfigured()",
	}}
)

// resourceOverlay guarantees a stable `id` attribute on every resource. The id
// may be chosen by the user on create (SDP accepts client generated UUIDs),
// otherwise it is assigned by the API and kept across plans.
func resourceOverlay(attrs []any) []any {
	id := findAttr(attrs, "id")
	if id == nil {
		id = map[string]any{"name": "id", "string": map[string]any{"computed_optional_required": computed, "description": "ID of the object."}}
		attrs = insertSorted(attrs, id)
	}
	body := id["string"].(map[string]any)
	mods := []any{useStateForUnknown}
	if body["computed_optional_required"] == computedOptional {
		mods = append(mods, requiresReplaceIfConfigured)
	}
	body["plan_modifiers"] = mods
	return attrs
}

// dataSourceOverlay lets data sources look objects up by id or by name.
func dataSourceOverlay(attrs []any) []any {
	id := findAttr(attrs, "id")
	if id == nil {
		attrs = insertSorted(attrs, map[string]any{"name": "id", "string": map[string]any{"computed_optional_required": computed, "description": "ID of the object."}})
	} else {
		id["string"].(map[string]any)["computed_optional_required"] = computedOptional
	}
	if name := findAttr(attrs, "name"); name != nil {
		if body, ok := name["string"].(map[string]any); ok {
			body["computed_optional_required"] = computedOptional
		}
	}
	return attrs
}

// readOnlyOverlay gives read-only data sources an id attribute (their
// parameters, including a path parameter called id, stay as generated).
func readOnlyOverlay(attrs []any) []any {
	if findAttr(attrs, "id") == nil {
		attrs = insertSorted(attrs, map[string]any{"name": "id", "string": map[string]any{"computed_optional_required": computed, "description": "Identifier of the data source."}})
	}
	return attrs
}

// reservedNames are root attribute names Terraform reserves for meta-arguments.
var reservedNames = []string{"connection", "count", "depends_on", "for_each", "lifecycle", "provider", "provisioner"}

func checkReservedNames(attrs []any) error {
	for _, a := range attrs {
		name := a.(map[string]any)["name"].(string)
		if slices.Contains(reservedNames, name) {
			return fmt.Errorf("attribute %q is a reserved Terraform name; add it to codegen/resources.yml ignore list", name)
		}
	}
	return nil
}

func findAttr(attrs []any, name string) map[string]any {
	for _, a := range attrs {
		if am := a.(map[string]any); am["name"] == name {
			return am
		}
	}
	return nil
}

func insertSorted(attrs []any, a map[string]any) []any {
	i, _ := slices.BinarySearchFunc(attrs, a["name"].(string), func(e any, name string) int {
		return strings.Compare(e.(map[string]any)["name"].(string), name)
	})
	return slices.Insert(attrs, i, any(a))
}

func deepCopy(m map[string]any) map[string]any {
	b, _ := json.Marshal(m)
	var out map[string]any
	_ = json.Unmarshal(b, &out)
	return out
}
