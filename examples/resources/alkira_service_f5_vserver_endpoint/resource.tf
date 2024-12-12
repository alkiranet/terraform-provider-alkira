resource "alkira_service_f5_vserver_endpoint" "example_vserver" {
  name                    = "example"
  f5_service_id           = alkira_service_f5_lb.example_lb.id
  f5_service_instance_ids = [alkira_service_f5_lb.example_lb.instances[0].id]
  type                    = "ELB"
  segment_id              = alkira_segment.example_segment.id
  fqdn_prefix             = "example_prefix"
  protocol                = "UDP"
  port_ranges             = ["-1"]
  snat                    = "AUTOMAP"
}
