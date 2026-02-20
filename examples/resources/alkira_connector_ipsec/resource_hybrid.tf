resource "alkira_connector_ipsec" "hybrid_routing" {
  name        = "ipsec-connector-hybrid"
  description = "IPSec connector with both static and dynamic routing"
  cxp         = "US-WEST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "MEDIUM"
  enabled     = true

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type                 = "BOTH"
    prefix_list_id       = alkira_policy_prefix_list.remote_subnets.id
    customer_gateway_asn = "65320"
    availability         = "IPSEC_INTERFACE_PING"
  }

  endpoint {
    name                = "hybrid-site"
    customer_gateway_ip = "203.0.113.20"
    preshared_keys      = ["hybrid-key"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
  }
}
