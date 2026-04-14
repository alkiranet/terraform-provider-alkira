resource "alkira_connector_azure_vhub" "test1" {
  name           = "test1"
  virtual_hub_id = "/subscriptions/XXXX/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-vhub"
  credential_id  = alkira_credential_azure_vnet.test1.id
  cxp            = "US-WEST"
  group          = "test"
  segment_id     = alkira_segment.segment1.id
  size           = "SMALL"

  vhub_routing {}
}
