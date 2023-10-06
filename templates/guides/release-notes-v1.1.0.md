---
subcategory: "Release Notes"
page_title: "v1.1.0"
description: |-
    Release notes for v1.1.0
---

The main feature of this release is the authentication with API
Key. Now, to initialize the provider, the API key could be used
instead of the old `username` & `password`:

```
provider "alkira" {
  portal  = "test.portal.alkira.com"
  api_key = "xxxxxxxxxx"
}
```

The API key could be managed through *Settings* -> *User Management*
on Alkira Portal.


## Resources & Data Sources

#### resource `alkira_connector_fortinet_sdwan` (**NEW**)

New resource to manage connector to Fortinet SD-WAN.

#### resource `alkira_connector_aruba_edge`

* Optional new field `enabled`.

#### resource `alkira_connector_aws_vpc`

* Optional new field `tgw_connect_enabled`.

#### resource `alkira_connector_vmware_sdwan`

* Optional new field `enabled`.

#### resource `alkira_internet_application`

* Optional new field `bi_directional_az`.

#### resource `alkira_segment`

* Optional new field `service_traffic_distribution`.

#### resource `alkira_segment_resource_share`

* Optional new field `policy_rule_list_id`.

#### resource `alkira_policy_nat`

* Optional new field `allow_overlapping_translated_source_addresses`.


