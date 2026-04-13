#
# Drop all routes — fully block redistribution from source to destination
# CXP while keeping the policy object in place (useful for temporarily
# quarantining a CXP pair without deleting the policy).
#
resource "alkira_policy_inter_cxp_routing" "drop_all" {
  name        = "block-us-west-to-eu-west"
  description = "Temporarily block inter-CXP routes from US-WEST-1 to EU-WEST-1"
  enabled     = true
  direction   = "OUTBOUND"
  segment_id  = alkira_segment.seg1.id
  source_cxps = ["US-WEST-1"]
  dest_cxps   = ["EU-WEST-1"]

  rule {
    name      = "drop-all"
    action    = "DENY"
    match_all = true
  }
}
