resource "alkira_segment_resource" "test" {
  name       = "test"
  segment_id = alkira_segment.tftest.id

  group_prefix {
    connector_group_id = -1
    prefix_list_id     = -1
  }
}
