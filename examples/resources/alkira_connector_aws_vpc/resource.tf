resource "alkira_connector_aws_vpc" "connector" {
  name           = "connector-test"
  vpc_id         = "your_vpc_id"

  aws_account_id = "your_aws_account_id"
  aws_region     = "us-east-1"

  credential_id  = alkira_credential_aws_vpc.account1.id
  cxp            = "US-WEST"
  group          = alkira_group.group1.name
  segment_id     = alkira_segment.segment1.id
  size           = "SMALL"
}
