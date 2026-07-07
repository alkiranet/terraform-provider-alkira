#
# Community-based matching with community tagging. Routes matching the
# inbound community list are re-tagged with new community values before
# being redistributed to the destination CXP.
#
resource "alkira_policy_inter_cxp_routing" "community_tagging" {
  name        = "tag-prod-routes"
  description = "Tag production routes crossing CXP boundaries"
  enabled     = true
  direction   = "OUTBOUND"
  segment_id  = alkira_segment.seg1.id
  source_cxps = ["US-WEST-1"]
  dest_cxps   = ["US-EAST-1"]

  rule {
    name                     = "tag-with-prod-community"
    action                   = "ALLOW"
    match_community_list_ids = [alkira_list_community.prod_inbound.id]
    set_community            = "65512:20 65512:21"
  }
}
