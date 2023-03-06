//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/appgate/sdp-api-client-go/api/v18/openapi"
)

type Resource struct {
	Name, Service, Model, ServiceGetMethod, ServiceIDGetMethod string
}

type templateStub struct {
	Imports  []string
	Resource []Resource
}

var (
	verbose = flag.Bool("v", false, "Print verbose log messages")

	stub       = templateStub{}
	generators = []Resource{
		{
			Name: "Entitlement",
		},
		{
			Name:    "AdministrativeRole",
			Service: "AdminRolesApi",
		},
		{
			Name: "ApplianceCustomization",
		},
		{
			Name: "Appliance",
		},
		{
			Name: "Condition",
		},
	}
)

func logf(fmt string, args ...interface{}) {
	if *verbose {
		log.Printf(fmt, args...)
	}
}

func main() {
	flag.Parse()

	client := openapi.APIClient{}
	val := reflect.ValueOf(client)
	stub.Imports = append(stub.Imports, val.Type().PkgPath())
	t := reflect.TypeOf(client)

	for i := 0; i < t.NumField(); i++ {
		for k, generator := range generators {
			plural := generator.Name + "s"
			// If we have already defined a service, its likely a alias
			guess := plural + "Api"
			if len(generator.Service) > 0 {
				guess = generator.Service
			}

			if guess == t.Field(i).Name {
				child := t.Field(i)
				generator.Service = fmt.Sprintf("%s", child.Type.Elem())

				// TODO get reflect | go analysis to get the exact method name and return value
				generator.ServiceGetMethod = fmt.Sprintf("%sGet", plural)
				generator.ServiceIDGetMethod = fmt.Sprintf("%sIdGet", plural)
				generator.Model = fmt.Sprintf("openapi.%s", generator.Name)
				generators[k] = generator

			}
		}
	}
	// sanity check; make sure all generators has a Method
	for _, generator := range generators {
		if len(generator.Service) < 1 {
			fmt.Printf("Name: %s\nService: %s\nnServiceGetMethod: %s\nServiceIDGetMethod: %s\nModel: %s",
				generator.Name,
				generator.Service,
				generator.ServiceGetMethod,
				generator.ServiceIDGetMethod,
				generator.Model,
			)
			die(fmt.Errorf("generator %s did not get correctly mapped", generator.Name))
		}
	}
	stub.Resource = generators

	f, err := os.Create("appgate/find_resource_by_name.go")
	if err != nil {
		die(err)
	}
	defer f.Close()

	funcs := map[string]any{
		"Title":     strings.Title,
		"Lowercase": strings.ToLower,
	}

	goTemplate, err := template.New("").Funcs(funcs).Parse(packageTemplate)
	if err != nil {
		die(fmt.Errorf("template New err %w", err))
	}
	var buf bytes.Buffer
	if err := goTemplate.Execute(&buf, stub); err != nil {
		die(fmt.Errorf("template err %w", err))
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		die(fmt.Errorf("format err %w", err))
	}
	if _, err := f.Write(p); err != nil {
		die(fmt.Errorf("write err %w", err))
	}
	logf("Done")
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

const packageTemplate = `// Code generated by go generate; DO NOT EDIT.
package appgate

import (
	"context"
	"log"

	{{- range .Imports }}
	"{{ . }}"
	{{- end }}
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)


{{- range .Resource }}


func find{{ .Name | Title }}ByUUID(ctx context.Context, api *{{ .Service }}, id, token string) (*{{ .Model }}, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source {{ .Name }} get by UUID %s", id)
	resource, _, err := api.{{ .ServiceIDGetMethod }}(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func find{{ .Name | Title }}ByName(ctx context.Context, api *{{ .Service }}, name, token string) (*{{ .Model }}, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source {{ .Name }} get by name %s", name)

	resource, _, err := api.{{ .ServiceGetMethod }}(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple {{ .Name }} matched; use additional constraints to reduce matches to a single {{ .Name }}")
	}
	for _, r := range resource.GetData() {
		return &r, nil
	}
	return nil, AppendErrorf(diags, "Failed to find {{ .Name }} %s", name)
}


func Resolve{{ .Name | Title}}FromResourceData(ctx context.Context, d *schema.ResourceData, api *{{ .Service }}, token string) (*{{ .Model }}, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("{{ .Name | Lowercase }}_id")
	resourceName, nok := d.GetOk("{{ .Name | Lowercase }}_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of {{ .Name | Lowercase }}_id or {{ .Name | Lowercase }}_name attributes")
	}

	if iok {
		return find{{ .Name | Title}}ByUUID(ctx, api, resourceID.(string), token)
	}
	return find{{ .Name | Title}}ByName(ctx, api, resourceName.(string), token)
}

{{- end }}

`
