---
subcategory: "Release Notes"
page_title: "v1.3.1"
description: |-
    Release notes for v1.3.1
---

This minor release contains mainly enhancements around Policy related
resources and minor fixes across all resources.

## Resources & Data Sources


#### resource `alkira_policy_prefix_list`

* New block `prefix`.

Add new optional block `prefix` to allow specifying `description` with
prefix like this:

```
prefix {
    cidr        = "0.0.0.0/0"
    description = "test prefix"
```

The old `prefixes` field will be deprecated in the future
release. Please update your current Terraform config accordingly.

* New optional field `description` in `prefix_range`.


#### resource `alkira_policy_routing`

* New optional field `set_med`.
* New optional field `source_routes_prefix_list_id`.
* New optional field `set_as_path_replace_with_segment_asn`.


#### resource `alkira_connector_aws_tgw`

* New read-only field `scale_group_id`.


#### resource `alkira_connector_azure_expressroute`

* Remove default value of field `ike_version`.
* Update default value of field `initiator`.


#### resource `alkira_service_f5`

* New value `BETTER` and `BEST` for field `deployment_type`.



