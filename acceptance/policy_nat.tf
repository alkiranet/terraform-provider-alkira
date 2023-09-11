resource "alkira_policy_nat_rule" "test1" {
  name        = "acceptance-basic"
  description = "acceptance basic NAT rule"
  enabled     = false

  match {
    src_prefixes = ["any"]
    dst_prefixes = ["any"]
    protocol     = "any"
  }

  action {
    src_addr_translation_type = "NONE"
    dst_addr_translation_type = "NONE"
  }
}

resource "alkira_policy_nat" "test1" {
  name               = "acceptance-basic"
  description        = "terraform test NAT policy"
  type               = "INTRA_SEGMENT"
  segment_id         = alkira_segment.test1.id
  included_group_ids = [alkira_group.test5.id]
  nat_rule_ids       = [alkira_policy_nat_rule.test1.id]
}
