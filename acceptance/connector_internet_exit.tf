resource "alkira_connector_internet_exit" "test1" {
  name       = "acceptance-inet-test1"
  cxp        = var.cxp
  group      = alkira_group.test1.name
  segment_id = alkira_segment.test1.id
}
