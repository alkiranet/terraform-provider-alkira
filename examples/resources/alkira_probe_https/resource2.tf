resource "alkira_probe_https" "no_cert_validation" {
  name    = "https-no-cert-validation"
  enabled = true

  network_entity {
    type = "INTERNET_APPLICATION"
    id   = alkira_internet_application.example_application.id
  }

  uri                     = "www.alkira.net/api/dashboard"
  disable_cert_validation = true

  validators {
    type  = "RESPONSE_BODY"
    regex = ".*Dashboard Ready.*"
  }

  period_seconds  = 45
  timeout_seconds = 10
}
