# Lookup by name
data "alkira_peering_gateway_cxp" "by_name" {
  name = "my-peering-gateway"
}

# Lookup by ID (useful when ID is auto-generated from a connector)
data "alkira_peering_gateway_cxp" "by_id" {
  id = 12345
}

# Example: Using auto-generated peering_gateway_cxp_id from a connector
resource "alkira_connector_azure_vnet" "example" {
  name              = "example-connector"
  azure_vnet_id     = "/subscriptions/xxx/resourceGroups/xxx/providers/Microsoft.Network/virtualNetworks/xxx"
  credential_id     = "credential-id"
  cxp               = "US-WEST"
  segment_id        = "segment-id"
  # peering_gateway_cxp_id will be auto-populated by the API
}

data "alkira_peering_gateway_cxp" "from_connector" {
  id = alkira_connector_azure_vnet.example.peering_gateway_cxp_id
}
