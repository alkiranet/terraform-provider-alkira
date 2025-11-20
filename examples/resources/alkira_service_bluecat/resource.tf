resource "alkira_service_bluecat" "test" {
  name                 = "test-bluecat"
  cxp                  = "US-WEST"
  description          = "Test Bluecat service"
  global_cidr_list_id  = alkira_list_global_cidr.test.id
  license_type         = "BRING_YOUR_OWN"
  segment_ids          = [alkira_segment.test.id]
  service_group_name   = "bluecat-group"

  bddsAnycast {
    ips = ["10.0.1.100"]
    backup_cxps = ["US-EAST"]
  }

  edgeAnycast {
    ips = ["10.0.1.101"]
    backup_cxps = ["US-EAST"]
  }

  instance {
    name = "bdds-instance"
    type = "BDDS"
    
    bddsOptions {
      hostname = "bdds.example.com"
      model    = "BDDS-50"
      version  = "9.4.0"
    }
  }

  instance {
    name = "edge-instance"
    type = "EDGE"
    
    edgeOptions {
      hostname = "edge.example.com"
      version  = "4.0.0"
    }
  }
}