# Multiple groups for different environments
resource "alkira_group" "dev_connectors" {
  name        = "development-connectors"
  description = "Development environment connectors"
}

resource "alkira_group" "staging_connectors" {
  name        = "staging-connectors"
  description = "Staging environment connectors"
}

resource "alkira_group" "prod_connectors" {
  name        = "production-connectors"
  description = "Production environment connectors"
}

# Output the group IDs for reference in other configurations
output "dev_group_id" {
  description = "ID of the development group"
  value       = alkira_group.dev_connectors.id
}

output "staging_group_id" {
  description = "ID of the staging group"
  value       = alkira_group.staging_connectors.id
}

output "prod_group_id" {
  description = "ID of the production group"
  value       = alkira_group.prod_connectors.id
}
