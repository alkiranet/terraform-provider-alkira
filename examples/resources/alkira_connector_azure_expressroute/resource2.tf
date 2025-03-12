

resource "alkira_connector_azure_expressroute" "multi_segment" {
  name            = "multi-segment-ipsec"
  description     = "ExpressRoute connector with multiple segments and IPsec tunnels"
  size            = "2LARGE"
  enabled         = true
  vhub_prefix     = "10.132.0.0/23"
  cxp             = "USWEST-AZURE-1"
  tunnel_protocol = "IPSEC" # Using IPsec for tunnel protocol
  group           = alkira_group.security.name

  instances {
    name                    = "multi-segment-instance"
    expressroute_circuit_id = "/subscriptions/12345678-abcd-efgh-ijkl-1234567890ab/resourceGroups/network-rg/providers/Microsoft.Network/expressRouteCircuits/multi-segment-circuit"
    redundant_router        = true
    loopback_subnet         = "192.168.22.0/26"
    credential_id           = alkira_credential_azure_vnet.security.id

    # First segment options with primary gateway
    segment_options {
      segment_name = alkira_segment.dmz.name
      customer_gateways {
        name = "dmz-primary-gateway"
        tunnels {
          name              = "primary-tunnel"
          ike_version       = "IKEv2"
          initiator         = true
          pre_shared_key    = "psk-dmz-primary-123!"
          profile_id        = 10
          remote_auth_type  = "FQDN"
          remote_auth_value = "dmz-gateway.example.com"
        }
      }
    }

    # Second segment options with primary and backup gateways
    segment_options {
      segment_name = alkira_segment.internal.name

      # Primary gateway for internal segment
      customer_gateways {
        name = "internal-primary-gateway"
        tunnels {
          name              = "primary-tunnel"
          ike_version       = "IKEv2"
          initiator         = true
          pre_shared_key    = "psk-internal-primary-456!"
          profile_id        = 11
          remote_auth_type  = "FQDN"
          remote_auth_value = "internal-primary.example.com"
        }
      }

      # Backup gateway for internal segment
      customer_gateways {
        name = "internal-backup-gateway"
        tunnels {
          name              = "backup-tunnel"
          ike_version       = "IKEv2"
          initiator         = false # Waiting for initiation from customer side
          pre_shared_key    = "psk-internal-backup-789!"
          profile_id        = 12
          remote_auth_type  = "FQDN"
          remote_auth_value = "internal-backup.example.com"
        }
      }
    }
  }

  # Global segment options for DMZ
  segment_options {
    segment_name             = alkira_segment.dmz.name
    customer_asn             = "65003"
    disable_internet_exit    = true
    advertise_on_prem_routes = false
  }

  # Global segment options for internal
  segment_options {
    segment_name             = alkira_segment.internal.name
    customer_asn             = "65004"
    disable_internet_exit    = false
    advertise_on_prem_routes = true
  }
}
