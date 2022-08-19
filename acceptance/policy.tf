
resource "alkira_group" "group1" {
  name = "tf-test-policy"
  description = "test policy"
}

resource "alkira_policy_rule" "test1" {
  name          = "tf-test-policy"
  description   = "Terraform Test Policy"
  src_ip        = "any"
  dst_ip        = "172.16.0.0/16"
  dscp          = "any"
  protocol      = "any"
  src_ports     = ["any"]
  dst_ports     = ["any"]
  rule_action   = "DROP"
}

resource "alkira_policy_rule_list" "test1" {
  name        = "tf-test-policy"
  description = "terraform test policy rule list"

  rules {
    priority = 1
    rule_id  = alkira_policy_rule.test1.id
  }
}

resource "alkira_policy" "tf_policy" {
  name         = "tf-test-policy"
  description  = "terraform test policy"
  enabled      = "false"
  from_groups  = ["-1"]
  to_groups    = ["-1"]
  rule_list_id = alkira_policy_rule_list.test1.id
  segment_ids  = [alkira_segment.seg1.id]
}
