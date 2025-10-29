---
page_title: "alkira_connector_ipsec Resource - terraform-provider-alkira"
subcategory: ""
description: |-
  Manage IPSec Connector.
---

# alkira_connector_ipsec (Resource)

Manage IPSec Connector.

## VPN Mode

`vpn_mode` can be either `ROUTE_BASED` or `POLICY_BASED`:

- **`ROUTE_BASED`**: Uses routing protocols to determine traffic flow. Requires `routing_options` block with one of the following routing types:
  - `STATIC`: Static routes defined via prefix lists
  - `DYNAMIC`: BGP routing with customer gateway ASN
  - `BOTH`: Combination of static and dynamic routing

- **`POLICY_BASED`**: Uses traffic selectors (prefix lists) to determine which traffic should be encrypted. Requires `policy_options` block with both `on_prem_prefix_list_ids` and `cxp_prefix_list_ids`.

## Routing Options

The `routing_options` block supports the following availability check methods:
- `IKE_STATUS`: Uses IKE tunnel status to determine route availability
- `IPSEC_INTERFACE_PING`: Pings the IPSec interface to verify connectivity (default)
- `PING`: Simple ping-based availability check

For dynamic routing, you can optionally specify `bgp_auth_key` for BGP MD5 authentication.

## High Availability

Multiple endpoints can be configured for redundancy:
- `ha_mode = "ACTIVE"`: Endpoint actively carries traffic
- `ha_mode = "STANDBY"`: Endpoint only used when all active endpoints are down
- `enable_tunnel_redundancy`: When `true`, all tunnels must be UP for the connector to be considered healthy

## Example Usage

### Basic Route-Based with Dynamic Routing (BGP)

```terraform
resource "alkira_connector_ipsec" "basic_dynamic" {
  name        = "ipsec-connector-basic"
  description = "Basic route-based IPSec connector with BGP"
  cxp         = "US-WEST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "SMALL"
  enabled     = true

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type                 = "DYNAMIC"
    customer_gateway_asn = "65310"
    availability         = "IPSEC_INTERFACE_PING"
  }

  endpoint {
    name                = "remote-site"
    customer_gateway_ip = "203.0.113.1"
    preshared_keys      = ["your-preshared-key-here"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
  }
}
```

### Route-Based with Static Routing

```terraform
resource "alkira_connector_ipsec" "route_based_static" {
  name        = "ipsec-connector-static"
  description = "Route-based IPSec connector with static routes"
  cxp         = "US-EAST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "SMALL"
  enabled     = true

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type           = "STATIC"
    prefix_list_id = alkira_list_global_cidr.remote_subnets.id
    availability   = "IKE_STATUS"
  }

  endpoint {
    name                = "branch-office"
    customer_gateway_ip = "203.0.113.10"
    preshared_keys      = ["branch-key-1", "branch-key-2"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
  }
}
```

### Hybrid Routing (Both Static and Dynamic)

```terraform
resource "alkira_connector_ipsec" "hybrid_routing" {
  name        = "ipsec-connector-hybrid"
  description = "IPSec connector with both static and dynamic routing"
  cxp         = "US-WEST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "MEDIUM"
  enabled     = true

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type                 = "BOTH"
    prefix_list_id       = alkira_list_global_cidr.remote_subnets.id
    customer_gateway_asn = "65320"
    availability         = "IPSEC_INTERFACE_PING"
  }

  endpoint {
    name                = "hybrid-site"
    customer_gateway_ip = "203.0.113.20"
    preshared_keys      = ["hybrid-key"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
  }
}
```

### BGP with MD5 Authentication

```terraform
resource "alkira_connector_ipsec" "bgp_auth" {
  name        = "ipsec-connector-bgp-auth"
  description = "IPSec connector with BGP MD5 authentication"
  cxp         = "US-WEST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "SMALL"
  enabled     = true

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type                 = "DYNAMIC"
    customer_gateway_asn = "65330"
    bgp_auth_key         = "my-bgp-secret-key"
    availability         = "PING"
  }

  endpoint {
    name                = "secured-site"
    customer_gateway_ip = "203.0.113.30"
    preshared_keys      = ["secured-psk"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
  }
}
```

### Policy-Based IPSec

When using `vpn_mode = "POLICY_BASED"`, you must specify traffic selectors for both on-premises and CXP networks:

```terraform
resource "alkira_connector_ipsec" "policy_based" {
  name        = "ipsec-connector-policy-based"
  description = "Policy-based IPSec connector with traffic selectors"
  cxp         = "US-WEST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "SMALL"
  enabled     = true

  vpn_mode = "POLICY_BASED"

  policy_options {
    on_prem_prefix_list_ids = [alkira_list_global_cidr.on_prem_subnets.id]
    cxp_prefix_list_ids     = [alkira_list_global_cidr.cxp_subnets.id]
  }

  endpoint {
    name                = "policy-site"
    customer_gateway_ip = "203.0.113.40"
    preshared_keys      = ["policy-key"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
  }
}
```

### High Availability with Active and Standby Endpoints

```terraform
resource "alkira_connector_ipsec" "ha_active_standby" {
  name        = "ipsec-connector-ha"
  description = "High availability IPSec with active and standby endpoints"
  cxp         = "US-WEST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "MEDIUM"
  enabled     = true

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type                 = "DYNAMIC"
    customer_gateway_asn = "65340"
  }

  endpoint {
    name                     = "primary-active"
    customer_gateway_ip      = "203.0.113.50"
    preshared_keys           = ["primary-key-1", "primary-key-2"]
    billing_tag_ids          = [alkira_billing_tag.tag1.id]
    ha_mode                  = "ACTIVE"
    enable_tunnel_redundancy = true
  }

  endpoint {
    name                     = "secondary-active"
    customer_gateway_ip      = "203.0.113.51"
    preshared_keys           = ["secondary-key-1", "secondary-key-2"]
    billing_tag_ids          = [alkira_billing_tag.tag1.id]
    ha_mode                  = "ACTIVE"
    enable_tunnel_redundancy = true
  }

  endpoint {
    name                = "standby"
    customer_gateway_ip = "203.0.113.52"
    preshared_keys      = ["standby-key"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
    ha_mode             = "STANDBY"
  }
}
```

### Dynamic Customer Gateway IP

For sites with dynamic public IP addresses, set `customer_ip_type = "DYNAMIC"` and `customer_gateway_ip = "0.0.0.0"`. The `advanced_options` block with `remote_auth_type` is required, and `initiator` must be `true`:

```terraform
resource "alkira_connector_ipsec" "dynamic_gateway" {
  name        = "ipsec-connector-dynamic-gw"
  description = "IPSec connector with dynamic customer gateway IP"
  cxp         = "US-WEST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "SMALL"
  enabled     = true

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type                 = "DYNAMIC"
    customer_gateway_asn = "65350"
  }

  endpoint {
    name                = "dynamic-site"
    customer_gateway_ip = "0.0.0.0"
    customer_ip_type    = "DYNAMIC"
    preshared_keys      = ["dynamic-key"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]

    advanced_options {
      esp_dh_group_numbers      = ["MODP2048"]
      esp_encryption_algorithms = ["AES256CBC"]
      esp_integrity_algorithms  = ["SHA256"]

      ike_dh_group_numbers      = ["MODP2048"]
      ike_encryption_algorithms = ["AES256CBC"]
      ike_integrity_algorithms  = ["SHA256"]
      ike_version               = "IKEv2"

      initiator = true

      remote_auth_type  = "FQDN"
      remote_auth_value = "remote-site.example.com"
    }
  }
}
```

### Advanced Cryptographic Options

Customize ESP and IKE encryption, integrity, and Diffie-Hellman algorithms:

```terraform
resource "alkira_connector_ipsec" "advanced_crypto" {
  name        = "ipsec-connector-advanced-crypto"
  description = "IPSec connector with custom cryptographic algorithms"
  cxp         = "US-WEST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "SMALL"
  enabled     = true

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type                 = "DYNAMIC"
    customer_gateway_asn = "65360"
  }

  endpoint {
    name                = "crypto-site"
    customer_gateway_ip = "203.0.113.60"
    preshared_keys      = ["crypto-key"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]

    advanced_options {
      esp_dh_group_numbers      = ["MODP4096", "ECP384"]
      esp_encryption_algorithms = ["AES256GCM16", "AES256CBC"]
      esp_integrity_algorithms  = ["SHA512", "SHA384"]

      ike_dh_group_numbers      = ["MODP4096", "ECP384"]
      ike_encryption_algorithms = ["AES256CBC", "AES192CBC"]
      ike_integrity_algorithms  = ["SHA512", "SHA384"]
      ike_version               = "IKEv1"

      initiator = false

      remote_auth_type  = "IP_ADDR"
      remote_auth_value = "203.0.113.60"
    }
  }
}
```

### Segment Options

Control internet access and route advertisement per segment:

```terraform
resource "alkira_connector_ipsec" "segment_options" {
  name        = "ipsec-connector-segment-opts"
  description = "IPSec connector with segment-specific options"
  cxp         = "US-WEST"
  group       = alkira_group.group1.name
  segment_id  = alkira_segment.segment1.id
  size        = "SMALL"
  enabled     = true

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type                 = "DYNAMIC"
    customer_gateway_asn = "65370"
  }

  segment_options {
    name                     = alkira_segment.segment1.name
    advertise_default_route  = true
    advertise_on_prem_routes = true
  }

  endpoint {
    name                = "segment-site"
    customer_gateway_ip = "203.0.113.70"
    preshared_keys      = ["segment-key"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
  }
}
```

### Multi-Site with All Features

Complex deployment with scale group, multiple sites, and mixed configurations:

```terraform
resource "alkira_connector_ipsec" "multi_site_advanced" {
  name           = "ipsec-connector-multi-site"
  description    = "Multi-site IPSec connector with scale group"
  cxp            = "US-WEST"
  group          = alkira_group.group1.name
  segment_id     = alkira_segment.segment1.id
  size           = "LARGE"
  enabled        = true
  scale_group_id = alkira_scale_group.ipsec_scale.id

  vpn_mode = "ROUTE_BASED"

  routing_options {
    type                 = "BOTH"
    prefix_list_id       = alkira_list_global_cidr.remote_subnets.id
    customer_gateway_asn = "65380"
    bgp_auth_key         = "multi-site-bgp-key"
    availability         = "IPSEC_INTERFACE_PING"
  }

  segment_options {
    name                     = alkira_segment.segment1.name
    advertise_default_route  = false
    advertise_on_prem_routes = true
  }

  endpoint {
    name                     = "site1-primary"
    customer_gateway_ip      = "203.0.113.80"
    preshared_keys           = ["site1-key-1", "site1-key-2"]
    billing_tag_ids          = [alkira_billing_tag.tag1.id]
    ha_mode                  = "ACTIVE"
    enable_tunnel_redundancy = true

    advanced_options {
      esp_dh_group_numbers      = ["MODP3072", "ECP256"]
      esp_encryption_algorithms = ["AES256CBC", "AES256GCM16"]
      esp_integrity_algorithms  = ["SHA256", "SHA384"]

      ike_dh_group_numbers      = ["MODP3072", "ECP256"]
      ike_encryption_algorithms = ["AES256CBC"]
      ike_integrity_algorithms  = ["SHA256"]
      ike_version               = "IKEv2"

      initiator = true

      remote_auth_type  = "FQDN"
      remote_auth_value = "site1.example.com"
    }
  }

  endpoint {
    name                     = "site2-secondary"
    customer_gateway_ip      = "203.0.113.81"
    preshared_keys           = ["site2-key-1", "site2-key-2"]
    billing_tag_ids          = [alkira_billing_tag.tag2.id]
    ha_mode                  = "ACTIVE"
    enable_tunnel_redundancy = false

    advanced_options {
      esp_dh_group_numbers      = ["MODP2048"]
      esp_encryption_algorithms = ["AES128CBC"]
      esp_integrity_algorithms  = ["SHA1"]

      ike_dh_group_numbers      = ["MODP2048"]
      ike_encryption_algorithms = ["AES128CBC"]
      ike_integrity_algorithms  = ["SHA1"]
      ike_version               = "IKEv1"

      initiator = true

      remote_auth_type  = "KEYID"
      remote_auth_value = "site2-identifier"
    }
  }

  endpoint {
    name                = "site3-backup"
    customer_gateway_ip = "203.0.113.82"
    preshared_keys      = ["site3-key"]
    billing_tag_ids     = [alkira_billing_tag.tag1.id]
    ha_mode             = "STANDBY"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cxp` (String) The CXP where the connector should be provisioned.
- `endpoint` (Block List, Min: 1) The endpoint. (see [below for nested schema](#nestedblock--endpoint))
- `name` (String) The name of the connector.
- `segment_id` (String) The ID of the segment associated with the connector.
- `size` (String) The size of the connector, one of `SMALL`, `MEDIUM`, `LARGE`, `2LARGE`, `5LARGE`, `10LARGE`.
- `vpn_mode` (String) The mode can be configured either as `ROUTE_BASED` or `POLICY_BASED`.

### Optional

- `description` (String) The description of the connector.
- `enabled` (Boolean) Is the connector enabled. Default is `true`.
- `group` (String) The group of the connector. (see resource `alkira_group`)
- `policy_options` (Block Set) Policy options, both `on_prem_prefix_list_ids` and `cxp_prefix_list_ids` must be provided if `vpn_mode` is `POLICY_BASED`. (see [below for nested schema](#nestedblock--policy_options))
- `routing_options` (Block Set) Routing options, type is `STATIC`, `DYNAMIC`, or`BOTH` must be provided if `vpn_mode` is `ROUTE_BASED` (see [below for nested schema](#nestedblock--routing_options))
- `scale_group_id` (String) The ID of the scale group associated with the connector.
- `segment_options` (Block Set) Additional options for each segment associated with the connector. (see [below for nested schema](#nestedblock--segment_options))

### Read-Only

- `id` (String) The ID of this resource.
- `implicit_group_id` (Number) The ID of implicit group automaticaly created with the connector.
- `provision_state` (String) The provision state of the connector.

<a id="nestedblock--endpoint"></a>
### Nested Schema for `endpoint`

Required:

- `customer_gateway_ip` (String) The IP address of the customer gateway.
- `name` (String) The name of the endpoint.
- `preshared_keys` (List of String) An array of preshared keys, one per tunnel. The value needs to be provided explictly.

Optional:

- `advanced_options` (Block List) Advanced options for IPSec endpoint. (see [below for nested schema](#nestedblock--endpoint--advanced_options))
- `billing_tag_ids` (Set of Number) Billing tags to be associated with the resource. (see resource `alkira_billing_tag`).
- `customer_ip_type` (String) The type of `customer_gateway_ip`. It could be either `STATIC` or `DYNAMIC`. Default value is `STATIC`. When it's `DYNAMIC`, `customer_gateway_ip` should be set to `0.0.0.0`. `remote_auth_type` in `advanced_options` is required as well.
- `enable_tunnel_redundancy` (Boolean) Disable this if all tunnels will not be configured or enabled on the on-premise device. If it's set to `false`, connector health will be shown as `UP` if at least one of the tunnels is `UP`. If enabled, all tunnels need to be `UP` for the connectorhealth to be shown as `UP`.
- `ha_mode` (String) The value could be `ACTIVE` or `STANDBY`. A endpoint in `STANDBY` mode will not be used for traffic unless all other endpoints for the connector are down. There can only be one endpoint in `STANDBY` mode per connector and there must be at least one endpoint that isn't in `STANDBY` mode per connector.

Read-Only:

- `id` (Number) The ID of the endpoint.

<a id="nestedblock--endpoint--advanced_options"></a>
### Nested Schema for `endpoint.advanced_options`

Required:

- `esp_dh_group_numbers` (List of String) Diffie Hellman groups to use for IPsec SA. Value could `MODP1024`, `MODP2048`, `MODP3072`, `MODP4096`, `MODP6144`, `MODP8192`, `ECP256`, `ECP384`, `ECP521`, `CURVE25519` and `NONE`.
- `esp_encryption_algorithms` (List of String) Encryption algorithms to use for IPsec SA. Value could be `AES256CBC`, `AES192CBC`, `AES128CBC`, `AES256GCM16` `3DESCBC`, or `NULL`.
- `esp_integrity_algorithms` (List of String) Integrity algorithms to use for IPsec SA. Value could `SHA1`, `SHA256`, `SHA384`, `SHA512` or `MD5`.
- `ike_dh_group_numbers` (List of String) Diffie Hellman groups to use for IKE SA, one of `MODP1024`, `MODP2048`, `MODP3072`, `MODP4096`, `MODP6144`, `MODP8192`, `ECP256`, `ECP384`, `ECP521`, or `CURVE25519`.
- `ike_encryption_algorithms` (List of String) Encryption algorithms to use for IKE SA, one of `AES256CBC`, `AES192CBC`, `AES128CBC` and `3DESCBC`.
- `ike_integrity_algorithms` (List of String) Integrity algorithms to use for IKE SA, one of `SHA1`, `SHA256`, `SHA384`, `SHA512`.
- `ike_version` (String) IKE version, either `IKEv1` or `IKEv2`
- `initiator` (Boolean) When the value is `false`, CXP will initiate the IKE connection and in all other cases the customer gateway should initiate IKE connection. When `gateway_ip_type` is `DYNAMIC`, initiator must be `true`.
- `remote_auth_type` (String) IKE identity to use for authentication round, one of `FQDN`, `USER_FQDN`, `KEYID`, or `IP_ADDR`.
- `remote_auth_value` (String) Remote-ID value.



<a id="nestedblock--policy_options"></a>
### Nested Schema for `policy_options`

Required:

- `cxp_prefix_list_ids` (List of Number) CXP Prefix List IDs.
- `on_prem_prefix_list_ids` (List of Number) On Prem Prefix List IDs.


<a id="nestedblock--routing_options"></a>
### Nested Schema for `routing_options`

Required:

- `type` (String) Routing type, one of `STATIC`, `DYNAMIC`, or `BOTH`.

Optional:

- `availability` (String) The method to determine the availability of the routes. The value could be `IKE_STATUS` or `IPSEC_INTERFACE_PING`. Default value is `IPSEC_INTERFACE_PING`.
- `bgp_auth_key` (String) BGP MD5 auth key for Alkira to authenticate Alkira CXP (On Premise Gateway).
- `customer_gateway_asn` (String) The customer gateway ASN to use for dynamic route propagation.
- `prefix_list_id` (Number) The ID of prefix list to use for static route propagation.


<a id="nestedblock--segment_options"></a>
### Nested Schema for `segment_options`

Required:

- `name` (String) Segment Name.

Optional:

- `advertise_default_route` (Boolean) Enable or disable access to the internet when traffic arrives via this connector. Default is `false`.
- `advertise_on_prem_routes` (Boolean) Additional options for each segment associated with the connector. Default is `false`.

## Import

Import is supported using the following syntax:

```shell
terraform import alkira_connector_ipsec.example CONNECTOR_ID
```
