resource "alkira_connector_aws_tgw" "test1" {
  name                                  = "tftest-test1"
  cxp                                   = "US-EAST"
  segment_id                            = alkira_segment.test.id
  size                                  = "SMALL"
  peering_gateway_aws_tgw_attachment_id = alkira_peering_gateway_aws_tgw_attachment.test.id
}
