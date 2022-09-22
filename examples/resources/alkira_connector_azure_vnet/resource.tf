#
# The following examples assumes that `alkira_segment` and
# `alkira_credential_azure_vnet` are already created.
#

#
# A simple connector could be created like this:
#
resource "alkira_connector_azure_vnet" "test1" {
  name           = "test1"
  azure_vnet_id  = "/subscriptions/XXXX/resourceGroups/Test/providers/Microsoft.Network/virtualNetworks/test1"
  credential_id  = alkira_credential_azure_vnet.test1.id
  cxp            = "US-WEST"
  group          = "test"
  segment_id     = alkira_segment.segment1.id
  size           = "SMALL"
}


#
# You could adjust routing options on VNET level by using
# `routing_options` and `routing_prefix_list_ids` along with resource
# `alkira_policy_prefix_list`:
#
resource "alkira_connector_azure_vnet" "test2" {
  name           = "test2"
  azure_vnet_id  = "/subscriptions/XXXX/resourceGroups/Test/providers/Microsoft.Network/virtualNetworks/test-vnet2"
  credential_id  = alkira_credential_azure_vnet.yours.id
  cxp            = "US-WEST"
  group          = "test"
  segment_id     = alkira_segment.segment1.id
  size           = "SMALL"

  routing_options         = "ADVERTISE_CUSTOM_PREFIX"
  routing_prefix_list_ids = [alkira_policy_prefix_list.test.id]
}


#
# Moreover, to adjust routing options for CIDR or subnet of the VNET,
# you could use `vnet_cidr` or `vnet_subnet` block:
#
resource "alkira_connector_azure_vnet" "cidr" {
  name                    = "test-cidr"
  azure_vnet_id           = azurerm_virtual_network.vnet.id
  credential_id           = alkira_credential_azure_vnet.vnet.id
  cxp                     = "US-WEST"
  group                   = alkira_group.test.name
  segment_id              = alkira_segment.test.id
  size                    = "SMALL"

  vnet_cidr {
    cidr            = "10.0.0.0/16"
    prefix_list_ids = [alkira_policy_prefix_list.azure.id]
    routing_options = "ADVERTISE_CUSTOM_PREFIX"
    service_tags    = ["ApiManagement", "AppService"]
  }
}

#
# There could be multi `vnet_subnet` blocks specified for each subnet if
# needed:
#
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
