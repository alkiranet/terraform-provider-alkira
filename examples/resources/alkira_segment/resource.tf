resource "alkira_segment" "test" {
  name  = "test-segment"
  asn   = "65513"
  cidrs = ["10.16.1.0/24", "10.1.1.0/24"]
}
