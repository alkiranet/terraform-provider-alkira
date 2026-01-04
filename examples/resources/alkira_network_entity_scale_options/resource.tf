resource "alkira_network_entity_scale_options" "example" {
  name        = "example-scale-options"
  description = "Example description for network entity scale options"
  entity_id   = alkira_service_fortinet.example.id
  entity_type = "SERVICE"
  segment_scale_options {
    additional_tunnels_per_node = 2
    segment_id                  = alkira_segment.example.id
    zone_name                   = "ZoneA"
  }
}
