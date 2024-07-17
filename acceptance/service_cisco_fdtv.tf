resource "alkira_service_cisco_ftdv" "cisco_ftdv_test" {
  name                = "acceptance-ftdv-test1"
  auto_scale          = "OFF"
  size                = "SMALL"
  tunnel_protocol     = "IPSEC"
  cxp                 = var.cxp
  global_cidr_list_id = alkira_list_global_cidr.ciscofdtv.id
  max_instance_count  = 1
  min_instance_count  = 1

  segment_ids = [alkira_segment.test3.id, alkira_segment.test4.id]

  firepower_management_center {
    username      = "admin"
    password      = "Abcd1234567"
    server_ip     = "1.1.1.2"
    segment_id    = alkira_segment.test4.id
    ip_allow_list = ["192.168.3.3"]
  }

  instance {
    hostname             = "ins-1"
    admin_password       = "Test@2018"
    fmc_registration_key = "abcd12345"
    ftdv_nat_id          = "abcd1234"
    version              = "7.2.7-500"
    license_type         = "BRING_YOUR_OWN"
  }

  segment_options {
    segment_id = alkira_segment.test3.id
    zone_name  = "zone1"
    groups     = [alkira_group.test4.name]
  }

  segment_options {
    segment_id = alkira_segment.test4.id
    zone_name  = "zone2"
  }
}
