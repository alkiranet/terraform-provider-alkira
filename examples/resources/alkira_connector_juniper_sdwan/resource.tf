resource "alkira_connector_juniper_sdwan" "juniper" {
  name    = "test"
  cxp     = "US-EAST"
  size    = "SMALL"
  version = "6.3.4"
  group   = alkira_group.test.name

  instance {
    hostname          = "host"
    username          = "user"
    password          = "password"
    registration_key  = "registrationKey"
  }
  
  juniper_ssr_vrf_mapping {
    segment_id = alkira_segment.test.id
  }
}