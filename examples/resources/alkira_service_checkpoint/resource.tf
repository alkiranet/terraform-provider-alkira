resource "alkira_service_checkpoint" "tf_checkpoint" {
  auto_scale   = "ON"
  cxp          = "US-WEST"
  name         = "chkpfw-1"
  license_type = "PAY_AS_YOU_GO"
  size         = "SMALL"
  version      = "R80.30"
  pdp_ips      = ["10.0.0.1"]
  password     = "abcd1234"

  # max_instance_count and min_instance_count must equal each other when auto_scale is off.
  max_instance_count = 2
  min_instance_count = 2


  instances {
    sic_key = "abcd1234"
  }
  instances {
    sic_key = "abcd12345"
  }

  management_server {
    type                = "MDS"
    configuration_mode  = "AUTOMATED"
    reachability        = "PRIVATE"
    ips                 = ["192.168.3.3"]
    global_cidr_list_id = alkira_list_global_cidr.checkpoint_cidr.id
    segment_id          = alkira_segment.checkpoint_seg.id

    # user_name and management_server_password required only when configuration_mode is AUTOMATED.
    user_name                  = "checkpoint_user"
    management_server_password = "abcd1234"

    # domain only required when configuration_mode is AUTOMATED and when type is MDS.
    domain = "test.alkira.com"
  }

  # only one segment allowed.    
  segment_ids = [alkira_segment.checkpoint_seg.id]
  segment_options {
    segment_id = alkira_segment.checkpoint_seg.id
    zone_name  = "DEFAULT"
    groups     = ["checkpoint_test"]
  }
}