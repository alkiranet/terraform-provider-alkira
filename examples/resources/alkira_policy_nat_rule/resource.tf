# Example 1: Basic NAT rule with no translation
resource "alkira_policy_nat_rule" "basic" {
  name        = "basic-nat-rule"
  description = "Basic NAT rule example"
  enabled     = true
  category    = "DEFAULT"

  match {
    src_prefixes = ["10.0.0.0/8"]
    dst_prefixes = ["192.168.0.0/16"]
    protocol     = "tcp"
    src_ports    = ["80", "443"]
    dst_ports    = ["80", "443"]
  }

  action {
    src_addr_translation_type = "NONE"
    dst_addr_translation_type = "NONE"
  }
}

# Example 2: NAT rule with static IP translation
resource "alkira_policy_nat_rule" "static_ip" {
  name        = "static-ip-nat-rule"
  description = "NAT rule with static IP translation"
  enabled     = true
  category    = "DEFAULT"

  match {
    src_prefixes = ["10.0.0.0/8"]
    dst_prefixes = ["any"]
    protocol     = "any"
  }

  action {
    src_addr_translation_type     = "STATIC_IP"
    src_addr_translation_prefixes = ["192.168.1.0/24"]
    dst_addr_translation_type     = "NONE"
    egress_type                   = "ALKIRA_PUBLIC_IP"
  }
}

# Example 3: NAT rule with dynamic IP and port translation
resource "alkira_policy_nat_rule" "dynamic_ip_port" {
  name        = "dynamic-ip-port-nat-rule"
  description = "NAT rule with dynamic IP and port translation"
  enabled     = true
  category    = "INTERNET_CONNECTOR"

  match {
    src_prefixes = ["10.0.0.0/8"]
    dst_prefixes = ["any"]
    protocol     = "tcp"
  }

  action {
    src_addr_translation_type = "DYNAMIC_IP_AND_PORT"
    dst_addr_translation_type = "NONE"
    egress_type               = "BYOIP"
  }
}

# Example 4: NAT rule with destination translation and port mapping
resource "alkira_policy_nat_rule" "dest_translation" {
  name        = "dest-translation-nat-rule"
  description = "NAT rule with destination translation"
  enabled     = true

  match {
    src_prefixes = ["any"]
    dst_prefixes = ["203.0.113.0/24"]
    protocol     = "tcp"
    dst_ports    = ["80"]
  }

  action {
    src_addr_translation_type                   = "NONE"
    dst_addr_translation_type                   = "STATIC_IP_AND_PORT"
    dst_addr_translation_prefixes               = ["10.0.1.100/32"]
    dst_addr_translation_ports                  = ["8080"]
    dst_addr_translation_advertise_to_connector = true
  }
}

# Example 5: Advanced NAT rule with prefix lists and routing tracking
resource "alkira_policy_prefix_list" "source_prefixes" {
  name        = "source-prefixes-list"
  description = "Source prefixes for NAT rule"

  prefix {
    cidr = "10.0.0.0/8"
  }
  prefix {
    cidr = "172.16.0.0/12"
  }
}

resource "alkira_policy_prefix_list" "dest_prefixes" {
  name        = "dest-prefixes-list"
  description = "Destination prefixes for NAT rule"

  prefix {
    cidr = "192.168.0.0/16"
  }
}

resource "alkira_policy_prefix_list" "track_prefixes" {
  name        = "track-prefixes-list"
  description = "Prefixes to track for routing"

  prefix {
    cidr = "203.0.113.0/24"
  }
}

resource "alkira_policy_nat_rule" "advanced" {
  name        = "advanced-nat-rule"
  description = "Advanced NAT rule with prefix lists and routing tracking"
  enabled     = true
  category    = "DEFAULT"

  match {
    src_prefix_list_ids = [alkira_policy_prefix_list.source_prefixes.id]
    dst_prefix_list_ids = [alkira_policy_prefix_list.dest_prefixes.id]
    protocol            = "tcp"
  }

  action {
    src_addr_translation_type                              = "STATIC_IP"
    src_addr_translation_prefixes                          = ["198.51.100.0/24"]
    src_addr_translation_routing_track_prefix_list_ids     = [alkira_policy_prefix_list.track_prefixes.id]
    src_addr_translation_routing_track_invalidate_prefixes = true
    dst_addr_translation_type                              = "NONE"
    egress_type                                            = "ALKIRA_PUBLIC_IP"
  }
}

# Example 6: OUTBOUND NAT rule with dynamic IP and port translation
resource "alkira_policy_nat_rule" "outbound_nat" {
  name        = "outbound-nat-rule"
  description = "Outbound NAT rule with dynamic IP and port translation"
  enabled     = true
  direction   = "OUTBOUND"

  match {
    src_prefixes = ["192.168.0.0/16"]
    dst_prefixes = ["any"]
    src_ports    = ["any"]
    dst_ports    = ["any"]
    protocol     = "any"
  }

  action {
    src_addr_translation_type                = "DYNAMIC_IP_AND_PORT"
    src_addr_translation_prefixes            = ["10.1.55.100/31"]
    src_addr_translation_match_and_invalidate = false
    dst_addr_translation_type                = "NONE"
  }
}