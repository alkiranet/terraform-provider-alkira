---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "alkira_list_global_cidr Data Source - terraform-provider-alkira"
subcategory: ""
description: |-
  Use this data source to get an existing Global CIDR List.
---

# alkira_list_global_cidr (Data Source)

Use this data source to get an existing Global CIDR List.

## Example Usage

```terraform
data "alkira_list_global_cidr" "test" {
  name = "test"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the Global CIDR List.

### Optional

- `values` (String) The values of the list.

### Read-Only

- `id` (String) The ID of this resource.
