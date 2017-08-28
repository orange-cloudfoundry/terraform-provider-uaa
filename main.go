package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-uaa/uaa"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: uaa.Provider})

}
