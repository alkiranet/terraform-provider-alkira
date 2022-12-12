resource "alkira_service_checkpoint" "test" {
  name         = "test"

  auto_scale   = "ON"
  cxp          = "US-WEST"
  license_type = "PAY_AS_YOU_GO"
  size         = "SMALL"
  version      = "R80.30"
  pdp_ips      = ["10.0.0.1"]
  password     = "abcd1234"
  segment_id   = alkira_segment.test.id

  # max_instance_count and min_instance_count must equal each other
  # when auto_scale is off.
  max_instance_count = 2
  min_instance_count = 2

  instance {
    name    = "ins1"
    sic_key = "abcd1234"
  }

  instance {
    name    = "ins2"
    sic_key = "abcd12345"
  }

  management_server {
    type                = "MDS"
    configuration_mode  = "AUTOMATED"
    reachability        = "PRIVATE"
    ips                 = ["192.168.3.3"]
    global_cidr_list_id = alkira_list_global_cidr.test.id
    segment_id          = alkira_segment.test.id

    # user_name and management_server_password required only when
    # configuration_mode is AUTOMATED.
    user_name                  = "checkpoint_user"
    management_server_password = "abcd1234"

    # domain only required when configuration_mode is AUTOMATED and
    # when type is MDS.
    domain = "test.alkira.com"
  }
}
