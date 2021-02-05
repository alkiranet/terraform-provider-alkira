#
# Create a segment
#
resource "alkira_segment" "segment1" {
  name = "seg1"
  asn  = "65513"
  cidr = "10.16.1.0/24"
}


#
# Create the credential to store the access to the AWS account that
# VPCs belongs two. In this example, both VPCs belong to this AWS
# account.
#
resource "alkira_credential_aws_vpc" "account1" {
  name           = "customer-aws-1"
  aws_access_key = "your_aws_acccess_key"
  aws_secret_key = "your_secret_key"
  type           = "ACCESS_KEY"
}


#
# Create AWS-VPC connector for the first VPC and attach it with
# segment 1
#
resource "alkira_connector_aws_vpc" "connector-vpc1" {
  name           = "customer-vpc1"
  vpc_id         = "your_vpc_id"

  aws_account_id = "your_aws_account_id"
  aws_region     = "us-east-1"

  credential_id  = alkira_credential_aws_vpc.account1.id
  cxp            = "US-WEST"
  group          = "test"
  segment        = alkira_segment.segment1.name
  size           = "SMALL"
}
