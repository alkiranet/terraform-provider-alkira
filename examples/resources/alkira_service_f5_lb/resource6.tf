resource "alkira_service_f5_lb" "example-aws-vxlan-ilb" {
  name                = "example-aws-vxlan-ilb"
  description         = "F5 LB on an AWS CXP using VXLAN, with ILB enabled."
  cxp                 = "US-EAST"
  global_cidr_list_id = alkira_list_global_cidr.example-global-cidr.id
  tunnel_protocol     = "VXLAN"
  instance {
    deployment_type = "BEST"
    hostname_fqdn   = "exampleawsvxlan.f5.local"
    license_type    = "PAY_AS_YOU_GO"
    name            = "example-aws-vxlan-ilb-instance-1"
    version         = "17.1.1.1-0.0.2"
    f5_username     = "admin"
    f5_password     = "passwordispassword"
  }
  segment_ids = [alkira_segment.example-segment.id]
  segment_options {
    elb_nic_count = 1
    segment_id    = alkira_segment.example-segment.id
    lb_type       = ["ILB"]
  }
  service_group_name     = "example-service-group-6"
  ilb_service_group_name = "example-ilb-service-group-6"
  size                   = "SMALL"
}
