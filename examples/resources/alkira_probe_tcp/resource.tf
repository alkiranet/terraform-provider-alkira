resource "alkira_probe_tcp" "basic_tcp" {
  name = "basic-tcp-probe"

  network_entity {
    type = "INTERNET_APPLICATION"
    id   = alkira_internet_application.example_application.id
  }

  port = 80
}
