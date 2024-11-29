package subscriptions

import (
	"context"
	"fmt"

	"terraform-provider-bytes/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider - Initialize all vars for Provider config
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"identity_api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BYTES_IDENTITY_HOST", nil),
				Description: "The identity API URL provided by the host",
			},
			"commerce_api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BYTES_COMMERCE_HOST", nil),
				Description: "The commerce API URL provided by the host",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BYTES_USERNAME", nil),
				Description: "Username used for authentication to API Endpoints",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("BYTES_PASSWORD", nil),
				Description: "Password used for authentication to API Endpoints",
			},
			"contract_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("BYTES_CONTRACT_ID", nil),
				Description: "Contract ID used for authentication to API Endpoints",
			},
		},
		// Define the function to call the resource.
		ResourcesMap: map[string]*schema.Resource{
			"bytes_subscription": resourceSubscription(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"bytes_order": datasourceOrder(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure - Configure Provider
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	// Get credentials and prepare it for provider
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	// Get contract ID
	contract_id := d.Get("contract_id").(int)

	// Prepare URL variable for Customclient
	var identity_api_url *string
	hVal, ok := d.GetOk("identity_api_url")
	if ok {
		tempHost := hVal.(string)
		identity_api_url = &tempHost
	}
	fmt.Printf("Identity API URL is: %s", *identity_api_url)

	var commerce_api_url *string
	cVal, ok := d.GetOk("commerce_api_url")
	if ok {
		tempCommerce := cVal.(string)
		commerce_api_url = &tempCommerce
	}
	fmt.Printf("Commerce API URL is: %s", *commerce_api_url)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// If all values are provided then create client
	if (username != "") && (password != "") {
		c, err := client.NewClient(identity_api_url, commerce_api_url, &username, &password, contract_id)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create RestApi Client",
				Detail:   fmt.Sprintf("Something wrong with Provider to create client. Error: %s", err),
			})

			return nil, diags
		}

		return c, diags
	}

	// if values are missing, then create client and return the response
	c, err := client.NewClient(nil, nil, nil, nil, 0)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create RestApi Client",
			Detail:   fmt.Sprintf("Unable to authenticate user for authenticated RestApi client: %s", err),
		})
		return nil, diags
	}
	return c, diags
}
