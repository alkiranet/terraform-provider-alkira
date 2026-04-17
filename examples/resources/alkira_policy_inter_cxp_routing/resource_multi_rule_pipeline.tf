#
# Multi-rule pipeline with mixed actions. Rules are evaluated top-down
# (first match wins) using system-assigned sequence numbers starting at
# 1000. Order your rules from most specific to most general.
#
resource "alkira_policy_inter_cxp_routing" "multi_rule_pipeline" {
  name        = "tiered-redistribution"
  description = "Tiered inter-CXP redistribution with explicit precedence"
  enabled     = true
  direction   = "OUTBOUND"
  segment_id  = alkira_segment.seg1.id
  source_cxps = ["US-WEST-1"]
  dest_cxps   = ["US-EAST-1", "EU-WEST-1"]

  # Highest precedence: drop blacklisted prefixes.
  rule {
    name                  = "drop-blacklist"
    action                = "DENY"
    match_prefix_list_ids = [alkira_policy_prefix_list.blacklist.id]
  }

  # Allow tagged prod routes with AS-path prepend.
  rule {
    name                     = "allow-prod-with-prepend"
    action                   = "ALLOW"
    match_community_list_ids = [alkira_list_community.prod_inbound.id]
    set_as_path_prepend      = "65100"
  }

  # Allow everything else from approved prefix list untouched.
  rule {
    name                  = "allow-approved"
    action                = "ALLOW"
    match_prefix_list_ids = [alkira_policy_prefix_list.approved.id]
  }

  # Default deny.
  rule {
    name      = "default-deny"
    action    = "DENY"
    match_all = true
  }
}
