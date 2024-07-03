resource "alkira_connector_aruba_edge" "test1" {
  boost_mode      = false
  cxp             = "US-WEST"
  gateway_gbp_asn = 22
  group           = alkira_group.test.name
  name            = "thisisanewname"
  segment_ids     = [alkira_segment.test1.id]
  size            = "SMALL"
  tunnel_protocol = "IPSEC"
  version         = "9.0.3.3"

  aruba_edge_vrf_mapping {
    segment_id                 = alkira_segment.test1.id
    aruba_edge_connect_segment = "aruba_edge_segment_name"
    gateway_gbp_asn            = 88
  }

  instances {
    account_name = "aruba-edge-account-name3"
    account_key  = "accountkey555555555"
    host_name    = "alkira.net2"
    name         = "instance2-aruba-edge1"
    site_tag     = "site_tag2"
  }

  instances {
    account_name = "aruba-edge-account-name5"
    account_key  = "accountkey33333333"
    host_name    = "alkira.net5"
    name         = "instance1-aruba-edge5"
    site_tag     = "site_tag5"
  }
}
