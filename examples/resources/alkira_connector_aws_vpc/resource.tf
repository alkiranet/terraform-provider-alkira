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
