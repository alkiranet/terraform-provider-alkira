---
subcategory: "Release Notes"
page_title: "v1.0.0"
description: |-
    Release notes for v1.0.0
---

This is the first official release with support of provisioning (BETA)
through Terraform Provider.

When initializing provider with `provision` flag, all eligible
resources in the configuration will be automatically provisioned, e.g.

```
provider "alkira" {
  portal    = "xxx"
  username  = "xxx"
  password  = "xxx"
  provision = true
}
```

~> **NOTE:** this feature is still in BETA. It may need to be manually
enabled for the tenant. Please consult Alkira support for more
information.


## Resources & Data Sources

#### resource `alkira_connector_gcp_vpc`

Update the resource to work better with Google Provider.

* Change the value of `id` in the subnet block to be consistent with
  Google Provider.

* `gcp_vpc_id` is not needed anymore.

* Update examples with latest changes.

#### resource `alkira_connector_cisco_sdwan`

* Tweak the default value of `allow_nat_exit`.

#### resource `alkira_segment_resource`

* Add optional `description`.

#### resource `alkira_service_checkpoint`

* Rename `user_name` to `username` in `management_server` block.
* Rename `management_server_password` to simple `password` in `management_server` block.
