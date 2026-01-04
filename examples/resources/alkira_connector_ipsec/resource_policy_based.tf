resource "alkira_connector_ipsec" "policy_based" {
  name        = "ipsec-connector-policy-based"
  description = "Policy-based IPSec connector with traffic selectors"
  cxp         = "US-WEST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "SMALL"
  enabled     = true

  vpn_mode = "POLICY_BASED"

  policy_options {
    on_prem_prefix_list_ids = [alkira_list_global_cidr.on_prem_subnets.id]
    cxp_prefix_list_ids     = [alkira_list_global_cidr.cxp_subnets.id]
  }

  endpoint {
    name                = "policy-site"
    customer_gateway_ip = "203.0.113.40"
    preshared_keys      = ["policy-key"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
  }
}
