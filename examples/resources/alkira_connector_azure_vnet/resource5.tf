resource "alkira_connector_azure_vnet" "peering" {
  name            = "azure-vnet-peering"
  azure_vnet_id   = "/subscriptions/XXXX/resourceGroups/Test/providers/Microsoft.Network/virtualNetworks/test-vnet"
  credential_id   = alkira_credential_azure_vnet.peering.id
  cxp             = "USEAST-AZURE-1"
  connection_mode = "VNET_PEERING"
  group           = "production"
  segment_id      = alkira_segment.segment1.id
  size            = "SMALL"
  peering_gateway_cxp_id = alkira_peering_gateway_cxp.azure_gateway.id
}
