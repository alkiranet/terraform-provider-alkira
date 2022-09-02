#
# A simple NAT policy was constructed with segment, policy_nat_rule
# and policy_nat.
#
resource "alkira_segment" "test" {
  name  = "test-segment"
  asn   = "65513"
  cidrs = ["10.16.1.0/24"]
}

resource "alkira_policy_nat_rule" "test" {
  name          = "test"
  description   = "test nat rule"
  enabled       = false

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

resource "alkira_policy_nat" "test" {
  name               = "tftest"
  description        = "terraform test NAT policy"
  type               = "INTRA_SEGMENT"
  segment_id         = alkira_segment.test.id
  included_group_ids = [alkira_group.test.id]
  nat_rule_ids       = [alkira_policy_nat_rule.test.id]
}
