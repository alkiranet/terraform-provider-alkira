resource "alkira_connector_remote_access" "example" {
  name                = "example"
  cxp                 = "US-WEST"
  segment_ids         = [alkira_segment.example.id]
  size                = "SMALL"
  authentication_mode = "LOCAL"

  authorization {
    user_group_name = alkira_group_user.example.name
    segment_id    = alkira_segment.example.id
    subnet        = "172.16.1.0/24"
  }
}
