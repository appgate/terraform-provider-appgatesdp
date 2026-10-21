package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
	"github.com/appgate/terraform-provider-appgatesdp/internal/sdp"
)

// hooks hold SDP behavior that the OpenAPI specification cannot express.
// Keep them small: anything expressible in codegen/resources.yml or the
// specification normalization belongs there instead.
type hooks struct {
	beforeDelete func(ctx context.Context, c *sdp.Client, m *apimeta.Resource, id string) error
}

var resourceHooks = map[string]hooks{
	"appliance": {beforeDelete: deactivateAppliance},
}

// deactivateAppliance deactivates (and wipes) an activated appliance; the API
// refuses to delete activated appliances.
func deactivateAppliance(ctx context.Context, c *sdp.Client, m *apimeta.Resource, id string) error {
	var appliance map[string]any
	p := m.Path + "/" + url.PathEscape(id)
	if err := c.Do(ctx, http.MethodGet, p, nil, nil, &appliance); err != nil {
		if sdp.IsNotFound(err) {
			return nil
		}
		return err
	}
	if appliance["activated"] != true {
		return nil
	}
	if err := c.Do(ctx, http.MethodPost, p+"/deactivate", url.Values{"wipe": {"true"}}, nil, nil); err != nil {
		return fmt.Errorf("deactivate appliance before delete: %w", err)
	}
	return nil
}
