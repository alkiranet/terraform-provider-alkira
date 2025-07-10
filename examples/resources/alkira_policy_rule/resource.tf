# Basic Policy Rule Example
resource "alkira_policy_rule" "basic" {
  name          = "basic-drop-rule"
  description   = "Basic policy rule to drop traffic"
  src_ip        = "any"
  dst_ip        = "172.16.0.0/16"
  dscp          = "any"
  protocol      = "any"
  src_ports     = ["any"]
  dst_ports     = ["any"]
  rule_action   = "DROP"
}

# Advanced Policy Rule with Prefix Lists
resource "alkira_policy_rule" "advanced" {
  name                  = "advanced-service-rule"
  description           = "Advanced policy rule with service chaining"
  src_prefix_list_id    = alkira_list_global_cidr.internal_networks.id
  dst_prefix_list_id    = alkira_list_global_cidr.dmz_networks.id
  dscp                  = "any"
  protocol              = "tcp"
  src_ports             = ["any"]
  dst_ports             = ["80", "443"]
  rule_action           = "ALLOW"

  # Service chaining
  rule_action_service_types = ["FIREWALL", "IDS"]
  rule_action_service_ids   = [alkira_service_pan.firewall.id, alkira_service_checkpoint.ids.id]

  # Flow collection
  rule_action_flow_collector_ids = [alkira_flow_collector.security.id]
}

# Policy Rule with Internet Application
resource "alkira_policy_rule" "internet_app" {
  name                    = "internet-app-rule"
  description             = "Policy rule for internet application access"
  src_ip                  = "10.0.0.0/8"
  internet_application_id = alkira_internet_application.office365.id
  dscp                    = "any"
  protocol                = "tcp"
  src_ports               = ["any"]
  dst_ports               = ["443"]
  rule_action             = "ALLOW"
}

# Policy Rule with Application IDs
resource "alkira_policy_rule" "application_based" {
  name           = "application-based-rule"
  description    = "Policy rule based on application identification"
  src_ip         = "any"
  dst_ip         = "any"
  dscp           = "any"
  protocol       = "any"
  src_ports      = ["any"]
  dst_ports      = ["any"]
  application_ids = ["HTTP", "HTTPS", "SSH"]
  rule_action    = "ALLOW"

  # Route through security services
  rule_action_service_types = ["FIREWALL"]
  rule_action_service_ids   = [alkira_service_pan.firewall.id]
}
