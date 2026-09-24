resource "alkira_connector_prisma_sdwan" "test1" {
  cxp             = "US-WEST"
  group           = alkira_group.test.name
  name            = "prisma-sdwan-connector1"
  size            = "SMALL"
  tunnel_protocol = "IPSEC"

  target_segment {
    segment_id      = alkira_segment.test1.id
    vrf_name        = "prisma_vrf_1"
    gateway_bgp_asn = 65500
  }

  instances {
    credential_id = alkira_credential_prisma_sdwan.test1.id
    host_name     = "prisma-ion-1.alkira.net"
    ion_model     = "ion7108v"
    version       = "6.4.1"
  }

  instances {
    credential_id = alkira_credential_prisma_sdwan.test2.id
    host_name     = "prisma-ion-2.alkira.net"
    ion_model     = "ion7108v"
    version       = "6.4.1"
  }
}
