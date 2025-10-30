# Service group for security services
resource "alkira_group" "security_services" {
  name        = "security-services-group"
  description = "Group for security services like firewalls and IDS/IPS"
}
