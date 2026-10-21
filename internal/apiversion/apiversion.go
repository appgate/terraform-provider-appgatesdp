// Package apiversion is the central registry of supported SDP Admin API
// versions. versions.json maps each API version to its SDP release, the
// sdp-api-specification branch and the pinned commit the code was generated
// from. The code generator and the provider both read it, so registering a new
// version is a one-line change followed by `make generate`.
package apiversion

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
)

//go:embed versions.json
var registryJSON []byte

// Registry lists every supported API version.
type Registry struct {
	Repository string    `json:"repository"`
	Versions   []Version `json:"versions"`
}

// Version is one supported Admin API version.
type Version struct {
	APIVersion string `json:"api_version"`
	SDPVersion string `json:"sdp_version"`
	Branch     string `json:"branch"`
	Commit     string `json:"commit"`
}

// Number returns the numeric API version, e.g. 25 for v25.
func (v Version) Number() int {
	n, _ := strconv.Atoi(strings.TrimPrefix(v.APIVersion, "v"))
	return n
}

// AcceptHeader is the media type the Admin API expects for this version.
func (v Version) AcceptHeader() string {
	return fmt.Sprintf("application/vnd.appgate.peer-v%d+json", v.Number())
}

var (
	apiVersionRE = regexp.MustCompile(`^v[1-9][0-9]*$`)
	commitRE     = regexp.MustCompile(`^[0-9a-f]{40}$`)
)

// Parse decodes and validates a registry document.
func Parse(b []byte) (Registry, error) {
	var r Registry
	if err := json.Unmarshal(b, &r); err != nil {
		return r, fmt.Errorf("parse version registry: %w", err)
	}
	if r.Repository == "" {
		return r, fmt.Errorf("version registry: repository is required")
	}
	if len(r.Versions) == 0 {
		return r, fmt.Errorf("version registry: no versions registered")
	}
	seen := map[string]bool{}
	for _, v := range r.Versions {
		switch {
		case !apiVersionRE.MatchString(v.APIVersion):
			return r, fmt.Errorf("version registry: invalid api_version %q, expected v<N>", v.APIVersion)
		case seen[v.APIVersion]:
			return r, fmt.Errorf("version registry: duplicate api_version %s", v.APIVersion)
		case v.Branch != fmt.Sprintf("version-%d", v.Number()):
			return r, fmt.Errorf("version registry: %s must use branch version-%d, got %q", v.APIVersion, v.Number(), v.Branch)
		case !commitRE.MatchString(v.Commit):
			return r, fmt.Errorf("version registry: %s commit %q is not a full 40 character SHA", v.APIVersion, v.Commit)
		case v.SDPVersion == "":
			return r, fmt.Errorf("version registry: %s sdp_version is required", v.APIVersion)
		}
		seen[v.APIVersion] = true
	}
	slices.SortFunc(r.Versions, func(a, b Version) int { return a.Number() - b.Number() })
	return r, nil
}

// Lookup returns the registered version with the given name.
func (r Registry) Lookup(name string) (Version, error) {
	for _, v := range r.Versions {
		if v.APIVersion == name {
			return v, nil
		}
	}
	return Version{}, fmt.Errorf("unsupported API version %q; supported versions: %s", name, strings.Join(r.Names(), ", "))
}

// Names returns all API version names in ascending order.
func (r Registry) Names() []string {
	names := make([]string, len(r.Versions))
	for i, v := range r.Versions {
		names[i] = v.APIVersion
	}
	return names
}

var (
	defaultOnce sync.Once
	defaultReg  Registry
	defaultErr  error
)

// Default returns the registry compiled into the provider.
func Default() Registry {
	defaultOnce.Do(func() { defaultReg, defaultErr = Parse(registryJSON) })
	if defaultErr != nil {
		// The registry is validated by unit tests; failing here is a build defect.
		panic(defaultErr)
	}
	return defaultReg
}
