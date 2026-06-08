# A /30 reservation with an explicit prefix.
# For a /30, both first_ip_assignment and node_id are required.
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

# A /30 reservation without an explicit prefix — the backend assigns it based
# on prefix_type and prefix_len. A /30 still requires first_ip_assignment and
# node_id even when the backend assigns the prefix.
resource "alkira_ip_reservation" "backend_assigned_prefix" {
  name                = "test-backend"
  type                = "OVERLAY"
  prefix_type         = "APIPA"
  prefix_len          = 30
  first_ip_assignment = "CUSTOMER"
  node_id             = "d70503d2-1a99-4084-8aae-8268e2764365"
  scale_group_id      = "99a6f3db-02d5-4189-8b0a-352eaeda2e10"
  segment_id          = alkira_segment.test.id
  cxp                 = "US-WEST"
}

# A single-IP (/32) reservation. first_ip_assignment only governs which IP of a
# /30 pair is assigned, so it is omitted here.
resource "alkira_ip_reservation" "single_ip" {
  name           = "test-single-ip"
  type           = "OVERLAY"
  prefix         = "10.1.0.10/32"
  prefix_type    = "SEGMENT"
  scale_group_id = "99a6f3db-02d5-4189-8b0a-352eaeda2e10"
  segment_id     = alkira_segment.test.id
  cxp            = "US-WEST"
}
