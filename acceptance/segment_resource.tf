resource "alkira_segment_resource" "test1" {
  name       = "acceptance-seg-res-test1"
  segment_id = alkira_segment.test1.id

  group_prefix {
    group_id       = -1
    prefix_list_id = -1
  }
}

resource "alkira_segment_resource" "test2" {
  name       = "acceptance-seg-res-test2"
  segment_id = alkira_segment.test2.id

  group_prefix {
    group_id       = -1
    prefix_list_id = -1
  }
}
