resource "alkira_service_checkpoint" "test1" {
  auto_scale         = "OFF"
  cxp                = "US-WEST"
  credential_id      = alkira_credential_checkpoint.tf_test_checkpoint.id
  license_type       = "PAY_AS_YOU_GO"
  max_instance_count = 2
  min_instance_count = 2
  name               = "testname"
  segment_names      = [alkira_segment.test-seg-1.name]
  size               = "LARGE"
  tunnel_protocol    = "IPSEC"
  version            = "R80.30"

  segment_options {
    segment_id = alkira_segment.test-seg-1.id
    zone_name  = "DEFAULT"
    groups     = [alkira_group.test.name]
  }

  management_server {
    configuration_mode  = "MANUAL"
    global_cidr_list_id = 22
    ips                 = ["10.2.0.3"]
    reachability        = "PRIVATE"
    segment_id          = alkira_segment.test-seg-1.id
    type                = "SMS"
    user_name           = "admin"
  }
}