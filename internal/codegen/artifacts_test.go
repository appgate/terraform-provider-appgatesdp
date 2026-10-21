package codegen

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
)

// TestCommittedArtifactsMatchRegistry checks that the committed generated
// files belong to the registered versions and pinned commits. When the
// specification cache from a previous `make generate` is present, the
// metadata is also rebuilt and compared; CI regenerates everything and fails
// on any diff.
func TestCommittedArtifactsMatchRegistry(t *testing.T) {
	root := moduleRoot(t)
	reg, err := LoadRegistry(root)
	if err != nil {
		t.Fatal(err)
	}
	defs, err := LoadDefinitions(root)
	if err != nil {
		t.Fatal(err)
	}
	doc, err := os.ReadFile(filepath.Join(root, VersionsDoc))
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range reg.Versions {
		t.Run(v.APIVersion, func(t *testing.T) {
			report, err := os.ReadFile(filepath.Join(root, ReportDir, v.APIVersion+".md"))
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range []string{"`" + v.Branch + "`", "`" + v.Commit + "`", "| SDP version | " + v.SDPVersion + " |"} {
				if !strings.Contains(string(report), want) {
					t.Errorf("report does not document %s", want)
				}
			}
			if !strings.Contains(string(doc), v.Commit[:12]) {
				t.Errorf("%s does not document commit %s", VersionsDoc, v.Commit)
			}
			if _, err := os.Stat(filepath.Join(root, IRDir, v.APIVersion+".json")); err != nil {
				t.Error(err)
			}
			var committed apimeta.API
			b, err := os.ReadFile(filepath.Join(root, MetadataDir, v.APIVersion+".json"))
			if err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(b, &committed); err != nil || committed.Version != v.APIVersion {
				t.Fatalf("metadata version %q, err %v", committed.Version, err)
			}
			cache := filepath.Join(root, CacheDir, "specs", v.APIVersion, v.Commit)
			if _, err := os.Stat(filepath.Join(cache, ".complete")); err != nil {
				t.Log("specification cache missing; skipping rebuild comparison")
				return
			}
			res, err := build(v.APIVersion, cache, defs)
			if err != nil {
				t.Fatal(err)
			}
			rebuilt, _ := json.Marshal(res.api)
			var fresh apimeta.API
			_ = json.Unmarshal(rebuilt, &fresh)
			if !reflect.DeepEqual(&fresh, &committed) {
				t.Errorf("committed metadata is stale; run `make generate API_VERSION=%s`", v.APIVersion)
			}
		})
	}
	for _, dir := range []string{IRDir, MetadataDir, ReportDir} {
		entries, _ := os.ReadDir(filepath.Join(root, dir))
		for _, e := range entries {
			ext := filepath.Ext(e.Name())
			base := strings.TrimSuffix(e.Name(), ext)
			if ext == ".go" || base == "merged" {
				continue
			}
			if _, err := reg.Lookup(base); err != nil {
				t.Errorf("%s/%s belongs to no registered version", dir, e.Name())
			}
		}
	}
}
