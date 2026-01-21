resource "alkira_service_bluecat" "bdds_only" {
  name                = "bluecat-bdds-only"
  cxp                 = "US-WEST"
  description         = "Bluecat service with BDDS instances only"
  global_cidr_list_id = alkira_list_global_cidr.dns_allowed.id
  license_type        = "BRING_YOUR_OWN"
  segment_ids         = [alkira_segment.corp.id]
  service_group_name  = "dns-services"

  bdds_anycast {
    ips         = ["10.0.100.10"]
    backup_cxps = ["US-EAST"]
  }

  instance {
    name = "bdds-primary"
    type = "BDDS"
    
    bdds_options {
      hostname       = "bdds-primary.corp.local"
      model          = "cBDDS50"
      version        = "9.4.0"
      client_id      = "bdds-client-001"
      activation_key = "ABCD1234EFGH5678IJKL9012"
    }
  }

  instance {
    name = "bdds-secondary"
    type = "BDDS"
    
    bdds_options {
      hostname       = "bdds-secondary.corp.local"
      model          = "cBDDS50"
      version        = "9.4.0"
      client_id      = "bdds-client-002"
      activation_key = "MNOP3456QRST7890UVWX1234"
    }
  }
}