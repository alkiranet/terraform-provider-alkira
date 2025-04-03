resource "alkira_connector_azure_expressroute" "multi_instance" {
  name            = "multi-instance-connector"
  description     = "ExpressRoute connector with multiple circuit instances"
  size            = "5LARGE"
  enabled         = true
  vhub_prefix     = "10.133.0.0/23"
  cxp             = "USEAST-AZURE-1"
  tunnel_protocol = "IPSEC"
  group           = alkira_group.global_network.name

  # First ExpressRoute circuit instance
  instances {
    name                    = "primary-circuit"
    expressroute_circuit_id = "/subscriptions/12345678-abcd-efgh-ijkl-1234567890ab/resourceGroups/network-rg/providers/Microsoft.Network/expressRouteCircuits/primary-circuit"
    redundant_router        = false
    loopback_subnet         = "192.168.23.0/26"
    credential_id           = alkira_credential_azure_vnet.primary.id

    iposec_customer_gateway {
      segment_id = alkira_segment.prod.id
      customer_gateway {
        name = "prod-gateway"
        tunnel {
          name              = "prod-tunnel"
          ike_version       = "IKEv2"
          initiator         = true
          pre_shared_key    = "psk-prod-primary-123!"
          profile_id        = 20
          remote_auth_type  = "FQDN"
          remote_auth_value = "prod-gateway.example.com"
        }
      }
    }
  }

  # Second ExpressRoute circuit instance
  instances {
    name                    = "backup-circuit"
    expressroute_circuit_id = "/subscriptions/12345678-abcd-efgh-ijkl-1234567890ab/resourceGroups/network-rg/providers/Microsoft.Network/expressRouteCircuits/backup-circuit"
    redundant_router        = false
    loopback_subnet         = "192.168.24.0/26"
    credential_id           = alkira_credential_azure_vnet.backup.id

    ipsec_customer_gateway {
      segment_id = alkira_segment.prod.id
      customer_gateway {
        name = "backup-gateway"
        tunnel {
          name              = "backup-tunnel"
          ike_version       = "IKEv2"
          initiator         = false
          pre_shared_key    = "psk-backup-456!"
          profile_id        = 21
          remote_auth_type  = "FQDN"
          remote_auth_value = "backup-gateway.example.com"
        }
      }
    }
  }

  # Global segment options
  segment_options {
    segment_id               = alkira_segment.prod.id
    customer_asn             = "65005"
    disable_internet_exit    = false
    advertise_on_prem_routes = true
  }
}
