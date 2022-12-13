resource "alkira_service_checkpoint" "test" {
  name       = "acceptance-checkpoint"
  auto_scale = "ON"
  cxp        = "US-WEST-1"

  license_type = "PAY_AS_YOU_GO"
  size         = "SMALL"
  version      = "R80.30"
  pdp_ips      = ["10.0.0.1"]
  password     = "abcd1234"

  max_instance_count = 2
  min_instance_count = 2

  segment_id = alkira_segment.test1.id

  management_server {
    type                = "MDS"
    configuration_mode  = "AUTOMATED"
    reachability        = "PRIVATE"
    ips                 = ["192.168.3.3"]
    global_cidr_list_id = alkira_list_global_cidr.checkpoint.id
    segment_id          = alkira_segment.test1.id

    user_name                  = "checkpoint_user"
    management_server_password = "abcd1234"

    # domain only required when configuration_mode is AUTOMATED and when type is MDS.
    domain = "test.alkira.com"
  }

  instance {
    name    = "ins1"
    sic_key = "abcd1234"
  }

  instance {
    name    = "ins2"
    sic_key = "abcd12345"
  }
}
