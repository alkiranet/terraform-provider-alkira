resource "alkira_connector_ipsec_adv" "test" {
  name       = "test"
  segment_id = alkira_segment.test.id
  cxp        = "US-EAST"
  size       = "SMALL"
  vpn_mode   = "ROUTE_BASED"
  enabled    = true

  advertise_default_route  = false
  advertise_on_prem_routes = false
  tunnels_per_gateway      = 1

  gateway {
    name                = "site1"
    customer_gateway_ip = "xx.xxx.xxx.xxx"
    ha_mode             = "ACTIVE"

    tunnel {
      preshared_key = "1234"

      customer_end_overlay_ip_reservation_id = "151e8723-xxxx-4d6d-be90-xxxxxxxxxxxx"
      cxp_end_overlay_ip_reservation_id      = "151e8723-xxxx-4d6d-be90-xxxxxxxxxxxx"
      cxp_end_public_ip_reservation_id       = "f9f05b7a-xxxx-48eb-93e2-xxxxxxxxxxxx"
    }
  }
}
