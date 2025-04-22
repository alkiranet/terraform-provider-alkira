resource "alkira_probe_https" "custom_cert_https" {
  name    = "https-custom-cert"
  enabled = true

  network_entity {
    type = "INTERNET_APPLICATION"
    id   = alkira_internet_application.example_application.id
  }

  uri         = "/secure/endpoint"
  server_name = "api.example.com"

  ca_certificate = file("${path.module}/certs/exmaple_ca.pem")

  headers = {
    "Authorization" = "Basic dXNlcjpwYXNzd29yZA=="
  }

  validators {
    type        = "STATUS_CODE"
    status_code = "200"
  }

  failure_threshold = 2
  success_threshold = 1
  period_seconds    = 60
  timeout_seconds   = 15
}
