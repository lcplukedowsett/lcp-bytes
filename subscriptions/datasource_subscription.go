package subscriptions

import (
	"context"
	"fmt"

	"terraform-provider-bytes/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This datasource is used to get information about a known Bytes order
func datasourceOrder() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceOrderRead,

		// Initialise all vars for datasource
		Schema: map[string]*schema.Schema{
			"order_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Existing Bytes order ID",
			},
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Existing Bytes ID",
			},
			"contract_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Bytes contract name used for the order",
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Existing Subscription ID to query",
			},
			"friendly_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Friendly name of the subscription",
			},
			"po_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Purchase order number for the subscription order",
			},
			"create_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the order was created",
			},
		},
		Description: "Get information about a known Bytes order.\n\n" +
			"Use this data source to get information such as subscription name, id and creation date.",
	}
}

// datasourceOrderRead is used to read the datasource and set the schema
func datasourceOrderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	orderID := d.Get("order_id").(string)
	order, err := c.GetOrderDetails(orderID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get order with id %s: %s", orderID, err))
	}

	d.SetId(fmt.Sprintf("%d", order.ID))
	d.Set("id", order.ID)
	d.Set("contract_name", order.ContractName)
	d.Set("create_date", order.CreateDate)

	if len(order.Items) > 0 {
		d.Set("subscription_id", order.Items[0].SubscriptionID)
		d.Set("friendly_name", order.Items[0].FriendlyName)
		d.Set("po_number", order.Items[0].PONumber)
	}

	return nil
}
