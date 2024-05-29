resource "alkira_service_cisco_ftdv" "example" {
  name               = "example"
  auto_scale         = "OFF"
  size               = "SMALL"
  tunnel_protocol    = "IPSEC"
  cxp                = "US-WEST"
  max_instance_count = 1
  min_instance_count = 1

  global_cidr_list_id = alkira_list_global_cidr.example.id
  segment_ids         = [alkira_segment.test1.id, alkira_segment.example.id]

  firepower_management_center {
    username      = "user1"
    password      = "Abcd1234567"
    server_ip     = "1.1.1.1"
    segment_id    = alkira_segment.example.id
    ip_allow_list = ["192.168.3.3"]
  }

  instance {
    hostname             = "instance1"
    admin_password       = "Abcd@12345"
    fmc_registration_key = "abcd12345"
    ftdv_nat_id          = "abcd1234"
    version              = "7.2.1-40"
    license_type         = "BRING_YOUR_OWN"
  }

  segment_options {
    segment_id = alkira_segment.example.id
    zone_name  = "zone1"
    groups     = [alkira_group.example.name]
  }

  segment_options {
    segment_id = alkira_segment.example.id
    zone_name  = "zone2"
  }
}
