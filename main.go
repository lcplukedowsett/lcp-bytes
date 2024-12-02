package main

import (
	"terraform-provider-bytesnew/subscriptions"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// Main function, calling the provider
func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return subscriptions.Provider()
		},
	})
}
