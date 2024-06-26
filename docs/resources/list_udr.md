---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "alkira_list_udr Resource - terraform-provider-alkira"
subcategory: ""
description: |-
  User Defined Routes (UDR) list.
---

# alkira_list_udr (Resource)

User Defined Routes (UDR) list.

## Example Usage

```terraform
resource "alkira_list_udr" "test1" {
  name               = "tf-test-1"
  description        = "terraform test UDR list 1"
  cloud_provider     = "AZURE"

  route {
    prefix = "10.0.0.0/24"
    description = "test route 1"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the list.
- `route` (Block Set, Min: 1) ID of `list_dns_server` resource. (see [below for nested schema](#nestedblock--route))

### Optional

- `cloud_provider` (String) Cloud provider. Only `AZURE` is supported for now.
- `description` (String) Description for the list.

### Read-Only

- `id` (String) The ID of this resource.
- `provision_state` (String) The provisioning state of the resource.

<a id="nestedblock--route"></a>
### Nested Schema for `route`

Required:

- `prefix` (String) The prefix of the route. This prefix must be in the CIDR format (`x.x.x.x/mask`). The mask can be between `8-32`.

Optional:

- `description` (String) Description for the route.

## Import

Import is supported using the following syntax:

```shell
terraform import alkira_list_udr.example LIST_ID
```
