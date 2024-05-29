resource "alkira_peering_gateway_aws_tgw_attachment" "test" {
  name                          = "test"
  description                   = "test"
  peering_gateway_aws_tgw_id    = alkira_peering_gateway_aws_tgw.test.id
  requestor                     = "CXP"
  peer_aws_region               = "us-west-1"
  peer_aws_tgw_id               = "tgw-08ffaxxxx"
  peer_aws_account_id           = "0123456789"
}
