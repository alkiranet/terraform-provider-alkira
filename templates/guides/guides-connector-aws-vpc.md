---
subcategory: "Guides"
page_title: "Connector - AWS VPC"
description: |-
    Guide for using AWS VPC Connector
---

This guide will cover building a complete environment using [alkira_connector_aws_vpc](https://registry.terraform.io/providers/alkiranet/alkira/latest/docs/resources/connector_aws_vpc).


## Getting started
Before getting started, a few resources need to exist.

### AWS resources
First, an [aws_vpc](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/vpc) and one or more of [aws_subnet](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/subnet) need to exist. We will use the official [AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs).

```terraform
resource "aws_vpc" "vpc" {
  cidr_block = "10.5.0.0/20"

  tags       = {
    Name     = "aws-vpc-east"
  }

}

resource "aws_subnet" "subnet" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = "10.5.1.0/24"

  tags = {
    Name = "app-subnet"
  }

}
```

### Alkira resources
In Alkira, we need to create types [alkira_credential_aws_vpc](https://registry.terraform.io/providers/alkiranet/alkira/latest/docs/resources/credential_aws_vpc) and [alkira_segment](https://registry.terraform.io/providers/alkiranet/alkira/latest/docs/resources/segment). If these resources already exist, you can reference them with the corresponding _data sources_.

```terraform
resource "alkira_credential_aws_vpc" "credential" {
  name           = "aws-credential"
  aws_access_key = var.aws_access_key
  aws_secret_key = var.aws_secret_key
  type           = "ACCESS_KEY"
}

resource "alkira_segment" "segment" {
  name = "corporate"
  cidr = "10.250.1.0/24"
}
```

## Connecting AWS VPC to Alkira
There are various options to choose from when connecting a _VPC_ to Alkira.

### Onboard complete CIDR with defaults
This example will connect the _aws_vpc_ we defined above and place it in the _alkira_segment_ we created using the _alkira_credential_aws_vpc_ for authentication. The following configuration will onboard the entire _VPC CIDR_ with defaults.


```terraform
resource "alkira_connector_aws_vpc" "connector" {

  # ID and CIDR from aws_vpc resource
  vpc_id          = aws_vpc.vpc.id
  vpc_cidr        = aws_vpc.vpc.cidr_block

  # AWS account_id and region
  aws_account_id  = "12345678"
  aws_region      = "us-east-2"

  # ID from credential and segment resource
  segment_id      = alkira_segment.segment.id
  credential_id   = alkira_credential_aws_vpc.credential.id

  # Connector configuration
  name            = "aws-connector-east"
  cxp             = "US-EAST-2"
  size            = "SMALL"

}
```

### Create group and associate it with AWS connector
_Micro-Segmentation_ can be accomplished in Alkira by using the [alkira_group_connector](https://registry.terraform.io/providers/alkiranet/alkira/latest/docs/resources/group_connector) resource. The following example creates a new group called _dev_ and applies it to the _connector_.


```terraform
resource "alkira_group_connector" "group" {
  name        = "dev"
  description = "Created by Terraform"
}

resource "alkira_connector_aws_vpc" "connector" {

  # ID and CIDR from aws_vpc resource
  vpc_id          = aws_vpc.vpc.id
  vpc_cidr        = aws_vpc.vpc.cidr_block

  # AWS account_id and region
  aws_account_id  = "12345678"
  aws_region      = "us-east-2"

  # ID from credential and segment resource
  segment_id      = alkira_segment.segment.id
  credential_id   = alkira_credential_aws_vpc.credential.id

  # Connector configuration
  name            = "aws-connector-east"
  cxp             = "US-EAST-2"
  size            = "SMALL"

  # Group for micro-segmentation
  group           = data.alkira_group.group.name

}
```

### Onboard specific subnets only
This example will onboard only specific subnets to Alkira. The **vpc_cidr** option gets replaced with the **vpc_subnet** option.

```terraform
resource "alkira_connector_aws_vpc" "connector" {

  # ID from aws_vpc resource
  vpc_id          = aws_vpc.vpc.id

  # ID and CIDR from aws_subnet resource to be onboarded
  vpc_subnet {
    id   = aws_subnet.subnet.id
    cidr = aws_subnet.subnet.cidr_block
  }

  # AWS account_id and region
  aws_account_id  = "12345678"
  aws_region      = "us-east-2"

  # ID from credential and segment resource
  segment_id      = alkira_segment.segment.id
  credential_id   = alkira_credential_aws_vpc.credential.id

  # Connector configuration
  name            = "aws-connector-east"
  cxp             = "US-EAST-2"
  size            = "SMALL"

  # Group for micro-segmentation
  group           = data.alkira_group.group.name

}
```

### Advertise Custom Prefix
By default, Alkira will override the existing default route and route the traffic to the CXP. As an alternative, you can provide a list of prefixes for which traffic must be routed. To do this, you must create an [alkira_policy_prefix_list](https://registry.terraform.io/providers/alkiranet/alkira/latest/docs/resources/policy_prefix_list) resource and then add the **vpc_route_table** block to the connector configuration.

```terraform
resource "alkira_policy_prefix_list" "prefix" {
  name        = "custom-prefix"
  description = "Created by Terraform"
  prefixes    = ["10.50.1.0/24"]
}

resource "alkira_connector_aws_vpc" "connector" {

  # ID and CIDR from aws_vpc resource
  vpc_id           = aws_vpc.vpc.id
  vpc_cidr         = aws_vpc.vpc.cidr_block

  # AWS account_id and region
  aws_account_id   = "12345678"
  aws_region       = "us-east-2"

  # ID from credential and segment resource
  segment_id       = alkira_segment.segment.id
  credential_id    = alkira_credential_aws_vpc.credential.id

  # Connector configuration
  name             = "aws-connector-east"
  cxp              = "US-EAST-2"
  size             = "SMALL"

  # Group for micro-segmentation
  group            = data.alkira_group.group.name

  # Route to custom prefixes
  vpc_route_table = {
    id             = aws_vpc.vpc.default_route_table_id
    options        = "ADVERTISE_CUSTOM_PREFIX"
    prefix_list_id = alkira_policy_prefix_list.prefix.id
  }

}
```

### Direct _Inter-VPC_ Communication
You can also set **direct_inter_vpc_communication** to _true_ for _VPCs_ to communication directly. This is _false_ by default.

```terraform
resource "alkira_connector_aws_vpc" "connector" {

  # ID and CIDR from aws_vpc resource
  vpc_id           = aws_vpc.vpc.id
  vpc_cidr         = aws_vpc.vpc.cidr_block

  # AWS account_id and region
  aws_account_id   = "12345678"
  aws_region       = "us-east-2"

  # ID from credential and segment resource
  segment_id       = alkira_segment.segment.id
  credential_id    = alkira_credential_aws_vpc.credential.id

  # Connector configuration
  name             = "aws-connector-east"
  cxp              = "US-EAST-2"
  size             = "SMALL"

  # Group for micro-segmentation
  group            = data.alkira_group.group.name

  # Direct inter-vpc communication
  direct_inter_vpc_communication = true

}
```

### Billing Tags
Billing tags can be created and applied to resources for _cost optimization_ purposes. Create the resource [alkira_billing_tag](https://registry.terraform.io/providers/alkiranet/alkira/latest/docs/resources/billing_tag) and then add **billing_tag_ids** to the connector configuration to apply it.

```terraform
resource "alkira_billing_tag" "tag" {
  name           = "digital"
  description    = "Create by Terraform"
}

resource "alkira_connector_aws_vpc" "connector" {

  # ID and CIDR from aws_vpc resource
  vpc_id           = aws_vpc.vpc.id
  vpc_cidr         = aws_vpc.vpc.cidr_block

  # AWS account_id and region
  aws_account_id   = "12345678"
  aws_region       = "us-east-2"

  # ID from credential and segment resource
  segment_id       = alkira_segment.segment.id
  credential_id    = alkira_credential_aws_vpc.credential.id

  # Connector configuration
  name             = "aws-connector-east"
  cxp              = "US-EAST-2"
  size             = "SMALL"

  # Group for micro-segmentation
  group            = data.alkira_group.group.name

  # Apply billing tags
  billing_tag_ids  = [alkira_billing_tag.tag.id]

}
```

## Simplifying configuration with modules
You can also leverage Alkira's _pre-built_ [aws_vpc module](https://github.com/alkiranet/terraform-alkira-aws-vpc) to simplify configuration and ease of deployment. This module will build the AWS VPC/subnets, add them to an existing Alkira segment/group, apply existing billing tags, and offer all the advanced options for customized routing. All options can be provided with _human readable_ values so you don't have to dig up _resource IDs_. Basic usage would work as follows:

```terraform
module "aws_vpc" {
  source = "alkiranet/aws-vpc/alkira"

  name    = "aws-vpc-east"
  cidr    = "10.5.0.0/20"

  subnets = [
    {
      name = "app-subnet-a"
      cidr = "10.5.1.0/24"
      zone = "us-east-2a"
    },
    {
      name = "app-subnet-b"
      cidr = "10.5.2.0/24"
      zone = "us-east-2b"
    }
  ]

  cxp          = "US-EAST-2"
  segment      = "corporate"
  group        = "nonprod"
  billing_tags = ["cloud", "network"]
  credential   = "aws-auth"

}
```