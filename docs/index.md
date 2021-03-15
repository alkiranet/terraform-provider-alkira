---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "alkira Provider"
subcategory: ""
description: |-
    The Alkira provider is used to configure and manage your Alkira network infrastructure
---

# Alkira Provider

Lifecycle management of Alkira Cloud Services Exchange.

A typical provider configuration looks like this:

```hcl
provider "alkira" {
  portal   = "your_tenant_name.portal.alkira.com"
  username = "your_name@email.com"
  password = "your_password"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **portal** (String) The URL of Alkira Customer Portal.
- **username** (String) The tenant username.
- **password** (String) The tenant password.
