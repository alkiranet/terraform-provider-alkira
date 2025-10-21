resource "alkira_controller_scale_options" "example" {
  name        = "example-scale-options"
  description = "Example description for controller scale options"
  entity_id   = alkira_service_fortinet.id
  entity_type = "SERVICE"
  segment_scale_options {
    additional_tunnels_per_node = 5
    additional_nodes            = 2
    segment_id                  = alkira_segment.id
    zone_name                   = "ZoneA"
  }
}

