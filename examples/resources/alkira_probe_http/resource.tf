resource "alkira_probe_http" "basic_http" {
  name = "basic-http-probe"

  network_entity_id = alkira_internet_application.example_application.id

  uri = "www.alkira.net/api/health"
}
