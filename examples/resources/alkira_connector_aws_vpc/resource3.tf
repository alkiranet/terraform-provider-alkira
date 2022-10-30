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
