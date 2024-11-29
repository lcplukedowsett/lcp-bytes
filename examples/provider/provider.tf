# This example fetches a known order ID
terraform {
  required_providers {
    bytes = {
      version = "~> 0.2.5"
      source  = "lcp-llp/bytes"
    }
  }
}

provider "bytes" {
  username = "example"
  password = "example"
  identity_api_url = "https://example.com/identity"
  commerce_api_url = "https://example.com/commerce"
  contract_id = "12345"
}

data "bytes_order" "example" {
  order_id = "12345"
}