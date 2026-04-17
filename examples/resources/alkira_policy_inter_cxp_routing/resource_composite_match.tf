#
# Composite rule — combine multiple match criteria (prefix + community +
# AS path) on a single rule and apply several set actions at once. All
# match_* fields on the same rule are AND-ed together.
#
resource "alkira_policy_inter_cxp_routing" "composite_match" {
  name        = "composite-match-and-set"
  description = "AND-match prefix, community, and AS path; apply multiple sets"
  enabled     = true
  direction   = "OUTBOUND"
  segment_id  = alkira_segment.seg1.id
  source_cxps = ["US-WEST-1"]
  dest_cxps   = ["US-EAST-1"]

  rule {
    name                     = "match-all-three-set-all-three"
    action                   = "ALLOW"
    match_prefix_list_ids    = [alkira_policy_prefix_list.corp_prefixes.id]
    match_community_list_ids = [alkira_list_community.prod_inbound.id]
    match_as_path_list_ids   = [alkira_list_as_path.trusted_as.id]
    set_as_path_prepend      = "65100 65100"
    set_community            = "65512:50"
    set_extended_community   = "rt:65512:50"
  }
}
