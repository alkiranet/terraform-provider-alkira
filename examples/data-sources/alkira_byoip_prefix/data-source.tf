data "alkira_byoip_prefix" "byoip" {
  name = "byoip-pfx-01"
}

resource "alkira_connector_internet_exit" "cn" {
  name       = "inet-exit"
  cxp        = "US-WEST"
  group      = "group1"
  segment_id = alkira_segment.test.id
  byoip_id   = data.alkira_byoip.byoip.id
}

