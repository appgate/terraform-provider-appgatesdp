package codegen

import (
	"encoding/json"
	"strings"
	"testing"
)

func attr(name, typ, computability string, extra map[string]any) map[string]any {
	body := map[string]any{"computed_optional_required": computability}
	for k, v := range extra {
		body[k] = v
	}
	return map[string]any{"name": name, typ: body}
}

func TestMergeAttributes(t *testing.T) {
	enumAB := []any{map[string]any{"custom": map[string]any{"schema_definition": `stringvalidator.OneOf("a","b")`}}}
	enumABC := []any{map[string]any{"custom": map[string]any{"schema_definition": `stringvalidator.OneOf("a","b","c")`}}}
	v1 := []any{
		attr("name", "string", required, nil),
		attr("old", "string", optional, nil),
		attr("mode", "string", computedOptional, map[string]any{"default": map[string]any{"static": "a"}, "validators": enumAB}),
		attr("size", "string", optional, nil),
		attr("both", "string", required, nil),
	}
	v2 := []any{
		attr("name", "string", required, nil),
		attr("mode", "string", computedOptional, map[string]any{"default": map[string]any{"static": "a"}, "validators": enumABC}),
		attr("size", "int64", optional, nil),
		attr("both", "string", required, nil),
		attr("new", "string", required, map[string]any{"default": map[string]any{"static": "x"}}),
		attr("secret", "string", optional, map[string]any{"sensitive": true}),
	}
	got, err := mergeAttributes([]versionedAttrs{{"v1", v1}, {"v2", v2}}, []string{"v1", "v2"})
	if err != nil {
		t.Fatal(err)
	}
	byName := map[string]map[string]any{}
	for _, a := range got {
		am := a.(map[string]any)
		typ, body := attrType(am)
		body["_type"] = typ
		byName[am["name"].(string)] = body
	}
	check := func(name, key string, want any) {
		t.Helper()
		b, _ := json.Marshal(byName[name][key])
		w, _ := json.Marshal(want)
		if string(b) != string(w) {
			t.Errorf("%s.%s = %s, want %s", name, key, b, w)
		}
	}
	check("name", "computed_optional_required", required)
	check("both", "computed_optional_required", required)
	// Present only in some versions: never required, documented availability.
	check("new", "computed_optional_required", optional)
	check("old", "computed_optional_required", optional)
	if !strings.Contains(byName["new"]["description"].(string), "Supported in API versions: v2.") {
		t.Error("availability note missing")
	}
	// Defaults only survive when every version has the attribute.
	check("mode", "default", map[string]any{"static": "a"})
	check("new", "default", nil)
	// Validators that differ between versions are dropped (the API validates).
	check("mode", "validators", nil)
	// Type changes: the newest type wins; older versions are type-incompatible.
	check("size", "_type", "int64")
	if !strings.Contains(byName["size"]["description"].(string), "Supported in API versions: v2.") {
		t.Error("type-incompatible versions must not be listed as supported")
	}
	check("secret", "sensitive", true)
}

func TestMergeComputability(t *testing.T) {
	for _, tc := range []struct {
		values     []string
		everywhere bool
		want       string
	}{
		{[]string{required}, true, required},
		{[]string{required}, false, optional},
		{[]string{required, optional}, true, optional},
		{[]string{computed}, false, computed},
		{[]string{computed, optional}, true, computedOptional},
		{[]string{computedOptional, required}, true, computedOptional},
	} {
		set := map[string]bool{}
		for _, v := range tc.values {
			set[v] = true
		}
		if got := mergeComputability(set, tc.everywhere); got != tc.want {
			t.Errorf("%v everywhere=%v: got %s, want %s", tc.values, tc.everywhere, got, tc.want)
		}
	}
}

func TestOverlays(t *testing.T) {
	res := resourceOverlay([]any{attr("name", "string", required, nil)})
	id := findAttr(res, "id")
	if id == nil || len(id["string"].(map[string]any)["plan_modifiers"].([]any)) != 1 {
		t.Fatalf("resources without id get a computed id with UseStateForUnknown: %v", res)
	}
	res = resourceOverlay([]any{attr("id", "string", computedOptional, nil)})
	if mods := findAttr(res, "id")["string"].(map[string]any)["plan_modifiers"].([]any); len(mods) != 2 {
		t.Fatal("configurable ids must require replacement when changed")
	}
	if single := dataSourceOverlay([]any{attr("mode", "string", computed, nil)}); findAttr(single, "id")["string"].(map[string]any)["computed_optional_required"] != computed {
		t.Error("singleton data sources get a computed id")
	}
	ds := dataSourceOverlay([]any{attr("id", "string", computedOptional, nil), attr("name", "string", computed, nil), attr("other", "string", computed, nil)})
	for _, name := range []string{"id", "name"} {
		if c := findAttr(ds, name)["string"].(map[string]any)["computed_optional_required"]; c != computedOptional {
			t.Errorf("data source %s must be a lookup argument, got %v", name, c)
		}
	}
	if err := checkReservedNames([]any{attr("count", "int64", optional, nil)}); err == nil {
		t.Error("reserved names must be rejected")
	}
}

func TestTerraformIdentifier(t *testing.T) {
	// Cases from terraform-plugin-codegen-openapi's TerraformIdentifier tests
	// plus SDP property names.
	for in, want := range map[string]string{
		"thing":               "thing",
		"camelCase":           "camel_case",
		"PascalCase":          "pascal_case",
		"kebab-case":          "kebabcase",
		"123thing":            "thing",
		"onBoarding2FA":       "on_boarding2fa",
		"ipv4":                "ipv4",
		"adminInterfaceHTTPS": "admin_interface_https",
		"vpn.weight":          "vpnweight",
	} {
		if got := TerraformIdentifier(in); got != want {
			t.Errorf("TerraformIdentifier(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPruneFile(t *testing.T) {
	src := []byte(`// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package resource_x

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func XResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{Attributes: map[string]schema.Attribute{
		"a": schema.SingleNestedAttribute{
			Attributes: map[string]schema.Attribute{},
			CustomType: AType{ObjectType: types.ObjectType{}},
		},
	}}
}

type XModel struct{ A types.Object }
type AType struct{ types.ObjectType }
func (t AType) String() string { return "" }
`)
	out, err := pruneFile(src)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, gone := range []string{"CustomType", "XModel", "AType", "terraform-plugin-framework/types\""} {
		if strings.Contains(s, gone) {
			t.Errorf("pruned output still contains %q:\n%s", gone, s)
		}
	}
	if !strings.Contains(s, "func XResourceSchema") || !strings.HasPrefix(s, "// Code generated") {
		t.Errorf("schema function or generated header lost:\n%s", s)
	}
}

func TestMergeIRNestedOverlay(t *testing.T) {
	nested := attr("vpn", "single_nested", computedOptional, map[string]any{"attributes": []any{
		attr("snat", "bool", computedOptional, map[string]any{"default": map[string]any{"static": false}}),
		attr("name", "string", required, nil),
	}})
	root := attr("mode", "string", computedOptional, map[string]any{"default": map[string]any{"static": "a"}})
	ir := map[string]any{"resources": []any{map[string]any{"name": "site", "schema": map[string]any{"attributes": []any{nested, root}}}}}
	merged, err := mergeIR([]versionedIR{{"v1", ir}}, mergeKinds{})
	if err != nil {
		t.Fatal(err)
	}
	attrs := listOf(merged["resources"].([]any)[0].(map[string]any)["schema"].(map[string]any)["attributes"])
	_, vpn := attrType(findAttr(attrs, "vpn"))
	_, snat := attrType(vpn["attributes"].([]any)[0].(map[string]any))
	if _, ok := snat["default"]; ok {
		t.Error("nested defaults must be removed")
	}
	if _, name := attrType(vpn["attributes"].([]any)[1].(map[string]any)); name["computed_optional_required"] != computedOptional {
		t.Error("children of computed parents must be computed")
	}
	if _, mode := attrType(findAttr(attrs, "mode")); mode["default"] == nil {
		t.Error("root defaults must be kept")
	}
}

func TestMergeIRNoDefaultsForExistingObjects(t *testing.T) {
	ir := map[string]any{"resources": []any{
		map[string]any{"name": "settings", "schema": map[string]any{"attributes": []any{attr("admin", "bool", computedOptional, map[string]any{"default": map[string]any{"static": false}})}}},
		map[string]any{"name": "fresh", "schema": map[string]any{"attributes": []any{attr("admin", "bool", computedOptional, map[string]any{"default": map[string]any{"static": false}})}}},
	}}
	merged, err := mergeIR([]versionedIR{{"v1", ir}}, mergeKinds{existing: map[string]bool{"settings": true}})
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range merged["resources"].([]any) {
		rm := r.(map[string]any)
		_, body := attrType(findAttr(listOf(rm["schema"].(map[string]any)["attributes"]), "admin"))
		if _, has := body["default"]; has != (rm["name"] == "fresh") {
			t.Errorf("%s: default kept=%v", rm["name"], has)
		}
	}
}
