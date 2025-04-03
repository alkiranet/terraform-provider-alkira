

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

    # First ipsec customer gateway with primary gateway
    ipsec_customer_gateway {
      segment_id = alkira_segment.dmz.id
      customer_gateway {
        name = "dmz-primary-gateway"
        tunnel {
          name              = "primary-tunnel"
          ike_version       = "IKEv2"
          initiator         = true
          pre_shared_key    = "psk-dmz-primary-123!"
          profile_id        = 10
          remote_auth_type  = "FQDN"
          remote_auth_value = "dmz-gateway.example.com"
        }
      }
      customer_gateway {
        name = "dmz-backup-gateway"
        tunnel {
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

    # Second ipsec customer gateway with primary and backup gateways
    ipsec_customer_gateway {
      segment_id = alkira_segment.internal.id

      # Primary gateway for internal segment
      customer_gateway {
        name = "internal-primary-gateway"
        tunnel {
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
      customer_gateway {
        name = "internal-backup-gateway"
        tunnel {
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
    segment_id               = alkira_segment.dmz.id
    customer_asn             = "65003"
    disable_internet_exit    = true
    advertise_on_prem_routes = false
  }

  # Global segment options for internal
  segment_options {
    segment_id               = alkira_segment.internal.id
    customer_asn             = "65004"
    disable_internet_exit    = false
    advertise_on_prem_routes = true
  }
}
