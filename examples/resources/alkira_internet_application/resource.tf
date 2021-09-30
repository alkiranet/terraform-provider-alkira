#
# A simple internet facing application that assumes that
# connector_aws_vpc and segment was already created.
#
resource "alkira_internet_application" "app1" {
  name           = "app1"

  connector_id   = alkira_connector_aws_vpc.vpc1.id
  connector_type = "AWS_VPC"

  fqdn_prefix    = "tfexample"
  segment        = alkira_segment.seg1.name

  private_ip     = "10.0.0.1"
  private_port   = "80"
  size           = "SMALL"
}
