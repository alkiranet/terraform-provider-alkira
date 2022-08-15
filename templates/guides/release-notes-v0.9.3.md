---
subcategory: "Release Notes"
page_title: "v0.9.3"
description: |-
    Release notes for v0.9.3
---

This release brings new resource to manage Routing Policy,
enhancements and various bug fixes.

**NOTE:** In resource `alkira_segment`, the old `cidr` argument was
replaced by `cidrs` now. Please update your configurations before
updating to this version.


## New Resources

### resource `alkira_policy_routing` (**NEW**)

    New resource to manage Routing Policy.


## New Data Sources

    New data sources for policy related resources

    + alkira_policy
    + alkira_policy_nat_rule
    + alkira_policy_prefix_list
    + alkira_policy_rule
    + alkira_policy_rule_list


## Enhancements

    Several enhancements may need configuration updates. Please refer the
    resource documentation for more details.

### resource `alkira_segment`

    * `cidr` was replaced by `cidrs` to allow specifying multiple CIDRs.

### resource `alkira_internet_application`

    * Add `internet_protocol` to support IPv6.

### resource `alkira_service_pan`

    * Add `license_sub_type` to support MODEL based and CREDIT based license.

### resource `alkira_service_fortinet`

    * Allow multi zone block to allow zone segment mapping.


## Fixes

    * Enable and update the `advanced_options` block of `alkira_connector_ipsec` to allow specifying advanced options for endpoints.
    * Retain the Vedge ID of `alkira_connector_cisco_sdwan` to avoid creating instance again during `terraform update`.
    * Retain the Endpoint ID of `alkira_connector_ipsec` to avoid issue when doing `terraform update`.
