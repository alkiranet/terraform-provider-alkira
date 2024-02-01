resource "alkira_connector_remote_access" "test" {
  name                = "acceptance-remote-access"
  cxp                 = var.cxp
  segment_ids         = [alkira_segment.test1.id]
  size                = "SMALL"
  authentication_mode = "LOCAL"

  authorization {
    user_group_name = alkira_group_user.test1.name
    segment_id    = alkira_segment.test1.id
    subnet        = "172.16.1.0/24"
  }
}
