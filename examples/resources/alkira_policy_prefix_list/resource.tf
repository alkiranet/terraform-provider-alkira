resource "alkira_policy_prefix_list" "test" {
  name        = "test-prefix-list"
  description = "test policy prefix list"
  prefixes    = ["0.0.0.0/0"]
}
