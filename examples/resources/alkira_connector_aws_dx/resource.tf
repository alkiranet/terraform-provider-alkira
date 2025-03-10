resource "alkira_connector_aws_dx" "test" {
  name            = "example"
  description     = "example"
  cxp             = "US-WEST"
  size            = "2LARGE"
  tunnel_protocol = "GRE"
  group           = alkira_group.example.name
  billing_tag_ids = [alkira_billing_tag.example.id]


  instance {
    name          = "instance1"
    connection_id = "test-id"

    dx_asn        = 64850
    dx_gateway_ip = "169.254.199.1"

    on_prem_asn        = 65000
    on_prem_gateway_ip = "169.254.199.2"

    underlay_prefix = "169.254.199.0/30"

    bgp_auth_key        = "Alkira2018"
    bgp_auth_key_alkira = "Alkira2018"

    vlan_id       = 305
    aws_region    = "us-west-1"
    credential_id = alkira_credential_aws_vpc.example.id

    segment_options {
      segment_id          = alkira_segment.example.id
      on_prem_segment_asn = 64303

      customer_loopback_ip = "192.168.23.243"
      alkira_loopback_ip1  = "192.168.23.188"
      alkira_loopback_ip2  = "192.168.23.205"
      loopback_subnet      = "192.168.23.0/24"

      advertise_on_prem_routes = false
    }
  }
}
