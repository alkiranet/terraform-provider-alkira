resource "alkira_service_zscaler" "test1" {
  name        = "acceptance-zia-1"
  description = "acceptance zia"
  cxp         = "US-WEST-1"

  connector_internet_exit_id = alkira_connector_internet_exit.test1.id
  primary_public_edge_ip     = "11.11.11.11"
  secondary_public_edge_ip   = "12.12.12.12"
  segment_ids                = [alkira_segment.test1.id]
  size                       = "MEDIUM"
  tunnel_protocol            = "IPSEC"

  ipsec_configuration {
    esp_dh_group_number      = "MODP1024"
    esp_encryption_algorithm = "AES256CBC"
    esp_integrity_algorithm  = "MD5"
    health_check_type        = "IKE_STATUS"
    http_probe_url           = "probe.url"
    ike_dh_group_number      = "MODP2048"
    ike_encryption_algorithm = "AES256CBC"
    ike_integrity_algorithm  = "SHA256"
    local_fpdn_id            = "local_fpdn_id"
    pre_shared_key           = "pre_shared_key"
  }
}
