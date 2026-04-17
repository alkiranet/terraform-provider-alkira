#
# AS-path list matching. Filters routes originating from or transiting
# specific AS numbers, then DENYs them — commonly used to keep
# partner/peer routes from leaking across CXPs.
#
resource "alkira_policy_inter_cxp_routing" "as_path_filter" {
  name        = "drop-partner-as-routes"
  description = "Drop routes transiting specific partner AS numbers"
  enabled     = true
  direction   = "OUTBOUND"
  segment_id  = alkira_segment.seg1.id
  source_cxps = ["US-WEST-1"]
  dest_cxps   = ["US-EAST-1"]

  rule {
    name                   = "drop-partner-as"
    action                 = "DENY"
    match_as_path_list_ids = [alkira_list_as_path.partner_as.id]
  }

  rule {
    name      = "allow-rest"
    action    = "ALLOW"
    match_all = true
  }
}
