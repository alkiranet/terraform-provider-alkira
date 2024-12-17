resource "alkira_service_f5_vserver_endpoint" "example-vserver-1" {
  name                    = "example1"
  f5_service_id           = alkira_service_f5_lb.example-lb.id
  f5_service_instance_ids = [alkira_service_f5_lb.example-lb.instances[0].id, alkira_service_f5_lb.example-lb.instances[1].id]
  type                    = "ELB"
  segment_id              = alkira_segment.example_segment.id
  fqdn_prefix             = "example-prefix"
  protocol                = "TCP"
  port_ranges             = ["8000-8010", "443", "80"]
  snat                    = "AUTOMAP"
}
