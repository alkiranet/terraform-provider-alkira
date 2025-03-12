resource "alkira_connector_azure_expressroute" "basic" {
  name            = "basic-expressroute"
  description     = "Basic ExpressRoute connector with VXLAN_GPE"
  size            = "MEDIUM"
  enabled         = true
  vhub_prefix     = "10.130.0.0/23"
  cxp             = "USWEST-AZURE-1"
  tunnel_protocol = "VXLAN_GPE" # Default tunnel protocol
  group           = alkira_group.networking.name

  instances {
    name                    = "primary-instance"
    expressroute_circuit_id = "/subscriptions/12345678-abcd-efgh-ijkl-1234567890ab/resourceGroups/network-rg/providers/Microsoft.Network/expressRouteCircuits/primary-circuit"
    redundant_router        = false # Single router configuration
    loopback_subnet         = "192.168.20.0/26"
    credential_id           = alkira_credential_azure_vnet.prod.id
    segment_options {
      segment_name = alkira_segment.prod.name
      customer_gateways {
        name = "gateway1"
      }
    }
  }

  segment_options {
    segment_name             = alkira_segment.prod.name
    customer_asn             = "65001"
    disable_internet_exit    = true
    advertise_on_prem_routes = true
  }
}
