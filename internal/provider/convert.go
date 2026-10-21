package provider

import (
	"encoding/json"
	"fmt"
	"math/big"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
)

// Values move between Terraform and the API as tftypes values on one side and
// generic JSON on the other. The generated per-version metadata provides the
// JSON property name, direction (read-only / write-only) and type of every
// attribute, so one implementation serves every resource and API version.

// typeMatches reports whether a Terraform type can hold the field's JSON
// type. Attributes whose type changed between API versions are only
// supported by the versions that match the merged (newest) schema.
func typeMatches(t tftypes.Type, f *apimeta.Field) bool {
	switch {
	case t.Is(tftypes.String):
		return f.Type == "string"
	case t.Is(tftypes.Number):
		return f.Type == "int64" || f.Type == "number"
	case t.Is(tftypes.Bool):
		return f.Type == "bool"
	case t.Is(tftypes.Object{}):
		return f.Type == "object"
	case t.Is(tftypes.List{}):
		return f.Type == "list" && elemMatches(t.(tftypes.List).ElementType, f)
	case t.Is(tftypes.Set{}):
		return f.Type == "set" && elemMatches(t.(tftypes.Set).ElementType, f)
	case t.Is(tftypes.Map{}):
		return f.Type == "map" && elemMatches(t.(tftypes.Map).ElementType, f)
	}
	return false
}

func elemMatches(t tftypes.Type, f *apimeta.Field) bool {
	return typeMatches(t, elemField(f))
}

// elemField describes the elements of a list, set or map field.
func elemField(f *apimeta.Field) *apimeta.Field {
	if f.Fields != nil {
		return &apimeta.Field{JSON: f.JSON, Type: "object", Fields: f.Fields}
	}
	return &apimeta.Field{JSON: f.JSON, Type: f.Elem}
}

// objectToJSON converts a Terraform object to a request body. Null and
// unknown attributes, read-only fields and fields the API version does not
// know are left out.
func objectToJSON(v tftypes.Value, fields map[string]*apimeta.Field) (map[string]any, error) {
	var attrs map[string]tftypes.Value
	if err := v.As(&attrs); err != nil {
		return nil, err
	}
	out := map[string]any{}
	for _, name := range sortedNames(attrs) {
		av := attrs[name]
		f := fields[name]
		if f == nil || f.ReadOnly || av.IsNull() || !av.IsKnown() || !typeMatches(av.Type(), f) {
			continue
		}
		j, err := valueToJSON(av, f)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		out[f.JSON] = j
	}
	return out, nil
}

func valueToJSON(v tftypes.Value, f *apimeta.Field) (any, error) {
	t := v.Type()
	switch {
	case t.Is(tftypes.String):
		var s string
		err := v.As(&s)
		return s, err
	case t.Is(tftypes.Bool):
		var b bool
		err := v.As(&b)
		return b, err
	case t.Is(tftypes.Number):
		n := new(big.Float)
		if err := v.As(&n); err != nil {
			return nil, err
		}
		if f.Type == "int64" {
			i, acc := n.Int64()
			if acc != big.Exact {
				return nil, fmt.Errorf("%s is not an integer", n.Text('g', -1))
			}
			return i, nil
		}
		return json.Number(n.Text('g', -1)), nil
	case t.Is(tftypes.Object{}):
		return objectToJSON(v, f.Fields)
	case t.Is(tftypes.List{}), t.Is(tftypes.Set{}):
		var elems []tftypes.Value
		if err := v.As(&elems); err != nil {
			return nil, err
		}
		out := make([]any, 0, len(elems))
		for _, e := range elems {
			if e.IsNull() || !e.IsKnown() {
				continue
			}
			j, err := valueToJSON(e, elemField(f))
			if err != nil {
				return nil, err
			}
			out = append(out, j)
		}
		return out, nil
	case t.Is(tftypes.Map{}):
		var elems map[string]tftypes.Value
		if err := v.As(&elems); err != nil {
			return nil, err
		}
		out := make(map[string]any, len(elems))
		for k, e := range elems {
			if e.IsNull() || !e.IsKnown() {
				continue
			}
			j, err := valueToJSON(e, elemField(f))
			if err != nil {
				return nil, err
			}
			out[k] = j
		}
		return out, nil
	}
	return nil, fmt.Errorf("unsupported Terraform type %s", t)
}

// mode controls how API responses are reconciled with the prior value.
type mode int

const (
	// afterApply reconciles a create/update response with the planned value.
	afterApply mode = iota
	// refresh reconciles a read with the prior state; fields the API no
	// longer returns become null so drift is detected.
	refresh
)

// objectFromJSON converts an API object to a Terraform object value.
//
// SDP specific reconciliation rules:
//   - write-only fields (secrets) are never returned, the prior value is kept;
//   - after apply, fields the API does not echo keep the planned value;
//   - an empty API value (empty string/collection, false, 0) for an attribute
//     that was null stays null, the API fills unset fields with zero values;
//   - lists of primitives keep their prior order if the API returns the same
//     elements reordered or re-cased (e.g. tags are lowercased and sorted).
func objectFromJSON(data map[string]any, typ tftypes.Object, fields map[string]*apimeta.Field, sh shape, prior tftypes.Value, m mode) (tftypes.Value, error) {
	priorAttrs := map[string]tftypes.Value{}
	if prior.IsKnown() && !prior.IsNull() {
		if err := prior.As(&priorAttrs); err != nil {
			return tftypes.Value{}, err
		}
	}
	vals := make(map[string]tftypes.Value, len(typ.AttributeTypes))
	for name, at := range typ.AttributeTypes {
		p, hasPrior := priorAttrs[name]
		if !hasPrior {
			p = tftypes.NewValue(at, nil)
		}
		f := fields[name]
		if f == nil || !typeMatches(at, f) || f.WriteOnly {
			vals[name] = knownOrNull(p, at)
			continue
		}
		raw, present := data[f.JSON]
		if !present || raw == nil {
			if m == afterApply {
				vals[name] = knownOrNull(p, at)
			} else {
				vals[name] = tftypes.NewValue(at, nil)
			}
			continue
		}
		v, err := valueFromJSON(raw, at, f, sh.nested(name), p, m)
		if err != nil {
			return tftypes.Value{}, fmt.Errorf("%s: %w", name, err)
		}
		// Computed attributes take whatever the API returns; only purely
		// optional ones must match the (null) configuration.
		if p.IsNull() && sh.optionalOnly(name) && isEmpty(v) {
			v = tftypes.NewValue(at, nil)
		}
		vals[name] = v
	}
	return tftypes.NewValue(typ, vals), nil
}

func valueFromJSON(raw any, t tftypes.Type, f *apimeta.Field, sh shape, prior tftypes.Value, m mode) (tftypes.Value, error) {
	switch {
	case t.Is(tftypes.String):
		switch r := raw.(type) {
		case string:
			return tftypes.NewValue(t, r), nil
		case json.Number, bool:
			return tftypes.NewValue(t, fmt.Sprint(r)), nil
		}
	case t.Is(tftypes.Bool):
		if b, ok := raw.(bool); ok {
			return tftypes.NewValue(t, b), nil
		}
	case t.Is(tftypes.Number):
		var s string
		switch r := raw.(type) {
		case json.Number:
			s = r.String()
		case string:
			s = r
		case float64:
			s = fmt.Sprint(r)
		}
		if n, _, err := big.ParseFloat(s, 10, 512, big.ToNearestEven); err == nil {
			return tftypes.NewValue(t, n), nil
		}
	case t.Is(tftypes.Object{}):
		if obj, ok := raw.(map[string]any); ok {
			return objectFromJSON(obj, t.(tftypes.Object), f.Fields, sh, prior, m)
		}
	case t.Is(tftypes.List{}), t.Is(tftypes.Set{}):
		items, ok := raw.([]any)
		if !ok {
			break
		}
		var et tftypes.Type
		if l, ok := t.(tftypes.List); ok {
			et = l.ElementType
		} else {
			et = t.(tftypes.Set).ElementType
		}
		var priorElems []tftypes.Value
		if prior.IsKnown() && !prior.IsNull() {
			_ = prior.As(&priorElems)
		}
		elems := make([]tftypes.Value, 0, len(items))
		for i, item := range items {
			p := tftypes.NewValue(et, nil)
			if t.Is(tftypes.List{}) && i < len(priorElems) {
				p = priorElems[i]
			}
			e, err := valueFromJSON(item, et, elemField(f), sh, p, m)
			if err != nil {
				return tftypes.Value{}, fmt.Errorf("[%d]: %w", i, err)
			}
			elems = append(elems, e)
		}
		if t.Is(tftypes.List{}) && samePrimitives(priorElems, elems) {
			elems = priorElems
		}
		if t.Is(tftypes.Set{}) {
			elems = dedupe(elems)
		}
		return tftypes.NewValue(t, elems), nil
	case t.Is(tftypes.Map{}):
		obj, ok := raw.(map[string]any)
		if !ok {
			break
		}
		et := t.(tftypes.Map).ElementType
		var priorElems map[string]tftypes.Value
		if prior.IsKnown() && !prior.IsNull() {
			_ = prior.As(&priorElems)
		}
		elems := make(map[string]tftypes.Value, len(obj))
		for k, item := range obj {
			p, ok := priorElems[k]
			if !ok {
				p = tftypes.NewValue(et, nil)
			}
			e, err := valueFromJSON(item, et, elemField(f), sh, p, m)
			if err != nil {
				return tftypes.Value{}, fmt.Errorf("[%q]: %w", k, err)
			}
			elems[k] = e
		}
		return tftypes.NewValue(t, elems), nil
	}
	return tftypes.Value{}, fmt.Errorf("unexpected API value %v (%T) for %s field", raw, raw, f.Type)
}

// knownOrNull returns v if it is fully known, otherwise null. State may not
// contain unknown values after apply.
func knownOrNull(v tftypes.Value, t tftypes.Type) tftypes.Value {
	if v.IsFullyKnown() {
		return v
	}
	return tftypes.NewValue(t, nil)
}

// isEmpty reports whether v is the zero value of its type.
func isEmpty(v tftypes.Value) bool {
	if v.IsNull() {
		return true
	}
	t := v.Type()
	switch {
	case t.Is(tftypes.String):
		var s string
		_ = v.As(&s)
		return s == ""
	case t.Is(tftypes.Bool):
		var b bool
		_ = v.As(&b)
		return !b
	case t.Is(tftypes.Number):
		n := new(big.Float)
		_ = v.As(&n)
		return n.Sign() == 0
	case t.Is(tftypes.List{}), t.Is(tftypes.Set{}):
		var elems []tftypes.Value
		_ = v.As(&elems)
		return len(elems) == 0
	case t.Is(tftypes.Map{}), t.Is(tftypes.Object{}):
		var elems map[string]tftypes.Value
		_ = v.As(&elems)
		for _, e := range elems {
			if !isEmpty(e) {
				return false
			}
		}
		return true
	}
	return false
}

// samePrimitives reports whether two lists of strings hold the same elements,
// ignoring order and letter case.
func samePrimitives(a, b []tftypes.Value) bool {
	if len(a) == 0 || len(a) != len(b) {
		return false
	}
	key := func(vs []tftypes.Value) []string {
		out := make([]string, 0, len(vs))
		for _, v := range vs {
			var s string
			if !v.Type().Is(tftypes.String) || !v.IsKnown() || v.IsNull() || v.As(&s) != nil {
				return nil
			}
			out = append(out, strings.ToLower(s))
		}
		slices.Sort(out)
		return out
	}
	ka, kb := key(a), key(b)
	return ka != nil && kb != nil && slices.Equal(ka, kb)
}

// dedupe removes duplicate set elements; the API may return lists with
// duplicates for fields Terraform models as sets.
func dedupe(elems []tftypes.Value) []tftypes.Value {
	out := elems[:0]
	for _, e := range elems {
		if !slices.ContainsFunc(out, e.Equal) {
			out = append(out, e)
		}
	}
	return out
}

func sortedNames[V any](m map[string]V) []string {
	names := make([]string, 0, len(m))
	for k := range m {
		names = append(names, k)
	}
	slices.Sort(names)
	return names
}
