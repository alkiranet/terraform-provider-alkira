resource "alkira_policy" "test" {
  name         = "test-policy"
  description  = "test policy"
  enabled      = "false"
  from_groups  = ["-1"]
  to_groups    = ["-1"]
  rule_list_id = alkira_policy_rule_list.test.id
  segment_ids  = [alkira_segment.test.id]
}
