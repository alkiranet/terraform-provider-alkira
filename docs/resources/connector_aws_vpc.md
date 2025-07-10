---
page_title: "alkira_connector_aws_vpc Resource - terraform-provider-alkira"
subcategory: ""
description: |-
  Provide AWS VPC Connector resource.
---

# alkira_connector_aws_vpc (Resource)

Provide AWS VPC Connector resource.

This resource is usually used along with `terraform-provider-aws`.

## Routing Options

Either `vpc_cidr` or `vpc_subnet` needs to be specified for routing
purpose.  If `vpc_cidr` is provided, it will automatically select all
associated subnets of the given VPC. Otherwise, you can select
certain subnets by specifying `vpc_subnet`.

`vpc_route_tables` can be used to adjust the routing options against
the specified route tables. When `OVERRIDE_DEFAULT_ROUTE` is
specified, the existing default route will be overwritten and the
traffic will be routed to Alkira CXP.  When `ADVERTISE_CUSTOM_PREFIX`
is specified, you need to provide a list of prefixes for which traffic
must be routed to Alkira CXP.

When `vpc_cidr` is used, `vpc_route_tables` should be also specified
to ensure that the traffic is attracted to the CXP.


## Tips

* Changing an existing connector to a new AWS VPC is not supported at
  this point. You need to create a new connector for a new AWS VPC.

* Updating an existing connector requires the tenant network to be
  re-provisioned to make the change effective, e.g. changing the
  segment the connector is associated.

* When direct inter VPC communication is enabled, several other
  functionalities won't work, like NAT policy, segment resource share,
  internet-facing applications and traffic policies.


## Example Usage

This is one simple minimal example to create an AWS VPC connector. One
`alkira_segment` and `alkira_credential_aws_vpc` are always required.

```terraform
# Basic AWS VPC Connector Example
resource "alkira_connector_aws_vpc" "basic" {
  name           = "aws-vpc-basic"
  description    = "Basic AWS VPC connector example"
  vpc_id         = "vpc-12345678"
  aws_account_id = "123456789012"
  aws_region     = "us-east-1"
  credential_id  = alkira_credential_aws_vpc.account1.id
  cxp            = "US-WEST"
  group          = alkira_group.group1.name
  segment_id     = alkira_segment.segment1.id
  size           = "SMALL"
  enabled        = true
}

# Advanced AWS VPC Connector with High Availability
resource "alkira_connector_aws_vpc" "advanced" {
  name            = "aws-vpc-advanced"
  description     = "Advanced AWS VPC connector with high availability"
  vpc_id          = "vpc-87654321"
  aws_account_id  = "123456789012"
  aws_region      = "us-east-1"
  credential_id   = alkira_credential_aws_vpc.account1.id
  cxp             = "US-WEST"
  failover_cxps   = ["US-EAST"]
  group           = alkira_group.group1.name
  segment_id      = alkira_segment.segment1.id
  size            = "MEDIUM"
  enabled         = true
  billing_tag_ids = [alkira_billing_tag.tag1.id]

  # VPC CIDR blocks for routing
  vpc_cidr = [
    "10.0.0.0/16",
    "10.1.0.0/16"
  ]

  # TGW Connect for high performance
  tgw_connect_enabled = true

  # TGW attachment configuration
  tgw_attachment {
    subnet_id = "subnet-12345678"
    az        = "us-east-1a"
  }

  tgw_attachment {
    subnet_id = "subnet-87654321"
    az        = "us-east-1b"
  }

  # Enable direct inter-VPC communication
  direct_inter_vpc_communication_enabled = true
  direct_inter_vpc_communication_group   = "production-vpcs"
}

# VPC Connector with Subnet-based Routing
resource "alkira_connector_aws_vpc" "subnet_routing" {
  name            = "aws-vpc-subnet-routing"
  description     = "AWS VPC connector with subnet-based routing"
  vpc_id          = "vpc-abcdef12"
  aws_account_id  = "123456789012"
  aws_region      = "us-west-2"
  credential_id   = alkira_credential_aws_vpc.account1.id
  cxp             = "US-WEST"
  group           = alkira_group.group1.name
  segment_id      = alkira_segment.segment1.id
  size            = "SMALL"
  enabled         = true
  billing_tag_ids = [alkira_billing_tag.tag1.id]

  # Use specific subnets instead of VPC CIDR
  vpc_subnet {
    id   = "subnet-production1"
    cidr = "10.0.1.0/24"
  }

  vpc_subnet {
    id   = "subnet-production2"
    cidr = "10.0.2.0/24"
  }

  # Route table configuration
  vpc_route_table {
    id      = "rtb-12345678"
    options = "ADVERTISE_DEFAULT_ROUTE"
  }

  vpc_route_table {
    id      = "rtb-87654321"
    options = "ADVERTISE_CUSTOM_PREFIX"
  }
}

# Required supporting resources for the examples
resource "alkira_billing_tag" "tag1" {
  name        = "aws-vpc-connector-tag"
  description = "Billing tag for AWS VPC connectors"
}

resource "alkira_credential_aws_vpc" "account1" {
  name           = "aws-account-credentials"
  aws_account_id = "123456789012"
  aws_access_key = "your_access_key"
  aws_secret_key = "your_secret_key"
}

resource "alkira_group" "group1" {
  name        = "aws-vpc-group"
  description = "Group for AWS VPC connectors"
}

resource "alkira_segment" "segment1" {
  name  = "production-segment"
  asn   = "65001"
  cidrs = ["10.0.0.0/8"]
}
```

To create a connector and adjust the routing to use the default
route. Multiple `vpc_route_table` blocks can be used for additional
route tables:

```terraform
resource "alkira_connector_aws_vpc" "connector" {
  name           = "vpc"

  aws_account_id = local.aws_account_id
  aws_region     = local.aws_region
  cxp            = local.cxp

  vpc_id         = aws_vpc.vpc2.id
  vpc_cidr       = [aws_vpc.vpc2.cidr_block]

  credential_id  = alkira_credential_aws_vpc.account1.id
  group          = alkira_group.group1.name
  segment_id     = alkira_segment.segment1.id
  size           = "SMALL"

  vpc_route_table {
    id              = aws_vpc.vpc2.default_route_table_id
    options         = "ADVERTISE_DEFAULT_ROUTE"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `aws_account_id` (String) AWS Account ID.
- `aws_region` (String) AWS Region where VPC resides.
- `credential_id` (String) ID of resource `credential_aws_vpc`.
- `cxp` (String) The CXP where the connector should be provisioned.
- `name` (String) The name of the connector.
- `segment_id` (String) The ID of segments associated with the connector. Currently, only `1` segment is allowed.
- `size` (String) The size of the connector, one of `5XSMALL`,`XSMALL`,`SMALL`, `MEDIUM`, `LARGE`, `2LARGE`, `5LARGE`, `10LARGE`, `20LARGE`.
- `vpc_id` (String) The ID of the target VPC.

### Optional

- `billing_tag_ids` (Set of Number) Billing tags to be associated with the resource. (see resource `alkira_billing_tag`).
- `description` (String) The description of the connector.
- `direct_inter_vpc_communication_enabled` (Boolean) Enable direct inter-vpc communication. Default is set to `false`.
- `direct_inter_vpc_communication_group` (String) Direct inter-vpc communication group.
- `enabled` (Boolean) Whether the connector is enabled. Default is `true`.
- `failover_cxps` (Set of String) A list of additional CXPs where the connector should be provisioned for failover.
- `group` (String) The group of the connector.
- `overlay_subnets` (List of String) Overlay subnet.
- `scale_group_id` (String) The ID of the scale group associated with the connector.
- `tgw_attachment` (Block List) TGW attachment. (see [below for nested schema](#nestedblock--tgw_attachment))
- `tgw_connect_enabled` (Boolean) When it's set to `true`, Alkira will use TGW Connect attachments to build connection to AWS Transit Gateway. Connect Attachments suppport GRE tunnel protocol for high performance and BGP for dynamic routing. This applies to all TGW attachments. This field can be set to `true` only if the VPC is in the same AWS region as the Alkira CXP it is being deployed onto.
- `vpc_cidr` (List of String) The list of CIDR attached to the target VPC for routing purpose. It could be only specified if `vpc_subnet` is not specified.
- `vpc_route_table` (Block Set) VPC route table (see [below for nested schema](#nestedblock--vpc_route_table))
- `vpc_subnet` (Block Set) The list of subnets of the target VPC for routing purpose. It could only specified if `vpc_cidr` is not specified. (see [below for nested schema](#nestedblock--vpc_subnet))

### Read-Only

- `id` (String) The ID of this resource.
- `implicit_group_id` (Number) The ID of implicit group automaticaly created with the connector.
- `provision_state` (String) The provisioning state of connector.

<a id="nestedblock--tgw_attachment"></a>
### Nested Schema for `tgw_attachment`

Required:

- `az` (String) The availability zone of the subnet.
- `subnet_id` (String) The Id of the subnet.


<a id="nestedblock--vpc_route_table"></a>
### Nested Schema for `vpc_route_table`

Optional:

- `id` (String) The Id of the route table
- `options` (String) Routing options, one of `ADVERTISE_DEFAULT_ROUTE`, `OVERRIDE_DEFAULT_ROUTE` or `ADVERTISE_CUSTOM_PREFIX`.
- `prefix_list_ids` (Set of Number) Prefix List IDs


<a id="nestedblock--vpc_subnet"></a>
### Nested Schema for `vpc_subnet`

Optional:

- `cidr` (String) The CIDR of the subnet.
- `id` (String) The Id of the subnet.

## Import

Import is supported using the following syntax:

```shell
terraform import alkira_connector_aws_vpc.example CONNECTOR_ID
```
