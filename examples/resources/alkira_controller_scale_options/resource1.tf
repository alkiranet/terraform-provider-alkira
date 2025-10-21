resource "alkira_controller_scale_options" "another_example" {
  name        = "another-scale-options"
  description = "Another example description"
  entity_id   = alkira_connector_aws_vpc.id
  entity_type = "CONNECTOR"
  segment_scale_options {
    additional_tunnels_per_node = 10
    segment_id                  = alkira_segment.example.id
  }
}