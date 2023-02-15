# Create resources under Management/Segments

resource "alkira_segment" "test1" {
  name  = "acceptance-seg-test1"
  asn   = "65514"
  cidrs = ["10.16.1.0/24"]
}

resource "alkira_segment" "test2" {
  name  = "acceptance-seg-test2"
  asn   = "65514"
  cidrs = ["10.16.1.0/24"]
}

resource "alkira_segment" "test3" {
  name  = "acceptance-seg-test3"
  asn   = "65513"
  cidrs = ["10.1.1.0/24"]
}

resource "alkira_segment" "test4" {
  name        = "acceptance-seg-test4"
  description = "test segment 4"
  asn         = "65513"
  cidrs       = ["10.2.1.0/24"]
}
