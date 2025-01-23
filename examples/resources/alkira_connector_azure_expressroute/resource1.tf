resource "alkira_connector_azure_expressroute" "example2" {
  name            = "example2"
  size            = "LARGE"
  enabled         = true
  vhub_prefix     = "10.129.0.0/23"
  cxp             = "USWEST-AZURE-1"
  tunnel_protocol = "VXLAN_GPE"
  group           = alkira_group.example.name

  instances {
    name                      = "InstanceName"
    expressroute_circuit_id   = "/subscriptions/<Id>/resourceGroups/<GroupName>/providers/Microsoft.Network/expressRouteCircuits/<CircuitName2>"
    redundant_router          = true
    loopback_subnet           = "192.168.19.0/26"
    credential_id             = alkira_credential_azure_vnet.example2.id
    gateway_mac_address       = ["00:1A:2B:3C:4D:5F", "00:6F:7G:8H:9I:0K"]
    virtual_network_interface = [16773026, 16773027]

    segment_options {
      segment_name = alkira_segment.example.name
      customer_gateways {
        name = "gateway1"
        tunnels {
          name              = "tunnel1"
          ike_version       = "v2"
          initiator         = true
          pre_shared_key    = "secretkey456"
          profile_id        = 2
          remote_auth_type  = "psk"
          remote_auth_value = "authvalue456"
        }
        tunnels {
          name              = "tunnel2"
          ike_version       = "v2"
          initiator         = false
          pre_shared_key    = "secretkey789"
          profile_id        = 3
          remote_auth_type  = "psk"
          remote_auth_value = "authvalue789"
        }
      }
      customer_gateways {
        name = "gateway2"
        tunnels {
          name              = "tunnel1"
          ike_version       = "v2"
          initiator         = true
          pre_shared_key    = "secretkeyabc"
          profile_id        = 4
          remote_auth_type  = "psk"
          remote_auth_value = "authvalueabc"
        }
      }
    }

    segment_options {
      segment_name = alkira_segment.example.name
      customer_gateways {
        name = "gateway3"
        tunnels {
          name              = "tunnel1"
          ike_version       = "v2"
          initiator         = true
          pre_shared_key    = "secretkeyxyz"
          profile_id        = 5
          remote_auth_type  = "psk"
          remote_auth_value = "authvaluexyz"
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
