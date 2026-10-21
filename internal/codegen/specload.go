package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// specLoader resolves $ref values across the multi-file SDP specification.
// The SDP spec is split into one YAML file per domain (condition.yml,
// policy.yml, ...) and uses non-standard locations such as
// `./common.yml#/definitions/BaseEntity`, which tfplugingen-openapi cannot
// resolve on its own.
type specLoader struct {
	dir   string
	files map[string]any
}

func newSpecLoader(dir string) *specLoader {
	return &specLoader{dir: dir, files: map[string]any{}}
}

func (l *specLoader) file(name string) (any, error) {
	if doc, ok := l.files[name]; ok {
		return doc, nil
	}
	b, err := os.ReadFile(filepath.Join(l.dir, name))
	if err != nil {
		return nil, fmt.Errorf("read spec file %s: %w", name, err)
	}
	var root yaml.Node
	if err := yaml.Unmarshal(b, &root); err != nil {
		return nil, fmt.Errorf("parse spec file %s: %w", name, err)
	}
	doc := nodeValue(&root)
	l.files[name] = doc
	return doc, nil
}

// nodeValue converts a YAML node to plain Go values. Unlike yaml.Unmarshal into
// `any` it tolerates duplicate mapping keys (last one wins), which occur in
// some SDP spec branches.
func nodeValue(n *yaml.Node) any {
	switch n.Kind {
	case yaml.DocumentNode:
		if len(n.Content) == 0 {
			return nil
		}
		return nodeValue(n.Content[0])
	case yaml.AliasNode:
		return nodeValue(n.Alias)
	case yaml.MappingNode:
		m := make(map[string]any, len(n.Content)/2)
		for i := 0; i+1 < len(n.Content); i += 2 {
			m[n.Content[i].Value] = nodeValue(n.Content[i+1])
		}
		return m
	case yaml.SequenceNode:
		s := make([]any, 0, len(n.Content))
		for _, c := range n.Content {
			s = append(s, nodeValue(c))
		}
		return s
	default:
		var v any
		if err := n.Decode(&v); err != nil {
			return n.Value
		}
		return v
	}
}

// resolve returns the file and node a $ref points to, relative to file cur.
func (l *specLoader) resolve(cur, ref string) (string, any, error) {
	target, pointer, _ := strings.Cut(ref, "#")
	file := cur
	if target != "" {
		file = filepath.ToSlash(filepath.Clean(filepath.Join(filepath.Dir(cur), target)))
	}
	doc, err := l.file(file)
	if err != nil {
		return "", nil, err
	}
	node := doc
	for _, tok := range strings.Split(strings.TrimPrefix(pointer, "/"), "/") {
		if tok == "" {
			continue
		}
		tok = strings.NewReplacer("~1", "/", "~0", "~").Replace(tok)
		m, ok := node.(map[string]any)
		if !ok {
			return "", nil, fmt.Errorf("unresolvable $ref %q from %s", ref, cur)
		}
		if node, ok = m[tok]; !ok {
			return "", nil, fmt.Errorf("unresolvable $ref %q from %s: %q not found", ref, cur, tok)
		}
	}
	return file, node, nil
}

// deref follows $ref chains until it reaches a non-reference node.
func (l *specLoader) deref(cur string, node any) (string, any, error) {
	for range 32 {
		m, ok := node.(map[string]any)
		if !ok {
			return cur, node, nil
		}
		ref, ok := m["$ref"].(string)
		if !ok {
			return cur, node, nil
		}
		var err error
		if cur, node, err = l.resolve(cur, ref); err != nil {
			return "", nil, err
		}
	}
	return "", nil, fmt.Errorf("$ref chain too deep in %s", cur)
}
