resource "alkira_network_entity_scale_options" "another_example" {
  name        = "another-scale-options"
  description = "Another example description"
  entity_id   = alkira_connector_aws_vpc.example.id
  entity_type = "CONNECTOR"
  segment_scale_options {
    additional_nodes = 2
    segment_id       = alkira_segment.example.id
  }
}
