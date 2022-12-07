# Create resources under Policies/

resource "alkira_policy" "test" {
  name         = "acceptance-test"
  description  = "terraform test policy"
  enabled      = "false"
  from_groups  = ["-1"]
  to_groups    = ["-1"]
  rule_list_id = alkira_policy_rule_list.test.id
  segment_ids  = [alkira_segment.test1.id]
}

resource "alkira_policy_rule" "test" {
  name        = "acceptance-test"
  description = "Terraform Test Policy"
  src_ip      = "any"
  dst_ip      = "172.16.0.0/16"
  dscp        = "any"
  protocol    = "any"
  src_ports   = ["any"]
  dst_ports   = ["any"]
  rule_action = "DROP"
}

resource "alkira_policy_rule_list" "test" {
  name        = "acceptance-test"
  description = "terraform test policy rule list"

  rules {
    priority = 1
    rule_id  = alkira_policy_rule.test.id
  }
}

resource "alkira_policy_nat_rule" "test" {
  name        = "acceptance-test"
  description = "tftest nat rule"
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
