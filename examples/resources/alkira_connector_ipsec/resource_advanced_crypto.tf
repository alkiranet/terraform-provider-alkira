resource "alkira_connector_ipsec" "advanced_crypto" {
  name        = "ipsec-connector-advanced-crypto"
  description = "IPSec connector with custom cryptographic algorithms"
  cxp         = "US-WEST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "SMALL"
  enabled     = true

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type                 = "DYNAMIC"
    customer_gateway_asn = "65360"
  }

  endpoint {
    name                = "crypto-site"
    customer_gateway_ip = "203.0.113.60"
    preshared_keys      = ["crypto-key"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]

    advanced_options {
      esp_dh_group_numbers      = ["MODP4096", "ECP384"]
      esp_encryption_algorithms = ["AES256GCM16", "AES256CBC"]
      esp_integrity_algorithms  = ["SHA512", "SHA384"]

      ike_dh_group_numbers      = ["MODP4096", "ECP384"]
      ike_encryption_algorithms = ["AES256CBC", "AES192CBC"]
      ike_integrity_algorithms  = ["SHA512", "SHA384"]
      ike_version               = "IKEv1"

      initiator = false

      remote_auth_type  = "IP_ADDR"
      remote_auth_value = "203.0.113.60"
    }
  }
}
