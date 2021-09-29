resource "alkira_policy_rule" "rule" {
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
