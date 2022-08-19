resource "alkira_segment" "test1" {
  name  = "tftest-basic"
  asn   = "65514"
  cidrs = ["10.16.1.0/24"]
}

resource "alkira_segment" "test2" {
  name                = "tftest-basic-public-ip"
  asn                 = "65514"
  cidrs               = ["10.16.1.0/24"]
  reserve_public_ips = true
}


resource "alkira_segment" "seg1" {
  name  = "tftest-segment1"
  asn   = "65513"
  cidrs = ["10.1.1.0/24"]
}

resource "alkira_segment" "seg2" {
  name  = "tftest-segment2"
  asn   = "65513"
  cidrs = ["10.2.1.0/24"]
}
