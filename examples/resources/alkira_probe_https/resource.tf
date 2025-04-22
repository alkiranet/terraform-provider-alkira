resource "alkira_probe_https" "basic_https" {
  name = "basic-https-probe"

  network_entity {
    type = "INTERNET_APPLICATION"
    id   = alkira_internet_application.example_application.id
  }

  uri = "/api/status"
}
