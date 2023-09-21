resource "alkira_connector_fortinet_sdwan" "test" {
  name              = "test"
  cxp               = var.cxp
  group             = alkira_group.test.name
  size              = "SMALL"

  wan_edge {
    hostname      = "edge1"
    username      = "admin"
    password      = "Test1234"
    license_type  = "PAY_AS_YOU_GO"
    version       = "7.0.11"
  }

  target_segment {
    segment_id  = alkira_segment.test.id
    vrf_id      = 1
  }
}
