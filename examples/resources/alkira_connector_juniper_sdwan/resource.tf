resource "alkira_connector_juniper_sdwan" "juniper" {
  name    = "test"
  cxp     = "US-EAST"
  size    = "SMALL"
  juniper_ssr_version = "6.3.4"
  group   = alkira_group.test.name
  availability_zone = 0

  instance {
    hostname          = "host"
    registration_key  = "registrationKey"
  }
  
  juniper_ssr_vrf_mapping {
    segment_id = alkira_segment.test.id
  }
}