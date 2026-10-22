// Package codegen implements the provider generation pipeline:
//
//	versions.json ─► fetch pinned spec per branch ─► normalize ─► tfplugingen-openapi (IR per version)
//	                                                               │
//	                     per-version API metadata ◄────────────────┤
//	                                                               ▼
//	                        merge IRs (superset) ─► tfplugingen-framework ─► internal/generated
package codegen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion"
	"gopkg.in/yaml.v3"
)

// Layout of generated and hand-maintained files, relative to the repository root.
const (
	RegistryFile    = "internal/apiversion/versions.json"
	DefinitionsFile = "codegen/resources.yml"
	IRDir           = "codegen/ir"
	ReportDir       = "codegen/reports"
	MetadataDir     = "internal/apiversion/generated"
	FrameworkDir    = "internal/generated"
	SchemaIndexFile = "internal/provider/zz_generated_schemas.go"
	VersionsDoc     = "templates/guides/api-versions.md.tmpl"
	CacheDir        = ".cache"
)

// Options configures a generation run.
type Options struct {
	// Root is the repository root that is read and written.
	Root string
	// ToolDir is the Go module directory the pinned generator tools run in.
	// Defaults to Root.
	ToolDir string
	// Versions limits generation to these API versions; empty means all.
	Versions []string
	// Log receives progress output.
	Log func(format string, args ...any)
}

func (o Options) path(rel string) string { return filepath.Join(o.Root, rel) }

func (o Options) toolDir() string {
	if o.ToolDir != "" {
		return o.ToolDir
	}
	return o.Root
}

func (o Options) logf(format string, args ...any) {
	if o.Log != nil {
		o.Log(format, args...)
	}
}

// LoadRegistry reads the version registry from the repository.
func LoadRegistry(root string) (apiversion.Registry, error) {
	b, err := os.ReadFile(filepath.Join(root, RegistryFile))
	if err != nil {
		return apiversion.Registry{}, err
	}
	return apiversion.Parse(b)
}

// LoadDefinitions reads codegen/resources.yml.
func LoadDefinitions(root string) (Definitions, error) {
	var defs Definitions
	b, err := os.ReadFile(filepath.Join(root, DefinitionsFile))
	if err != nil {
		return defs, err
	}
	dec := yaml.NewDecoder(bytes.NewReader(b))
	dec.KnownFields(true)
	if err := dec.Decode(&defs); err != nil {
		return defs, fmt.Errorf("%s: %w", DefinitionsFile, err)
	}
	for name, d := range defs.Resources {
		if !strings.HasPrefix(d.Path, "/") {
			return defs, fmt.Errorf("%s: resource %s: path must start with /", DefinitionsFile, name)
		}
		if TerraformIdentifier(name) != name {
			return defs, fmt.Errorf("%s: resource name %q must be snake_case", DefinitionsFile, name)
		}
	}
	return defs, nil
}

// Generate runs the pipeline for the selected versions and regenerates the
// merged provider schema from the committed IR of every registered version.
func Generate(ctx context.Context, o Options) error {
	root, err := filepath.Abs(o.Root)
	if err != nil {
		return err
	}
	o.Root = root
	reg, err := LoadRegistry(o.Root)
	if err != nil {
		return err
	}
	defs, err := LoadDefinitions(o.Root)
	if err != nil {
		return err
	}
	selected := reg.Versions
	if len(o.Versions) > 0 {
		selected = nil
		for _, name := range o.Versions {
			v, err := reg.Lookup(name)
			if err != nil {
				return err
			}
			selected = append(selected, v)
		}
	}
	f := &Fetcher{Repository: reg.Repository, CacheDir: o.path(CacheDir)}
	for _, v := range selected {
		if err := generateVersion(ctx, o, f, v, defs); err != nil {
			return err
		}
	}
	if err := removeStale(o, reg); err != nil {
		return err
	}
	return generateMerged(ctx, o, reg)
}

func generateVersion(ctx context.Context, o Options, f *Fetcher, v apiversion.Version, defs Definitions) error {
	o.logf("%s: fetching %s@%s", v.APIVersion, v.Branch, v.Commit[:12])
	dir, err := f.Fetch(ctx, v)
	if err != nil {
		return err
	}
	o.logf("%s: normalizing specification", v.APIVersion)
	res, err := build(v.APIVersion, dir, defs)
	if err != nil {
		return err
	}
	work := o.path(filepath.Join(CacheDir, "build", v.APIVersion))
	if err := os.MkdirAll(work, 0o755); err != nil {
		return err
	}
	specFile := filepath.Join(work, "generator-input.yml")
	configFile := filepath.Join(work, "generator-config.yml")
	if err := writeYAML(specFile, res.spec); err != nil {
		return err
	}
	if err := writeYAML(configFile, res.config); err != nil {
		return err
	}
	irFile := o.path(filepath.Join(IRDir, v.APIVersion+".json"))
	if err := os.MkdirAll(filepath.Dir(irFile), 0o755); err != nil {
		return err
	}
	o.logf("%s: running tfplugingen-openapi", v.APIVersion)
	genLog, err := goTool(ctx, o.toolDir(), "tfplugingen-openapi", "generate", "--config", configFile, "--output", irFile, specFile)
	if err != nil {
		return fmt.Errorf("%s: tfplugingen-openapi: %w", v.APIVersion, err)
	}
	ir, err := readJSON(irFile)
	if err != nil {
		return err
	}
	if err := verifyIR(res.api, ir); err != nil {
		return fmt.Errorf("%s: %w\ngenerator output:\n%s", v.APIVersion, err, genLog)
	}
	// Re-encode for stable formatting independent of the generator's writer.
	if err := writeJSON(irFile, ir); err != nil {
		return err
	}
	if err := writeJSON(o.path(filepath.Join(MetadataDir, v.APIVersion+".json")), res.api); err != nil {
		return err
	}
	return writeReport(o.path(filepath.Join(ReportDir, v.APIVersion+".md")), v, res)
}

// verifyIR fails when the generator skipped a resource or produced attributes
// the runtime cannot map to JSON properties.
func verifyIR(api *apimeta.API, ir map[string]any) error {
	check := func(kind string, want map[string]*apimeta.Resource) error {
		got := map[string]map[string]any{}
		for _, r := range listOf(ir[kind]) {
			rm := r.(map[string]any)
			got[rm["name"].(string)] = rm
		}
		for _, name := range sortedKeys(want) {
			r, ok := got[name]
			if !ok {
				return fmt.Errorf("tfplugingen-openapi did not generate %s %s; see generator output", kind, name)
			}
			schema, _ := r["schema"].(map[string]any)
			extra := []string{"id"}
			for _, p := range want[name].Params {
				extra = append(extra, p.TF)
			}
			if err := verifyAttrs(name, listOf(schema["attributes"]), want[name].Fields, extra); err != nil {
				return err
			}
		}
		return nil
	}
	if err := check("resources", api.Resources); err != nil {
		return err
	}
	return check("datasources", api.DataSources)
}

// verifyAttrs checks generated attributes against the API fields; extra
// lists root attributes that are not JSON fields (id, request parameters).
func verifyAttrs(at string, attrs []any, fields map[string]*apimeta.Field, extra []string) error {
	for _, a := range attrs {
		am := a.(map[string]any)
		name := am["name"].(string)
		f, ok := fields[name]
		if !ok {
			if slices.Contains(extra, name) {
				continue
			}
			return fmt.Errorf("%s.%s: generated attribute has no API field mapping", at, name)
		}
		_, body := attrType(am)
		if _, nested := nestedAttributes(body); nested != nil {
			if err := verifyAttrs(at+"."+name, nested, f.Fields, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

func generateMerged(ctx context.Context, o Options, reg apiversion.Registry) error {
	var irs []versionedIR
	apis := map[string]*apimeta.API{}
	for _, v := range reg.Versions {
		ir, err := readJSON(o.path(filepath.Join(IRDir, v.APIVersion+".json")))
		if err != nil {
			return fmt.Errorf("%s: no generated IR (%w); run `make generate API_VERSION=%s`", v.APIVersion, err, v.APIVersion)
		}
		irs = append(irs, versionedIR{v.APIVersion, ir})
		var api apimeta.API
		if err := readJSONInto(o.path(filepath.Join(MetadataDir, v.APIVersion+".json")), &api); err != nil {
			return fmt.Errorf("%s: no generated metadata (%w); run `make generate API_VERSION=%s`", v.APIVersion, err, v.APIVersion)
		}
		apis[v.APIVersion] = &api
	}
	o.logf("merging %d API versions", len(irs))
	readOnly := map[string]bool{}
	for _, api := range apis {
		for name, ds := range api.DataSources {
			if ds.ReadOnly {
				readOnly[name] = true
			}
		}
	}
	merged, err := mergeIR(irs, readOnly)
	if err != nil {
		return err
	}
	mergedFile := o.path(filepath.Join(IRDir, "merged.json"))
	if err := writeJSON(mergedFile, merged); err != nil {
		return err
	}
	out := o.path(FrameworkDir)
	if err := os.RemoveAll(out); err != nil {
		return err
	}
	o.logf("running tfplugingen-framework")
	for _, kind := range []string{"resources", "data-sources"} {
		if _, err := goTool(ctx, o.toolDir(), "tfplugingen-framework", "generate", kind, "--input", mergedFile, "--output", out); err != nil {
			return fmt.Errorf("tfplugingen-framework generate %s: %w", kind, err)
		}
	}
	if err := pruneFrameworkOutput(out); err != nil {
		return err
	}
	if err := writeSchemaIndex(o.path(SchemaIndexFile), merged); err != nil {
		return err
	}
	return writeVersionsDoc(o.path(VersionsDoc), reg, apis, merged)
}

// removeStale deletes generated files of versions no longer registered.
func removeStale(o Options, reg apiversion.Registry) error {
	known := reg.Names()
	for _, dir := range []string{IRDir, MetadataDir, ReportDir} {
		entries, err := os.ReadDir(o.path(dir))
		if err != nil {
			continue
		}
		for _, e := range entries {
			base := strings.TrimSuffix(strings.TrimSuffix(e.Name(), ".json"), ".md")
			if strings.HasPrefix(base, "v") && !slices.Contains(known, base) {
				o.logf("removing stale %s", filepath.Join(dir, e.Name()))
				if err := os.Remove(filepath.Join(o.path(dir), e.Name())); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// goTool runs a tool pinned by the `tool` directives in go.mod.
func goTool(ctx context.Context, root, tool string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "go", append([]string{"tool", tool}, args...)...)
	cmd.Dir = root
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return buf.String(), fmt.Errorf("%w\n%s", err, buf.String())
	}
	return buf.String(), nil
}

func writeYAML(path string, v any) error {
	b, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func writeJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

func readJSON(path string) (map[string]any, error) {
	var m map[string]any
	return m, readJSONInto(path, &m)
}

func readJSONInto(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// Pin updates the registry so the selected versions point at the current head
// of their specification branch. Generation always uses the pinned commit.
func Pin(ctx context.Context, o Options) error {
	reg, err := LoadRegistry(o.Root)
	if err != nil {
		return err
	}
	f := &Fetcher{Repository: reg.Repository, CacheDir: o.path(CacheDir)}
	for i, v := range reg.Versions {
		if len(o.Versions) > 0 && !slices.Contains(o.Versions, v.APIVersion) {
			continue
		}
		head, err := f.BranchHead(ctx, v.Branch)
		if err != nil {
			return fmt.Errorf("%s: %w", v.APIVersion, err)
		}
		if head != v.Commit {
			o.logf("%s: %s %s -> %s", v.APIVersion, v.Branch, v.Commit[:12], head[:12])
		}
		reg.Versions[i].Commit = head
	}
	for _, name := range o.Versions {
		if _, err := reg.Lookup(name); err != nil {
			return err
		}
	}
	return writeJSON(o.path(RegistryFile), reg)
}
