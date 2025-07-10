---
subcategory: "Release Notes"
page_title: "v1.3.2"
description: |-
    Release notes for v1.3.2
---

This minor release contains enhancements and bug fixes across various resources
and data sources.

## Resources & Data Sources


#### data source `alkira_zta_profile`

* New data source added for retrieving ZTA (Zero Trust Architecture) profiles.

This data source allows you to reference existing ZTA profiles in your
Terraform configurations:

```
data "alkira_zta_profile" "example" {
    name = "example-zta-profile"
}
```


#### resource `alkira_connector_aws_vpc`

* Enhanced configuration options and improved validation.


#### resource `alkira_connector_ipsec`

* Updated configuration parameters and improved error handling.


#### resource `alkira_policy_rule`

* Additional configuration options and enhanced validation.


