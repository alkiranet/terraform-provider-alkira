resource "alkira_service_infoblox" "test" {
  cxp                 = "US-WEST-1"
  global_cidr_list_id = alkira_list_global_cidr.testcidr.id
  license_type        = "BRING_YOUR_OWN"
  name                = "alkiraServiceInfoblox5"
  segment_ids         = [alkira_segment.test1.id]
  service_group_name  = "serviceGroupName"
  shared_secret       = "thisisanewsecredet"

  instances {
    anycast_enabled = false
    name            = "instance3"
    hostname        = "hostname.localdomain"
    model           = "TE-V1425"
    password        = "password1234"
    type            = "MASTER_CANDIDATE"
    version         = "8.5.2"
  }

  instances {
    anycast_enabled = false
    name            = "instance2"
    hostname        = "hostname.localdomain"
    model           = "TE-V1425"
    password        = "password1234"
    type            = "MASTER_CANDIDATE"
    version         = "8.5.2"
  }

  anycast {
    enabled = false
  }

  grid_master {
    external = false
    ip       = "10.10.10.10"
    name     = "newGridName2"
    username = "admin"
    password = "admin1234"
  }
}

