resource "alkira_policy_prefix_list" "test" {
  name        = "test-list"
  description = "Prefix list with detailed descriptions"

  prefix {
    cidr        = "10.2.3.0/24"
    description = "Production subnet for US-East-1"
  }
  prefix {
    cidr        = "10.4.6.0/24"
    description = "Staging environment subnet"
  }

}

