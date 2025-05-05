resource "alkira_probe_https" "custom_cert_https" {
  name    = "https-custom-cert"
  enabled = true

  network_entity {
    type = "INTERNET_APPLICATION"
    id   = alkira_internet_application.example_application.id
  }

  uri         = "www.alkira.net/secure/endpoint"
  server_name = "api.example.com"

  # we can either pass the path of the file 
  # ca_certificate = file("${path.module}/certs/exmaple_ca.pem")

  # or the whole certificate as a string.
  ca_certificate = "-----BEGIN CERTIFICATE-----ZXhhbXBsZS1jYS1jZXJ0aWZpY2F0ZS4gSGVsbG8gY3VyaW91cyBwZXJzb24=-----END CERTIFICATE-----"

  validators {
    type        = "STATUS_CODE"
    status_code = "200"
  }

  failure_threshold = 2
  success_threshold = 1
  period_seconds    = 60
  timeout_seconds   = 15
}
