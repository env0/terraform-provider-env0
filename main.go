package main

import (
	"github.com/env0/terraform-provider-env0/env0"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: env0.Provider,
	})
}
