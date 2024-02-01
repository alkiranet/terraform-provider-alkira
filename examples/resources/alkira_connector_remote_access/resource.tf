resource "alkira_connector_remote_access" "test1" {
  name                = "tftest-test1"
  cxp                 = "US-WEST"
  segment_ids         = [alkira_segment.test1.id]
  size                = "SMALL"
  authentication_mode = "LOCAL"

  authorization {
    user_group_name = alkira_group_user.test.name
    segment_id    = alkira_segment.test1.id
    subnet        = "172.16.1.0/24"
  }
}
