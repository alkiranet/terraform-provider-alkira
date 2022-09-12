resource "alkira_service_checkpoint" "test1" {
  auto_scale         = "OFF"
  cxp                = "US-WEST-1"
  credential_id      = alkira_credential_checkpoint.tf_test_checkpoint.id
  license_type       = "PAY_AS_YOU_GO"
  max_instance_count = 2
  min_instance_count = 2
  name               = "testname"
  segment_ids        = [alkira_segment.segment.id]
  size               = "LARGE"
  tunnel_protocol    = "IPSEC"
  version            = "R80.30"

  segment_options {
    segment_id = alkira_segment.segment.id
    zone_name  = "DEFAULT"
    groups     = [alkira_group.group.name, alkira_group.group1.name]
  }

  segment_options {
    segment_id = alkira_segment.segment1.id
    zone_name  = "zonename1"
    groups     = [alkira_group.group2.name]
  }

  management_server {
    configuration_mode  = "MANUAL"
    global_cidr_list_id = alkira_list_global_cidr.cidr.id
    ips                 = ["10.2.0.3"]
    reachability        = "PRIVATE"
    segment_id          = alkira_segment.segment1.id
    type                = "SMS"
    user_name           = "admin"
  }
}

