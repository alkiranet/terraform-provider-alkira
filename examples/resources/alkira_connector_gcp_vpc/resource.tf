#
# Creating two segments and each one is for each connector
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
resource "alkira_credential_gcp_vpc" "terraform_gcp_account" {
  name                 = "customer-gcp"
  auth_provider        = "https://www.googleapis.com/oauth2/v1/certs"
  auth_uri             = "https://accounts.google.com/o/oauth2/auth"
  client_email         = "tenant@tenant.iam.gserviceaccount.com"
  client_id            = "tenant_client_id"
  client_x509_cert_url = "https://www.googleapis.com/robot/v1/metadata/x509/customer"
  private_key          = "tenant_private_key"
  private_key_id       = "tenant_private_key_id"
  project_id           = "test"
  token_uri            = "https://oauth2.googleapis.com/token"
  type                 = "service_account"
}

#
# Create GCP-VPC connector
#
resource "alkira_connector_gcp_vpc" "gcp_vpc1" {
  name           = "customer-vpc1"
  gcp_region     = "us-west1"
  gcp_vpc_id     = "0000000000000"
  gcp_vpc_name   = "customer-vpc1"
  credential_id  = alkira_credential_gcp_vpc.terraform_gcp_account.id
  cxp            = "US-WEST"
  group          = "test"
  segment        = alkira_segment.segment1.name
  size           = "SMALL"
}
