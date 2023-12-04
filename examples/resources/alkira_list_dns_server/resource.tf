resource "alkira_list_dns_server" "test" {
  name            = "test"
  description     = "test dns server list"
  dns_server_ips  = ["8.8.8.8"]
  segment_id      = alkira_segment.test.id
}
