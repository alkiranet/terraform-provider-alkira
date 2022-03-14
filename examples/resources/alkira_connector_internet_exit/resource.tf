resource "alkira_connector_internet_exit" "test1" {
  name           = "test1"
  cxp            = "US-WEST"
  group          = "group1"
  segment_id     = alkira_segment.test.id
  size           = "SMALL"
}

