resource "alkira_connector_azure_expressroute" "ha_config" {
  name            = "ha-expressroute"
  description     = "High availability ExpressRoute connector"
  size            = "LARGE"
  enabled         = true
  vhub_prefix     = "10.131.0.0/23"
  cxp             = "USEAST-AZURE-1"
  tunnel_protocol = "VXLAN"
  group           = alkira_group.core_network.name
  billing_tag_ids = [alkira_billing_tag.production.id, alkira_billing_tag.networking.id]

  instances {
    name                    = "ha-instance"
    expressroute_circuit_id = "/subscriptions/12345678-abcd-efgh-ijkl-1234567890ab/resourceGroups/network-rg/providers/Microsoft.Network/expressRouteCircuits/ha-circuit"
    redundant_router        = true # Enable redundant routers for high availability
    loopback_subnet         = "192.168.21.0/26"
    credential_id           = alkira_credential_azure_vnet.prod.id

    # MAC addresses of customer VXLAN gateways
    gateway_mac_address = ["00:1A:2B:3C:4D:5E", "00:6F:7G:8H:9I:0J"]

    # Optional virtual network interface IDs
    virtual_network_interface = [16774000, 16774001]

    segment_options {
      segment_name = alkira_segment.prod.name
      customer_gateways {
        name = "gateway1"
      }
    }

  }

  segment_options {
    segment_name             = alkira_segment.prod.name
    customer_asn             = "65002"
    disable_internet_exit    = false
    advertise_on_prem_routes = true
  }
}
