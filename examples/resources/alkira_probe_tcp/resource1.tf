resource "alkira_probe_tcp" "full_tcp" {
  name              = "tcp-full-options"
  enabled           = true
  network_entity_id = alkira_internet_application.example_application.id
  port              = 443
  failure_threshold = 5
  success_threshold = 3
  period_seconds    = 20
  timeout_seconds   = 5
}
