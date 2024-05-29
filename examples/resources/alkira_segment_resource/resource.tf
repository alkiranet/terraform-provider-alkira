resource "alkira_segment_resource" "example" {
  name       = "example"
  segment_id = alkira_segment.example.id

  group_prefix {
    group_id       = -1
    prefix_list_id = -1
  }
}
