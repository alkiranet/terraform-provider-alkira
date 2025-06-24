resource "alkira_probe_https" "basic_https" {
  name = "basic-https-probe"

  network_entity_id = alkira_internet_application.example_application.id

  uri = "www.alkira.net/api/status"
}
