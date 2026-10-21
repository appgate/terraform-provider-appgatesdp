package codegen

import (
	"fmt"
	"slices"
	"strings"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
)

// Definitions is the hand-maintained list of Terraform resources and data
// sources, shared by every API version (codegen/resources.yml).
type Definitions struct {
	Resources map[string]Definition `yaml:"resources"`
	// DataSources lists resource names that also get a data source.
	DataSources []string `yaml:"data_sources"`
}

// Definition locates a Terraform resource in the SDP Admin API.
type Definition struct {
	// Path is the collection path, e.g. /conditions. Objects live at Path/{id}.
	Path string `yaml:"path"`
	// Singleton resources have exactly one object at Path (GET/PUT only).
	Singleton bool `yaml:"singleton"`
	// Variant narrows a polymorphic collection to one discriminator value.
	Variant *apimeta.Discriminator `yaml:"variant"`
	// Ignore lists JSON property paths (dot separated) excluded from the schema.
	Ignore []string `yaml:"ignore"`
	// Adopt marks built-in objects that always exist (e.g. the local identity
	// provider). Create adopts the object with the configured name and delete
	// only removes it from state.
	Adopt bool `yaml:"adopt"`
	// DataSourceOnly definitions only generate a data source.
	DataSourceOnly bool `yaml:"data_source_only"`

	// The fields below describe APIs that are not plain collections.

	// CreateMethod overrides the create method on Path (default POST, PUT for singletons).
	CreateMethod string `yaml:"create_method"`
	// ReadList reads objects from this list property of GET Path instead of GET Path/{id}.
	ReadList string `yaml:"read_list"`
	// IDProperty is the JSON property identifying an object (default id).
	IDProperty string `yaml:"id_property"`
	// DeletePath is the specification path deleting an object, e.g. /blacklist/{distinguished-name}.
	DeletePath string `yaml:"delete_path"`
	// Immutable objects cannot be updated; every change replaces them.
	Immutable bool `yaml:"immutable"`
	// Sensitive lists JSON properties to mark sensitive beyond the generic rules.
	Sensitive []string `yaml:"sensitive"`
	// Required lists request JSON properties the API requires although the
	// specification does not say so.
	Required []string `yaml:"required"`
}

// noObjectPath reports whether objects have no GET Path/{id} endpoint.
func (d Definition) noObjectPath() bool { return d.Singleton || d.ReadList != "" }

func (d Definition) createMethod() string {
	switch {
	case d.CreateMethod != "":
		return strings.ToLower(d.CreateMethod)
	case d.Singleton:
		return "put"
	}
	return "post"
}

// dataSourceNames returns every definition that gets a data source.
func (d Definitions) dataSourceNames() []string {
	names := slices.Clone(d.DataSources)
	for name, def := range d.Resources {
		if def.DataSourceOnly && !slices.Contains(names, name) {
			names = append(names, name)
		}
	}
	slices.Sort(names)
	return names
}

// buildResult is everything derived from one API version's specification.
type buildResult struct {
	// spec is the generator input: one synthetic path per resource.
	spec map[string]any
	// config is the tfplugingen-openapi generator config.
	config map[string]any
	api    *apimeta.API
	// unsupported maps definitions missing from this version to the reason.
	unsupported map[string]string
	warnings    []string
}

// build derives the generator input and API metadata for one version from the
// specification extracted into dir.
func build(version, dir string, defs Definitions) (*buildResult, error) {
	l := newSpecLoader(dir)
	root, err := l.file("api_specs.yml")
	if err != nil {
		return nil, err
	}
	rootMap, _ := root.(map[string]any)
	paths, _ := rootMap["paths"].(map[string]any)
	if len(paths) == 0 {
		return nil, fmt.Errorf("%s: api_specs.yml has no paths", version)
	}
	for _, name := range defs.DataSources {
		if _, ok := defs.Resources[name]; !ok {
			return nil, fmt.Errorf("data source %s has no resource definition", name)
		}
	}
	n := &normalizer{l: l}
	res := &buildResult{
		api:         &apimeta.API{Version: version, Resources: map[string]*apimeta.Resource{}, DataSources: map[string]*apimeta.Resource{}},
		unsupported: map[string]string{},
	}
	specPaths := map[string]any{}
	genResources := map[string]any{}
	genDataSources := map[string]any{}
	schemas := map[string]*resourceSchemas{}

	for _, name := range sortedKeys(defs.Resources) {
		def := defs.Resources[name]
		r, err := n.resource(paths, name, def)
		if err != nil {
			return nil, fmt.Errorf("%s: %s: %w", version, name, err)
		}
		if r.unsupported != "" {
			res.unsupported[name] = r.unsupported
			continue
		}
		schemas[name] = r
		if def.DataSourceOnly {
			continue
		}
		res.api.Resources[name] = r.meta
		base, item := syntheticPaths(name, def)
		create := operation(r.request, r.response, false)
		specPaths[base] = map[string]any{"post": create}
		idParam := !def.noObjectPath()
		itemOps := map[string]any{
			"get":    operation(nil, r.response, idParam),
			"put":    operation(r.request, r.response, idParam),
			"delete": operation(nil, nil, idParam),
		}
		if !idParam {
			itemOps["post"] = create
		}
		specPaths[item] = itemOps
		genResources[name] = withIgnores(map[string]any{
			"create": map[string]any{"path": base, "method": "POST"},
			"read":   map[string]any{"path": item, "method": "GET"},
			"update": map[string]any{"path": item, "method": "PUT"},
			"delete": map[string]any{"path": item, "method": "DELETE"},
		}, def)
	}
	for _, name := range defs.dataSourceNames() {
		def := defs.Resources[name]
		r, ok := schemas[name]
		if !ok {
			continue
		}
		res.api.DataSources[name] = r.meta
		_, item := syntheticPaths(name, def)
		ops, _ := specPaths[item].(map[string]any)
		if ops == nil {
			ops = map[string]any{}
			specPaths[item] = ops
		}
		ops["get"] = operation(nil, r.response, !def.Singleton)
		genDataSources[name] = withIgnores(map[string]any{"read": map[string]any{"path": item, "method": "GET"}}, def)
	}

	res.warnings = n.warnings
	res.spec = map[string]any{
		"openapi": "3.0.0",
		"info":    map[string]any{"title": "SDP Admin API " + version + " (normalized for code generation)", "version": version},
		"paths":   specPaths,
	}
	res.config = map[string]any{
		"provider":     map[string]any{"name": "appgatesdp"},
		"resources":    genResources,
		"data_sources": genDataSources,
	}
	return res, nil
}

// syntheticPaths returns the generator-input paths of a definition. Every
// definition gets its own paths, so variants sharing an API collection (such
// as the policy types) can be mapped independently.
func syntheticPaths(name string, def Definition) (base, item string) {
	base = "/" + name
	if def.noObjectPath() {
		return base, base
	}
	return base, base + "/{id}"
}

func withIgnores(op map[string]any, def Definition) map[string]any {
	if len(def.Ignore) > 0 {
		op["schema"] = map[string]any{"ignores": tfPaths(def.Ignore)}
	}
	return op
}

type resourceSchemas struct {
	request, response map[string]any
	meta              *apimeta.Resource
	// unsupported is set when the version does not provide the definition.
	unsupported string
}

// resource extracts the normalized request/response schemas of one definition.
func (n *normalizer) resource(paths map[string]any, name string, def Definition) (*resourceSchemas, error) {
	coll, ok := paths[def.Path]
	if !ok {
		return &resourceSchemas{unsupported: "path " + def.Path + " not in specification"}, nil
	}
	item := coll
	if !def.noObjectPath() {
		if item, ok = paths[def.Path+"/{id}"]; !ok {
			return &resourceSchemas{unsupported: "path " + def.Path + "/{id} not in specification"}, nil
		}
	}
	deleteItem := item
	if def.DeletePath != "" {
		if deleteItem, ok = paths[def.DeletePath]; !ok {
			return &resourceSchemas{unsupported: "path " + def.DeletePath + " not in specification"}, nil
		}
	}
	collFile, collItem, err := n.l.deref("api_specs.yml", coll)
	if err != nil {
		return nil, err
	}
	itemFile, itemItem, err := n.l.deref("api_specs.yml", item)
	if err != nil {
		return nil, err
	}
	respFile, respNode, err := n.body(itemFile, itemItem, "get", "responses")
	if err != nil {
		return nil, err
	}
	reqFile, reqNode := respFile, respNode
	if !def.DataSourceOnly {
		if reqFile, reqNode, err = n.body(collFile, collItem, def.createMethod(), "requestBody"); err != nil {
			return nil, err
		}
	}
	if def.Variant != nil {
		var found bool
		if reqFile, reqNode, found, err = n.selectVariant(reqFile, reqNode, def.Variant.Value); err != nil {
			return nil, err
		} else if !found {
			return &resourceSchemas{unsupported: fmt.Sprintf("%s=%s is not a variant of %s", def.Variant.Property, def.Variant.Value, def.Path)}, nil
		}
		if respFile, respNode, _, err = n.selectVariant(respFile, respNode, def.Variant.Value); err != nil {
			return nil, err
		}
	}
	reqSchema := n.schema(reqFile, reqNode, request, name+"(request)", nil)
	respSchema := n.schema(respFile, respNode, response, name+"(response)", nil)
	if def.ReadList != "" && respSchema != nil {
		list, _ := respSchema["properties"].(map[string]any)[def.ReadList].(map[string]any)
		respSchema, _ = list["items"].(map[string]any)
	}
	if reqSchema != nil && len(def.Required) > 0 {
		reqSchema["required"] = append(stringList(reqSchema["required"]), def.Required...)
	}
	for _, s := range []map[string]any{reqSchema, respSchema} {
		props, _ := s["properties"].(map[string]any)
		for _, p := range def.Sensitive {
			if ps, ok := props[p].(map[string]any); ok && ps["type"] == "string" {
				ps["format"] = "password"
			}
		}
	}
	if reqSchema == nil || respSchema == nil || reqSchema["type"] != "object" || respSchema["type"] != "object" {
		return nil, fmt.Errorf("request or response of %s is not an object schema", def.Path)
	}
	meta := &apimeta.Resource{
		Path:         def.Path,
		Singleton:    def.Singleton,
		Adopt:        def.Adopt,
		CreateMethod: strings.ToUpper(def.CreateMethod),
		ReadList:     def.ReadList,
		IDProperty:   def.IDProperty,
		DeletePath:   def.DeletePath,
		Immutable:    def.Immutable,
	}
	_, del, err := n.l.deref("api_specs.yml", deleteItem)
	if err != nil {
		return nil, err
	}
	_, meta.Delete = del.(map[string]any)["delete"]
	if def.Variant != nil {
		if prop, ok := reqSchema["properties"].(map[string]any)[def.Variant.Property].(map[string]any); ok {
			if enum := listOf(prop["enum"]); enum != nil && !slices.Contains(enum, any(def.Variant.Value)) {
				return &resourceSchemas{unsupported: fmt.Sprintf("%s=%s is not allowed by %s", def.Variant.Property, def.Variant.Value, def.Path)}, nil
			}
			meta.Discriminator = def.Variant
		}
	}
	for _, s := range []map[string]any{reqSchema, respSchema} {
		if def.Variant != nil {
			removeProperty(s, def.Variant.Property)
		}
		for _, ig := range def.Ignore {
			removeProperty(s, ig)
		}
	}
	if def.DataSourceOnly {
		// Data sources only read; nothing is sent.
		reqSchema = respSchema
	}
	meta.Fields = fields(reqSchema, respSchema)
	return &resourceSchemas{request: reqSchema, response: respSchema, meta: meta}, nil
}

// body returns the JSON schema of an operation's request body or 2xx response.
func (n *normalizer) body(file string, pathItem any, method, kind string) (string, any, error) {
	op, _ := pathItem.(map[string]any)[method].(map[string]any)
	if op == nil {
		return "", nil, fmt.Errorf("operation %s not found", strings.ToUpper(method))
	}
	node := op["requestBody"]
	if kind == "responses" {
		rs, _ := op["responses"].(map[string]any)
		node = nil
		for _, code := range []string{"200", "201"} {
			if r, ok := rs[code]; ok {
				node = r
				break
			}
		}
	}
	if node == nil {
		return "", nil, fmt.Errorf("%s in %s has no %s", strings.ToUpper(method), file, kind)
	}
	f, resolved, err := n.l.deref(file, node)
	if err != nil {
		return "", nil, err
	}
	content, _ := resolved.(map[string]any)["content"].(map[string]any)
	for _, ct := range append([]string{"application/json"}, sortedKeys(content)...) {
		if c, ok := content[ct].(map[string]any); ok && c["schema"] != nil {
			return f, c["schema"], nil
		}
	}
	return "", nil, fmt.Errorf("%s in %s has no JSON %s schema", strings.ToUpper(method), file, kind)
}

// operation builds a synthetic OpenAPI operation for the generator input.
func operation(req, resp map[string]any, idParam bool) map[string]any {
	op := map[string]any{"responses": map[string]any{"204": map[string]any{"description": "No content"}}}
	if resp != nil {
		op["responses"] = map[string]any{"200": map[string]any{
			"description": "OK",
			"content":     map[string]any{"application/json": map[string]any{"schema": resp}},
		}}
	}
	if req != nil {
		op["requestBody"] = map[string]any{
			"required": true,
			"content":  map[string]any{"application/json": map[string]any{"schema": req}},
		}
	}
	if idParam {
		op["parameters"] = []any{map[string]any{
			"name": "id", "in": "path", "required": true,
			"description": "ID of the object.",
			"schema":      map[string]any{"type": "string"},
		}}
	}
	return op
}

// removeProperty deletes a dot separated JSON property path from a schema.
func removeProperty(s map[string]any, path string) {
	head, rest, nested := strings.Cut(path, ".")
	if items, ok := s["items"].(map[string]any); ok && s["type"] == "array" {
		s = items
	}
	props, _ := s["properties"].(map[string]any)
	if !nested {
		delete(props, head)
		if req, ok := s["required"].([]string); ok {
			s["required"] = slices.DeleteFunc(slices.Clone(req), func(r string) bool { return r == head })
		}
		return
	}
	if child, ok := props[head].(map[string]any); ok {
		removeProperty(child, rest)
	}
}

// tfPaths converts dot separated JSON property paths to generator ignores.
func tfPaths(paths []string) []string {
	out := make([]string, len(paths))
	for i, p := range paths {
		parts := strings.Split(p, ".")
		for j := range parts {
			parts[j] = TerraformIdentifier(parts[j])
		}
		out[i] = strings.Join(parts, ".")
	}
	return out
}

// fields builds the field metadata from the request and response schemas.
func fields(req, resp map[string]any) map[string]*apimeta.Field {
	out := map[string]*apimeta.Field{}
	reqProps, _ := req["properties"].(map[string]any)
	respProps, _ := resp["properties"].(map[string]any)
	required := map[string]bool{}
	for _, r := range stringList(req["required"]) {
		required[r] = true
	}
	names := map[string]bool{}
	for k := range reqProps {
		names[k] = true
	}
	for k := range respProps {
		names[k] = true
	}
	for _, name := range sortedKeys(names) {
		rq, _ := reqProps[name].(map[string]any)
		rs, _ := respProps[name].(map[string]any)
		s := rq
		if s == nil {
			s = rs
		}
		f := &apimeta.Field{
			JSON:      name,
			ReadOnly:  rq == nil,
			WriteOnly: rs == nil,
			Required:  required[name],
		}
		f.Type, f.Elem = fieldType(s)
		if f.Type == "string" {
			f.Enum = stringList(s["enum"])
		}
		nestedReq, nestedResp := rq, rs
		if f.Type != "object" {
			nestedReq, nestedResp = elemSchema(rq), elemSchema(rs)
		}
		if (nestedReq != nil && nestedReq["type"] == "object") || (nestedResp != nil && nestedResp["type"] == "object") {
			f.Fields = fields(orEmpty(nestedReq), orEmpty(nestedResp))
		}
		tf := TerraformIdentifier(name)
		if _, dup := out[tf]; dup || tf == "" {
			continue
		}
		out[tf] = f
	}
	return out
}

func elemSchema(s map[string]any) map[string]any {
	if s == nil {
		return nil
	}
	if items, ok := s["items"].(map[string]any); ok {
		return items
	}
	ap, _ := s["additionalProperties"].(map[string]any)
	return ap
}

func orEmpty(s map[string]any) map[string]any {
	if s == nil {
		return map[string]any{}
	}
	return s
}

// fieldType maps a normalized schema to the metadata type names.
func fieldType(s map[string]any) (typ, elem string) {
	switch s["type"] {
	case "string":
		return "string", ""
	case "integer":
		return "int64", ""
	case "number":
		return "number", ""
	case "boolean":
		return "bool", ""
	case "array":
		e, _ := fieldType(elemSchema(s))
		if s["format"] == "set" {
			return "set", e
		}
		return "list", e
	case "object":
		if _, ok := s["properties"]; !ok {
			e, _ := fieldType(elemSchema(s))
			return "map", e
		}
		return "object", ""
	}
	return "", ""
}
