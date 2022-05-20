resource "alkira_service_infoblox" "test" {
  cxp                 = "US-WEST-1"
  global_cidr_list_id = alkira_list_global_cidr.testcidr.id
  license_type        = "BRING_YOUR_OWN"
  name                = "alkiraServiceInfoblox5"
  segment_names       = [alkira_segment.test1.name]
  service_group_name  = "serviceGroupName"
  shared_secret       = "thisisanewsecredetshhhhh"
  size                = "SMALL"

  instances {
    any_cast_enabled = false
    name             = "instance2"
    host_name        = "host_name.localdomain"
    model            = "TE-V1425"
    password         = "password1234"
    type             = "MASTER_CANDIDATE"
    version          = "8.5.2"
  }

  instances {
    any_cast_enabled = false
    name             = "instance2"
    host_name        = "host_name.localdomain"
    model            = "TE-V1425"
    password         = "password1234"
    type             = "MASTER_CANDIDATE"
    version          = "8.5.2"
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
