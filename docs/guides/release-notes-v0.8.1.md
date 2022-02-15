---
subcategory: "Release Notes"
page_title: "v0.8.1"
description: |-
    Release notes for v0.8.1
---

~> **NOTE:** This release contains changes that breaks backward
compatbility of the old configurations
(`connector_azure_vnet`). Please read the release notes carefully and
test out the old configuration before upgrading.

This release brings several enhancements and minor bug fixes.


### resource `alkira_connector_azure_vnet`

Argument `azure_region` and `azure_subscription_id` are deprecated and
removed.

* New argument `service_tags`.


### resource `alkira_service_pan` with global protect support (**NEW**)

Initial support of global protect is added to
`resource_alkira_service_pan`.


### New optional `enabled` argument is added to all connector resources

New optional `enabled` argument is added to all connnector
resources. The default value will be `true`. This argument will allow
all connnectors to be created as disabled initially.


### resource `alkira_connector_aws_vpc`

New argument `direct_inter_vpc_communication` was added to allow
Inter-Vpc communications. There are some limitations when using
it. Please refer to resource documentation for more details.


### resource `alkira_connector_internet_exit`

New arguments were added:

* New argument `public_ip_number`
* New argument `traffic_distribution_algorithm`
* New argument `traffic_distribution_algorithm_attributes`

Please refer the resource documentation for more details.

