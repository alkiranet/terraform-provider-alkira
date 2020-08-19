package main

import (
	"github.com/alkiranet/terraform-provider-alkira/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: alkira.Provider,
	})
}
