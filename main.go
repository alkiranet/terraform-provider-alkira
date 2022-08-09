package main

import (
    "context"
    "flag"
    "log"

	"github.com/alkiranet/terraform-provider-alkira/alkira"
    "github.com/hashicorp/terraform-plugin-go/tfprotov5"
    "github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
    "github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

var (
    // Version can be updated by goreleaser on release
    version string = "dev"
)

func main() {
    debugFlag := flag.Bool("debug", false, "Start provider in debug mode.")
    flag.Parse()

    ctx := context.Background()
    providers := []func() tfprotov5.ProviderServer{
        alkira.Provider.New(version).GRPCProvider,
    }

    muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)

    if err != nil {
        log.Fatal(err)
    }

    var serveOpts []tf5server.ServeOpt

    if *debugFlag {
        serveOpts = append(serveOpts, tf5server.WithManagedDebug())
    }

    err = tf5server.Serve(
        "registry.terraform.io/alkiranet/alkira",
        muxServer.ProviderServer,
        serveOpts...,
    )

    if err != nil {
        log.Fatal(err)
    }
}
