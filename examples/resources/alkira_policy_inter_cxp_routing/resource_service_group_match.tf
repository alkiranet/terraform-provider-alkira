#
# Match by service group. Routes learned from the BlueCat service are
# advertised with extra AS hops, and the local preference is reset to
# the default. Use service_group_implicit_group_id, not service_group_id.
#
resource "alkira_policy_inter_cxp_routing" "service_group_match" {
  name        = "bluecat-backup-path"
  description = "Make US-EAST-1 a backup path for BlueCat routes"
  enabled     = true
  direction   = "OUTBOUND"
  segment_id  = alkira_segment.seg1.id
  source_cxps = ["US-WEST-1"]
  dest_cxps   = ["US-EAST-1"]

  rule {
    name   = "prepend-bluecat"
    action = "ALLOW"
    match_service_group_ids = [
      alkira_service_bluecat.bluecat1.service_group_implicit_group_id,
    ]
    set_as_path_prepend             = "65100 65100 65100"
    set_reset_local_pref_to_default = true
  }
}
