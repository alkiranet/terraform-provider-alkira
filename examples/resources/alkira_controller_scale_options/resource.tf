resource "alkira_controller_scale_options" "example" {
  name        = "example-scale-options"
  description = "Example description for controller scale options"
  entity_id   = alkira_service_fortinet.id
  entity_type = "SERVICE"
  segment_scale_options {
    additional_tunnels_per_node = 5
    segment_id                  = alkira_segment.id
    zone_name                   = "ZoneA"
  }
}

// zone is only applicable for services
// remove the network_entity_type and hardcode it as "SCALE_OPTION", they plan on adding more 'network_entity_types' later
// network_entity_id also needs to be removed as it is same as id of the resource.
