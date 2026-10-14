resource "alkira_multi_cloud_gateway" "example" {
  name        = "example-gateway"
  tier        = "10GBPS"
  description = "Multi-Cloud Gateway managed by Terraform"

  timeouts {
    create = "30m"
  }
}
