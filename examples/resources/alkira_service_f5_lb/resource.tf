resource "alkira_service_f5_lb" "example_lb" {
  name                = "example_lb"
  description         = "example_lb description."
  cxp                 = "US-WEST"
  global_cidr_list_id = alkira_list_global_cidr.example_global_cidr.id
  instance {
    deployment_type     = "ALL"
    hostname_fqdn       = "example_lb.hostname"
    license_type        = "BRING_YOUR_OWN"
    name                = "example_lb_instance_1"
    version             = "17.1.1.1-0.0.2"
    deployment_option   = "ONE_BOOT_LOCATION"
    f5_registration_key = "key"
    f5_username         = "admin"
    f5_password         = "verysecretpassword"

  }
  segment_ids = [alkira_segment.example_segment.id]
  segment_options {
    elb_nic_count = 2
    segment_id    = alkira_segment.example_segment.id
  }
  service_group_name = "example_service_group"
  size               = "LARGE"
}
