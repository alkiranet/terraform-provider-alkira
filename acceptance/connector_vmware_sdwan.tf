resource "alkira_connector_vmware_sdwan" "test" {
  name              = "acceptance-test-vmware-sdwan"
  cxp               = "US-WEST-1"
  group             = alkira_group.test1.name
  orchestrator_host = "http://alkiratest.alkira3.net/test"
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
