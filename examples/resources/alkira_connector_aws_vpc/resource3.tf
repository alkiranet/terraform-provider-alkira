provider "aws" {
  region     = local.aws_region
  access_key = local.aws_access_key
  secret_key = local.aws_secret_key
}

resource "aws_vpc" "vpc_test" {
  cidr_block = "10.2.0.0/16"

  tags = {
    Name = "vpc-test"
  }
}

resource "aws_subnet" "vpc_subnet1" {
  vpc_id     = aws_vpc.vpc_test.id
  cidr_block = "10.2.0.0/24"
}

resource "aws_subnet" "vpc_subnet2" {
  vpc_id     = aws_vpc.vpc_test.id
  cidr_block = "10.3.16.0/20"
}

resource "alkira_connector_aws_vpc" "test" {
  name           = "test"

  aws_account_id = local.aws_account_id
  aws_region     = local.aws_region
  cxp            = local.cxp

  vpc_id         = aws_vpc.vpc_test.id

  credential_id  = alkira_credential_aws_vpc.test.id
  group          = alkira_group.test.name
  segment_id     = alkira_segment.test.id
  size           = "SMALL"

  vpc_route_table {
    id              = aws_vpc.vpc_test.default_route_table_id
    options         = "ADVERTISE_DEFAULT_ROUTE"
  }

  vpc_subnet {
    id   = aws_subnet.vpc_subnet1.id
    cidr = aws_subnet.vpc_subnet1.cidr_block
  }

  vpc_subnet {
    id   = aws_subnet.vpc_subnet2.id
    cidr = aws_subnet.vpc_subnet2.cidr_block
  }
}
