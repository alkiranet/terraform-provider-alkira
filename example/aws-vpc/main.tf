provider "alkira" {
  portal   = "spike1.returntrue.dev.alkira2.net"
  username = "sanchit@alkira.net"
  password = "Alkira2018"
}

resource "alkira_segment" "segment2" {
  name = "seg2"
  asn  = "65513"
  cidr = "10.16.1.0/24"
}

resource "alkira_segment" "segment3" {
  name = "seg3"
  asn  = "65513"
  cidr = "10.16.1.0/24"
}

resource "alkira_connector_aws_vpc" "connector1" {
  vpc_1_id       =
  vpc_1_name     =
  vpc_1_owner_id =

  vpc_2_id       =
  vpc_2_name     =
  vpc_2_owner_id =

  size           = "SMALL"
  segments       = ["seg2"]
}
