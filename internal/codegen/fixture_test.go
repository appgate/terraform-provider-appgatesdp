package codegen

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion"
)

// fixtureRepo is a local sdp-api-specification look-alike with branches
// version-98 and version-99 built from testdata/spec.
type fixtureRepo struct {
	dir string
	// commits by name: v98, v99, v99next (newer commit on version-99),
	// orphan (commit on another branch), v97bad (declares the wrong version).
	commits map[string]string
}

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

// writeSpec copies testdata/spec into dir, applying string replacements per file.
func writeSpec(t *testing.T, dir string, edits map[string][][2]string) {
	t.Helper()
	entries, err := os.ReadDir("testdata/spec")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		b, err := os.ReadFile(filepath.Join("testdata/spec", e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		for _, r := range edits[e.Name()] {
			if !strings.Contains(s, r[0]) {
				t.Fatalf("fixture edit %q not found in %s", r[0], e.Name())
			}
			s = strings.Replace(s, r[0], r[1], 1)
		}
		if err := os.WriteFile(filepath.Join(dir, e.Name()), []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func commitAll(t *testing.T, dir, msg string) string {
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-q", "-m", msg)
	return git(t, dir, "rev-parse", "HEAD")
}

func newFixtureRepo(t *testing.T) *fixtureRepo {
	t.Helper()
	dir := t.TempDir()
	r := &fixtureRepo{dir: dir, commits: map[string]string{}}
	git(t, dir, "init", "-q", "-b", "version-98")
	writeSpec(t, dir, nil)
	r.commits["v98"] = commitAll(t, dir, "v98")

	git(t, dir, "checkout", "-q", "-b", "version-99")
	v99 := map[string][][2]string{
		"api_specs.yml": {
			{"API version 98.0", "API version 99.0"},
			{"  /settings:", "  /gadgets:\n    $ref: './widget.yml#/widgets'\n  /gadgets/{id}:\n    $ref: './widget.yml#/widgets-id'\n  /settings:"},
		},
		"widget.yml": {{"        required: [radius]\n        properties:\n", "        required: [radius]\n        properties:\n          color:\n            type: string\n"}},
	}
	writeSpec(t, dir, v99)
	r.commits["v99"] = commitAll(t, dir, "v99")
	v99["common.yml"] = [][2]string{{"description: Name of the object.", "description: Name of the object (v99 wording)."}}
	writeSpec(t, dir, v99)
	r.commits["v99next"] = commitAll(t, dir, "v99 docs")

	git(t, dir, "checkout", "-q", "-b", "version-97", r.commits["v98"])
	writeSpec(t, dir, map[string][][2]string{"api_specs.yml": {{"API version 98.0", "API version 96.0"}}})
	r.commits["v97bad"] = commitAll(t, dir, "wrong version")

	git(t, dir, "checkout", "-q", "-b", "unrelated", r.commits["v98"])
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	r.commits["orphan"] = commitAll(t, dir, "not on a version branch")
	git(t, dir, "checkout", "-q", "version-98")
	return r
}

const fixtureDefinitions = `resources:
  round_widget:
    path: /widgets
    variant: {property: type, value: Round}
  square_widget:
    path: /widgets
    variant: {property: type, value: Square}
  oval_widget:
    path: /widgets
    variant: {property: type, value: Oval}
  settings:
    path: /settings
    singleton: true
  token:
    path: /tokens
    read_list: tokens
    id_property: tokenId
    delete_path: /tokens/{token-id}
    immutable: true
    sensitive: [value]
    required: [owner]
names:
  /widgets/{id}/status: widget_status_report
exclude:
  /login: authentication is handled by the provider configuration
`

// newFixtureRoot creates a repository root with a registry and definitions.
func newFixtureRoot(t *testing.T, repo *fixtureRepo, versions ...apiversion.Version) string {
	t.Helper()
	root := t.TempDir()
	writeRegistry(t, root, repo.dir, versions...)
	if err := os.MkdirAll(filepath.Join(root, "codegen"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, DefinitionsFile), []byte(fixtureDefinitions), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func writeRegistry(t *testing.T, root, repoDir string, versions ...apiversion.Version) {
	t.Helper()
	b, err := json.Marshal(apiversion.Registry{Repository: repoDir, Versions: versions})
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(root, RegistryFile)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, b, 0o644); err != nil {
		t.Fatal(err)
	}
}

func fixtureVersion(n, commit string) apiversion.Version {
	return apiversion.Version{APIVersion: "v" + n, SDPVersion: "9." + n, Branch: "version-" + n, Commit: commit}
}

// moduleRoot is the repository root, where the pinned generator tools run.
func moduleRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	return root
}
