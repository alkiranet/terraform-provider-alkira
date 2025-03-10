resource "alkira_connector_azure_expressroute" "example" {
  name            = "example"
  description     = "example connector"
  size            = "LARGE"
  enabled         = true
  vhub_prefix     = "10.129.0.0/23"
  cxp             = "USWEST-AZURE-1"
  tunnel_protocol = "VXLAN_GPE"
  group           = alkira_group.example.name

  # You can add more instances blocks.
  instances {
    name                    = "InstanceName"
    expressroute_circuit_id = "/subscriptions/<Id>/resourceGroups/<GroupName>/providers/Microsoft.Network/expressRouteCircuits/<CircuitName>"
    redundant_router        = false
    loopback_subnet         = "192.168.18.0/26"
    credential_id           = alkira_credential_azure_vnet.example.id
  }

  segment_options {
    segment_name             = alkira_segment.example.name
    customer_asn             = "65514"
    disable_internet_exit    = false
    advertise_on_prem_routes = false
  }
}
