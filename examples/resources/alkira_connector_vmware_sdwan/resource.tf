resource "alkira_connector_vmware_sdwan" "test" {
  name              = "test"
  cxp               = "US-WEST"
  group             = alkira_group.test1.name
  orchestrator_host = "http://test.alkira.com/portal"
  size              = "SMALL"
  version           = "4.3.1"

  virtual_edge {
    hostname        = "vedge1"
    activation_code = "12345678"
  }

  target_segment {
    segment_id                = alkira_segment.test1.id
    vmware_sdwan_segment_name = "segment1"
  }
}
