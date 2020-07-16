provider "alkira" {
  portal   = "spike1.returntrue.dev.alkira2.net"
  username = "sanchit@alkira.net"
  password = "Alkira2018"
}

resource "alkira_segment" "my_segment" {
  name = "seg3"
  asn  = "65513"
  cidr = "10.16.1.0/24"
}
