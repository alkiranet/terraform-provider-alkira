#
# Create a billing tag that could be used later on
#
resource "alkira_billing_tag" "tag1" {
  name           = "tag1"
}

resource "alkira_billing_tag" "tag2" {
  name           = "tag2"
}

#
# Now you could use those two billing tags by Id when creating
# supported resources.
#
# For example, you could create an AWS-VPC connector with above 2
# billing tags.
#
resource "alkira_connector_aws_vpc" "connector-aws-vpc1" {
  name           = "customer-aws-vpc1"
  billing_tags   = [alkira_billing_tag.tag1.id, alkira_billing_tag.tag2.id]

  # All other variables...
  # ...
}
