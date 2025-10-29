resource "alkira_connector_ipsec" "basic_dynamic" {
  name        = "ipsec-connector-basic"
  description = "Basic route-based IPSec connector with BGP"
  cxp         = "US-WEST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "SMALL"
  enabled     = true

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type                 = "DYNAMIC"
    customer_gateway_asn = "65310"
    availability         = "IPSEC_INTERFACE_PING"
  }

  endpoint {
    name                = "remote-site"
    customer_gateway_ip = "203.0.113.1"
    preshared_keys      = ["your-preshared-key-here"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
  }
}
