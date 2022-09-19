resource "alkira_connector_azure_er" "test_basic" {
  name            = "AzureErName"
  size            = "LARGE"
  enabled         = true
  vhub_prefix     = "10.129.0.0/23"
  cxp             = "USWEST-AZURE-1"
  tunnel_protocol = "VXLAN_GPE"
  group           = alkira_group.tftest.name

  # You can add more instances blocks.
  instances {
    name                     = "instance13"
    express_route_circuit_id = "/subscriptions/45060700-1949-4d0f-ba2c-4241274e8fa1/resourceGroups/robin-test/providers/Microsoft.Network/expressRouteCircuits/er-automation"
    redundant_router         = false
    loopback_subnet          = "192.168.18.0/26"
    credential_id            = alkira_credential_azure_vnet.tftest.id
  }

  segment_options {
    segment_name             = alkira_segment.tftest.name
    customer_asn             = "65514"
    disable_internet_exit    = false
    advertise_on_prem_routes = false
  }
}