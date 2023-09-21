resource "alkira_flow_collector" "test" {
  name        = "test"
  description = "test"
  enabled     = false

  cxps             = [var.cxp]
  destination_ip   = "172.16.0.1"
  destination_port = "2379"
}
