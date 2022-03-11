resource "alkira_fortinet" "test1" {
  credential_id             = alkira_credential_fortinet.tf_test_fortinet.id
  cxp                       = "US-WEST"
  license_type              = "PAY_AS_YOU_GO"
  management_server_ip      = ""
  management_server_segment = alkira_segment.test1.name
  max_instance_count        = 1
  min_instance_count        = 1
  name                      = "test1-update"
  segment_names             = [alkira_segment.test1.name, alkira_segment.test2.name]
  size                      = "SMALL"
  tunnel_protocol           = "IPSEC"
  version                   = "7.0.2"

  instances {
    name          = "tf-fortinet-instance1"
    serial_number = "test-instance-1"
    credential_id = alkira_credential_fortinet_instance.tf_test_fortinet_instance.id
  }
}
