---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "alkira_connector_aws_dx Resource - terraform-provider-alkira"
subcategory: ""
description: |-
  Manage AWS Direct Connect (DX) connector.
---

# alkira_connector_aws_dx (Resource)

Manage AWS Direct Connect (DX) connector.

## Example Usage

```terraform
resource "alkira_connector_aws_dx" "test" {
  name            = "example"
  description     = "example"
  cxp             = "US-WEST"
  size            = "2LARGE"
  tunnel_protocol = "GRE"
  group           = alkira_group.example.name
  billing_tag_ids = [alkira_billing_tag.example.id]


  instance {
    name          = "instance1"
    connection_id = "test-id"

    dx_asn        = 64850
    dx_gateway_ip = "169.254.199.1"

    on_prem_asn        = 65000
    on_prem_gateway_ip = "169.254.199.2"

    underlay_prefix = "169.254.199.0/30"

    bgp_auth_key        = "Alkira2018"
    bgp_auth_key_alkira = "Alkira2018"

    vlan_id       = 305
    aws_region    = "us-west-1"
    credential_id = alkira_credential_aws_vpc.example.id

    segment_options {
      segment_id          = alkira_segment.example.id
      on_prem_segment_asn = 64303

      customer_loopback_ip = "192.168.23.243"
      alkira_loopback_ip1  = "192.168.23.188"
      alkira_loopback_ip2  = "192.168.23.205"
      loopback_subnet      = "192.168.23.0/24"

      advertise_on_prem_routes = false
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cxp` (String) The CXP where the connector should be provisioned.
- `instance` (Block List, Min: 1) AWS DirectConnect (DX) instance. (see [below for nested schema](#nestedblock--instance))
- `name` (String) The name of the connector.
- `size` (String) The size of the connector, one of `SMALL`, `MEDIUM`, `LARGE`, `2LARGE`, `5LARGE` or `10LARGE`.
- `tunnel_protocol` (String) The tunnel protocol used by the connector.The value should be one of `GRE`, `IPSEC`, `VXLAN`, `VXLAN_GPE`.

### Optional

- `billing_tag_ids` (Set of Number) Billing tags to be associated with the resource. (see resource `alkira_billing_tag`).
- `description` (String) The description of the connector.
- `enabled` (Boolean) Is the connector enabled. Default is `true`.
- `group` (String) The group of the connector.

### Read-Only

- `id` (String) The ID of this resource.
- `implicit_group_id` (Number) ID of implicit group created for the connector.

<a id="nestedblock--instance"></a>
### Nested Schema for `instance`

Required:

- `aws_region` (String) AWS region of the Direct Connect.
- `connection_id` (String) AWS DirctConnect connection ID.
- `credential_id` (String) ID of AWS credential.
- `dx_asn` (Number) The ASN of AWS side of the connection.
- `name` (String) The name of the instance.
- `on_prem_asn` (Number) The customer underlay ASN.
- `segment_options` (Block Set, Min: 1) Options for each segment associated with the instance. (see [below for nested schema](#nestedblock--instance--segment_options))
- `vlan_id` (Number) This is the ID of customer facing VLAN provided by the co-location provider, configured for the link between colo provider and the customer router.

Optional:

- `bgp_auth_key` (String) The BGP MD5 authentication key forDirect Connect Gateway to verify peer.
- `bgp_auth_key_alkira` (String) The BGP MD5 authentication key forAlkira to authenticate CXP.
- `dx_gateway_ip` (String) Valid IP from underlay_prefix network used on AWS Direct Connect gateway.
- `gateway_mac_address` (String) The MAC address of the gateway.It's required if the `tunnel_protocol` is `VXLAN`.
- `on_prem_gateway_ip` (String) Valid IP from customer gateway.
- `underlay_prefix` (String) A `/30` IP prefix for on-premise gateway and DirectConnect gateway.
- `vni` (Number) Customer provided VXLAN Network Identifier (VNI). This field is required only when `tunnel_protocol` is `VXLAN`.

Read-Only:

- `id` (Number) ID of the instance.

<a id="nestedblock--instance--segment_options"></a>
### Nested Schema for `instance.segment_options`

Required:

- `loopback_subnet` (String) Prefix of all loopback IPs, helps to identify the block to reserve IPs from.
- `on_prem_segment_asn` (Number) The ASN of customer on-prem side.
- `segment_id` (String) The ID of the segment.

Optional:

- `advertise_default_route` (Boolean) Enable or disable access to the internet when traffic arrives via this connector. Default value is `false`.
- `advertise_on_prem_routes` (Boolean) Advertise on-prem routes. Default value is `false`.
- `alkira_loopback_ip1` (String) Alkira loopback IP which is set as tunnel 1. The field is applicable only when `tunnel_protocol` is not `IPSEC`.
- `alkira_loopback_ip2` (String) Alkira loopback IP which is set as tunnel 2. The field is applicable only when `tunnel_protocol` is not `IPSEC`.
- `customer_loopback_ip` (String) Customer loopback IP which is set as tunnel source. The field is applicable only when `tunnel_protocol` is not `IPSEC`.
- `number_of_customer_loopback_ips` (Number) The number of customer loopback IPs needs to be generated by Alkira from `loopback_subnet`.The field is only applicable when `tunnel_protocol` is `IPSEC`.
- `tunnel_count_per_customer_loopback_ip` (Number) The number of tunnels needs to be created for each customer loopback IP. The value must be multiple of `2` (one tunnel per AZ). The field is only applicable when `tunnel_protocol` is `IPSEC`.

## Import

Import is supported using the following syntax:

```shell
terraform import alkira_connector_aws_dx.example CONNECTOR_ID
```