resource "alkira_byoip" "test" {
  prefix      = "172.16.1.2"
  cxp         = "US-WEST"
  description = "simple test"
  message     = "1|aws|0123456789AB|198.51.100.0/24|20211231|SHA256|RSAPSS"
  signature   = "signature from AWS BYOIP"
  public_key  = "public key from AWS BYOIP"
}
