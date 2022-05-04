---
subcategory: "Release Notes"
page_title: "v0.9.1"
description: |-
    Release notes for v0.9.1
---

This release brings new resources, enhancements and bug fixes. The
main focus is to bring support for all available features up to date.

**NOTE:** There are several changes in `alkira_service_pan` and
`alkira_connector_ipsec` that require updates to the existing
configuration. Please check the resource's documentation before
attempting an upgrade.


## New Resources

### resource `alkira_connector_aruba_edge` (**NEW**)

New connector for connecting with Aruba Edge.

### resource `alkira_service_checkpoint` (**NEW**)

New service integration with Checkpoint.

### resource `alkira_service_zscaler` (**NEW**)

New service integration with Zscaler.


## Enhancements

Several enhancements may need configuration updates. Please refer the
resource documentation for more details.


### `size` is updated across connectors and services

    There are more available `size` now across many connectors and
    services: `4LARGE`, `5LARGE`, `10LARGE` and `20LARGE`.


### resource `alkira_connector_aws_vpc`

    New `tgw_attachment` block was added to the connector.


### resource `alkira_connector_azure_vent`

    Introduce two more blocks `vnet_subnet` and `vnet_cidr` to support
    routing options on more granular level of subnet and CIDR of the VNET.


### resource `alkira_connector_cisco_sdwan`

    * New argument `customer_asn` was added.


### resoruce `alkira_service_pan` (**NEEDS CONFIG UPDATE**)

    `panorama_ip_address` was changed to `panorama_ip_addresses` to
    support multiple Panorama IP addresses.



## Fixes

* Fix the missing `type` in `alkira_connector_cisco_sdwan`.
* Fix the inconsistent argument `disable_internet_exit` to `allow_nat_exit` in `alkira_connector_ipsec`.
* Fix the inconsistent argument `disable_advertise_on_prem_routes` to `advertise_on_prem_routes` in `alkira_connector_ipsec`.
* Fix the missing `availability` in `alkira_connector_ipsec`.
* Fix the `billing_tag_ids` in `alkira_connector_ipsec`.
