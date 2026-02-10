resource "alkira_connector_azure_vnet_third_party" "example" {
  name        = "azure-third-party-example"
  description = "Example Azure VNET Third Party"
  cxp         = "USEAST-AZURE-1"
  segment_id  = alkira_segment.example.id
  size        = "MEDIUM"
  enabled     = true
  group       = "production"

  azure_vnet_third_party_connector_attachment_id = alkira_peering_gateway_azure_vnet_third_party_connector_attachment.example.id

  billing_tag_ids = [
    alkira_billing_tag.tag1.id,
    alkira_billing_tag.tag2.id,
  ]

  static_route_prefix_list_ids = [
    alkira_policy_prefix_list.routes.id,
  ]
}
