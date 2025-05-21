resource "alkira_probe_tcp" "basic_tcp" {
  name              = "basic-tcp-probe"
  network_entity_id = alkira_internet_application.example_application.id
  port              = 80
}
