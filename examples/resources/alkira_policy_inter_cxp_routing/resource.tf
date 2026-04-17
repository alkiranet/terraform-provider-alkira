#
# Minimal policy — allow all routes from one CXP to another.
# Assume segment1 has been created and CXPs US-WEST-1 and US-EAST-1 exist.
#
resource "alkira_policy_inter_cxp_routing" "allow_all" {
  name        = "us-west-to-us-east"
  description = "Outbound inter-CXP policy from US-WEST-1 to US-EAST-1"
  enabled     = true
  direction   = "OUTBOUND"
  segment_id  = alkira_segment.seg1.id
  source_cxps = ["US-WEST-1"]
  dest_cxps   = ["US-EAST-1"]

  rule {
    name      = "allow-all"
    action    = "ALLOW"
    match_all = true
  }
}
