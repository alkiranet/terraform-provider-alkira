---
subcategory: "Release Notes"
page_title: "v1.2.4"
description: |-
    Release notes for v1.2.4
---

This release contains enhancments and various fixes. Documentation has
also been updated with various fixes.


### RESOURCES

#### resource `alkira_connector_aws_vpc`

* New optional field `overlay_subnets`.

#### resource `alkira_connector_remote_access`

* The missing optional field `prefix_list_id` has been added.

#### resource `alkira_policy_routing`

Two new optional fields are added:

* `enable_as_override`
* `match_segment_resource_ids`


### DATA SOURCES

The following data sources have been added to make it easier to create
policies with with existing resources:

* alkira_connector_remote_access
* alkira_connector_vmware_sdwan
* alkira_internet_application


