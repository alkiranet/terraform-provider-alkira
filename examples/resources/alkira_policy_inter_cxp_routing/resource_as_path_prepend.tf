#
# AS-path prepend for traffic engineering. Routes matching the prefix
# list are advertised with three extra AS hops so the destination CXP
# prefers alternate paths.
#
resource "alkira_policy_inter_cxp_routing" "as_path_prepend" {
  name        = "us-east-backup-path"
  description = "Make US-EAST-1 a backup path by prepending AS hops"
  enabled     = true
  direction   = "OUTBOUND"
  segment_id  = alkira_segment.seg1.id
  source_cxps = ["US-WEST-1"]
  dest_cxps   = ["US-EAST-1"]

  rule {
    name                  = "prepend-for-backup"
    action                = "ALLOW"
    match_prefix_list_ids = [alkira_policy_prefix_list.corp_prefixes.id]
    set_as_path_prepend   = "65100 65100 65100"
  }
}
