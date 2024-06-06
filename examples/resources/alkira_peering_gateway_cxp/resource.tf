resource "alkira_peering_gateway_cxp" "test1" {
  name         = "tf-test-1"
  description  = "Test CXP Peering Gatewat"
  cloud_region = "us-east-1"
  cxp          = "US_EAST_1"
  segment      = "prod-seg"
}
