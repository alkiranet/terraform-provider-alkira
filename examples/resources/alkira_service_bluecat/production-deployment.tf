resource "alkira_service_bluecat" "production" {
  name                = "bluecat-production"
  cxp                 = "ASIA-PACIFIC"
  description         = "Production Bluecat service for enterprise DNS"
  global_cidr_list_id = alkira_list_global_cidr.enterprise.id
  license_type        = "BRING_YOUR_OWN"
  segment_ids         = [alkira_segment.production.id]
  service_group_name  = "enterprise-dns"

  billing_tag_ids = [
    alkira_billing_tag.networking.id,
    alkira_billing_tag.production.id,
    alkira_billing_tag.asia_pacific.id
  ]

  bdds_anycast {
    ips = ["10.100.1.10"]
    backup_cxps = ["US-WEST"]
  }

  edge_anycast {
    ips = ["10.200.1.10"]
    backup_cxps = ["US-WEST"]
  }

  # Primary BDDS for enterprise services
  instance {
    name = "bdds-enterprise-primary"
    type = "BDDS"
    
    bdds_options {
      hostname       = "bdds-ent-01.asia.enterprise.com"
      model          = "cBDDS50"
      version        = "9.5.2"
      client_id      = "enterprise-asia-001"
      activation_key = "ENT_ASIA_PRIMARY_KEY_2024_001"
    }
  }

  # Secondary BDDS for redundancy
  instance {
    name = "bdds-enterprise-secondary"
    type = "BDDS"
    
    bdds_options {
      hostname       = "bdds-ent-02.asia.enterprise.com"
      model          = "cBDDS50"
      version        = "9.5.2"
      client_id      = "enterprise-asia-002"
      activation_key = "ENT_ASIA_SECONDARY_KEY_2024_002"
    }
  }

  # Edge for distributed locations
  instance {
    name = "edge-asia-primary"
    type = "EDGE"
    
    edge_options {
      hostname    = "edge-primary.asia.enterprise.com"
      version     = "4.2.1"
      config_data = "ASIA_PRIMARY_EDGE_CONFIG_BASE64"
    }
  }

  # Edge for backup services
  instance {
    name = "edge-asia-backup"
    type = "EDGE"
    
    edge_options {
      hostname    = "edge-backup.asia.enterprise.com"
      version     = "4.2.1"
      config_data = "ASIA_BACKUP_EDGE_CONFIG_BASE64"
    }
  }
}