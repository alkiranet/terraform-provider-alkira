resource "alkira_service_fortinet" "fortinet" {

  credential_id             = alkira_credential_fortinet.credential.id
  cxp                       = "US-WEST-1"
  license_type              = "PAY_AS_YOU_GO"
  management_server_ip      = ""
  management_server_segment = alkira_segment.segment.name
  max_instance_count        = 1
  min_instance_count        = 1
  name                      = "test1-update"
  segment_ids               = [alkira_segment.segment.id, alkira_segment.segment1.id]
  size                      = "SMALL"
  tunnel_protocol           = "IPSEC"
  version                   = "7.0.2"

  instances {
    name          = "tf-fortinet-instance-1"
    serial_number = "mactest-instance-1"
    credential_id = alkira_credential_fortinet.credential.id
  }

  segment_options {
    segment_id = alkira_segment.segment.id
    zone_name  = "zonename"
    groups     = [alkira_group.group.name, alkira_group.group1.name]
  }

  segment_options {
    segment_id = alkira_segment.segment1.id
    zone_name  = "zonename1"
    groups     = [alkira_group.group2.name]
  }
}
