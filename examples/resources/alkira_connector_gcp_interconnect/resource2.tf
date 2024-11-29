resource "alkira_connector_gcp_interconnect" "example_gcp_interconnect_2" {
  name              = "example_gcp_interconnect_1"
  size              = "LARGE"
  description       = "example connector with multiple segment options"
  cxp               = "US-WEST"
  group             = alkira_group.group1.name
  tunnel_protocol   = "IPSEC"
  loopback_prefixes = ["10.40.0.0/24"]
  instances {
    name                     = "instance1"
    edge_availability_domain = "AVAILABILITY_DOMAIN_1"
    customer_asn             = 56009
    bgp_auth_key             = "key"
    segment_options {
      segment_id               = alkira_segment.segment1.id
      advertise_on_prem_routes = true
      advertise_default_routes = false
      customer_gateways {
        tunnel_count = 2
      }
    }
    segment_options {
      segment_id               = alkira_segment.segment2.id
      advertise_on_prem_routes = true
      advertise_default_routes = false
      customer_gateways {
        tunnel_count = 2
      }
    }
  }
}
