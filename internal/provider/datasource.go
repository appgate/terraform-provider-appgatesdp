package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
)

var _ datasource.DataSourceWithConfigure = &genericDataSource{}

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
	var out map[string]any
	var err error
	switch {
	case m.Singleton:
		err = c.Do(ctx, http.MethodGet, m.Path, nil, nil, &out)
	case id != "":
		err = c.Do(ctx, http.MethodGet, m.Path+"/"+url.PathEscape(id), nil, nil, &out)
		if err == nil && !matchesDiscriminator(out, m) {
			err = fmt.Errorf("object %s is not of %s %q", id, m.Discriminator.Property, m.Discriminator.Value)
		}
	case name != "":
		out, err = findByName(ctx, c, m, name)
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
