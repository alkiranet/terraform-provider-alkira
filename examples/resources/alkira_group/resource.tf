# Basic connector group
resource "alkira_group" "connector_group" {
  name        = "production-connectors"
  description = "Group of production connectors"
}
