// Package generated embeds the per-version API metadata written by sdpgen.
package generated

import (
	"embed"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
)

//go:embed v*.json
var files embed.FS

var cache sync.Map

// API returns the generated metadata of an API version, e.g. "v25".
func API(version string) (*apimeta.API, error) {
	if v, ok := cache.Load(version); ok {
		return v.(*apimeta.API), nil
	}
	b, err := files.ReadFile(version + ".json")
	if err != nil {
		return nil, fmt.Errorf("no generated metadata for API %s; run `make generate API_VERSION=%s`", version, version)
	}
	var api apimeta.API
	if err := json.Unmarshal(b, &api); err != nil {
		return nil, fmt.Errorf("decode generated metadata for API %s: %w", version, err)
	}
	v, _ := cache.LoadOrStore(version, &api)
	return v.(*apimeta.API), nil
}
