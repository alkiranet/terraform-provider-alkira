---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "alkira_connector_versa_sdwan Resource - terraform-provider-alkira"
subcategory: ""
description: |-
  Manage Versa SD-WAN Connector. (BETA)
---

# alkira_connector_versa_sdwan (Resource)

Manage Versa SD-WAN Connector. (**BETA**)

## Example Usage

```terraform
resource "alkira_connector_versa_sdwan" "test" {
  name    = "test"
  cxp     = "US-WEST"
  group   = alkira_group.test.name
  size    = "SMALL"

  versa_controller_host = "172.16.0.1"
  local_id  = 1
  remote_id = 2

  versa_vos_device {
    hostname                   = "dev1"
    local_device_serial_number = "12345678"
    version                    = "21.2.3-B"
  }

  vrf_segment_mapping {
    segment_id     = alkira_segment.test.id
    vrf_name       = "test"
    versa_bgp_asn  = 1203403435
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cxp` (String) The CXP where the connector should be provisioned.
- `local_id` (String) The local ID.
- `name` (String) The name of the connector.
- `remote_id` (String) The remote ID.
- `size` (String) The size of the connector, one of `SMALL`, `MEDIUM`, `LARGE`, `2LARGE`, `5LARGE`.
- `versa_controller_host` (String) The Versa controller IP/FQDN.
- `versa_vos_device` (Block List, Min: 1) Versa VOS Device. (see [below for nested schema](#nestedblock--versa_vos_device))
- `vrf_segment_mapping` (Block Set, Min: 1) Specify target segment for VRF. (see [below for nested schema](#nestedblock--vrf_segment_mapping))

### Optional

- `billing_tag_ids` (Set of Number) Billing tags to be associated with the resource. (see resource `alkira_billing_tag`).
- `description` (String) The description of the connector.
- `enabled` (Boolean) Is the connector enabled. Default value is `true`.
- `global_tenant_id` (Number) The global tenant ID of Versa SD-WAN. Default value is `1`.
- `group` (String) The group of the connector.
- `local_public_shared_key` (String) The local public shared key. Default value is`1234`.
- `remote_public_shared_key` (String) The remote public shared key. Default value is`1234`.
- `tunnel_protocol` (String) The tunnel protocol of Versa SD-WAN.

### Read-Only

- `id` (String) The ID of this resource.
- `implicit_group_id` (Number) The ID of implicit group automaticaly created with the connector.
- `provision_state` (String) The provision state of the connector.

<a id="nestedblock--versa_vos_device"></a>
### Nested Schema for `versa_vos_device`

Required:

- `hostname` (String) The hostname of the VOS Device.
- `local_device_serial_number` (String) Local device serial number.
- `version` (String) Versa version.

Read-Only:

- `id` (Number) The ID of the VOS device.


<a id="nestedblock--vrf_segment_mapping"></a>
### Nested Schema for `vrf_segment_mapping`

Required:

- `segment_id` (Number) Segment ID.
- `versa_bgp_asn` (Number) BGP ASN on the Versa. A typical value for 2 byte segment is `64523` and `4200064523` for 4 byte segment.
- `vrf_name` (String) VRF Name.

Optional:

- `advertise_default_route` (Boolean) Whether advertise default route of internet connector. Default value is `false`.
- `advertise_on_prem_routes` (Boolean) Advertise On Prem Routes. Default value is `false`.

## Import

Import is supported using the following syntax:

```shell
terraform import alkira_connector_versa_sdwan.example CONNECTOR_ID
```
