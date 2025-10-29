resource "alkira_connector_ipsec" "bgp_auth" {
  name        = "ipsec-connector-bgp-auth"
  description = "IPSec connector with BGP MD5 authentication"
  cxp         = "US-WEST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "SMALL"
  enabled     = true

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type                 = "DYNAMIC"
    customer_gateway_asn = "65330"
    bgp_auth_key         = "my-bgp-secret-key"
    availability         = "PING"
  }

  endpoint {
    name                = "secured-site"
    customer_gateway_ip = "203.0.113.30"
    preshared_keys      = ["secured-psk"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
  }
}
