---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "alkira_cxp_peering_gateway Data Source - terraform-provider-alkira"
subcategory: ""
description: |-
  This data source allows to retrieve an existing Cxp Peering Gateway by its name.
---

# alkira_cxp_peering_gateway (Data Source)

This data source allows to retrieve an existing Cxp Peering Gateway by its name.

## Example Usage

```terraform
data "alkira_cxp_peering_gateway" "test-gateway" {
  name = "test-gateway"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the resource.

### Read-Only

- `id` (String) The ID of this resource.