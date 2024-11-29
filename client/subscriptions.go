package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SubscriptionDetails Struct for creating a subscription
type SubscriptionDetails struct {
	FriendlyName string
	PrincipalID  string
	PONumber     string
	BudgetCode   string
	DivisionID   int
}

// Basket Struct for Basket Post Request response
type BasketDetails struct {
	ID    int `json:"id"`
	Items []struct {
		ID          int    `json:"id"`
		PONumber    string `json:"poNumber"`
		PrincipalID string `json:"principalId"`
		BudgetCode  string `json:"budgetCode"`
	} `json:"items"`
}

// Checkout Struct for Checkout Post Request response
type Checkout struct {
	ID    int `json:"id"`
	Items []struct {
		ID           int    `json:"id"`
		PONumber     string `json:"poNumber"`
		FriendlyName string `json:"friendlyName"`
		PrincipalID  string `json:"principalId"`
	} `json:"items"`
}

type BasketPayload struct {
	Quantity     int    `json:"quantity"`
	FriendlyName string `json:"friendlyName"`
	ProductID    string `json:"productId"`
	SkuID        string `json:"skuId"`
	PrincipalID  string `json:"principalId"`
	PriceID      int    `json:"priceId"`
	PONumber     string `json:"poNumber"`
	BillingFreq  string `json:"billingFrequency"`
	Term         string `json:"term"`
	DivisionID   *int   `json:"divisionId"`
	BudgetCode   string `json:"budgetCode"`
}

func (c *Client) CreateBasket(friendlyName string, principalId string, poNumber string, budgetCode string) (*BasketDetails, error) {
	return c.createBasketHelper(friendlyName, principalId, poNumber, budgetCode, 0)
}

// CreateBasket creates a basket ready for checkout
func (c *Client) createBasketHelper(friendlyName string, principalId string, poNumber string, budgetCode string, retryCount int) (*BasketDetails, error) {
	url := fmt.Sprintf("%s/api/v2/contracts/%d/baskets", c.CommerceAPIURL, c.ContractID)

	// Create a payload for the request
	payload := BasketPayload{
		Quantity:     1,
		FriendlyName: friendlyName,
		ProductID:    "ENTITLEMENT",
		SkuID:        "ENTITLEMENT",
		PrincipalID:  principalId,
		PriceID:      24492277,
		PONumber:     poNumber,
		BillingFreq:  "monthly",
		Term:         "Perpetual",
		DivisionID:   nil,
		BudgetCode:   budgetCode,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %s", err)
	}

	// Submit the request with the payload
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.CustomToken))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Content-Length", "0")

	res, err := c.CustomHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %s", err)
	}
	defer res.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Print the response body for debugging
	fmt.Printf("Response Body: %s\n", string(bodyBytes))

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status code: %d, response body: %s", res.StatusCode, string(bodyBytes))
	}

	// Unmarshal the response body into a struct
	var basketdetails BasketDetails
	err = json.Unmarshal(bodyBytes, &basketdetails)
	if err != nil {
		return nil, err
	}

	// If there's more than 2 items in the basket, then clear the basket and try submitting the request again
	// This is due to the API creating a single basket until it's checked out or empty and prevents failures from creating multiple subscriptions
	itemsLength := len(basketdetails.Items)
	if itemsLength >= 2 {
		for _, item := range basketdetails.Items {
			deleteURL := fmt.Sprintf("%s/api/v1/CloudDashboard/DeleteBasketItem", c.CommerceAPIURL)
			fmt.Printf("Deleting Item: %d", item.ID)

			deletePayload := map[string]int{
				"basketItemId": item.ID,
			}

			deleteData, err := json.Marshal(deletePayload)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal delete payload: %s", err)
			}

			deleteReq, err := http.NewRequest("POST", deleteURL, bytes.NewBuffer(deleteData))
			if err != nil {
				return nil, fmt.Errorf("failed to create delete request: %s", err)
			}

			deleteReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.CustomToken))
			deleteReq.Header.Add("Content-Type", "application/json")

			deleteRes, err := c.CustomHTTPClient.Do(deleteReq)
			if err != nil || deleteRes.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("failed to delete basket item with id %d: %s", item.ID, err)
			}
		}

		// Check if we're below the maximum retry count (e.g., 3 retries)
		if retryCount < 3 {
			return c.createBasketHelper(friendlyName, principalId, poNumber, budgetCode, retryCount+1)
		} else {
			return nil, fmt.Errorf("max retries reached while trying to create basket")
		}
	}

	return &basketdetails, nil
}

// Checkout basket from previous CreateBasket Step
func (c *Client) CheckoutBasket(basketdetails *BasketDetails) (*Checkout, error) {
	url := fmt.Sprintf("%s/api/v2/contracts/%d/baskets/%d/checkout", c.CommerceAPIURL, c.ContractID, basketdetails.ID)

	// Submit the request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.CustomToken))

	res, err := c.CustomHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %s", err)
	}
	defer res.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Print the response body for debugging
	fmt.Printf("Response Body: %s\n", string(bodyBytes))

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("unexpected HTTP status code: %d, response body: %s", res.StatusCode, string(bodyBytes))
	}

	// Unmarshal the response body into a struct
	var checkout Checkout
	err = json.Unmarshal(bodyBytes, &checkout)
	if err != nil {
		return nil, err
	}

	return &checkout, nil
}

// CreateSubscription creates a subscription by first creating a basket and then proceeding to checkout
func (c *Client) CreateSubscription(details SubscriptionDetails) (*OrderDetails, error) {
	// First, create a basket
	basketInfo, err := c.CreateBasket(details.FriendlyName, details.PrincipalID, details.PONumber, details.BudgetCode)
	if err != nil {
		return nil, fmt.Errorf("failed to create basket: %s", err)
	}

	// Then, checkout the basket using the returned ID
	checkoutInfo, err := c.CheckoutBasket(basketInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to checkout basket: %s", err)
	}

	// Wait for a maximum of 20 minutes for subscriptionId to be not null.
	// Check every 30 seconds
	maxRetries := 40
	retryInterval := 30 * time.Second
	var subscriptionInfo *OrderDetails
	for i := 0; i < maxRetries; i++ {
		// Lastly, check the status of the order using the Checkout ID
		subscriptionInfo, err = c.GetOrderDetails(fmt.Sprintf("%d", checkoutInfo.ID))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch order details: %s", err)
		}

		// Check if subscriptionId is not null for the first item
		if subscriptionInfo.Items[0].SubscriptionID != "" {
			fmt.Printf("SubscriptionId: %s\n", subscriptionInfo.Items[0].SubscriptionID)
			break
		}
		fmt.Printf("SubscriptionId is blank, waiting...")
		// Wait for the retry interval before the next check
		time.Sleep(retryInterval)
	}

	if subscriptionInfo.Items[0].SubscriptionID == "" {
		return nil, fmt.Errorf("subscriptionId did not update after waiting for the maximum allowed time")
	}

	return subscriptionInfo, nil
}
