resource "alkira_service_bluecat" "hybrid_deployment" {
  name                = "bluecat-hybrid"
  cxp                 = "US-EAST"
  description         = "Hybrid Bluecat deployment with both BDDS and Edge instances"
  global_cidr_list_id = alkira_list_global_cidr.global_dns.id
  license_type        = "BRING_YOUR_OWN"
  segment_ids         = [alkira_segment.production.id]
  service_group_name  = "hybrid-dns-services"

  billing_tag_ids = [
    alkira_billing_tag.dns_infrastructure.id,
    alkira_billing_tag.production.id
  ]

  bdds_anycast {
    ips         = ["192.168.10.50"]
    backup_cxps = ["US-WEST"]
  }

  edge_anycast {
    ips         = ["192.168.20.50"]
    backup_cxps = ["US-WEST"]
  }

  # Core BDDS instances for centralized management
  instance {
    name = "bdds-core-primary"
    type = "BDDS"
    
    bdds_options {
      hostname       = "bdds-core-01.enterprise.local"
      model          = "cBDDS50"
      version        = "9.5.1"
      client_id      = "enterprise-core-001"
      activation_key = "CORE1234ABCD5678EFGH9012IJKL"
    }
  }

  instance {
    name = "bdds-core-secondary"
    type = "BDDS"
    
    bdds_options {
      hostname       = "bdds-core-02.enterprise.local"
      model          = "cBDDS50"
      version        = "9.5.1"
      client_id      = "enterprise-core-002"
      activation_key = "CORE5678MNOP9012QRST3456UVWX"
    }
  }

  # Edge instances for distributed locations
  instance {
    name = "edge-datacenter-east"
    type = "EDGE"
    
    edge_options {
      hostname    = "edge-dc-east.enterprise.local"
      version     = "4.2.0"
      config_data = "EDGE_DC_EAST_CONFIG_BASE64_ENCODED_STRING"
    }
  }

  instance {
    name = "edge-datacenter-west"
    type = "EDGE"
    
    edge_options {
      hostname    = "edge-dc-west.enterprise.local"
      version     = "4.2.0"
      config_data = "EDGE_DC_WEST_CONFIG_BASE64_ENCODED_STRING"
    }
  }
}