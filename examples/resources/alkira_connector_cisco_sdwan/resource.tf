resource "alkira_segment" "test" {
  name = "test"
  asn  = "65513"
  cidr = "10.1.1.0/24"
}

resource "alkira_credential_cisco_sdwan" "test" {
  name           = "test"
  username       = "xxxxx" # fill in proper username
  password       = "xxxxx" # fill in proper password
}

resource "alkira_connector_cisco_sdwan" "test" {
  name          = "test"
  cxp           = "US-WEST"
  size          = "SMALL"
  version       = "18.4.0"

  vedge {
    hostname        = "vedge1"
    cloud_init_file = "xxxxxxxxxxxxx"
    credential_id   = alkira_credential_cisco_sdwan.test.id
  }

  vrf_segment_mapping {
    segment_id = alkira_segment.test.id
    vrf_id     = 1 # fill in the proper vrf_id
  }
}

