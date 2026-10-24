package provider

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
	"github.com/appgate/terraform-provider-appgatesdp/internal/sdp"
)

var _ datasource.DataSourceWithConfigure = &genericDataSource{}

// readEndpoint reads a read-only data source: path parameters and set query
// parameters come from the configuration.
func readEndpoint(ctx context.Context, c *sdp.Client, m *apimeta.Resource, config tftypes.Value) (map[string]any, error) {
	path := m.Path
	query := url.Values{}
	var attrs map[string]tftypes.Value
	if err := config.As(&attrs); err != nil {
		return nil, err
	}
	for _, p := range m.Params {
		v, ok := attrs[p.TF]
		if !ok || v.IsNull() || !v.IsKnown() {
			if p.In == "path" {
				return nil, fmt.Errorf("%s must be set", p.TF)
			}
			continue
		}
		s, err := paramString(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", p.TF, err)
		}
		if p.In == "path" {
			path = strings.Replace(path, "{"+p.Name+"}", url.PathEscape(s), 1)
		} else {
			query.Set(p.Name, s)
		}
	}
	var out map[string]any
	err := c.DoAccept(ctx, m.Accept, http.MethodGet, path, query, nil, &out)
	return out, err
}

// paramString renders a primitive Terraform value as a request parameter.
func paramString(v tftypes.Value) (string, error) {
	switch t := v.Type(); {
	case t.Is(tftypes.String):
		var s string
		err := v.As(&s)
		return s, err
	case t.Is(tftypes.Bool):
		var b bool
		err := v.As(&b)
		return strconv.FormatBool(b), err
	case t.Is(tftypes.Number):
		n := new(big.Float)
		err := v.As(&n)
		return n.Text('f', -1), err
	}
	return "", fmt.Errorf("unsupported parameter type %s", v.Type())
}

// genericDataSource implements every generated data source. Objects are
// looked up by id, or by exact name; singletons need neither.
type genericDataSource struct {
	name     string
	schemaFn func(context.Context) dschema.Schema
	data     *providerData
}

func (d *genericDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.name
}

func (d *genericDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = d.schemaFn(ctx)
}

func (d *genericDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	data, ok := req.ProviderData.(*providerData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("got %T", req.ProviderData))
		return
	}
	d.data = data
}

func (d *genericDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.data == nil {
		resp.Diagnostics.AddError("Provider not configured", "The appgatesdp provider must be configured before reading "+TypePrefix+d.name+".")
		return
	}
	m, ok := d.data.api.DataSources[d.name]
	if !ok {
		resp.Diagnostics.AddError("Data source not supported by API version", unsupportedDetail(TypePrefix+d.name, d.data.version, func(a *apimeta.API) bool {
			_, ok := a.DataSources[d.name]
			return ok
		}))
		return
	}
	c := d.data.client
	id, name := stringAttr(req.Config.Raw, "id"), stringAttr(req.Config.Raw, "name")
	lookup := &genericResource{name: d.name, data: d.data}
	var out map[string]any
	var err error
	switch {
	case m.ReadOnly:
		out, err = readEndpoint(ctx, c, m, req.Config.Raw)
	case m.Singleton:
		out, err = lookup.getObject(ctx, m, "")
	case id != "":
		out, err = lookup.getObject(ctx, m, id)
		if err == nil && !matchesDiscriminator(out, m) {
			err = fmt.Errorf("object %s is not of %s %q", id, m.Discriminator.Property, m.Discriminator.Value)
		}
	case name != "" && m.ReadList == "":
		out, err = findByName(ctx, c, m, name)
		out = m.Wrap(out)
	default:
		resp.Diagnostics.AddError("Missing lookup attribute", "Set either id or name.")
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Cannot read "+TypePrefix+d.name, err.Error())
		return
	}
	typ := req.Config.Raw.Type().(tftypes.Object)
	v, err := objectFromJSON(out, typ, m.Fields, nil, req.Config.Raw, refresh)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected API response for "+TypePrefix+d.name, err.Error())
		return
	}
	if stringAttr(v, "id") == "" {
		v = withAttr(v, "id", tftypes.NewValue(tftypes.String, d.name))
	}
	resp.State.Raw = v
}
