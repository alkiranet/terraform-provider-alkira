---
subcategory: "Release Notes"
page_title: "v1.2.11"
description: |-
    Release notes for v1.2.11
---

This release contains fixes and enhancements.

## Documentation

* Fix various inconsistent cases across all documentations.
* Update documentation with better indication of resources references from different resources.
* Update documentation with latest functionality changes.

## Resources & Data Sources

#### resource `alkira_byoip_prefix`

* Add optional field `cloud_provider`.

#### resource `alkira_connector_akamai_prolexic`

* Add optional field `description`.

#### resource `alkira_connector_aruba_edge`

* Make `aruba_edge_vrf_mapping` block as `required` now.

#### resource `alkira_connector_aws_vpc`

* Add optional field `description`.

#### resource `alkira_connector_azure_vnet`

* Add optional field `native_services` to specify on global or `vnet_cidr` or `vnet_subet` level.
* Add optional field `description`.
* Add optional field `udr_list_ids`.
* Remove validation of `size` field to allow using new supported connector sizes.

#### resource `alkira_connector_azure_expressroute`

* Fix the wrong payload naming of `gateway_mac_address` field.

#### resource `alkira_connector_cisco_sdwan`

* Add optional field `description`.
* Add optional field `tunnel_protocol`.

#### resource `alkira_connector_ipsec`

* Add optional field `customer_ip_type` to allow switching type of field `customer_gateway_ip`.

#### resource `alkira_service_checkpoint`

* Support updating credentials.

#### resource `alkira_service_fortinet`

* Support updating credentials.
* Add additional validation for `management_server_ip` avoid crash caused by empty value.

#### resource `alkira_service_infoblox`

* Support updating credentials.
* Add new computed field `service_group_implicit_group_id`.
* Add new computed field `service_group_id`.

#### resource `alkira_service_pan`

* Support updating credentials.

