resource "alkira_azure_vnet_third_party_connector" "advanced" {
  name        = "azure-third-party-advanced"
  description = "Advanced configuration with billing and routing"
  cxp         = "US-EAST"
  segment_id  = alkira_segment.example.id
  size        = "MEDIUM"
  enabled     = true
  group       = "production"

  azure_vnet_third_party_connector_attachment_id = alkira_azure_vnet_third_party_connector_attachment.example.id

  billing_tag_ids = [
    alkira_billing_tag.tag1.id,
    alkira_billing_tag.tag2.id,
  ]

  static_route_prefix_list_ids = [
    alkira_policy_prefix_list.routes.id,
  ]
}
