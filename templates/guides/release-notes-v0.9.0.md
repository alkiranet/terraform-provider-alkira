---
subcategory: "Release Notes"
page_title: "v0.9.0"
description: |-
    Release notes for v0.9.0
---

This release brings new resources for supporting various new features,
enhancements and bug fixes. This is the first release that we try to
release provider with new feature support along with product.

## New Resources

### resource `alkira_byoip_prefix` (**NEW**)

Manage BYOIP prefix. This resource is needed by `alkira_connector_akamai_prolexic`.

### resource `alkira_connector_akamai_prolexic` (**BETA**)

New connector for connecting to Akamai Prolexic. This is still under
active development and may have changes in the future. This resource
will need `alkira_byoipd_prefix` and `alkira_segments` to work. Please
refer to its documentation for details.

### resource `alkira_group_connector` (**NEW**)

This resource replaces the old `alkira_group` resource. Now, there are
3 groups: **connector group**, **user group** and **segment resource
groups**.

### resource `alkira_group_user` (**NEW**)

Manage the new user group.

### resource `alkira_segment_resource` (**NEW**)

Manage segment resource. It's usually used along with
`alkira_segment_resource_share`.

### resource `alkira_segment_resource_share` (**NEW**)

This resource will allow to manage segment resource share from
Terraform. It's usually used along with `alkira_segment_resource`

### resource `alkira_service_fortinet` (**NEW**)

Manage the new Fortinet service.


## Enhancements

Please refer the resource documentation for more details.

### resource `alkira_connector_gcp_vpc`

Support GCP Routing Options now.

### resource `alkira_connector_internet_exit`

* New arguments were added to support `connector_akamai_prolexic`.
* New arguments for specifying targets (ILB/IP).
* Update all arguments to be more consistent with UI.

### resource `alkira_segment`

* Add new argument `reserve_public_ips`.
* Make `asn` argument optional with default value of `65514`.

### policy resources

Bug fixes and optimizations.
