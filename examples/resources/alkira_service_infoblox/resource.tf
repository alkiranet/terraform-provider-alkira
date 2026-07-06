resource "alkira_service_infoblox" "test" {
  cxp                 = "US-WEST-1"
  global_cidr_list_id = alkira_list_global_cidr.testcidr.id
  license_type        = "BRING_YOUR_OWN"
  name                = "alkiraServiceInfoblox5"
  segment_ids         = [alkira_segment.test1.id]
  service_group_name  = "serviceGroupName"
  shared_secret       = "thisisanewsecredet"
  size                = "SMALL" # drives NIOS-X instance sizing; NIOS instances ignore it

  # NIOS instance (grid-managed appliance): model + member-role type required.
  instance {
    anycast_enabled = false
    hostname        = "hostname.localdomain"
    model           = "TE-V1425"
    password        = "password1234"
    type            = "MASTER_CANDIDATE"
    version         = "8.5.2"
  }

  # NIOS-X instance (SaaS-managed): platform = NIOS_X, registered via a join token.
  # No model / member-role type. One service may mix NIOS and NIOS_X instances.
  instance {
    anycast_enabled = false
    hostname        = "niosx1.localdomain"
    platform        = "NIOS_X"
    version         = "4.0.1"
    join_token      = "REPLACE_WITH_JOIN_TOKEN"
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

