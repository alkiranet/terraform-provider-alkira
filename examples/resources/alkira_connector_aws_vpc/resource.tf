#
# EXAMPLE 1
#
# This is simple example to show how to create an AWS-VPC connector.
#
# One segment and credential are needed for a connector and you could
# also adjust routing preferences by specifying `vpc_cidr` or
# `vpc_subnet` or `vpc_route_tables`.
#
resource "alkira_segment" "segment1" {
  name  = "seg1"
  asn   = "65513"
  cidrs = ["10.16.1.0/24"]
}

resource "alkira_group" "group1" {
  name        = "group1"
  description = "test group"
}

resource "alkira_credential_aws_vpc" "account1" {
  name           = "customer-aws-1"
  aws_access_key = "your_aws_acccess_key"
  aws_secret_key = "your_secret_key"
  type           = "ACCESS_KEY"
}

#
# EXAMPLE 2
#
# Create one connector for a VPC and attach it with segment1
#
resource "alkira_connector_aws_vpc" "connector1" {
  name           = "vpc1"
  vpc_id         = "your_vpc_id"

  aws_account_id = "your_aws_account_id"
  aws_region     = "us-east-1"

  credential_id  = alkira_credential_aws_vpc.account1.id
  cxp            = "US-WEST"
  group          = alkira_group.group1.name
  segment_id     = alkira_segment.segment1.id
  size           = "SMALL"
}

#
# EXAMPLE 3
#
# Create a VPC and create a aws-vpc connector to connect to it.
#
resource "aws_vpc" "vpc2" {
  cidr_block = "10.2.0.0/16"

  tags = {
    Name = "vpc2"
  }
}

resource "aws_subnet" "vpc2_subnet1" {
  vpc_id     = aws_vpc.vpc2.id
  cidr_block = "10.2.0.0/24"
}

#
# EXAMPLE 4
#
# Create a connector and adjust the routing to use the default
# route. There could be multiple vpc_route_table sections for
# additional route tables.
#
resource "alkira_connector_aws_vpc" "connector2" {
  name           = "vpc2"

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
