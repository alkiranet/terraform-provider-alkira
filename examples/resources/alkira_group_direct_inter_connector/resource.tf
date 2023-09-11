resource "alkira_group_direct_inter_connector" "test" {
  name                      = "test"
  description               = "test"
  cxp                       = "US-EAST"
  segment_id                = alkira_segment.test.id
  connector_type            = "AWS_VPC"
  connector_provider_region = "us-east-1"
}
