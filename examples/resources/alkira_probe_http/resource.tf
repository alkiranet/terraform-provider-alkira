resource "alkira_probe_http" "basic_http" {
  name = "basic-http-probe"

  network_entity {
    type = "INTERNET_APPLICATION"
    id   = alkira_internet_application.example_application.id
  }

  uri = "/health"
}
