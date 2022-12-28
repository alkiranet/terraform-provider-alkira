resource "alkira_service_infoblox" "test" {
  name                = "alkiraServiceInfoblox5"
  cxp                 = "US-WEST-1"
  global_cidr_list_id = alkira_list_global_cidr.testcidr.id
  segment_ids         = [alkira_segment.test1.id]
  service_group_name  = "serviceGroupName"

  grid_master {
    # existing = false means create a new grid master.
    existing      = false
    name          = "newGridName2"
    username      = "admin"
    password      = "Abcd12345"
    shared_secret = "thisisanewsecredetshhhhh"
  }

  # Only one instance is allowed when creating a new grid master.
  instance {
    anycast_enabled = true
    hostname        = "hostname.localdomain"
    model           = "TE-V1425"
    password        = "Abcd12345"
    type            = "MASTER"
    version         = "8.5.2"
  }

  anycast {
    enabled = false
  }
}