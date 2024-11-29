package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Get Bearer Token for CustomAPI
func (c *Client) GetCustomClientToken() (*CustomAuthResponse, error) {
	// Check if credentials provided are empty
	if c.CustomAuth.Username == "" || c.CustomAuth.Password == "" {
		return nil, fmt.Errorf("define CustomAPI username and password")
	}

	// Convert Credentials Struct to URL-encoded form
	data := url.Values{}
	data.Set("client_id", c.CustomAuth.Username)
	data.Set("client_secret", c.CustomAuth.Password)
	data.Set("grant_type", "client_credentials")

	// Prepare Client for Getting the Token
	url := fmt.Sprintf("%s/api/v1/oauth/token", c.CustomHostURL)
	method := "POST"
	payload := strings.NewReader(data.Encode()) // Convert url.Values to string

	// Initialize the Request client
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return nil, err
	}

	// Add Headers for the request
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	fmt.Println("payload is: ", payload)
	fmt.Println("URL: ", url)
	fmt.Println("Request Headers: ", req.Header)

	// Make the request
	res, err := c.CustomHTTPClient.Do(req)
	if err != nil {
		fmt.Println("Error executing HTTP request:", err)
		return nil, err
	}
	defer res.Body.Close()

	// Read the Body from the response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	// Check status code
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	// Prepare the response as a struct
	ctAr := CustomAuthResponse{}

	fmt.Println("Response Status Code:", res.StatusCode)
	fmt.Println("Response Headers:", res.Header)
	fmt.Printf("Response Body: %s\n", body)
	// Convert Json response to struct
	err = json.Unmarshal(body, &ctAr)
	if err != nil {
		fmt.Println("Error unmarshalling JSON response:", err)
		return nil, err
	}

	// Return the Struct
	return &ctAr, nil
}
