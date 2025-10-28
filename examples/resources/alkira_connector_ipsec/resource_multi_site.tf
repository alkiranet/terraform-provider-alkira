resource "alkira_connector_ipsec" "multi_site_advanced" {
  name           = "ipsec-connector-multi-site"
  description    = "Multi-site IPSec connector with scale group"
  cxp            = "US-WEST"
  group          = alkira_group.group1.name
  segment_id     = alkira_segment.segment1.id
  size           = "LARGE"
  enabled        = true
  scale_group_id = alkira_scale_group.ipsec_scale.id

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type                 = "BOTH"
    prefix_list_id       = alkira_list_global_cidr.remote_subnets.id
    customer_gateway_asn = "65380"
    bgp_auth_key         = "multi-site-bgp-key"
    availability         = "IPSEC_INTERFACE_PING"
  }

  segment_options {
    name                     = alkira_segment.segment1.name
    advertise_default_route  = false
    advertise_on_prem_routes = true
  }

  endpoint {
    name                     = "site1-primary"
    customer_gateway_ip      = "203.0.113.80"
    preshared_keys           = ["site1-key-1", "site1-key-2"]
    billing_tag_ids          = [alkira_billing_tag.tag1.id]
    ha_mode                  = "ACTIVE"
    enable_tunnel_redundancy = true

    advanced_options {
      esp_dh_group_numbers      = ["MODP3072", "ECP256"]
      esp_encryption_algorithms = ["AES256CBC", "AES256GCM16"]
      esp_integrity_algorithms  = ["SHA256", "SHA384"]

      ike_dh_group_numbers      = ["MODP3072", "ECP256"]
      ike_encryption_algorithms = ["AES256CBC"]
      ike_integrity_algorithms  = ["SHA256"]
      ike_version               = "IKEv2"

      initiator = true

      remote_auth_type  = "FQDN"
      remote_auth_value = "site1.example.com"
    }
  }

  endpoint {
    name                     = "site2-secondary"
    customer_gateway_ip      = "203.0.113.81"
    preshared_keys           = ["site2-key-1", "site2-key-2"]
    billing_tag_ids          = [alkira_billing_tag.tag2.id]
    ha_mode                  = "ACTIVE"
    enable_tunnel_redundancy = false

    advanced_options {
      esp_dh_group_numbers      = ["MODP2048"]
      esp_encryption_algorithms = ["AES128CBC"]
      esp_integrity_algorithms  = ["SHA1"]

      ike_dh_group_numbers      = ["MODP2048"]
      ike_encryption_algorithms = ["AES128CBC"]
      ike_integrity_algorithms  = ["SHA1"]
      ike_version               = "IKEv1"

      initiator = true

      remote_auth_type  = "KEYID"
      remote_auth_value = "site2-identifier"
    }
  }

  endpoint {
    name                = "site3-backup"
    customer_gateway_ip = "203.0.113.82"
    preshared_keys      = ["site3-key"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
    ha_mode             = "STANDBY"
  }
}
