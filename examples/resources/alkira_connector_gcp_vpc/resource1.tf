resource "alkira_connector_gcp_vpc" "gcp_routing" {
  name           = "example-vpc1"
  gcp_region     = "us-west1"
  gcp_vpc_name   = "example-vpc1"
  cxp            = "US-WEST"
  size           = "SMALL"
  credential_id  = alkira_credential_gcp_vpc.terraform_gcp_account.id
  group          = alkira_group.group1.name
  segment_id     = alkira_segment.segment1.id

  gcp_routing {
    prefix_list_ids = [alkira_policy_prefix_list.azure.id]
    custom_prefix = "ADVERTISE_CUSTOM_PREFIX"
  }
}
