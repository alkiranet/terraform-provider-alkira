# Normal
resource "alkira_credential_aws_vpc" "account1" {
  name           = "customer-aws-1"
  aws_access_key = "your_aws_acccess_key"
  aws_secret_key = "your_secret_key"
  type           = "ACCESS_KEY"
}

# Using environment variables
export AWS_ACCESS_KEY_ID=XXX
export AWS_SECRET_ACCESS_KEY=XXX

resource "alkira_credential_aws_vpc" "account1" {
  name           = "customer-aws-1"
  type           = "ACCESS_KEY"
}

