#
# Fan-out to multiple destination CXPs with prefix filtering. Only routes
# matching the allow-list prefix list are redistributed; the catch-all
# DENY rule at the end blocks everything else.
#
resource "alkira_policy_inter_cxp_routing" "prefix_filter" {
  name        = "us-west-hub-prefix-filter"
  description = "Allow only approved prefixes from US-WEST-1 to multiple CXPs"
  enabled     = true
  direction   = "OUTBOUND"
  segment_id  = alkira_segment.seg1.id
  source_cxps = ["US-WEST-1"]
  dest_cxps   = ["US-EAST-1", "EU-WEST-1", "AP-SOUTH-1"]

  rule {
    name                  = "allow-approved-prefixes"
    action                = "ALLOW"
    match_prefix_list_ids = [alkira_policy_prefix_list.approved.id]
  }

  rule {
    name      = "drop-everything-else"
    action    = "DENY"
    match_all = true
  }
}
