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
