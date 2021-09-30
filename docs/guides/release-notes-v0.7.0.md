---
subcategory: "Release Notes"
page_title: "v0.7.0"
description: |-
    Release notes for v0.7.0
---

Major release with new resource to support Cisco SDWAN connectors and
enhancement across policy related resources.

~> **NOTE:** This release contains changes that breaks backward
compatbility of the old configurations. Please read this release notes
carefully and test out the old configuration before upgrading.

### resource `alkira_connector_cisco_sdwan` (**NEW**)

New resource is added to support managing Cisco SDWAN connector
through Terraform.


### resource `credential_cisco_sdwan` (**NEW**)

New credential resource is added to support managing Cisco SDWAN
credential.


### resource `alkira_connector_aws_vpc`

* `billing_tags` argument is renamed to `billing_tags_id` for consistency.
* `segment` argument is renamed to `segment_id` for consistency.
* Bug fixes.


### resource `alkira_connector_azure_vnet`

* `billing_tags` argument is renamed to `billing_tags_id` for consistency.
* `segment` argument is renamed to `segment_id` for consistency.
* Bug fixes.


### resource `alkira_connector_gcp_vpc`

* `billing_tags` argument is renamed to `billing_tags_id` for consistency.
* `segment` argument is renamed to `segment_id` for consistency.
* Bug fixes.


### resource `alkira_connector_internet_exit` (old `alkira_connector_internet`)

The resource was renamed to `alkira_connector_internet_exit` to be
consistent with the official documentation.

* `segment` argument is renamed to `segment_id` for consistency.
* Bug fixes and documentation update.


### resource `alkira_connector_ipsec`

* Argument `availability` was added for static routing. (**NEW**)
* Argument `big_auth_key` was added. (**NEW**)
* `billing_tags` argument is renamed to `billing_tags_id` for consistency.
* `segment` argument is renamed to `segment_id` for consistency.
* Bug fixes.

### resource `alkira_internet_application`

* `billing_tags` argument is renamed to `billing_tags_id` for consistency.
* `segment` argument is renamed to `segment_id` for consistency.
* Bug fixes and documentation update.

### resource `alkira_policy`

* Bug fixes.
* Dcoumentation update.

### resource `alkira_policy_prefix_list`

* Bug fixes.
* Documentation update.

### resource `alkira_policy_rule`

New argument was added to support using `policy_prefix_list` in rules.

* `src_prefix_list_id` was added.
* `dst_prefix_list_id` was added.
* `application_list` argument was renamed to `application_ids`.
* `application_family_list` argument was renamed to `application_family_ids`.
* Bug fixes.
* Documentation update.

### resource `alkira_policy_rule_list`

* Bug fixes.
* Documentation update.

### resource `alkira_service_pan`

* `group` argument was removed.
* `max_instance_count` was required argument now.
* Bug fixes.
* Documentation update.
