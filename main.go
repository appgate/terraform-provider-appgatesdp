package main

import (
	"github.com/appgate/terraform-provider-appgate/appgate"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: appgate.Provider,
	})
}
