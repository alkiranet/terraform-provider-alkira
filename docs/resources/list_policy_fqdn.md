---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "alkira_list_policy_fqdn Resource - terraform-provider-alkira"
subcategory: ""
description: |-
  Policy FQDN list.
---

# alkira_list_policy_fqdn (Resource)

Policy FQDN list.

## Example Usage

```terraform
resource "alkira_list_policy_fqdn" "test" {
  name               = "test"
  description        = "test policy fqdn list"
  fqdns              = ["test.alkira.com"]
  list_dns_server_id = alkira_list_dns_server.test.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `fqdns` (Set of String) A list of FQDNs.
- `list_dns_server_id` (Number) ID of `list_dns_server` resource.
- `name` (String) Name of the list.

### Optional

- `description` (String) Description for the list.

### Read-Only

- `id` (String) The ID of this resource.
- `provision_state` (String) The provisioning state of the resource.

## Import

Import is supported using the following syntax:

```shell
terraform import alkira_list_policy_fqdn.example LIST_ID
```