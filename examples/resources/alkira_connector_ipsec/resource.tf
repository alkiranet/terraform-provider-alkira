#
# Create a segment
#
resource "alkira_segment" "segment1" {
  name = "seg1"
  asn  = "65513"
  cidr = "10.16.1.0/24"
}

#
# Create a group
#
resource "alkira_group" "group1" {
  name        = "group1"
  description = "group1"
}

#
# Create IPSec connector and attach it with one segment and group.
#
resource "alkira_connector_ipsec" "connector-ipsec1" {
  name           = "connector-ipsec1"
  cxp            = "US-WEST"
  group          = alkira_group.group1.name
  segment_id     = alkira_segment.segment1.id
  size           = "SMALL"

  vpn_mode       = "ROUTE_BASED"

  routing_options {
    type = "DYNAMIC"
    customer_gateway_asn = "65310"
  }
}
