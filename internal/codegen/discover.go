package codegen

import (
	"fmt"
	"slices"
	"strings"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
)

// Every path of a version's specification is classified, so nothing an API
// version offers is silently left out:
//
//   - collection: POST <path> plus GET <path>/{x}. One resource per
//     discriminator variant of the request body, otherwise one resource.
//   - singleton: GET and PUT on <path> (e.g. /global-settings).
//   - registration: GET, POST and DELETE on <path> without object paths
//     (e.g. /replication-source): an immutable singleton created with POST.
//   - read-only: any other GET returning JSON becomes a data source whose
//     path and query parameters are attributes.
//   - action: remaining operations (bulk upserts, reboots, tests, ...).
//
// Definitions in codegen/resources.yml take precedence (they keep names and
// tune behavior); codegen/resources.yml `exclude` skips paths with a reason.

// pathClass is the classification of one specification path.
type pathClass struct {
	Kind   string // resource, data source, action, excluded
	Detail string
}

// readOnlyDef is a discovered read-only data source.
type readOnlyDef struct {
	path string
}

// discovery is the effective set of definitions of one API version.
type discovery struct {
	resources map[string]Definition
	readOnly  map[string]readOnlyDef
	classes   map[string]pathClass
}

func (n *normalizer) discover(paths map[string]any, defs Definitions) (*discovery, error) {
	d := &discovery{resources: map[string]Definition{}, readOnly: map[string]readOnlyDef{}, classes: map[string]pathClass{}}
	ops := map[string]map[string]bool{}
	files := map[string]string{}
	items := map[string]map[string]any{}
	for _, p := range sortedKeys(paths) {
		f, item, err := n.l.deref("api_specs.yml", paths[p])
		if err != nil {
			return nil, fmt.Errorf("path %s: %w", p, err)
		}
		m, _ := item.(map[string]any)
		ops[p] = map[string]bool{}
		for _, op := range []string{"get", "post", "put", "delete"} {
			if _, ok := m[op]; ok {
				ops[p][op] = true
			}
		}
		files[p], items[p] = f, m
	}

	// claims records which resources and data sources use a path.
	claims := map[string][]string{}
	claim := func(p, kind, name string) {
		entry := kind + " " + TypePrefix + name
		if !slices.Contains(claims[p], entry) {
			claims[p] = append(claims[p], entry)
		}
	}
	// Explicit definitions take precedence over discovery.
	explicitVariants := map[string][]string{}
	for _, name := range sortedKeys(defs.Resources) {
		def := defs.Resources[name]
		d.resources[name] = def // missing paths are reported as unsupported by build
		kind := "resource"
		if def.DataSourceOnly {
			kind = "data source"
		}
		claim(def.Path, kind, name)
		if !def.noObjectPath() {
			claim(def.Path+"/{id}", kind, name)
		}
		if def.DeletePath != "" {
			claim(def.DeletePath, kind, name)
		}
		switch {
		case def.Variant != nil:
			explicitVariants[def.Path] = append(explicitVariants[def.Path], def.Variant.Value)
		case !def.DataSourceOnly:
			explicitVariants[def.Path] = append(explicitVariants[def.Path], "*")
		}
	}

	add := func(name string, def Definition) error {
		if prev, ok := d.resources[name]; ok {
			return fmt.Errorf("discovered resource %s for %s collides with the definition for %s; add a name to codegen/resources.yml", name, def.Path, prev.Path)
		}
		d.resources[name] = def
		claim(def.Path, "resource", name)
		if def.itemPath != "" {
			claim(def.itemPath, "resource", name)
		}
		return nil
	}

	for _, p := range sortedKeys(paths) {
		if reason, ok := defs.Exclude[p]; ok {
			d.classes[p] = pathClass{"excluded", reason}
			continue
		}
		o := ops[p]
		if strings.Contains(p, "{") {
			continue // object and parent-scoped paths are classified below
		}
		item := objectPath(p, ops)
		switch {
		case o["post"] && item != "" && ops[item]["get"]:
			if slices.Contains(explicitVariants[p], "*") {
				continue
			}
			if !jsonRequest(items[p], "post") {
				d.classes[p] = pathClass{"excluded", "request body is not JSON"}
				d.classes[item] = d.classes[p]
				continue
			}
			values, err := n.variants(files[p], items[p])
			if err != nil {
				return nil, fmt.Errorf("path %s: %w", p, err)
			}
			if len(values) == 0 {
				if len(explicitVariants[p]) > 0 {
					continue // explicit variants cover versions without a discriminator
				}
				if err := add(defs.name(p, "", false), Definition{Path: p, discovered: true, itemPath: item}); err != nil {
					return nil, err
				}
			}
			for _, v := range values {
				if slices.Contains(explicitVariants[p], v.Value) {
					continue
				}
				if err := add(defs.name(p, v.Value, false), Definition{Path: p, Variant: v, discovered: true, itemPath: item}); err != nil {
					return nil, err
				}
			}
		case o["get"] && o["put"] && !o["post"] && item == "":
			if len(claims[p]) > 0 {
				continue
			}
			if err := add(defs.name(p, "", true), Definition{Path: p, Singleton: true, discovered: true}); err != nil {
				return nil, err
			}
		case o["get"] && o["post"] && o["delete"] && !o["put"] && item == "":
			if len(claims[p]) > 0 {
				continue
			}
			if err := add(defs.name(p, "", true), Definition{Path: p, Singleton: true, CreateMethod: "POST", Immutable: true, discovered: true}); err != nil {
				return nil, err
			}
		}
	}

	for _, p := range sortedKeys(paths) {
		if _, ok := d.classes[p]; ok {
			continue
		}
		if c := claims[p]; len(c) > 0 {
			kind := "data source"
			if slices.ContainsFunc(c, func(e string) bool { return strings.HasPrefix(e, "resource ") }) {
				kind = "resource"
			}
			d.classes[p] = pathClass{kind, strings.Join(c, ", ")}
			continue
		}
		if ops[p]["get"] {
			if !jsonResponse(items[p]) {
				d.classes[p] = pathClass{"excluded", "response is not JSON"}
				continue
			}
			name := readOnlyName(p)
			if n, ok := defs.Names[p]; ok {
				name = n
			}
			if _, dup := d.resources[name]; dup {
				name += "_info"
			}
			if _, dup := d.readOnly[name]; dup {
				return nil, fmt.Errorf("read-only data source name %s is used twice; add a name for %s to codegen/resources.yml", name, p)
			}
			d.readOnly[name] = readOnlyDef{path: p}
			detail := "data source " + TypePrefix + name
			if ops[p]["post"] || ops[p]["put"] || ops[p]["delete"] {
				detail += "; other methods are actions"
			}
			d.classes[p] = pathClass{"data source", detail}
			continue
		}
		d.classes[p] = pathClass{"action", "not representable as Terraform state"}
	}
	return d, nil
}

// objectPath returns the object path p/{param} of a collection, if any.
func objectPath(p string, ops map[string]map[string]bool) string {
	for q := range ops {
		rest, ok := strings.CutPrefix(q, p+"/{")
		if ok && strings.HasSuffix(rest, "}") && !strings.Contains(rest, "/") {
			return q
		}
	}
	return ""
}

// variants returns the discriminator values of a collection's request body.
func (n *normalizer) variants(file string, item map[string]any) ([]*apimeta.Discriminator, error) {
	f, node, err := n.body(file, item, "post", "requestBody")
	if err != nil {
		return nil, nil //nolint:nilerr // collections without a JSON body are reported as excluded
	}
	_, target, err := n.l.deref(f, node)
	if err != nil {
		return nil, err
	}
	m, _ := target.(map[string]any)
	mapping := discriminatorMapping(m)
	disc, _ := m["discriminator"].(map[string]any)
	prop, _ := disc["propertyName"].(string)
	var out []*apimeta.Discriminator
	for _, v := range sortedKeys(mapping) {
		out = append(out, &apimeta.Discriminator{Property: prop, Value: v})
	}
	return out, nil
}

func jsonRequest(item map[string]any, method string) bool {
	op, _ := item[method].(map[string]any)
	rb, ok := op["requestBody"].(map[string]any)
	if !ok {
		return false
	}
	if _, isRef := rb["$ref"]; isRef {
		return true // shared request bodies in the SDP spec are JSON
	}
	content, _ := rb["content"].(map[string]any)
	_, ok = content["application/json"]
	return ok
}

func jsonResponse(item map[string]any) bool {
	op, _ := item["get"].(map[string]any)
	rs, _ := op["responses"].(map[string]any)
	r, ok := rs["200"].(map[string]any)
	if !ok {
		return false
	}
	if _, isRef := r["$ref"]; isRef {
		return true
	}
	content, _ := r["content"].(map[string]any)
	_, ok = content["application/json"]
	return ok
}

// name returns the Terraform name of a discovered resource: the singular of
// the collection's last path segment, prefixed by the variant if any
// (/secrets + SshKey → ssh_key_secret), or the joined path of singletons.
func (d Definitions) name(path, variant string, singleton bool) string {
	key := path
	if variant != "" {
		key += "#" + variant
	}
	if n, ok := d.Names[key]; ok {
		return n
	}
	segs := strings.Split(strings.Trim(path, "/"), "/")
	base := TerraformIdentifier(strings.ReplaceAll(segs[len(segs)-1], "-", "_"))
	if variant == "" && (singleton || !strings.HasSuffix(base, "s")) {
		return TerraformIdentifier(strings.ReplaceAll(strings.Join(segs, "_"), "-", "_"))
	}
	base = singular(base)
	if variant != "" {
		return TerraformIdentifier(variant) + "_" + base
	}
	return base
}

// readOnlyName names a read-only data source after its literal path
// segments; a trailing parameter makes the last segment singular
// (/discovered-apps/{id} → discovered_app).
func readOnlyName(path string) string {
	segs := strings.Split(strings.Trim(path, "/"), "/")
	var lit []string
	for i, s := range segs {
		if strings.HasPrefix(s, "{") {
			if i == len(segs)-1 && len(lit) > 0 {
				lit[len(lit)-1] = singular(lit[len(lit)-1])
			}
			continue
		}
		lit = append(lit, TerraformIdentifier(strings.ReplaceAll(s, "-", "_")))
	}
	return strings.Join(lit, "_")
}

func singular(s string) string {
	switch {
	case strings.HasSuffix(s, "ies"):
		return strings.TrimSuffix(s, "ies") + "y"
	case strings.HasSuffix(s, "sses"), strings.HasSuffix(s, "settings"):
		return s
	case strings.HasSuffix(s, "s") && !strings.HasSuffix(s, "ss"):
		return strings.TrimSuffix(s, "s")
	}
	return s
}

// TypePrefix is the Terraform type prefix of every resource and data source.
const TypePrefix = "appgatesdp_"
