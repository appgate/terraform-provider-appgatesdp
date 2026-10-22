// Package apimeta describes, per SDP Admin API version, how Terraform
// attributes map onto API endpoints and JSON fields. The descriptions are
// produced by the code generator from each version's OpenAPI specification and
// consumed at runtime by the version-aware resource implementations.
package apimeta

import "strings"

// API is the generated description of one Admin API version.
type API struct {
	Version     string               `json:"version"`
	Resources   map[string]*Resource `json:"resources"`
	DataSources map[string]*Resource `json:"data_sources"`
}

// Resource describes the endpoints and fields of a Terraform resource or data
// source for one API version.
type Resource struct {
	// Path is the collection path (e.g. /conditions) or, for singletons, the
	// object path (e.g. /global-settings).
	Path      string `json:"path"`
	Singleton bool   `json:"singleton,omitempty"`
	// Delete reports whether the API supports DELETE on the object.
	Delete bool `json:"delete,omitempty"`
	// Adopt marks built-in objects: create adopts the existing object with
	// the configured name and delete only removes it from state.
	Adopt bool `json:"adopt,omitempty"`
	// CreateMethod is the HTTP method creating the object on Path; empty means
	// POST for collections and PUT for singletons.
	CreateMethod string `json:"create_method,omitempty"`
	// ReadList names the list property of GET Path holding the objects, for
	// APIs without GET Path/{id} (e.g. licenses).
	ReadList string `json:"read_list,omitempty"`
	// IDProperty is the JSON property identifying objects; empty means "id".
	IDProperty string `json:"id_property,omitempty"`
	// DeletePath overrides the object path used by DELETE; its {...}
	// parameter is replaced by the object id.
	DeletePath string `json:"delete_path,omitempty"`
	// Immutable objects have no update; every change replaces them.
	Immutable bool `json:"immutable,omitempty"`
	// Wrapped resources have a free-form JSON body, exposed as the single
	// attribute WrapProperty.
	Wrapped bool `json:"wrapped,omitempty"`
	// ReadOnly data sources read one GET endpoint (Path, with Params filled
	// in) instead of looking objects up by id or name.
	ReadOnly bool `json:"read_only,omitempty"`
	// Params are the path and query parameters of read-only data sources.
	Params []Param `json:"params,omitempty"`
	// Discriminator is a fixed JSON property injected into requests and used
	// to filter responses, e.g. type=Access for access policies.
	Discriminator *Discriminator    `json:"discriminator,omitempty"`
	Fields        map[string]*Field `json:"fields"`
}

// WrapProperty is the attribute holding the body of Wrapped resources.
const WrapProperty = "value"

// Wrap returns the body of an API object in the shape of the field metadata.
func (r *Resource) Wrap(body map[string]any) map[string]any {
	if !r.Wrapped {
		return body
	}
	return map[string]any{WrapProperty: body}
}

// Unwrap returns the request body for a Terraform-shaped object.
func (r *Resource) Unwrap(obj map[string]any) map[string]any {
	if !r.Wrapped {
		return obj
	}
	body, _ := obj[WrapProperty].(map[string]any)
	if body == nil {
		body = map[string]any{}
	}
	return body
}

// Param is a path or query parameter set from a Terraform attribute.
type Param struct {
	Name string `json:"name"`
	// In is "path" or "query".
	In       string   `json:"in"`
	TF       string   `json:"tf"`
	Required bool     `json:"required,omitempty"`
	Enum     []string `json:"enum,omitempty"`
}

// Discriminator is a fixed JSON property value identifying a variant.
type Discriminator struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

// Field maps one Terraform attribute onto a JSON property.
type Field struct {
	JSON string `json:"json"`
	// Type is one of string, int64, number, bool, object, list, set, map, or
	// json for free-form values held as JSON encoded strings.
	Type string `json:"type"`
	// Elem is the element type of list, set and map fields.
	Elem      string `json:"elem,omitempty"`
	ReadOnly  bool   `json:"read_only,omitempty"`
	WriteOnly bool   `json:"write_only,omitempty"`
	Required  bool   `json:"required,omitempty"`
	// Enum lists the allowed values of string fields.
	Enum []string `json:"enum,omitempty"`
	// Fields holds nested attributes of objects and of collections of objects.
	Fields map[string]*Field `json:"fields,omitempty"`
}

// ID returns the JSON property identifying objects.
func (r *Resource) ID() string {
	if r.IDProperty != "" {
		return r.IDProperty
	}
	return "id"
}

// ObjectPath returns the path of the object with the given escaped id.
func (r *Resource) ObjectPath(escapedID string) string {
	if r.Singleton {
		return r.Path
	}
	return r.Path + "/" + escapedID
}

// DeleteObjectPath returns the path deleting the object with the given
// escaped id.
func (r *Resource) DeleteObjectPath(escapedID string) string {
	start, end := strings.Index(r.DeletePath, "{"), strings.Index(r.DeletePath, "}")
	if r.DeletePath == "" || start < 0 || end < start {
		return r.ObjectPath(escapedID)
	}
	return r.DeletePath[:start] + escapedID + r.DeletePath[end+1:]
}

// ListProperty is the property of GET Path holding the objects.
func (r *Resource) ListProperty() string {
	if r.ReadList != "" {
		return r.ReadList
	}
	return "data"
}
