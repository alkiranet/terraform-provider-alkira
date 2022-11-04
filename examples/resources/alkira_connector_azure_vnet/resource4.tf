resource "alkira_connector_azure_vnet" "subnet" {
  name                    = "test-subnet"
  azure_vnet_id           = azurerm_virtual_network.vnet.id
  credential_id           = alkira_credential_azure_vnet.testtest.id
  cxp                     = "US-WEST"
  group                   = alkira_group.test.name
  segment_id              = alkira_segment.test.id
  size                    = "SMALL"

  vnet_subnet {
    subnet_id       = data.azurerm_subnet.vnet.id
    prefix_list_ids = [alkira_policy_prefix_list.azure.id]
    routing_options = "ADVERTISE_CUSTOM_PREFIX"
    service_tags    = ["ApiManagement", "AppService"]
  }
}
