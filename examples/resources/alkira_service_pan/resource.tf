resource "alkira_service_pan" "test1" {
  name                  = "test1"
  credential_id         = alkira_credential_pan.test1.id
  cxp                   = "US-WEST"
  license_type          = "BRING_YOUR_OWN"
  management_segment_id = alkira_segment.segment1.id
  panorama_enabled      = "false"
  panorama_device_group = "alkira-test"
  panorama_ip_address   = "172.16.0.8"
  panorama_template     = "test"
  segment_ids           = [alkira_segment.segment1.id, alkira_segment.segment2.id]
  size                  = "SMALL"
  type                  = "VM-300"
  version               = "9.0.5-xfr"

  instance {
    name = "instance1"
    credential_id = alkira_credential_pan_instance.test1.id
  }

  zones_to_groups {
    segment_name = alkira_segment.segment1.name
    zone_name    = "Zone11"
    groups       = [alkira_group.pan.name]
  }
}
