package apiversion

import (
	"strings"
	"testing"
)

// TestRegistryMappings pins the supported SDP releases to their API versions
// and specification branches.
func TestRegistryMappings(t *testing.T) {
	want := map[string]struct{ sdp, branch string }{
		"v21": {"6.4", "version-21"},
		"v22": {"6.5", "version-22"},
		"v23": {"6.6", "version-23"},
		"v24": {"6.7", "version-24"},
		"v25": {"7.0", "version-25"},
	}
	reg := Default()
	if reg.Repository != "https://github.com/appgate/sdp-api-specification.git" {
		t.Errorf("repository = %s", reg.Repository)
	}
	if len(reg.Versions) != len(want) {
		t.Fatalf("registered %v", reg.Names())
	}
	for _, v := range reg.Versions {
		w, ok := want[v.APIVersion]
		if !ok || v.SDPVersion != w.sdp || v.Branch != w.branch {
			t.Errorf("%+v, want %+v", v, w)
		}
	}
	if got := strings.Join(reg.Names(), ","); got != "v21,v22,v23,v24,v25" {
		t.Errorf("versions must be sorted, got %s", got)
	}
	v, err := reg.Lookup("v24")
	if err != nil || v.Number() != 24 || v.AcceptHeader() != "application/vnd.appgate.peer-v24+json" {
		t.Errorf("Lookup(v24) = %+v, %v", v, err)
	}
	if _, err := reg.Lookup("v20"); err == nil || !strings.Contains(err.Error(), "supported versions: v21, v22, v23, v24, v25") {
		t.Errorf("unsupported version error: %v", err)
	}
}

func TestParseRejectsInvalidRegistries(t *testing.T) {
	const sha = "4d21f5ac71adea3f2f1ff15699939e74b2d4f606"
	entry := func(api, sdp, branch, commit string) string {
		return `{"api_version":"` + api + `","sdp_version":"` + sdp + `","branch":"` + branch + `","commit":"` + commit + `"}`
	}
	doc := func(entries ...string) string {
		return `{"repository":"r","versions":[` + strings.Join(entries, ",") + `]}`
	}
	for name, tc := range map[string]struct{ doc, want string }{
		"invalid json":  {`{`, "parse version registry"},
		"no repository": {`{"versions":[` + entry("v25", "7.0", "version-25", sha) + `]}`, "repository is required"},
		"no versions":   {doc(), "no versions registered"},
		"bad name":      {doc(entry("25", "7.0", "version-25", sha)), `invalid api_version "25"`},
		"wrong branch":  {doc(entry("v25", "7.0", "version-24", sha)), "must use branch version-25"},
		"short commit":  {doc(entry("v25", "7.0", "version-25", "4d21f5a")), "not a full 40 character SHA"},
		"duplicate":     {doc(entry("v25", "7.0", "version-25", sha), entry("v25", "7.0", "version-25", sha)), "duplicate api_version v25"},
		"missing sdp":   {doc(entry("v25", "", "version-25", sha)), "sdp_version is required"},
	} {
		if _, err := Parse([]byte(tc.doc)); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: got %v, want %q", name, err, tc.want)
		}
	}
}
