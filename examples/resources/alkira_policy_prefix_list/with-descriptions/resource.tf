resource "alkira_policy_prefix_list" "descriptive" {
  name        = "descriptive-prefixes"
  description = "Prefix list with detailed descriptions"

  prefixes = [
    {
      prefix       = "10.2.3.0/24"
      description = "Production subnet for US-East-1"
    },
    {
      prefix       = "10.4.6.0/24" 
      description = "Staging environment subnet"
    }
  ]
}