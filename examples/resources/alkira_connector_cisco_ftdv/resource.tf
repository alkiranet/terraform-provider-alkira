resource "alkira_connector_cisco_ftdv" "cisco_ftdv" {
  name                = "FtdvFirewalll"
  auto_scale          = "OFF"
  size                = "SMALL"
  tunnel_protocol     = "IPSEC"
  cxp                 = "US-WEST"
  max_instance_count  = 1
  min_instance_count  = 1
  ip_allow_list       = ["192.168.3.3"]

  global_cidr_list_id = alkira_list_global_cidr.checkpoint_cidr.id
  segment_ids = [alkira_segment.tftest1.id, alkira_segment.tftest2.id]

  username = "user1"
  password = "Abcd1234567"

  management_server {
    fmc_ip       = "1.1.1.1"
    segment_id   = alkira_segment.tftest1.id
  }

  instances {
    admin_password       = "Abcd@12345"
    fmc_registration_key = "abcd12345"
    ftdv_nat_id          = "abcd1234"
    version              = "7.2.1-40"
    license_type         = "BRING_YOUR_OWN"
  }

  segment_options {
    segment_id = alkira_segment.tftest1.id
    zone_name  = "zone1"
  }
  segment_options {
    segment_id = alkira_segment.tftest2.id
    zone_name  = "zone2"
  }
}