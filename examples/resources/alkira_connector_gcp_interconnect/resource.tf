resource "alkira_connector_gcp_interconnect" "example_gcp_interconnect" {
  name              = "example_gcp_interconnect"
  size              = "SMALL"
  description       = "example connector"
  cxp               = "US-WEST"
  group             = alkira_group.group1.name
  tunnel_protocol   = "IPSEC"
  loopback_prefixes = ["1.2.3.4"]
  instances = {
    name                     = "instance1"
    edge_availibility_domain = "AVAILIBILITY_DOMAIN_1"
    customer_asn             = 56009
    bgp_auth_key             = "key"
  }
  segment_options {
    segment_id               = alkira_segment.segment1.id
    instance_name            = "instance1"
    advertise_on_prem_routes = true
    disable_internet_exit    = false
    customer_gateways {
      loopback_ip  = "192.192.192.192"
      tunnel_count = 2
    }
  }
}
