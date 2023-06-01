resource "alkira_connector_ipsec" "test" {
  name       = "acceptance-ipsec-test1"
  cxp        = "US-WEST-1"
  segment_id = alkira_segment.test1.id
  size       = "SMALL"
  vpn_mode   = "ROUTE_BASED"

  routing_options {
    type                 = "DYNAMIC"
    customer_gateway_asn = "65512"
  }

  endpoint {
    name                     = "e10"
    customer_gateway_ip      = "8.8.8.8"
    preshared_keys           = ["12345", "1235"]
    billing_tag_ids          = [alkira_billing_tag.test1.id]
    enable_tunnel_redundancy = false
    ha_mode                  = "ACTIVE"

    advanced_options {
      initiator                 = true
      dpd_delay                 = 30
      dpd_timeout               = 150
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
      local_auth_type           = "IP_ADDR"
      local_auth_value          = "172.16.1.1"
      remote_auth_type          = "IP_ADDR"
      remote_auth_value         = "54.70.233.220"
      replay_window_size        = 32
    }
  }

  endpoint {
    name                = "e2"
    customer_gateway_ip = "9.9.9.9"
    preshared_keys      = ["1234", "1235"]
    billing_tag_ids     = [alkira_billing_tag.test1.id]
  }

  endpoint {
    name                = "e3"
    customer_gateway_ip = "9.9.9.1"
    preshared_keys      = ["1234", "1235"]
    billing_tag_ids     = [alkira_billing_tag.test1.id]
  }

  endpoint {
    name                = "e4"
    customer_gateway_ip = "6.6.6.6"
    preshared_keys      = ["1234", "1235"]
    billing_tag_ids     = [alkira_billing_tag.test1.id]
  }

  endpoint {
    name                = "e5"
    customer_gateway_ip = "5.5.5.5"
    preshared_keys      = ["1234", "1235"]
    billing_tag_ids     = [alkira_billing_tag.test1.id]
  }
}
