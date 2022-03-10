#
# This example assumes that alkira_segment "test1", alkira_segment_resource "test1" and
# "test2" are created separately.
#
resource "alkira_segment_resource_share" "test" {
  name                       = "simple-test"
  designated_segment_id      = alkira_segment.test1.id
  service_ids                = [-1]
  end_a_segment_resource_ids = [alkira_segment_resource.test1.id]
  end_b_segment_resource_ids = [alkira_segment_resource.test2.id]
}

