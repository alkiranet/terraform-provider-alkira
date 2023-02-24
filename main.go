package main

import (
	"context"
	"flag"
	"log"

	"github.com/alkiranet/terraform-provider-alkira/alkira"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	err := providerserver.Serve(context.Background(), alkira.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/providers/alkiranet",
		Debug:   debug,
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
