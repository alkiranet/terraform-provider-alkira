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
