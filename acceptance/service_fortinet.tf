resource "alkira_service_fortinet" "test1" {
  username                     = "admin"
  password                     = "Ak12345678"
  cxp                          = "US-WEST"
  license_type                 = "BRING_YOUR_OWN"
  management_server_segment_id = alkira_segment.test1.id
  max_instance_count           = 1
  min_instance_count           = 1
  name                         = "acceptance-fortinet-test1"
  segment_ids                  = [alkira_segment.test1.id, alkira_segment.test2.id]
  size                         = "SMALL"
  tunnel_protocol              = "IPSEC"
  version                      = "7.0.3"

  instances {
    name                  = "acceptance-fortinet-instance-1"
    serial_number         = "test-instance-1"
    license_key_file_path = "fortinet-test.license"
  }

  segment_options {
    segment_id = alkira_segment.test1.id
    zone_name  = "zonename"
    groups     = [alkira_group.test1.name]
  }

  segment_options {
    segment_id = alkira_segment.test1.id
    zone_name  = "zonename3"
    groups     = [alkira_group.test2.name]
  }

}
