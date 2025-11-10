resource "alkira_connector_ipsec" "ha_active_standby" {
  name        = "ipsec-connector-ha"
  description = "High availability IPSec with active and standby endpoints"
  cxp         = "US-WEST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "MEDIUM"
  enabled     = true

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type                 = "DYNAMIC"
    customer_gateway_asn = "65340"
  }

  endpoint {
    name                     = "primary-active"
    customer_gateway_ip      = "203.0.113.50"
    preshared_keys           = ["primary-key-1", "primary-key-2"]
    billing_tag_ids          = [alkira_billing_tag.tag1.id]
    ha_mode                  = "ACTIVE"
    enable_tunnel_redundancy = true
  }

  endpoint {
    name                     = "secondary-active"
    customer_gateway_ip      = "203.0.113.51"
    preshared_keys           = ["secondary-key-1", "secondary-key-2"]
    billing_tag_ids          = [alkira_billing_tag.tag1.id]
    ha_mode                  = "ACTIVE"
    enable_tunnel_redundancy = true
  }

  endpoint {
    name                = "standby"
    customer_gateway_ip = "203.0.113.52"
    preshared_keys      = ["standby-key"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
    ha_mode             = "STANDBY"
  }
}
