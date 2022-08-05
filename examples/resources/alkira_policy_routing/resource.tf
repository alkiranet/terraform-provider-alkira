#
# Assume segment1 and group1 has been created
#
resource "alkira_policy_routing" "test" {
  name                = "minimal"
  description         = "minimal routing policy"
  enabled             = "false"
  direction           = "INBOUND"
  included_group_ids  = [alkira_group.group1.id]
  segment_id          = alkira_segment.seg1.id

  rule {
    name   = "test"
    action = "ALLOW"

    match {
      all = true
    }

    inter_cxp_routes_redistribution {
      distribution_type = "ALL"
    }
  }
}
