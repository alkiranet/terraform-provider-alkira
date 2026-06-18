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

# Override the customer-end overlay IP for cases where the IP does not
# fit into the ranges available with the IP reservation. Mutually
# exclusive with customer_end_overlay_ip_reservation_id — supply exactly
# one. When customer_end_overlay_ip is set, the cxp-end overlay
# reservation must be /32.
resource "alkira_connector_ipsec_adv" "test_with_overlay_ip_override" {
  name       = "test-overlay-ip-override"
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

      customer_end_overlay_ip           = "10.20.30.40"
      cxp_end_overlay_ip_reservation_id = "151e8723-xxxx-4d6d-be90-xxxxxxxxxxxx"
      cxp_end_public_ip_reservation_id  = "f9f05b7a-xxxx-48eb-93e2-xxxxxxxxxxxx"
    }
  }
}
