resource "alkira_multi_cloud_gateway" "example" {
  name                   = "example-gateway"
  tier                   = "10GBPS"
  description            = "Multi-Cloud Gateway managed by Terraform"
  billing_account_number = "1-1ABCDE-F"
}
