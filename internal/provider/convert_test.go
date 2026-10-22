package provider

import (
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
)

var (
	objType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"name":     tftypes.String,
		"port":     tftypes.Number,
		"weight":   tftypes.Number,
		"enabled":  tftypes.Bool,
		"tags":     tftypes.List{ElementType: tftypes.String},
		"hosts":    tftypes.Set{ElementType: tftypes.String},
		"labels":   tftypes.Map{ElementType: tftypes.String},
		"secret":   tftypes.String,
		"created":  tftypes.String,
		"newer":    tftypes.String,
		"settings": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"mode": tftypes.String, "extra": tftypes.Bool}},
		"rules":    tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"action": tftypes.String}}},
	}}
	testFields = map[string]*apimeta.Field{
		"name":     {JSON: "name", Type: "string", Required: true},
		"port":     {JSON: "port", Type: "int64"},
		"weight":   {JSON: "weight", Type: "number"},
		"enabled":  {JSON: "enabled", Type: "bool"},
		"tags":     {JSON: "tags", Type: "list", Elem: "string"},
		"hosts":    {JSON: "hosts", Type: "set", Elem: "string"},
		"labels":   {JSON: "labels", Type: "map", Elem: "string"},
		"secret":   {JSON: "sharedSecret", Type: "string", WriteOnly: true},
		"created":  {JSON: "created", Type: "string", ReadOnly: true},
		"settings": {JSON: "settings", Type: "object", Fields: map[string]*apimeta.Field{"mode": {JSON: "mode", Type: "string", Required: true}}},
		"rules":    {JSON: "rules", Type: "list", Fields: map[string]*apimeta.Field{"action": {JSON: "action", Type: "string"}}},
		// "newer" is not part of this API version.
	}
)

// testShape marks enabled and labels as optional only; the rest is computed.
var testShape = shape{"enabled": {optionalOnly: true}, "labels": {optionalOnly: true}}

func str(s string) tftypes.Value { return tftypes.NewValue(tftypes.String, s) }

func value(t *testing.T, attrs map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	full := map[string]tftypes.Value{}
	for name, at := range objType.AttributeTypes {
		if v, ok := attrs[name]; ok {
			full[name] = v
		} else {
			full[name] = tftypes.NewValue(at, nil)
		}
	}
	return tftypes.NewValue(objType, full)
}

func strList(elems ...string) tftypes.Value {
	vs := make([]tftypes.Value, len(elems))
	for i, e := range elems {
		vs[i] = str(e)
	}
	return tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, vs)
}

func TestObjectToJSON(t *testing.T) {
	v := value(t, map[string]tftypes.Value{
		"name":    str("n"),
		"port":    tftypes.NewValue(tftypes.Number, big.NewFloat(443)),
		"weight":  tftypes.NewValue(tftypes.Number, big.NewFloat(0.5)),
		"enabled": tftypes.NewValue(tftypes.Bool, false),
		"tags":    strList("a", "b"),
		"labels":  tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{"k": str("v")}),
		"secret":  str("s"),
		"created": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"newer":   str("not sent"),
		"settings": tftypes.NewValue(objType.AttributeTypes["settings"], map[string]tftypes.Value{
			"mode": str("m"), "extra": tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue),
		}),
	})
	got, err := objectToJSON(v, testFields)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"name": "n", "port": int64(443), "weight": json.Number("0.5"), "enabled": false,
		"tags": []any{"a", "b"}, "labels": map[string]any{"k": "v"}, "sharedSecret": "s",
		"settings": map[string]any{"mode": "m"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got  %#v\nwant %#v", got, want)
	}
	if _, err := objectToJSON(value(t, map[string]tftypes.Value{"port": tftypes.NewValue(tftypes.Number, big.NewFloat(1.5))}), testFields); err == nil {
		t.Fatal("fractional value for an integer field must fail")
	}
}

func decode(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatal(err)
	}
	return m
}

func attr(t *testing.T, v tftypes.Value, name string) tftypes.Value {
	t.Helper()
	var attrs map[string]tftypes.Value
	if err := v.As(&attrs); err != nil {
		t.Fatal(err)
	}
	return attrs[name]
}

func TestObjectFromJSONReconciliation(t *testing.T) {
	prior := value(t, map[string]tftypes.Value{
		"name":   str("n"),
		"tags":   strList("Prod", "api"),
		"secret": str("kept"),
		"newer":  tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"port":   tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
	})
	data := decode(t, `{"name":"n","tags":["api","prod"],"enabled":false,"hosts":["a","a","b"],"created":"2026","port":8443,
		"labels":{},"settings":{"mode":"x","unknownToTerraform":1},"rules":[{"action":"allow"}]}`)
	got, err := objectFromJSON(data, objType, testFields, testShape, prior, afterApply)
	if err != nil {
		t.Fatal(err)
	}
	checks := map[string]tftypes.Value{
		// reordered and re-cased by the API: the configured value is kept
		"tags": strList("Prod", "api"),
		// write-only: never returned, prior value kept
		"secret": str("kept"),
		// unknown in plan and unsupported by the version: becomes null
		"newer": tftypes.NewValue(tftypes.String, nil),
		// zero values for attributes that were null stay null
		"enabled": tftypes.NewValue(tftypes.Bool, nil),
		"labels":  tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
		"created": str("2026"),
		"port":    tftypes.NewValue(tftypes.Number, big.NewFloat(8443)),
		"hosts":   tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{str("a"), str("b")}),
	}
	for name, want := range checks {
		if g := attr(t, got, name); !g.Equal(want) {
			t.Errorf("%s = %s, want %s", name, g, want)
		}
	}
	if !got.IsFullyKnown() {
		t.Fatal("state must not contain unknown values")
	}
}

func TestObjectFromJSONRefreshDetectsRemovedValues(t *testing.T) {
	prior := value(t, map[string]tftypes.Value{"name": str("n"), "port": tftypes.NewValue(tftypes.Number, big.NewFloat(1))})
	got, err := objectFromJSON(decode(t, `{"name":"n"}`), objType, testFields, testShape, prior, refresh)
	if err != nil {
		t.Fatal(err)
	}
	if p := attr(t, got, "port"); !p.IsNull() {
		t.Fatalf("refresh must surface removed values as drift, got %s", p)
	}
	got, _ = objectFromJSON(decode(t, `{"name":"n"}`), objType, testFields, testShape, prior, afterApply)
	if p := attr(t, got, "port"); p.IsNull() {
		t.Fatal("after apply, values the API does not echo keep the planned value")
	}
}

func TestObjectFromJSONTypeErrors(t *testing.T) {
	for _, body := range []string{`{"name":["x"]}`, `{"port":"abc"}`, `{"enabled":"yes"}`, `{"settings":[1]}`, `{"tags":"x"}`} {
		if _, err := objectFromJSON(decode(t, body), objType, testFields, testShape, tftypes.NewValue(objType, nil), refresh); err == nil {
			t.Errorf("%s: expected type error", body)
		}
	}
}

func TestValidateForVersion(t *testing.T) {
	config := value(t, map[string]tftypes.Value{
		"newer":   str("x"),
		"created": str("x"),
		"settings": tftypes.NewValue(objType.AttributeTypes["settings"], map[string]tftypes.Value{
			"mode": tftypes.NewValue(tftypes.String, nil), "extra": tftypes.NewValue(tftypes.Bool, true),
		}),
	})
	diags := validateForVersion(config, config, testFields, "v21")
	got := map[string]string{}
	for _, d := range diags {
		got[d.(interface{ Path() path.Path }).Path().String()] = d.Summary()
	}
	want := map[string]string{
		"newer":          "Attribute not supported by API version",
		"created":        "Attribute is read-only in API version",
		"settings.extra": "Attribute not supported by API version",
		"settings.mode":  "Missing attribute required by API version",
		"name":           "Missing attribute required by API version",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got  %v\nwant %v", got, want)
	}
}

func TestOverlayJSONKeepsUnmanagedFields(t *testing.T) {
	current := decode(t, `{"id":"1","name":"old","serverOnly":true,"settings":{"mode":"a","hidden":1},"port":1}`)
	plan := value(t, map[string]tftypes.Value{
		"name": str("new"),
		"port": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
		"settings": tftypes.NewValue(objType.AttributeTypes["settings"], map[string]tftypes.Value{
			"mode": str("b"), "extra": tftypes.NewValue(tftypes.Bool, nil),
		}),
	})
	if err := overlayJSON(current, plan, testFields); err != nil {
		t.Fatal(err)
	}
	want := decode(t, `{"id":"1","name":"new","serverOnly":true,"settings":{"mode":"b","hidden":1},"port":1}`)
	if !reflect.DeepEqual(current, want) {
		t.Fatalf("got  %v\nwant %v", current, want)
	}
}

func TestJSONFields(t *testing.T) {
	typ := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"prefs": tftypes.String}}
	fields := map[string]*apimeta.Field{"prefs": {JSON: "prefs", Type: "json"}}
	plan := tftypes.NewValue(typ, map[string]tftypes.Value{"prefs": str(`{ "b": [1, 2],  "a": "x" }`)})
	body, err := objectToJSON(plan, fields)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"prefs": map[string]any{"a": "x", "b": []any{json.Number("1"), json.Number("2")}}}
	if !reflect.DeepEqual(body, want) {
		t.Fatalf("got %#v", body)
	}
	// The API reorders keys: the configured formatting is kept.
	got, err := objectFromJSON(decode(t, `{"prefs":{"a":"x","b":[1,2]}}`), typ, fields, nil, plan, afterApply)
	if err != nil {
		t.Fatal(err)
	}
	if !attr(t, got, "prefs").Equal(str(`{ "b": [1, 2],  "a": "x" }`)) {
		t.Fatalf("semantically equal JSON must keep the prior value, got %s", attr(t, got, "prefs"))
	}
	// A real change is reported in canonical form.
	got, _ = objectFromJSON(decode(t, `{"prefs":{"a":"y"}}`), typ, fields, nil, plan, refresh)
	if !attr(t, got, "prefs").Equal(str(`{"a":"y"}`)) {
		t.Fatalf("got %s", attr(t, got, "prefs"))
	}
	if _, err := objectToJSON(tftypes.NewValue(typ, map[string]tftypes.Value{"prefs": str("{")}), fields); err == nil {
		t.Fatal("invalid JSON must fail")
	}
}

func TestValidateEnumPerVersion(t *testing.T) {
	typ := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"mode": tftypes.String}}
	fields := map[string]*apimeta.Field{"mode": {JSON: "mode", Type: "string", Enum: []string{"a", "b"}}}
	bad := tftypes.NewValue(typ, map[string]tftypes.Value{"mode": str("c")})
	diags := validateForVersion(bad, bad, fields, "v21")
	if len(diags) != 1 || diags[0].Summary() != "Value not supported by API version" {
		t.Fatalf("got %v", diags)
	}
	ok := tftypes.NewValue(typ, map[string]tftypes.Value{"mode": str("b")})
	if diags := validateForVersion(ok, ok, fields, "v21"); diags.HasError() {
		t.Fatal(diags)
	}
}
