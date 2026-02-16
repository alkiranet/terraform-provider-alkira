resource "alkira_service_bluecat" "minimal" {
  name                = "bluecat-minimal"
  cxp                 = "US-WEST"
  global_cidr_list_id = alkira_list_global_cidr.basic.id
  segment_ids         = [alkira_segment.default.id]
  service_group_name  = "dns-basic"

  instance {
    type = "BDDS"
    
    bdds_options {
      hostname       = "bdds"
      model          = "cBDDS50"
      version        = "9.4.0"
      client_id      = "basic-client"
      activation_key = "BASIC1234567890ABCDEF"
    }
  }
}