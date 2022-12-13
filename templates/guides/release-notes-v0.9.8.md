---
subcategory: "Release Notes"
page_title: "v0.9.8"
description: |-
    Release notes for v0.9.8
---

This release introduces new resource to integrate with Cisco Firepower
Firewall Service with various enhancements and bug fixes.

## Resources & Data Sources


#### resource `alkira_service_cisco_ftdv` (**NEW**)

(**BETA**) New service for integration with Cisco Firepower Fireall Service.


#### resource `alkira_connector_cisco_sdwan`

* Mark `customer_asn` to be required.
* Mark `endpoint` to be required.


#### resource `alkira_connector_ipsec`

* Fix `availability` field when `routing_type` is set to `BOTH`.


#### resource `alkira_internet_application`

* Enable `internet_protocol`.


#### resource `alkira_list_global_cidr`

* Add optional `tags` field.


#### resource `alkira_policy_prefix_list`

* Add optional `prefix_range`.


#### resource `alkira_policy_routing`

* Update documentation with more examples for both `INBOUND` and
  `OUTBOUND` policy.

* Fix the broken `enabled` field.


#### resource `alkira_segment`

* Add new optional `description`.


#### resource `alkira_service_checkpoint`

* Add `id` in instance block for tracking instance state.
* Add `name` in instance block for tracking instance.
* Update documentation.


#### resource `alkira_service_pan`

* Add new optional `pan_license_key`.
