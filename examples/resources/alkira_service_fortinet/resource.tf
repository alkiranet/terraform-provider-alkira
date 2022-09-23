
resource "alkira_service_fortinet" "test1" {
  username                  = "admin"
  password                  = "Ak12345678"
  cxp                       = "US-WEST"
  license_type              = "BRING_YOUR_OWN"
  management_server_ip      = ""
  management_server_segment_id = alkira_segment.test1.id
  max_instance_count        = 1
  min_instance_count        = 1
  name                      = "test1-update"
  segment_ids               = [alkira_segment.test1.id, alkira_segment.test2.id]
  size                      = "SMALL"
  tunnel_protocol           = "IPSEC"
  version                   = "7.0.2"

  # You can add more instance blocks. Make sure to change "max_instance_count".
  instances {
    name                  = "tf-fortinet-instance-1"
    serial_number         = "licensekey"
    license_key_file_path = "/path/to/license.lic"
  }
  segment_options {
    segment_id = alkira_segment.test.id
    zone_name  = "zonename"
    groups     = [alkira_group.test.name]
  }

  segment_options {
    segment_id = alkira_segment.test.id
    zone_name  = "zonename1"
    groups     = [alkira_group.test.name]
  }
}