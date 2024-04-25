---
subcategory: "Release Notes"
page_title: "v1.2.5"
description: |-
    Release notes for v1.2.5
---

This release contains enhancments and various fixes. Documentation has
also been updated with various fixes.


* The API client have been improved with the new authentication method
  to avoid the session limit problem.

* The error handling of API client failures have been improved to show
  more meaingful errors and messages.


### RESOURCES

#### resource `alkira_list_udr` (**NEW**)

New resource for defining `User Defined Routes`. The feature is still
in **BETA** phase and please contact DevOps team for enabling it.

#### resource `alkira_connector_azure_vnet`

* New optional field for working with `alkira_list_udr`.

#### resource `alkira_connector_remote_access`

* New optional field `search_scope_domain`.

#### resource `alkira_service_fortinet`

* New optional field `license_scheme`.

### DATA SOURCES

* Fix the problem of `alkira_ip_reservation`.


