resource "alkira_service_checkpoint" "test" {
  name       = "acceptance-checkpoint"
  auto_scale = "ON"
  cxp        = var.cxp

  license_type = "PAY_AS_YOU_GO"
  size         = "SMALL"
  version      = "R80.30"
  pdp_ips      = ["10.0.0.1"]
  password     = "abcd1234"

  max_instance_count = 2
  min_instance_count = 2

  segment_id = alkira_segment.test1.id

  management_server {
    type                = "MDS"
    configuration_mode  = "AUTOMATED"
    reachability        = "PRIVATE"
    ips                 = ["192.168.3.3"]
    global_cidr_list_id = alkira_list_global_cidr.checkpoint.id
    segment_id          = alkira_segment.test1.id

    username = "checkpoint-user"
    password = "abcd1234"

    # domain only required when configuration_mode is AUTOMATED and
    # when type is MDS.
    domain = "test.alkira.com"
  }

  instance {
    name    = "ins1"
    sic_key = "abcd1234"
  }

  instance {
    name    = "ins2"
    sic_key = "abcd12345"
  }
}

resource "alkira_service_checkpoint" "test2" {
  name       = "acceptance-checkpoint-2"
  auto_scale = "OFF"
  cxp        = var.cxp

  billing_tag_ids    = [alkira_billing_tag.test1.id]
  license_type       = "BRING_YOUR_OWN"
  max_instance_count = 1
  min_instance_count = 1

  password        = "xxxxxxxx"
  pdp_ips         = ["10.1.1.116"]
  segment_id      = alkira_segment.test1.id
  size            = "LARGE"
  tunnel_protocol = "IPSEC"
  version         = "R81"

  instance {
    name    = "acceptance-checkpoint-2"
    sic_key = "ak12345678"
  }

  management_server {
    configuration_mode  = "AUTOMATED"
    global_cidr_list_id = alkira_list_global_cidr.checkpoint.id
    ips                 = ["54.69.129.30"]
    username            = "checkpoint-user"
    password            = "Alkira2023"
    reachability        = "PUBLIC"
    type                = "SMS"
  }

  segment_options {
    groups     = [alkira_group.test1.name]
    segment_id = alkira_segment.test1.id
    zone_name  = "DEFAULT"
  }
}
