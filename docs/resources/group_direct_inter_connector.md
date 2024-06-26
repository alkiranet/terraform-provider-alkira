---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "alkira_group_direct_inter_connector Resource - terraform-provider-alkira"
subcategory: ""
description: |-
  Provide direct inter-connector group resource.
---

# alkira_group_direct_inter_connector (Resource)

Provide direct inter-connector group resource.

## Example Usage

```terraform
resource "alkira_group_direct_inter_connector" "test" {
  name                      = "test"
  description               = "test"
  cxp                       = "US-EAST"
  segment_id                = alkira_segment.test.id
  connector_type            = "AWS_VPC"
  connector_provider_region = "us-east-1"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `connector_type` (String) The type of the connector.
- `name` (String) The name of the group.
- `segment_id` (String) The segment ID of the group.

### Optional

- `azure_network_manager_id` (Number) The Azure Virtual Network Manager's Alkira ID.
- `connector_provider_region` (String) The region of the connector.
- `cxp` (String) The CXP of the group.
- `description` (String) The description of the group.

### Read-Only

- `id` (String) The ID of this resource.
- `provision_state` (String) The provisioning state of the resource.

## Import

Import is supported using the following syntax:

```shell
terraform import alkira_group_direct_inter_connector.example GROUP_ID
```
