resource "alkira_connector_ipsec" "ipsec" {
  name           = "connector-test-ipsec"
  cxp            = "US-WEST"
  group          = alkira_group.group1.name
  segment_id     = alkira_segment.segment1.id
  size           = "SMALL"

  vpn_mode       = "ROUTE_BASED"

  routing_options {
    type = "DYNAMIC"
    customer_gateway_asn = "65310"
  }
}
