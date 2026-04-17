#
# Match by connector group. Only routes learned from connectors in the
# specified groups are redistributed. match_group_ids is mutually
# exclusive with match_segment_resource_ids.
#
resource "alkira_policy_inter_cxp_routing" "group_match" {
  name        = "allow-aws-connectors-only"
  description = "Only redistribute routes learned from AWS connector groups"
  enabled     = true
  direction   = "OUTBOUND"
  segment_id  = alkira_segment.seg1.id
  source_cxps = ["US-WEST-1"]
  dest_cxps   = ["US-EAST-1"]

  rule {
    name   = "allow-aws-groups"
    action = "ALLOW"
    match_group_ids = [
      alkira_group.aws_prod.id,
      alkira_group.aws_staging.id,
    ]
  }

  rule {
    name      = "drop-others"
    action    = "DENY"
    match_all = true
  }
}
