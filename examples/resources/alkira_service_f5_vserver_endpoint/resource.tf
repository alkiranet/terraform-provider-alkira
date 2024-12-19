resource "alkira_service_f5_vserver_endpoint" "example-vserver" {
  name                    = "example"
  f5_service_id           = alkira_service_f5_lb.example-lb.id
  f5_service_instance_ids = [alkira_service_f5_lb.example-lb.instances[0].id]
  type                    = "ELB"
  segment_id              = alkira_segment.example-segment.id
  fqdn_prefix             = "example-prefix"
  protocol                = "UDP"
  port_ranges             = ["8000"]
  snat                    = "AUTOMAP"
}
