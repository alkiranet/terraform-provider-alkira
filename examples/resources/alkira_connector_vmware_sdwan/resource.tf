resource "alkira_connector_vmware_sdwan" "test" {
  name          = "test"
  cxp           = "US-WEST"
  size          = "SMALL"
  version       = "18.4.0"

  virtual_vedge {
    hostname        = "vedge1"
    activation_code = "12345678"
  }

  vrf_segment_mapping {
    segment_id = alkira_segment.test.id
    vmware_sdwan_segment_name = "segment1"
  }
}

