resource "alkira_connector_cisco_ftdv" "cisco_ftdv_test" {
  name                = "ftdvFirewalll"
  auto_scale          = "OFF"
  size                = "SMALL"
  tunnel_protocol     = "IPSEC"
  cxp                 = "US-WEST-1"
  global_cidr_list_id = alkira_list_global_cidr.test.id
  max_instance_count  = 1
  min_instance_count  = 1
  ip_allow_list       = ["192.168.3.3"]

  segment_ids = [alkira_segment.test4.id]

  username = "admin"
  password = "Abcd1234567"

  management_server {
    fmc_ip       = "1.1.1.2"
    segment_name = alkira_segment.test4.name
    segment_id   = alkira_segment.test4.id
  }

  instances {
    admin_password       = "Test@2018"
    fmc_registration_key = "abcd12345"
    ftdv_nat_id          = "abcd1234"
    version              = "7.2.1-40"
    license_type         = "BRING_YOUR_OWN"
  }

  segment_options {
    segment_id = alkira_segment.test4.id
    zone_name  = "zone1"
  }
}
