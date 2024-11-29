# Create a new subscription
resource "bytes_subscription" "example" {
  friendly_name = "examplesub"
  po_number = "13102023-example"
  default_admin = "username@domain.uk.com"
  budget_code = "12345"
}