#
# Disabled policy — kept in configuration but not enforced. Handy for
# staging a policy change without activating it.
#
resource "alkira_policy_inter_cxp_routing" "staged_disabled" {
  name        = "staged-policy"
  description = "Staged policy, disabled until cutover"
  enabled     = false
  direction   = "OUTBOUND"
  segment_id  = alkira_segment.seg1.id
  source_cxps = ["US-WEST-1"]
  dest_cxps   = ["US-EAST-1"]

  rule {
    name                  = "future-allow"
    action                = "ALLOW"
    match_prefix_list_ids = [alkira_policy_prefix_list.corp_prefixes.id]
  }
}
