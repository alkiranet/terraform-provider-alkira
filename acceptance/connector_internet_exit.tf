resource "alkira_connector_internet_exit" "test1" {
  name           = "acceptance-test1"
  cxp            = "US-WEST-1"
  group          = alkira_group.test1.name
  segment_id     = alkira_segment.test1.id
}
