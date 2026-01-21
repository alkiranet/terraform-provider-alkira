resource "alkira_service_bluecat" "edge_only" {
  name                = "bluecat-edge-only"
  cxp                 = "EU-CENTRAL"
  description         = "Bluecat service with Edge instances only"
  global_cidr_list_id = alkira_list_global_cidr.branch_dns.id
  license_type        = "BRING_YOUR_OWN"
  segment_ids         = [alkira_segment.branch.id, alkira_segment.dmz.id]
  service_group_name  = "edge-dns-services"

  edge_anycast {
    ips         = ["172.16.50.10"]
    backup_cxps = ["US-WEST"]
  }

  instance {
    name = "edge-branch-01"
    type = "EDGE"
    
    edge_options {
      hostname    = "edge-branch-01.example.com"
      version     = "4.1.2"
      config_data = "CONFIG_DATA_STRING_BRANCH_01_ENCODED_BASE64"
    }
  }

  instance {
    name = "edge-branch-02"
    type = "EDGE"
    
    edge_options {
      hostname    = "edge-branch-02.example.com"
      version     = "4.1.2"
      config_data = "CONFIG_DATA_STRING_BRANCH_02_ENCODED_BASE64"
    }
  }

  instance {
    name = "edge-dmz"
    type = "EDGE"
    
    edge_options {
      hostname    = "edge-dmz.example.com"
      version     = "4.0.5"
      config_data = "CONFIG_DATA_STRING_DMZ_ENCODED_BASE64"
    }
  }
}