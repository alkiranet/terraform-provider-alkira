resource "alkira_internet_application" "test" {
  name           = "test-ifa"
  connector_id   = alkira_connector_aws_vpc.test.id
  connector_type = "AWS_VPC"
  fqdn_prefix    = "tfexample"
  segment_id     = alkira_segment.seg1.id
  size           = "SMALL"

  target {
    type        = "IP"
    value       = "192.168.1.1"
    port_ranges = ["1200"]
  }
}

resource "alkira_internet_application" "test_internal_dns" {
  name           = "test-ifa-internal-dns"
  connector_id   = alkira_connector_aws_vpc.test2.id
  connector_type = "AWS_VPC"
  fqdn_prefix    = "tfexample-dns"
  segment_id     = alkira_segment.seg1.id
  size           = "SMALL"

  target {
    type                = "INTERNAL_DNS"
    policy_fqdn_list_id = alkira_list_policy_fqdn.test.id
    port_ranges         = ["1200"]
  }
}
