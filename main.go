package main

import (
	"context"
	"log"

	"github.com/alkiranet/terraform-provider-alkira/alkira"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	err := providerserver.Serve(context.Background(), alkira.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/providers/alkiranet",
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
