#
# Extended community matching with Site-of-Origin (SOO) tagging. Useful
# for preventing route loops in multi-CXP topologies.
#
resource "alkira_policy_inter_cxp_routing" "ext_community_soo" {
  name        = "soo-tagging-us-west"
  description = "Stamp SOO on routes leaving US-WEST-1"
  enabled     = true
  direction   = "OUTBOUND"
  segment_id  = alkira_segment.seg1.id
  source_cxps = ["US-WEST-1"]
  dest_cxps   = ["US-EAST-1"]

  rule {
    name                              = "stamp-soo"
    action                            = "ALLOW"
    match_extended_community_list_ids = [alkira_list_extended_community.internal.id]
    set_extended_community            = "soo:65512:21 soo:10.1.1.1:1234"
  }
}
