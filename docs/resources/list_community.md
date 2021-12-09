---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "alkira_list_community Resource - terraform-provider-alkira"
subcategory: ""
description: |-
  This list could be used to matches a route when all values in the list are present on the route. A route matches a list when any of the values match.
---

# alkira_list_community (Resource)

This list could be used to matches a route when all values in the list are present on the route. A route matches a list when any of the values match.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) Name of the list.
- **values** (List of String) A list of communities to match on routes. Each community in the list is a tag value in the format of `AA:NN` format (where AA and NN are `0-65535`). AA denotes a AS number.

### Optional

- **description** (String) Description for the list.
- **id** (String) The ID of this resource.

