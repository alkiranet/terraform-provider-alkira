resource "alkira_peering_gateway_cxp" "test1" {
  name         = "tf-test-1"
  description  = "Test CXP Peering Gatewat"
  cloud_region = "useast"
  cxp          = "US_EAST_1"
  segment      = "prod-seg"
}
