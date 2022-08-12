resource "alkira_connector_ipsec" "ipsec" {
  name           = "connector-test-ipsec"
  cxp            = "US-WEST"
  group          = alkira_group.group1.name
  segment_id     = alkira_segment.segment1.id
  size           = "SMALL"

  vpn_mode       = "ROUTE_BASED"

  routing_options {
    type = "DYNAMIC"
    customer_gateway_asn = "65310"
  }

  # There could be multiple endpoints defined.
  endpoint {
    name                     = "Site1"
    customer_gateway_ip      = "8.8.8.8"
    preshared_keys           = ["1234", "1235"]
    billing_tag_ids          = [alkira_billing_tag.tag1.id]
    enable_tunnel_redundancy = false

    # Optional advanced options could be specified per endpoint.
    advanced_options {
      dpd_delay   = 30
      dpd_timeout = 150

      esp_dh_group_numbers      = ["MODP3072"]
      esp_encryption_algorithms = ["AES256CBC"]
      esp_integrity_algorithms  = ["SHA256"]
      esp_life_time             = 3960
      esp_random_time           = 360
      esp_rekey_time            = 3600

      ike_dh_group_numbers      = ["MODP3072"]
      ike_encryption_algorithms = ["AES256CBC"]
      ike_integrity_algorithms  = ["SHA256"]
      ike_over_time             = 2880
      ike_random_time           = 2880
      ike_rekey_time            = 28800
      ike_version               = "IKEv2"

      initiator          = true

      local_auth_type    = "IP_ADDR"
      local_auth_value   = "172.16.1.1"

      remote_auth_type   = "IP_ADDR"
      remote_auth_value  = "54.70.233.220"

      replay_window_size = 32
    }
  }

  endpoint {
      name                 = "Site2"
      customer_gateway_ip  = "9.9.9.9"
      preshared_keys       = ["1234", "1235"]
      billing_tag_ids      = [alkira_billing_tag.tag1.id]
  }
}
