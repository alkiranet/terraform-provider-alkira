resource "alkira_connector_ipsec" "route_based_static" {
  name        = "ipsec-connector-static"
  description = "Route-based IPSec connector with static routes"
  cxp         = "US-EAST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "SMALL"
  enabled     = true

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type           = "STATIC"
    prefix_list_id = alkira_list_global_cidr.remote_subnets.id
    availability   = "IKE_STATUS"
  }

  endpoint {
    name                = "branch-office"
    customer_gateway_ip = "203.0.113.10"
    preshared_keys      = ["branch-key-1", "branch-key-2"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
  }
}
