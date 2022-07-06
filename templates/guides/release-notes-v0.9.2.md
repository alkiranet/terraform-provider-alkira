---
subcategory: "Release Notes"
page_title: "v0.9.2"
description: |-
    Release notes for v0.9.2
---

This maintainence release brings one new resource `connector-infoblox`
, small resource updates and various bug fixes.

**NOTE:** The registration pin changes to `alkira-service-pan` is
mandatory for newer versions of PAN.


## New Resources

### resource `alkira_service_infoblox` (**NEW**)

    New service integration with Infoblox.


## New Data Sources

    New data sources for lists

    + alkira_list_as_path
    + alkira_list_community
    + alkira_list_extended_community
    + alkira_list_global_cidr


## Enhancements

    Several enhancements may need configuration updates. Please refer the
    resource documentation for more details.


### resource `alkira_service_pan`

    Palo Alto stopped supporting older `8.x` and `9.0.x` PAN. All customers need
    to update to newer version. To work with newer versions of PAN, the following
    changes were introduced:

    * New argument `registration_pin_id`, `registration_pin_value`
      and `registration_pin_expiry` were added as required by newer version
      of PAN.

    * New argument `masterkey_enabled`, `masterkey_username` and `masterkey_password`
      are introduced to support `masterkey` for PAN instances.

### resource `alkira_service_fortinet`

    * Add support for segment options.

### resource `internet_application`

    * Add support for BYOIP.


## Fixes

    * Retain the instance ID of `alkira_service_pan` to avoid creating instance again during `terraform update`.
    * Retain the instance iF of `alkira_service_fortinet` to avoid issue when doing `terraform update`.
    * Fix the routing options problem in `alkira_connector_azure_vnet`.
