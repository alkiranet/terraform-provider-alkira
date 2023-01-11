resource "alkira_internet_application" "test" {
  name           = "test-ifa"
  connector_id   = alkira_connector_aws_vpc.test.id
  connector_type = "AWS_VPC"
  fqdn_prefix    = "tfexample"
  segment_id     = alkira_segment.seg1.id
  size           = "SMALL"

  target {
    type = "IP"
    value = "192.168.1.1"
    ports = [1200]
  }
}
