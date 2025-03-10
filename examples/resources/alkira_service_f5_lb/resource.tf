resource "alkira_service_f5_lb" "example-lb" {
  name                = "example-lb"
  description         = "example-lb description."
  cxp                 = "US-WEST"
  global_cidr_list_id = alkira_list_global_cidr.example-global-cidr.id
  instance {
    deployment_type     = "LTM_DNS"
    hostname_fqdn       = "examplelb.hostname"
    license_type        = "BRING_YOUR_OWN"
    name                = "example-lb-instance-1"
    version             = "17.1.1.1-0.0.2"
    deployment_option   = "TWO_BOOT_LOCATION"
    f5_registration_key = "key"
    f5_username         = "admin"
    f5_password         = "verysecretpassword"

  }
  segment_ids = [alkira_segment.example-segment.id]
  segment_options {
    elb_nic_count = 2
    segment_id    = alkira_segment.example-segment.id
  }
  service_group_name = "example-service-group"
  size               = "LARGE"
}
