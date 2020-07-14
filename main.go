package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/alkiranet/terraform-provider-alkira/alkira"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: alkira.Provider,
	})
}
