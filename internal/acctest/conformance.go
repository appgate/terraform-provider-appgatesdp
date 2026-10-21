package acctest

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
)

// Finding is a difference between a live API response and the generated
// metadata of its API version.
type Finding struct {
	// Severity is "error" for values the provider cannot represent and
	// "info" for fields the specification does not describe.
	Severity string
	Path     string
	Message  string
}

func (f Finding) String() string { return fmt.Sprintf("%s %s: %s", f.Severity, f.Path, f.Message) }

// Conformance compares an API object with the field metadata.
//   - error: a field has a different JSON type than specified, or a value
//     is not in the specified enum;
//   - info: the API returned a field the specification does not describe
//     (the provider ignores it, but the spec may be incomplete).
func Conformance(path string, obj map[string]any, fields map[string]*apimeta.Field) []Finding {
	var out []Finding
	byJSON := map[string]*apimeta.Field{}
	for _, f := range fields {
		byJSON[f.JSON] = f
	}
	for _, k := range sortedKeys(obj) {
		p := path + "." + k
		f, ok := byJSON[k]
		if !ok {
			out = append(out, Finding{"info", p, "field not in specification"})
			continue
		}
		out = append(out, conformValue(p, obj[k], f.Type, f)...)
	}
	return out
}

func conformValue(p string, v any, typ string, f *apimeta.Field) []Finding {
	if v == nil {
		return nil
	}
	mismatch := func() []Finding {
		return []Finding{{"error", p, fmt.Sprintf("expected %s, got %T %v", typ, v, v)}}
	}
	switch typ {
	case "string":
		s, ok := v.(string)
		if !ok {
			return mismatch()
		}
		if len(f.Enum) > 0 && !slices.Contains(f.Enum, s) {
			return []Finding{{"error", p, fmt.Sprintf("value %q not in enum %v", s, f.Enum)}}
		}
	case "bool":
		if _, ok := v.(bool); !ok {
			return mismatch()
		}
	case "int64":
		n, ok := v.(json.Number)
		if !ok {
			return mismatch()
		}
		if _, err := n.Int64(); err != nil {
			return []Finding{{"error", p, fmt.Sprintf("expected integer, got %s", n)}}
		}
	case "number":
		if _, ok := v.(json.Number); !ok {
			return mismatch()
		}
	case "object":
		m, ok := v.(map[string]any)
		if !ok {
			return mismatch()
		}
		return Conformance(p, m, f.Fields)
	case "list", "set":
		items, ok := v.([]any)
		if !ok {
			return mismatch()
		}
		var out []Finding
		for i, item := range items {
			out = append(out, conformElem(fmt.Sprintf("%s[%d]", p, i), item, f)...)
		}
		return out
	case "map":
		m, ok := v.(map[string]any)
		if !ok {
			return mismatch()
		}
		var out []Finding
		for _, k := range sortedKeys(m) {
			out = append(out, conformElem(p+"."+k, m[k], f)...)
		}
		return out
	}
	return nil
}

func conformElem(p string, v any, f *apimeta.Field) []Finding {
	if f.Fields != nil {
		m, ok := v.(map[string]any)
		if !ok {
			return []Finding{{"error", p, fmt.Sprintf("expected object, got %T", v)}}
		}
		return Conformance(p, m, f.Fields)
	}
	return conformValue(p, v, f.Elem, &apimeta.Field{Type: f.Elem})
}

// Summarize renders findings for test logs.
func Summarize(findings []Finding) string {
	var b strings.Builder
	for _, f := range findings {
		b.WriteString(f.String())
		b.WriteByte('\n')
	}
	return b.String()
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
