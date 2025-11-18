resource "alkira_service_f5_vserver_endpoint" "example-vserver-2" {
  name                              = "example2"
  f5_service_id                     = alkira_service_f5_lb.example-ilb.id
  f5_service_instance_ids           = [alkira_service_f5_lb.example-ilb.instances[0].id, alkira_service_f5_lb.example-ilb.instances[1].id]
  type                              = "ILB"
  segment_id                        = alkira_segment.example_segment.id
  protocol                          = "TCP"
  destination_endpoint_port_ranges  = ["8000-8010", "443", "80"]
  destination_endpoint_ip_addresses = ["1.2.3.4", "1.2.4.4"]
  snat                              = "NONE"
}
