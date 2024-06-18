resource "alkira_virtual_network_manager_azure" "exampleAvnm" {
  name                   = "exampleAvnm"
  region                 = "eastus"
  cxp                    = "US-WEST"
  subscription_id        = "azure-subscription-id"
  description            = "Example description"
  resource_group         = azurerm_resource_group.exampleResourceGroup.name
  credential_id          = alkira_credential_azure_vnet.exampleCredentials.id
  subscriptions_in_scope = ["azure-subscription-id"]
}
