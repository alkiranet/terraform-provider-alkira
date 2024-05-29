resource "alkira_segment" "example" {
  name  = "example"
  asn   = "65513"
  cidrs = ["10.16.1.0/24"]
}

resource "alkira_policy_nat_rule" "example" {
  name          = "example"
  description   = "example nat rule"
  enabled       = false

  match {
    src_prefixes = ["any"]
    dst_prefixes = ["any"]
    protocol     = "any"
  }

  action {
    src_addr_translation_type = "NONE"
    dst_addr_translation_type = "NONE"
  }
}

resource "alkira_policy_nat" "example" {
  name               = "example"
  description        = "terraform example NAT policy"
  type               = "INTRA_SEGMENT"
  segment_id         = alkira_segment.example.id
  included_group_ids = [alkira_group.example.id]
  nat_rule_ids       = [alkira_policy_nat_rule.example.id]
}
