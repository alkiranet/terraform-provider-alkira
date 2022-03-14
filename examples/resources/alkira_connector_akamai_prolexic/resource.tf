#
# This example assumes that resource "alkira_group", "alkira_segment"
# and "alkira_byoip_prefix" are created separately.
#
resource "alkira_connector_akamai_prolexic" "test" {
  name           = "tftest"
  cxp            = "US-WEST"
  group          = alkira_group.test.name
  segment_id     = alkira_segment.test.id
  size           = "SMALL"

  akamai_bgp_asn = 65514
  akamai_bgp_authentication_key = "xxxxxx"

  byoip_options {
    byoip_prefix_id = alkira_byoip_prefix.test.id
    enable_route_advertisement = true
  }

  tunnel_configuration {
    alkira_public_ip = "192.168.1.1"

    tunnel_ips {
      ran_tunnel_ip = "172.16.1.10"
      alkira_overlay_tunnel_ip  = "8.8.8.8"
      akamai_overlay_tunnel_ip  = "192.168.1.1"
    }

    tunnel_ips {
      ran_tunnel_ip = "172.16.1.20"
      alkira_overlay_tunnel_ip  = "8.8.8.8"
      akamai_overlay_tunnel_ip  = "192.168.1.1"
    }

    tunnel_ips {
      ran_tunnel_ip = "172.16.1.30"
      alkira_overlay_tunnel_ip  = "8.8.8.8"
      akamai_overlay_tunnel_ip  = "192.168.1.1"
    }
  }
}

