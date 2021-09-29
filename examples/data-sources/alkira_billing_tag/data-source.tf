#
# Refer an existing bill tag by name
#
data "alkira_billing_tag" "tag1" {
  name = "tag1"
}


#
# Use it inside any other resources like this connector_aws_vpc.
#
resource "alkira_connector_aws_vpc" "test_connector" {
  name            = "test_vpc"
  vpc_id          = "test_vpc_id"

  aws_account_id  = "test_vpc_aws_account_id"
  aws_region      = "us-east-1"

  credential_id   = alkira_credential_aws_vpc.account1.id
  cxp             = "US-WEST"
  group           = "test"
  segment_id      = alkira_segment.test_segment.id
  size            = "SMALL"

  billing_tag_ids = [data.alkira_billing_tag.tag1]
}
