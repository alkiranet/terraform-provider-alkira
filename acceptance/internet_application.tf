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

  # Internet applications require an internet exit in the same segment/CXP.
  # Without this, Terraform may delete the internet exit before the
  # internet application during destroy, causing an API rejection.
  depends_on = [alkira_connector_internet_exit.test1]
}
