resource "alkira_connector_cisco_sdwan" "test" {
  name          = "test"
  cxp           = "US-WEST"
  size          = "SMALL"
  version       = "18.4.0"

  vedge {
    hostname        = "vedge1"
    cloud_init_file = "xxxxxxxxxxxxx"
    username        = "username"
    password        = "password"
  }

  vrf_segment_mapping {
    segment_id = alkira_segment.test.id
    vrf_id     = 1 # fill in the proper vrf_id
  }
}

