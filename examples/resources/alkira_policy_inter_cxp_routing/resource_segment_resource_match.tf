#
# Match by segment resource. Useful when route selection needs to align
# with segment-scoped resources already shared into the policy segment.
# Mutually exclusive with match_group_ids.
#
resource "alkira_policy_inter_cxp_routing" "segment_resource_match" {
  name        = "allow-shared-segment-resources"
  description = "Redistribute routes for specific segment resources only"
  enabled     = true
  direction   = "OUTBOUND"
  segment_id  = alkira_segment.seg1.id
  source_cxps = ["US-WEST-1"]
  dest_cxps   = ["US-EAST-1"]

  rule {
    name   = "allow-shared-resources"
    action = "ALLOW"
    match_segment_resource_ids = [
      alkira_segment_resource.shared_vpc.id,
      alkira_segment_resource.shared_vnet.id,
    ]
    set_community = "65512:100"
  }
}
