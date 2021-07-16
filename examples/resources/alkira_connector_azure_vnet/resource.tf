resource "alkira_connector_azure_vnet" "test1" {
  name           = "test1"
  azure_region   = "westus2"
  azure_vnet_id  = "/subscriptions/XXXX/resourceGroups/Test/providers/Microsoft.Network/virtualNetworks/test1"
  credential_id  = alkira_credential_azure_vnet.test1.id
  cxp            = "US-WEST"
  group          = "test"
  segment        = alkira_segment.segment1.name
  size           = "SMALL"
}
