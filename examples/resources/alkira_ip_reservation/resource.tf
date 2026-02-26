# With an explicit prefix
resource "alkira_ip_reservation" "explicit_prefix" {
  name                = "test"
  type                = "OVERLAY"
  prefix              = "169.254.200.0/30"
  prefix_type         = "APIPA"
  first_ip_assignment = "CUSTOMER"
  node_id             = "d70503d2-1a99-4084-8aae-8268e2764365"
  scale_group_id      = "99a6f3db-02d5-4189-8b0a-352eaeda2e10"
  segment_id          = alkira_segment.test.id
  cxp                 = "US-WEST"
}

# Without a prefix — the backend assigns it based on prefix_type and prefix_len
resource "alkira_ip_reservation" "backend_assigned_prefix" {
  name           = "test-backend"
  type           = "OVERLAY"
  prefix_type    = "APIPA"
  prefix_len     = 30
  scale_group_id = "99a6f3db-02d5-4189-8b0a-352eaeda2e10"
  segment_id     = alkira_segment.test.id
  cxp            = "US-WEST"
}
