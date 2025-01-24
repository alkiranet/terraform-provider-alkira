resource "alkira_connector_azure_expressroute" "example" {
  name            = "example"
  description     = "example connector"
  size            = "LARGE"
  enabled         = true
  vhub_prefix     = "10.129.0.0/23"
  cxp             = "USWEST-AZURE-1"
  tunnel_protocol = "VXLAN_GPE"
  group           = alkira_group.example.name

  instances {
    name                      = "InstanceName"
    expressroute_circuit_id   = "/subscriptions/<Id>/resourceGroups/<GroupName>/providers/Microsoft.Network/expressRouteCircuits/<CircuitName>"
    redundant_router          = true
    loopback_subnet           = "192.168.18.0/26"
    credential_id             = alkira_credential_azure_vnet.example.id
    gateway_mac_address       = ["00:1A:2B:3C:4D:5E", "00:6F:7G:8H:9I:0J"]
    virtual_network_interface = [16773024, 16773025]

    segment_options {
      segment_name = alkira_segment.example.name
      customer_gateways {
        name = "gateway1"
        tunnels {
          name              = "tunnel1"
          ike_version       = "IKEv2"
          initiator         = true
          pre_shared_key    = "secretkey123"
          profile_id        = 1
          remote_auth_type  = "FQDN"
          remote_auth_value = "authvalue123"
        }
      }
    }
  }

  segment_options {
    segment_name             = alkira_segment.example.name
    customer_asn             = "65514"
    disable_internet_exit    = false
    advertise_on_prem_routes = false
  }
}
