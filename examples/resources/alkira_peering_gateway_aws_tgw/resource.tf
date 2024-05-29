resource "alkira_peering_gateway_aws_tgw" "test" {
  name         = "test"
  description  = "test"
  cxp          = "US-EAST"
  asn          = "64512"
  aws_region   = "us-east-1"
}
