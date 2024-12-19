resource "alkira_service_f5_lb" "example-lb-1" {
  name                = "example-lb-1"
  description         = "example-lb-1 description."
  cxp                 = "US-WEST"
  global_cidr_list_id = alkira_list_global_cidr.example-global-cidr.id
  prefix_list_id      = alkira_list_prefix_list.example-prefix-list.id
  instance {
    deployment_type = "GOOD"
    hostname_fqdn   = "examplelb.hostname.1"
    license_type    = "PAY_AS_YOU_GO"
    name            = "example-lb-1-instance-1"
    version         = "17.1.1.1-0.0.2"
    f5_password     = "passwordispassword"
    f5_username     = "admin"

  }
  instance {
    deployment_type = "GOOD"
    hostname_fqdn   = "examplelb.hostname.1"
    license_type    = "PAY_AS_YOU_GO"
    name            = "example-lb-1-instance-2"
    version         = "17.1.1.1-0.0.2"
    f5_password     = "passwordispassword"
    f5_username     = "admin"

  }
  segment_ids = [alkira_segment.example-segment.id]
  segment_options {
    elb_nic_count = 2
    segment_id    = alkira_segment.example-segment.id
  }
  service_group_name = "example-service-group-1"
  size               = "2LARGE"
}
