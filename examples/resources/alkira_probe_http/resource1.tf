resource "alkira_probe_http" "advanced_http" {
  name              = "http-with-validators"
  enabled           = true
  network_entity_id = alkira_internet_application.example_application.id
  uri               = "www.alkira.net/api/health"

  validators {
    type        = "STATUS_CODE"
    status_code = "200-299"
  }

  validators {
    type  = "RESPONSE_BODY"
    regex = ".*status.*:.*OK.*"
  }

  failure_threshold = 3
  success_threshold = 2
  period_seconds    = 30
  timeout_seconds   = 10
}
