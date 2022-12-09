resource "alkira_policy_routing" "test1" {
  name               = "acceptance-inbound"
  description        = "aceeptance test of inbound routing policy"
  enabled            = "false"
  direction          = "INBOUND"
  included_group_ids = [alkira_group.test2.id]
  segment_id         = alkira_segment.test2.id

  rule {
    name                     = "test-rule-1"
    action                   = "ALLOW"
    match_all                = true
    routes_distribution_type = "ALL"
  }
}

resource "alkira_policy_routing" "test2" {
  name               = "acceptance-outbound"
  enabled            = false
  direction          = "OUTBOUND"
  included_group_ids = [alkira_group.test2.id]
  segment_id         = alkira_segment.test2.id

  rule {
    name            = "test-rule-2"
    action          = "ALLOW"
    match_all       = false
    match_group_ids = [alkira_group.test2.id]

    set_as_path_prepend = 65300
  }
}
