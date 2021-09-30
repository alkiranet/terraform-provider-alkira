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
