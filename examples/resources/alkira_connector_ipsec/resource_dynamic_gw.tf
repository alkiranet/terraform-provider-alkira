resource "alkira_connector_ipsec" "dynamic_gateway" {
  name        = "ipsec-connector-dynamic-gw"
  description = "IPSec connector with dynamic customer gateway IP"
  cxp         = "US-WEST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "SMALL"
  enabled     = true

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type                 = "DYNAMIC"
    customer_gateway_asn = "65350"
  }

  endpoint {
    name                = "dynamic-site"
    customer_gateway_ip = "0.0.0.0"
    customer_ip_type    = "DYNAMIC"
    preshared_keys      = ["dynamic-key"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]

    advanced_options {
      esp_dh_group_numbers      = ["MODP2048"]
      esp_encryption_algorithms = ["AES256CBC"]
      esp_integrity_algorithms  = ["SHA256"]

      ike_dh_group_numbers      = ["MODP2048"]
      ike_encryption_algorithms = ["AES256CBC"]
      ike_integrity_algorithms  = ["SHA256"]
      ike_version               = "IKEv2"

      initiator = true

      remote_auth_type  = "FQDN"
      remote_auth_value = "remote-site.example.com"
    }
  }
}
