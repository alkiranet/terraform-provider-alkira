resource "alkira_internet_application" "test" {
  name              = "acceptance-ifa"
  connector_id      = alkira_connector_ipsec.test.id
  connector_type    = "IP_SEC"
  fqdn_prefix       = "acceptance"
  internet_protocol = "IPV4"
  segment_id        = alkira_segment.test1.id
  size              = "SMALL"

  target {
    type        = "IP"
    value       = "192.168.1.1"
    port_ranges = [-1]
  }
}
