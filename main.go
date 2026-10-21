package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/appgate/terraform-provider-appgatesdp/internal/provider"
	"github.com/appgate/terraform-provider-appgatesdp/version"
)

//go:generate go run ./cmd/sdpgen generate

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), provider.New(version.ProviderVersion), providerserver.ServeOpts{
		Address: "registry.terraform.io/appgate/appgatesdp",
		Debug:   debug,
	})
	if err != nil {
		log.Fatal(err)
	}
}
