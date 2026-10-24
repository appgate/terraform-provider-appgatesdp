package codegen

import (
	"context"
	"encoding/json"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
)

func TestFetchUsesPinnedCommitOfBranch(t *testing.T) {
	repo := newFixtureRepo(t)
	f := &Fetcher{Repository: repo.dir, CacheDir: t.TempDir()}
	ctx := context.Background()

	dir, err := f.Fetch(ctx, fixtureVersion("99", repo.commits["v99"]))
	if err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(filepath.Join(dir, "common.yml"))
	if strings.Contains(string(b), "v99 wording") {
		t.Fatal("fetched the branch head instead of the pinned commit")
	}
	if !strings.Contains(dir, repo.commits["v99"]) {
		t.Fatalf("cache dir %s must be keyed by commit", dir)
	}
	head, err := f.BranchHead(ctx, "version-99")
	if err != nil || head != repo.commits["v99next"] {
		t.Fatalf("BranchHead = %s, %v", head, err)
	}

	// A second fetch is served from the cache, even without the repository.
	f.Repository = filepath.Join(t.TempDir(), "gone")
	if _, err := f.Fetch(ctx, fixtureVersion("99", repo.commits["v99"])); err != nil {
		t.Fatalf("cached fetch failed: %v", err)
	}
}

func TestFetchErrors(t *testing.T) {
	repo := newFixtureRepo(t)
	ctx := context.Background()
	tests := map[string]struct {
		version string
		commit  string
		want    string
	}{
		"missing branch":                     {"96", repo.commits["v98"], "branch version-96 not found"},
		"commit not on branch":               {"98", repo.commits["orphan"], "is not part of branch version-98"},
		"commit from another version branch": {"98", repo.commits["v99"], "is not part of branch version-98"},
		"unknown commit":                     {"98", strings.Repeat("ab", 20), "is not part of branch version-98"},
		"wrong spec version":                 {"97", repo.commits["v97bad"], `declares "API version 96.0", expected API version 97`},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f := &Fetcher{Repository: repo.dir, CacheDir: t.TempDir()}
			_, err := f.Fetch(ctx, fixtureVersion(tc.version, tc.commit))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("got %v, want %q", err, tc.want)
			}
		})
	}
}

func TestBuildFixture(t *testing.T) {
	res, err := build("v98", "testdata/spec", mustDefinitions(t))
	if err != nil {
		t.Fatal(err)
	}
	if res.unsupported["oval_widget"] == "" {
		t.Fatalf("oval_widget (missing variant) must be unsupported: %v", res.unsupported)
	}
	for name := range res.api.Resources {
		if strings.Contains(name, "gadget") {
			t.Errorf("v98 has no /gadgets, got %s", name)
		}
	}
	status := res.api.DataSources["widget_status_report"]
	switch {
	case status == nil || !status.ReadOnly || status.Path != "/widgets/{id}/status":
		t.Fatalf("read-only endpoint must become a renamed data source, got %+v", status)
	case len(status.Params) != 3:
		t.Errorf("params %+v", status.Params)
	case status.Fields["details"] == nil || status.Fields["details"].Type != "json":
		t.Error("free-form objects must become JSON attributes")
	}
	if status.Accept != "application/json" || status.Fields["healthy"].Type != "json" {
		t.Errorf("correction not applied: accept %q, healthy %+v", status.Accept, status.Fields["healthy"])
	}
	if !slices.ContainsFunc(res.warnings, func(w string) bool { return strings.Contains(w, "corrected: healthy is returned as a string") }) {
		t.Error("corrections must be reported")
	}
	for _, p := range status.Params {
		if p.TF == "level" && (!p.Required || len(p.Enum) != 2) {
			t.Errorf("required enum query parameter not recorded: %+v", p)
		}
	}
	for path, want := range map[string]string{"/login": "excluded", "/widgets": "resource", "/widgets/{id}/status": "data source", "/tokens/{token-id}": "resource"} {
		if got := res.classes[path].Kind; got != want {
			t.Errorf("%s classified %q, want %q", path, got, want)
		}
	}
	if _, ok := res.api.DataSources["round_widget"]; !ok {
		t.Error("every resource needs a data source")
	}
	round := res.api.Resources["round_widget"]
	if round.Discriminator == nil || round.Discriminator.Value != "Round" || !round.Delete {
		t.Fatalf("round_widget meta %+v", round)
	}
	f := round.Fields
	switch {
	case f["type"] != nil:
		t.Error("discriminator property must not be a Terraform attribute")
	case f["name"] == nil || !f["name"].Required:
		t.Error("allOf base properties must be merged and keep required")
	case f["radius"] == nil || !f["radius"].Required || f["radius"].Type != "int64":
		t.Error("variant properties missing")
	case f["side"] != nil:
		t.Error("properties of other variants must not leak into a variant")
	case !f["created"].ReadOnly || !f["secret"].WriteOnly:
		t.Error("readOnly/writeOnly not detected")
	case f["tags"].Type != "set":
		t.Error("uniqueItems arrays of primitives must become sets")
	case f["labels"].Type != "map" || f["labels"].Elem != "string":
		t.Error("additionalProperties must become a map")
	case f["parent"] != nil:
		t.Error("recursive references must be dropped")
	}
	tok := res.api.Resources["token"]
	switch {
	case tok == nil:
		t.Fatal("token resource missing")
	case tok.ReadList != "tokens" || tok.ID() != "tokenId" || !tok.Immutable || !tok.Delete:
		t.Errorf("token meta %+v", tok)
	case tok.DeleteObjectPath("abc") != "/tokens/abc":
		t.Errorf("delete path %s", tok.DeleteObjectPath("abc"))
	case !tok.Fields["value"].WriteOnly || !tok.Fields["owner"].Required || !tok.Fields["expires"].ReadOnly:
		t.Error("token fields must come from the POST body and the list items")
	}
	spec, _ := json.Marshal(res.spec)
	if !strings.Contains(string(spec), `"value":{"format":"password"`) {
		t.Error("definition sensitive list must mark value sensitive")
	}
	if res.api.Resources["settings"] == nil || !res.api.Resources["settings"].Singleton {
		t.Error("settings must be a singleton")
	}
	if _, ok := res.api.DataSources["round_widget"]; !ok {
		t.Error("data source missing")
	}

	for _, bad := range []string{"allOf", "oneOf", "$ref", "discriminator", `(?!bad)`} {
		if strings.Contains(string(spec), bad) {
			t.Errorf("generator input still contains %q", bad)
		}
	}
	joined := strings.Join(res.warnings, "\n")
	for _, want := range []string{"dropped recursive reference", "unsupported by Go regexp"} {
		if !strings.Contains(joined, want) {
			t.Errorf("missing warning %q in\n%s", want, joined)
		}
	}
}

func TestBuildIsDeterministic(t *testing.T) {
	defs := mustDefinitions(t)
	first, err := build("v98", "testdata/spec", defs)
	if err != nil {
		t.Fatal(err)
	}
	for range 5 {
		next, _ := build("v98", "testdata/spec", defs)
		a, _ := json.Marshal([]any{first.spec, first.config, first.api, first.warnings})
		b, _ := json.Marshal([]any{next.spec, next.config, next.api, next.warnings})
		if string(a) != string(b) {
			t.Fatal("build output differs between runs")
		}
	}
}

func mustDefinitions(t *testing.T) Definitions {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "codegen"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, DefinitionsFile), []byte(fixtureDefinitions), 0o644); err != nil {
		t.Fatal(err)
	}
	defs, err := LoadDefinitions(root)
	if err != nil {
		t.Fatal(err)
	}
	return defs
}

// snapshot reads every file under root.
func snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(root, p)
		if strings.HasPrefix(rel, CacheDir) {
			return nil
		}
		b, err := os.ReadFile(p)
		out[rel] = string(b)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// TestGeneratePipeline runs the real generators on the fixture branches.
func TestGeneratePipeline(t *testing.T) {
	if testing.Short() {
		t.Skip("runs the generator tools")
	}
	repo := newFixtureRepo(t)
	root := newFixtureRoot(t, repo, fixtureVersion("98", repo.commits["v98"]), fixtureVersion("99", repo.commits["v99"]))
	ctx := context.Background()
	opts := Options{Root: root, ToolDir: moduleRoot(t)}
	if err := Generate(ctx, opts); err != nil {
		t.Fatal(err)
	}
	first := snapshot(t, root)
	for _, want := range []string{
		"codegen/ir/v98.json", "codegen/ir/v99.json", "codegen/ir/merged.json",
		"codegen/reports/v98.md", "codegen/reports/v99.md",
		"internal/apiversion/generated/v98.json", "internal/apiversion/generated/v99.json",
		"internal/generated/resource_round_widget/round_widget_resource_gen.go",
		"internal/generated/resource_round_gadget/round_gadget_resource_gen.go",
		"internal/generated/datasource_widget_status_report/widget_status_report_data_source_gen.go",
		"internal/generated/datasource_round_widget/round_widget_data_source_gen.go",
		SchemaIndexFile, VersionsDoc,
	} {
		if _, ok := first[want]; !ok {
			t.Errorf("missing %s", want)
		}
	}
	if !strings.Contains(first["codegen/reports/v99.md"], repo.commits["v99"]) || !strings.Contains(first[VersionsDoc], repo.commits["v98"][:12]) {
		t.Error("reports must document the pinned commit")
	}
	var v98 apimeta.API
	_ = json.Unmarshal([]byte(first["internal/apiversion/generated/v98.json"]), &v98)
	if _, ok := v98.Resources["round_gadget"]; ok {
		t.Error("v98 must not use the v99 specification")
	}
	if v98.Resources["round_widget"].Fields["color"] != nil {
		t.Error("v98 round_widget must not have the v99 color field")
	}
	merged := first["internal/generated/resource_round_widget/round_widget_resource_gen.go"]
	if !strings.Contains(merged, `"color"`) || !strings.Contains(merged, "Supported in API versions: v99.") {
		t.Error("merged schema must contain v99 additions with their availability")
	}
	if !strings.Contains(merged, `"secret"`) || !strings.Contains(merged, "Sensitive:") {
		t.Error("write-only strings must be sensitive")
	}

	t.Run("deterministic", func(t *testing.T) {
		if err := Generate(ctx, opts); err != nil {
			t.Fatal(err)
		}
		if second := snapshot(t, root); !maps.Equal(first, second) {
			for k := range first {
				if first[k] != second[k] {
					t.Errorf("%s changed on regeneration", k)
				}
			}
		}
	})

	t.Run("single version leaves other versions untouched", func(t *testing.T) {
		v98Files := []string{"codegen/ir/v98.json", "codegen/reports/v98.md", "internal/apiversion/generated/v98.json"}
		before := map[string]os.FileInfo{}
		for _, f := range v98Files {
			before[f], _ = os.Stat(filepath.Join(root, f))
		}
		// Move the v99 pin to the newer commit and regenerate only v99.
		writeRegistry(t, root, repo.dir, fixtureVersion("98", repo.commits["v98"]), fixtureVersion("99", repo.commits["v99next"]))
		if err := Generate(ctx, Options{Root: root, ToolDir: moduleRoot(t), Versions: []string{"v99"}}); err != nil {
			t.Fatal(err)
		}
		for _, f := range v98Files {
			after, _ := os.Stat(filepath.Join(root, f))
			if !after.ModTime().Equal(before[f].ModTime()) {
				t.Errorf("%s was rewritten", f)
			}
		}
		now := snapshot(t, root)
		if !strings.Contains(now["codegen/reports/v99.md"], repo.commits["v99next"]) {
			t.Error("v99 report must reference the new pin")
		}
		if !strings.Contains(now["codegen/ir/v99.json"], "v99 wording") {
			t.Error("v99 must be generated from the new pinned commit")
		}
	})

	t.Run("unsupported version", func(t *testing.T) {
		err := Generate(ctx, Options{Root: root, ToolDir: moduleRoot(t), Versions: []string{"v42"}})
		if err == nil || !strings.Contains(err.Error(), `unsupported API version "v42"; supported versions: v98, v99`) {
			t.Fatalf("got %v", err)
		}
	})

	t.Run("stale versions are removed", func(t *testing.T) {
		writeRegistry(t, root, repo.dir, fixtureVersion("99", repo.commits["v99next"]))
		if err := Generate(ctx, Options{Root: root, ToolDir: moduleRoot(t), Versions: []string{"v99"}}); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(filepath.Join(root, "codegen/ir/v98.json")); !os.IsNotExist(err) {
			t.Error("v98 IR must be removed once v98 is unregistered")
		}
	})
}

func TestPin(t *testing.T) {
	repo := newFixtureRepo(t)
	root := newFixtureRoot(t, repo, fixtureVersion("98", repo.commits["v98"]), fixtureVersion("99", repo.commits["v99"]))
	if err := Pin(context.Background(), Options{Root: root, Versions: []string{"v99"}}); err != nil {
		t.Fatal(err)
	}
	reg, err := LoadRegistry(root)
	if err != nil {
		t.Fatal(err)
	}
	if reg.Versions[0].Commit != repo.commits["v98"] || reg.Versions[1].Commit != repo.commits["v99next"] {
		t.Fatalf("pin must only move the selected version: %+v", reg.Versions)
	}
}
