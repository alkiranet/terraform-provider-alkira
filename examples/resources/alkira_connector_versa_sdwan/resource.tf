resource "alkira_connector_versa_sdwan" "test" {
  name    = "test"
  cxp     = "US-WEST"
  group   = alkira_group.test.name
  size    = "SMALL"

  versa_controller_host = "172.16.0.1"
  local_id  = 1
  remote_id = 2

  versa_vos_device {
    hostname                   = "dev1"
    local_device_serial_number = "12345678"
    version                    = "21.2.3-B"
  }

  vrf_segment_mapping {
    segment_id     = alkira_segment.test.id
    vrf_name       = "test"
    versa_bgp_asn  = 1203403435
  }
}


