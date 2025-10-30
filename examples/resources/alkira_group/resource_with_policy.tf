# Create groups and reference them in a policy
resource "alkira_group" "source_group" {
  name        = "source-connectors"
  description = "Source connector group"
}

resource "alkira_group" "destination_group" {
  name        = "destination-services"
  description = "Destination service group"
}

# Reference group IDs in a policy
resource "alkira_policy" "example" {
  name        = "group-policy-example"
  description = "Policy using connector and service groups"
  enabled     = true

  # Reference the group IDs as outputs
  from_groups = [alkira_group.source_group.id]
  to_groups   = [alkira_group.destination_group.id]

  rule_list_id = alkira_policy_rule_list.example.id
  segment_ids  = [alkira_segment.example.id]
}
