resource "alkira_segment" "segment" {
  name  = "seg-test"
  asn   = "65513"
  cidrs = ["10.16.1.0/24"]
}

resource "alkira_group" "group" {
  name        = "group-test"
  description = "test group"
}

resource "alkira_credential_aws_vpc" "account" {
  name           = "test-aws"
  aws_access_key = "your_aws_acccess_key"
  aws_secret_key = "your_secret_key"
  type           = "ACCESS_KEY"
}

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
