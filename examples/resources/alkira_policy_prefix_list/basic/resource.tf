resource "alkira_policy_prefix_list" "example" {
  name        = "example-prefix-list"
  description = "Basic example prefix list"

  prefix {
    cidr = "10.0.0.0/24"
  }
  prefix {
    cidr = "192.168.1.0/24"
  }

}

