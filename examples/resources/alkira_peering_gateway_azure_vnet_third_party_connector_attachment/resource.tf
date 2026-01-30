resource "alkira_peering_gateway_azure_vnet_third_party_connector_attachment" "example" {
  name                    = "azure-vnet-attachment"
  description             = "Azure VNET Third Party Attachment"
  cxp_peering_gateway_id  = alkira_peering_gateway_cxp.example.id
  azure_vnet_id           = "/subscriptions/xxxx-xxxx-xxxx-xxxx/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/myVnet"
}
