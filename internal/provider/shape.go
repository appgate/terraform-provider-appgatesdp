package provider

import (
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// shape records schema properties that tftypes values do not carry. A nil
// shape treats every attribute as computed (data sources).
type shape map[string]*attrShape

type attrShape struct {
	// optionalOnly is set for Optional attributes that are not Computed: the
	// only ones whose state must stay null when the configuration is null.
	optionalOnly bool
	nested       shape
}

func (s shape) optionalOnly(name string) bool {
	a, ok := s[name]
	return ok && a.optionalOnly
}

func (s shape) nested(name string) shape {
	if a, ok := s[name]; ok {
		return a.nested
	}
	return nil
}

// resourceShape builds the shape of a resource schema.
func resourceShape(attrs map[string]rschema.Attribute) shape {
	out := shape{}
	for name, a := range attrs {
		s := &attrShape{optionalOnly: a.IsOptional() && !a.IsComputed()}
		switch n := a.(type) {
		case rschema.SingleNestedAttribute:
			s.nested = resourceShape(n.Attributes)
		case rschema.ListNestedAttribute:
			s.nested = resourceShape(n.NestedObject.Attributes)
		case rschema.SetNestedAttribute:
			s.nested = resourceShape(n.NestedObject.Attributes)
		case rschema.MapNestedAttribute:
			s.nested = resourceShape(n.NestedObject.Attributes)
		}
		out[name] = s
	}
	return out
}
