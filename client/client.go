package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client Configuration-
type Client struct {
	CommerceAPIURL   string
	CustomHostURL    string
	CustomHTTPClient *http.Client
	CustomToken      string
	CustomAuth       CustomAuthStruct
	ContractID       int
}

// CustomHostURL  Default URL, empty one
const CustomHostURL string = ""

// CustomAuthStruct Auth Credentials -
type CustomAuthStruct struct {
	Username string `url:"client_id"`
	Password string `url:"client_secret"`
}

// CustomAuthResponse Auth Response -
type CustomAuthResponse struct {
	Token     string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
}

type OrderDetails struct {
	ID           int    `json:"id"`
	ContractName string `json:"contractName"`
	CreateDate   string `json:"createDate"`
	Items        []struct {
		SubscriptionID      string `json:"subscriptionId"`
		PONumber            string `json:"poNumber"`
		FriendlyName        string `json:"friendlyName"`
		PrincipalID         string `json:"principalId"`
		CloudSubscriptionID *int   `json:"cloudSubscriptionId"`
	} `json:"items"`
}

// NewClient Initialize a new Client
func NewClient(identity_api_url, commerce_api_url, username, password *string, contract_id int) (*Client, error) {

	// Initialize the client
	client := Client{
		CustomHTTPClient: &http.Client{Timeout: 120 * time.Second},
		// Set Default URLs
		CustomHostURL: CustomHostURL,
	}

	// Add CommerceAPIURL for Custom client
	if commerce_api_url != nil {
		client.CommerceAPIURL = *commerce_api_url
	}

	// Add Host URL for Custom client
	if identity_api_url != nil {
		client.CustomHostURL = *identity_api_url
	}

	// Add Contract ID for Custom client
	client.ContractID = contract_id

	// If necessary variables are empty then return and error
	if username == nil || password == nil || contract_id == 0 {
		return &client, nil
	}

	// Prepare struct for Cloud Auth Credentials
	client.CustomAuth = CustomAuthStruct{
		Username: *username,
		Password: *password,
	}

	// Get The token for Custom Client
	ctt, err := client.GetCustomClientToken()
	if err != nil {
		return nil, err
	}

	// Set the received token to the CustomClient
	client.CustomToken = ctt.Token

	// Return the client
	return &client, nil
}

// GetOrderDetails fetches the details of an order
func (c *Client) GetOrderDetails(orderID string) (*OrderDetails, error) {
	url := fmt.Sprintf("%s/api/v2/contracts/%d/orders/%s", c.CommerceAPIURL, c.ContractID, orderID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.CustomToken))
	res, err := c.CustomHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("unexpected HTTP status code: %d, response body: %s", res.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var order OrderDetails
	err = json.Unmarshal(body, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
