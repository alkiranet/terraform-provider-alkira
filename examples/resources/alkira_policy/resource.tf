#
# A simple policy was constructed with segment, policy_rule and
# policy_rule_list.
#
resource "alkira_segment" "test" {
  name  = "test-segment"
  asn   = "65513"
  cidrs = ["10.16.1.0/24"]
}

resource "alkira_policy_rule" "test" {
  name          = "test-rule"
  description   = "Test Rule"
  src_ip        = "any"
  dst_ip        = "172.16.0.0/16"
  dscp          = "any"
  protocol      = "any"
  src_ports     = ["any"]
  dst_ports     = ["any"]
  rule_action   = "DROP"
}

resource "alkira_policy_rule_list" "test" {
  name        = "test-rule-list"
  description = "test policy rule list"

  rules {
    priority = 1
    rule_id  = alkira_policy_rule.test.id
  }
}

resource "alkira_policy" "test" {
  name         = "test-policy"
  description  = "test policy"
  enabled      = "false"
  from_groups  = ["-1"]
  to_groups    = ["-1"]
  rule_list_id = alkira_policy_rule_list.test.id
  segment_ids  = [alkira_segment.test.id]
}
