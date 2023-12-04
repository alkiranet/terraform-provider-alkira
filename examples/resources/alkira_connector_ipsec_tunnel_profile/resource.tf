resource "alkira_connector_ipsec_tunnel_profile" "tunnel1" {
  name        = "tunnel1"
  description = "tunnel1"

  ipsec_encryption_algorithm = "AES256CBC"
  ipsec_integrity_algorithm  = "SHA1"
  ipsec_dh_group             = "MODP1024"

  ike_encryption_algorithm = "AES256CBC"
  ike_integrity_algorithm  = "SHA1"
  ike_dh_group             = "MODP1024"
}

