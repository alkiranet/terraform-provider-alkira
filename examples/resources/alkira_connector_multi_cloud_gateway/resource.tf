resource "alkira_connector_multi_cloud_gateway" "example" {
  name                   = "example-mcgw-connector"
  size                   = "SMALL"
  description            = "example connector"
  cxp                    = "US-WEST"
  group                  = alkira_group.group1.name
  multi_cloud_gateway_id = alkira_multi_cloud_gateway.example.id
  loopback_prefix        = "10.30.0.0/26"

  segment_options {
    segment_id               = alkira_segment.segment1.id
    advertise_on_prem_routes = true
    disable_internet_exit    = false
  }

  instance {
    name           = "instance1"
    on_ramp_region = "us-sanjose-1 (San Jose)"
  }

  instance {
    name           = "instance2"
    on_ramp_region = "us-ashburn-1 (Ashburn)"
  }
}
