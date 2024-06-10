resource "alkira_peering_gateway_cxp" "example-cxp-gateway" {
  name         = "example-cxp-gateway"
  description  = "Example CXP Peering Gateway"
  cloud_region = "useast"
  cxp          = "US-EAST-1"
  segment_id   = alkira_segment.example-segment.id
}
