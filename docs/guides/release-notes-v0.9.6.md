---
subcategory: "Release Notes"
page_title: "v0.9.6"
description: |-
    Release notes for v0.9.6
---

This release contains major improvements around resource
`alkira_service_pan` to make it easier to use. It requires update of
the current configurations.

**NOTE: Resource `credential_pan` and `credential_pan_instance` were
deprecated and it was combined into `alkira_service_pan` directly.**


## Changes

    * `credential_pan` is deprecated now and you could directly specify
      `pan_username` and `pan_password` in `alkira_service_pan` when
       Panorama is enabled.

    * `credential_pan_instance` is deprecated now and you could directly
      specify `auth_key` in `instance`.

    * New block `segment_options` is introduced across service resources
      to make configuring segment options more consistent.

    * New `implict_group_id` is added for all connector resources to be used
      with policy resources.

    * New `ha_mode` support for `alkira_connector_ipsec`.

    * New `gcp_project_id` support for `alkira_connector_gcp_vpc`.

    * New data sources.

## Fixes

    * More validations and fixes in resource `alkira_connector_ipsec`.
