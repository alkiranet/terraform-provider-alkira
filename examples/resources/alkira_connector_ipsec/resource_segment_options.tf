resource "alkira_connector_ipsec" "segment_options" {
  name        = "ipsec-connector-segment-opts"
  description = "IPSec connector with segment-specific options"
  cxp         = "US-WEST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "SMALL"
  enabled     = true

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type                 = "DYNAMIC"
    customer_gateway_asn = "65370"
  }

  segment_options {
    name                     = alkira_segment.segment1.name
    advertise_default_route  = true
    advertise_on_prem_routes = true
  }

  endpoint {
    name                = "segment-site"
    customer_gateway_ip = "203.0.113.70"
    preshared_keys      = ["segment-key"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
  }
}
