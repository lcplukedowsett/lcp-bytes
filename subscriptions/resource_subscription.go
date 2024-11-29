package subscriptions

import (
	"context"
	"fmt"
	"terraform-provider-bytes/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This resource is used to create a new Azure subscription
func resourceSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubscriptionCreate,
		ReadContext:   resourceSubscriptionRead,
		DeleteContext: resourceSubscriptionDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID assigned by Bytes to the order",
			},
			"friendly_name": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Friendly name of the subscription to create. This is used as the name of the subscription in the Bytes/Azure Portal",
			},
			"po_number": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The PO number which can be used to assign a cost to a purchase for billing purposes",
			},
			"default_admin": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The default admin which is assigned to a newly created subscription",
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The automatically generated subscription ID returned by Azure",
			},
			"budget_code": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The budget code to use for subscription billing",
			},
			"division_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The division ID to use for subscription billing",
			},
		},
		Description: "Creates a new Azure subscription.\n\n" +
			"This resources is intended to be used to create a new Azure subscription",
	}
}

func resourceSubscriptionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	c := m.(*client.Client)

	// Prepare the JSON Data for API Payload
	subscriptionDetails := client.SubscriptionDetails{
		FriendlyName: d.Get("friendly_name").(string),
		PONumber:     d.Get("po_number").(string),
		PrincipalID:  d.Get("default_admin").(string),
		BudgetCode:   d.Get("budget_code").(string),
		DivisionID:   d.Get("division_id").(int),
	}

	// Call the function create the subscription with payload
	subscription, err := c.CreateSubscription(subscriptionDetails)

	// If error, print it out
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Something wrong with Provider to Delete record: %s", err),
			Detail:   fmt.Sprintf("Something wrong with Provider to create record: %s", err),
		})
		return diags
	}

	d.SetId(fmt.Sprintf("%d", subscription.ID))
	d.Set("contract_name", subscription.ContractName)
	d.Set("create_date", subscription.CreateDate)

	if len(subscription.Items) > 0 {
		d.Set("subscription_id", subscription.Items[0].SubscriptionID)
		d.Set("friendly_name", subscription.Items[0].FriendlyName)
		d.Set("po_number", subscription.Items[0].PONumber)
		d.Set("default_admin", subscription.Items[0].PrincipalID)
	}
	// Call the resourceSubscriptionRead function
	resourceSubscriptionRead(ctx, d, m)

	return nil
}
func resourceSubscriptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	return diags
}
func resourceSubscriptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// No-op, do nothing when deleting, not currently supported by Bytes API
	return nil
}
