---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "alkira_byoip_prefix Resource - terraform-provider-alkira"
subcategory: ""
description: |-
  Manage BYOIP Prefix.
---

# alkira_byoip_prefix (Resource)

Manage BYOIP Prefix.

## Example Usage

```terraform
resource "alkira_byoip" "test" {
  prefix      = "172.16.1.2"
  cxp         = "US-WEST"
  description = "simple test"
  message     = "1|aws|0123456789AB|198.51.100.0/24|20211231|SHA256|RSAPSS"
  signature   = "signature from AWS BYOIP"
  public_key  = "public key from AWS BYOIP"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **cxp** (String) CXP region.
- **message** (String) Message from AWS BYOIP.
- **prefix** (String) Prefix for BYOIP.
- **public_key** (String) Public Key from AWS BYOIP.
- **signature** (String) Signautre from AWS BYOIP.

### Optional

- **description** (String) Description for the list.
- **do_not_advertise** (Boolean) Do not advertise.
- **id** (String) The ID of this resource.

