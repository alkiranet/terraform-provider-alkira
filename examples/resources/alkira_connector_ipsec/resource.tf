# Basic IPSec Connector Example
resource "alkira_connector_ipsec" "basic" {
  name           = "ipsec-connector-basic"
  description    = "Basic IPSec connector example"
  cxp            = "US-WEST"
  group          = alkira_group.group1.name
  segment_id     = alkira_segment.segment1.id
  size           = "SMALL"
  enabled        = true
  billing_tag_ids = [alkira_billing_tag.tag1.id]

  vpn_mode       = "ROUTE_BASED"

  routing_options {
    type = "DYNAMIC"
    customer_gateway_asn = "65310"
  }

  endpoint {
    name                = "remote-site"
    customer_gateway_ip = "203.0.113.1"
    preshared_keys      = ["your-preshared-key"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
  }
}

# Advanced IPSec Connector with High Availability
resource "alkira_connector_ipsec" "advanced" {
  name           = "ipsec-connector-advanced"
  description    = "Advanced IPSec connector with high availability"
  cxp            = "US-WEST"
  failover_cxps  = ["US-EAST"]
  group          = alkira_group.group1.name
  segment_id     = alkira_segment.segment1.id
  size           = "MEDIUM"
  enabled        = true
  billing_tag_ids = [alkira_billing_tag.tag1.id]
  scale_group_id = alkira_scale_group.ipsec_scale.id

  vpn_mode       = "ROUTE_BASED"
  health_check_type = "IPSEC"

  routing_options {
    type = "DYNAMIC"
    customer_gateway_asn = "65310"
  }

  # Primary site with redundancy
  endpoint {
    name                     = "primary-site"
    customer_gateway_ip      = "203.0.113.1"
    preshared_keys           = ["primary-key-1", "primary-key-2"]
    billing_tag_ids          = [alkira_billing_tag.tag1.id]
    enable_tunnel_redundancy = true

    # Advanced IPSec options
    advanced_options {
      esp_dh_group_numbers      = ["MODP3072"]
      esp_encryption_algorithms = ["AES256CBC"]
      esp_integrity_algorithms  = ["SHA256"]

      ike_dh_group_numbers      = ["MODP3072"]
      ike_encryption_algorithms = ["AES256CBC"]
      ike_integrity_algorithms  = ["SHA256"]
      ike_version               = "IKEv2"

      initiator          = true

      remote_auth_type   = "IP_ADDR"
      remote_auth_value  = "203.0.113.1"
    }
  }

  # Secondary site for redundancy
  endpoint {
    name                = "secondary-site"
    customer_gateway_ip = "203.0.113.2"
    preshared_keys      = ["secondary-key-1", "secondary-key-2"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
    enable_tunnel_redundancy = true
  }
}

# Policy-based IPSec Connector
resource "alkira_connector_ipsec" "policy_based" {
  name           = "ipsec-connector-policy"
  description    = "Policy-based IPSec connector"
  cxp            = "US-WEST"
  group          = alkira_group.group1.name
  segment_id     = alkira_segment.segment1.id
  size           = "SMALL"
  enabled        = true
  billing_tag_ids = [alkira_billing_tag.tag1.id]

  vpn_mode       = "POLICY_BASED"

  routing_options {
    type = "STATIC"
    prefix_list_ids = [alkira_list_global_cidr.remote_subnets.id]
  }

  endpoint {
    name                = "branch-office"
    customer_gateway_ip = "203.0.113.10"
    preshared_keys      = ["branch-office-key"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
  }
}

# Supporting resources
resource "alkira_billing_tag" "tag1" {
  name = "ipsec-connector-tag"
  description = "Billing tag for IPSec connectors"
}

resource "alkira_scale_group" "ipsec_scale" {
  name = "ipsec-scale-group"
  description = "Scale group for IPSec connectors"
  min_instance_count = 2
  max_instance_count = 4
}

resource "alkira_list_global_cidr" "remote_subnets" {
  name = "remote-subnets"
  description = "Remote subnets for static routing"
  values = ["192.168.1.0/24", "192.168.2.0/24"]
}

resource "alkira_group" "group1" {
  name = "ipsec-group"
  description = "Group for IPSec connectors"
}

resource "alkira_segment" "segment1" {
  name = "production-segment"
  asn  = "65001"
  cidrs = ["10.0.0.0/8"]
}
