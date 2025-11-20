resource "alkira_service_bluecat" "test" {
  name                 = "test-bluecat"
  cxp                  = "US-WEST"
  description          = "Test Bluecat service"
  global_cidr_list_id  = alkira_list_global_cidr.test.id
  license_type         = "BRING_YOUR_OWN"
  segment_ids          = [alkira_segment.test.id]
  service_group_name   = "bluecat-group"

  bdds_anycast {
    ips = ["10.0.1.100"]
    backup_cxps = ["US-EAST"]
  }

  edge_anycast {
    ips = ["10.0.1.101"]
    backup_cxps = ["US-EAST"]
  }

  instance {
    name = "bdds-instance"
    type = "BDDS"
    
    bdds_options {
      hostname = "bdds.example.com"
      model    = "BDDS-50"
      version  = "9.4.0"
      client_id = "0924124nfds3"
      activation_key = "asfasfvdgregterg"
    }
  }

  instance {
    name = "edge-instance"
    type = "EDGE"
    
    edge_options {
      hostname = "edge.example.com"
      version  = "4.0.0"
      config_data = "ASDSFDDFGDGHRTGHRTHBDS234325SDFVVD"
    }
  }
}