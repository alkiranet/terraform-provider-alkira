resource "alkira_azure_vnet_third_party_connector" "example" {
  name        = "azure-third-party-connector"
  description = "Azure VNET Third Party Connector"
  cxp         = "US-WEST"
  segment_id  = alkira_segment.example.id
  size        = "SMALL"
  enabled     = true
  group       = "production"

  azure_vnet_third_party_connector_attachment_id = alkira_peering_gateway_azure_vnet_third_party_connector_attachment.example.id
}
