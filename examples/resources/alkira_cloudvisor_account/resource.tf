resource "alkira_cloudvisor_account" "test" {
  name           = "test"
  credential_id  = alkira_credential_aws_vpc.test.id # using a credential_aws_vpc
  cloud_provider = "AWS"
  auto_sync      = "NONE"
}
