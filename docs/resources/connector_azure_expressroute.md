---
page_title: "alkira_connector_azure_expressroute Resource - terraform-provider-alkira"
subcategory: ""
description: |-
  Manage Azure ExpressRoute Connector. (BETA)
---

# alkira_connector_azure_expressroute (Resource)

Manage Azure ExpressRoute Connector. (**BETA**)


This example demonstrates a straightforward Azure ExpressRoute connector configuration with minimal settings
```terraform
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
  }

  segment_options {
    segment_id               = alkira_segment.prod.id
    customer_asn             = "65001"
    disable_internet_exit    = true
    advertise_on_prem_routes = true
  }
}
```

This example demonstrates a high-availability configuration with redundant routers and multiple gateway MAC addresses for VXLAN connectivity
```terraform
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

  }

  segment_options {
    segment_id               = alkira_segment.prod.id
    customer_asn             = "65002"
    disable_internet_exit    = false
    advertise_on_prem_routes = true
  }
}
```

This example demonstrates a configuration with multiple segments and IPsec tunnels with various authentication and security settings
```terraform
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
```

This example demonstrates configuring multiple ExpressRoute circuit instances within a single connector resource
```terraform
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
```


## Understanding Connector Components

### 1. Connector Base Configuration

The connector requires basic information such as:
- Name and description
- Size (from SMALL to 10LARGE)
- CXP location (where the connector is provisioned)
- VHUB prefix (a /23 CIDR block for the virtual hub)
- Tunnel protocol (VXLAN, VXLAN_GPE, or IPSEC)

### 2. ExpressRoute Circuit Instances

Each connector must have at least one ExpressRoute circuit instance, which specifies:
- Circuit identifier from Azure
- Loopback subnet (/26) for establishing VXLAN GPE tunnels
- Azure credentials
- Optional redundant router configuration
- For VXLAN: Gateway MAC addresses and optional virtual network interfaces
- For IPsec: Segment-specific gateway and tunnel configurations

### 3. Segment Options

Segment options define routing parameters for each network segment, including:
- Customer ASN (Autonomous System Number)
- Internet exit controls
- On-premises route advertisement options

## Important Notes

- The VHUB prefix must be a `/23` CIDR block
- The loopback subnet must be a `/26` CIDR block
- For VXLAN with redundant routers, two gateway MAC addresses are required
- For IPsec tunnels, at least one customer gateway with one tunnel is required per segment
- Currently, only IKEv2 and FQDN authentication are supported for IPsec tunnels
- **Segment mapping requirement**: There must be a one-to-one correspondence between segments defined in the `ipsec_customer_gateway` and global-level `segment_options`.
  Every segment referenced within an instance must also have a corresponding global segment options entry with the same segment name.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cxp` (String) The CXP where the connector should be provisioned.
- `instances` (Block List, Min: 1) (see [below for nested schema](#nestedblock--instances))
- `name` (String) The name of the connector.
- `segment_options` (Block List, Min: 1) (see [below for nested schema](#nestedblock--segment_options))
- `size` (String) The size of the connector, one of `SMALL`, `MEDIUM`, `LARGE`, `2LARGE`, `5LARGE`, `10LARGE`.
- `vhub_prefix` (String) IP address prefix for VWAN Hub. This should be a `/23` prefix.

### Optional

- `billing_tag_ids` (Set of Number) Billing tags to be associated with the resource. (see resource `alkira_billing_tag`).
- `description` (String) The description of the connector.
- `enabled` (Boolean) Is the connector enabled. Default is `true`.
- `group` (String) The group of the connector.
- `tunnel_protocol` (String) The tunnel protocol. One of `VXLAN`, `VXLAN_GPE`, `IPSEC`. Default is `VXLAN_GPE`

### Read-Only

- `id` (String) The ID of this resource.
- `provision_state` (String) The provision state of the connector.

<a id="nestedblock--instances"></a>
### Nested Schema for `instances`

Required:

- `credential_id` (String) An opaque identifier generated when storing Azure VNET credentials.
- `expressroute_circuit_id` (String) ExpressRoute circuit ID from Azure. ExpressRoute Circuit should have a private peering connection provisioned, also an valid authorization key associated with it.
- `loopback_subnet` (String) A `/26` subnet from which loopback IPs would be used to establish underlay VXLAN GPE tunnels.
- `name` (String) User provided connector instance name.

Optional:

- `gateway_mac_address` (List of String) An array containing the mac addresses of VXLAN gateways reachable through ExpressRoute circuit. The field is only expected if VXLAN tunnel protocol is selected, and 2 gateway MAC addresses are expected only if `redundant_router` is enabled.
- `ipsec_customer_gateway` (Block List) IPSec customer gateway configuration. The block is only required when tunnel_protocol is IPSEC. All segments defined in the segment_options should be configured here as well. (see [below for nested schema](#nestedblock--instances--ipsec_customer_gateway))
- `redundant_router` (Boolean) Indicates if ExpressRoute Circuit terminates on redundant routers on customer side.
- `virtual_network_interface` (List of Number) This is an optional field if the `tunnel_protocol` is `VXLAN`. If not specified Alkira allocates unique VNI from the range `[16773023, 16777215]`.

Read-Only:

- `id` (Number)

<a id="nestedblock--instances--ipsec_customer_gateway"></a>
### Nested Schema for `instances.ipsec_customer_gateway`

Required:

- `customer_gateway` (Block List, Min: 1) Customer gateway configurations for `IPSEC` tunnels. (see [below for nested schema](#nestedblock--instances--ipsec_customer_gateway--customer_gateway))
- `segment_id` (String) The ID of a segment.

<a id="nestedblock--instances--ipsec_customer_gateway--customer_gateway"></a>
### Nested Schema for `instances.ipsec_customer_gateway.customer_gateway`

Required:

- `name` (String) A unique name for the customer gateway.
- `tunnel` (Block List, Min: 1) Tunnel configurations for the gateway. At least one tunnel is required for `IPSEC`. (see [below for nested schema](#nestedblock--instances--ipsec_customer_gateway--customer_gateway--tunnel))

<a id="nestedblock--instances--ipsec_customer_gateway--customer_gateway--tunnel"></a>
### Nested Schema for `instances.ipsec_customer_gateway.customer_gateway.tunnel`

Required:

- `name` (String) A unique name for the tunnel.

Optional:

- `ike_version` (String) The IKE protocol version. Currently, only `IKEv2` is supported.
- `initiator` (Boolean) Whether this endpoint initiates the tunnel connection. Default value is `true`.
- `pre_shared_key` (String, Sensitive) The pre-shared key for tunnel authentication. This field is sensitive and will not be displayed in logs.
- `profile_id` (Number) The ID of the IPSec Tunnel Profile (`connector_ipsec_tunnel_profile`).
- `remote_auth_type` (String) The authentication type for the remote endpoint. Only `FQDN` iscurrently supported.
- `remote_auth_value` (String, Sensitive) The authentication value for the remote endpoint. This field is sensitive.

Read-Only:

- `id` (String) The ID of the tunnel.





<a id="nestedblock--segment_options"></a>
### Nested Schema for `segment_options`

Required:

- `customer_asn` (Number) ASN on the customer premise side.
- `segment_id` (String) The ID of the segment.

Optional:

- `advertise_on_prem_routes` (Boolean) Allow routes from the branch/premises to be advertised to the cloud.
- `disable_internet_exit` (Boolean) Enable or disable access to the internet when traffic arrives via this connector.
