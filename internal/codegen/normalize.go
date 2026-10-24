package codegen

import (
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"
)

// direction selects which view of a schema is produced. Request schemas drop
// readOnly properties, response schemas drop writeOnly properties.
type direction int

const (
	request direction = iota
	response
)

// normalizer converts SDP schemas into self-contained schemas that
// tfplugingen-openapi can map. It resolves every $ref, flattens allOf, merges
// oneOf/anyOf/discriminator variants into a single object and drops
// constructs the generator cannot express. Every lossy decision is recorded in
// warnings so it shows up in the generation report.
type normalizer struct {
	l           *specLoader
	corrections map[string]Correction
	warnings    []string
}

// keptKeys are the schema keywords forwarded to the generator.
var keptKeys = []string{
	"type", "format", "description", "enum", "default", "minimum", "maximum",
	"minLength", "maxLength", "pattern", "minItems", "maxItems", "readOnly",
	"writeOnly", "nullable",
}

func (n *normalizer) warnf(format string, args ...any) {
	n.warnings = append(n.warnings, fmt.Sprintf(format, args...))
}

// flat is a schema with references resolved and compositions flattened, but
// not yet validated. allOf members are usually fragments (e.g. only a
// description), so they are merged as flat schemas and validated once.
type flat struct {
	keys     map[string]any
	props    map[string]any
	required []string
}

// schema returns the normalized schema or nil if it cannot be represented.
func (n *normalizer) schema(file string, node any, dir direction, at string, stack []string) map[string]any {
	f := n.collect(file, node, dir, at, stack)
	if f == nil {
		return nil
	}
	return n.finalize(f, dir, at)
}

func (n *normalizer) collect(file string, node any, dir direction, at string, stack []string) *flat {
	m, ok := node.(map[string]any)
	if !ok {
		return nil
	}
	if ref, ok := m["$ref"].(string); ok {
		f, target, err := n.l.resolve(file, ref)
		if err != nil {
			n.warnf("%s: %v", at, err)
			return nil
		}
		key := f + "#" + ref[strings.Index(ref, "#")+1:]
		if slices.Contains(stack, key) {
			n.warnf("%s: dropped recursive reference to %s", at, ref)
			return nil
		}
		return n.collect(f, target, dir, at, append(stack, key))
	}

	out := map[string]any{}
	props := map[string]any{}
	var required []string

	merge := func(s map[string]any) {
		for _, k := range keptKeys {
			if v, ok := s[k]; ok {
				out[k] = v
			}
		}
		if p, ok := s["properties"].(map[string]any); ok {
			for k, v := range p {
				props[k] = v
			}
		}
		for _, k := range []string{"items", "additionalProperties"} {
			if v, ok := s[k]; ok {
				out[k] = v
			}
		}
		required = append(required, stringList(s["required"])...)
	}

	for i, sub := range listOf(m["allOf"]) {
		if f := n.collect(file, sub, dir, fmt.Sprintf("%s/allOf[%d]", at, i), stack); f != nil {
			merge(f.keys)
			for k, v := range f.props {
				props[k] = v
			}
			required = append(required, f.required...)
		}
	}

	variants := listOf(m["oneOf"])
	variants = append(variants, listOf(m["anyOf"])...)
	if len(variants) == 0 {
		// A bare discriminator (no oneOf) lists its variants in mapping.
		for _, ref := range sortedVariantRefs(discriminatorMapping(m)) {
			variants = append(variants, map[string]any{"$ref": ref})
		}
	}
	if len(variants) > 0 {
		n.mergeVariants(file, variants, dir, at, stack, merge, &required)
	}

	own := map[string]any{}
	for _, k := range keptKeys {
		if v, ok := m[k]; ok {
			own[k] = v
		}
	}
	own["required"] = m["required"]
	if raw, ok := m["properties"].(map[string]any); ok {
		p := map[string]any{}
		for _, name := range sortedKeys(raw) {
			if child := n.schema(file, raw[name], dir, at+"/"+name, stack); child != nil {
				markSecret(name, child)
				p[name] = child
			}
		}
		own["properties"] = p
	}
	if raw, ok := m["items"]; ok {
		if items := n.schema(file, raw, dir, at+"/items", stack); items != nil {
			own["items"] = items
		}
	}
	if raw, ok := m["additionalProperties"].(map[string]any); ok {
		if ap := n.schema(file, raw, dir, at+"/additionalProperties", stack); ap != nil {
			own["additionalProperties"] = ap
		}
	}
	merge(own)
	if m["uniqueItems"] == true {
		out["uniqueItems"] = true
	}
	return &flat{keys: out, props: props, required: required}
}

// finalize infers the type of a flat schema and drops what the generator
// cannot represent.
func (n *normalizer) finalize(f *flat, dir direction, at string) map[string]any {
	out, props, required := clone(f.keys), f.props, f.required
	resolvedProps := map[string]any{}
	for name, v := range props {
		child := v.(map[string]any)
		if (dir == request && child["readOnly"] == true) || (dir == response && child["writeOnly"] == true) {
			continue
		}
		resolvedProps[name] = child
	}

	typ, _ := out["type"].(string)
	if typ == "" {
		switch {
		case len(resolvedProps) > 0, out["additionalProperties"] != nil:
			typ = "object"
		case out["items"] != nil:
			typ = "array"
		case out["enum"] != nil:
			typ = "string"
		}
	}

	switch typ {
	case "object":
		if _, ok := out["additionalProperties"].(map[string]any); ok && len(resolvedProps) == 0 {
			delete(out, "properties")
		} else {
			delete(out, "additionalProperties")
			if len(resolvedProps) == 0 {
				if len(props) == 0 {
					return freeForm(out)
				}
				// Every property is read-only (request) or write-only (response).
				n.warnf("%s: dropped object without properties in this direction", at)
				return nil
			}
			out["properties"] = resolvedProps
		}
		var req []string
		for _, r := range required {
			if _, ok := resolvedProps[r]; ok && !slices.Contains(req, r) {
				req = append(req, r)
			}
		}
		sort.Strings(req)
		if len(req) > 0 {
			out["required"] = req
		} else {
			delete(out, "required")
		}
	case "array":
		items, ok := out["items"].(map[string]any)
		if !ok {
			n.warnf("%s: dropped array with unsupported items", at)
			return nil
		}
		if items["type"] != "object" && items["type"] != "array" && out["uniqueItems"] == true {
			out["format"] = "set"
		}
		delete(out, "additionalProperties")
	case "string", "integer", "number", "boolean":
		delete(out, "items")
		delete(out, "additionalProperties")
	default:
		// No type, properties or items: any JSON value.
		return freeForm(out)
	}
	out["type"] = typ
	delete(out, "uniqueItems")

	// SDP patterns are ECMA regular expressions; the generated validators
	// use regexp.MustCompile, so anything RE2 rejects must not reach them.
	if p, ok := out["pattern"].(string); ok {
		if _, err := regexp.Compile(p); err != nil {
			n.warnf("%s: dropped pattern %q unsupported by Go regexp", at, p)
			delete(out, "pattern")
		}
	}
	return out
}

var secretName = regexp.MustCompile(`(?i)(password|secret|passphrase)|(token|privatekey|apikey)$`)

// markSecret marks secrets as sensitive. tfplugingen-openapi only derives
// sensitivity from `format: password`, while the SDP specification marks
// many secrets only as writeOnly (e.g. LDAP adminPassword) or not at all.
func markSecret(name string, s map[string]any) {
	if s["type"] != "string" || (s["format"] != nil && s["format"] != "password") || s["format"] == "json" {
		return
	}
	if s["writeOnly"] == true || secretName.MatchString(name) {
		s["format"] = "password"
	}
}

// freeForm represents schemas without declared structure (free-form objects,
// untyped values) as JSON encoded strings, Terraform's convention for
// arbitrary JSON (jsonencode/jsondecode). The runtime compares them
// semantically.
func freeForm(s map[string]any) map[string]any {
	out := map[string]any{"type": "string", "format": "json"}
	desc, _ := s["description"].(string)
	out["description"] = strings.TrimSpace(desc + " JSON encoded value; use jsonencode() and jsondecode().")
	for _, k := range []string{"readOnly", "writeOnly"} {
		if v, ok := s[k]; ok {
			out[k] = v
		}
	}
	return out
}

// mergeVariants flattens polymorphic schemas into one object holding the
// union of all variant properties. Only properties required by every variant
// stay required.
func (n *normalizer) mergeVariants(file string, variants []any, dir direction, at string, stack []string, merge func(map[string]any), required *[]string) {
	var objs []map[string]any
	for i, v := range variants {
		s := n.schema(file, v, dir, fmt.Sprintf("%s/variant[%d]", at, i), stack)
		if s == nil {
			continue
		}
		if s["type"] != "object" {
			n.warnf("%s: polymorphic schema mixes non-object variants, using the first variant", at)
			merge(s)
			return
		}
		objs = append(objs, s)
	}
	var common []string
	for i, s := range objs {
		req := stringList(s["required"])
		if i == 0 {
			common = req
		} else {
			common = slices.DeleteFunc(common, func(r string) bool { return !slices.Contains(req, r) })
		}
		s = clone(s)
		delete(s, "required")
		merge(s)
	}
	*required = append(*required, common...)
}

// selectVariant narrows a discriminated schema to the variant with the given
// discriminator value. Schemas without a discriminator mapping are returned
// unchanged; found is false if a mapping exists but lacks the value.
func (n *normalizer) selectVariant(file string, node any, value string) (string, any, bool, error) {
	f, target, err := n.l.deref(file, node)
	if err != nil {
		return "", nil, false, err
	}
	m, _ := target.(map[string]any)
	mapping := discriminatorMapping(m)
	if len(mapping) == 0 {
		return f, target, true, nil
	}
	ref, ok := mapping[value]
	return f, map[string]any{"$ref": ref}, ok, nil
}

func discriminatorMapping(m map[string]any) map[string]string {
	d, _ := m["discriminator"].(map[string]any)
	mapping, _ := d["mapping"].(map[string]any)
	out := map[string]string{}
	for k, v := range mapping {
		if s, ok := v.(string); ok {
			out[k] = s
		}
	}
	return out
}

// sortedVariantRefs returns mapping refs in a deterministic order.
func sortedVariantRefs(m map[string]string) []string {
	var refs []string
	for _, k := range sortedKeys(m) {
		refs = append(refs, m[k])
	}
	return refs
}

func listOf(v any) []any {
	l, _ := v.([]any)
	return l
}

func stringList(v any) []string {
	if s, ok := v.([]string); ok {
		return s
	}
	var out []string
	for _, x := range listOf(v) {
		if s, ok := x.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func clone(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
