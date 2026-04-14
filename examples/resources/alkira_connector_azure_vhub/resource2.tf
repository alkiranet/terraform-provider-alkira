resource "alkira_connector_azure_vhub" "test2" {
  name            = "test2"
  virtual_hub_id  = "/subscriptions/XXXX/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-vhub"
  credential_id   = alkira_credential_azure_vnet.yours.id
  cxp             = "US-WEST"
  group           = "test"
  segment_id      = alkira_segment.segment1.id
  size            = "SMALL"
  scale_group_id  = alkira_group.test.id

  vhub_routing {
    route_import_mode = "ADVERTISE_CUSTOM_PREFIX"
    prefix_list_ids   = [alkira_policy_prefix_list.test.id]
  }
}
