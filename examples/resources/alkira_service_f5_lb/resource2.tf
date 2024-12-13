resource "alkira_service_f5_lb" "example_lb_2" {
  name                = "example-2"
  description         = "example-2 description."
  cxp                 = "US-WEST"
  global_cidr_list_id = alkira_list_global_cidr.example_global_cidr.id
  prefix_list_id      = alkira_list_prefix_list.example_prefix_list.id
  instance {
    deployment_type = "GOOD"
    hostname_fqdn   = "example2.hostname.1"
    license_type    = "PAY_AS_YOU_GO"
    name            = "example-2-instance-1"
    version         = "17.1.1.1-0.0.2"
    f5_password     = "passwordispassword"
    f5_username     = "admin"

  }
  instance {
    deployment_type = "GOOD"
    hostname_fqdn   = "example2.hostname.2"
    license_type    = "PAY_AS_YOU_GO"
    name            = "example-2-instance-2"
    version         = "17.1.1.1-0.0.2"
    f5_password     = "passwordispassword"
    f5_username     = "admin"

  }
  segment_ids = [alkira_segment.example_segment.id]
  segment_options {
    elb_nic_count = 2
    segment_id    = alkira_segment.example_segment.id
  }
  segment_options {
    elb_nic_count = 2
    segment_id    = alkira_segment_1.example_segment.id
  }
  service_group_name = "example_service_group_1"
  size               = "2LARGE"
}
