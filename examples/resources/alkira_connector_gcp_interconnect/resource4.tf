resource "alkira_connector_gcp_interconnect" "example_gcp_interconnect_4" {
  name              = "example_gcp_interconnect_4"
  size              = "SMALL"
  description       = "example connector"
  cxp               = "US-WEST"
  group             = alkira_group.group1.name
  tunnel_protocol   = "GRE"
  loopback_prefixes = ["10.30.0.0/24"]
  instances {
    name                     = "instance1"
    edge_availability_domain = "AVAILABILITY_DOMAIN_1"
    customer_asn             = 56009
    bgp_auth_key             = "key"
    segment_options {
      segment_id               = alkira_segment.segment1.id
      advertise_on_prem_routes = true
      advertise_default_route = false
      customer_gateways {
        tunnel_count = 2
      }
      customer_gateways {
        tunnel_count = 2
      }
    }
  }
}
